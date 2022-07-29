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
	DefaultPort = 8107

	MiningDifficulty          = 3
	MiningRewardSenderAddress = "MINING REWARD SENDER ADDRESS"
	MiningReward              = 10.0
	MiningTimerSec            = 60

	NeighborSynchronizationTimeSecond  = 10
	HostConnectionTimeoutSecond        = 10
	NeighborClientFindingTimeoutSecond = 5
)

type Blockchain struct {
	transactions   []*Transaction
	blocks         []*Block
	blockResponses []*BlockResponse
	address        string
	mineMutex      sync.Mutex
	miningStarted  bool
	miningStopped  bool

	ip     string
	port   uint16
	logger *log.Logger

	neighbors         []*Node
	neighborsMutex    sync.Mutex
	neighborsByTarget map[string]*Node
}

func NewBlockchain(address string, ip string, port uint16, logger *log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.address = address
	blockchain.ip = ip
	blockchain.port = port
	blockchain.logger = logger
	blockchain.createBlock(0, new(Block).Hash())
	seedsIps := []string{
		"89.82.76.241",
	}
	blockchain.neighborsByTarget = map[string]*Node{}
	for _, seedIp := range seedsIps {
		seed := NewNode(seedIp, DefaultPort, logger)
		blockchain.neighborsByTarget[seed.Target()] = seed
	}
	return blockchain
}

func (blockchain *Blockchain) Run() {
	blockchain.StartNeighborsSynchronization()
	blockchain.ResolveConflicts()
}

func (blockchain *Blockchain) SynchronizeNeighbors() {
	blockchain.neighborsMutex.Lock()
	defer blockchain.neighborsMutex.Unlock()
	blockchain.FindNeighbors()
}

func (blockchain *Blockchain) StartNeighborsSynchronization() {
	blockchain.SynchronizeNeighbors()
	_ = time.AfterFunc(time.Second*NeighborSynchronizationTimeSecond, blockchain.StartNeighborsSynchronization)
}

func (blockchain *Blockchain) FindNeighbors() {
	go func(neighborsByTarget map[string]*Node) {
		var neighbors []*Node
		for _, neighbor := range neighborsByTarget {
			neighborsIps, err := net.LookupIP(neighbor.Ip())
			if err != nil {
				blockchain.logger.Error(fmt.Sprintf("ERROR: DNS discovery failed on addresse %s: %v", neighbor.Ip(), err))
				return
			}

			numNeighbors := len(neighborsIps)
			if numNeighbors != 1 {
				blockchain.logger.Error(fmt.Sprintf("ERROR: DNS discovery did not find a single address (%d addresses found) for the given IP %s", numNeighbors, neighbor.Ip()))
				return
			}
			neighborIp := neighborsIps[0]
			if (neighborIp.String() != blockchain.ip || neighbor.Port() != blockchain.port) && neighborIp.String() == neighbor.Ip() && neighbor.IsFound() {
				neighbors = append(neighbors, neighbor)
				kind := PostTargetRequest
				request := TargetRequest{
					Kind: &kind,
					Ip:   &blockchain.ip,
					Port: &blockchain.port,
				}
				_ = neighbor.SendTarget(request)
			}
		}
		blockchain.neighbors = neighbors
	}(blockchain.neighborsByTarget)
}

func (blockchain *Blockchain) AddTarget(ip string, port uint16) {
	neighbor := NewNode(ip, port, blockchain.logger)
	blockchain.neighborsByTarget[neighbor.Target()] = neighbor
}

func (blockchain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blocks"`
	}{
		Blocks: blockchain.blocks,
	})
}

func (blockchain *Blockchain) CreateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) {
	go func() {
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
		go func(neighbors []*Node) {
			for _, neighbor := range neighbors {
				_ = neighbor.UpdateTransactions(transactionRequest)
			}
		}(blockchain.neighbors)
	}()
}

func (blockchain *Blockchain) UpdateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) {
	transaction := NewTransaction(senderAddress, senderPublicKey, recipientAddress, value)
	err := blockchain.addTransaction(transaction, signature)
	if err != nil {
		blockchain.logger.Error(fmt.Sprintf("ERROR: Failed to add transaction\n%v", err))
	}
}

func (blockchain *Blockchain) addTransaction(transaction *Transaction, signature *Signature) (err error) {
	if transaction.SenderAddress() == MiningRewardSenderAddress {
		blockchain.transactions = append(blockchain.transactions, transaction)
	} else if transaction.VerifySignature(signature) {
		if blockchain.CalculateTotalAmount(transaction.SenderAddress()) < transaction.Value() {
			err = errors.New("not enough balance in the sender wallet")
		} else {
			blockchain.transactions = append(blockchain.transactions, transaction)
		}
	} else {
		err = errors.New("failed to verify transaction")
	}
	return
}

func (blockchain *Blockchain) Mine() {
	go func() {
		blockchain.mineMutex.Lock()
		defer blockchain.mineMutex.Unlock()

		transaction := NewTransaction(MiningRewardSenderAddress, nil, blockchain.address, MiningReward)
		if err := blockchain.addTransaction(transaction, nil); err != nil {
			blockchain.logger.Error(fmt.Sprintf("ERROR: Failed to mine, error: %v", err))
		}
		nonce := blockchain.proofOfWork()
		previousHash := blockchain.lastBlock().Hash()
		blockchain.createBlock(nonce, previousHash)

		go func(neighbors []*Node) {
			for _, neighbor := range neighbors {
				_ = neighbor.Consensus()
			}
		}(blockchain.neighbors)
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
	return blockchain.transactions
}

func (blockchain *Blockchain) Blocks() []*BlockResponse {
	return blockchain.blockResponses
}

func (blockchain *Blockchain) clearTransactions() {
	blockchain.transactions = nil
}

func (blockchain *Blockchain) GetValidBlocks(blocks []*BlockResponse) (validBlocks []*Block) {
	previousBlock := NewBlockFromDto(blocks[0])
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
	var longestChainResponse []*BlockResponse
	var longestChain []*Block
	maxLength := len(blockchain.blocks)

	go func(neighbors []*Node) {
		for _, neighbor := range neighbors {
			neighborBlocks, err := neighbor.GetBlocks()
			if err == nil && len(neighborBlocks) > maxLength {
				validBlocks := blockchain.GetValidBlocks(neighborBlocks)
				if validBlocks != nil {
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
		}
		blockchain.logger.Info("Conflicts resolved: blockchain kept")
	}(blockchain.neighbors)
}

func (blockchain *Blockchain) createBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, blockchain.transactions)
	blockchain.blocks = append(blockchain.blocks, block)
	blockResponse := block.GetDto()
	blockchain.blockResponses = append(blockchain.blockResponses, blockResponse)
	blockchain.clearTransactions()
	return block
}

func (blockchain *Blockchain) lastBlock() *Block {
	return blockchain.blocks[len(blockchain.blocks)-1]
}

func (blockchain *Blockchain) copyTransactions() (transactions []*Transaction) {
	for _, transaction := range blockchain.transactions {
		transactions = append(transactions,
			NewTransaction(transaction.SenderAddress(),
				transaction.SenderPublicKey(),
				transaction.RecipientAddress(),
				transaction.Value()))
	}
	return
}

func (blockchain *Blockchain) proofOfWork() int {
	transactions := blockchain.copyTransactions()
	previousHash := blockchain.lastBlock().Hash()
	var nonce int
	for NewBlock(nonce, previousHash, transactions).IsInValid(MiningDifficulty) {
		nonce++
	}
	return nonce
}
