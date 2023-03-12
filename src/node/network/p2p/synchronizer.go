package p2p

import (
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Synchronizer struct {
	clientFactory       ClientFactory
	hostTarget          *Target
	maxOutboundsCount   int
	neighbors           []network.Neighbor
	neighborsMutex      sync.RWMutex
	scoresBySeedTarget  map[string]int
	scoresByTarget      map[string]int
	scoresByTargetMutex sync.RWMutex
	watch               clock.Watch
}

func NewSynchronizer(clientFactory ClientFactory, hostIp string, hostPort string, maxOutboundsCount int, scoresBySeedTarget map[string]int, watch clock.Watch) *Synchronizer {
	synchronizer := new(Synchronizer)
	synchronizer.clientFactory = clientFactory
	synchronizer.hostTarget = NewTarget(hostIp, hostPort)
	synchronizer.maxOutboundsCount = maxOutboundsCount
	synchronizer.scoresBySeedTarget = scoresBySeedTarget
	synchronizer.scoresByTarget = map[string]int{}
	synchronizer.watch = watch
	return synchronizer
}

func (synchronizer *Synchronizer) AddTargets(targetRequests []network.TargetRequest) {
	synchronizer.scoresByTargetMutex.Lock()
	defer synchronizer.scoresByTargetMutex.Unlock()
	for _, targetRequest := range targetRequests {
		_, isTargetAlreadyKnown := synchronizer.scoresByTarget[*targetRequest.Target]
		target, err := NewTargetFromValue(*targetRequest.Target)
		if err != nil {
			continue
		}
		isTargetOnSameNetwork := synchronizer.hostTarget.IsSameNetworkId(target)
		if !isTargetAlreadyKnown && isTargetOnSameNetwork {
			synchronizer.scoresByTarget[*targetRequest.Target] = 0
		}
	}
}

func (synchronizer *Synchronizer) HostTarget() string {
	return synchronizer.hostTarget.Value()
}

func (synchronizer *Synchronizer) Incentive(target string) {
	synchronizer.scoresByTargetMutex.Lock()
	defer synchronizer.scoresByTargetMutex.Unlock()
	synchronizer.scoresByTarget[target] += 1
}

func (synchronizer *Synchronizer) Neighbors() []network.Neighbor {
	return synchronizer.neighbors
}

func (synchronizer *Synchronizer) Synchronize(int64) {
	synchronizer.scoresByTargetMutex.Lock()
	var scoresByTarget map[string]int
	if len(synchronizer.scoresByTarget) == 0 {
		scoresByTarget = synchronizer.scoresBySeedTarget
	} else {
		scoresByTarget = synchronizer.scoresByTarget
	}
	synchronizer.scoresByTarget = map[string]int{}
	synchronizer.scoresByTargetMutex.Unlock()
	neighborsByScore := map[int][]network.Neighbor{}
	var targetRequests []network.TargetRequest
	hostTargetValue := synchronizer.hostTarget.Value()
	hostTargetRequest := network.TargetRequest{
		Target: &hostTargetValue,
	}
	targetRequests = append(targetRequests, hostTargetRequest)
	for target, score := range scoresByTarget {
		if target != synchronizer.hostTarget.Value() {
			neighborTarget, err := NewTargetFromValue(target)
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
	outbounds := synchronizer.selectOutbounds(neighborsByScore, len(scoresByTarget))
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

func (synchronizer *Synchronizer) selectOutbounds(neighborsByScore map[int][]network.Neighbor, targetsCount int) []network.Neighbor {
	var keys []int
	for k := range neighborsByScore {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	outboundsCount := min(targetsCount, synchronizer.maxOutboundsCount)
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

func min(first, second int) int {
	var result int
	if first < second {
		result = first
	} else {
		result = second
	}
	return result
}
