package node

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"ruthenium/src/chain"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	MiningDifficulty          = 3
	MiningRewardSenderAddress = "MINING REWARD SENDER ADDRESS"
	MiningReward              = 1.0
	MiningTimerSec            = 30

	StartPort     uint16 = 5000
	EndPort       uint16 = 5002
	StartIpSuffix uint8  = 0
	EndIpSuffix   uint8  = 0

	NeighborSynchronizationTimeSecond = 5
)

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

type Blockchain struct {
	transactions []*chain.Transaction
	blocks       []*chain.Block
	address      string
	mineMutex    sync.Mutex

	ip   string
	port uint16

	neighbors      []*Neighbor
	neighborsMutex sync.Mutex
}

func NewBlockchain(address string, port uint16) *Blockchain {
	blockchain := new(Blockchain)
	blockchain.address = address
	blockchain.port = port
	blockchain.createBlock(0, new(chain.Block).Hash())
	return blockchain
}

func (blockchain *Blockchain) Run() {
	blockchain.StartNeighborsSynchronization()
	blockchain.ResolveConflicts()
	blockchain.StartMining()
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

func (blockchain *Blockchain) FindNeighbors() []*Neighbor {
	address := fmt.Sprintf("%s:%d", blockchain.ip, blockchain.port)

	m := PATTERN.FindStringSubmatch(blockchain.ip)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, err := strconv.Atoi(m[len(m)-1])
	if err != nil {
		fmt.Printf("ERROR: Failed to parse IP %s, err:%v\n", m[len(m)-1], err)
	}
	neighbors := make([]*Neighbor, 0)

	for port := StartPort; port <= EndPort; port += 1 {
		for ipSuffix := StartIpSuffix; ipSuffix <= EndIpSuffix; ipSuffix += 1 {
			guessIp := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ipSuffix))
			neighbor := NewNeighbor(guessIp, port)
			guessTarget := neighbor.IpAndPort()
			if guessTarget != address && neighbor.IsFound() {
				neighbor.StartClient()
				neighbors = append(neighbors, neighbor)
			}
		}
	}
	return neighbors
}

func (blockchain *Blockchain) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Blocks []*chain.Block `json:"blocks"`
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

func (blockchain *Blockchain) CreateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *chain.Signature) bool {
	isTransacted := blockchain.UpdateTransaction(senderAddress, recipientAddress, senderPublicKey, value, signature)

	if isTransacted {
		publicKeyStr := fmt.Sprintf("%064x%064x", senderPublicKey.X.Bytes(), senderPublicKey.Y.Bytes())
		signatureStr := signature.String()
		transactionRequest := &chain.PutTransactionRequest{
			&senderAddress,
			&recipientAddress,
			&publicKeyStr,
			&value,
			&signatureStr}
		//marshaledTransactionRequest, err := json.Marshal(transactionRequest)
		//if err != nil {
		//	log.Println("ERROR: Failed to marshal transaction request for neighbors")
		//}
		for _, neighbor := range blockchain.neighbors {
			//	// TODO extract http logic
			//	endpoint := fmt.Sprintf("http://%s/transactions", neighbor.IpAndPort())
			//	client := &http.Client{}
			//	buffer := bytes.NewBuffer(marshaledTransactionRequest)
			//	request, requestError := http.NewRequest("PUT", endpoint, buffer)
			//	if requestError != nil {
			//		log.Printf("ERROR: %v", requestError)
			//	}
			//	response, responseError := client.Do(request)
			//	if responseError != nil {
			//		log.Printf("ERROR: %v", responseError)
			//	}
			//	log.Printf("%v\n", response)
			neighbor.PutTransactions(transactionRequest)
		}
	}

	return isTransacted
}

func (blockchain *Blockchain) UpdateTransaction(senderAddress string, recipientAddress string, senderPublicKey *ecdsa.PublicKey, value float32, signature *chain.Signature) (isTransacted bool) {
	// FIXME nil private key
	sender := chain.PopWallet(nil, senderPublicKey, senderAddress)
	transaction := chain.NewTransaction(sender.Address(), sender.PublicKey(), recipientAddress, value)
	return blockchain.addTransaction(transaction, signature)
}

func (blockchain *Blockchain) addTransaction(transaction *chain.Transaction, signature *chain.Signature) bool {
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

	// TODO decide if we should get a reward when there is no transaction in the pool
	//if len(blockchain.transactionPool) == 0 {
	//	return false
	//}

	sender := chain.PopWallet(nil, nil, MiningRewardSenderAddress)
	transaction := chain.NewTransaction(sender.Address(), sender.PublicKey(), blockchain.address, MiningReward)
	blockchain.addTransaction(transaction, nil)
	nonce := blockchain.proofOfWork()
	previousHash := blockchain.lastBlock().Hash()
	blockchain.createBlock(nonce, previousHash)
	log.Println("action=mining, status=success")

	for _, neighbor := range blockchain.neighbors {
		//	// TODO extract http logic
		//	endpoint := fmt.Sprintf("http://%s/consensus", neighbor.IpAndPort())
		//	client := &http.Client{}
		//	request, requestError := http.NewRequest("PUT", endpoint, nil)
		//	if requestError != nil {
		//		log.Printf("ERROR: %v", requestError)
		//	}
		//	response, responseError := client.Do(request)
		//	if responseError != nil {
		//		log.Printf("ERROR: %v", responseError)
		//	}
		//	log.Printf("%v\n", response)
		neighbor.Consensus()
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

func (blockchain *Blockchain) Transactions() []*chain.Transaction {
	return blockchain.transactions
}

func (blockchain *Blockchain) Blocks() []*chain.Block {
	// TODO improve copy
	var blocks []*chain.Block
	for _, transaction := range blockchain.blocks {
		blocks = append(blocks, transaction)
	}
	return blocks
}

func (blockchain *Blockchain) ClearTransactions() {
	blockchain.transactions = blockchain.transactions[:0]
}

func (blockchain *Blockchain) IsValid(blocks []*chain.Block) bool {
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
	var longestChain []*chain.Block
	maxLength := len(blockchain.blocks)

	for _, neighbor := range blockchain.neighbors {
		//TODO cleaning
		//endpoint := fmt.Sprintf("http://%s/chain", neighbor.IpAndPort())
		//response, responseError := http.Get(endpoint)
		//if responseError != nil {
		//	log.Printf("ERROR: %v", responseError)
		//} else if response.StatusCode == 200 {
		//	var neighborBlockchain Blockchain
		//	decoder := json.NewDecoder(response.Body)
		//	err := decoder.Decode(&neighborBlockchain)
		//	if err != nil {
		//		log.Printf("ERROR: Failed to decode neighbor blockchain\n%v", responseError)
		//	}
		neighborBlocks := neighbor.ReadBlocks()
		if len(neighborBlocks) > maxLength && blockchain.IsValid(neighborBlocks) {
			maxLength = len(neighborBlocks)
			longestChain = neighborBlocks
		}
		//}
	}

	if longestChain != nil {
		blockchain.blocks = longestChain
		log.Println("Conflicts resolved: blockchain replaced")
		return true
	}
	log.Println("Conflicts resolved: blockchain kept")
	return false
}

func (blockchain *Blockchain) createBlock(nonce int, previousHash [32]byte) *chain.Block {
	block := chain.NewBlock(nonce, previousHash, blockchain.transactions)
	blockchain.blocks = append(blockchain.blocks, block)
	blockchain.ClearTransactions()
	for _, neighbor := range blockchain.neighbors {
		//	// TODO extract http logic
		//	endpoint := fmt.Sprintf("http://%s/transactions", neighbor.IpAndPort())
		//	client := &http.Client{}
		//	// FIXME don't delete transactions if the block is not validated by peers
		//	request, requestError := http.NewRequest("DELETE", endpoint, nil)
		//	if requestError != nil {
		//		log.Printf("ERROR: %v", requestError)
		//	}
		//	response, responseError := client.Do(request)
		//	if responseError != nil {
		//		log.Printf("ERROR: %v", responseError)
		//	}
		//	log.Printf("%v\n", response)
		neighbor.DeleteTransactions()
	}
	return block
}

func (blockchain *Blockchain) lastBlock() *chain.Block {
	return blockchain.blocks[len(blockchain.blocks)-1]
}

func (blockchain *Blockchain) copyTransactions() []*chain.Transaction {
	transactions := make([]*chain.Transaction, 0)
	for _, transaction := range blockchain.transactions {
		transactions = append(transactions,
			chain.NewTransaction(transaction.SenderAddress(),
				transaction.SenderPublicKey(),
				transaction.RecipientAddress(),
				transaction.Value()))
	}
	return transactions
}

func (blockchain *Blockchain) isProofValid(nonce int, previousHash [32]byte, transactions []*chain.Transaction, difficulty int) bool {
	zeros := strings.Repeat("0", difficulty)
	guessBlock := chain.NewBlock(nonce, previousHash, transactions)
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
