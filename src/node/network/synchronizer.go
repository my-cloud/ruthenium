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
	"net/http"
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

	time          clock.Time
	senderFactory SenderFactory

	neighbors             []neighborhood.Neighbor
	neighborsMutex        sync.RWMutex
	neighborsTargets      map[string]*Target
	neighborsTargetsMutex sync.RWMutex
	seedsTargets          map[string]*Target

	waitGroup *sync.WaitGroup
	logger    *log.Logger
}

func NewSynchronizer(hostPort uint16, time clock.Time, senderFactory SenderFactory, configurationPath string, logger *log.Logger) (synchronizer *Synchronizer, err error) {
	synchronizer = new(Synchronizer)
	synchronizer.hostIp, err = findPublicIp(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to find the public IP: %w", err)
	}
	synchronizer.hostPort = hostPort
	synchronizer.time = time
	synchronizer.senderFactory = senderFactory
	synchronizer.logger = logger
	var waitGroup sync.WaitGroup
	synchronizer.waitGroup = &waitGroup
	seedsIps, err := readSeedsIps(configurationPath, logger)
	if err != nil {
		return nil, err
	}
	synchronizer.seedsTargets = map[string]*Target{}
	for _, seedIp := range seedsIps {
		seedTarget := NewTarget(seedIp, DefaultPort)
		synchronizer.seedsTargets[seedTarget.Value()] = seedTarget
	}
	synchronizer.neighborsTargets = map[string]*Target{}
	return synchronizer, nil
}

func findPublicIp(logger *log.Logger) (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return
	}
	defer func() {
		if bodyCloseError := resp.Body.Close(); bodyCloseError != nil {
			logger.Error(fmt.Errorf("failed to close public IP request body: %w", bodyCloseError).Error())
		}
	}()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ip = string(body)
	return
}

func readSeedsIps(configurationPath string, logger *log.Logger) ([]string, error) {
	jsonFile, err := os.Open(configurationPath + "/seeds.json")
	if err != nil {
		return nil, fmt.Errorf("unable to open seeds IPs configuration file: %w", err)
	}
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read seeds IPs configuration file: %w", err)
	}
	if err = jsonFile.Close(); err != nil {
		logger.Error(fmt.Errorf("unable to close seeds IPs configuration file: %w", err).Error())
	}
	var seedsIps []string
	if err = json.Unmarshal(byteValue, &seedsIps); err != nil {
		return nil, fmt.Errorf("unable to unmarshal seeds IPs: %w", err)
	}
	return seedsIps, nil
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
				neighbor, err := NewNeighbor(target, synchronizer.senderFactory, synchronizer.logger)
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
		rand.Seed(synchronizer.time.Now().UnixNano())
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
