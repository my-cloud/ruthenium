package main

import (
	"io"
	"log"
	"net/http"
	"ruthenium/src/chain"
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
			log.Printf("ERROR: Failed to marshal blockchain")
		}
		i, err := io.WriteString(writer, string(marshaledBlockchain[:]))
		if err != nil || i == 0 {
			log.Printf("ERROR: Failed to write blockchain")
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (blockchainServer *BlockchainServer) Run() {
	http.HandleFunc("/", blockchainServer.GetChain)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(blockchainServer.port)), nil))
}
