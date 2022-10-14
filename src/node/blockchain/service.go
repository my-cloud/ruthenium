package blockchain

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
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
	IncomeExponent             = 0.54692829
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
	watch             clock.Time

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

func NewService(address string, ip string, port uint16, miningTimer time.Duration, watch clock.Time, logger *log.Logger) *Service {
	service := new(Service)
	service.address = address
	service.ip = ip
	service.port = port
	service.miningTimer = miningTimer
	service.miningTicker = time.NewTicker(miningTimer)
	service.watch = watch
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
	now := service.watch.Now()
	parsedStartDate := now.Truncate(service.miningTimer).Add(service.miningTimer)
	deadline := parsedStartDate.Sub(now)
	service.miningTicker.Reset(deadline)
	<-service.miningTicker.C
	if service.blocks == nil {
		timestamp := parsedStartDate.UnixNano()
		rewardTransaction := NewTransaction(service.address, RewardSenderAddress, nil, timestamp, GenesisAmount)
		transactions := []*Transaction{rewardTransaction}
		block := NewBlock(timestamp, [32]byte{}, transactions, nil)
		service.blocks = append(service.blocks, block)
		blockResponse := block.GetResponse()
		service.blockResponses = append(service.blockResponses, blockResponse)
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
		rand.Seed(service.watch.Now().UnixNano())
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
	service.blocksMutex.Lock()
	defer service.blocksMutex.Unlock()
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
	if err = transaction.VerifySignature(); err != nil {
		err = errors.New("failed to verify transaction")
		return
	}
	var senderWalletAmount uint64
	if len(service.blocks) > 0 {
		senderWalletAmount = service.calculateTotalAmount(service.blocks[len(service.blocks)-1].Timestamp(), transaction.SenderAddress(), service.blocks)
	}
	insufficientBalance := senderWalletAmount < transaction.Value()+transaction.Fee()
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
	startTime := service.watch.Now()
	parsedStartDate := startTime.Truncate(service.miningTimer).Add(service.miningTimer)
	deadline := parsedStartDate.Sub(startTime)
	service.miningTicker.Reset(deadline)
	service.mineRequested = true
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		<-service.miningTicker.C
		service.mine()
		service.mineRequested = false
		if service.miningStarted {
			now := service.watch.Now()
			newParsedStartDate := now.Truncate(service.miningTimer).Add(service.miningTimer)
			newDeadline := newParsedStartDate.Sub(now)
			service.miningTicker.Reset(newDeadline)
		} else {
			service.miningTicker.Stop()
		}
	}()
}

func (service *Service) mine() {
	service.blocksMutex.Lock()
	defer service.blocksMutex.Unlock()
	now := service.watch.Now().Round(service.miningTimer).UnixNano()
	service.addBlock(now)
}

func (service *Service) StartMining() {
	if !service.miningStarted {
		service.miningStarted = true
		go service.mining()
	}
}

func (service *Service) mining() {
	now := service.watch.Now()
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
	service.blocksMutex.RLock()
	defer service.blocksMutex.RUnlock()
	return service.calculateTotalAmount(currentTimestamp, blockchainAddress, service.blocks)
}

func (service *Service) calculateTotalAmount(currentTimestamp int64, blockchainAddress string, blocks []*Block) uint64 {
	var totalAmount uint64
	var lastTimestamp int64
	for _, block := range blocks {
		for _, registeredAddress := range block.RegisteredAddresses() {
			if blockchainAddress == registeredAddress {
				if totalAmount > 0 {
					totalAmount = service.decay(lastTimestamp, block.Timestamp(), totalAmount)
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
					totalAmount = service.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				totalAmount += value
				lastTimestamp = block.Timestamp()
			} else if blockchainAddress == transaction.SenderAddress() {
				if totalAmount > 0 {
					totalAmount = service.decay(lastTimestamp, block.Timestamp(), totalAmount)
				}
				if totalAmount < value+transaction.Fee() {
					service.logger.Error(fmt.Sprintf("historical transaction have not been properly validated: wallet amount=%d, transaction value=%d", totalAmount, value))
					totalAmount = 0
				} else {
					totalAmount -= value + transaction.Fee()
				}
				lastTimestamp = block.Timestamp()
			}
		}
	}
	return service.decay(lastTimestamp, currentTimestamp, totalAmount)
}

func calculateIncome(amount uint64) uint64 {
	return uint64(math.Round(math.Pow(float64(amount), IncomeExponent)))
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

func (service *Service) getValidBlocks(neighborBlocks []*neighborhood.BlockResponse) (validBlocks []*Block, err error) {
	if len(neighborBlocks) < 2 || len(neighborBlocks) < len(service.blocks) {
		return nil, errors.New("neighbor's blockchain is too short")
	}
	lastNeighborBlock, err := NewBlockFromResponse(neighborBlocks[len(neighborBlocks)-1])
	if err != nil {
		return nil, fmt.Errorf("failed to instantiate last neighbor block: %w", err)
	}
	isValidatorPohValid, err := lastNeighborBlock.IsProofOfHumanityValid()
	if err != nil {
		return nil, fmt.Errorf("failed to get validator proof of humanity: %w", err)
	}
	if !isValidatorPohValid {
		return nil, fmt.Errorf("validator address is not registered in Proof of Humanity registry")
	}

	penultimateBlock := neighborBlocks[len(neighborBlocks)-2]
	registeredAddresses := penultimateBlock.RegisteredAddresses
	registeredAddressesMap := make(map[string]bool)
	for _, address := range registeredAddresses {
		registeredAddressesMap[address] = true
	}
	for _, transaction := range lastNeighborBlock.Transactions() {
		if transaction.SenderAddress() != RewardSenderAddress && transaction.Value() > 0 {
			if _, isRegistered := registeredAddressesMap[transaction.SenderAddress()]; !isRegistered {
				var isPohValid bool
				isPohValid, err = NewHuman(transaction.SenderAddress()).IsRegistered()
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
		isPohValid, err = NewHuman(address).IsRegistered()
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
	now := service.watch.Now().UnixNano()
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
		if i >= len(service.blocks) {
			isNewBlock = true
		} else if len(service.blocks) > 2 {
			var hostBlockHash [32]byte
			var currentBlockHash [32]byte
			currentBlockHash, err = currentBlock.Hash()
			if err != nil {
				return nil, fmt.Errorf("failed to calculate neighbor block hash: %w", err)
			}
			hostBlockHash, err = service.blocks[i].Hash()
			if err != nil {
				service.logger.Error(fmt.Errorf("failed to calculate host block hash: %w", err).Error())
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
				if transaction.SenderAddress() == RewardSenderAddress {
					// Check that there is only one reward by block
					if rewarded {
						return nil, errors.New("multiple rewards attempt for the same neighbor block")
					}
					rewarded = true
					previousBlockTimestamp := previousBlock.Timestamp()
					if currentBlockTimestamp != previousBlockTimestamp+int64(service.miningTimer) || currentBlockTimestamp > now {
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
				if totalTransactionsValue > service.calculateTotalAmount(currentBlockTimestamp, senderAddress, validBlocks) {
					return nil, errors.New("neighbor block total transactions value exceeds its sender wallet amount")
				}
			}
		}
		validBlocks = append(validBlocks, currentBlock)
		previousBlock = currentBlock
	}
	return validBlocks, nil
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
				var validBlocks []*Block
				validBlocks, err = service.getValidBlocks(neighborBlocks)
				if err != nil || validBlocks == nil {
					service.logger.Debug(fmt.Errorf("failed to validate blocks for neighbor %s: %w", neighbor.Target(), err).Error())
				} else {
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

		var selectedBlocksResponse []*neighborhood.BlockResponse
		var selectedBlocks []*Block
		if selectedNeighbors != nil {
			// Keep blockchains with consensus for the previous hash (prevent forks)
			service.sortByBlocksLength(selectedNeighbors, blocksByNeighbor)
			halfNeighborsCount := len(blocksByNeighbor) / 2
			minLength := len(blocksByNeighbor[selectedNeighbors[len(selectedNeighbors)-1]])
			var rejectedNeighbors []*neighborhood.Neighbor
			for neighbor, blocks := range blocksByNeighbor {
				var samePreviousHashCount int
				for _, otherBlocks := range blocksByNeighbor {
					if blocks[minLength-1].PreviousHash() == otherBlocks[minLength-1].PreviousHash() {
						samePreviousHashCount++
					}
				}
				if samePreviousHashCount <= halfNeighborsCount {
					// The previous hash of the blockchain used to compare is shared by at least 51% neighbors, reject other neighbors
					rejectedNeighbors = append(rejectedNeighbors, neighbor)
				}
			}
			for _, rejectedNeighbor := range rejectedNeighbors {
				delete(blocksByNeighbor, rejectedNeighbor)
				delete(blockResponsesByNeighbor, rejectedNeighbor)
				removeNeighbor(selectedNeighbors, rejectedNeighbor)
			}
			// Keep the longest blockchains
			maxLength := len(blocksByNeighbor[selectedNeighbors[0]])
			rejectedNeighbors = nil
			for neighbor, blocks := range blocksByNeighbor {
				if len(blocks) < maxLength {
					rejectedNeighbors = append(rejectedNeighbors, neighbor)
				}
			}
			for _, rejectedNeighbor := range rejectedNeighbors {
				delete(blocksByNeighbor, rejectedNeighbor)
				delete(blockResponsesByNeighbor, rejectedNeighbor)
				removeNeighbor(selectedNeighbors, rejectedNeighbor)
			}
			// Select the oldest reward recipient's blockchain
			var maxRewardRecipientAddressAge uint64
			for neighbor, blocks := range blocksByNeighbor {
				var rewardRecipientAddressAge uint64
				var neighborLastBlockRewardRecipientAddress string
				for _, transaction := range blocks[len(blocks)-1].transactions {
					if transaction.SenderAddress() == RewardSenderAddress {
						neighborLastBlockRewardRecipientAddress = transaction.RecipientAddress()
					}
				}
				var isAgeCalculated bool
				for i := len(blocks) - 2; i > 0; i-- {
					for _, transaction := range blocks[i].transactions {
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
			var blockchainReplaced bool
			// Check if blockchain is replaced
			if len(service.blocks) < 2 {
				blockchainReplaced = true
			} else if len(selectedBlocks) >= 2 {
				lastNewBlockHash, newBlockHashError := selectedBlocks[len(selectedBlocks)-1].Hash()
				if newBlockHashError != nil {
					service.logger.Error("failed to calculate new block hash")
				} else {
					lastOldBlockHash, oldBlockHashError := service.blocks[len(service.blocks)-1].Hash()
					if oldBlockHashError != nil {
						service.logger.Error("failed to calculate old block hash")
						blockchainReplaced = true
					} else {
						blockchainReplaced = lastOldBlockHash != lastNewBlockHash
					}
				}
			}
			if blockchainReplaced && selectedBlocks != nil {
				service.transactions = nil
				service.blockResponses = selectedBlocksResponse
				service.blocks = selectedBlocks
				service.logger.Debug("conflicts resolved: blockchain replaced")
			} else {
				service.logger.Debug("conflicts resolved: blockchain kept")
			}
		}
	}()
}

func (service *Service) sortByBlocksLength(selectedNeighbors []*neighborhood.Neighbor, blocksByNeighbor map[*neighborhood.Neighbor][]*Block) {
	sort.Slice(selectedNeighbors, func(i, j int) bool {
		return len(blocksByNeighbor[selectedNeighbors[i]]) > len(blocksByNeighbor[selectedNeighbors[j]])
	})
}

func removeNeighbor(neighbors []*neighborhood.Neighbor, removedNeighbor *neighborhood.Neighbor) []*neighborhood.Neighbor {
	for i := 0; i < len(neighbors); i++ {
		if neighbors[i] == removedNeighbor {
			neighbors = append(neighbors[:i], neighbors[i+1:]...)
			return neighbors
		}
	}
	return neighbors
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

func (service *Service) addBlock(timestamp int64) {
	lastBlock := service.blocks[len(service.blocks)-1]
	if lastBlock.Timestamp() == timestamp {
		service.logger.Error("unable to add block, a block with the same timestamp already is in the blockchain")
		return
	}
	lastBlockHash, err := lastBlock.Hash()
	if err != nil {
		service.logger.Error(fmt.Errorf("failed calculate last block hash: %w", err).Error())
		return
	}
	service.transactionsMutex.Lock()
	defer service.transactionsMutex.Unlock()
	totalTransactionsValueBySenderAddress := make(map[string]uint64)
	transactions := service.transactions
	var reward uint64
	for _, transaction := range transactions {
		fee := transaction.Fee()
		totalTransactionsValueBySenderAddress[transaction.SenderAddress()] += transaction.Value() + fee
		reward += fee
	}
	registeredAddresses := lastBlock.registeredAddresses
	registeredAddressesMap := make(map[string]bool)
	for _, address := range registeredAddresses {
		registeredAddressesMap[address] = true
	}
	for senderAddress, totalTransactionsValue := range totalTransactionsValueBySenderAddress {
		senderTotalAmount := service.calculateTotalAmount(timestamp, senderAddress, service.blocks)
		if totalTransactionsValue > senderTotalAmount {
			var rejectedTransactions []*Transaction
			rand.Seed(service.watch.Now().UnixNano())
			rand.Shuffle(len(transactions), func(i, j int) { transactions[i], transactions[j] = transactions[j], transactions[i] })
			for _, transaction := range transactions {
				if transaction.SenderAddress() == senderAddress {
					rejectedTransactions = append(rejectedTransactions, transaction)
					fee := transaction.Fee()
					totalTransactionsValue -= transaction.Value() + fee
					reward -= fee
					if totalTransactionsValue <= senderTotalAmount {
						break
					}
				}
			}
			for _, transaction := range rejectedTransactions {
				transactions = removeTransactions(transactions, transaction)
			}
			service.logger.Warn("transactions removed from the transactions pool, total transactions value exceeds its sender wallet amount")
		}
		if totalTransactionsValue > 0 {
			if _, isRegistered := registeredAddressesMap[senderAddress]; !isRegistered {
				registeredAddressesMap[senderAddress] = true
			}
		}
	}
	var newRegisteredAddresses []string
	for registeredAddress := range registeredAddressesMap {
		var isPohValid bool
		isPohValid, err = NewHuman(registeredAddress).IsRegistered()
		if err != nil {
			service.logger.Error(fmt.Errorf("failed to get proof of humanity: %w", err).Error())
		} else if isPohValid {
			newRegisteredAddresses = append(newRegisteredAddresses, registeredAddress)
		}
	}
	rewardTransaction := NewTransaction(service.address, RewardSenderAddress, nil, timestamp, reward)
	transactions = append(transactions, rewardTransaction)
	block := NewBlock(timestamp, lastBlockHash, transactions, newRegisteredAddresses)
	service.transactions = nil
	service.blocks = append(service.blocks, block)
	blockResponse := block.GetResponse()
	service.blockResponses = append(service.blockResponses, blockResponse)
	service.logger.Debug(fmt.Sprintf("reward: %d", reward))
}
