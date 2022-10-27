package protocol

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"math"
	"sort"
	"sync"
	"time"
)

const (
	incomeExponent                  = 0.54692829
	halfLifeInDays                  = 373.59
	verificationsCountPerValidation = 6
)

type Blockchain struct {
	blocks         []*Block
	blockResponses []*neighborhood.BlockResponse
	mutex          sync.RWMutex

	registrable Registry

	validationTimer time.Duration
	timeable        clock.Time

	synchronizable Synchronizer

	lambda float64

	logger     *log.Logger
	isReplaced bool
}

func NewBlockchain(registrable Registry, validationTimer time.Duration, timeable clock.Time, synchronizable Synchronizer, logger *log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.registrable = registrable
	blockchain.validationTimer = validationTimer
	blockchain.timeable = timeable
	blockchain.synchronizable = synchronizable
	blockchain.logger = logger
	const hoursADay = 24
	halfLife := halfLifeInDays * hoursADay * float64(time.Hour.Nanoseconds())
	blockchain.lambda = math.Log(2) / halfLife
	return blockchain
}

func (blockchain *Blockchain) AddBlock(blockResponse *neighborhood.BlockResponse) {
	block, err := NewBlockFromResponse(blockResponse)
	if err != nil {
		blockchain.logger.Error(fmt.Errorf("unable to add block: %w", err).Error())
		return
	}
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.blocks = append(blockchain.blocks, block)
}

func (blockchain *Blockchain) IsEmpty() bool {
	return blockchain.blocks == nil
}

func (blockchain *Blockchain) Blocks() []*neighborhood.BlockResponse {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	return blockchain.blockResponses
}

func (blockchain *Blockchain) getValidBlocks(neighborBlocks []*neighborhood.BlockResponse) (validBlocks []*Block, err error) {
	if len(neighborBlocks) < 2 || len(neighborBlocks) < len(blockchain.blocks) {
		return nil, errors.New("neighbor's blockchain is too short")
	}
	lastNeighborBlock, err := NewBlockFromResponse(neighborBlocks[len(neighborBlocks)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate last neighbor block: %w", err)
	}
	validatorAddress := lastNeighborBlock.ValidatorAddress()
	isValidatorRegistered, err := blockchain.registrable.IsRegistered(validatorAddress)
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
				isPohValid, err = blockchain.registrable.IsRegistered(transaction.SenderAddress())
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
		isPohValid, err = blockchain.registrable.IsRegistered(address)
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
	now := blockchain.timeable.Now().UnixNano()
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
		if i >= len(blockchain.blocks) {
			isNewBlock = true
		} else if len(blockchain.blocks) > 2 {
			var hostBlockHash [32]byte
			var currentBlockHash [32]byte
			currentBlockHash, err = currentBlock.Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to calculate neighbor block hash: %w", err)
			}
			hostBlockHash, err = blockchain.blocks[i].Hash()
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
			var reward uint64
			var totalTransactionsFees uint64
			for _, transaction := range currentBlock.Transactions() {
				if transaction.IsReward() {
					// Check that there is only one reward by block
					if rewarded {
						return nil, errors.New("multiple rewards attempt for the same neighbor block")
					}
					rewarded = true
					previousBlockTimestamp := previousBlock.Timestamp()
					if currentBlockTimestamp != previousBlockTimestamp+blockchain.validationTimer.Nanoseconds() || currentBlockTimestamp > now {
						return nil, errors.New("neighbor block reward timestamp is invalid")
					}
					reward = transaction.Value()
				} else {
					if err = transaction.VerifySignature(); err != nil {
						return nil, fmt.Errorf("neighbor transaction is invalid: %w", err)
					}
					fee := transaction.Fee()
					totalTransactionsValueBySenderAddress[transaction.SenderAddress()] += transaction.Value() + fee
					totalTransactionsFees += fee
				}
			}
			if !rewarded {
				return nil, errors.New("neighbor block has not been rewarded")
			}
			if reward > totalTransactionsFees {
				return nil, errors.New("neighbor block reward exceeds the consented one")
			}
			for senderAddress, totalTransactionsValue := range totalTransactionsValueBySenderAddress {
				if totalTransactionsValue > blockchain.calculateTotalAmount(currentBlockTimestamp, senderAddress, validBlocks) {
					return nil, errors.New("neighbor block total transactions value exceeds its sender wallet amount")
				}
			}
		}
		validBlocks = append(validBlocks, currentBlock)
		previousBlock = currentBlock
	}
	return validBlocks, nil
}

func (blockchain *Blockchain) Verify() {
	// Select valid blocks
	neighbors := blockchain.synchronizable.Neighbors()
	blockchain.mutex.RLock()
	blockResponsesByTarget := make(map[string][]*neighborhood.BlockResponse)
	blocksByTarget := make(map[string][]*Block)
	var selectedTargets []string
	for _, neighbor := range neighbors {
		neighborBlocks, err := neighbor.GetBlocks()
		target := neighbor.Target()
		blockResponsesByTarget[target] = neighborBlocks
		if err == nil {
			var validBlocks []*Block
			validBlocks, err = blockchain.getValidBlocks(neighborBlocks)
			if err != nil || validBlocks == nil {
				blockchain.logger.Debug(fmt.Errorf("failed to verify blocks for neighbor %s: %w", target, err).Error())
			} else {
				blocksByTarget[target] = validBlocks
				selectedTargets = append(selectedTargets, target)
			}
		}
	}

	if len(blockchain.blocks) > 2 {
		hostTarget := "host"
		blockResponsesByTarget[hostTarget] = blockchain.blockResponses
		blocksByTarget[hostTarget] = blockchain.blocks
		selectedTargets = append(selectedTargets, hostTarget)
	}

	var selectedBlocksResponse []*neighborhood.BlockResponse
	var selectedBlocks []*Block
	var isReplaced bool
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
		// Select the oldest reward recipient's blockchain
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
				selectedBlocksResponse = blockResponsesByTarget[target]
				selectedBlocks = blocks
			}
		}
		// Check if blockchain is replaced
		if len(blockchain.blocks) < 2 && selectedBlocks != nil || len(blockchain.blocks) < len(selectedBlocks) {
			isReplaced = true
		} else if len(selectedBlocks) >= 2 {
			lastNewBlockHash, newBlockHashError := selectedBlocks[len(selectedBlocks)-1].Hash()
			if newBlockHashError != nil {
				blockchain.logger.Error("failed to calculate new block hash")
			} else {
				lastOldBlockHash, oldBlockHashError := blockchain.blocks[len(blockchain.blocks)-1].Hash()
				if oldBlockHashError != nil {
					blockchain.logger.Error("failed to calculate old block hash")
					isReplaced = true
				} else {
					isReplaced = lastOldBlockHash != lastNewBlockHash
				}
			}
		}
	}
	blockchain.mutex.RUnlock()
	if isReplaced {
		blockchain.mutex.Lock()
		blockchain.blockResponses = selectedBlocksResponse
		blockchain.blocks = selectedBlocks
		blockchain.mutex.Unlock()
		blockchain.logger.Debug("verification done: blockchain replaced")
	} else {
		blockchain.logger.Debug("verification done: blockchain kept")
	}
	blockchain.isReplaced = isReplaced
}

func (blockchain *Blockchain) IsReplaced() bool {
	return blockchain.isReplaced
}

func (blockchain *Blockchain) sortByBlocksLength(selectedTargets []string, blocksByTarget map[string][]*Block) {
	sort.Slice(selectedTargets, func(i, j int) bool {
		return len(blocksByTarget[selectedTargets[i]]) > len(blocksByTarget[selectedTargets[j]])
	})
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

func removeTransactions(transactions []*Transaction, removedTransaction *Transaction) []*Transaction {
	for i := 0; i < len(transactions); i++ {
		if transactions[i] == removedTransaction {
			transactions = append(transactions[:i], transactions[i+1:]...)
			return transactions
		}
	}
	return transactions
}

func (blockchain *Blockchain) CalculateTotalAmount(currentTimestamp int64, blockchainAddress string) uint64 {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	return blockchain.calculateTotalAmount(currentTimestamp, blockchainAddress, blockchain.blocks)
}

func (blockchain *Blockchain) calculateTotalAmount(currentTimestamp int64, blockchainAddress string, blocks []*Block) uint64 {
	var totalAmount uint64
	var lastTimestamp int64
	for _, block := range blocks {
		for _, registeredAddress := range block.RegisteredAddresses() {
			if blockchainAddress == registeredAddress {
				if totalAmount > 0 {
					totalAmount = blockchain.decay(lastTimestamp, block.Timestamp(), totalAmount)
					totalAmount += calculateIncome(totalAmount)
					lastTimestamp = block.Timestamp()
				}
				break
			}
		}
		for _, transaction := range block.Transactions() {
			value := transaction.Value()
			if blockchainAddress == transaction.RecipientAddress() {
				if totalAmount > 0 {
					totalAmount = blockchain.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				totalAmount += value
				lastTimestamp = block.Timestamp()
			} else if blockchainAddress == transaction.SenderAddress() {
				if totalAmount > 0 {
					totalAmount = blockchain.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				if totalAmount < value+transaction.Fee() {
					blockchain.logger.Error(fmt.Sprintf("historical transaction have not been properly validated: wallet amount=%d, transaction value=%d", totalAmount, value))
					totalAmount = 0
				} else {
					totalAmount -= value + transaction.Fee()
				}
				lastTimestamp = block.Timestamp()
			}
		}
	}
	return blockchain.decay(lastTimestamp, currentTimestamp, totalAmount)
}

func calculateIncome(amount uint64) uint64 {
	return uint64(math.Round(math.Pow(float64(amount), incomeExponent)))
}

func (blockchain *Blockchain) decay(lastTimestamp int64, newTimestamp int64, amount uint64) uint64 {
	elapsedTimestamp := newTimestamp - lastTimestamp
	return uint64(math.Floor(float64(amount) * math.Exp(-blockchain.lambda*float64(elapsedTimestamp))))
}

func (blockchain *Blockchain) StartVerification() {
	ticker := time.NewTicker(blockchain.validationTimer / verificationsCountPerValidation)
	go func() {
		for {
			for i := 0; i < verificationsCountPerValidation; i++ {
				if i > 0 {
					go blockchain.Verify()
				}
				<-ticker.C
			}
		}
	}()
}
