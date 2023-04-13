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

const halfLifeInDays = 373.59

type Blockchain struct {
	blocks                []*Block
	blockResponses        []*network.BlockResponse
	genesisTimestamp      int64
	lambda                float64
	minimalTransactionFee uint64
	mutex                 sync.RWMutex
	registry              protocol.Registry
	synchronizer          network.Synchronizer
	utxosByAddress        map[string][]*network.WalletOutputResponse
	utxosById             map[string][]*network.OutputResponse
	validationTimestamp   int64
	logger                log.Logger
}

func NewBlockchain(genesisTimestamp int64, genesisTransaction *network.TransactionResponse, minimalTransactionFee uint64, registry protocol.Registry, validationTimer time.Duration, synchronizer network.Synchronizer, logger log.Logger) *Blockchain {
	utxosByAddress := make(map[string][]*network.WalletOutputResponse)
	utxosById := make(map[string][]*network.OutputResponse)
	blockchain := newBlockchain(nil, genesisTimestamp, minimalTransactionFee, registry, utxosByAddress, utxosById, validationTimer.Nanoseconds(), synchronizer, logger)
	blockchain.addGenesisBlock(genesisTransaction)
	return blockchain
}

func newBlockchain(blockResponses []*network.BlockResponse, genesisTimestamp int64, minimalTransactionFee uint64, registry protocol.Registry, utxosByAddress map[string][]*network.WalletOutputResponse, utxosById map[string][]*network.OutputResponse, validationTimestamp int64, synchronizer network.Synchronizer, logger log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.blockResponses = blockResponses
	blockchain.genesisTimestamp = genesisTimestamp
	blockchain.minimalTransactionFee = minimalTransactionFee
	blockchain.registry = registry
	blockchain.validationTimestamp = validationTimestamp
	blockchain.synchronizer = synchronizer
	blockchain.utxosByAddress = utxosByAddress
	blockchain.utxosById = utxosById
	const hoursADay = 24
	halfLife := halfLifeInDays * hoursADay * float64(time.Hour.Nanoseconds())
	blockchain.lambda = math.Log(2) / halfLife
	blockchain.logger = logger
	return blockchain
}

func (blockchain *Blockchain) AddBlock(timestamp int64, transactions []*network.TransactionResponse, newAddresses []string) error {
	var previousHash [32]byte
	var err error
	blockchain.mutex.Lock()
	defer blockchain.mutex.Unlock()
	var lastRegisteredAddresses []string
	if !blockchain.isEmpty() {
		previousBlock := blockchain.blocks[len(blockchain.blocks)-1]
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

	blockResponse := NewBlockResponse(timestamp, previousHash, transactions, addedRegisteredAddresses, removedRegisteredAddresses)
	block, err := NewBlockFromResponse(blockResponse, lastRegisteredAddresses)
	if err != nil {
		return fmt.Errorf("unable to instantiate block: %w", err)
	}
	if !blockchain.isEmpty() {
		newBlocks := []*network.BlockResponse{blockchain.blockResponses[len(blockchain.blockResponses)-1]}
		blockchain.addUtxos(newBlocks)
	}
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.blocks = append(blockchain.blocks, block)
	return nil
}

func (blockchain *Blockchain) UtxosByAddress(address string) []*network.WalletOutputResponse {
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
	blockchainCopy := new(Blockchain)
	blockchainCopy.genesisTimestamp = blockchain.genesisTimestamp
	blockchainCopy.registry = blockchain.registry
	blockchainCopy.validationTimestamp = blockchain.validationTimestamp
	blockchainCopy.synchronizer = blockchain.synchronizer
	blockchainCopy.utxosByAddress = blockchain.utxosByAddress
	blockchainCopy.utxosById = blockchain.utxosById
	blockchainCopy.lambda = blockchain.lambda
	blockchainCopy.logger = blockchain.logger
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	blockchainCopy.blocks = blockchain.blocks
	blockchainCopy.blockResponses = blockchain.blockResponses
	return blockchainCopy
}

func (blockchain *Blockchain) Lambda() float64 {
	return blockchain.lambda
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
	hostBlockResponses := blockchain.blockResponses
	var oldHostBlockResponses []*network.BlockResponse
	var oldHostBlocks []*Block
	var lastHostBlocks []*Block
	var lastRegisteredAddresses []string
	if len(hostBlocks) > 2 {
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
			target := neighbor.Target()
			startingBlockHeight := uint64(len(hostBlocks) - 1)
			lastNeighborBlockResponses, err := neighbor.GetLastBlocks(startingBlockHeight)
			if err == nil && lastNeighborBlockResponses != nil {
				verifiedBlocks, err := blockchain.verify(lastHostBlocks, lastNeighborBlockResponses, lastRegisteredAddresses, oldHostBlockResponses, timestamp)
				if err != nil || verifiedBlocks == nil {
					blockchain.logger.Debug(fmt.Errorf("failed to verify blocks for neighbor %s: %w", target, err).Error())
				} else {
					blocksByTarget[target] = append(oldHostBlocks, verifiedBlocks...)
					blockResponsesByTarget[target] = append(oldHostBlockResponses, lastNeighborBlockResponses...)
					selectedTargets = append(selectedTargets, target)
				}
			}
		}
	}
	var isFork bool
	if len(selectedTargets) < 2 && len(neighbors) > 0 {
		isFork = true
		blockchain.logger.Debug("all neighbor blockchains are forks, verifying the whole blockchains")
		for _, neighbor := range neighbors {
			target := neighbor.Target()
			neighborBlockResponses, err := neighbor.GetBlocks()
			if err == nil && neighborBlockResponses != nil {
				verifiedBlocks, err := blockchain.verify(hostBlocks, neighborBlockResponses, nil, nil, timestamp)
				if err != nil || verifiedBlocks == nil {
					blockchain.logger.Debug(fmt.Errorf("failed to verify blocks for neighbor %s: %w", target, err).Error())
				} else {
					blocksByTarget[target] = verifiedBlocks
					blockResponsesByTarget[target] = neighborBlockResponses
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
	if isDifferent {
		blockchain.mutex.Lock()
		defer blockchain.mutex.Unlock()
		if isFork {
			blockchain.utxosById = make(map[string][]*network.OutputResponse)
			blockchain.addUtxos(selectedBlockResponses[:len(selectedBlockResponses)-2])
		} else if len(hostBlocks) < len(selectedBlocks) {
			blockchain.addUtxos(selectedBlockResponses[len(blockchain.blockResponses)-1 : len(selectedBlockResponses)-2])
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

func removeUtxo(utxos []*network.WalletOutputResponse, transactionId string, outputIndex uint16) []*network.WalletOutputResponse {
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

func (blockchain *Blockchain) verify(lastHostBlocks []*Block, lastNeighborBlockResponses []*network.BlockResponse, lastRegisteredAddresses []string, oldHostBlockResponses []*network.BlockResponse, timestamp int64) ([]*Block, error) {
	// TODO verify double spend
	if len(lastNeighborBlockResponses) < len(lastHostBlocks) {
		return nil, errors.New("neighbor's blockchain is too short")
	}
	err := blockchain.verifyLastBlock(lastHostBlocks, lastNeighborBlockResponses)
	if err != nil {
		return nil, err
	}
	previousNeighborBlock, err := NewBlockFromResponse(lastNeighborBlockResponses[0], lastRegisteredAddresses)
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate first neighbor block: %w", err)
	}
	verifiedBlocks := []*Block{previousNeighborBlock}
	neighborBlockchain := newBlockchain(
		oldHostBlockResponses,
		blockchain.genesisTimestamp,
		blockchain.minimalTransactionFee,
		blockchain.registry,
		blockchain.utxosByAddress,
		blockchain.utxosById,
		blockchain.validationTimestamp,
		blockchain.synchronizer,
		blockchain.logger,
	)
	for i := 1; i < len(lastNeighborBlockResponses); i++ {
		lastNeighborBlockResponse := lastNeighborBlockResponses[i]
		neighborBlock, err := NewBlockFromResponse(lastNeighborBlockResponse, previousNeighborBlock.RegisteredAddresses())
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
		if i >= len(lastHostBlocks) {
			isNewBlock = true
		} else if len(lastHostBlocks) > 2 {
			neighborBlockHash, err := neighborBlock.Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to calculate neighbor block hash: %w", err)
			}
			hostBlockHash, err := lastHostBlocks[i].Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to calculate host block hash: %w", err)
			}
			if neighborBlockHash != hostBlockHash {
				isNewBlock = true
			}
		}
		err = neighborBlockchain.AddBlock(lastNeighborBlockResponse.Timestamp, lastNeighborBlockResponse.Transactions, lastNeighborBlockResponse.AddedRegisteredAddresses)
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
				return fmt.Errorf("failed to find a neighbor block transaction fee: %w", err)
			}
			if fee < blockchain.minimalTransactionFee {
				return fmt.Errorf("a neighbor block transaction fee is too low, fee: %d, minimal fee: %d", fee, blockchain.minimalTransactionFee)
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
	if lastHostBlocks[0].PreviousHash() != lastNeighborBlockResponses[0].PreviousHash {
		return errors.New("neighbor's blockchain is is a fork")
	}
	lastNeighborBlock, err := NewBlockFromResponse(lastNeighborBlockResponses[len(lastNeighborBlockResponses)-1], lastHostBlocks[0].RegisteredAddresses())
	if err != nil {
		return fmt.Errorf("failed to instantiate last neighbor block: %w", err)
	}
	validatorAddress := lastNeighborBlock.ValidatorAddress()
	isValidatorRegistered, err := blockchain.registry.IsRegistered(validatorAddress)
	if err != nil {
		return fmt.Errorf("failed to get validator proof of humanity: %w", err)
	}
	if !isValidatorRegistered {
		return fmt.Errorf("validator address is not registered in Proof of Humanity registry")
	}
	return nil
}

func (blockchain *Blockchain) addUtxos(blocks []*network.BlockResponse) {
	for _, block := range blocks {
		for _, transaction := range block.Transactions {
			var outputs = make([]*network.OutputResponse, len(transaction.Outputs))
			copy(outputs, transaction.Outputs)
			blockchain.utxosById[transaction.Id] = outputs
			for i, output := range transaction.Outputs {
				if output.Value > 0 {
					walletOutput := &network.WalletOutputResponse{
						Address:       output.Address,
						BlockHeight:   output.BlockHeight,
						HasReward:     output.HasReward,
						HasIncome:     output.HasIncome,
						OutputIndex:   uint16(i),
						TransactionId: transaction.Id,
						Value:         output.Value,
					}
					blockchain.utxosByAddress[output.Address] = append(blockchain.utxosByAddress[output.Address], walletOutput)
				}
			}
			for _, input := range transaction.Inputs {
				utxos := blockchain.utxosById[input.TransactionId]
				if utxos == nil {
					blockchain.logger.Error(fmt.Errorf("failed to find utxo, input: %v", input).Error())
					return
				}
				utxo := utxos[input.OutputIndex]
				if utxo == nil {
					blockchain.logger.Error(fmt.Errorf("failed to find utxo, input: %v", input).Error())
					return
				}
				blockchain.utxosByAddress[utxo.Address] = removeUtxo(blockchain.utxosByAddress[utxo.Address], input.TransactionId, input.OutputIndex)
				blockchain.utxosById[input.TransactionId][input.OutputIndex] = nil
				isEmpty := true
				for _, output := range blockchain.utxosById[input.TransactionId] {
					if output != nil {
						isEmpty = false
					}
				}
				if isEmpty {
					delete(blockchain.utxosById, input.TransactionId)
				}
				if len(blockchain.utxosByAddress[utxo.Address]) == 0 {
					delete(blockchain.utxosByAddress, utxo.Address)
				}
			}
		}
	}
}

func (blockchain *Blockchain) FindFee(transaction *network.TransactionResponse, timestamp int64) (uint64, error) {
	var inputValues uint64
	var outputValues uint64
	for _, input := range transaction.Inputs {
		utxos := blockchain.utxosById[input.TransactionId]
		if utxos == nil {
			return 0, fmt.Errorf("failed to find utxo, input: %v", input)
		}
		utxo := utxos[input.OutputIndex]
		if utxo == nil {
			return 0, fmt.Errorf("failed to find utxo, input: %v", input)
		}
		output := validation.NewOutputFromResponse(utxo, blockchain.lambda, blockchain.validationTimestamp, blockchain.genesisTimestamp)
		inputValues += output.Value(timestamp)
	}
	for _, output := range transaction.Outputs {
		outputValues += output.Value
	}
	if inputValues < outputValues {
		return 0, errors.New("fee is negative")
	}
	return inputValues - outputValues, nil
}
