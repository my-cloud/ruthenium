package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"ruthenium/src/chain"
	"strconv"
)

const templateDir = "src/wallet_server/templates"

type WalletServer struct {
	port    uint16
	gateway string
}

func NewWalletServer(port uint16, gateway string) *WalletServer {
	return &WalletServer{port, gateway}
}

func (walletServer *WalletServer) Port() uint16 {
	return walletServer.port
}

func (walletServer *WalletServer) Gateway() string {
	return walletServer.gateway
}

func (walletServer *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(templateDir, "index.html"))
		if err != nil {
			log.Printf("ERROR: Failed to parse the template")
		} else if err := t.Execute(w, ""); err != nil {
			log.Printf("ERROR: Failed to execute the template")
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) Wallet(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		writer.Header().Add("Content-Type", "application/json")
		wallet := chain.NewWallet()
		marshaledWallet, err := wallet.MarshalJSON()
		if err != nil {
			log.Println("ERROR: Failed to marshal wallet")
		}
		i, err := io.WriteString(writer, string(marshaledWallet[:]))
		if err != nil || i == 0 {
			log.Printf("ERROR: Failed to write wallet")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) CreateTransaction(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		decoder := json.NewDecoder(req.Body)
		var transactionRequest TransactionRequest
		err := decoder.Decode(&transactionRequest)
		if err != nil {
			log.Printf("ERROR: %v", err)
			i, err := io.WriteString(writer, NewStatus("fail").StringValue())
			if err != nil || i == 0 {
				log.Printf("ERROR: Failed to write status")
				return
			}
		}
		if transactionRequest.IsInvalid() {
			log.Println("ERROR: missing fields in transaction request")
			i, err := io.WriteString(writer, NewStatus("fail").StringValue())
			if err != nil || i == 0 {
				log.Printf("ERROR: Failed to write status")
				return
			}
		}

		publicKey := chain.NewPublicKey(*transactionRequest.SenderPublicKey)
		privateKey := chain.NewPublicKey(*transactionRequest.SenderPrivateKey)
		value, err := strconv.ParseFloat(*transactionRequest.Value, 32)
		if err != nil {
			log.Println("ERROR: transaction value parsing error")
			i, err := io.WriteString(writer, NewStatus("fail").StringValue())
			if err != nil || i == 0 {
				log.Printf("ERROR: Failed to write status")
				return
			}
		}
		value32 := float32(value)

		fmt.Println(publicKey)
		fmt.Println(privateKey)
		fmt.Printf("%.1f\n", value32)

		writer.Header().Add("Content-type", "application/json")

	default:
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) Run() {
	http.HandleFunc("/", walletServer.Index)
	http.HandleFunc("/wallet", walletServer.Wallet)
	http.HandleFunc("/transaction", walletServer.CreateTransaction)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(walletServer.Port())), nil))
}
