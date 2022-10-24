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
	"os"
	"sync"
	"time"
)

const (
	DefaultPort                  = 8105
	synchronizationTimeInSeconds = 10
	maxOutboundsCount            = 8
)

type Neighborhood struct {
	hostIp    string
	hostPort  uint16
	logger    *log.Logger
	waitGroup *sync.WaitGroup

	timeable              clock.Timeable
	neighbors             []network.Requestable
	neighborsMutex        sync.RWMutex
	neighborsTargets      map[string]*Target
	neighborsTargetsMutex sync.RWMutex
	seedsTargets          map[string]*Target
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
	neighborhood.seedsTargets = map[string]*Target{}
	for _, seedIp := range seedsIps {
		seedTarget := NewTarget(seedIp, DefaultPort)
		neighborhood.seedsTargets[seedTarget.Value()] = seedTarget
	}
	neighborhood.neighborsTargets = map[string]*Target{}
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
	neighborhood.neighborsTargetsMutex.Lock()
	var neighborsTargets map[string]*Target
	if len(neighborhood.neighborsTargets) == 0 {
		neighborsTargets = neighborhood.seedsTargets
	} else {
		neighborsTargets = neighborhood.neighborsTargets
	}
	neighborhood.neighborsTargets = map[string]*Target{}
	neighborhood.neighborsTargetsMutex.Unlock()
	neighborhood.waitGroup.Add(1)
	go func(neighborsTargets map[string]*Target) {
		defer neighborhood.waitGroup.Done()
		var neighbors []network.Requestable
		var targetRequests []network.TargetRequest
		hostTargetRequest := network.TargetRequest{
			Ip:   &neighborhood.hostIp,
			Port: &neighborhood.hostPort,
		}
		targetRequests = append(targetRequests, hostTargetRequest)
		neighborhood.neighborsMutex.RLock()
		for _, target := range neighborsTargets {
			targetIp := target.Ip()
			targetPort := target.Port()
			if err := target.Reach(); err == nil && targetIp != neighborhood.hostIp || targetPort != neighborhood.hostPort {
				var neighbor network.Requestable
				neighbor, err = NewNeighbor(targetIp, targetPort, neighborhood.logger)
				if err == nil {
					neighbors = append(neighbors, neighbor)
					targetRequest := network.TargetRequest{
						Ip:   &targetIp,
						Port: &targetPort,
					}
					targetRequests = append(targetRequests, targetRequest)
				}
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
	}(neighborsTargets)
}

func (neighborhood *Neighborhood) AddTargets(targetRequests []network.TargetRequest) {
	neighborhood.waitGroup.Add(1)
	go func() {
		defer neighborhood.waitGroup.Done()
		neighborhood.neighborsTargetsMutex.Lock()
		defer neighborhood.neighborsTargetsMutex.Unlock()
		for _, targetRequest := range targetRequests {
			target := NewTarget(*targetRequest.Ip, *targetRequest.Port)
			if err := target.Reach(); err == nil {
				if _, ok := neighborhood.neighborsTargets[target.Value()]; !ok {
					neighborhood.neighborsTargets[target.Value()] = target
				}
			}
		}
	}()
}

func (neighborhood *Neighborhood) Neighbors() []network.Requestable {
	return neighborhood.neighbors
}
