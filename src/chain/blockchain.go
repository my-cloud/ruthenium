package chain

import (
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"ruthenium/src/log"
	"sync"
	"time"
)

const (
	DefaultPort = 8106

	MiningDifficulty          = 3
	MiningRewardSenderAddress = "MINING REWARD SENDER ADDRESS"
	MiningReward              = 10.0
	MiningTimerSec            = 60

	NeighborSynchronizationTimeSecond  = 10
	HostConnectionTimeoutSecond        = 10
	NeighborClientFindingTimeoutSecond = 5
)

type Blockchain struct {
	transactions      []*Transaction
	transactionsMutex sync.RWMutex
	blocks            []*Block
	blockResponses    []*BlockResponse
	address           string
	mineMutex         sync.Mutex
	miningStarted     bool
	miningStopped     bool

	ip        string
	port      uint16
	logger    *log.Logger
	waitGroup *sync.WaitGroup

	neighbors              []*Node // TODO manage max neighbors count (Outbound/Inbound)
	neighborsMutex         sync.RWMutex
	neighborsByTarget      map[string]*Node
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
	blockchain.addBlock(new(Block))
	seedsIps := []string{
		"89.82.76.241",
	}
	blockchain.neighborsByTarget = map[string]*Node{}
	for _, seedIp := range seedsIps {
		seed := NewNode(seedIp, DefaultPort, logger)
		blockchain.neighborsByTarget[seed.Target()] = seed
		if err := seed.StartClient(); err != nil {
			blockchain.logger.Error(fmt.Sprintf("Failed to start neighbor client for target %s\n%v", seed.Target(), err))
		}
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
	go func(neighborsByTarget map[string]*Node) {
		defer blockchain.waitGroup.Done()
		var neighbors []*Node
		var targetRequests []TargetRequest
		hostTargetRequest := TargetRequest{
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
				blockchain.logger.Error(fmt.Sprintf("ERROR: DNS discovery failed on addresse %s: %v", neighborIp, err))
				return
			}

			neighborsCount := len(lookedUpNeighborsIps)
			if neighborsCount != 1 {
				blockchain.logger.Error(fmt.Sprintf("ERROR: DNS discovery did not find a single address (%d addresses found) for the given IP %s", neighborsCount, neighborIp))
				return
			}
			lookedUpNeighborIp := lookedUpNeighborsIps[0]
			lookedUpNeighborIpString := lookedUpNeighborIp.String()
			if (lookedUpNeighborIpString != blockchain.ip || neighborPort != blockchain.port) && lookedUpNeighborIpString == neighborIp && neighbor.IsFound() {
				neighbors = append(neighbors, neighbor)
				targetRequest := TargetRequest{
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
			var neighborTargetRequests []TargetRequest
			for _, targetRequest := range targetRequests {
				neighborIp := neighbor.Ip()
				neighborPort := neighbor.Port()
				if neighborIp != *targetRequest.Ip || neighborPort != *targetRequest.Port {
					neighborTargetRequests = append(neighborTargetRequests, targetRequest)
				}
			}
			go func(neighbor *Node) {
				_ = neighbor.SendTargets(neighborTargetRequests)
			}(neighbor)
		}
	}(blockchain.neighborsByTarget)
}

func (blockchain *Blockchain) AddTargets(targetRequests []TargetRequest) {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		blockchain.neighborsByTargetMutex.Lock()
		defer blockchain.neighborsByTargetMutex.Unlock()
		for _, targetRequest := range targetRequests {
			neighbor := NewNode(*targetRequest.Ip, *targetRequest.Port, blockchain.logger)
			blockchain.neighborsByTarget[neighbor.Target()] = neighbor
			if _, ok := blockchain.neighborsByTarget[neighbor.Target()]; !ok {
				blockchain.neighborsByTarget[neighbor.Target()] = neighbor
				if err := neighbor.StartClient(); err != nil {
					blockchain.logger.Error(fmt.Sprintf("Failed to start neighbor client for target %s\n%v", neighbor.Target(), err))
				}
			}
		}
	}()
}

func (blockchain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blocks"`
	}{
		Blocks: blockchain.blocks,
	})
}

func (blockchain *Blockchain) CreateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		blockchain.UpdateTransaction(senderAddress, recipientAddress, senderPublicKey, value, signature)
		publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
		signatureStr := signature.String()
		var verb = PUT
		transactionRequest := TransactionRequest{
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
			go func(neighbor *Node) {
				_ = neighbor.UpdateTransactions(transactionRequest)
			}(neighbor)
		}
	}()
}

func (blockchain *Blockchain) UpdateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		blockchain.neighborsMutex.Lock()
		defer blockchain.neighborsMutex.Unlock()
		transaction := NewTransaction(senderAddress, senderPublicKey, recipientAddress, value)
		err := blockchain.addTransaction(transaction, signature)
		if err != nil {
			blockchain.logger.Error(fmt.Sprintf("ERROR: Failed to add transaction\n%v", err))
		}
	}()
}

func (blockchain *Blockchain) addTransaction(transaction *Transaction, signature *Signature) (err error) {
	if transaction.SenderAddress() == MiningRewardSenderAddress {
		blockchain.transactionsMutex.Lock()
		defer blockchain.transactionsMutex.Unlock()
		blockchain.transactions = append(blockchain.transactions, transaction)
	} else if transaction.Verify(signature) {
		if blockchain.CalculateTotalAmount(transaction.SenderAddress()) < transaction.Value() {
			err = errors.New("not enough balance in the sender wallet")
		} else {
			blockchain.transactionsMutex.Lock()
			defer blockchain.transactionsMutex.Unlock()
			blockchain.transactions = append(blockchain.transactions, transaction)
		}
	} else {
		err = errors.New("failed to verify transaction")
	}
	return
}

func (blockchain *Blockchain) Mine() {
	blockchain.waitGroup.Add(1)
	go func() {
		defer blockchain.waitGroup.Done()
		blockchain.mineMutex.Lock()
		defer blockchain.mineMutex.Unlock()

		transaction := NewTransaction(MiningRewardSenderAddress, nil, blockchain.address, MiningReward)
		if err := blockchain.addTransaction(transaction, nil); err != nil {
			blockchain.logger.Error(fmt.Sprintf("ERROR: Failed to mine, error: %v", err))
		}
		block := blockchain.createBlock()
		blockchain.addBlock(block)

		blockchain.neighborsMutex.RLock()
		defer blockchain.neighborsMutex.RUnlock()
		for _, neighbor := range blockchain.neighbors {
			go func(neighbor *Node) {
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

func (blockchain *Blockchain) Transactions() []*Transaction {
	blockchain.transactionsMutex.RLock()
	defer blockchain.transactionsMutex.RUnlock()
	return blockchain.transactions
}

func (blockchain *Blockchain) Blocks() []*BlockResponse {
	return blockchain.blockResponses
}

func (blockchain *Blockchain) clearTransactions() {
	blockchain.transactionsMutex.Lock()
	defer blockchain.transactionsMutex.Unlock()
	blockchain.transactions = nil
}

func (blockchain *Blockchain) GetValidBlocks(blocks []*BlockResponse) (validBlocks []*Block) {
	previousBlock := NewBlockFromDto(blocks[0])
	validBlocks = append(validBlocks, previousBlock)
	currentIndex := 1
	for currentIndex < len(blocks) {
		currentBlock := blocks[currentIndex]
		isPreviousHashValid := currentBlock.PreviousHash == previousBlock.Hash()
		if !isPreviousHashValid {
			return nil
		}

		block := NewBlockFromDto(currentBlock)
		if block.IsInValid(MiningDifficulty) {
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
		var longestChainResponse []*BlockResponse
		var longestChain []*Block
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
			blockchain.blockResponses = longestChainResponse
			blockchain.blocks = longestChain
			blockchain.clearTransactions()
			blockchain.logger.Info("Conflicts resolved: blockchain replaced")
		} else {
			blockchain.logger.Info("Conflicts resolved: blockchain kept")
		}
	}()
}

func (blockchain *Blockchain) addBlock(block *Block) *Block {
	blockchain.blocks = append(blockchain.blocks, block)
	blockResponse := block.GetDto()
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.clearTransactions()
	return block
}

func (blockchain *Blockchain) lastBlock() *Block {
	return blockchain.blocks[len(blockchain.blocks)-1]
}

func (blockchain *Blockchain) createBlock() *Block {
	previousHash := blockchain.lastBlock().Hash()
	var nonce int
	blockchain.transactionsMutex.RLock()
	defer blockchain.transactionsMutex.RUnlock()
	for {
		block := NewBlock(nonce, previousHash, blockchain.transactions)
		if block.IsInValid(MiningDifficulty) {
			nonce++
		} else {
			return block
		}
	}
}
