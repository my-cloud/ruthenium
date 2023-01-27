package verification

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"math"
	"sort"
	"sync"
	"time"
)

const (
	incomeExponent = 0.54692829
	halfLifeInDays = 373.59
)

type Blockchain struct {
	blocks          []*Block
	blockResponses  []*network.BlockResponse
	lambda          float64
	mutex           sync.RWMutex
	registry        protocol.Registry
	synchronizer    network.Synchronizer
	validationTimer time.Duration
	logger          log.Logger
}

func NewBlockchain(genesisAmount uint64, initialTimestamp int64, registry protocol.Registry, validatorAddress string, validationTimer time.Duration, synchronizer network.Synchronizer, logger log.Logger) *Blockchain {
	blockchain := newBlockchain(nil, registry, validationTimer, synchronizer, logger)
	blockchain.addGenesisBlock(initialTimestamp, genesisAmount, validatorAddress)
	return blockchain
}

func newBlockchain(blockResponses []*network.BlockResponse, registry protocol.Registry, validationTimer time.Duration, synchronizer network.Synchronizer, logger log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.blockResponses = blockResponses
	blockchain.registry = registry
	blockchain.validationTimer = validationTimer
	blockchain.synchronizer = synchronizer
	const hoursADay = 24
	halfLife := halfLifeInDays * hoursADay * float64(time.Hour.Nanoseconds())
	blockchain.lambda = math.Log(2) / halfLife
	blockchain.logger = logger
	return blockchain
}

func (blockchain *Blockchain) AddBlock(timestamp int64, transactions []*network.TransactionResponse, registeredAddresses []string) {
	var previousHash [32]byte
	var err error
	blockchain.mutex.Lock()
	defer blockchain.mutex.Unlock()
	if !blockchain.isEmpty() {
		previousHash, err = blockchain.blocks[len(blockchain.blocks)-1].Hash()
		if err != nil {
			blockchain.logger.Error(fmt.Errorf("unable to calculate last block hash: %w", err).Error())
			return
		}
	}
	blockResponse := NewBlockResponse(timestamp, previousHash, transactions, registeredAddresses)
	block, err := NewBlockFromResponse(blockResponse)
	if err != nil {
		blockchain.logger.Error(fmt.Errorf("unable to instantiate block: %w", err).Error())
		return
	}
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.blocks = append(blockchain.blocks, block)
}

func (blockchain *Blockchain) Blocks() []*network.BlockResponse {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	return blockchain.blockResponses
}

func (blockchain *Blockchain) CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64 {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	var totalAmount uint64
	var lastTimestamp int64
	for _, block := range blockchain.blockResponses {
		for _, registeredAddress := range block.RegisteredAddresses {
			if blockchainAddress == registeredAddress {
				if totalAmount > 0 {
					totalAmount = blockchain.decay(lastTimestamp, block.Timestamp, totalAmount)
					totalAmount += calculateIncome(totalAmount)
					lastTimestamp = block.Timestamp
				}
				break
			}
		}
		for _, transaction := range block.Transactions {
			value := transaction.Value
			if blockchainAddress == transaction.RecipientAddress {
				if totalAmount > 0 {
					totalAmount = blockchain.decay(lastTimestamp, block.Timestamp, totalAmount)
				}
				totalAmount += value
				lastTimestamp = block.Timestamp
			} else if blockchainAddress == transaction.SenderAddress {
				if totalAmount > 0 {
					totalAmount = blockchain.decay(lastTimestamp, block.Timestamp, totalAmount)
				}
				if totalAmount < value+transaction.Fee {
					blockchain.logger.Error(fmt.Sprintf("historical transaction have not been properly validated: wallet amount=%d, transaction value=%d", totalAmount, value))
					totalAmount = 0
				} else {
					totalAmount -= value + transaction.Fee
				}
				lastTimestamp = block.Timestamp
			}
		}
	}
	return blockchain.decay(lastTimestamp, currentTimestamp, totalAmount)
}

func (blockchain *Blockchain) Copy() protocol.Blockchain {
	blockchainCopy := new(Blockchain)
	blockchainCopy.registry = blockchain.registry
	blockchainCopy.validationTimer = blockchain.validationTimer
	blockchainCopy.synchronizer = blockchain.synchronizer
	blockchainCopy.logger = blockchain.logger
	blockchainCopy.lambda = blockchain.lambda
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	blockchainCopy.blocks = blockchain.blocks
	blockchainCopy.blockResponses = blockchain.blockResponses
	return blockchainCopy
}

func (blockchain *Blockchain) LastBlocks(startingBlockHash [32]byte) []*network.BlockResponse {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	for i := len(blockchain.blockResponses); i > 0; i-- {
		currentBlockPreviousHash := &blockchain.blockResponses[i-1].PreviousHash
		if bytes.Equal(currentBlockPreviousHash[:], startingBlockHash[:]) {
			return blockchain.blockResponses[i-2:]
		}
	}
	return nil
}

func (blockchain *Blockchain) Update(timestamp int64) {
	// Verify neighbor blockchains
	neighbors := blockchain.synchronizer.Neighbors()
	blockResponsesByTarget := make(map[string][]*network.BlockResponse)
	blocksByTarget := make(map[string][]*Block)
	var selectedTargets []string
	hostBlocks := blockchain.blocks
	hostBlockResponses := blockchain.blockResponses
	var lastHostBlocks []*Block
	var oldHostBlocks []*network.BlockResponse
	if len(hostBlocks) > 2 {
		hostTarget := "host"
		blockResponsesByTarget[hostTarget] = hostBlockResponses
		blocksByTarget[hostTarget] = hostBlocks
		selectedTargets = append(selectedTargets, hostTarget)
		lastHostBlocks = hostBlocks[len(hostBlocks)-3:]
		oldHostBlocks = hostBlockResponses[:len(hostBlocks)-3]
	}
	for _, neighbor := range neighbors {
		var err error
		var neighborBlocks []*network.BlockResponse
		if len(hostBlocks) > 2 {
			startingBlockHash := hostBlocks[len(hostBlocks)-2].PreviousHash()
			lastBlocksRequest := network.LastBlocksRequest{StartingBlockHash: &startingBlockHash}
			neighborBlocks, err = neighbor.GetLastBlocks(lastBlocksRequest)
			if err == nil {
				target := neighbor.Target()
				var validBlocks []*Block
				validBlocks, err = blockchain.verify(neighborBlocks, lastHostBlocks, oldHostBlocks, timestamp)
				if err != nil || validBlocks == nil {
					blockchain.logger.Debug(fmt.Errorf("failed to verify blocks for neighbor %s: %w", target, err).Error())
				} else {
					blocksByTarget[target] = append(hostBlocks[:len(hostBlocks)-3], validBlocks...)
					blockResponsesByTarget[target] = append(hostBlockResponses[:len(hostBlockResponses)-3], neighborBlocks...)
					selectedTargets = append(selectedTargets, target)
				}
			}
		} else {
			neighborBlocks, err = neighbor.GetBlocks()
			if err == nil {
				target := neighbor.Target()
				var validBlocks []*Block
				validBlocks, err = blockchain.verify(neighborBlocks, hostBlocks, oldHostBlocks, timestamp)
				if err != nil || validBlocks == nil {
					blockchain.logger.Debug(fmt.Errorf("failed to verify blocks for neighbor %s: %w", target, err).Error())
				} else {
					blocksByTarget[target] = validBlocks
					blockResponsesByTarget[target] = neighborBlocks
					selectedTargets = append(selectedTargets, target)
				}
			}
		}
	}

	var selectedBlockResponses []*network.BlockResponse
	var selectedBlocks []*Block
	var isDifferent bool
	if selectedTargets != nil {
		// Keep blockchains with consensus for the previous hash (prevent forks)
		blockchain.sortByBlocksLength(selectedTargets, blocksByTarget)
		halfNeighborsCount := len(blocksByTarget) / 2
		minLength := len(blocksByTarget[selectedTargets[len(selectedTargets)-1]])
		var rejectedTargets []string
		for target, blocks := range blocksByTarget {
			var samePreviousHashCount int
			for _, otherBlocks := range blocksByTarget {
				if blocks[minLength-1].PreviousHash() == otherBlocks[minLength-1].PreviousHash() {
					samePreviousHashCount++
				}
			}
			if samePreviousHashCount <= halfNeighborsCount {
				// The previous hash of the blockchain used to compare is shared by at least 51% neighbors, reject other neighbors
				rejectedTargets = append(rejectedTargets, target)
			}
		}
		for _, rejectedTarget := range rejectedTargets {
			delete(blocksByTarget, rejectedTarget)
			delete(blockResponsesByTarget, rejectedTarget)
			removeTarget(selectedTargets, rejectedTarget)
		}
		// Keep the longest blockchains
		maxLength := len(blocksByTarget[selectedTargets[0]])
		rejectedTargets = nil
		for target, blocks := range blocksByTarget {
			if len(blocks) < maxLength {
				rejectedTargets = append(rejectedTargets, target)
			}
		}
		for _, rejectedTarget := range rejectedTargets {
			delete(blocksByTarget, rejectedTarget)
			delete(blockResponsesByTarget, rejectedTarget)
			removeTarget(selectedTargets, rejectedTarget)
		}
		// Select the blockchain of the oldest reward recipient
		var maxRewardRecipientAddressAge uint64
		for target, blocks := range blocksByTarget {
			var rewardRecipientAddressAge uint64
			var lastBlockRewardRecipientAddress string
			for _, transaction := range blocks[len(blocks)-1].transactions {
				if transaction.IsReward() {
					lastBlockRewardRecipientAddress = transaction.RecipientAddress()
				}
			}
			var isAgeCalculated bool
			for i := len(blocks) - 2; i >= 0; i-- {
				for _, transaction := range blocks[i].transactions {
					if transaction.IsReward() {
						if transaction.RecipientAddress() == lastBlockRewardRecipientAddress {
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
				selectedBlockResponses = blockResponsesByTarget[target]
				selectedBlocks = blocks
			}
		}
		// Check if blockchain is different to know if it should be updated
		if len(hostBlocks) < 2 && selectedBlocks != nil || len(hostBlocks) < len(selectedBlocks) {
			isDifferent = true
		} else if len(selectedBlocks) >= 2 {
			lastNewBlockHash, newBlockHashError := selectedBlocks[len(selectedBlocks)-1].Hash()
			if newBlockHashError != nil {
				blockchain.logger.Error("failed to calculate new block hash")
			} else {
				lastOldBlockHash, oldBlockHashError := hostBlocks[len(hostBlocks)-1].Hash()
				if oldBlockHashError != nil {
					blockchain.logger.Error("failed to calculate old block hash")
					isDifferent = true
				} else {
					isDifferent = lastOldBlockHash != lastNewBlockHash
				}
			}
		}
	}
	if isDifferent {
		blockchain.mutex.Lock()
		defer blockchain.mutex.Unlock()
		blockchain.blockResponses = selectedBlockResponses
		blockchain.blocks = selectedBlocks
		blockchain.logger.Debug("verification done: blockchain replaced")
	} else {
		blockchain.logger.Debug("verification done: blockchain kept")
	}
}

func (blockchain *Blockchain) addGenesisBlock(timestamp int64, genesisAmount uint64, validatorAddress string) {
	genesisTransaction := validation.NewRewardTransaction(validatorAddress, timestamp, genesisAmount)
	transactions := []*network.TransactionResponse{genesisTransaction}
	blockResponse := NewBlockResponse(timestamp, [32]byte{}, transactions, nil)
	block, err := NewBlockFromResponse(blockResponse)
	if err != nil {
		blockchain.logger.Error(fmt.Errorf("unable to instantiate block: %w", err).Error())
		return
	}
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.blocks = append(blockchain.blocks, block)
}

func calculateIncome(amount uint64) uint64 {
	return uint64(math.Round(math.Pow(float64(amount), incomeExponent)))
}

func (blockchain *Blockchain) decay(lastTimestamp int64, newTimestamp int64, amount uint64) uint64 {
	elapsedTimestamp := newTimestamp - lastTimestamp
	return uint64(math.Floor(float64(amount) * math.Exp(-blockchain.lambda*float64(elapsedTimestamp))))
}

func (blockchain *Blockchain) isEmpty() bool {
	return blockchain.blocks == nil
}

func removeTarget(targets []string, removedTarget string) []string {
	for i := 0; i < len(targets); i++ {
		if targets[i] == removedTarget {
			targets = append(targets[:i], targets[i+1:]...)
			return targets
		}
	}
	return targets
}

func (blockchain *Blockchain) sortByBlocksLength(selectedTargets []string, blocksByTarget map[string][]*Block) {
	sort.Slice(selectedTargets, func(i, j int) bool {
		return len(blocksByTarget[selectedTargets[i]]) > len(blocksByTarget[selectedTargets[j]])
	})
}

func (blockchain *Blockchain) verify(neighborBlocks []*network.BlockResponse, hostBlocks []*Block, oldHostBlocks []*network.BlockResponse, timestamp int64) (validBlocks []*Block, err error) {
	if len(neighborBlocks) < 2 || len(neighborBlocks) < len(hostBlocks) {
		return nil, errors.New("neighbor's blockchain is too short")
	}
	lastNeighborBlock, err := NewBlockFromResponse(neighborBlocks[len(neighborBlocks)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate last neighbor block: %w", err)
	}
	validatorAddress := lastNeighborBlock.ValidatorAddress()
	isValidatorRegistered, err := blockchain.registry.IsRegistered(validatorAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to get validator proof of humanity: %w", err)
	}
	if !isValidatorRegistered {
		return nil, fmt.Errorf("validator address is not registered in Proof of Humanity registry")
	}

	penultimateBlock := neighborBlocks[len(neighborBlocks)-2]
	registeredAddresses := penultimateBlock.RegisteredAddresses
	registeredAddressesMap := make(map[string]bool)
	for _, address := range registeredAddresses {
		registeredAddressesMap[address] = true
	}
	for _, transaction := range lastNeighborBlock.Transactions() {
		if !transaction.IsReward() && transaction.Value() > 0 {
			if _, isRegistered := registeredAddressesMap[transaction.SenderAddress()]; !isRegistered {
				var isPohValid bool
				isPohValid, err = blockchain.registry.IsRegistered(transaction.SenderAddress())
				if err != nil {
					return nil, fmt.Errorf("failed to get proof of humanity: %w", err)
				} else if isPohValid {
					registeredAddressesMap[transaction.SenderAddress()] = true
				}
			}
		}
	}
	if len(registeredAddressesMap) != len(lastNeighborBlock.RegisteredAddresses()) {
		if len(registeredAddressesMap) > len(lastNeighborBlock.RegisteredAddresses()) {
			return nil, fmt.Errorf("a registered address is missing in the neighbor block")
		} else if len(registeredAddressesMap) < len(lastNeighborBlock.RegisteredAddresses()) {
			return nil, fmt.Errorf("a registered address is one too many in the neighbor block")
		}
	}
	for _, address := range lastNeighborBlock.RegisteredAddresses() {
		var isPohValid bool
		isPohValid, err = blockchain.registry.IsRegistered(address)
		if err != nil {
			return nil, fmt.Errorf("failed to get proof of humanity: %w", err)
		} else if !isPohValid {
			return nil, fmt.Errorf("an address is not registered in Proof of Humanity registry")
		}
		if _, isRegistered := registeredAddressesMap[address]; !isRegistered {
			return nil, fmt.Errorf("a registered address is is wrong in the neighbor block")
		}
	}
	previousBlock, err := NewBlockFromResponse(neighborBlocks[0])
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate first neighbor block: %w", err)
	}
	validBlocks = append(validBlocks, previousBlock)
	neighborBlockchain := newBlockchain(
		append(oldHostBlocks, neighborBlocks...),
		blockchain.registry,
		blockchain.validationTimer,
		blockchain.synchronizer,
		blockchain.logger,
	)
	for i := 1; i < len(neighborBlocks); i++ {
		var currentBlock *Block
		currentBlock, err = NewBlockFromResponse(neighborBlocks[i])
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate last neighbor block: %w", err)
		}
		var previousBlockHash [32]byte
		previousBlockHash, err = previousBlock.Hash()
		if err != nil {
			return nil, fmt.Errorf("failed to calculate previous neighbor block hash: %w", err)
		}
		isPreviousHashValid := currentBlock.PreviousHash() == previousBlockHash
		if !isPreviousHashValid {
			return nil, errors.New("a previous neighbor block hash is invalid")
		}
		var isNewBlock bool
		if i >= len(hostBlocks) {
			isNewBlock = true
		} else if len(hostBlocks) > 2 {
			var hostBlockHash [32]byte
			var currentBlockHash [32]byte
			currentBlockHash, err = currentBlock.Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to calculate neighbor block hash: %w", err)
			}
			hostBlockHash, err = hostBlocks[i].Hash()
			if err != nil {
				blockchain.logger.Error(fmt.Errorf("failed to calculate host block hash: %w", err).Error())
			}
			if currentBlockHash != hostBlockHash {
				isNewBlock = true
			}
		}

		if isNewBlock {
			var rewarded bool
			totalTransactionsValueBySenderAddress := make(map[string]uint64)
			currentBlockTimestamp := currentBlock.Timestamp()
			previousBlockTimestamp := previousBlock.Timestamp()
			expectedBlockTimestamp := previousBlockTimestamp + blockchain.validationTimer.Nanoseconds()
			if currentBlockTimestamp != expectedBlockTimestamp {
				blockDate := time.Unix(0, currentBlockTimestamp)
				expectedDate := time.Unix(0, expectedBlockTimestamp)
				return nil, fmt.Errorf("neighbor block timestamp is invalid: block date is %v, expected is %v", blockDate, expectedDate)
			}
			if currentBlockTimestamp > timestamp {
				blockDate := time.Unix(0, currentBlockTimestamp)
				nowDate := time.Unix(0, timestamp)
				return nil, fmt.Errorf("neighbor block timestamp is in the future: block date is %v, now is %v", blockDate, nowDate)
			}
			var reward uint64
			var totalTransactionsFees uint64
			for _, transaction := range currentBlock.Transactions() {
				if transaction.IsReward() {
					// Check that there is only one reward by block
					if rewarded {
						return nil, errors.New("multiple rewards attempt for the same neighbor block")
					}
					rewarded = true
					reward = transaction.Value()
				} else {
					if err = transaction.VerifySignature(); err != nil {
						return nil, fmt.Errorf("neighbor transaction is invalid: %w", err)
					}
					fee := transaction.Fee()
					totalTransactionsValueBySenderAddress[transaction.SenderAddress()] += transaction.Value() + fee
					totalTransactionsFees += fee
					if currentBlockTimestamp+blockchain.validationTimer.Nanoseconds() < transaction.Timestamp() {
						return nil, fmt.Errorf("a neighbor block transaction timestamp is too far in the future, transaction: %v", transaction)
					}
					if i > 0 && transaction.Timestamp() < neighborBlocks[i-1].Timestamp {
						return nil, fmt.Errorf("a neighbor block transaction timestamp is too old, transaction: %v", transaction)
					}
				}
			}
			if !rewarded {
				return nil, errors.New("neighbor block has not been rewarded")
			}
			if reward > totalTransactionsFees {
				return nil, errors.New("neighbor block reward exceeds the consented one")
			}
			for senderAddress, totalTransactionsValue := range totalTransactionsValueBySenderAddress {
				amount := neighborBlockchain.CalculateTotalAmount(currentBlockTimestamp, senderAddress)
				if totalTransactionsValue > amount {
					return nil, errors.New("neighbor block total transactions value exceeds its sender wallet amount")
				}
			}
		}
		validBlocks = append(validBlocks, currentBlock)
		previousBlock = currentBlock
	}
	return validBlocks, nil
}
