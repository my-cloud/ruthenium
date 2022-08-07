package chain

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"ruthenium/src/log"
	"ruthenium/src/node/authentication"
	"ruthenium/src/node/chain/mine"
	"ruthenium/src/node/neighborhood"
	"sync"
	"time"
)

const (
	DefaultPort = 8106

	MiningRewardSenderAddress         = "MINING REWARD SENDER ADDRESS"
	MiningReward              float32 = 10.
	MiningTimerSec                    = 60

	NeighborSynchronizationTimeSecond = 10
	HostConnectionTimeoutSecond       = 10
)

type Blockchain struct {
	transactions      []*mine.Transaction
	transactionsMutex sync.RWMutex
	blocks            []*mine.Block
	blockResponses    []*neighborhood.BlockResponse
	blocksMutex       sync.RWMutex
	address           string
	mineMutex         sync.Mutex
	miningStarted     bool
	miningStopped     bool

	ip        string
	port      uint16
	logger    *log.Logger
	waitGroup *sync.WaitGroup

	neighbors              []*neighborhood.Node // TODO manage max neighbors count (Outbound/Inbound)
	neighborsMutex         sync.RWMutex
	neighborsByTarget      map[string]*neighborhood.Node
	neighborsByTargetMutex sync.RWMutex
}

func NewBlockchain(address string, ip string, port uint16, logger *log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.address = address
	blockchain.ip = ip
	blockchain.port = port
	blockchain.logger = logger
	var waitGroup sync.WaitGroup
	blockchain.waitGroup = &waitGroup
	blockchain.addBlock(new(mine.Block))
	seedsIps := []string{
		"89.82.76.241",
	}
	blockchain.neighborsByTarget = map[string]*neighborhood.Node{}
	for _, seedIp := range seedsIps {
		seed := neighborhood.NewNode(seedIp, DefaultPort, logger)
		blockchain.neighborsByTarget[seed.Target()] = seed
	}
	return blockchain
}

func (blockchain *Blockchain) WaitGroup() *sync.WaitGroup {
	return blockchain.waitGroup
}

func (blockchain *Blockchain) Run() {
	blockchain.StartNeighborsSynchronization()
}

func (blockchain *Blockchain) SynchronizeNeighbors() {
	blockchain.FindNeighbors()
}

func (blockchain *Blockchain) StartNeighborsSynchronization() {
	blockchain.SynchronizeNeighbors()
	_ = time.AfterFunc(time.Second*NeighborSynchronizationTimeSecond, blockchain.StartNeighborsSynchronization)
}

func (blockchain *Blockchain) FindNeighbors() {
	blockchain.waitGroup.Add(1)
	go func(neighborsByTarget map[string]*neighborhood.Node) {
		defer blockchain.waitGroup.Done()
		var neighbors []*neighborhood.Node
		var targetRequests []neighborhood.TargetRequest
		hostTargetRequest := neighborhood.TargetRequest{
			Ip:   &blockchain.ip,
			Port: &blockchain.port,
		}
		targetRequests = append(targetRequests, hostTargetRequest)
		var newNeighborFound bool
		blockchain.neighborsMutex.RLock()
		blockchain.neighborsByTargetMutex.RLock()
		for _, neighbor := range neighborsByTarget {
			neighborIp := neighbor.Ip()
			neighborPort := neighbor.Port()
			lookedUpNeighborsIps, err := net.LookupIP(neighborIp)
			if err != nil {
				blockchain.logger.Error(fmt.Errorf("DNS discovery failed on addresse %s: %w", neighborIp, err).Error())
				return
			}

			neighborsCount := len(lookedUpNeighborsIps)
			if neighborsCount != 1 {
				blockchain.logger.Error(fmt.Errorf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", neighborsCount, neighborIp).Error())
				return
			}
			lookedUpNeighborIp := lookedUpNeighborsIps[0]
			lookedUpNeighborIpString := lookedUpNeighborIp.String()
			if (lookedUpNeighborIpString != blockchain.ip || neighborPort != blockchain.port) && lookedUpNeighborIpString == neighborIp && neighbor.IsFound() {
				neighbors = append(neighbors, neighbor)
				targetRequest := neighborhood.TargetRequest{
					Ip:   &neighborIp,
					Port: &neighborPort,
				}
				targetRequests = append(targetRequests, targetRequest)
				newNeighborFound = true
				for _, oldNeighbor := range blockchain.neighbors {
					if oldNeighbor.Ip() == neighbor.Ip() && oldNeighbor.Port() == neighbor.Port() {
						newNeighborFound = false
						break
					}
				}
			}
		}
		blockchain.neighborsMutex.RUnlock()
		blockchain.neighborsByTargetMutex.RUnlock()
		blockchain.neighborsMutex.Lock()
		blockchain.neighbors = neighbors
		blockchain.neighborsMutex.Unlock()
		// TODO handle case where a known neighbor have been disconnected for a while (consider it as a new neighbor)
		if newNeighborFound {
			blockchain.ResolveConflicts()
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
			go func(neighbor *neighborhood.Node) {
				_ = neighbor.SendTargets(neighborTargetRequests)
			}(neighbor)
		}
	}(blockchain.neighborsByTarget)
}

func (blockchain *Blockchain) AddTargets(targetRequests []neighborhood.TargetRequest) {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		blockchain.neighborsByTargetMutex.Lock()
		defer blockchain.neighborsByTargetMutex.Unlock()
		for _, targetRequest := range targetRequests {
			neighbor := neighborhood.NewNode(*targetRequest.Ip, *targetRequest.Port, blockchain.logger)
			if _, ok := blockchain.neighborsByTarget[neighbor.Target()]; !ok {
				blockchain.neighborsByTarget[neighbor.Target()] = neighbor
			}
		}
	}()
}

func (blockchain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*mine.Block `json:"blocks"`
	}{
		Blocks: blockchain.blocks,
	})
}

func (blockchain *Blockchain) CreateTransaction(transaction *mine.Transaction, publicKey *ecdsa.PublicKey, signature *authentication.Signature) {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		err := blockchain.addTransaction(transaction, publicKey, signature)
		if err != nil {
			blockchain.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
			return
		}
		senderAddress := transaction.SenderAddress()
		recipientAddress := transaction.RecipientAddress()
		value := transaction.Value()
		// TODO create an object Public key holding the following statement
		publicKeyStr := fmt.Sprintf("%064x%064x", publicKey.X.Bytes(), publicKey.Y.Bytes())
		signatureStr := signature.String()
		var verb = neighborhood.PUT
		transactionRequest := neighborhood.TransactionRequest{
			Verb:             &verb,
			SenderAddress:    &senderAddress,
			RecipientAddress: &recipientAddress,
			SenderPublicKey:  &publicKeyStr,
			Value:            &value,
			Signature:        &signatureStr,
		}
		blockchain.neighborsMutex.RLock()
		defer blockchain.neighborsMutex.RUnlock()
		for _, neighbor := range blockchain.neighbors {
			go func(neighbor *neighborhood.Node) {
				_ = neighbor.AddTransaction(transactionRequest)
			}(neighbor)
		}
	}()
}

func (blockchain *Blockchain) AddTransaction(transaction *mine.Transaction, publicKey *ecdsa.PublicKey, signature *authentication.Signature) {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		err := blockchain.addTransaction(transaction, publicKey, signature)
		if err != nil {
			blockchain.logger.Error(fmt.Errorf("failed to add transaction: %w", err).Error())
		}
	}()
}

func (blockchain *Blockchain) addTransaction(transaction *mine.Transaction, publicKey *ecdsa.PublicKey, signature *authentication.Signature) (err error) {
	marshaledTransaction, err := json.Marshal(transaction)
	if err != nil {
		err = fmt.Errorf("failed to marshal transaction, %w", err)
		return
	}
	if !signature.Verify(marshaledTransaction, publicKey, transaction.SenderAddress()) {
		err = errors.New("failed to verify transaction")
		return
	}
	insufficientBalance := blockchain.CalculateTotalAmount(transaction.SenderAddress()) < transaction.Value()
	if insufficientBalance {
		err = errors.New("not enough balance in the sender wallet")
		return
	}
	blockchain.transactionsMutex.Lock()
	defer blockchain.transactionsMutex.Unlock()
	blockchain.transactions = append(blockchain.transactions, transaction)
	return
}

func (blockchain *Blockchain) Mine() {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		blockchain.mineMutex.Lock()
		defer blockchain.mineMutex.Unlock()
		transaction := mine.NewTransaction(MiningRewardSenderAddress, blockchain.address, MiningReward)
		blockchain.transactionsMutex.Lock()
		blockchain.transactions = append(blockchain.transactions, transaction)
		blockchain.transactionsMutex.Unlock()
		blockchain.blocksMutex.Lock()
		block, err := blockchain.createBlock()
		if err != nil {
			blockchain.logger.Error(fmt.Errorf("failed to create newly mined block: %v", err).Error())
			return
		}
		blockchain.addBlock(block)
		blockchain.blocksMutex.Unlock()
		blockchain.neighborsMutex.RLock()
		defer blockchain.neighborsMutex.RUnlock()
		for _, neighbor := range blockchain.neighbors {
			go func(neighbor *neighborhood.Node) {
				_ = neighbor.Consensus()
			}(neighbor)
		}
	}()
}

func (blockchain *Blockchain) StartMining() {
	if !blockchain.miningStarted {
		blockchain.miningStarted = true
		blockchain.miningStopped = false
		blockchain.mining()
	}
}

func (blockchain *Blockchain) mining() {
	if !blockchain.miningStopped {
		blockchain.Mine()
		_ = time.AfterFunc(time.Second*MiningTimerSec, blockchain.mining)
	}
}

func (blockchain *Blockchain) StopMining() {
	blockchain.miningStarted = false
	blockchain.miningStopped = true
}

func (blockchain *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32
	blockchain.blocksMutex.RLock()
	defer blockchain.blocksMutex.RUnlock()
	for _, block := range blockchain.blocks {
		for _, transaction := range block.Transactions() {
			value := transaction.Value()
			if blockchainAddress == transaction.RecipientAddress() {
				totalAmount += value
			}
			if blockchainAddress == transaction.SenderAddress() {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (blockchain *Blockchain) Transactions() []*mine.Transaction {
	blockchain.transactionsMutex.RLock()
	defer blockchain.transactionsMutex.RUnlock()
	return blockchain.transactions
}

func (blockchain *Blockchain) Blocks() []*neighborhood.BlockResponse {
	return blockchain.blockResponses
}

func (blockchain *Blockchain) clearTransactions() {
	blockchain.transactionsMutex.Lock()
	defer blockchain.transactionsMutex.Unlock()
	blockchain.transactions = nil
}

func (blockchain *Blockchain) GetValidBlocks(blocks []*neighborhood.BlockResponse) (validBlocks []*mine.Block) {
	previousBlock := mine.NewBlockFromDto(blocks[0])
	validBlocks = append(validBlocks, previousBlock)
	currentIndex := 1
	for currentIndex < len(blocks) {
		currentBlock := blocks[currentIndex]
		previousBlockHash, err := previousBlock.Hash()
		if err != nil {
			blockchain.logger.Error(fmt.Errorf("failed calculate previous block hash: %w", err).Error())
		}
		isPreviousHashValid := currentBlock.PreviousHash == previousBlockHash
		if !isPreviousHashValid {
			return nil
		}
		block := mine.NewBlockFromDto(currentBlock)
		var pow *mine.ProofOfWork
		if pow, err = block.ProofOfWork(); err != nil {
			blockchain.logger.Error(fmt.Errorf("failed to get proof of work: %w", err).Error())
			return nil
		}
		if pow.IsInValid() {
			blockchain.logger.Info("proof of work is invalid")
			return nil
		}
		previousBlock = block
		validBlocks = append(validBlocks, block)
		currentIndex++
	}
	return validBlocks
}

func (blockchain *Blockchain) ResolveConflicts() {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		blockchain.neighborsMutex.RLock()
		defer blockchain.neighborsMutex.RUnlock()
		var longestChainResponse []*neighborhood.BlockResponse
		var longestChain []*mine.Block
		maxLength := len(blockchain.blocks)
		for _, neighbor := range blockchain.neighbors {
			neighborBlocks, err := neighbor.GetBlocks()
			if err == nil && len(neighborBlocks) > maxLength {
				validBlocks := blockchain.GetValidBlocks(neighborBlocks)
				if len(validBlocks) > 1 {
					maxLength = len(neighborBlocks)
					longestChainResponse = neighborBlocks
					longestChain = validBlocks
				}
			}
		}
		if longestChain != nil {
			blockchain.blocksMutex.Lock()
			defer blockchain.blocksMutex.Unlock()
			blockchain.blockResponses = longestChainResponse
			blockchain.blocks = longestChain
			blockchain.clearTransactions()
			blockchain.logger.Info("conflicts resolved: blockchain replaced")
		} else {
			blockchain.logger.Info("conflicts resolved: blockchain kept")
		}
	}()
}

func (blockchain *Blockchain) addBlock(block *mine.Block) *mine.Block {
	blockchain.blocks = append(blockchain.blocks, block)
	blockResponse := block.GetDto()
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.clearTransactions()
	return block
}

func (blockchain *Blockchain) createBlock() (block *mine.Block, err error) {
	lastBlock := blockchain.blocks[len(blockchain.blocks)-1]
	lastBlockHash, err := lastBlock.Hash()
	if err != nil {
		err = fmt.Errorf("failed calculate last block hash: %w", err)
		return
	}
	blockchain.transactionsMutex.RLock()
	defer blockchain.transactionsMutex.RUnlock()
	block, err = mine.NewBlock(lastBlockHash, blockchain.transactions)
	return
}
