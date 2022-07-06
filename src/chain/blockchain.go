package chain

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

const (
	MiningDifficulty = 3
	MiningSender     = "THE BLOCKCHAIN"
	MiningReward     = 1.0
)

type Blockchain struct {
	transactions []*Transaction
	blocks       []*Block
	address      string
	port         uint16
}

func NewBlockchain(address string, port uint16) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.address = address
	blockchain.port = port
	blockchain.createBlock(0, new(Block).Hash())
	return blockchain
}

func (blockchain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blocks"`
	}{
		Blocks: blockchain.blocks,
	})
}

func (blockchain *Blockchain) Print() {
	for i, block := range blockchain.blocks {
		fmt.Printf("%s Block  %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

func (blockchain *Blockchain) CreateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) bool {
	// FIXME nil private key
	sender := &Wallet{nil, senderPublicKey, senderAddress}
	transaction := NewTransaction(sender, recipientAddress, value)
	isTransacted := blockchain.AddTransaction(transaction, signature)

	// TODO
	// Sync

	return isTransacted
}

func (blockchain *Blockchain) AddTransaction(transaction *Transaction, signature *Signature) bool {
	if transaction.Sender().Address() == MiningSender {
		blockchain.transactions = append(blockchain.transactions, transaction)
		return true
	}

	if blockchain.verifyTransactionSignature(transaction.Sender().PublicKey(), signature, transaction) {
		/*
			if blockchain.CalculateTotalAmount(sender) < value {
				log.Println("ERROR: Not enough balance in a wallet")
				return false
			}
		*/
		blockchain.transactions = append(blockchain.transactions, transaction)
		return true
	} else {
		log.Println("ERROR: Failed to verify transaction")
	}
	return false

}

func (blockchain *Blockchain) Mining() bool {
	transaction := NewTransaction(&Wallet{nil, nil, MiningSender}, blockchain.address, MiningReward)
	blockchain.AddTransaction(transaction, nil)
	nonce := blockchain.proofOfWork()
	previousHash := blockchain.lastBlock().Hash()
	blockchain.createBlock(nonce, previousHash)
	log.Println("action=mining, status=success")
	return true
}

func (blockchain *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, block := range blockchain.blocks {
		for _, transaction := range block.transactions {
			value := transaction.Value()
			if blockchainAddress == transaction.RecipientAddress() {
				totalAmount += value
			}

			if blockchainAddress == transaction.Sender().Address() {
				totalAmount -= value
			}
		}
	}
	return totalAmount
}

func (blockchain *Blockchain) Transactions() []*Transaction {
	return blockchain.transactions
}

func (blockchain *Blockchain) createBlock(nonce int, previousHash [32]byte) *Block {
	block := NewBlock(nonce, previousHash, blockchain.transactions)
	blockchain.blocks = append(blockchain.blocks, block)
	blockchain.transactions = []*Transaction{}
	return block
}

func (blockchain *Blockchain) lastBlock() *Block {
	return blockchain.blocks[len(blockchain.blocks)-1]
}

func (blockchain *Blockchain) verifyTransactionSignature(
	senderPublicKey *ecdsa.PublicKey, signature *Signature, t *Transaction) bool {
	marshaledBlockchain, err := json.Marshal(t)
	if err != nil {
		log.Println("ERROR: Failed to marshal blockchain")
	}
	hash := sha256.Sum256(marshaledBlockchain)
	return ecdsa.Verify(senderPublicKey, hash[:], signature.r, signature.s)
}

func (blockchain *Blockchain) copyTransactions() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range blockchain.transactions {
		transactions = append(transactions,
			NewTransaction(transaction.Sender(),
				transaction.RecipientAddress(),
				transaction.Value()))
	}
	return transactions
}

func (blockchain *Blockchain) validProof(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
	guessHashStr := fmt.Sprintf("%x", guessBlock.Hash())
	return guessHashStr[:difficulty] == zeros
}

func (blockchain *Blockchain) proofOfWork() int {
	transactions := blockchain.copyTransactions()
	previousHash := blockchain.lastBlock().Hash()
	var nonce int
	for !blockchain.validProof(nonce, previousHash, transactions, MiningDifficulty) {
		nonce++
	}
	return nonce
}
