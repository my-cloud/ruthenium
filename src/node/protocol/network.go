package protocol

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/api/node"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const (
	neighborSynchronizationTimeInSeconds = 10
	maxOutboundsCount                    = 8
)

type Network struct {
	ip        string
	port      uint16
	logger    *log.Logger
	waitGroup *sync.WaitGroup

	timeable               clock.Timeable
	neighbors              []node.Requestable
	neighborsMutex         sync.RWMutex
	neighborsByTarget      map[string]node.Requestable
	neighborsByTargetMutex sync.RWMutex
	seedsByTarget          map[string]node.Requestable
}

func NewNetwork(ip string, port uint16, timeable clock.Timeable, configurationPath string, logger *log.Logger) *Network {
	network := new(Network)
	network.ip = ip
	network.port = port
	network.timeable = timeable
	network.logger = logger
	var waitGroup sync.WaitGroup
	network.waitGroup = &waitGroup
	seedsIps := readSeedsIps(configurationPath, logger)
	network.seedsByTarget = map[string]node.Requestable{}
	for _, seedIp := range seedsIps {
		seed := neighborhood.NewNeighbor(seedIp, DefaultPort, logger)
		network.seedsByTarget[seed.Target()] = seed
	}
	network.neighborsByTarget = map[string]node.Requestable{}
	return network
}

func readSeedsIps(configurationPath string, logger *log.Logger) []string {
	jsonFile, err := os.Open(configurationPath + "/seeds.json")
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to open seeds IPs configuration file: %w", err).Error())
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	if err = jsonFile.Close(); err != nil {
		logger.Error(fmt.Errorf("unable to close seeds IPs configuration file: %w", err).Error())
	}
	var seedsIps []string
	if err = json.Unmarshal(byteValue, &seedsIps); err != nil {
		logger.Fatal(fmt.Errorf("unable to unmarshal seeds IPs: %w", err).Error())
	}
	return seedsIps
}

func (network *Network) Wait() {
	network.waitGroup.Wait()
}

func (network *Network) StartNeighborsSynchronization() {
	network.SynchronizeNeighbors()
	_ = time.AfterFunc(time.Second*neighborSynchronizationTimeInSeconds, network.StartNeighborsSynchronization)
}

func (network *Network) SynchronizeNeighbors() {
	network.neighborsByTargetMutex.Lock()
	var neighborsByTarget map[string]node.Requestable
	if len(network.neighborsByTarget) == 0 {
		neighborsByTarget = network.seedsByTarget
	} else {
		neighborsByTarget = network.neighborsByTarget
	}
	network.neighborsByTarget = map[string]node.Requestable{}
	network.neighborsByTargetMutex.Unlock()
	network.waitGroup.Add(1)
	go func(neighborsByTarget map[string]node.Requestable) {
		defer network.waitGroup.Done()
		var neighbors []node.Requestable
		var targetRequests []node.TargetRequest
		hostTargetRequest := node.TargetRequest{
			Ip:   &network.ip,
			Port: &network.port,
		}
		targetRequests = append(targetRequests, hostTargetRequest)
		network.neighborsMutex.RLock()
		for _, neighbor := range neighborsByTarget {
			neighborIp := neighbor.Ip()
			neighborPort := neighbor.Port()
			lookedUpNeighborsIps, err := net.LookupIP(neighborIp)
			if err != nil {
				network.logger.Error(fmt.Errorf("DNS discovery failed on addresse %s: %w", neighborIp, err).Error())
				return
			}

			neighborsCount := len(lookedUpNeighborsIps)
			if neighborsCount != 1 {
				network.logger.Error(fmt.Sprintf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", neighborsCount, neighborIp))
				return
			}
			lookedUpNeighborIp := lookedUpNeighborsIps[0]
			lookedUpNeighborIpString := lookedUpNeighborIp.String()
			if (lookedUpNeighborIpString != network.ip || neighborPort != network.port) && lookedUpNeighborIpString == neighborIp && neighbor.IsFound() {
				neighbors = append(neighbors, neighbor)
				targetRequest := node.TargetRequest{
					Ip:   &neighborIp,
					Port: &neighborPort,
				}
				targetRequests = append(targetRequests, targetRequest)
			}
		}
		network.neighborsMutex.RUnlock()
		rand.Seed(network.timeable.Now().UnixNano())
		rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
		outboundsCount := int(math.Min(float64(len(neighbors)), maxOutboundsCount))
		network.neighborsMutex.Lock()
		network.neighbors = neighbors[:outboundsCount]
		network.neighborsMutex.Unlock()
		for _, neighbor := range neighbors[:outboundsCount] {
			var neighborTargetRequests []node.TargetRequest
			for _, targetRequest := range targetRequests {
				neighborIp := neighbor.Ip()
				neighborPort := neighbor.Port()
				if neighborIp != *targetRequest.Ip || neighborPort != *targetRequest.Port {
					neighborTargetRequests = append(neighborTargetRequests, targetRequest)
				}
			}
			go func(neighbor node.Requestable) {
				_ = neighbor.SendTargets(neighborTargetRequests)
			}(neighbor)
		}
	}(neighborsByTarget)
}

func (network *Network) AddTargets(targetRequests []node.TargetRequest) {
	network.waitGroup.Add(1)
	go func() {
		defer network.waitGroup.Done()
		network.neighborsByTargetMutex.Lock()
		defer network.neighborsByTargetMutex.Unlock()
		for _, targetRequest := range targetRequests {
			neighbor := neighborhood.NewNeighbor(*targetRequest.Ip, *targetRequest.Port, network.logger)
			if _, ok := network.neighborsByTarget[neighbor.Target()]; !ok {
				network.neighborsByTarget[neighbor.Target()] = neighbor
			}
		}
	}()
}

func (network *Network) Neighbors() []node.Requestable {
	return network.neighbors
}
