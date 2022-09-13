package blockchain

import (
	"errors"
	"fmt"
	"gitlab.com/coinsmaster/ruthenium/src/log"
	"gitlab.com/coinsmaster/ruthenium/src/node/neighborhood"
	"math"
	"math/rand"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	DefaultPort = 8106

	RewardSenderAddress        = "REWARD SENDER ADDRESS"
	ParticlesCount             = 100000000
	GenesisAmount       uint64 = 100000 * ParticlesCount
	RewardExponent             = 1 / 1.828393264
	HalfLifeInDays             = 373.59

	NeighborSynchronizationTimeInSeconds = 10
	MaxOutboundsCount                    = 8
	ConflictResolutionPerValidation      = 6
)

type Service struct {
	transactions      []*Transaction
	transactionsMutex sync.RWMutex
	blocks            []*Block
	blockResponses    []*neighborhood.BlockResponse
	blocksMutex       sync.RWMutex
	address           string
	miningStarted     bool
	mineRequested     bool
	miningTicker      *time.Ticker
	miningTimer       time.Duration

	ip        string
	port      uint16
	logger    *log.Logger
	waitGroup *sync.WaitGroup

	neighbors              []*neighborhood.Neighbor
	neighborsMutex         sync.RWMutex
	neighborsByTarget      map[string]*neighborhood.Neighbor
	neighborsByTargetMutex sync.RWMutex
	seedsByTarget          map[string]*neighborhood.Neighbor

	lambda float64
}

func NewService(address string, ip string, port uint16, miningTimer time.Duration, logger *log.Logger) *Service {
	service := new(Service)
	service.address = address
	service.ip = ip
	service.port = port
	service.miningTimer = miningTimer
	service.miningTicker = time.NewTicker(miningTimer)
	service.logger = logger
	var waitGroup sync.WaitGroup
	service.waitGroup = &waitGroup
	seedsIps := []string{
		"89.82.76.241",
	}
	service.seedsByTarget = map[string]*neighborhood.Neighbor{}
	for _, seedIp := range seedsIps {
		seed := neighborhood.NewNeighbor(seedIp, DefaultPort, logger)
		service.seedsByTarget[seed.Target()] = seed
	}
	service.neighborsByTarget = map[string]*neighborhood.Neighbor{}
	const hoursADay = 24
	halfLife := HalfLifeInDays * hoursADay * float64(time.Hour.Nanoseconds())
	service.lambda = math.Log(2) / halfLife
	return service
}

func (service *Service) WaitGroup() *sync.WaitGroup {
	return service.waitGroup
}

func (service *Service) Run() {
	go func() {
		service.logger.Info("updating the blockchain...")
		service.StartNeighborsSynchronization()
		service.WaitGroup().Wait()
		service.resolveConflicts()
		service.WaitGroup().Wait()
		service.logger.Info("the blockchain is now up to date")
		service.AddGenesisBlock()
		service.StartMining()
		service.StartConflictsResolution()
	}()
}

func (service *Service) AddGenesisBlock() {
	now := time.Now()
	parsedStartDate := now.Truncate(service.miningTimer).Add(service.miningTimer)
	deadline := parsedStartDate.Sub(now)
	service.miningTicker.Reset(deadline)
	<-service.miningTicker.C
	if service.blocks == nil {
		genesisTransaction := NewTransaction(service.address, RewardSenderAddress, nil, parsedStartDate.UnixNano(), GenesisAmount)
		service.addBlock(parsedStartDate.UnixNano(), genesisTransaction)
		service.logger.Debug("genesis block added")
	}
}

func (service *Service) StartNeighborsSynchronization() {
	service.SynchronizeNeighbors()
	_ = time.AfterFunc(time.Second*NeighborSynchronizationTimeInSeconds, service.StartNeighborsSynchronization)
}

func (service *Service) SynchronizeNeighbors() {
	service.neighborsByTargetMutex.Lock()
	var neighborsByTarget map[string]*neighborhood.Neighbor
	if len(service.neighborsByTarget) == 0 {
		neighborsByTarget = service.seedsByTarget
	} else {
		neighborsByTarget = service.neighborsByTarget
	}
	service.neighborsByTarget = map[string]*neighborhood.Neighbor{}
	service.neighborsByTargetMutex.Unlock()
	service.waitGroup.Add(1)
	go func(neighborsByTarget map[string]*neighborhood.Neighbor) {
		defer service.waitGroup.Done()
		var neighbors []*neighborhood.Neighbor
		var targetRequests []neighborhood.TargetRequest
		hostTargetRequest := neighborhood.TargetRequest{
			Ip:   &service.ip,
			Port: &service.port,
		}
		targetRequests = append(targetRequests, hostTargetRequest)
		service.neighborsMutex.RLock()
		for _, neighbor := range neighborsByTarget {
			neighborIp := neighbor.Ip()
			neighborPort := neighbor.Port()
			lookedUpNeighborsIps, err := net.LookupIP(neighborIp)
			if err != nil {
				service.logger.Error(fmt.Errorf("DNS discovery failed on addresse %s: %w", neighborIp, err).Error())
				return
			}

			neighborsCount := len(lookedUpNeighborsIps)
			if neighborsCount != 1 {
				service.logger.Error(fmt.Sprintf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", neighborsCount, neighborIp))
				return
			}
			lookedUpNeighborIp := lookedUpNeighborsIps[0]
			lookedUpNeighborIpString := lookedUpNeighborIp.String()
			if (lookedUpNeighborIpString != service.ip || neighborPort != service.port) && lookedUpNeighborIpString == neighborIp && neighbor.IsFound() {
				neighbors = append(neighbors, neighbor)
				targetRequest := neighborhood.TargetRequest{
					Ip:   &neighborIp,
					Port: &neighborPort,
				}
				targetRequests = append(targetRequests, targetRequest)
			}
		}
		service.neighborsMutex.RUnlock()
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(neighbors), func(i, j int) { neighbors[i], neighbors[j] = neighbors[j], neighbors[i] })
		outboundsCount := int(math.Min(float64(len(neighbors)), MaxOutboundsCount))
		service.neighborsMutex.Lock()
		service.neighbors = neighbors[:outboundsCount]
		service.neighborsMutex.Unlock()
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

func (service *Service) AddTargets(targetRequests []neighborhood.TargetRequest) {
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		service.neighborsByTargetMutex.Lock()
		defer service.neighborsByTargetMutex.Unlock()
		for _, targetRequest := range targetRequests {
			neighbor := neighborhood.NewNeighbor(*targetRequest.Ip, *targetRequest.Port, service.logger)
			if _, ok := service.neighborsByTarget[neighbor.Target()]; !ok {
				service.neighborsByTarget[neighbor.Target()] = neighbor
			}
		}
	}()
}

func (service *Service) AddTransaction(transaction *Transaction) {
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		err := service.addTransaction(transaction)
		if err != nil {
			service.logger.Debug(fmt.Errorf("failed to add transaction: %w", err).Error())
			return
		}
		transactionRequest := transaction.GetRequest()
		service.neighborsMutex.RLock()
		defer service.neighborsMutex.RUnlock()
		for _, neighbor := range service.neighbors {
			go func(neighbor *neighborhood.Neighbor) {
				_ = neighbor.AddTransaction(transactionRequest)
			}(neighbor)
		}
	}()
}

func (service *Service) addTransaction(transaction *Transaction) (err error) {
	service.blocksMutex.RLock()
	if len(service.blocks) > 2 {
		if transaction.Timestamp() < service.blocks[len(service.blocks)-2].Timestamp() {
			err = errors.New("the transaction timestamp is invalid")
			return
		}
		for i := len(service.blocks) - 2; i < len(service.blocks); i++ {
			for _, validatedTransaction := range service.blocks[i].transactions {
				if validatedTransaction.Equals(transaction) {
					err = errors.New("the transaction already is in the blockchain")
					return
				}
			}
		}
	}
	for _, pendingTransaction := range service.transactions {
		if pendingTransaction.Equals(transaction) {
			err = errors.New("the transaction already is in the transactions pool")
			return
		}
	}
	service.blocksMutex.RUnlock()
	if err = transaction.Verify(); err != nil {
		err = errors.New("failed to verify transaction")
		return
	}
	var senderWalletAmount uint64
	if len(service.blocks) > 0 {
		senderWalletAmount = service.CalculateTotalAmount(service.blocks[len(service.blocks)-1].Timestamp(), transaction.SenderAddress())
	}
	insufficientBalance := senderWalletAmount < transaction.Value()
	if insufficientBalance {
		err = errors.New("not enough balance in the sender wallet")
		return
	}
	service.transactionsMutex.Lock()
	defer service.transactionsMutex.Unlock()
	service.transactions = append(service.transactions, transaction)
	return
}

func (service *Service) Mine() {
	if service.miningStarted || service.mineRequested {
		return
	}
	now := time.Now()
	parsedStartDate := now.Truncate(service.miningTimer).Add(service.miningTimer)
	deadline := parsedStartDate.Sub(now)
	service.miningTicker.Reset(deadline)
	service.mineRequested = true
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		<-service.miningTicker.C
		service.mine()
		service.mineRequested = false
		if service.miningStarted {
			newNow := time.Now()
			newParsedStartDate := newNow.Truncate(service.miningTimer).Add(service.miningTimer)
			newDeadline := newParsedStartDate.Sub(newNow)
			service.miningTicker.Reset(newDeadline)
		} else {
			service.miningTicker.Stop()
		}
	}()
}

func (service *Service) mine() {
	service.blocksMutex.Lock()
	defer service.blocksMutex.Unlock()
	now := time.Now().Round(service.miningTimer).UnixNano()
	reward := service.calculateTotalReward(now, service.address, service.blocks)
	service.logger.Debug(fmt.Sprintf("reward: %d", reward))
	rewardTransaction := NewTransaction(service.address, RewardSenderAddress, nil, now, reward)
	service.addBlock(now, rewardTransaction)
}

func (service *Service) StartMining() {
	if !service.miningStarted {
		service.miningStarted = true
		go service.mining()
	}
}

func (service *Service) mining() {
	now := time.Now()
	parsedStartDate := now.Truncate(service.miningTimer).Add(service.miningTimer)
	deadline := parsedStartDate.Sub(now)
	service.miningTicker.Reset(deadline)
	<-service.miningTicker.C
	service.miningTicker.Reset(service.miningTimer)
	for {
		if !service.miningStarted {
			service.miningTicker.Stop()
			return
		}
		service.mine()
		<-service.miningTicker.C
	}
}

func (service *Service) StopMining() {
	service.miningStarted = false
	service.miningTicker.Reset(time.Nanosecond)
}

func (service *Service) CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64 {
	var totalAmount uint64
	var lastTimestamp int64
	service.blocksMutex.RLock()
	defer service.blocksMutex.RUnlock()
	for _, block := range service.blocks {
		for _, transaction := range block.Transactions() {
			value := transaction.Value()
			if blockchainAddress == transaction.RecipientAddress() {
				if totalAmount > 0 {
					totalAmount = service.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				totalAmount += value
				lastTimestamp = block.Timestamp()
			}
			if blockchainAddress == transaction.SenderAddress() {
				if totalAmount > 0 {
					totalAmount = service.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				if totalAmount < value {
					service.logger.Error(fmt.Sprintf("historical transaction should not have been validated: wallet amount=%d, transaction value=%d", totalAmount, value))
				}
				totalAmount -= value
				lastTimestamp = block.Timestamp()
			}
		}
	}
	return service.decay(lastTimestamp, currentTimestamp, totalAmount)
}
func (service *Service) calculateTotalReward(currentTimestamp int64, blockchainAddress string, blocks []*Block) uint64 {
	var totalAmount uint64
	var lastTimestamp int64
	var isValidatorKnown bool
	var totalReward uint64
	var rewardRecipientAddresses []string
	for _, block := range blocks {
		for _, transaction := range block.Transactions() {
			value := transaction.Value()
			if blockchainAddress == transaction.RecipientAddress() {
				if totalAmount > 0 {
					totalAmount = service.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				totalAmount += value
				lastTimestamp = block.Timestamp()
				if transaction.SenderAddress() == RewardSenderAddress {
					totalReward = 0
					rewardRecipientAddresses = nil
					if !isValidatorKnown {
						isValidatorKnown = true
					}
				}
			} else if transaction.SenderAddress() == RewardSenderAddress {
				for _, rewardRecipientAddress := range rewardRecipientAddresses {
					if transaction.RecipientAddress() == rewardRecipientAddress {
						isValidatorKnown = false
					}
				}
				rewardRecipientAddresses = append(rewardRecipientAddresses, transaction.RecipientAddress())
				if isValidatorKnown {
					reward := calculateReward(totalAmount)
					totalReward = service.decay(lastTimestamp, block.Timestamp(), totalReward) + reward
				} else {
					totalReward = 0
				}
			}
			if blockchainAddress == transaction.SenderAddress() {
				if totalAmount > 0 {
					totalAmount = service.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				if totalAmount < value {
					service.logger.Error(fmt.Sprintf("historical transaction should not have been validated: wallet amount=%d, transaction value=%d", totalAmount, value))
				}
				totalAmount -= value
				lastTimestamp = block.Timestamp()
			}
		}
	}
	totalAmount = service.decay(lastTimestamp, currentTimestamp, totalAmount)
	totalReward = service.decay(lastTimestamp, currentTimestamp, totalReward)
	return totalReward + calculateReward(totalAmount)
}

func calculateReward(amount uint64) uint64 {
	return uint64(math.Round(math.Pow(float64(amount), RewardExponent)))
}

func (service *Service) decay(lastTimestamp int64, newTimestamp int64, amount uint64) uint64 {
	elapsedTimestamp := newTimestamp - lastTimestamp
	return uint64(math.Floor(float64(amount) * math.Exp(-service.lambda*float64(elapsedTimestamp))))
}

func (service *Service) Transactions() []*Transaction {
	service.transactionsMutex.RLock()
	defer service.transactionsMutex.RUnlock()
	return service.transactions
}

func (service *Service) Blocks() []*neighborhood.BlockResponse {
	return service.blockResponses
}

func (service *Service) getValidBlocks(neighborBlocks []*neighborhood.BlockResponse) (validBlocks []*Block) {
	if len(neighborBlocks) == 0 || len(neighborBlocks) < len(service.blocks) {
		return
	}
	lastNeighborBlock, err := NewBlockFromResponse(neighborBlocks[len(neighborBlocks)-1])
	if err != nil {
		// TODO return error and log it at upper level precising the involved neighbor
		service.logger.Error(fmt.Errorf("failed to instantiate last neighbor block: %w", err).Error())
		return
	}
	if err = lastNeighborBlock.IsProofOfHumanityValid(); err != nil {
		service.logger.Error(fmt.Errorf("failed to get valid proof of humanity: %w", err).Error())
		return
	}
	previousBlock, err := NewBlockFromResponse(neighborBlocks[0])
	if err != nil {
		service.logger.Error(fmt.Errorf("failed to instantiate first neighbor block: %w", err).Error())
		return
	}
	validBlocks = append(validBlocks, previousBlock)
	for i := 1; i < len(neighborBlocks); i++ {
		var currentBlock *Block
		currentBlock, err = NewBlockFromResponse(neighborBlocks[i])
		if err != nil {
			service.logger.Error(fmt.Errorf("failed to instantiate last neighbor block: %w", err).Error())
			return
		}
		var previousBlockHash [32]byte
		previousBlockHash, err = previousBlock.Hash()
		if err != nil {
			service.logger.Error(fmt.Errorf("failed to calculate previous block hash: %w", err).Error())
			return
		}
		isPreviousHashValid := currentBlock.PreviousHash() == previousBlockHash
		if !isPreviousHashValid {
			service.logger.Debug("a hash is invalid for a neighbor")
			return
		}
		var isNewBlock bool
		if i >= len(service.blocks) {
			isNewBlock = true
		} else if len(service.blocks) > 2 {
			var hostBlockHash [32]byte
			var currentBlockHash [32]byte
			currentBlockHash, err = currentBlock.Hash()
			if err != nil {
				service.logger.Error(fmt.Errorf("failed to calculate neighbor block hash: %w", err).Error())
				return
			}
			hostBlockHash, err = service.blocks[i].Hash()
			if err != nil {
				service.logger.Error(fmt.Errorf("failed to calculate host block hash: %w", err).Error())
			} else if currentBlockHash != hostBlockHash {
				isNewBlock = true
			}
		}

		if isNewBlock {
			var rewarded bool
			for _, transaction := range currentBlock.Transactions() {
				if transaction.SenderAddress() == RewardSenderAddress {
					// Check that there is only one reward by block
					if rewarded {
						service.logger.Error("multiple rewards attempt for the same block")
						return
					}
					rewarded = true
					currentBlockTimestamp := currentBlock.Timestamp()
					previousBlockTimestamp := previousBlock.Timestamp()
					now := time.Now().UnixNano()
					if currentBlockTimestamp != previousBlockTimestamp+int64(service.miningTimer) || currentBlockTimestamp > now {
						service.logger.Error("reward timestamp is invalid")
						return
					}
					if transaction.Value() > service.calculateTotalReward(currentBlockTimestamp, transaction.RecipientAddress(), validBlocks) {
						service.logger.Error("reward exceeds the consented one")
						return
					}
				}
			}
		}
		validBlocks = append(validBlocks, currentBlock)
		previousBlock = currentBlock
	}
	return
}

func (service *Service) StartConflictsResolution() {
	consensusTimer := service.miningTimer / ConflictResolutionPerValidation
	consensusTicker := time.Tick(consensusTimer)
	go func() {
		for {
			for i := 0; i < ConflictResolutionPerValidation; i++ {
				if i > 0 || (!service.miningStarted && !service.mineRequested) {
					service.resolveConflicts()
				}
				<-consensusTicker
			}
		}
	}()
}

func (service *Service) resolveConflicts() {
	service.WaitGroup().Add(1)
	go func() {
		defer service.WaitGroup().Done()
		// Select valid blocks
		blockResponsesByNeighbor := make(map[*neighborhood.Neighbor][]*neighborhood.BlockResponse)
		blocksByNeighbor := make(map[*neighborhood.Neighbor][]*Block)
		var selectedNeighbors []*neighborhood.Neighbor
		service.neighborsMutex.RLock()
		for _, neighbor := range service.neighbors {
			neighborBlocks, err := neighbor.GetBlocks()
			blockResponsesByNeighbor[neighbor] = neighborBlocks
			if err == nil {
				validBlocks := service.getValidBlocks(neighborBlocks)
				if validBlocks != nil {
					blocksByNeighbor[neighbor] = validBlocks
					selectedNeighbors = append(selectedNeighbors, neighbor)
				}
			}
		}
		service.neighborsMutex.RUnlock()

		service.blocksMutex.Lock()
		defer service.blocksMutex.Unlock()
		service.transactionsMutex.Lock()
		defer service.transactionsMutex.Unlock()
		if len(service.blocks) > 2 {
			host := neighborhood.NewNeighbor(service.ip, service.port, service.logger)
			blockResponsesByNeighbor[host] = service.blockResponses
			blocksByNeighbor[host] = service.blocks
			selectedNeighbors = append(selectedNeighbors, host)
		}

		// Select blockchain with consensus for the previous hash
		var selectedBlocksResponse []*neighborhood.BlockResponse
		var selectedBlocks []*Block
		if selectedNeighbors != nil && service.neighbors != nil {
			service.sortByBlockLength(selectedNeighbors, blocksByNeighbor)
			for len(blocksByNeighbor) > 0 {
				fiftyOnePercent := len(blocksByNeighbor)/2 + 1
				shortestBlocks := blocksByNeighbor[selectedNeighbors[fiftyOnePercent-1]]
				shortestBlocksLength := len(shortestBlocks)
				shortestBlocksSlicePreviousHash := shortestBlocks[len(shortestBlocks)-1].PreviousHash()
				var samePreviousHashesCount int
				for i := 0; i < fiftyOnePercent; i++ {
					if blocksByNeighbor[selectedNeighbors[i]][shortestBlocksLength-1].PreviousHash() == shortestBlocksSlicePreviousHash {
						samePreviousHashesCount++
					}
				}
				if samePreviousHashesCount == fiftyOnePercent {
					// Keep the longest chain with the oldest reward recipient address
					maxLength := len(service.blocks)
					var maxRewardRecipientAddressAge uint64
					for neighbor, blocks := range blocksByNeighbor {
						var rewardRecipientAddressAge uint64
						var neighborLastBlockRewardRecipientAddress string
						for _, transaction := range blocks[len(blocks)-1].transactions {
							if transaction.SenderAddress() == RewardSenderAddress {
								neighborLastBlockRewardRecipientAddress = transaction.RecipientAddress()
							}
						}
						if len(blocks) > maxLength {
							maxLength = len(blocks)
							selectedBlocksResponse = blockResponsesByNeighbor[neighbor]
							selectedBlocks = blocks
							maxRewardRecipientAddressAge = 0
						} else if len(blocks) == maxLength {
							var isAgeCalculated bool
							for i := len(service.blocks) - 2; i > 0; i-- {
								for _, transaction := range service.blocks[i].transactions {
									if transaction.SenderAddress() == RewardSenderAddress {
										if transaction.RecipientAddress() == neighborLastBlockRewardRecipientAddress {
											isAgeCalculated = true
										}
										rewardRecipientAddressAge++
										break
									}
								}
								if isAgeCalculated {
									break
								}
							}
							if rewardRecipientAddressAge > maxRewardRecipientAddressAge {
								maxRewardRecipientAddressAge = rewardRecipientAddressAge
								selectedBlocksResponse = blockResponsesByNeighbor[neighbor]
								selectedBlocks = blocks
							}
						}
					}
					break
				}
				var rejectedNeighbors []*neighborhood.Neighbor
				if samePreviousHashesCount < fiftyOnePercent/2+1 {
					// The previous hash of the blockchain used to compare is not shared by less than 51% neighbors, reject it and all neighbors with the same previous hash
					for i := 0; i < fiftyOnePercent; i++ {
						selectedNeighbor := selectedNeighbors[i]
						if blocksByNeighbor[selectedNeighbor][shortestBlocksLength-1].PreviousHash() == shortestBlocksSlicePreviousHash {
							delete(blocksByNeighbor, selectedNeighbor)
							delete(blockResponsesByNeighbor, selectedNeighbor)
							rejectedNeighbors = append(rejectedNeighbors, selectedNeighbor)
						}
					}
				} else {
					// The previous hash of the blockchain used to compare is shared by at least 51% neighbors, reject other neighbors
					for i := 0; i < fiftyOnePercent; i++ {
						selectedNeighbor := selectedNeighbors[i]
						if blocksByNeighbor[selectedNeighbor][shortestBlocksLength-1].PreviousHash() != shortestBlocksSlicePreviousHash {
							delete(blocksByNeighbor, selectedNeighbor)
							delete(blockResponsesByNeighbor, selectedNeighbor)
							rejectedNeighbors = append(rejectedNeighbors, selectedNeighbor)
						}
					}
				}
				for _, rejectedNeighbor := range rejectedNeighbors {
					remove(selectedNeighbors, rejectedNeighbor)
				}
			}
			if selectedBlocks != nil {
				// Check if blockchain is replaced
				var blockchainReplaced bool
				if len(selectedBlocks) >= 2 {
					if len(service.blocks) < 2 {
						blockchainReplaced = true
					} else {
						lastNewBlockHash, newBlockHashError := selectedBlocks[len(selectedBlocks)-1].Hash()
						penultimateNewBlockHash, newBlockHashError := selectedBlocks[len(selectedBlocks)-2].Hash()
						if newBlockHashError != nil {
							service.logger.Error("failed to calculate new block hash")
						} else {
							lastOldBlockHash, oldBlockHashError := service.blocks[len(service.blocks)-1].Hash()
							penultimateOldBlockHash, oldBlockHashError := service.blocks[len(service.blocks)-2].Hash()
							if oldBlockHashError != nil {
								service.logger.Error("failed to calculate old block hash")
								blockchainReplaced = true
							} else {
								blockchainReplaced = penultimateOldBlockHash != penultimateNewBlockHash || lastOldBlockHash != lastNewBlockHash
							}
						}
					}
				}
				if blockchainReplaced {
					if len(service.blocks) > 2 {
						transactions := service.transactions
						// Add transactions which are not in the new blocks but the rewards
						for i := len(service.blocks) - 2; i < len(service.blocks); i++ {
							for _, invalidatedTransaction := range service.blocks[i].transactions {
								if invalidatedTransaction.SenderAddress() != RewardSenderAddress {
									transactions = append(transactions, invalidatedTransaction)
								}
							}
						}
						// Remove transactions which are in the new blocks
						for i := len(service.blocks) - 2; i < len(selectedBlocks); i++ {
							for _, validatedTransaction := range selectedBlocks[i].transactions {
								for j := 0; j < len(transactions); j++ {
									// FIXME verify that len(transactions) don't change at each iteration
									if validatedTransaction.Equals(transactions[j]) {
										transactions[j] = transactions[len(transactions)-1]
										transactions = transactions[:len(transactions)-1]
										j--
									}
								}
							}
						}
						service.transactions = transactions
					}
					service.blockResponses = selectedBlocksResponse
					service.blocks = selectedBlocks
					service.logger.Debug("conflicts resolved: blockchain replaced")
				} else {
					service.logger.Debug("conflicts resolved: blockchain kept")
				}
			} else {
				service.logger.Debug("conflicts resolved: blockchain kept")
			}
		}
	}()
}

func (service *Service) sortByBlockLength(selectedNeighbors []*neighborhood.Neighbor, blocksByNeighbor map[*neighborhood.Neighbor][]*Block) {
	sort.Slice(selectedNeighbors, func(i, j int) bool {
		return len(blocksByNeighbor[selectedNeighbors[i]]) > len(blocksByNeighbor[selectedNeighbors[j]])
	})
}

func remove(neighbors []*neighborhood.Neighbor, removedNeighbor *neighborhood.Neighbor) []*neighborhood.Neighbor {
	for i := 0; i < len(neighbors); i++ {
		if neighbors[i] == removedNeighbor {
			neighbors = append(neighbors[:i], neighbors[i+1:]...)
			i--
			return neighbors
		}
	}
	return neighbors
}

func (service *Service) addBlock(timestamp int64, rewardTransaction *Transaction) {
	var lastBlockHash [32]byte
	if service.blocks != nil {
		lastBlock := service.blocks[len(service.blocks)-1]
		if lastBlock.Timestamp() == timestamp {
			service.logger.Error("unable to add block, a block with the same timestamp already is in the blockchain")
			return
		}
		var err error
		lastBlockHash, err = lastBlock.Hash()
		if err != nil {
			service.logger.Error(fmt.Errorf("failed calculate last block hash: %w", err).Error())
			return
		}
	}
	service.transactionsMutex.Lock()
	defer service.transactionsMutex.Unlock()
	service.transactions = append(service.transactions, rewardTransaction)
	block := NewBlock(timestamp, lastBlockHash, service.transactions)
	service.transactions = nil
	service.blocks = append(service.blocks, block)
	blockResponse := block.GetResponse()
	service.blockResponses = append(service.blockResponses, blockResponse)
	return
}
