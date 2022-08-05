package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path"
	"ruthenium/src/chain"
	"ruthenium/src/log"
	"ruthenium/src/rest"
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
	if !blockchainClient.IsFound() {
		logger.Error("Unable to find blockchain node")
	}
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
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to parse the template\n%v", err))
		} else if err = t.Execute(w, ""); err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to execute the template\n%v", err))
		}
	default:
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) Wallet(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		writer.Header().Add("Content-Type", "application/json")
		var wallet *chain.Wallet
		if walletServer.publicKey == "" || walletServer.privateKey == "" {
			privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
			if err != nil {
				walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to generate private key\n%v", err))
			} else {
				wallet = chain.NewWallet(privateKey)
			}
		} else {
			publicKey := chain.NewPublicKey(walletServer.publicKey)
			privateKey := chain.NewPrivateKey(walletServer.privateKey, publicKey)
			wallet = chain.PopWallet(publicKey, privateKey)
		}

		if wallet != nil {
			marshaledWallet, err := wallet.MarshalJSON()
			if err != nil {
				walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to marshal wallet\n%v", err))
			}
			i, err := io.WriteString(writer, string(marshaledWallet[:]))
			if err != nil || i == 0 {
				walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to write wallet\n%v", err))
			}
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
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
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to decode transaction request\n%v", err))
			jsonWriter.WriteStatus("fail")
		}
		if transactionRequest.IsInvalid() {
			walletServer.logger.Error("ERROR: Field(s) are missing in transaction request to wallet server")
			jsonWriter.WriteStatus("fail")
		}

		publicKey := chain.NewPublicKey(*transactionRequest.SenderPublicKey)
		privateKey := chain.NewPrivateKey(*transactionRequest.SenderPrivateKey, publicKey)
		value, err := strconv.ParseFloat(*transactionRequest.Value, 32)
		if err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to parse transaction value\n%v", err))
			jsonWriter.WriteStatus("fail")
		}
		value32 := float32(value)

		writer.Header().Add("Content-type", "application/json")

		transaction := chain.NewTransaction(*transactionRequest.SenderAddress, publicKey, *transactionRequest.RecipientAddress, value32)
		signature := chain.NewSignature(transaction, privateKey)

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
		err = walletServer.blockchainClient.UpdateTransactions(blockchainTransactionRequest)
		if err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to create transaction\n%v", err))
			jsonWriter.WriteStatus("fail")
		} else {
			jsonWriter.WriteStatus("success")
			return
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) GetTransactions(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := walletServer.blockchainClient.GetTransactions()

		jsonWriter := rest.NewJsonWriter(writer)
		writer.Header().Add("Content-Type", "application/json")
		if err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to get transactions\n%v", err))
			jsonWriter.WriteStatus("fail")
		} else {
			var marshaledTransactions []byte
			if transactions == nil {
				marshaledTransactions, err = json.Marshal([]chain.TransactionResponse{})
			} else {
				marshaledTransactions, err = json.Marshal(transactions)
			}
			if err != nil {
				walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to marshal transactions\n%v", err))
			}
			i, err := io.WriteString(writer, string(marshaledTransactions[:]))
			if err != nil || i == 0 {
				walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to write transactions\n%v", err))
			}
		}
	default:
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (walletServer *WalletServer) Mine(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := walletServer.blockchainClient.Mine()
		if err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to mine\n%v", err))
			jsonWriter := rest.NewJsonWriter(writer)
			jsonWriter.WriteStatus("fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) StartMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := walletServer.blockchainClient.StartMining()
		if err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to start mining\n%v", err))
			jsonWriter := rest.NewJsonWriter(writer)
			jsonWriter.WriteStatus("fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) StopMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := walletServer.blockchainClient.StopMining()
		if err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to stop mining\n%v", err))
			jsonWriter := rest.NewJsonWriter(writer)
			jsonWriter.WriteStatus("fail")
		}
	default:
		writer.WriteHeader(http.StatusBadRequest)
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
	}
}

func (walletServer *WalletServer) WalletAmount(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")

		amountRequest := chain.AmountRequest{
			Address: &address,
		}

		jsonWriter := rest.NewJsonWriter(writer)
		if amountRequest.IsInvalid() {
			walletServer.logger.Error("ERROR: Field(s) are missing in amount request to wallet server")
			jsonWriter.WriteStatus("fail")
			return
		}

		amount, err := walletServer.blockchainClient.GetAmount(amountRequest)

		writer.Header().Add("Content-Type", "application/json")
		if err != nil {
			walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to get amount\n%v", err))
			jsonWriter.WriteStatus("fail")
		} else {
			marshaledAmount, err := json.Marshal(struct {
				Message string  `json:"message"`
				Amount  float32 `json:"amount"`
			}{
				Message: "success",
				Amount:  amount.Amount,
			})
			if err != nil {
				walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to marshal amount\n%v", err))
			}
			i, err := io.WriteString(writer, string(marshaledAmount[:]))
			if err != nil || i == 0 {
				walletServer.logger.Error(fmt.Sprintf("ERROR: Failed to write amount\n%v", err))
			}
		}
	default:
		walletServer.logger.Error("ERROR: Invalid HTTP Method")
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
