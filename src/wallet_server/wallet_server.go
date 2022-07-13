package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"ruthenium/src/chain"
	"ruthenium/src/rest"
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
			log.Println("ERROR: Failed to parse the template")
		} else if err := t.Execute(w, ""); err != nil {
			log.Println("ERROR: Failed to execute the template")
		}
	default:
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) Wallet(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		writer.Header().Add("Content-Type", "application/json")
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			fmt.Printf("ERROR: Failed to generate private key, err%v\n", err)
		} else {
			wallet := chain.NewWallet(privateKey)
			marshaledWallet, err := wallet.MarshalJSON()
			if err != nil {
				log.Println("ERROR: Failed to marshal wallet")
			}
			i, err := io.WriteString(writer, string(marshaledWallet[:]))
			if err != nil || i == 0 {
				log.Println("ERROR: Failed to write wallet")
			}
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
		jsonWriter := rest.NewJsonWriter(writer)
		if err != nil {
			log.Printf("ERROR: %v", err)
			jsonWriter.WriteStatus("fail")
		}
		if transactionRequest.IsInvalid() {
			log.Println("ERROR: Field(s) are missing in transaction request to wallet server")
			jsonWriter.WriteStatus("fail")
		}

		publicKey := chain.NewPublicKey(*transactionRequest.SenderPublicKey)
		privateKey := chain.NewPrivateKey(*transactionRequest.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*transactionRequest.Value, 32)
		if err != nil {
			log.Println("ERROR: Failed to parse transaction value")
			jsonWriter.WriteStatus("fail")
		}
		value32 := float32(value)

		fmt.Println(publicKey)
		fmt.Println(privateKey)
		fmt.Printf("%.1f\n", value32)

		writer.Header().Add("Content-type", "application/json")

		sender := chain.PopWallet(privateKey, publicKey, *transactionRequest.SenderAddress)
		transaction := chain.NewTransaction(sender.Address(), sender.PublicKey(), *transactionRequest.RecipientAddress, value32)
		signature := chain.NewSignature(transaction, privateKey)

		signatureString := signature.String()
		marshaledTransaction, err := json.Marshal(&chain.TransactionRequest{
			SenderAddress:    transactionRequest.SenderAddress,
			RecipientAddress: transactionRequest.RecipientAddress,
			SenderPublicKey:  transactionRequest.SenderPublicKey,
			Value:            &value32,
			Signature:        &signatureString,
		})
		if err != nil {
			log.Println("ERROR: Failed to marshal transaction")
		}

		buffer := bytes.NewBuffer(marshaledTransaction)

		response, err := http.Post(walletServer.Gateway()+"/transactions", "application/json", buffer)
		if err != nil {
			log.Printf("ERROR: %v", err)
		} else if response.StatusCode == 201 {
			jsonWriter.WriteStatus("success")
			return
		}
		log.Printf("ERROR: status code %d", response.StatusCode)
		jsonWriter.WriteStatus("fail")
	default:
		writer.WriteHeader(http.StatusBadRequest)
		log.Println("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) WalletAmount(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		endpoint := fmt.Sprintf("%s/amount", walletServer.Gateway())

		jsonWriter := rest.NewJsonWriter(writer)
		bcsReq, err := http.NewRequest("GET", endpoint, nil)
		if err != nil {
			log.Printf("ERROR: %v", err)
			jsonWriter.WriteStatus("fail")
			return
		}
		q := bcsReq.URL.Query()
		q.Add("address", address)
		bcsReq.URL.RawQuery = q.Encode()

		client := &http.Client{}
		bcsResp, err := client.Do(bcsReq)
		if err != nil {
			log.Printf("ERROR: %v", err)
			jsonWriter.WriteStatus("fail")
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		if bcsResp.StatusCode == 200 {
			decoder := json.NewDecoder(bcsResp.Body)
			var bar chain.AmountResponse
			err := decoder.Decode(&bar)
			if err != nil {
				log.Printf("ERROR: %v", err)
				jsonWriter.WriteStatus("fail")
				return
			}

			marshaledAmount, err := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  bar.Amount,
			})
			if err != nil {
				log.Println("ERROR: Failed to marshal amount")
			}
			i, err := io.WriteString(writer, string(marshaledAmount[:]))
			if err != nil || i == 0 {
				log.Println("ERROR: Failed to write amount")
			}
		} else {
			jsonWriter.WriteStatus("fail")
		}
	default:
		log.Printf("ERROR: Invalid HTTP Method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (walletServer *WalletServer) Run() {
	http.HandleFunc("/", walletServer.Index)
	http.HandleFunc("/wallet", walletServer.Wallet)
	http.HandleFunc("/transaction", walletServer.CreateTransaction)
	http.HandleFunc("/wallet/amount", walletServer.WalletAmount)
	log.Fatal(http.ListenAndServe("localhost:"+strconv.Itoa(int(walletServer.Port())), nil))
}
