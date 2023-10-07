package verification

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"sync"
	"time"
)

type Blockchain struct {
	blocks              []*Block
	mutex               sync.RWMutex
	registeredAddresses map[string]bool
	registry            protocol.Registry
	synchronizer        network.Synchronizer
	settings            protocol.Settings
	utxosByAddress      map[string][]*Utxo
	utxosById           map[string][]*Utxo
	logger              log.Logger
}

func NewBlockchain(registry protocol.Registry, settings protocol.Settings, synchronizer network.Synchronizer, logger log.Logger) *Blockchain {
	utxosByAddress := make(map[string][]*Utxo)
	utxosById := make(map[string][]*Utxo)
	registeredAddresses := make(map[string]bool)
	blockchain := newBlockchain(nil, registeredAddresses, registry, settings, synchronizer, utxosByAddress, utxosById, logger)
	return blockchain
}

func newBlockchain(blocks []*Block, registeredAddresses map[string]bool, registry protocol.Registry, settings protocol.Settings, synchronizer network.Synchronizer, utxosByAddress map[string][]*Utxo, utxosById map[string][]*Utxo, logger log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.blocks = blocks
	blockchain.registeredAddresses = registeredAddresses
	blockchain.registry = registry
	blockchain.settings = settings
	blockchain.synchronizer = synchronizer
	blockchain.utxosByAddress = utxosByAddress
	blockchain.utxosById = utxosById
	blockchain.logger = logger
	return blockchain
}

func (blockchain *Blockchain) AddBlock(timestamp int64, transactionsBytes []byte, newAddresses []string) error {
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
	var addedRegisteredAddresses []string
	var removedRegisteredAddresses []string
	for address := range blockchain.registeredAddresses {
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
	for _, address := range newAddresses {
		isPohValid, err := blockchain.registry.IsRegistered(address)
		if err != nil {
			return fmt.Errorf("failed to get proof of humanity: %w", err)
		}
		if isPohValid {
			addedRegisteredAddresses = append(addedRegisteredAddresses, address)
		}
	}
	var transactions []*Transaction
	if transactionsBytes != nil {
		err := json.Unmarshal(transactionsBytes, &transactions)
		if err != nil {
			return fmt.Errorf("failed to unmarshal transactions: %w", err)
		}
	}
	block := NewBlock(timestamp, previousHash, transactions, addedRegisteredAddresses, removedRegisteredAddresses)
	return blockchain.addBlock(block)
}

func (blockchain *Blockchain) Blocks(startingBlockHeight uint64) []byte {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	var endingBlockHeight uint64
	blocksCountLimit := blockchain.settings.BlocksCountLimit()
	blocksCount := len(blockchain.blocks)
	if blockchain.isEmpty() || startingBlockHeight > uint64(blocksCount)-1 || blocksCountLimit == 0 {
		return marshalledEmptyArray()
	} else if startingBlockHeight+blocksCountLimit < uint64(blocksCount) {
		endingBlockHeight = startingBlockHeight + blocksCountLimit
	} else {
		endingBlockHeight = uint64(blocksCount)
	}
	blocks := blockchain.blocks[startingBlockHeight:endingBlockHeight]
	blocksBytes, err := json.Marshal(blocks)
	if err != nil {
		blockchain.logger.Error(err.Error())
		return marshalledEmptyArray()
	}
	return blocksBytes
}

func (blockchain *Blockchain) Copy() protocol.Blockchain {
	blockchain.mutex.RLock()
	defer blockchain.mutex.RUnlock()
	blocks := make([]*Block, len(blockchain.blocks))
	copy(blocks, blockchain.blocks)
	utxosByAddress := copyUtxosMap(blockchain.utxosByAddress)
	utxosById := copyUtxosMap(blockchain.utxosById)
	registeredAddresses := copyRegisteredAddressesMap(blockchain.registeredAddresses)
	blockchainCopy := newBlockchain(blocks, registeredAddresses, blockchain.registry, blockchain.settings, blockchain.synchronizer, utxosByAddress, utxosById, blockchain.logger)
	return blockchainCopy
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

func (blockchain *Blockchain) Update(timestamp int64) {
	// Verify neighbor blockchains
	neighbors := blockchain.synchronizer.Neighbors()
	blocksByTarget := make(map[string][]*Block)
	hostBlocks := blockchain.blocks
	var waitGroup sync.WaitGroup
	var mutex sync.RWMutex
	if len(hostBlocks) > 2 {
		hostTarget := "host"
		blocksByTarget[hostTarget] = hostBlocks
		oldHostBlocks := make([]*Block, len(hostBlocks)-1)
		copy(oldHostBlocks, hostBlocks[:len(hostBlocks)-1])
		lastHostBlocks := []*Block{hostBlocks[len(hostBlocks)-1]}
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
	var selectedBlocks []*Block
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
	if isDifferent && len(selectedBlocks) != 0 {
		blockchain.mutex.Lock()
		defer blockchain.mutex.Unlock()
		var newBlocks []*Block
		if isFork {
			blockchain.registeredAddresses = make(map[string]bool)
			blockchain.utxosById = make(map[string][]*Utxo)
			blockchain.utxosByAddress = make(map[string][]*Utxo)
			newBlocks = selectedBlocks[:len(selectedBlocks)-1]
		} else if len(hostBlocks) < len(selectedBlocks) {
			newBlocks = selectedBlocks[len(hostBlocks)-1 : len(selectedBlocks)-1]
		}
		var err error
		for _, newBlock := range newBlocks {
			err = blockchain.updateUtxos(newBlock, newBlock.Timestamp())
			if err != nil {
				blockchain.logger.Error(fmt.Errorf("verification failed: failed to add UTXO: %w", err).Error())
				blockchain.logger.Debug("verification done: blockchain kept")
				return
			}
			blockchain.updateRegisteredAddresses(newBlock.AddedRegisteredAddresses(), newBlock.RemovedRegisteredAddresses())
		}
		blockchain.blocks = selectedBlocks
		blockchain.logger.Debug("verification done: blockchain replaced")
	} else {
		blockchain.logger.Debug("verification done: blockchain kept")
	}
}

func (blockchain *Blockchain) Utxo(input protocol.InputInfo) (protocol.Utxo, error) {
	utxos, ok := blockchain.utxosById[input.TransactionId()]
	if !ok || int(input.OutputIndex()) > len(utxos)-1 {
		return nil, fmt.Errorf("failed to find UTXO, input: %v", input)
	}
	utxo := utxos[input.OutputIndex()]
	if utxo == nil {
		return nil, fmt.Errorf("failed to find UTXO, input: %v", input)
	}
	return utxo, nil
}

func (blockchain *Blockchain) Utxos(address string) []byte {
	utxos, ok := blockchain.utxosByAddress[address]
	if !ok {
		return marshalledEmptyArray()
	}
	marshaledUtxos, err := json.Marshal(utxos)
	if err != nil {
		blockchain.logger.Error(err.Error())
		return marshalledEmptyArray()
	}
	return marshaledUtxos
}

func (blockchain *Blockchain) addBlock(block *Block) error {
	if !blockchain.isEmpty() {
		lastBlock := blockchain.blocks[len(blockchain.blocks)-1]
		err := blockchain.updateUtxos(lastBlock, lastBlock.Timestamp())
		if err != nil {
			return fmt.Errorf("failed to add UTXO: %w", err)
		}
		blockchain.updateRegisteredAddresses(lastBlock.AddedRegisteredAddresses(), lastBlock.RemovedRegisteredAddresses())
	}
	blockchain.blocks = append(blockchain.blocks, block)
	return nil
}

func (blockchain *Blockchain) isEmpty() bool {
	return len(blockchain.blocks) == 0
}

func (blockchain *Blockchain) isRegistered(address string, addedRegisteredAddresses []string, removedRegisteredAddresses []string) error {
	var isRegistered bool
	for _, addedAddress := range addedRegisteredAddresses {
		isRegistered = addedAddress == address
		if isRegistered {
			break
		}
	}
	if !isRegistered {
		for _, addedAddress := range removedRegisteredAddresses {
			isRegistered = addedAddress != address
			if !isRegistered {
				break
			}
		}
		if !isRegistered {
			if _, ok := blockchain.registeredAddresses[address]; !ok {
				return fmt.Errorf("an incomed output address is not registered")
			}
		}
	}
	return nil
}

func (blockchain *Blockchain) updateRegisteredAddresses(addedRegisteredAddresses []string, removedRegisteredAddresses []string) {
	for _, address := range removedRegisteredAddresses {
		delete(blockchain.registeredAddresses, address)
	}
	for _, address := range addedRegisteredAddresses {
		blockchain.registeredAddresses[address] = true
	}
}

func (blockchain *Blockchain) updateUtxos(block *Block, timestamp int64) error {
	utxosByAddress := copyUtxosMap(blockchain.utxosByAddress)
	utxosById := copyUtxosMap(blockchain.utxosById)
	for _, transaction := range block.Transactions() {
		utxosForTransactionId, ok := utxosById[transaction.Id()]
		if ok {
			return fmt.Errorf("transaction ID already exists: %s", transaction.Id())
		}
		for j, output := range transaction.Outputs() {
			if output.InitialValue() > 0 {
				inputInfo := NewInputInfo(uint16(j), transaction.Id())
				utxo := NewUtxo(inputInfo, output, timestamp)
				utxosForTransactionId = append(utxosForTransactionId, utxo)
				utxosById[transaction.Id()] = utxosForTransactionId
				utxosByAddress[output.Address()] = append(utxosByAddress[output.Address()], utxo)
			}
		}
		for _, input := range transaction.Inputs() {
			utxosForInputTransactionId := utxosById[input.TransactionId()]
			if int(input.OutputIndex()) > len(utxosForInputTransactionId)-1 {
				return fmt.Errorf("failed to find UTXO, input: %v", input)
			}
			utxo := utxosForInputTransactionId[input.OutputIndex()]
			if utxo == nil {
				return fmt.Errorf("failed to find output index, input: %v", input)
			}
			utxosForUtxoAddress := utxosByAddress[utxo.Address()]
			utxosForUtxoAddress = removeUtxo(utxosForUtxoAddress, input.TransactionId(), input.OutputIndex())
			utxosByAddress[utxo.Address()] = utxosForUtxoAddress
			utxosById[input.TransactionId()][input.OutputIndex()] = nil
			isEmpty := true
			for _, output := range utxosForInputTransactionId {
				if output != nil {
					isEmpty = false
					break
				}
			}
			if isEmpty {
				delete(utxosById, input.TransactionId())
			}
			if len(utxosForUtxoAddress) == 0 {
				delete(utxosByAddress, utxo.Address())
			}
		}
	}
	if err := verifyIncomes(utxosByAddress); err != nil {
		return err
	}
	blockchain.utxosById = utxosById
	blockchain.utxosByAddress = utxosByAddress
	return nil
}

func (blockchain *Blockchain) verify(lastHostBlocks []*Block, neighborBlocks []*Block, oldHostBlocks []*Block, timestamp int64) ([]*Block, error) {
	if len(oldHostBlocks) == 0 && len(neighborBlocks) < 2 {
		return nil, errors.New("neighbor's blockchain is too short")
	} else if len(oldHostBlocks) > 0 && (len(neighborBlocks) == 0 || lastHostBlocks[0].PreviousHash() != neighborBlocks[0].PreviousHash()) {
		return nil, errors.New("neighbor's blockchain is a fork")
	}
	if neighborBlocks[len(neighborBlocks)-1].Timestamp() == timestamp {
		err := blockchain.verifyLastBlock(neighborBlocks)
		if err != nil {
			return nil, err
		}
	}
	var verifiedBlocks []*Block
	var utxosByAddress map[string][]*Utxo
	var utxosById map[string][]*Utxo
	var registeredAddresses map[string]bool
	if len(oldHostBlocks) == 0 {
		utxosByAddress = make(map[string][]*Utxo)
		utxosById = make(map[string][]*Utxo)
		registeredAddresses = make(map[string]bool)
	} else {
		utxosByAddress = copyUtxosMap(blockchain.utxosByAddress)
		utxosById = copyUtxosMap(blockchain.utxosById)
		registeredAddresses = copyRegisteredAddressesMap(blockchain.registeredAddresses)
	}
	neighborBlockchain := newBlockchain(oldHostBlocks, registeredAddresses, blockchain.registry, blockchain.settings, blockchain.synchronizer, utxosByAddress, utxosById, blockchain.logger)
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
			if err := blockchain.verifyBlock(neighborBlock, previousBlockTimestamp, timestamp, neighborBlockchain); err != nil {
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
	if err := blockchain.verifyRegisteredAddresses(lastNeighborBlock); err != nil {
		return nil, fmt.Errorf("failed to verify registered addresses: %w", err)
	}
	currentBlockTimestamp := lastNeighborBlock.Timestamp()
	nextBlockTimestamp := currentBlockTimestamp + blockchain.settings.ValidationTimestamp()
	if err := neighborBlockchain.AddBlock(nextBlockTimestamp, nil, nil); err != nil {
		return nil, err
	}
	return verifiedBlocks, nil
}

func (blockchain *Blockchain) verifyNeighborBlockchain(timestamp int64, neighbor network.Neighbor, startingBlockHeight uint64, lastHostBlocks []*Block, oldHostBlocks []*Block) ([]*Block, error) {
	type ChanResult struct {
		Blocks []*Block
		Err    error
	}
	blocksChannel := make(chan *ChanResult)
	go func(neighbor network.Neighbor) {
		defer close(blocksChannel)
		neighborBlocksBytes, err := neighbor.GetBlocks(startingBlockHeight)
		if err != nil {
			blocksChannel <- &ChanResult{Err: fmt.Errorf("failed to get neighbor's blockchain: %w", err)}
		}
		var neighborBlocks []*Block
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

func (blockchain *Blockchain) verifyBlock(neighborBlock *Block, previousBlockTimestamp int64, timestamp int64, neighborBlockchain *Blockchain) error {
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
	removedRegisteredAddresses := neighborBlock.RemovedRegisteredAddresses()
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
			fee, err := transaction.Fee(blockchain.settings, currentBlockTimestamp, neighborBlockchain.Utxo)
			if err != nil {
				return fmt.Errorf("failed to verify a neighbor block transaction fee: %w", err)
			}
			totalTransactionsFees += fee
			if currentBlockTimestamp < transaction.Timestamp() {
				return fmt.Errorf("a neighbor block transaction timestamp is too far in the future: transaction timestamp: %d, id: %s", transaction.Timestamp(), transaction.Id())
			}
			if transaction.Timestamp() < previousBlockTimestamp {
				return fmt.Errorf("a neighbor block transaction timestamp is too old: transaction timestamp: %d, id: %s", transaction.Timestamp(), transaction.Id())
			}
			for _, output := range transaction.Outputs() {
				if output.IsRegistered() {
					if err := blockchain.isRegistered(output.Address(), addedRegisteredAddresses, removedRegisteredAddresses); err != nil {
						return err
					}
				}
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

func (blockchain *Blockchain) verifyLastBlock(lastNeighborBlocks []*Block) error {
	lastNeighborBlock := lastNeighborBlocks[len(lastNeighborBlocks)-1]
	validatorAddress := lastNeighborBlock.ValidatorAddress()
	isValidatorRegistered, err := blockchain.registry.IsRegistered(validatorAddress)
	if err != nil {
		blockchain.logger.Debug(fmt.Errorf("failed to get validator proof of humanity: %w", err).Error())
	} else if !isValidatorRegistered {
		return fmt.Errorf("validator address is not registered in Proof of Humanity registry")
	}
	return nil
}

func (blockchain *Blockchain) verifyRegisteredAddresses(block *Block) error {
	for _, address := range block.RemovedRegisteredAddresses() {
		isPohValid, err := blockchain.registry.IsRegistered(address)
		if err != nil {
			blockchain.logger.Debug(fmt.Errorf("failed to get proof of humanity for address %s: %w", address, err).Error())
		} else if isPohValid {
			return fmt.Errorf("a removed address is registered")
		}
	}
	for _, address := range block.AddedRegisteredAddresses() {
		isPohValid, err := blockchain.registry.IsRegistered(address)
		if err != nil {
			blockchain.logger.Debug(fmt.Errorf("failed to get proof of humanity for address %s: %w", address, err).Error())
		} else if !isPohValid {
			return fmt.Errorf("an added address is not registered")
		}
	}
	return nil
}

func copyRegisteredAddressesMap(registeredAddresses map[string]bool) map[string]bool {
	registeredAddressesCopy := make(map[string]bool, len(registeredAddresses))
	for address := range registeredAddresses {
		registeredAddressesCopy[address] = true
	}
	return registeredAddressesCopy
}

func copyUtxosMap(utxosMap map[string][]*Utxo) map[string][]*Utxo {
	utxosMapCopy := make(map[string][]*Utxo, len(utxosMap))
	for address, utxos := range utxosMap {
		utxosCopy := make([]*Utxo, len(utxos))
		copy(utxosCopy, utxos)
		utxosMapCopy[address] = utxosCopy
	}
	return utxosMapCopy
}

func marshalledEmptyArray() []byte {
	return []byte{91, 93}
}

func removeUtxo(utxos []*Utxo, transactionId string, outputIndex uint16) []*Utxo {
	for i := 0; i < len(utxos); i++ {
		if utxos[i].TransactionId() == transactionId && utxos[i].OutputIndex() == outputIndex {
			utxos = append(utxos[:i], utxos[i+1:]...)
			return utxos
		}
	}
	return utxos
}

func verifyIncomes(utxosByAddress map[string][]*Utxo) error {
	for address, utxos := range utxosByAddress {
		var hasIncome bool
		for _, utxo := range utxos {
			if utxo.IsRegistered() {
				if hasIncome {
					return fmt.Errorf("income requested for several UTXOs for address: %s", address)
				}
				hasIncome = true
			}
		}
	}
	return nil
}
