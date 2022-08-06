package wallet

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"ruthenium/src/chain"
	"ruthenium/src/log"
	"strconv"
)

const templateDir = "src/wallet_server/templates"

type WalletServer struct {
	publicKey        string
	privateKey       string
	port             uint16
	blockchainClient *chain.Node
	logger           *log.Logger
}

func NewWalletServer(publicKey string, privateKey string, port uint16, hostIp string, hostPort uint16, level log.Level) *WalletServer {
	logger := log.NewLogger(level)
	blockchainClient := chain.NewNode(hostIp, hostPort, logger)
	return &WalletServer{publicKey, privateKey, port, blockchainClient, logger}
}

func (walletServer *WalletServer) Port() uint16 {
	return walletServer.port
}

func (walletServer *WalletServer) BlockchainClient() *chain.Node {
	return walletServer.blockchainClient
}

func (walletServer *WalletServer) Index(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		t, err := template.ParseFiles(path.Join(templateDir, "index.html"))
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to parse the template: %w", err).Error())
			return
		}
		if err = t.Execute(w, ""); err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to execute the template: %w", err).Error())
		}
	default:
		walletServer.logger.Error("invalid HTTP method")
	}
}

func (walletServer *WalletServer) Wallet(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		wallet, err := chain.NewWallet(walletServer.publicKey, walletServer.privateKey)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to create wallet: %w", err).Error())
			return
		}
		marshaledWallet, err := wallet.MarshalJSON()
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to marshal wallet: %w", err).Error())
		}
		writer.Header().Add("Content-Type", "application/json")
		walletServer.write(writer, string(marshaledWallet[:]))
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("invalid HTTP method")
	}
}

func (walletServer *WalletServer) CreateTransaction(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var transactionRequest TransactionRequest
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&transactionRequest)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to decode transaction request: %w", err).Error())
			walletServer.write(writer, "fail")
			return
		}
		if transactionRequest.IsInvalid() {
			walletServer.logger.Error("field(s) are missing in transaction request to wallet server")
			walletServer.write(writer, "fail")
			return
		}
		publicKey, err := chain.NewPublicKey(*transactionRequest.SenderPublicKey)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to decode transaction public key: %w", err).Error())
			return
		}
		privateKey, err := chain.NewPrivateKey(*transactionRequest.SenderPrivateKey, publicKey)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to decode transaction private key: %w", err).Error())
			return
		}
		value, err := strconv.ParseFloat(*transactionRequest.Value, 32)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to parse transaction value: %w", err).Error())
			walletServer.write(writer, "fail")
			return
		}
		value32 := float32(value)
		transaction := chain.NewTransaction(publicKey, *transactionRequest.SenderAddress, *transactionRequest.RecipientAddress, value32, walletServer.logger)
		signature, err := chain.NewSignature(transaction, privateKey)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to generate signature: %w", err).Error())
		}
		signatureString := signature.String()
		var verb = chain.POST
		blockchainTransactionRequest := chain.TransactionRequest{
			Verb:             &verb,
			SenderAddress:    transactionRequest.SenderAddress,
			RecipientAddress: transactionRequest.RecipientAddress,
			SenderPublicKey:  transactionRequest.SenderPublicKey,
			Value:            &value32,
			Signature:        &signatureString,
		}
		err = walletServer.blockchainClient.AddTransaction(blockchainTransactionRequest)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
			walletServer.write(writer, "fail")
			return
		}
		walletServer.write(writer, "success")
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("invalid HTTP method")
	}
}

func (walletServer *WalletServer) GetTransactions(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := walletServer.blockchainClient.GetTransactions()
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			walletServer.write(writer, "fail")
			return
		}
		var marshaledTransactions []byte
		marshaledTransactions, err = json.Marshal(transactions)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to marshal transactions: %w", err).Error())
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		walletServer.write(writer, string(marshaledTransactions[:]))
	default:
		walletServer.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (walletServer *WalletServer) Mine(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := walletServer.blockchainClient.Mine()
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to mine: %w", err).Error())
			walletServer.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("invalid HTTP method")
	}
}

func (walletServer *WalletServer) StartMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := walletServer.blockchainClient.StartMining()
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to start mining: %w", err).Error())
			walletServer.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("invalid HTTP method")
	}
}

func (walletServer *WalletServer) StopMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := walletServer.blockchainClient.StopMining()
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to stop mining: %w", err).Error())
			walletServer.write(writer, "fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("invalid HTTP method")
	}
}

func (walletServer *WalletServer) WalletAmount(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		amountRequest := chain.AmountRequest{
			Address: &address,
		}
		if amountRequest.IsInvalid() {
			walletServer.logger.Error("field(s) are missing in amount request to wallet server")
			walletServer.write(writer, "fail")
			return
		}
		amountResponse, err := walletServer.blockchainClient.GetAmount(amountRequest)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to get amountResponse: %w", err).Error())
			walletServer.write(writer, "fail")
			return
		}
		var marshaledAmount []byte
		marshaledAmount, err = json.Marshal(amountResponse)
		if err != nil {
			walletServer.logger.Error(fmt.Errorf("failed to marshal amountResponse: %w", err).Error())
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		walletServer.write(writer, string(marshaledAmount[:]))
	default:
		walletServer.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (walletServer *WalletServer) Run() {
	http.HandleFunc("/", walletServer.Index)
	http.HandleFunc("/wallet", walletServer.Wallet)
	http.HandleFunc("/transaction", walletServer.CreateTransaction)
	http.HandleFunc("/transactions", walletServer.GetTransactions)
	http.HandleFunc("/wallet/amount", walletServer.WalletAmount)
	http.HandleFunc("/mine", walletServer.Mine)
	http.HandleFunc("/mine/start", walletServer.StartMining)
	http.HandleFunc("/mine/stop", walletServer.StopMining)
	walletServer.logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(walletServer.Port())), nil).Error())
}

func (walletServer *WalletServer) write(writer http.ResponseWriter, message string) {
	i, err := io.WriteString(writer, message)
	if err != nil || i == 0 {
		walletServer.logger.Error(fmt.Errorf("failed to write message: %s", message).Error())
	}
}
