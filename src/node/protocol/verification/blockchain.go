package verification

import (
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

type Blockchain struct {
	blocks                []*Block
	blockResponses        []*network.BlockResponse
	genesisTimestamp      int64
	halfLifeInNanoseconds float64
	incomeLimit           uint64
	k                     float64
	minimalTransactionFee uint64
	mutex                 sync.RWMutex
	registry              protocol.Registry
	synchronizer          network.Synchronizer
	utxosByAddress        map[string][]*network.UtxoResponse
	utxosById             map[string][]*network.UtxoResponse
	validationTimestamp   int64
	logger                log.Logger
}

func NewBlockchain(genesisTimestamp int64, genesisTransaction *network.TransactionResponse, halfLifeInNanoseconds float64, incomeBase uint64, incomeLimit uint64, minimalTransactionFee uint64, registry protocol.Registry, validationTimer time.Duration, synchronizer network.Synchronizer, logger log.Logger) *Blockchain {
	k := math.Log(2) / math.Sqrt(-math.Log(1-float64(incomeBase)/float64(incomeLimit)))
	utxosByAddress := make(map[string][]*network.UtxoResponse)
	utxosById := make(map[string][]*network.UtxoResponse)
	blockchain := newBlockchain(nil, genesisTimestamp, halfLifeInNanoseconds, incomeLimit, k, minimalTransactionFee, registry, synchronizer, utxosByAddress, utxosById, validationTimer.Nanoseconds(), logger)
	blockchain.addGenesisBlock(genesisTransaction)
	return blockchain
}

func newBlockchain(blockResponses []*network.BlockResponse, genesisTimestamp int64, halfLifeInNanoseconds float64, incomeLimit uint64, k float64, minimalTransactionFee uint64, registry protocol.Registry, synchronizer network.Synchronizer, utxosByAddress map[string][]*network.UtxoResponse, utxosById map[string][]*network.UtxoResponse, validationTimestamp int64, logger log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.blockResponses = blockResponses
	blockchain.genesisTimestamp = genesisTimestamp
	blockchain.halfLifeInNanoseconds = halfLifeInNanoseconds
	blockchain.incomeLimit = incomeLimit
	blockchain.k = k
	blockchain.minimalTransactionFee = minimalTransactionFee
	blockchain.registry = registry
	blockchain.synchronizer = synchronizer
	blockchain.utxosByAddress = utxosByAddress
	blockchain.utxosById = utxosById
	blockchain.validationTimestamp = validationTimestamp
	blockchain.logger = logger
	return blockchain
}

func (blockchain *Blockchain) AddBlock(timestamp int64, transactions []*network.TransactionResponse, newAddresses []string) error {
	blockchain.mutex.Lock()
	defer blockchain.mutex.Unlock()
	var previousHash [32]byte
	var lastRegisteredAddresses []string
	if !blockchain.isEmpty() {
		previousBlock := blockchain.blocks[len(blockchain.blocks)-1]
		var err error
		previousHash, err = previousBlock.Hash()
		if err != nil {
			return fmt.Errorf("unable to calculate last block hash: %w", err)
		}
		lastRegisteredAddresses = previousBlock.RegisteredAddresses()
	}
	registeredAddressesMap := make(map[string]bool)
	for _, address := range append(lastRegisteredAddresses, newAddresses...) {
		if _, ok := registeredAddressesMap[address]; !ok {
			registeredAddressesMap[address] = false
		}
	}
	var addedRegisteredAddresses []string
	var removedRegisteredAddresses []string
	for address := range registeredAddressesMap {
		isPohValid, err := blockchain.registry.IsRegistered(address)
		if err != nil {
			return fmt.Errorf("failed to get proof of humanity: %w", err)
		}
		if isPohValid {
			addedRegisteredAddresses = append(addedRegisteredAddresses, address)
		} else {
			removedRegisteredAddresses = append(removedRegisteredAddresses, address)
		}
	}
	return blockchain.addBlock(timestamp, transactions, previousHash, addedRegisteredAddresses, removedRegisteredAddresses, lastRegisteredAddresses)
}

func (blockchain *Blockchain) addBlock(timestamp int64, transactions []*network.TransactionResponse, previousHash [32]byte, addedRegisteredAddresses []string, removedRegisteredAddresses []string, lastRegisteredAddresses []string) error {
	if !blockchain.isEmpty() {
		blockHeight := len(blockchain.blockResponses) - 1
		err := blockchain.updateUtxos(blockchain.blockResponses[blockHeight], blockHeight)
		if err != nil {
			return fmt.Errorf("failed to add UTXO: %w", err)
		}
	}
	blockResponse := NewBlockResponse(timestamp, previousHash, transactions, addedRegisteredAddresses, removedRegisteredAddresses)
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	block, err := NewBlockFromResponse(blockResponse, lastRegisteredAddresses)
	if err != nil {
		return fmt.Errorf("unable to instantiate block: %w", err)
	}
	blockchain.blocks = append(blockchain.blocks, block)
	return nil
}

func (blockchain *Blockchain) UtxosByAddress(address string) []*network.UtxoResponse {
	if _, ok := blockchain.utxosByAddress[address]; !ok {
		return nil
	}
	return blockchain.utxosByAddress[address]
}

func (blockchain *Blockchain) Block(blockHeight uint64) *network.BlockResponse {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	if blockHeight > uint64(len(blockchain.blockResponses)-1) {
		return nil
	}
	return blockchain.blockResponses[blockHeight]
}

func (blockchain *Blockchain) Blocks() []*network.BlockResponse {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	return blockchain.blockResponses
}

func (blockchain *Blockchain) Copy() protocol.Blockchain {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	blocks := make([]*Block, len(blockchain.blocks))
	copy(blocks, blockchain.blocks)
	blockResponses := make([]*network.BlockResponse, len(blockchain.blockResponses))
	copy(blockResponses, blockchain.blockResponses)
	utxosByAddress := copyUtxosMap(blockchain.utxosByAddress)
	utxosById := copyUtxosMap(blockchain.utxosById)
	blockchainCopy := newBlockchain(blockResponses, blockchain.genesisTimestamp, blockchain.halfLifeInNanoseconds, blockchain.incomeLimit, blockchain.k, blockchain.minimalTransactionFee, blockchain.registry, blockchain.synchronizer, utxosByAddress, utxosById, blockchain.validationTimestamp, blockchain.logger)
	blockchainCopy.blocks = blocks
	return blockchainCopy
}

func copyUtxosMap(utxosMap map[string][]*network.UtxoResponse) map[string][]*network.UtxoResponse {
	utxosMapCopy := make(map[string][]*network.UtxoResponse, len(utxosMap))
	for address, utxos := range utxosMap {
		utxosCopy := make([]*network.UtxoResponse, len(utxos))
		copy(utxosCopy, utxos)
		utxosMapCopy[address] = utxosCopy
	}
	return utxosMapCopy
}

func (blockchain *Blockchain) LastBlocks(startingBlockHeight uint64) []*network.BlockResponse {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	if startingBlockHeight > uint64(len(blockchain.blockResponses)) {
		return nil
	}
	lastBlocks := make([]*network.BlockResponse, uint64(len(blockchain.blockResponses))-startingBlockHeight)
	copy(lastBlocks, blockchain.blockResponses[startingBlockHeight:])
	return lastBlocks
}

func (blockchain *Blockchain) Update(timestamp int64) {
	// Verify neighbor blockchains
	neighbors := blockchain.synchronizer.Neighbors()
	blockResponsesByTarget := make(map[string][]*network.BlockResponse)
	blocksByTarget := make(map[string][]*Block)
	var selectedTargets []string
	hostBlocks := blockchain.blocks
	var mutex sync.RWMutex
	var waitGroup sync.WaitGroup
	if len(hostBlocks) > 2 {
		hostBlockResponses := blockchain.blockResponses
		var oldHostBlockResponses []*network.BlockResponse
		var oldHostBlocks []*Block
		var lastHostBlocks []*Block
		var lastRegisteredAddresses []string
		hostTarget := "host"
		blockResponsesByTarget[hostTarget] = hostBlockResponses
		blocksByTarget[hostTarget] = hostBlocks
		selectedTargets = append(selectedTargets, hostTarget)
		oldHostBlockResponses = make([]*network.BlockResponse, len(hostBlockResponses)-1)
		oldHostBlocks = make([]*Block, len(hostBlocks)-1)
		lastHostBlocks = []*Block{hostBlocks[len(hostBlocks)-1]}
		copy(oldHostBlockResponses, hostBlockResponses[:len(hostBlockResponses)-1])
		copy(oldHostBlocks, hostBlocks[:len(hostBlocks)-1])
		lastRegisteredAddresses = oldHostBlocks[len(oldHostBlocks)-1].RegisteredAddresses()
		for _, neighbor := range neighbors {
			waitGroup.Add(1)
			go func(neighbor network.Neighbor) {
				target := neighbor.Target()
				startingBlockHeight := uint64(len(hostBlocks) - 1)
				lastNeighborBlockResponses, err := neighbor.GetLastBlocks(startingBlockHeight)
				if err != nil || len(lastNeighborBlockResponses) == 0 || lastHostBlocks[0].PreviousHash() != lastNeighborBlockResponses[0].PreviousHash {
					blockchain.logger.Debug(errors.New("neighbor's blockchain is a fork").Error())
				} else {
					verifiedBlocks, err := blockchain.verify(lastHostBlocks, lastNeighborBlockResponses, lastRegisteredAddresses, oldHostBlockResponses, timestamp)
					if err != nil || verifiedBlocks == nil {
						blockchain.logger.Debug(fmt.Errorf("failed to verify blocks for neighbor %s: %w", target, err).Error())
					} else {
						mutex.Lock()
						blocksByTarget[target] = append(oldHostBlocks, verifiedBlocks...)
						blockResponsesByTarget[target] = append(oldHostBlockResponses, lastNeighborBlockResponses...)
						selectedTargets = append(selectedTargets, target)
						mutex.Unlock()
					}
				}
				waitGroup.Done()
			}(neighbor)
		}
		waitGroup.Wait()
	}
	var isFork bool
	if len(hostBlocks) > 0 && len(selectedTargets) < 2 && len(neighbors) > 0 {
		isFork = true
		blockchain.logger.Debug("all neighbor blockchains are forks, verifying the whole blockchains")
		for _, neighbor := range neighbors {
			waitGroup.Add(1)
			go func(neighbor network.Neighbor) {
				target := neighbor.Target()
				neighborBlockResponses, err := neighbor.GetBlocks()
				if err != nil && len(neighborBlockResponses) < 2 {
					blockchain.logger.Debug(errors.New("neighbor's blockchain is too short").Error())
				} else {
					verifiedBlocks, err := blockchain.verify(hostBlocks, neighborBlockResponses, nil, nil, timestamp)
					if err != nil || verifiedBlocks == nil {
						blockchain.logger.Debug(fmt.Errorf("failed to verify blocks for neighbor %s: %w", target, err).Error())
					} else {
						mutex.Lock()
						blocksByTarget[target] = verifiedBlocks
						blockResponsesByTarget[target] = neighborBlockResponses
						selectedTargets = append(selectedTargets, target)
						mutex.Unlock()
					}
				}
				waitGroup.Done()
			}(neighbor)
		}
	}
	waitGroup.Wait() // FIXME waits indefinitely if a neighbor doesn't respond
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
			if samePreviousHashCount < halfNeighborsCount {
				// The previous hash of the blockchain used to compare is shared by more than 50% neighbors, reject other neighbors
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
				if transaction.HasReward() {
					lastBlockRewardRecipientAddress = transaction.RewardRecipientAddress()
				}
			}
			var isAgeCalculated bool
			for i := len(blocks) - 2; i >= 0; i-- {
				for _, transaction := range blocks[i].transactions {
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
				selectedBlockResponses = blockResponsesByTarget[target]
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
	if isDifferent && selectedBlockResponses != nil {
		blockchain.mutex.Lock()
		defer blockchain.mutex.Unlock()
		var newBlocks []*network.BlockResponse
		if isFork {
			blockchain.utxosById = make(map[string][]*network.UtxoResponse)
			blockchain.utxosByAddress = make(map[string][]*network.UtxoResponse)
			blockchain.genesisTimestamp = selectedBlockResponses[0].Timestamp
			newBlocks = selectedBlockResponses[:len(selectedBlockResponses)-2]
		} else if len(hostBlocks) < len(selectedBlocks) {
			newBlocks = selectedBlockResponses[len(hostBlocks)-1 : len(selectedBlockResponses)-2]
		}
		var err error
		for i, newBlock := range newBlocks {
			err = blockchain.updateUtxos(newBlock, i)
			if err != nil {
				blockchain.logger.Error(fmt.Errorf("verification failed: faild to add UTXO: %w", err).Error())
				blockchain.logger.Debug("verification done: blockchain kept")
				return
			}
		}
		blockchain.blockResponses = selectedBlockResponses
		blockchain.blocks = selectedBlocks
		blockchain.logger.Debug("verification done: blockchain replaced")
	} else {
		blockchain.logger.Debug("verification done: blockchain kept")
	}
}

func (blockchain *Blockchain) addGenesisBlock(genesisTransaction *network.TransactionResponse) {
	var transactions []*network.TransactionResponse
	var addedAddresses []string
	if genesisTransaction != nil {
		transactions = []*network.TransactionResponse{genesisTransaction}
		addedAddresses = []string{genesisTransaction.Outputs[0].Address}
	}
	blockResponse := NewBlockResponse(blockchain.genesisTimestamp, [32]byte{}, transactions, addedAddresses, nil)
	block, err := NewBlockFromResponse(blockResponse, nil)
	if err != nil {
		blockchain.logger.Error(fmt.Errorf("unable to instantiate block: %w", err).Error())
		return
	}
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.blocks = append(blockchain.blocks, block)
}

func (blockchain *Blockchain) isEmpty() bool {
	return blockchain.blocks == nil
}

func removeAddress(addresses []string, removedAddress string) []string {
	for i := 0; i < len(addresses); i++ {
		if addresses[i] == removedAddress {
			addresses = append(addresses[:i], addresses[i+1:]...)
			return addresses
		}
	}
	return addresses
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

func removeUtxo(utxos []*network.UtxoResponse, transactionId string, outputIndex uint16) []*network.UtxoResponse {
	for i := 0; i < len(utxos); i++ {
		if utxos[i].TransactionId == transactionId && utxos[i].OutputIndex == outputIndex {
			utxos = append(utxos[:i], utxos[i+1:]...)
			return utxos
		}
	}
	return utxos
}

func (blockchain *Blockchain) sortByBlocksLength(selectedTargets []string, blocksByTarget map[string][]*Block) {
	sort.Slice(selectedTargets, func(i, j int) bool {
		return len(blocksByTarget[selectedTargets[i]]) > len(blocksByTarget[selectedTargets[j]])
	})
}

func (blockchain *Blockchain) verify(lastHostBlocks []*Block, neighborBlockResponses []*network.BlockResponse, lastRegisteredAddresses []string, oldHostBlockResponses []*network.BlockResponse, timestamp int64) ([]*Block, error) {
	err := blockchain.verifyLastBlock(lastHostBlocks, neighborBlockResponses)
	if err != nil {
		return nil, err
	}
	previousNeighborBlock, err := NewBlockFromResponse(neighborBlockResponses[0], lastRegisteredAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate first neighbor block: %w", err)
	}
	verifiedBlocks := []*Block{previousNeighborBlock}
	var utxosByAddress map[string][]*network.UtxoResponse
	var utxosById map[string][]*network.UtxoResponse
	if oldHostBlockResponses == nil {
		utxosByAddress = make(map[string][]*network.UtxoResponse)
		utxosById = make(map[string][]*network.UtxoResponse)
	} else {
		utxosByAddress = copyUtxosMap(blockchain.utxosByAddress)
		utxosById = copyUtxosMap(blockchain.utxosById)
	}
	neighborBlockchain := newBlockchain(oldHostBlockResponses, blockchain.genesisTimestamp, blockchain.halfLifeInNanoseconds, blockchain.incomeLimit, blockchain.k, blockchain.minimalTransactionFee, blockchain.registry, blockchain.synchronizer, utxosByAddress, utxosById, blockchain.validationTimestamp, blockchain.logger)
	for i := 1; i < len(neighborBlockResponses); i++ {
		neighborBlockResponse := neighborBlockResponses[i]
		neighborBlock, err := NewBlockFromResponse(neighborBlockResponse, previousNeighborBlock.RegisteredAddresses())
		if err != nil {
			return nil, fmt.Errorf("failed to instantiate last neighbor block: %w", err)
		}
		previousNeighborBlockHash, err := previousNeighborBlock.Hash()
		if err != nil {
			return nil, fmt.Errorf("failed to calculate previous neighbor block hash: %w", err)
		}
		neighborBlockPreviousHash := neighborBlock.PreviousHash()
		isPreviousHashValid := neighborBlockPreviousHash == previousNeighborBlockHash
		if !isPreviousHashValid {
			blockHeight := len(oldHostBlockResponses) + i
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
		currentNeighborBlockResponse := neighborBlockResponses[i-1]
		err = neighborBlockchain.addBlock(currentNeighborBlockResponse.Timestamp, currentNeighborBlockResponse.Transactions, currentNeighborBlockResponse.PreviousHash, currentNeighborBlockResponse.AddedRegisteredAddresses, currentNeighborBlockResponse.RemovedRegisteredAddresses, neighborBlock.RegisteredAddresses())
		if err != nil {
			return nil, err
		}
		if isNewBlock {
			err := blockchain.verifyBlock(neighborBlock, previousNeighborBlock, timestamp, neighborBlockchain)
			if err != nil {
				return nil, err
			}
		}
		verifiedBlocks = append(verifiedBlocks, neighborBlock)
		previousNeighborBlock = neighborBlock
	}
	lastNeighborBlockResponse := neighborBlockResponses[len(neighborBlockResponses)-1]
	currentBlockTimestamp := lastNeighborBlockResponse.Timestamp
	err = neighborBlockchain.addBlock(currentBlockTimestamp, lastNeighborBlockResponse.Transactions, lastNeighborBlockResponse.PreviousHash, lastNeighborBlockResponse.AddedRegisteredAddresses, lastNeighborBlockResponse.RemovedRegisteredAddresses, previousNeighborBlock.RegisteredAddresses())
	if err != nil {
		return nil, err
	}
	nextBlockTimestamp := currentBlockTimestamp + blockchain.validationTimestamp
	err = neighborBlockchain.AddBlock(nextBlockTimestamp, nil, nil)
	if err != nil {
		return nil, err
	}
	return verifiedBlocks, nil
}

func (blockchain *Blockchain) verifyBlock(neighborBlock *Block, previousBlock *Block, timestamp int64, neighborBlockchain *Blockchain) error {
	var rewarded bool
	currentBlockTimestamp := neighborBlock.Timestamp()
	previousBlockTimestamp := previousBlock.Timestamp()
	expectedBlockTimestamp := previousBlockTimestamp + blockchain.validationTimestamp
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
	for _, transaction := range neighborBlock.Transactions() {
		if transaction.HasReward() {
			// Check that there is only one reward by block
			if rewarded {
				return errors.New("multiple rewards attempt for the same neighbor block")
			}
			rewarded = true
			reward = transaction.RewardValue()
		} else {
			if err := transaction.VerifySignatures(); err != nil {
				return fmt.Errorf("neighbor transaction is invalid: %w", err)
			}
			fee, err := neighborBlockchain.FindFee(transaction.GetResponse(), neighborBlock.Timestamp())
			if err != nil {
				return fmt.Errorf("failed to verify a neighbor block transaction fee: %w", err)
			}
			totalTransactionsFees += fee
			if currentBlockTimestamp+blockchain.validationTimestamp < transaction.Timestamp() {
				return fmt.Errorf("a neighbor block transaction timestamp is too far in the future, transaction: %v", transaction.GetResponse())
			}
			if transaction.Timestamp() < previousBlock.Timestamp() {
				return fmt.Errorf("a neighbor block transaction timestamp is too old, transaction: %v", transaction.GetResponse())
			}
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

func (blockchain *Blockchain) verifyLastBlock(lastHostBlocks []*Block, lastNeighborBlockResponses []*network.BlockResponse) error {
	lastNeighborBlock, err := NewBlockFromResponse(lastNeighborBlockResponses[len(lastNeighborBlockResponses)-1], lastHostBlocks[0].RegisteredAddresses())
	if err != nil {
		return fmt.Errorf("failed to instantiate last neighbor block: %w", err)
	}
	validatorAddress := lastNeighborBlock.ValidatorAddress()
	isValidatorRegistered, err := blockchain.registry.IsRegistered(validatorAddress)
	if err != nil {
		blockchain.logger.Debug(fmt.Errorf("failed to get validator proof of humanity: %w", err).Error())
	} else if !isValidatorRegistered {
		return fmt.Errorf("validator address is not registered in Proof of Humanity registry")
	}
	return nil
}

func (blockchain *Blockchain) updateUtxos(block *network.BlockResponse, blockHeight int) error {
	utxosByAddress := copyUtxosMap(blockchain.utxosByAddress)
	utxosById := copyUtxosMap(blockchain.utxosById)
	for _, transaction := range block.Transactions {
		if _, ok := utxosById[transaction.Id]; ok {
			return fmt.Errorf("transaction ID already exists: %s", transaction.Id)
		}
		for j, output := range transaction.Outputs {
			if output.Value > 0 {
				utxo := &network.UtxoResponse{
					Address:       output.Address,
					BlockHeight:   blockHeight,
					HasReward:     output.HasReward,
					HasIncome:     output.HasIncome,
					OutputIndex:   uint16(j),
					TransactionId: transaction.Id,
					Value:         output.Value,
				}
				utxosById[transaction.Id] = append(utxosById[transaction.Id], utxo)
				utxosByAddress[output.Address] = append(utxosByAddress[output.Address], utxo)
			}
		}
		for _, input := range transaction.Inputs {
			utxos := utxosById[input.TransactionId]
			if utxos == nil {
				return fmt.Errorf("failed to find transaction ID, input: %v", input)
			}
			utxo := utxos[input.OutputIndex]
			if utxo == nil {
				return fmt.Errorf("failed to find output index, input: %v", input)
			}
			utxosByAddress[utxo.Address] = removeUtxo(utxosByAddress[utxo.Address], input.TransactionId, input.OutputIndex)
			utxosById[input.TransactionId][input.OutputIndex] = nil
			isEmpty := true
			for _, output := range utxosById[input.TransactionId] {
				if output != nil {
					isEmpty = false
				}
			}
			if isEmpty {
				delete(utxosById, input.TransactionId)
			}
			if len(utxosByAddress[utxo.Address]) == 0 {
				delete(utxosByAddress, utxo.Address)
			}
		}
	}
	if err := blockchain.verifyIncomes(utxosByAddress); err != nil {
		return err
	}
	blockchain.utxosById = utxosById
	blockchain.utxosByAddress = utxosByAddress
	return nil
}

func (blockchain *Blockchain) FindFee(transaction *network.TransactionResponse, timestamp int64) (uint64, error) {
	var inputsValue uint64
	var outputsValue uint64
	for _, input := range transaction.Inputs {
		utxos := blockchain.utxosById[input.TransactionId]
		if utxos == nil {
			return 0, fmt.Errorf("failed to find utxo, input: %v", input)
		}
		utxo := utxos[input.OutputIndex]
		if utxo == nil {
			return 0, fmt.Errorf("failed to find utxo, input: %v", input)
		}
		output := validation.NewOutputFromUtxoResponse(utxo, blockchain.genesisTimestamp, blockchain.halfLifeInNanoseconds, blockchain.incomeLimit, blockchain.k, blockchain.validationTimestamp)
		value := output.Value(timestamp)
		inputsValue += value
	}
	for _, output := range transaction.Outputs {
		outputsValue += output.Value
	}
	if inputsValue < outputsValue {
		return 0, errors.New("transaction fee is negative")
	}
	fee := inputsValue - outputsValue
	if fee < blockchain.minimalTransactionFee {
		return 0, fmt.Errorf("transaction fee is too low, fee: %d, minimal fee: %d", fee, blockchain.minimalTransactionFee)
	}
	return fee, nil
}

func (blockchain *Blockchain) verifyIncomes(utxosByAddress map[string][]*network.UtxoResponse) error {
	for address, utxos := range utxosByAddress {
		var hasIncome bool
		for _, utxo := range utxos {
			if utxo.HasIncome {
				if hasIncome {
					return fmt.Errorf("income requested for several UTXOs for address: %s", address)
				}
				hasIncome = true
			}
		}
	}
	return nil
}
