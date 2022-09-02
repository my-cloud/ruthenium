package blockchain

import (
	"errors"
	"fmt"
	"gitlab.com/coinsmaster/ruthenium/src/log"
	"gitlab.com/coinsmaster/ruthenium/src/node/encryption"
	"gitlab.com/coinsmaster/ruthenium/src/node/neighborhood"
	"math"
	"net"
	"sort"
	"sync"
	"time"
)

const (
	DefaultPort = 8106

	MiningRewardSenderAddress        = "MINING REWARD SENDER ADDRESS"
	ParticlesCount                   = 100000000
	GenesisAmount             uint64 = 100000 * ParticlesCount
	RewardExponent                   = 1 / 1.828393264
	HalfLifeInDays                   = 373.59

	NeighborSynchronizationTimeInSeconds = 10
	OutboundsCount                       = 6
)

type Service struct {
	transactions          []*Transaction
	transactionsMutex     sync.RWMutex
	blocks                []*Block
	blockResponses        []*neighborhood.BlockResponse
	blocksMutex           sync.RWMutex
	address               string
	mineMutex             sync.Mutex
	miningStarted         bool
	miningStopped         bool
	mineRequested         bool
	miningTimer           time.Duration
	miningTicker          *time.Ticker
	hasMiningTimerChanged bool

	ip        string
	port      uint16
	logger    *log.Logger
	waitGroup *sync.WaitGroup

	neighbors              []*neighborhood.Neighbor // TODO manage max neighbors count (Outbound/Inbound)
	neighborsMutex         sync.RWMutex
	neighborsByTarget      map[string]*neighborhood.Neighbor
	neighborsByTargetMutex sync.RWMutex

	lambda float64
}

func NewService(address string, ip string, port uint16, miningTimer time.Duration, logger *log.Logger) *Service {
	service := new(Service)
	service.miningTimer = miningTimer
	service.miningTicker = time.NewTicker(miningTimer)
	service.address = address
	service.ip = ip
	service.port = port
	service.logger = logger
	var waitGroup sync.WaitGroup
	service.waitGroup = &waitGroup
	now := time.Now().Unix() * time.Second.Nanoseconds()
	genesisTransaction := NewTransaction(now, MiningRewardSenderAddress, address, GenesisAmount)
	service.createBlock(genesisTransaction)
	seedsIps := []string{
		"89.82.76.241",
	}
	service.neighborsByTarget = map[string]*neighborhood.Neighbor{}
	for _, seedIp := range seedsIps {
		seed := neighborhood.NewNeighbor(seedIp, DefaultPort, logger)
		service.neighborsByTarget[seed.Target()] = seed
	}
	const hoursADay = 24
	halfLife := HalfLifeInDays * hoursADay * float64(time.Hour.Nanoseconds())
	service.lambda = math.Log(2) / halfLife
	return service
}

func (service *Service) WaitGroup() *sync.WaitGroup {
	return service.waitGroup
}

func (service *Service) Run() {
	service.StartNeighborsSynchronization()
	service.StartMining()
}

func (service *Service) SynchronizeNeighbors() {
	service.FindNeighbors()
}

func (service *Service) StartNeighborsSynchronization() {
	service.SynchronizeNeighbors()
	_ = time.AfterFunc(time.Second*NeighborSynchronizationTimeInSeconds, service.StartNeighborsSynchronization)
}

func (service *Service) FindNeighbors() {
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
		service.neighborsByTargetMutex.RLock()
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
		service.neighborsByTargetMutex.RUnlock()
		service.neighborsMutex.Lock()
		service.neighbors = neighbors
		service.neighborsMutex.Unlock()
		for _, neighbor := range neighbors {
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
	}(service.neighborsByTarget)
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

func (service *Service) CreateTransaction(transaction *Transaction, publicKey *encryption.PublicKey, signature *encryption.Signature) {
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		err := service.addTransaction(transaction, publicKey, signature)
		if err != nil {
			service.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
			return
		}
		timestamp := transaction.Timestamp()
		senderAddress := transaction.SenderAddress()
		recipientAddress := transaction.RecipientAddress()
		value := transaction.Value()
		publicKeyStr := publicKey.String()
		signatureStr := signature.String()
		var verb = neighborhood.PUT
		transactionRequest := neighborhood.TransactionRequest{
			Verb:             &verb,
			Timestamp:        &timestamp,
			SenderAddress:    &senderAddress,
			RecipientAddress: &recipientAddress,
			SenderPublicKey:  &publicKeyStr,
			Value:            &value,
			Signature:        &signatureStr,
		}
		service.neighborsMutex.RLock()
		defer service.neighborsMutex.RUnlock()
		for _, neighbor := range service.neighbors {
			go func(neighbor *neighborhood.Neighbor) {
				_ = neighbor.AddTransaction(transactionRequest)
			}(neighbor)
		}
	}()
}

func (service *Service) AddTransaction(transaction *Transaction, publicKey *encryption.PublicKey, signature *encryption.Signature) {
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		service.blocksMutex.RLock()
		if len(service.blocks) > 2 {
			for i := len(service.blocks) - 2; i < len(service.blocks); i++ {
				for _, validatedTransaction := range service.blocks[i].transactions {
					if validatedTransaction.Equals(transaction) {
						service.logger.Error("failed to add transaction: the transaction already is in the blockchain")
						return
					}
				}
			}
		}
		service.blocksMutex.RUnlock()

		err := service.addTransaction(transaction, publicKey, signature)
		if err != nil {
			service.logger.Error(fmt.Errorf("failed to add transaction: %w", err).Error())
		}
	}()
}

func (service *Service) addTransaction(transaction *Transaction, publicKey *encryption.PublicKey, signature *encryption.Signature) (err error) {
	marshaledTransaction, err := transaction.MarshalJSON()
	if err != nil {
		err = fmt.Errorf("failed to marshal transaction, %w", err)
		return
	}
	if !signature.Verify(marshaledTransaction, publicKey, transaction.SenderAddress()) {
		err = errors.New("failed to verify transaction")
		return
	}
	senderWalletAmount := service.CalculateTotalAmount(transaction.Timestamp(), transaction.SenderAddress())
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
	service.mineRequested = true
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		<-service.miningTicker.C
		if service.hasMiningTimerChanged {
			service.miningTicker.Reset(service.miningTimer)
			service.hasMiningTimerChanged = false
		}
		service.mine()
		service.mineRequested = false
	}()
}

func (service *Service) mine() {
	service.mineMutex.Lock()
	defer service.mineMutex.Unlock()
	now := time.Now().Unix() * time.Second.Nanoseconds()
	amount := service.CalculateTotalAmount(now, service.address)
	reward := uint64(math.Round(math.Pow(float64(amount), RewardExponent)))
	service.logger.Info(fmt.Sprintf("amount: %d - reward: %d - total: %d", amount, reward, amount+reward))
	rewardTransaction := NewTransaction(now, MiningRewardSenderAddress, service.address, reward)
	service.createBlock(rewardTransaction)
	service.neighborsMutex.RLock()
	defer service.neighborsMutex.RUnlock()
	for _, neighbor := range service.neighbors {
		go func(neighbor *neighborhood.Neighbor) {
			_ = neighbor.Consensus()
		}(neighbor)
	}
}

func (service *Service) StartMining() {
	if !service.miningStarted {
		service.miningStarted = true
		service.miningStopped = false
		go service.mining()
	}
}

func (service *Service) mining() {
	for {
		<-service.miningTicker.C
		if service.hasMiningTimerChanged {
			service.miningTicker.Reset(service.miningTimer)
			service.hasMiningTimerChanged = false
		}
		if service.miningStopped {
			break
		}
		service.mine()
	}
}

func (service *Service) StopMining() {
	service.miningStarted = false
	service.miningStopped = true
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
					totalAmount = service.decay(lastTimestamp, transaction.Timestamp(), totalAmount)
				}
				totalAmount += value
				lastTimestamp = transaction.Timestamp()
			}
			if blockchainAddress == transaction.SenderAddress() {
				if totalAmount > 0 {
					totalAmount = service.decay(lastTimestamp, transaction.Timestamp(), totalAmount)
				}
				if totalAmount < value {
					service.logger.Error(fmt.Sprintf("historical transaction should not have been validated: wallet amount=%d, transaction value=%d", totalAmount, value))
				}
				totalAmount -= value
				lastTimestamp = transaction.Timestamp()
			}
		}
	}
	return service.decay(lastTimestamp, currentTimestamp, totalAmount)
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
	if len(service.blocks) < 3 || len(neighborBlocks) < len(service.blocks) {
		return
	}
	lastNeighborBlock := NewBlockFromResponse(neighborBlocks[len(neighborBlocks)-1])
	if err := lastNeighborBlock.IsProofOfHumanityValid(); err != nil {
		service.logger.Error(fmt.Errorf("failed to get valid proof of humanity: %w", err).Error())
		return
	}
	newBlocks := neighborBlocks[len(service.blockResponses)-2:]
	for _, block := range newBlocks {
		var rewarded bool
		for _, transaction := range block.Transactions {
			if transaction.SenderAddress == MiningRewardSenderAddress {
				// Check that there is only one reward by block
				// FIXME check the reward amount
				if rewarded {
					service.logger.Error("multiple rewards attempt for the same block")
					return
				}
				rewarded = true
			}
		}
	}
	// TODO verify mining reward timestamp
	previousBlock := NewBlockFromResponse(neighborBlocks[0])
	validBlocks = append(validBlocks, previousBlock)
	for i := 1; i < len(neighborBlocks); i++ {
		currentBlock := neighborBlocks[i]
		previousBlockHash, err := previousBlock.Hash()
		if err != nil {
			service.logger.Error(fmt.Errorf("failed to calculate previous block hash: %w", err).Error())
			return
		}
		isPreviousHashValid := currentBlock.PreviousHash == previousBlockHash
		if !isPreviousHashValid {
			service.logger.Info("a hash is invalid for a neighbor")
			return
		}
		if len(service.blocks) > 3 && i == len(service.blocks)-2 {
			var antePenultimateServiceBlockHash [32]byte
			antePenultimateServiceBlockHash, err = service.blocks[len(service.blocks)-3].Hash()
			if err != nil {
				service.logger.Error(fmt.Errorf("failed to calculate antepenultimate block hash: %w", err).Error())
				return
			}
			if previousBlockHash != antePenultimateServiceBlockHash {
				service.logger.Error("blockchain replacement attack")
				return
			}
		}

		previousBlock = NewBlockFromResponse(currentBlock)
		validBlocks = append(validBlocks, previousBlock)
	}
	return
}

func (service *Service) ResolveConflicts() {
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()

		// Select valid blocks
		blockResponsesByNeighbor := make(map[*neighborhood.Neighbor][]*neighborhood.BlockResponse)
		blocksByNeighbor := make(map[*neighborhood.Neighbor][]*Block)
		var selectedNeighbors []*neighborhood.Neighbor
		if len(service.blocks) > 1 {
			host := neighborhood.NewNeighbor(service.ip, service.port, service.logger)
			blockResponsesByNeighbor[host] = service.blockResponses
			blocksByNeighbor[host] = service.blocks
			selectedNeighbors = append(selectedNeighbors, host)
		}
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

		// Select blockchain with consensus for the previous hash
		var selectedBlocksResponse []*neighborhood.BlockResponse
		var selectedBlocks []*Block
		if selectedNeighbors != nil {
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
					// Keep the longest chain (the 2 last blocks might be replaced by blocks that have a hash not shared by 51% neighbors)
					maxLength := len(service.blocks)
					for neighbor, blocks := range blocksByNeighbor {
						if len(blocks) > maxLength {
							maxLength = len(blocks)
							selectedBlocksResponse = blockResponsesByNeighbor[neighbor]
							selectedBlocks = blocks
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
					service.blocksMutex.RLock()
					if len(service.blocks) > 2 {
						oldTransactions := service.transactions
						// Add transactions which are not in the new blocks but the rewards
						for i := len(service.blocks) - 2; i < len(service.blocks); i++ {
							for _, invalidatedTransaction := range service.blocks[i].transactions {
								if invalidatedTransaction.SenderAddress() != MiningRewardSenderAddress {
									oldTransactions = append(oldTransactions, invalidatedTransaction)
								}
							}
						}
						// Reset mining ticker if the host was the last rewarded
						lastBlock := service.blocks[len(service.blocks)-1]
						// TODO avoid escaping error even if the hash already have been calculated and should not return error
						lastSelectedBlockHash, _ := selectedBlocks[len(service.blocks)-1].Hash()
						lastBlockHash, _ := lastBlock.Hash()
						if len(lastBlock.Transactions()) > 0 && lastSelectedBlockHash != lastBlockHash {
							lastTransaction := lastBlock.Transactions()[len(lastBlock.Transactions())-1]
							if lastTransaction.SenderAddress() == MiningRewardSenderAddress && lastTransaction.RecipientAddress() == service.address {
								service.miningTicker.Reset(service.miningTimer / OutboundsCount)
								service.hasMiningTimerChanged = true
								service.logger.Info("mining timer changed")
							}
						}
						// Remove transactions which are in the new blocks
						newTransactions := oldTransactions
						for i := len(service.blocks) - 2; i < len(selectedBlocks); i++ {
							for _, validatedTransaction := range selectedBlocks[i].transactions {
								for j, transaction := range newTransactions {
									if validatedTransaction.Equals(transaction) {
										newTransactions[j] = newTransactions[len(newTransactions)-1]
										newTransactions = newTransactions[:len(newTransactions)-1]
									}
								}
							}
						}
						service.transactionsMutex.Lock()
						service.transactions = newTransactions
						service.transactionsMutex.Unlock()
					}
					service.blocksMutex.RUnlock()
					service.blocksMutex.Lock()
					defer service.blocksMutex.Unlock()
					service.blockResponses = selectedBlocksResponse
					service.blocks = selectedBlocks
					service.logger.Info("conflicts resolved: blockchain replaced")
				} else {
					service.transactionsMutex.Lock()
					service.transactions = nil
					service.transactionsMutex.Unlock()
					service.logger.Info("conflicts resolved: blockchain kept")
				}
			} else {
				service.transactionsMutex.Lock()
				service.transactions = nil
				service.transactionsMutex.Unlock()
				service.logger.Info("conflicts resolved: blockchain kept")
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
	for i, neighbor := range neighbors {
		if neighbor == removedNeighbor {
			neighbors = append(neighbors[:i], neighbors[i+1:]...)
			return neighbors
		}
	}
	return neighbors
}

func (service *Service) createBlock(rewardTransaction *Transaction) *Block {
	service.transactionsMutex.Lock()
	defer service.transactionsMutex.Unlock()
	service.blocksMutex.Lock()
	defer service.blocksMutex.Unlock()
	service.transactions = append(service.transactions, rewardTransaction)
	var lastBlockHash [32]byte
	if service.blocks != nil {
		lastBlock := service.blocks[len(service.blocks)-1]
		var err error
		lastBlockHash, err = lastBlock.Hash()
		if err != nil {
			service.logger.Error(fmt.Errorf("failed calculate last block hash: %w", err).Error())
			return nil
		}
	}
	block := NewBlock(lastBlockHash, service.transactions)
	service.transactions = nil
	service.blocks = append(service.blocks, block)
	blockResponse := block.GetResponse()
	service.blockResponses = append(service.blockResponses, blockResponse)
	return block
}
