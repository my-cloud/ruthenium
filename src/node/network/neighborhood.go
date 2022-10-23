package network

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/api/node/network"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

const (
	DefaultPort                  = 8106
	synchronizationTimeInSeconds = 10
	maxOutboundsCount            = 8
)

type Neighborhood struct {
	hostIp    string
	hostPort  uint16
	logger    *log.Logger
	waitGroup *sync.WaitGroup

	timeable               clock.Timeable
	neighbors              []network.Requestable
	neighborsMutex         sync.RWMutex
	neighborsByTarget      map[string]network.Requestable
	neighborsByTargetMutex sync.RWMutex
	seedsByTarget          map[string]network.Requestable
}

func NewNeighborhood(hostIp string, hostPort uint16, timeable clock.Timeable, configurationPath string, logger *log.Logger) *Neighborhood {
	neighborhood := new(Neighborhood)
	neighborhood.hostIp = hostIp
	neighborhood.hostPort = hostPort
	neighborhood.timeable = timeable
	neighborhood.logger = logger
	var waitGroup sync.WaitGroup
	neighborhood.waitGroup = &waitGroup
	seedsIps := readSeedsIps(configurationPath, logger)
	neighborhood.seedsByTarget = map[string]network.Requestable{}
	for _, seedIp := range seedsIps {
		seed := NewNeighbor(seedIp, DefaultPort, logger)
		neighborhood.seedsByTarget[seed.Target()] = seed
	}
	neighborhood.neighborsByTarget = map[string]network.Requestable{}
	return neighborhood
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

func (neighborhood *Neighborhood) Wait() {
	neighborhood.waitGroup.Wait()
}

func (neighborhood *Neighborhood) StartSynchronization() {
	neighborhood.Synchronize()
	_ = time.AfterFunc(time.Second*synchronizationTimeInSeconds, neighborhood.StartSynchronization)
}

func (neighborhood *Neighborhood) Synchronize() {
	neighborhood.neighborsByTargetMutex.Lock()
	var neighborsByTarget map[string]network.Requestable
	if len(neighborhood.neighborsByTarget) == 0 {
		neighborsByTarget = neighborhood.seedsByTarget
	} else {
		neighborsByTarget = neighborhood.neighborsByTarget
	}
	neighborhood.neighborsByTarget = map[string]network.Requestable{}
	neighborhood.neighborsByTargetMutex.Unlock()
	neighborhood.waitGroup.Add(1)
	go func(neighborsByTarget map[string]network.Requestable) {
		defer neighborhood.waitGroup.Done()
		var neighbors []network.Requestable
		var targetRequests []network.TargetRequest
		hostTargetRequest := network.TargetRequest{
			Ip:   &neighborhood.hostIp,
			Port: &neighborhood.hostPort,
		}
		targetRequests = append(targetRequests, hostTargetRequest)
		neighborhood.neighborsMutex.RLock()
		for _, neighbor := range neighborsByTarget {
			neighborIp := neighbor.Ip()
			neighborPort := neighbor.Port()
			lookedUpNeighborsIps, err := net.LookupIP(neighborIp)
			if err != nil {
				neighborhood.logger.Error(fmt.Errorf("DNS discovery failed on addresse %s: %w", neighborIp, err).Error())
				return
			}

			neighborsCount := len(lookedUpNeighborsIps)
			if neighborsCount != 1 {
				neighborhood.logger.Error(fmt.Sprintf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", neighborsCount, neighborIp))
				return
			}
			lookedUpNeighborIp := lookedUpNeighborsIps[0]
			lookedUpNeighborIpString := lookedUpNeighborIp.String()
			if (lookedUpNeighborIpString != neighborhood.hostIp || neighborPort != neighborhood.hostPort) && lookedUpNeighborIpString == neighborIp && neighbor.IsFound() {
				neighbors = append(neighbors, neighbor)
				targetRequest := network.TargetRequest{
					Ip:   &neighborIp,
					Port: &neighborPort,
				}
				targetRequests = append(targetRequests, targetRequest)
			}
		}
		neighborhood.neighborsMutex.RUnlock()
		rand.Seed(neighborhood.timeable.Now().UnixNano())
		rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
		outboundsCount := int(math.Min(float64(len(neighbors)), maxOutboundsCount))
		neighborhood.neighborsMutex.Lock()
		neighborhood.neighbors = neighbors[:outboundsCount]
		neighborhood.neighborsMutex.Unlock()
		for _, neighbor := range neighbors[:outboundsCount] {
			var neighborTargetRequests []network.TargetRequest
			for _, targetRequest := range targetRequests {
				neighborIp := neighbor.Ip()
				neighborPort := neighbor.Port()
				if neighborIp != *targetRequest.Ip || neighborPort != *targetRequest.Port {
					neighborTargetRequests = append(neighborTargetRequests, targetRequest)
				}
			}
			go func(neighbor network.Requestable) {
				_ = neighbor.SendTargets(neighborTargetRequests)
			}(neighbor)
		}
	}(neighborsByTarget)
}

func (neighborhood *Neighborhood) AddTargets(targetRequests []network.TargetRequest) {
	neighborhood.waitGroup.Add(1)
	go func() {
		defer neighborhood.waitGroup.Done()
		neighborhood.neighborsByTargetMutex.Lock()
		defer neighborhood.neighborsByTargetMutex.Unlock()
		for _, targetRequest := range targetRequests {
			neighbor := NewNeighbor(*targetRequest.Ip, *targetRequest.Port, neighborhood.logger)
			if _, ok := neighborhood.neighborsByTarget[neighbor.Target()]; !ok {
				neighborhood.neighborsByTarget[neighbor.Target()] = neighbor
			}
		}
	}()
}

func (neighborhood *Neighborhood) Neighbors() []network.Requestable {
	return neighborhood.neighbors
}
