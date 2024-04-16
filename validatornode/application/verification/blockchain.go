package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"sync"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type Blockchain struct {
	blocks         []*protocol.Block
	mutex          sync.RWMutex
	registry       application.AddressesManager
	sendersManager application.SendersManager
	utxosManager   application.UtxosManager
	settings       application.SettingsProvider
	logger         log.Logger
}

func NewBlockchain(registry application.AddressesManager, settings application.SettingsProvider, sendersManager application.SendersManager, utxosManager application.UtxosManager, logger log.Logger) *Blockchain {
	blockchain := newBlockchain(nil, registry, settings, sendersManager, utxosManager, logger)
	return blockchain
}

func newBlockchain(blocks []*protocol.Block, registry application.AddressesManager, settings application.SettingsProvider, sendersManager application.SendersManager, utxosManager application.UtxosManager, logger log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.blocks = blocks
	blockchain.registry = registry
	blockchain.settings = settings
	blockchain.sendersManager = sendersManager
	blockchain.utxosManager = utxosManager
	blockchain.logger = logger
	return blockchain
}

func (blockchain *Blockchain) AddBlock(timestamp int64, transactions []*protocol.Transaction, newAddresses []string) error {
	blockchain.mutex.Lock()
	defer blockchain.mutex.Unlock()
	var previousHash [32]byte
	if !blockchain.isEmpty() {
		previousBlock := blockchain.blocks[len(blockchain.blocks)-1]
		var err error
		previousHash, err = previousBlock.Hash()
		if err != nil {
			return fmt.Errorf("unable to calculate last block hash: %w", err)
		}
	}
	addedAddresses := blockchain.registry.Filter(newAddresses)
	removedAddresses := blockchain.registry.RemovedAddresses()
	block := protocol.NewBlock(previousHash, addedAddresses, removedAddresses, timestamp, transactions)
	return blockchain.addBlock(block)
}

func (blockchain *Blockchain) Blocks(startingBlockHeight uint64) []*protocol.Block {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	var endingBlockHeight uint64
	blocksCountLimit := blockchain.settings.BlocksCountLimit()
	blocksCount := len(blockchain.blocks)
	if blockchain.isEmpty() || startingBlockHeight > uint64(blocksCount)-1 || blocksCountLimit == 0 {
		return []*protocol.Block{}
	} else if startingBlockHeight+blocksCountLimit < uint64(blocksCount) {
		endingBlockHeight = startingBlockHeight + blocksCountLimit
	} else {
		endingBlockHeight = uint64(blocksCount)
	}
	return blockchain.blocks[startingBlockHeight:endingBlockHeight]
}

func (blockchain *Blockchain) FirstBlockTimestamp() int64 {
	if blockchain.isEmpty() {
		return 0
	} else {
		return blockchain.blocks[0].Timestamp()
	}
}

func (blockchain *Blockchain) LastBlockTimestamp() int64 {
	if blockchain.isEmpty() {
		return 0
	} else {
		return blockchain.blocks[len(blockchain.blocks)-1].Timestamp()
	}
}

func (blockchain *Blockchain) LastBlockTransactions() []*protocol.Transaction {
	if blockchain.isEmpty() {
		return []*protocol.Transaction{}
	} else {
		return blockchain.blocks[len(blockchain.blocks)-1].Transactions()
	}
}

func (blockchain *Blockchain) Update(timestamp int64) {
	// Verify neighbor blockchains
	neighbors := blockchain.sendersManager.Senders()
	blocksByTarget := make(map[string][]*protocol.Block)
	hostBlocks := blockchain.blocks
	var waitGroup sync.WaitGroup
	var mutex sync.RWMutex
	if len(hostBlocks) > 2 {
		hostTarget := "host"
		blocksByTarget[hostTarget] = hostBlocks
		oldHostBlocks := make([]*protocol.Block, len(hostBlocks)-1)
		copy(oldHostBlocks, hostBlocks[:len(hostBlocks)-1])
		lastHostBlocks := []*protocol.Block{hostBlocks[len(hostBlocks)-1]}
		startingBlockHeight := uint64(len(hostBlocks) - 1)
		for _, neighbor := range neighbors {
			waitGroup.Add(1)
			verifiedBlocks, err := blockchain.verifyNeighborBlockchain(timestamp, neighbor, startingBlockHeight, lastHostBlocks, oldHostBlocks)
			target := neighbor.Target()
			if err != nil {
				blockchain.logger.Debug(fmt.Errorf("failed to verify last neighbor blocks for target %s: %w", target, err).Error())
			} else {
				mutex.Lock()
				blocksByTarget[target] = append(oldHostBlocks, verifiedBlocks...)
				mutex.Unlock()
			}
			waitGroup.Done()
		}
		waitGroup.Wait()
	}
	waitGroup.Wait()
	var isFork bool
	if len(hostBlocks) > 0 && len(blocksByTarget) < 2 && len(neighbors) > 0 {
		isFork = true
		blockchain.logger.Debug("all neighbor blockchains are forks, verifying the whole blockchains")
		lastHostBlocks := hostBlocks[:len(hostBlocks)-1]
		var startingBlockHeight uint64 = 0
		for _, neighbor := range neighbors {
			waitGroup.Add(1)
			verifiedBlocks, err := blockchain.verifyNeighborBlockchain(timestamp, neighbor, startingBlockHeight, lastHostBlocks, nil)
			target := neighbor.Target()
			if err != nil {
				blockchain.logger.Debug(fmt.Errorf("failed to verify whole neighbor blocks for target %s: %w", target, err).Error())
			} else {
				mutex.Lock()
				blocksByTarget[target] = verifiedBlocks
				mutex.Unlock()
			}
			waitGroup.Done()
		}
	}
	waitGroup.Wait()
	var selectedBlocks []*protocol.Block
	var isDifferent bool
	if len(blocksByTarget) > 0 {
		// Keep blockchains with consensus for the previous hash (prevent forks)
		minLength := len(hostBlocks)
		maxLength := len(hostBlocks)
		for _, blocks := range blocksByTarget {
			if len(blocks) < minLength {
				minLength = len(blocks)
			}
			if len(blocks) > maxLength {
				maxLength = len(blocks)
			}
		}
		halfNeighborsCount := len(blocksByTarget) / 2
		var rejectedTargets []string
		for target, blocks := range blocksByTarget {
			var samePreviousHashCount int
			for _, otherBlocks := range blocksByTarget {
				if blocks[minLength-1].PreviousHash() == otherBlocks[minLength-1].PreviousHash() {
					samePreviousHashCount++
				}
			}
			if samePreviousHashCount < halfNeighborsCount {
				// The previous hash of the blockchain used to compare is shared by more than 50% neighbors, reject other neighbors
				rejectedTargets = append(rejectedTargets, target)
			}
		}
		for _, rejectedTarget := range rejectedTargets {
			delete(blocksByTarget, rejectedTarget)
		}
		// Keep the longest blockchains
		rejectedTargets = nil
		for target, blocks := range blocksByTarget {
			if len(blocks) < maxLength {
				rejectedTargets = append(rejectedTargets, target)
			}
		}
		for _, rejectedTarget := range rejectedTargets {
			delete(blocksByTarget, rejectedTarget)
		}
		// Select the blockchain of the oldest reward recipient
		var maxRewardRecipientAddressAge uint64
		for _, blocks := range blocksByTarget {
			var rewardRecipientAddressAge uint64
			var lastBlockRewardRecipientAddress string
			for _, transaction := range blocks[len(blocks)-1].Transactions() {
				if transaction.HasReward() {
					lastBlockRewardRecipientAddress = transaction.RewardRecipientAddress()
				}
			}
			var isAgeCalculated bool
			for i := len(blocks) - 2; i >= 0; i-- {
				for _, transaction := range blocks[i].Transactions() {
					if transaction.HasReward() {
						if transaction.RewardRecipientAddress() == lastBlockRewardRecipientAddress {
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
				selectedBlocks = blocks
			}
		}
		// Check if blockchain is different to know if it should be updated
		if len(hostBlocks) < len(selectedBlocks) {
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
	isReplaced := isDifferent && len(selectedBlocks) != 0
	if isReplaced {
		blockchain.mutex.Lock()
		defer blockchain.mutex.Unlock()
		var newBlocks []*protocol.Block
		if isFork {
			blockchain.registry.Clear()
			blockchain.utxosManager.Clear()
			newBlocks = selectedBlocks[:len(selectedBlocks)-1]
		} else if len(hostBlocks) < len(selectedBlocks) {
			newBlocks = selectedBlocks[len(hostBlocks)-1 : len(selectedBlocks)-1]
		}
		for _, newBlock := range newBlocks {
			if err := blockchain.utxosManager.UpdateUtxos(newBlock.Transactions(), newBlock.Timestamp()); err != nil {
				blockchain.logger.Error(fmt.Errorf("verification failed: failed to add UTXO: %w", err).Error())
				isReplaced = false
			} else {
				blockchain.registry.Update(newBlock.AddedRegisteredAddresses(), newBlock.RemovedRegisteredAddresses())
			}
		}
	}
	if isReplaced {
		blockchain.blocks = selectedBlocks
		blockchain.logger.Debug("verification done: blockchain replaced")
	} else {
		blockchain.logger.Debug("verification done: blockchain kept")
	}
}

func (blockchain *Blockchain) addBlock(block *protocol.Block) error {
	if !blockchain.isEmpty() {
		lastBlock := blockchain.blocks[len(blockchain.blocks)-1]
		if err := blockchain.utxosManager.UpdateUtxos(lastBlock.Transactions(), lastBlock.Timestamp()); err != nil {
			return fmt.Errorf("failed to add UTXO: %w", err)
		}
		blockchain.registry.Update(lastBlock.AddedRegisteredAddresses(), lastBlock.RemovedRegisteredAddresses())
	}
	blockchain.blocks = append(blockchain.blocks, block)
	return nil
}

func (blockchain *Blockchain) isEmpty() bool {
	return len(blockchain.blocks) == 0
}

func (blockchain *Blockchain) verify(lastHostBlocks []*protocol.Block, neighborBlocks []*protocol.Block, oldHostBlocks []*protocol.Block, timestamp int64) ([]*protocol.Block, error) {
	if len(oldHostBlocks) == 0 && len(neighborBlocks) < 2 {
		return nil, errors.New("neighbor's blockchain is too short")
	} else if len(oldHostBlocks) > 0 && (len(neighborBlocks) == 0 || lastHostBlocks[0].PreviousHash() != neighborBlocks[0].PreviousHash()) {
		return nil, errors.New("neighbor's blockchain is a fork")
	}
	if neighborBlocks[len(neighborBlocks)-1].Timestamp() == timestamp {
		lastNeighborBlock := neighborBlocks[len(neighborBlocks)-1]
		addedRegisteredAddresses := append(lastNeighborBlock.AddedRegisteredAddresses(), lastNeighborBlock.ValidatorAddress())
		if err := blockchain.registry.Verify(addedRegisteredAddresses, lastNeighborBlock.RemovedRegisteredAddresses()); err != nil {
			return nil, fmt.Errorf("failed to verify registered addresses: %w", err)
		}
	}
	neighborUtxosPool := blockchain.utxosManager.Copy()
	neighborRegistry := blockchain.registry.Copy()
	if len(oldHostBlocks) == 0 {
		neighborUtxosPool.Clear()
		neighborRegistry.Clear()
	}
	neighborBlockchain := newBlockchain(oldHostBlocks, neighborRegistry, blockchain.settings, blockchain.sendersManager, neighborUtxosPool, blockchain.logger)
	var verifiedBlocks []*protocol.Block
	for i := 0; i < len(neighborBlocks); i++ {
		neighborBlock := neighborBlocks[i]
		var previousBlockTimestamp int64
		var previousNeighborBlockHash [32]byte
		var isGenesisBlock bool
		if i == 0 {
			if len(oldHostBlocks) == 0 {
				isGenesisBlock = true
			} else {
				previousNeighborBlock := oldHostBlocks[len(oldHostBlocks)-1]
				previousBlockTimestamp = previousNeighborBlock.Timestamp()
				var err error
				previousNeighborBlockHash, err = previousNeighborBlock.Hash()
				if err != nil {
					return nil, fmt.Errorf("failed to calculate last host block hash: %w", err)
				}
			}
		} else {
			previousNeighborBlock := neighborBlocks[i-1]
			previousBlockTimestamp = previousNeighborBlock.Timestamp()
			var err error
			previousNeighborBlockHash, err = previousNeighborBlock.Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to calculate previous neighbor block hash: %w", err)
			}
		}
		neighborBlockPreviousHash := neighborBlock.PreviousHash()
		isPreviousHashValid := neighborBlockPreviousHash == previousNeighborBlockHash
		if !isPreviousHashValid {
			blockHeight := len(oldHostBlocks) + i
			return nil, fmt.Errorf("a previous neighbor block hash is invalid: block height: %d, block previous hash: %v, previous block hash: %v", blockHeight, neighborBlockPreviousHash, previousNeighborBlockHash)
		}
		var isNewBlock bool
		if len(lastHostBlocks)-1 < i {
			isNewBlock = true
		} else {
			neighborBlockHash, err := neighborBlock.Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to calculate neighbor block hash: %w", err)
			}
			hostBlockHash, err := lastHostBlocks[i].Hash()
			if err != nil {
				blockchain.logger.Error(fmt.Errorf("failed to calculate host block hash: %w", err).Error())
			}
			if neighborBlockHash != hostBlockHash {
				isNewBlock = true
			}
		}
		if isNewBlock && !isGenesisBlock {
			if err := neighborBlockchain.verifyBlock(neighborBlock, previousBlockTimestamp, timestamp); err != nil {
				return nil, err
			}
		}
		if i == 0 {
			neighborBlockchain.blocks = append(neighborBlockchain.blocks, neighborBlock)
		} else if err := neighborBlockchain.addBlock(neighborBlock); err != nil {
			return nil, err
		}
		verifiedBlocks = append(verifiedBlocks, neighborBlock)
	}
	lastNeighborBlock := neighborBlocks[len(neighborBlocks)-1]
	currentBlockTimestamp := lastNeighborBlock.Timestamp()
	nextBlockTimestamp := currentBlockTimestamp + neighborBlockchain.settings.ValidationTimestamp()
	if err := neighborBlockchain.AddBlock(nextBlockTimestamp, nil, nil); err != nil {
		return nil, err
	}
	return verifiedBlocks, nil
}

func (blockchain *Blockchain) verifyNeighborBlockchain(timestamp int64, neighbor application.Sender, startingBlockHeight uint64, lastHostBlocks []*protocol.Block, oldHostBlocks []*protocol.Block) ([]*protocol.Block, error) {
	type ChanResult struct {
		Blocks []*protocol.Block
		Err    error
	}
	blocksChannel := make(chan *ChanResult)
	go func(neighbor application.Sender) {
		defer close(blocksChannel)
		neighborBlocksBytes, err := neighbor.GetBlocks(startingBlockHeight)
		if err != nil {
			blocksChannel <- &ChanResult{Err: fmt.Errorf("failed to get neighbor's blockchain: %w", err)}
		}
		var neighborBlocks []*protocol.Block
		err = json.Unmarshal(neighborBlocksBytes, &neighborBlocks)
		if err != nil {
			blocksChannel <- &ChanResult{Err: fmt.Errorf("failed to get neighbor's blockchain: %w", err)}
		} else {
			blocksChannel <- &ChanResult{Blocks: neighborBlocks}
		}
	}(neighbor)
	timeout := blockchain.settings.ValidationTimeout()
	select {
	case chanResult := <-blocksChannel:
		if chanResult.Err != nil {
			return nil, chanResult.Err
		}
		return blockchain.verify(lastHostBlocks, chanResult.Blocks, oldHostBlocks, timestamp)
	case <-time.After(timeout):
		return nil, errors.New("neighbor's response timeout")
	}
}

func (blockchain *Blockchain) verifyBlock(neighborBlock *protocol.Block, previousBlockTimestamp int64, timestamp int64) error {
	var rewarded bool
	currentBlockTimestamp := neighborBlock.Timestamp()
	expectedBlockTimestamp := previousBlockTimestamp + blockchain.settings.ValidationTimestamp()
	if currentBlockTimestamp != expectedBlockTimestamp {
		blockDate := time.Unix(0, currentBlockTimestamp)
		expectedDate := time.Unix(0, expectedBlockTimestamp)
		return fmt.Errorf("neighbor block timestamp is invalid: block date is %v, expected is %v", blockDate, expectedDate)
	}
	if currentBlockTimestamp > timestamp {
		blockDate := time.Unix(0, currentBlockTimestamp)
		nowDate := time.Unix(0, timestamp)
		return fmt.Errorf("neighbor block timestamp is in the future: block date is %v, now is %v", blockDate, nowDate)
	}
	var reward uint64
	var totalTransactionsFees uint64
	addedRegisteredAddresses := neighborBlock.AddedRegisteredAddresses()
	for _, transaction := range neighborBlock.Transactions() {
		if transaction.HasReward() {
			// Check that there is only one reward by block
			if rewarded {
				return errors.New("multiple rewards attempt for the same neighbor block")
			}
			rewarded = true
			reward = transaction.RewardValue()
		} else {
			if currentBlockTimestamp < transaction.Timestamp() {
				return fmt.Errorf("a neighbor block transaction timestamp is too far in the future: transaction timestamp: %d, id: %s", transaction.Timestamp(), transaction.Id())
			}
			if transaction.Timestamp() < previousBlockTimestamp {
				return fmt.Errorf("a neighbor block transaction timestamp is too old: transaction timestamp: %d, id: %s", transaction.Timestamp(), transaction.Id())
			}
			if err := transaction.VerifySignatures(); err != nil {
				return fmt.Errorf("neighbor transaction is invalid: %w", err)
			}
			for _, output := range transaction.Outputs() {
				if output.IsYielding() {
					var isNewlyRegistered bool
					for _, addedAddress := range addedRegisteredAddresses {
						isNewlyRegistered = addedAddress == output.Address()
						if isNewlyRegistered {
							break
						}
					}
					if !isNewlyRegistered && !blockchain.registry.IsRegistered(output.Address()) {
						return fmt.Errorf("a neighbor block transaction yielding output address is not registered")
					}
				}
			}
			fee, err := blockchain.utxosManager.CalculateFee(transaction, currentBlockTimestamp)
			if err != nil {
				return fmt.Errorf("failed to verify a neighbor block transaction fee: %w", err)
			}
			totalTransactionsFees += fee
		}
	}
	if !rewarded {
		return errors.New("neighbor block has not been rewarded")
	}
	if reward > totalTransactionsFees {
		return errors.New("neighbor block reward exceeds the consented one")
	}
	return nil
}
