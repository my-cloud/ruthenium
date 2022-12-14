package p2p

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"io"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	maxOutboundsCount = 8
)

type Synchronizer struct {
	hostIp   string
	hostPort uint16

	watch         clock.Watch
	clientFactory ClientFactory

	neighbors             []network.Neighbor
	neighborsMutex        sync.RWMutex
	neighborsTargets      map[string]*Target
	neighborsTargetsMutex sync.RWMutex
	seedsTargets          map[string]*Target

	waitGroup *sync.WaitGroup
}

func NewSynchronizer(hostPort uint16, watch clock.Watch, clientFactory ClientFactory, configurationPath string, logger log.Logger) (synchronizer *Synchronizer, err error) {
	synchronizer = new(Synchronizer)
	synchronizer.hostIp, err = findPublicIp(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to find the public IP: %w", err)
	}
	synchronizer.hostPort = hostPort
	synchronizer.watch = watch
	synchronizer.clientFactory = clientFactory
	var waitGroup sync.WaitGroup
	synchronizer.waitGroup = &waitGroup
	err = synchronizer.readSeedsTargets(configurationPath, logger)
	if err != nil {
		return nil, err
	}
	synchronizer.neighborsTargets = map[string]*Target{}
	return synchronizer, nil
}

func (synchronizer *Synchronizer) Synchronize(int64) {
	synchronizer.neighborsTargetsMutex.Lock()
	var neighborsTargets map[string]*Target
	if len(synchronizer.neighborsTargets) == 0 {
		neighborsTargets = synchronizer.seedsTargets
	} else {
		neighborsTargets = synchronizer.neighborsTargets
	}
	synchronizer.neighborsTargets = map[string]*Target{}
	synchronizer.neighborsTargetsMutex.Unlock()
	var neighbors []network.Neighbor
	var targetRequests []network.TargetRequest
	hostTargetRequest := network.TargetRequest{
		Ip:   &synchronizer.hostIp,
		Port: &synchronizer.hostPort,
	}
	targetRequests = append(targetRequests, hostTargetRequest)
	synchronizer.neighborsMutex.RLock()
	for _, target := range neighborsTargets {
		targetIp := target.Ip()
		targetPort := target.Port()
		if targetIp != synchronizer.hostIp || targetPort != synchronizer.hostPort {
			var neighbor network.Neighbor
			neighbor, err := NewNeighbor(target, synchronizer.clientFactory)
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
	synchronizer.neighborsMutex.RUnlock()
	rand.Seed(synchronizer.watch.Now().UnixNano())
	rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
	outboundsCount := int(math.Min(float64(len(neighbors)), maxOutboundsCount))
	synchronizer.neighborsMutex.Lock()
	synchronizer.neighbors = neighbors[:outboundsCount]
	synchronizer.neighborsMutex.Unlock()
	for _, neighbor := range neighbors[:outboundsCount] {
		var neighborTargetRequests []network.TargetRequest
		for _, targetRequest := range targetRequests {
			neighborIp := neighbor.Ip()
			neighborPort := neighbor.Port()
			if neighborIp != *targetRequest.Ip || neighborPort != *targetRequest.Port {
				neighborTargetRequests = append(neighborTargetRequests, targetRequest)
			}
		}
		go func(neighbor network.Neighbor) {
			_ = neighbor.SendTargets(neighborTargetRequests)
		}(neighbor)
	}
}

func (synchronizer *Synchronizer) AddTargets(targetRequests []network.TargetRequest) {
	synchronizer.neighborsTargetsMutex.Lock()
	defer synchronizer.neighborsTargetsMutex.Unlock()
	for _, targetRequest := range targetRequests {
		target := NewTarget(*targetRequest.Ip, *targetRequest.Port)
		if _, ok := synchronizer.neighborsTargets[target.Value()]; !ok {
			synchronizer.neighborsTargets[target.Value()] = target
		}
	}
}

func (synchronizer *Synchronizer) Neighbors() []network.Neighbor {
	return synchronizer.neighbors
}

func findPublicIp(logger log.Logger) (ip string, err error) {
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
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ip = string(body)
	return
}

func (synchronizer *Synchronizer) readSeedsTargets(configurationPath string, logger log.Logger) error {
	jsonFile, err := os.Open(configurationPath + "/seeds.json")
	if err != nil {
		return fmt.Errorf("unable to open seeds IPs configuration file: %w", err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("unable to read seeds IPs configuration file: %w", err)
	}
	if err = jsonFile.Close(); err != nil {
		logger.Error(fmt.Errorf("unable to close seeds IPs configuration file: %w", err).Error())
	}
	var seedsStringTargets []string
	if err = json.Unmarshal(byteValue, &seedsStringTargets); err != nil {
		return fmt.Errorf("unable to unmarshal seeds IPs: %w", err)
	}
	synchronizer.seedsTargets = map[string]*Target{}
	for _, seedStringTarget := range seedsStringTargets {
		separator := ":"
		separatorIndex := strings.Index(seedStringTarget, separator)
		if separatorIndex == -1 {
			return fmt.Errorf("seed target format is invalid (should be of the form 89.82.76.241:8106")
		}
		seedIp := seedStringTarget[:separatorIndex]
		seedPortString := seedStringTarget[separatorIndex+1:]
		var seedPort uint64
		seedPort, err = strconv.ParseUint(seedPortString, 10, 16)
		if err != nil {
			return err
		}
		seedTarget := NewTarget(seedIp, uint16(seedPort))
		synchronizer.seedsTargets[seedTarget.Value()] = seedTarget
	}
	return nil
}
