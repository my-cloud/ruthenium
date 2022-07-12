package chain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	MiningDifficulty          = 3
	MiningRewardSenderAddress = "MINING REWARD SENDER ADDRESS"
	MiningReward              = 1.0
	MiningTimerSec            = 20

	NeighborSynchronizationTimeSecond = 5
)

type Blockchain struct {
	transactions []*Transaction
	blocks       []*Block
	address      string
	mineMutex    sync.Mutex

	neighbors      []string
	neighborsMutex sync.Mutex
	hostNode       *Node
}

func NewBlockchain(address string, port uint16) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.address = address
	blockchain.hostNode = NewHostNode(port)
	blockchain.createBlock(0, new(Block).Hash())
	return blockchain
}

func (blockchain *Blockchain) Run() {
	blockchain.StartNeighborsSynchronization()
	blockchain.ResolveConflicts()
}

func (blockchain *Blockchain) SynchronizeNeighbors() {
	blockchain.neighborsMutex.Lock()
	defer blockchain.neighborsMutex.Unlock()
	blockchain.neighbors = blockchain.hostNode.FindNeighbors()
}

func (blockchain *Blockchain) StartNeighborsSynchronization() {
	blockchain.SynchronizeNeighbors()
	_ = time.AfterFunc(time.Second*NeighborSynchronizationTimeSecond, blockchain.StartNeighborsSynchronization)
}

func (blockchain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*Block `json:"blocks"`
	}{
		Blocks: blockchain.blocks,
	})
}

func (blockchain *Blockchain) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &struct {
		Blocks *[]*Block `json:"blocks"`
	}{
		Blocks: &blockchain.blocks,
	}); err != nil {
		return err
	}
	return nil
}

func (blockchain *Blockchain) Print() {
	for i, block := range blockchain.blocks {
		fmt.Printf("%s Block  %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
	fmt.Printf("%s\n", strings.Repeat("*", 60))
}

func (blockchain *Blockchain) CreateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) bool {
	isTransacted := blockchain.UpdateTransaction(senderAddress, recipientAddress, senderPublicKey, value, signature)

	if isTransacted {
		publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
		signatureStr := signature.String()
		transactionRequest := &TransactionRequest{
			&senderAddress,
			&recipientAddress,
			&publicKeyStr,
			&value,
			&signatureStr}
		marshaledTransactionRequest, err := json.Marshal(transactionRequest)
		if err != nil {
			log.Println("ERROR: Failed to marshal transaction request for neighbors")
		}
		for _, neighbor := range blockchain.neighbors {
			// TODO extract http logic
			endpoint := fmt.Sprintf("http://%s/transactions", neighbor)
			client := &http.Client{}
			buffer := bytes.NewBuffer(marshaledTransactionRequest)
			request, requestError := http.NewRequest("PUT", endpoint, buffer)
			if requestError != nil {
				log.Printf("ERROR: %v", requestError)
			}
			response, responseError := client.Do(request)
			if responseError != nil {
				log.Printf("ERROR: %v", responseError)
			}
			log.Printf("%v\n", response)
		}
	}

	return isTransacted
}

func (blockchain *Blockchain) UpdateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *Signature) (isTransacted bool) {
	// FIXME nil private key
	sender := &Wallet{nil, senderPublicKey, senderAddress}
	transaction := NewTransaction(sender.Address(), sender.PublicKey(), recipientAddress, value)
	return blockchain.addTransaction(transaction, signature)
}

func (blockchain *Blockchain) addTransaction(transaction *Transaction, signature *Signature) bool {
	if transaction.SenderAddress() == MiningRewardSenderAddress {
		blockchain.transactions = append(blockchain.transactions, transaction)
		return true
	}

	if blockchain.verifyTransactionSignature(transaction.SenderPublicKey(), signature, transaction) {
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

func (blockchain *Blockchain) Mine() bool {
	blockchain.mineMutex.Lock()
	defer blockchain.mineMutex.Unlock()

	//if len(blockchain.transactionPool) == 0 {
	//	return false
	//}

	sender := &Wallet{nil, nil, MiningRewardSenderAddress}
	transaction := NewTransaction(sender.Address(), sender.PublicKey(), blockchain.address, MiningReward)
	blockchain.addTransaction(transaction, nil)
	nonce := blockchain.proofOfWork()
	previousHash := blockchain.lastBlock().Hash()
	blockchain.createBlock(nonce, previousHash)
	log.Println("action=mining, status=success")

	for _, neighbor := range blockchain.neighbors {
		// TODO extract http logic
		endpoint := fmt.Sprintf("http://%s/consensus", neighbor)
		client := &http.Client{}
		request, requestError := http.NewRequest("PUT", endpoint, nil)
		if requestError != nil {
			log.Printf("ERROR: %v", requestError)
		}
		response, responseError := client.Do(request)
		if responseError != nil {
			log.Printf("ERROR: %v", responseError)
		}
		log.Printf("%v\n", response)
	}

	return true
}

func (blockchain *Blockchain) StartMining() {
	blockchain.Mine()
	_ = time.AfterFunc(time.Second*MiningTimerSec, blockchain.StartMining)
}

func (blockchain *Blockchain) CalculateTotalAmount(blockchainAddress string) float32 {
	var totalAmount float32 = 0.0
	for _, block := range blockchain.blocks {
		for _, transaction := range block.transactions {
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
	blockchain.transactions = blockchain.transactions[:0]
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
		// TODO extract http logic
		endpoint := fmt.Sprintf("http://%s/chain", neighbor)
		response, responseError := http.Get(endpoint)
		if responseError != nil {
			log.Printf("ERROR: %v", responseError)
		} else if response.StatusCode == 200 {
			var neighborBlockchain Blockchain
			decoder := json.NewDecoder(response.Body)
			err := decoder.Decode(&neighborBlockchain)
			if err != nil {
				log.Printf("ERROR: Failed to decode neighbor blockchain\n%v", responseError)
			}
			neighborBlocks := neighborBlockchain.Blocks()
			if len(neighborBlocks) > maxLength && blockchain.IsValid(neighborBlocks) {
				maxLength = len(neighborBlocks)
				longestChain = neighborBlocks
			}
		}
	}

	if longestChain != nil {
		blockchain.blocks = longestChain
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
		// TODO extract http logic
		endpoint := fmt.Sprintf("http://%s/transactions", neighbor)
		client := &http.Client{}
		request, requestError := http.NewRequest("DELETE", endpoint, nil)
		if requestError != nil {
			log.Printf("ERROR: %v", requestError)
		}
		response, responseError := client.Do(request)
		if responseError != nil {
			log.Printf("ERROR: %v", responseError)
		}
		log.Printf("%v\n", response)
	}
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
			NewTransaction(transaction.SenderAddress(),
				transaction.SenderPublicKey(),
				transaction.RecipientAddress(),
				transaction.Value()))
	}
	return transactions
}

func (blockchain *Blockchain) isProofValid(nonce int, previousHash [32]byte, transactions []*Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := Block{0, nonce, previousHash, transactions}
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
