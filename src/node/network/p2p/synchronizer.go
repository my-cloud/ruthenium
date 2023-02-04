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
	"sort"
	"sync"
	"time"
)

const (
	maxOutboundsCount = 8
)

type Synchronizer struct {
	hostTarget string

	watch         clock.Watch
	clientFactory ClientFactory

	neighbors             []network.Neighbor
	neighborsMutex        sync.RWMutex
	scoresByNeighbor      map[string]int
	scoresByNeighborMutex sync.RWMutex
	scoresBySeed          map[string]int

	waitGroup *sync.WaitGroup
}

func NewSynchronizer(hostPort string, watch clock.Watch, clientFactory ClientFactory, configurationPath string, logger log.Logger) (synchronizer *Synchronizer, err error) {
	synchronizer = new(Synchronizer)
	hostIp, err := findPublicIp(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to find the public IP: %w", err)
	}
	synchronizer.hostTarget = fmt.Sprint(hostIp, ":", hostPort)
	synchronizer.watch = watch
	synchronizer.clientFactory = clientFactory
	var waitGroup sync.WaitGroup
	synchronizer.waitGroup = &waitGroup
	err = synchronizer.readSeedsTargets(configurationPath, logger)
	if err != nil {
		return nil, err
	}
	synchronizer.scoresByNeighbor = map[string]int{}
	return synchronizer, nil
}

func (synchronizer *Synchronizer) Synchronize(int64) {
	synchronizer.scoresByNeighborMutex.Lock()
	var scoresByNeighbor map[string]int
	if len(synchronizer.scoresByNeighbor) == 0 {
		scoresByNeighbor = synchronizer.scoresBySeed
	} else {
		scoresByNeighbor = synchronizer.scoresByNeighbor
	}
	synchronizer.scoresByNeighbor = map[string]int{}
	synchronizer.scoresByNeighborMutex.Unlock()
	neighborsByScore := map[int][]network.Neighbor{}
	var targetRequests []network.TargetRequest
	hostTargetRequest := network.TargetRequest{
		Target: &synchronizer.hostTarget,
	}
	targetRequests = append(targetRequests, hostTargetRequest)
	for target, score := range scoresByNeighbor {
		if target != synchronizer.hostTarget {
			neighborTarget, err := NewTarget(target)
			if err != nil {
				continue
			}
			neighbor, err := NewNeighbor(neighborTarget, synchronizer.clientFactory)
			if err != nil {
				continue
			}
			neighborsByScore[score] = append(neighborsByScore[score], neighbor)
			targetRequest := network.TargetRequest{
				Target: &target,
			}
			targetRequests = append(targetRequests, targetRequest)
		}
	}
	outbounds := pickOutbounds(neighborsByScore)
	synchronizer.neighborsMutex.Lock()
	synchronizer.neighbors = outbounds
	synchronizer.neighborsMutex.Unlock()
	for _, neighbor := range outbounds {
		var neighborTargetRequests []network.TargetRequest
		for _, targetRequest := range targetRequests {
			neighborTarget := neighbor.Target()
			if neighborTarget != *targetRequest.Target {
				neighborTargetRequests = append(neighborTargetRequests, targetRequest)
			}
		}
		go func(neighbor network.Neighbor) {
			_ = neighbor.SendTargets(neighborTargetRequests)
		}(neighbor)
	}
}

func (synchronizer *Synchronizer) AddTargets(targetRequests []network.TargetRequest) {
	synchronizer.scoresByNeighborMutex.Lock()
	defer synchronizer.scoresByNeighborMutex.Unlock()
	for _, targetRequest := range targetRequests {
		if _, ok := synchronizer.scoresByNeighbor[*targetRequest.Target]; !ok {
			synchronizer.scoresByNeighbor[*targetRequest.Target] = 0
		}
		senderTarget := targetRequests[0].Target
		synchronizer.scoresByNeighbor[*senderTarget] += 1
	}
}

func (synchronizer *Synchronizer) Incentive(target string) {
	synchronizer.scoresByNeighborMutex.Lock()
	defer synchronizer.scoresByNeighborMutex.Unlock()
	synchronizer.scoresByNeighbor[target] += 1
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

func pickOutbounds(neighborsByScore map[int][]network.Neighbor) []network.Neighbor {
	keys := make([]int, len(neighborsByScore))
	for k := range neighborsByScore {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	outboundsCount := int(math.Min(float64(len(neighborsByScore)), maxOutboundsCount))
	var outbounds []network.Neighbor
	for i := len(keys) - 1; i >= 0; i-- {
		if len(outbounds)+len(neighborsByScore[keys[i]]) >= outboundsCount {
			temp := neighborsByScore[keys[i]]
			rand.Seed(time.Now().UnixNano())
			rand.Shuffle(len(temp), func(i, j int) { temp[i], temp[j] = temp[j], temp[i] })
			outbounds = append(outbounds, temp[:outboundsCount-len(outbounds)]...)
			break
		}
		outbounds = append(outbounds, neighborsByScore[keys[i]]...)
	}
	return outbounds
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
	synchronizer.scoresBySeed = map[string]int{}
	for _, seedStringTarget := range seedsStringTargets {
		synchronizer.scoresBySeed[seedStringTarget] = 0
	}
	return nil
}
