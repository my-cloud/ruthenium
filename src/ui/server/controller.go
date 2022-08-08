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
	"ruthenium/src/node/blockchain/mining"
	"ruthenium/src/node/neighborhood"
	"strconv"
)

type Controller struct {
	publicKey        string
	privateKey       string
	port             uint16
	blockchainClient *neighborhood.Neighbor
	templatesPath    string
	logger           *log.Logger
}

func NewController(publicKey string, privateKey string, port uint16, hostIp string, hostPort uint16, templatesPath string, level log.Level) *Controller {
	logger := log.NewLogger(level)
	blockchainClient := neighborhood.NewNeighbor(hostIp, hostPort, logger)
	return &Controller{publicKey, privateKey, port, blockchainClient, templatesPath, logger}
}

func (controller *Controller) Port() uint16 {
	return controller.port
}

func (controller *Controller) BlockchainClient() *neighborhood.Neighbor {
	return controller.blockchainClient
}

func (controller *Controller) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(controller.templatesPath, "index.html"))
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to parse the template: %w", err).Error())
			return
		}
		if err = t.Execute(w, ""); err != nil {
			controller.logger.Error(fmt.Errorf("failed to execute the template: %w", err).Error())
		}
	default:
		controller.logger.Error("invalid HTTP method")
	}
}

func (controller *Controller) Wallet(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		wallet, err := authentication.NewWallet(controller.publicKey, controller.privateKey)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to create wallet: %w", err).Error())
			return
		}
		marshaledWallet, err := wallet.MarshalJSON()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to marshal wallet: %w", err).Error())
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		controller.write(writer, string(marshaledWallet[:]))
	default:
		writer.WriteHeader(http.StatusBadRequest)
		controller.logger.Error("invalid HTTP method")
	}
}

func (controller *Controller) CreateTransaction(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var transactionRequest TransactionRequest
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&transactionRequest)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to decode transaction request: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		if transactionRequest.IsInvalid() {
			controller.logger.Error("field(s) are missing in transaction request")
			controller.write(writer, "fail")
			return
		}
		publicKey, err := authentication.NewPublicKey(*transactionRequest.SenderPublicKey)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to decode transaction public key: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		privateKey, err := authentication.NewPrivateKey(*transactionRequest.SenderPrivateKey, publicKey)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to decode transaction private key: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		value, err := strconv.ParseFloat(*transactionRequest.Value, 32)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to parse transaction value: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		value32 := float32(value)
		transaction := mining.NewTransaction(*transactionRequest.SenderAddress, *transactionRequest.RecipientAddress, value32)
		marshaledTransaction, err := json.Marshal(transaction)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to marshal transaction: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		signature, err := authentication.NewSignature(marshaledTransaction, privateKey)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to generate signature: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		signatureString := signature.String()
		var verb = neighborhood.POST
		timestamp := transaction.Timestamp()
		blockchainTransactionRequest := neighborhood.TransactionRequest{
			Verb:             &verb,
			Timestamp:        &timestamp,
			SenderAddress:    transactionRequest.SenderAddress,
			RecipientAddress: transactionRequest.RecipientAddress,
			SenderPublicKey:  transactionRequest.SenderPublicKey,
			Value:            &value32,
			Signature:        &signatureString,
		}
		err = controller.blockchainClient.AddTransaction(blockchainTransactionRequest)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		controller.write(writer, "success")
	default:
		writer.WriteHeader(http.StatusBadRequest)
		controller.logger.Error("invalid HTTP method")
	}
}

func (controller *Controller) GetTransactions(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := controller.blockchainClient.GetTransactions()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		var marshaledTransactions []byte
		marshaledTransactions, err = json.Marshal(transactions)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to marshal transactions: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		controller.write(writer, string(marshaledTransactions[:]))
	default:
		controller.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (controller *Controller) Mine(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := controller.blockchainClient.Mine()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to mine: %w", err).Error())
			controller.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		controller.logger.Error("invalid HTTP method")
	}
}

func (controller *Controller) StartMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := controller.blockchainClient.StartMining()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to start mining: %w", err).Error())
			controller.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		controller.logger.Error("invalid HTTP method")
	}
}

func (controller *Controller) StopMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := controller.blockchainClient.StopMining()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to stop mining: %w", err).Error())
			controller.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		controller.logger.Error("invalid HTTP method")
	}
}

func (controller *Controller) WalletAmount(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		amountRequest := neighborhood.AmountRequest{
			Address: &address,
		}
		if amountRequest.IsInvalid() {
			controller.logger.Error("field(s) are missing in amount request")
			controller.write(writer, "fail")
			return
		}
		amountResponse, err := controller.blockchainClient.GetAmount(amountRequest)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to get amountResponse: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		var marshaledAmount []byte
		marshaledAmount, err = json.Marshal(amountResponse)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to marshal amountResponse: %w", err).Error())
			controller.write(writer, "fail")
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		controller.write(writer, string(marshaledAmount[:]))
	default:
		controller.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (controller *Controller) Run() {
	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/wallet", controller.Wallet)
	http.HandleFunc("/transaction", controller.CreateTransaction)
	http.HandleFunc("/transactions", controller.GetTransactions)
	http.HandleFunc("/wallet/amount", controller.WalletAmount)
	http.HandleFunc("/mine", controller.Mine)
	http.HandleFunc("/mine/start", controller.StartMining)
	http.HandleFunc("/mine/stop", controller.StopMining)
	controller.logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(controller.Port())), nil).Error())
}

func (controller *Controller) write(writer http.ResponseWriter, message string) {
	i, err := io.WriteString(writer, message)
	if err != nil || i == 0 {
		controller.logger.Error(fmt.Errorf("failed to write message: %s", message).Error())
	}
}
