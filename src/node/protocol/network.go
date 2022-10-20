package protocol

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"math"
	"math/rand"
	"net"
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

	timing                 clock.Timing
	neighbors              []*neighborhood.Neighbor
	neighborsMutex         sync.RWMutex
	neighborsByTarget      map[string]*neighborhood.Neighbor
	neighborsByTargetMutex sync.RWMutex
	seedsByTarget          map[string]*neighborhood.Neighbor
}

func NewNetwork(ip string, port uint16, timing clock.Timing, logger *log.Logger) *Network {
	network := new(Network)
	network.ip = ip
	network.port = port
	network.timing = timing
	network.logger = logger
	var waitGroup sync.WaitGroup
	network.waitGroup = &waitGroup
	seedsIps := []string{
		"89.82.76.241",
	}
	network.seedsByTarget = map[string]*neighborhood.Neighbor{}
	for _, seedIp := range seedsIps {
		seed := neighborhood.NewNeighbor(seedIp, DefaultPort, logger)
		network.seedsByTarget[seed.Target()] = seed
	}
	network.neighborsByTarget = map[string]*neighborhood.Neighbor{}
	return network
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
	var neighborsByTarget map[string]*neighborhood.Neighbor
	if len(network.neighborsByTarget) == 0 {
		neighborsByTarget = network.seedsByTarget
	} else {
		neighborsByTarget = network.neighborsByTarget
	}
	network.neighborsByTarget = map[string]*neighborhood.Neighbor{}
	network.neighborsByTargetMutex.Unlock()
	network.waitGroup.Add(1)
	go func(neighborsByTarget map[string]*neighborhood.Neighbor) {
		defer network.waitGroup.Done()
		var neighbors []*neighborhood.Neighbor
		var targetRequests []neighborhood.TargetRequest
		hostTargetRequest := neighborhood.TargetRequest{
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
				targetRequest := neighborhood.TargetRequest{
					Ip:   &neighborIp,
					Port: &neighborPort,
				}
				targetRequests = append(targetRequests, targetRequest)
			}
		}
		network.neighborsMutex.RUnlock()
		rand.Seed(network.timing.Now().UnixNano())
		rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
		outboundsCount := int(math.Min(float64(len(neighbors)), maxOutboundsCount))
		network.neighborsMutex.Lock()
		network.neighbors = neighbors[:outboundsCount]
		network.neighborsMutex.Unlock()
		for _, neighbor := range neighbors[:outboundsCount] {
			var neighborTargetRequests []neighborhood.TargetRequest
			for _, targetRequest := range targetRequests {
				neighborIp := neighbor.Ip()
				neighborPort := neighbor.Port()
				if neighborIp != *targetRequest.Ip || neighborPort != *targetRequest.Port {
					neighborTargetRequests = append(neighborTargetRequests, targetRequest)
				}
			}
			go func(neighbor *neighborhood.Neighbor) {
				_ = neighbor.SendTargets(neighborTargetRequests)
			}(neighbor)
		}
	}(neighborsByTarget)
}

func (network *Network) AddTargets(targetRequests []neighborhood.TargetRequest) {
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

func (network *Network) Neighbors() []*neighborhood.Neighbor {
	return network.neighbors
}
