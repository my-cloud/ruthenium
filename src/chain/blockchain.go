package chain

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

const (
	DefaultPort = 8106

	MiningDifficulty          = 3
	MiningRewardSenderAddress = "MINING REWARD SENDER ADDRESS"
	MiningReward              = 1.0
	MiningTimerSec            = 1

	NeighborSynchronizationTimeSecond  = 5
	HostConnectionTimeoutSecond        = 10
	NeighborClientFindingTimeoutSecond = 1
)

type Blockchain struct {
	transactions  []*Transaction
	blocks        []*Block
	address       string
	mineMutex     sync.Mutex
	miningStopped bool

	ip   string
	port uint16

	neighbors      []*Node
	neighborsMutex sync.Mutex
	seeds          []string
}

func NewBlockchain(address string, ip string, port uint16) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.address = address
	blockchain.ip = ip
	blockchain.port = port
	blockchain.createBlock(0, new(Block).Hash())
	blockchain.seeds = []string{
		"89.82.76.241",
	}
	// 1ATviRXq6fG6QGBJcPJfJs6sU42SS1qAmM
	return blockchain
}

func (blockchain *Blockchain) Run() {
	blockchain.StartNeighborsSynchronization()
	blockchain.ResolveConflicts()
}

func (blockchain *Blockchain) SynchronizeNeighbors() {
	blockchain.neighborsMutex.Lock()
	defer blockchain.neighborsMutex.Unlock()
	blockchain.neighbors = blockchain.FindNeighbors()
}

func (blockchain *Blockchain) StartNeighborsSynchronization() {
	blockchain.SynchronizeNeighbors()
	_ = time.AfterFunc(time.Second*NeighborSynchronizationTimeSecond, blockchain.StartNeighborsSynchronization)
}

func (blockchain *Blockchain) FindNeighbors() []*Node {
	neighbors := make([]*Node, 0)
	for _, seed := range blockchain.seeds {
		go func(seed string) {
			seedPeers, err := net.LookupIP(seed)
			if err != nil {
				log.Printf("ERROR: DNS discovery failed on seed %s: %v", seed, err)
				return
			}

			numPeers := len(seedPeers)
			log.Printf("%d addresses found from DNS seed %s", numPeers, seed)
			if numPeers == 0 {
				return
			}
			for _, peer := range seedPeers {
				address := peer.String()
				if address != blockchain.ip {
					neighbor := NewNode(address, DefaultPort)
					neighbors = append(neighbors, neighbor)
					neighbor.StartClient()
				}
			}
		}(seed)
	}
	return neighbors
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
			neighbor.UpdateTransactions(transactionRequest)
		}
	}

	return isTransacted
}

func (blockchain *Blockchain) UpdateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) (isTransacted bool) {
	transaction := NewTransaction(senderAddress, senderPublicKey, recipientAddress, value)
	return blockchain.addTransaction(transaction, signature)
}

func (blockchain *Blockchain) addTransaction(transaction *Transaction, signature *Signature) bool {
	if transaction.SenderAddress() == MiningRewardSenderAddress {
		blockchain.transactions = append(blockchain.transactions, transaction)
		return true
	}

	if transaction.VerifySignature(signature) {
		if blockchain.CalculateTotalAmount(transaction.SenderAddress()) < transaction.Value() {
			log.Println("ERROR: Not enough balance in a wallet")
			return false
		}
		blockchain.transactions = append(blockchain.transactions, transaction)
		return true
	} else {
		log.Println("ERROR: Failed to verify transaction")
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
	log.Println("action=mining, status=success")

	for _, neighbor := range blockchain.neighbors {
		neighbor.Consensus()
	}

	return true
}

func (blockchain *Blockchain) StartMining() {
	if blockchain.miningStopped {
		blockchain.miningStopped = false
	} else {
		blockchain.Mine()
		_ = time.AfterFunc(time.Second*MiningTimerSec, blockchain.StartMining)
	}
}

func (blockchain *Blockchain) StopMining() {
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

func (blockchain *Blockchain) ResolveConflicts() bool {
	var longestChain []*Block
	maxLength := len(blockchain.blocks)

	for _, neighbor := range blockchain.neighbors {
		neighborBlocks := neighbor.GetBlocks()
		if len(neighborBlocks) > maxLength && blockchain.IsValid(neighborBlocks) {
			maxLength = len(neighborBlocks)
			longestChain = neighborBlocks
		}
	}

	if longestChain != nil {
		blockchain.blocks = longestChain
		// TODO clear transactions pool here
		log.Println("Conflicts resolved: blockchain replaced")
		return true
	}
	log.Println("Conflicts resolved: blockchain kept")
	return false
}

func (blockchain *Blockchain) createBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, blockchain.transactions)
	blockchain.blocks = append(blockchain.blocks, block)
	blockchain.ClearTransactions()
	for _, neighbor := range blockchain.neighbors {
		// FIXME don't delete transactions if the block is not validated by peers
		neighbor.DeleteTransactions()
	}
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
