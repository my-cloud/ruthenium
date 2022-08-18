package blockchain

import (
	"errors"
	"fmt"
	"math"
	"net"
	"ruthenium/src/log"
	"ruthenium/src/node/authentication"
	"ruthenium/src/node/blockchain/mining"
	"ruthenium/src/node/neighborhood"
	"sync"
	"time"
)

const (
	DefaultPort = 8106

	MiningRewardSenderAddress        = "MINING REWARD SENDER ADDRESS"
	ParticlesCount                   = 100000000
	GenesisAmount             uint64 = 100000 * ParticlesCount
	RewardExponent                   = 1 / 1.828393264
	MinutesCountPerDay               = 1440
	HalfLifeInDays                   = 373.59

	NeighborSynchronizationTimeInSeconds = 10
)

type Service struct {
	transactions      []*mining.Transaction
	transactionsMutex sync.RWMutex
	blocks            []*mining.Block
	blockResponses    []*neighborhood.BlockResponse
	blocksMutex       sync.RWMutex
	address           string
	mineMutex         sync.Mutex
	miningStarted     bool
	miningStopped     bool
	mineRequested     bool
	miningTicker      <-chan time.Time

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
	service.miningTicker = time.Tick(miningTimer)
	service.address = address
	service.ip = ip
	service.port = port
	service.logger = logger
	var waitGroup sync.WaitGroup
	service.waitGroup = &waitGroup
	var transactions []*mining.Transaction
	now := time.Now().Unix() * time.Second.Nanoseconds()
	transactions = append(transactions, mining.NewTransaction(now, MiningRewardSenderAddress, service.address, GenesisAmount))
	genesisBlock, _ := mining.NewBlock([32]byte{}, transactions)
	service.addBlock(genesisBlock)
	seedsIps := []string{
		"89.82.76.241",
	}
	service.neighborsByTarget = map[string]*neighborhood.Neighbor{}
	for _, seedIp := range seedsIps {
		seed := neighborhood.NewNeighbor(seedIp, DefaultPort, logger)
		service.neighborsByTarget[seed.Target()] = seed
	}
	halfLife := HalfLifeInDays * MinutesCountPerDay * float64(time.Minute.Nanoseconds())
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
		var newNeighborFound bool
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
				service.logger.Error(fmt.Errorf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", neighborsCount, neighborIp).Error())
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
				newNeighborFound = true
				for _, oldNeighbor := range service.neighbors {
					if oldNeighbor.Ip() == neighbor.Ip() && oldNeighbor.Port() == neighbor.Port() {
						newNeighborFound = false
						break
					}
				}
			}
		}
		service.neighborsMutex.RUnlock()
		service.neighborsByTargetMutex.RUnlock()
		service.neighborsMutex.Lock()
		service.neighbors = neighbors
		service.neighborsMutex.Unlock()
		// TODO handle case where a known neighbor have been disconnected for a while (consider it as a new neighbor)
		if newNeighborFound {
			service.ResolveConflicts()
		}
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

func (service *Service) CreateTransaction(transaction *mining.Transaction, publicKey *authentication.PublicKey, signature *authentication.Signature) {
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

func (service *Service) AddTransaction(transaction *mining.Transaction, publicKey *authentication.PublicKey, signature *authentication.Signature) {
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		err := service.addTransaction(transaction, publicKey, signature)
		if err != nil {
			service.logger.Error(fmt.Errorf("failed to add transaction: %w", err).Error())
		}
	}()
}

func (service *Service) addTransaction(transaction *mining.Transaction, publicKey *authentication.PublicKey, signature *authentication.Signature) (err error) {
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
		<-service.miningTicker
		service.mine()
		service.mineRequested = false
	}()
}

func (service *Service) mine() {
	service.mineMutex.Lock()
	defer service.mineMutex.Unlock()
	now := time.Now().Unix() * time.Second.Nanoseconds()
	amount := service.CalculateTotalAmount(now, service.address)
	var reward uint64
	if amount > 0 {
		reward = uint64(math.Round(math.Pow(float64(amount), RewardExponent)))
	} else {
		reward = 0
	}
	service.logger.Info(fmt.Sprintf("amount: %d - reward: %d - total: %d", amount, reward, amount+reward))
	transaction := mining.NewTransaction(now, MiningRewardSenderAddress, service.address, reward)
	service.transactionsMutex.Lock()
	service.transactions = append(service.transactions, transaction)
	service.transactionsMutex.Unlock()
	service.blocksMutex.Lock()
	block, err := service.createBlock()
	if err != nil {
		service.logger.Error(fmt.Errorf("failed to create newly mined block: %v", err).Error())
		return
	}
	service.addBlock(block)
	service.blocksMutex.Unlock()
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
		<-service.miningTicker
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
					service.logger.Error(fmt.Errorf("historical transaction should not have been validated: wallet amount=%d, transaction value=%d", totalAmount, value).Error())
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

func (service *Service) Transactions() []*mining.Transaction {
	service.transactionsMutex.RLock()
	defer service.transactionsMutex.RUnlock()
	return service.transactions
}

func (service *Service) Blocks() []*neighborhood.BlockResponse {
	return service.blockResponses
}

func (service *Service) clearTransactions() {
	service.transactionsMutex.Lock()
	defer service.transactionsMutex.Unlock()
	service.transactions = nil
}

func (service *Service) GetValidBlocks(blocks []*neighborhood.BlockResponse) (validBlocks []*mining.Block) {
	previousBlock := mining.NewBlockFromResponse(blocks[0])
	validBlocks = append(validBlocks, previousBlock)
	currentIndex := 1
	for currentIndex < len(blocks) {
		currentBlock := blocks[currentIndex]
		previousBlockHash, err := previousBlock.Hash()
		if err != nil {
			service.logger.Error(fmt.Errorf("failed calculate previous block hash: %w", err).Error())
		}
		isPreviousHashValid := currentBlock.PreviousHash == previousBlockHash
		if !isPreviousHashValid {
			return nil
		}
		block := mining.NewBlockFromResponse(currentBlock)
		var pow *mining.ProofOfWork
		if pow, err = block.ProofOfWork(); err != nil {
			service.logger.Error(fmt.Errorf("failed to get proof of work: %w", err).Error())
			return nil
		}
		if pow.IsInValid() {
			service.logger.Info("proof of work is invalid")
			return nil
		}
		previousBlock = block
		validBlocks = append(validBlocks, block)
		currentIndex++
	}
	return validBlocks
}

func (service *Service) ResolveConflicts() {
	service.waitGroup.Add(1)
	go func() {
		defer service.waitGroup.Done()
		service.neighborsMutex.RLock()
		defer service.neighborsMutex.RUnlock()
		var longestChainResponse []*neighborhood.BlockResponse
		var longestChain []*mining.Block
		maxLength := len(service.blocks)
		for _, neighbor := range service.neighbors {
			neighborBlocks, err := neighbor.GetBlocks()
			if err == nil && len(neighborBlocks) > maxLength {
				validBlocks := service.GetValidBlocks(neighborBlocks)
				if len(validBlocks) > 1 {
					maxLength = len(neighborBlocks)
					longestChainResponse = neighborBlocks
					longestChain = validBlocks
				}
			}
		}
		if longestChain != nil {
			service.blocksMutex.Lock()
			defer service.blocksMutex.Unlock()
			service.blockResponses = longestChainResponse
			service.blocks = longestChain
			service.clearTransactions()
			service.logger.Info("conflicts resolved: blockchain replaced")
		} else {
			service.logger.Info("conflicts resolved: blockchain kept")
		}
	}()
}

func (service *Service) addBlock(block *mining.Block) *mining.Block {
	service.blocks = append(service.blocks, block)
	blockResponse := block.GetResponse()
	service.blockResponses = append(service.blockResponses, blockResponse)
	service.clearTransactions()
	return block
}

func (service *Service) createBlock() (block *mining.Block, err error) {
	lastBlock := service.blocks[len(service.blocks)-1]
	lastBlockHash, err := lastBlock.Hash()
	if err != nil {
		err = fmt.Errorf("failed calculate last block hash: %w", err)
		return
	}
	service.transactionsMutex.RLock()
	defer service.transactionsMutex.RUnlock()
	block, err = mining.NewBlock(lastBlockHash, service.transactions)
	return
}
