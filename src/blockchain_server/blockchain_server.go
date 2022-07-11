package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"ruthenium/src/chain"
	"ruthenium/src/rest"
	"strconv"
)

var cache = make(map[string]*chain.Blockchain)

type BlockchainServer struct {
	port uint16
}

func NewBlockchainServer(port uint16) *BlockchainServer {
	return &BlockchainServer{port}
}

func (blockchainServer *BlockchainServer) Port() uint16 {
	return blockchainServer.port
}

func (blockchainServer *BlockchainServer) GetBlockchain() *chain.Blockchain {
	blockchain, ok := cache["blockchain"]
	if !ok {
		minerWallet := chain.NewWallet()
		blockchain = chain.NewBlockchain(minerWallet.Address(), blockchainServer.Port())
		cache["blockchain"] = blockchain
	}
	return blockchain
}

func (blockchainServer *BlockchainServer) GetChain(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodGet:
		writer.Header().Add("Content-type", "application/json")
		blockchain := blockchainServer.GetBlockchain()
		marshaledBlockchain, err := blockchain.MarshalJSON()
		if err != nil {
			log.Println("ERROR: Failed to marshal blockchain")
		}
		i, err := io.WriteString(writer, string(marshaledBlockchain[:]))
		if err != nil || i == 0 {
			log.Println("ERROR: Failed to write blockchain")
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (blockchainServer *BlockchainServer) Transactions(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		writer.Header().Add("Content-Type", "application/json")
		blockchain := blockchainServer.GetBlockchain()
		transactions := blockchain.Transactions()
		marshaledTransactions, err := json.Marshal(struct {
			Transactions []*chain.Transaction `json:"transactions"`
			Length       int                  `json:"length"`
		}{
			Transactions: transactions,
			Length:       len(transactions),
		})
		if err != nil {
			log.Println("ERROR: Failed to marshal transactions")
		}
		i, err := io.WriteString(writer, string(marshaledTransactions[:]))
		if err != nil || i == 0 {
			log.Println("ERROR: Failed to write transactions")
		}
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var transactionRequest chain.TransactionRequest
		err := decoder.Decode(&transactionRequest)
		jsonWriter := rest.NewJsonWriter(writer)
		if err != nil {
			log.Printf("ERROR: %v", err)
			jsonWriter.WriteStatus("fail")
			return
		}
		if transactionRequest.IsInvalid() {
			log.Println("ERROR: Field(s) are missing in transaction request to blockchain server")
			jsonWriter.WriteStatus("fail")
			return
		}
		publicKey := chain.NewPublicKey(*transactionRequest.SenderPublicKey)
		signature := chain.DecodeSignature(*transactionRequest.Signature)
		blockchain := blockchainServer.GetBlockchain()
		isCreated := blockchain.CreateTransaction(*transactionRequest.SenderAddress, *transactionRequest.RecipientAddress, publicKey, *transactionRequest.Value, signature)

		writer.Header().Add("Content-Type", "application/json")
		if !isCreated {
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.WriteStatus("fail")
		} else {
			writer.WriteHeader(http.StatusCreated)
			jsonWriter.WriteStatus("success")
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (blockchainServer *BlockchainServer) Mine(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchain := blockchainServer.GetBlockchain()
		isMined := blockchain.Mine()

		jsonWriter := rest.NewJsonWriter(writer)
		writer.Header().Add("Content-Type", "application/json")
		if isMined {
			jsonWriter.WriteStatus("success")
		} else {
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.WriteStatus("fail")
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (blockchainServer *BlockchainServer) StartMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchain := blockchainServer.GetBlockchain()
		blockchain.StartMining()

		jsonWriter := rest.NewJsonWriter(writer)
		writer.Header().Add("Content-Type", "application/json")
		jsonWriter.WriteStatus("success")
	default:
		log.Println("ERROR: Invalid HTTP Method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (blockchainServer *BlockchainServer) Amount(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		blockchainAddress := req.URL.Query().Get("address")
		amount := blockchainServer.GetBlockchain().CalculateTotalAmount(blockchainAddress)

		amountResponse := chain.NewAmountResponse(amount)
		marshaledAmountResponse, err := amountResponse.MarshalJSON()
		if err != nil {
			log.Println("ERROR: Failed to marshal amount response")
		}

		w.Header().Add("Content-Type", "application/json")
		i, err := io.WriteString(w, string(marshaledAmountResponse[:]))
		if err != nil || i == 0 {
			log.Println("ERROR: Failed to write amount response")
		}

	default:
		log.Println("ERROR: Invalid HTTP Method")
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (blockchainServer *BlockchainServer) Run() {
	blockchainServer.GetBlockchain().Run()

	http.HandleFunc("/", blockchainServer.GetChain)
	http.HandleFunc("/transactions", blockchainServer.Transactions)
	http.HandleFunc("/mine", blockchainServer.Mine)
	http.HandleFunc("/mine/start", blockchainServer.StartMining)
	http.HandleFunc("/amount", blockchainServer.Amount)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(blockchainServer.port)), nil))
}
