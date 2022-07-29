package chain

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"net"
	"ruthenium/src/log"
	"strings"
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
	NeighborClientFindingTimeoutSecond = 5
)

type Blockchain struct {
	transactions  []*Transaction
	blocks        []*Block
	address       string
	mineMutex     sync.Mutex
	miningStarted bool
	miningStopped bool

	ip     string
	port   uint16
	logger *log.Logger

	neighbors         []*Node
	neighborsMutex    sync.Mutex
	transactionMutex  sync.Mutex
	neighborsByTarget map[string]*Node
}

func NewBlockchain(address string, ip string, port uint16, logger *log.Logger) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.address = address
	blockchain.ip = ip
	blockchain.port = port
	blockchain.logger = logger
	blockchain.createBlock(0, new(Block).Hash())
	seeds := []string{
		"89.82.76.241",
	}
	blockchain.neighborsByTarget = map[string]*Node{}
	for _, seed := range seeds {
		blockchain.neighborsByTarget[fmt.Sprintf("%s:%d", seed, DefaultPort)] = NewNode(seed, DefaultPort, logger)
	}
	return blockchain
}

func (blockchain *Blockchain) Run() {
	blockchain.StartNeighborsSynchronization()
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
	blockchain.neighbors = nil
	for _, neighbor := range blockchain.neighborsByTarget {
		go func(neighbor *Node) {
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
				blockchain.neighbors = append(blockchain.neighbors, neighbor)
				kind := PostTargetRequest
				request := TargetRequest{
					Kind: &kind,
					Ip:   &blockchain.ip,
					Port: &blockchain.port,
				}
				go neighbor.SendTarget(request)
				blockchain.ResolveConflicts()
			}
		}(neighbor)
	}
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

func (blockchain *Blockchain) CreateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) bool {
	isTransacted := blockchain.UpdateTransaction(senderAddress, recipientAddress, senderPublicKey, value, signature)

	if isTransacted {
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
		for _, neighbor := range blockchain.neighbors {
			go neighbor.UpdateTransactions(transactionRequest)
		}
	}

	return isTransacted
}

func (blockchain *Blockchain) UpdateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) (isTransacted bool) {
	transaction := NewTransaction(senderAddress, senderPublicKey, recipientAddress, value)
	return blockchain.addTransaction(transaction, signature)
}

func (blockchain *Blockchain) addTransaction(transaction *Transaction, signature *Signature) bool {
	blockchain.transactionMutex.Lock()
	defer blockchain.transactionMutex.Unlock()
	if transaction.SenderAddress() == MiningRewardSenderAddress {
		blockchain.transactions = append(blockchain.transactions, transaction)
		return true
	}

	if transaction.VerifySignature(signature) {
		if blockchain.CalculateTotalAmount(transaction.SenderAddress()) < transaction.Value() {
			blockchain.logger.Error("ERROR: Not enough balance in the sender wallet")
			return false
		}
		blockchain.transactions = append(blockchain.transactions, transaction)
		return true
	} else {
		blockchain.logger.Error("ERROR: Failed to verify transaction")
	}
	return false

}

func (blockchain *Blockchain) Mine() bool {
	blockchain.mineMutex.Lock()
	defer blockchain.mineMutex.Unlock()

	transaction := NewTransaction(MiningRewardSenderAddress, nil, blockchain.address, MiningReward)
	blockchain.addTransaction(transaction, nil)
	nonce := blockchain.proofOfWork()
	previousHash := blockchain.lastBlock().Hash()
	blockchain.createBlock(nonce, previousHash)

	for _, neighbor := range blockchain.neighbors {
		go neighbor.Consensus()
	}

	return true
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
	var totalAmount float32 = 0.0
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

func (blockchain *Blockchain) Blocks() []*Block {
	return blockchain.blocks
}

func (blockchain *Blockchain) ClearTransactions() {
	blockchain.transactionMutex.Lock()
	defer blockchain.transactionMutex.Unlock()
	blockchain.transactions = nil
}

func (blockchain *Blockchain) IsValid(blocks []*Block) bool {
	previousBlock := blocks[0]
	currentIndex := 1
	for currentIndex > len(blocks) {
		currentBlock := blocks[currentIndex]
		isPreviousHashValid := currentBlock.PreviousHash() == previousBlock.Hash()
		if !isPreviousHashValid {
			return false
		}

		isProofValid := blockchain.isProofValid(currentBlock.Nonce(), currentBlock.PreviousHash(), currentBlock.Transactions(), MiningDifficulty)
		if !isProofValid {
			return false
		}

		previousBlock = currentBlock
		currentIndex++
	}
	return true
}

func (blockchain *Blockchain) ResolveConflicts() {
	var longestChain []*Block
	maxLength := len(blockchain.blocks)

	go func(neighbors []*Node) {
		blockchain.transactionMutex.Lock()
		defer blockchain.transactionMutex.Unlock()
		for _, neighbor := range blockchain.neighbors {
			neighborBlocks := neighbor.GetBlocks()
			if len(neighborBlocks) > maxLength && blockchain.IsValid(neighborBlocks) {
				maxLength = len(neighborBlocks)
				longestChain = neighborBlocks
			}
		}

		if longestChain != nil {
			blockchain.blocks = longestChain
			blockchain.ClearTransactions()
			blockchain.logger.Info("Conflicts resolved: blockchain replaced")
		}
		blockchain.logger.Info("Conflicts resolved: blockchain kept")
	}(blockchain.neighbors)
}

func (blockchain *Blockchain) createBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, blockchain.transactions)
	blockchain.blocks = append(blockchain.blocks, block)
	blockchain.ClearTransactions()
	return block
}

func (blockchain *Blockchain) lastBlock() *Block {
	return blockchain.blocks[len(blockchain.blocks)-1]
}

func (blockchain *Blockchain) copyTransactions() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range blockchain.transactions {
		transactions = append(transactions,
			NewTransaction(transaction.SenderAddress(),
				transaction.SenderPublicKey(),
				transaction.RecipientAddress(),
				transaction.Value()))
	}
	return transactions
}

func (blockchain *Blockchain) isProofValid(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := NewBlock(nonce, previousHash, transactions)
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (blockchain *Blockchain) proofOfWork() int {
	transactions := blockchain.copyTransactions()
	previousHash := blockchain.lastBlock().Hash()
	var nonce int
	for !blockchain.isProofValid(nonce, previousHash, transactions, MiningDifficulty) {
		nonce++
	}
	return nonce
}
