package network

import (
	"github.com/my-cloud/ruthenium/validatornode/application"
	"math/rand"
	"sort"
	"sync"
	"time"
)

type Neighborhood struct {
	senderCreator            application.SenderCreator
	hostTarget               *Target
	maxOutboundsCount        int
	senders                  []application.Sender
	sendersMutex             sync.RWMutex
	scoresBySeedTargetValue  map[string]int
	scoresByTargetValue      map[string]int
	scoresByTargetValueMutex sync.RWMutex
	watch                    application.TimeProvider
}

func NewNeighborhood(senderCreator application.SenderCreator, hostIp string, hostPort string, maxOutboundsCount int, scoresBySeedTargetValue map[string]int, watch application.TimeProvider) *Neighborhood {
	neighborhood := new(Neighborhood)
	neighborhood.senderCreator = senderCreator
	neighborhood.hostTarget = NewTarget(hostIp, hostPort)
	neighborhood.maxOutboundsCount = maxOutboundsCount
	neighborhood.scoresBySeedTargetValue = scoresBySeedTargetValue
	neighborhood.scoresByTargetValue = map[string]int{}
	neighborhood.watch = watch
	return neighborhood
}

func (neighborhood *Neighborhood) AddTargets(targetValues []string) {
	neighborhood.scoresByTargetValueMutex.Lock()
	defer neighborhood.scoresByTargetValueMutex.Unlock()
	for _, targetValue := range targetValues {
		_, isTargetAlreadyKnown := neighborhood.scoresByTargetValue[targetValue]
		target, err := NewTargetFromValue(targetValue)
		if err != nil {
			continue
		}
		isTargetOnSameNetwork := neighborhood.hostTarget.IsSameNetworkId(target)
		if !isTargetAlreadyKnown && isTargetOnSameNetwork {
			neighborhood.scoresByTargetValue[targetValue] = 0
		}
	}
}

func (neighborhood *Neighborhood) HostTarget() string {
	return neighborhood.hostTarget.Value()
}

func (neighborhood *Neighborhood) Incentive(targetValue string) {
	neighborhood.scoresByTargetValueMutex.Lock()
	defer neighborhood.scoresByTargetValueMutex.Unlock()
	neighborhood.scoresByTargetValue[targetValue] += 1
}

func (neighborhood *Neighborhood) Senders() []application.Sender {
	return neighborhood.senders
}

func (neighborhood *Neighborhood) Synchronize(_ int64) {
	neighborhood.scoresByTargetValueMutex.Lock()
	var scoresByTargetValue map[string]int
	if len(neighborhood.scoresByTargetValue) == 0 {
		scoresByTargetValue = neighborhood.scoresBySeedTargetValue
	} else {
		scoresByTargetValue = neighborhood.scoresByTargetValue
	}
	neighborhood.scoresByTargetValue = map[string]int{}
	neighborhood.scoresByTargetValueMutex.Unlock()
	neighborsByScore := map[int][]application.Sender{}
	var targetValues []string
	hostTargetValue := neighborhood.hostTarget.Value()
	targetValues = append(targetValues, hostTargetValue)
	for targetValue, score := range scoresByTargetValue {
		if targetValue != hostTargetValue {
			neighborTarget, err := NewTargetFromValue(targetValue)
			if err != nil {
				continue
			}
			neighbor, err := neighborhood.senderCreator.CreateSender(neighborTarget.Ip(), neighborTarget.Port())
			if err != nil {
				continue
			}
			neighborsByScore[score] = append(neighborsByScore[score], neighbor)
			targetValues = append(targetValues, targetValue)
		}
	}
	outbounds := neighborhood.selectOutbounds(neighborsByScore, len(scoresByTargetValue))
	neighborhood.sendersMutex.Lock()
	neighborhood.senders = outbounds
	neighborhood.sendersMutex.Unlock()
	for _, neighbor := range outbounds {
		var neighborTargetValues []string
		for _, targetValue := range targetValues {
			neighborTargetValue := neighbor.Target()
			if neighborTargetValue != targetValue {
				neighborTargetValues = append(neighborTargetValues, targetValue)
			}
		}
		go func(neighbor application.Sender) {
			_ = neighbor.SendTargets(neighborTargetValues)
		}(neighbor)
	}
}

func (neighborhood *Neighborhood) selectOutbounds(neighborsByScore map[int][]application.Sender, targetsCount int) []application.Sender {
	var keys []int
	for k := range neighborsByScore {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	outboundsCount := min(targetsCount, neighborhood.maxOutboundsCount)
	var outbounds []application.Sender
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
