package network

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	DefaultPort                  = 8106
	synchronizationTimeInSeconds = 10
	maxOutboundsCount            = 8
)

type Synchronizer struct {
	hostIp   string
	hostPort uint16

	timeable       clock.Time
	senderProvider SenderFactory

	neighbors             []neighborhood.Neighbor
	neighborsMutex        sync.RWMutex
	neighborsTargets      map[string]*Target
	neighborsTargetsMutex sync.RWMutex
	seedsTargets          map[string]*Target

	waitGroup *sync.WaitGroup
	logger    *log.Logger
}

func NewSynchronizer(hostIp string, hostPort uint16, timeable clock.Time, senderProvider SenderFactory, configurationPath string, logger *log.Logger) *Synchronizer {
	synchronizer := new(Synchronizer)
	synchronizer.hostIp = hostIp
	synchronizer.hostPort = hostPort
	synchronizer.timeable = timeable
	synchronizer.senderProvider = senderProvider
	synchronizer.logger = logger
	var waitGroup sync.WaitGroup
	synchronizer.waitGroup = &waitGroup
	seedsIps := readSeedsIps(configurationPath, logger)
	synchronizer.seedsTargets = map[string]*Target{}
	for _, seedIp := range seedsIps {
		seedTarget := NewTarget(seedIp, DefaultPort)
		synchronizer.seedsTargets[seedTarget.Value()] = seedTarget
	}
	synchronizer.neighborsTargets = map[string]*Target{}
	return synchronizer
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

func (synchronizer *Synchronizer) Wait() {
	synchronizer.waitGroup.Wait()
}

func (synchronizer *Synchronizer) StartSynchronization() {
	synchronizer.Synchronize()
	_ = time.AfterFunc(time.Second*synchronizationTimeInSeconds, synchronizer.StartSynchronization)
}

func (synchronizer *Synchronizer) Synchronize() {
	synchronizer.neighborsTargetsMutex.Lock()
	var neighborsTargets map[string]*Target
	if len(synchronizer.neighborsTargets) == 0 {
		neighborsTargets = synchronizer.seedsTargets
	} else {
		neighborsTargets = synchronizer.neighborsTargets
	}
	synchronizer.neighborsTargets = map[string]*Target{}
	synchronizer.neighborsTargetsMutex.Unlock()
	synchronizer.waitGroup.Add(1)
	go func(neighborsTargets map[string]*Target) {
		defer synchronizer.waitGroup.Done()
		var neighbors []neighborhood.Neighbor
		var targetRequests []neighborhood.TargetRequest
		hostTargetRequest := neighborhood.TargetRequest{
			Ip:   &synchronizer.hostIp,
			Port: &synchronizer.hostPort,
		}
		targetRequests = append(targetRequests, hostTargetRequest)
		synchronizer.neighborsMutex.RLock()
		for _, target := range neighborsTargets {
			targetIp := target.Ip()
			targetPort := target.Port()
			if targetIp != synchronizer.hostIp || targetPort != synchronizer.hostPort {
				var neighbor neighborhood.Neighbor
				neighbor, err := NewNeighbor(target, synchronizer.senderProvider, synchronizer.logger)
				if err == nil {
					neighbors = append(neighbors, neighbor)
					targetRequest := neighborhood.TargetRequest{
						Ip:   &targetIp,
						Port: &targetPort,
					}
					targetRequests = append(targetRequests, targetRequest)
				}
			}
		}
		synchronizer.neighborsMutex.RUnlock()
		rand.Seed(synchronizer.timeable.Now().UnixNano())
		rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
		outboundsCount := int(math.Min(float64(len(neighbors)), maxOutboundsCount))
		synchronizer.neighborsMutex.Lock()
		synchronizer.neighbors = neighbors[:outboundsCount]
		synchronizer.neighborsMutex.Unlock()
		for _, neighbor := range neighbors[:outboundsCount] {
			var neighborTargetRequests []neighborhood.TargetRequest
			for _, targetRequest := range targetRequests {
				neighborIp := neighbor.Ip()
				neighborPort := neighbor.Port()
				if neighborIp != *targetRequest.Ip || neighborPort != *targetRequest.Port {
					neighborTargetRequests = append(neighborTargetRequests, targetRequest)
				}
			}
			go func(neighbor neighborhood.Neighbor) {
				_ = neighbor.SendTargets(neighborTargetRequests)
			}(neighbor)
		}
	}(neighborsTargets)
}

func (synchronizer *Synchronizer) AddTargets(targetRequests []neighborhood.TargetRequest) {
	synchronizer.waitGroup.Add(1)
	go func() {
		defer synchronizer.waitGroup.Done()
		synchronizer.neighborsTargetsMutex.Lock()
		defer synchronizer.neighborsTargetsMutex.Unlock()
		for _, targetRequest := range targetRequests {
			target := NewTarget(*targetRequest.Ip, *targetRequest.Port)
			if _, ok := synchronizer.neighborsTargets[target.Value()]; !ok {
				synchronizer.neighborsTargets[target.Value()] = target
			}
		}
	}()
}

func (synchronizer *Synchronizer) Neighbors() []neighborhood.Neighbor {
	return synchronizer.neighbors
}
