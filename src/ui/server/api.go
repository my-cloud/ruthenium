package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"ruthenium/src/log"
	"ruthenium/src/node/authentication"
	"ruthenium/src/node/blockchain/mine"
	"ruthenium/src/node/neighborhood"
	"strconv"
)

type Api struct {
	publicKey        string
	privateKey       string
	port             uint16
	blockchainClient *neighborhood.Node
	templatesPath    string
	logger           *log.Logger
}

func NewApi(publicKey string, privateKey string, port uint16, hostIp string, hostPort uint16, templatesPath string, level log.Level) *Api {
	logger := log.NewLogger(level)
	blockchainClient := neighborhood.NewNode(hostIp, hostPort, logger)
	return &Api{publicKey, privateKey, port, blockchainClient, templatesPath, logger}
}

func (api *Api) Port() uint16 {
	return api.port
}

func (api *Api) BlockchainClient() *neighborhood.Node {
	return api.blockchainClient
}

func (api *Api) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(api.templatesPath, "index.html"))
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to parse the template: %w", err).Error())
			return
		}
		if err = t.Execute(w, ""); err != nil {
			api.logger.Error(fmt.Errorf("failed to execute the template: %w", err).Error())
		}
	default:
		api.logger.Error("invalid HTTP method")
	}
}

func (api *Api) Wallet(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		wallet, err := authentication.NewWallet(api.publicKey, api.privateKey)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to create wallet: %w", err).Error())
			return
		}
		marshaledWallet, err := wallet.MarshalJSON()
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to marshal wallet: %w", err).Error())
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		api.write(writer, string(marshaledWallet[:]))
	default:
		writer.WriteHeader(http.StatusBadRequest)
		api.logger.Error("invalid HTTP method")
	}
}

func (api *Api) CreateTransaction(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var transactionRequest TransactionRequest
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&transactionRequest)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to decode transaction request: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		if transactionRequest.IsInvalid() {
			api.logger.Error("field(s) are missing in transaction request")
			api.write(writer, "fail")
			return
		}
		publicKey, err := authentication.NewPublicKey(*transactionRequest.SenderPublicKey)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to decode transaction public key: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		privateKey, err := authentication.NewPrivateKey(*transactionRequest.SenderPrivateKey, publicKey)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to decode transaction private key: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		value, err := strconv.ParseFloat(*transactionRequest.Value, 32)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to parse transaction value: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		value32 := float32(value)
		transaction := mine.NewTransaction(*transactionRequest.SenderAddress, *transactionRequest.RecipientAddress, value32)
		marshaledTransaction, err := json.Marshal(transaction)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to marshal transaction: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		signature, err := authentication.NewSignature(marshaledTransaction, privateKey)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to generate signature: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		signatureString := signature.String()
		var verb = neighborhood.POST
		blockchainTransactionRequest := neighborhood.TransactionRequest{
			Verb:             &verb,
			SenderAddress:    transactionRequest.SenderAddress,
			RecipientAddress: transactionRequest.RecipientAddress,
			SenderPublicKey:  transactionRequest.SenderPublicKey,
			Value:            &value32,
			Signature:        &signatureString,
		}
		err = api.blockchainClient.AddTransaction(blockchainTransactionRequest)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		api.write(writer, "success")
	default:
		writer.WriteHeader(http.StatusBadRequest)
		api.logger.Error("invalid HTTP method")
	}
}

func (api *Api) GetTransactions(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := api.blockchainClient.GetTransactions()
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		var marshaledTransactions []byte
		marshaledTransactions, err = json.Marshal(transactions)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to marshal transactions: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		api.write(writer, string(marshaledTransactions[:]))
	default:
		api.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (api *Api) Mine(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := api.blockchainClient.Mine()
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to mine: %w", err).Error())
			api.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		api.logger.Error("invalid HTTP method")
	}
}

func (api *Api) StartMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := api.blockchainClient.StartMining()
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to start mining: %w", err).Error())
			api.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		api.logger.Error("invalid HTTP method")
	}
}

func (api *Api) StopMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := api.blockchainClient.StopMining()
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to stop mining: %w", err).Error())
			api.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		api.logger.Error("invalid HTTP method")
	}
}

func (api *Api) WalletAmount(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		amountRequest := neighborhood.AmountRequest{
			Address: &address,
		}
		if amountRequest.IsInvalid() {
			api.logger.Error("field(s) are missing in amount request")
			api.write(writer, "fail")
			return
		}
		amountResponse, err := api.blockchainClient.GetAmount(amountRequest)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to get amountResponse: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		var marshaledAmount []byte
		marshaledAmount, err = json.Marshal(amountResponse)
		if err != nil {
			api.logger.Error(fmt.Errorf("failed to marshal amountResponse: %w", err).Error())
			api.write(writer, "fail")
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		api.write(writer, string(marshaledAmount[:]))
	default:
		api.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (api *Api) Run() {
	http.HandleFunc("/", api.Index)
	http.HandleFunc("/wallet", api.Wallet)
	http.HandleFunc("/transaction", api.CreateTransaction)
	http.HandleFunc("/transactions", api.GetTransactions)
	http.HandleFunc("/wallet/amount", api.WalletAmount)
	http.HandleFunc("/mine", api.Mine)
	http.HandleFunc("/mine/start", api.StartMining)
	http.HandleFunc("/mine/stop", api.StopMining)
	api.logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(api.Port())), nil).Error())
}

func (api *Api) write(writer http.ResponseWriter, message string) {
	i, err := io.WriteString(writer, message)
	if err != nil || i == 0 {
		api.logger.Error(fmt.Errorf("failed to write message: %s", message).Error())
	}
}
