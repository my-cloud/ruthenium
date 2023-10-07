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
	clientFactory            ClientFactory
	hostTarget               *Target
	maxOutboundsCount        int
	neighbors                []network.Neighbor
	neighborsMutex           sync.RWMutex
	scoresBySeedTargetValue  map[string]int
	scoresByTargetValue      map[string]int
	scoresByTargetValueMutex sync.RWMutex
	watch                    clock.Watch
}

func NewSynchronizer(clientFactory ClientFactory, hostIp string, hostPort string, maxOutboundsCount int, scoresBySeedTargetValue map[string]int, watch clock.Watch) *Synchronizer {
	synchronizer := new(Synchronizer)
	synchronizer.clientFactory = clientFactory
	synchronizer.hostTarget = NewTarget(hostIp, hostPort)
	synchronizer.maxOutboundsCount = maxOutboundsCount
	synchronizer.scoresBySeedTargetValue = scoresBySeedTargetValue
	synchronizer.scoresByTargetValue = map[string]int{}
	synchronizer.watch = watch
	return synchronizer
}

func (synchronizer *Synchronizer) AddTargets(targetValues []string) {
	synchronizer.scoresByTargetValueMutex.Lock()
	defer synchronizer.scoresByTargetValueMutex.Unlock()
	for _, targetValue := range targetValues {
		_, isTargetAlreadyKnown := synchronizer.scoresByTargetValue[targetValue]
		target, err := NewTargetFromValue(targetValue)
		if err != nil {
			continue
		}
		isTargetOnSameNetwork := synchronizer.hostTarget.IsSameNetworkId(target)
		if !isTargetAlreadyKnown && isTargetOnSameNetwork {
			synchronizer.scoresByTargetValue[targetValue] = 0
		}
	}
}

func (synchronizer *Synchronizer) HostTarget() string {
	return synchronizer.hostTarget.Value()
}

func (synchronizer *Synchronizer) Incentive(targetValue string) {
	synchronizer.scoresByTargetValueMutex.Lock()
	defer synchronizer.scoresByTargetValueMutex.Unlock()
	synchronizer.scoresByTargetValue[targetValue] += 1
}

func (synchronizer *Synchronizer) Neighbors() []network.Neighbor {
	return synchronizer.neighbors
}

func (synchronizer *Synchronizer) Synchronize(int64) {
	synchronizer.scoresByTargetValueMutex.Lock()
	var scoresByTargetValue map[string]int
	if len(synchronizer.scoresByTargetValue) == 0 {
		scoresByTargetValue = synchronizer.scoresBySeedTargetValue
	} else {
		scoresByTargetValue = synchronizer.scoresByTargetValue
	}
	synchronizer.scoresByTargetValue = map[string]int{}
	synchronizer.scoresByTargetValueMutex.Unlock()
	neighborsByScore := map[int][]network.Neighbor{}
	var targetValues []string
	hostTargetValue := synchronizer.hostTarget.Value()
	targetValues = append(targetValues, hostTargetValue)
	for targetValue, score := range scoresByTargetValue {
		if targetValue != hostTargetValue {
			neighborTarget, err := NewTargetFromValue(targetValue)
			if err != nil {
				continue
			}
			neighbor, err := NewNeighbor(neighborTarget, synchronizer.clientFactory)
			if err != nil {
				continue
			}
			neighborsByScore[score] = append(neighborsByScore[score], neighbor)
			targetValues = append(targetValues, targetValue)
		}
	}
	outbounds := synchronizer.selectOutbounds(neighborsByScore, len(scoresByTargetValue))
	synchronizer.neighborsMutex.Lock()
	synchronizer.neighbors = outbounds
	synchronizer.neighborsMutex.Unlock()
	for _, neighbor := range outbounds {
		var neighborTargetValues []string
		for _, targetValue := range targetValues {
			neighborTargetValue := neighbor.Target()
			if neighborTargetValue != targetValue {
				neighborTargetValues = append(neighborTargetValues, targetValue)
			}
		}
		go func(neighbor network.Neighbor) {
			_ = neighbor.SendTargets(neighborTargetValues)
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
