package server

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/p2p"
	"html/template"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

const DefaultPort = 8080

type Controller struct {
	mnemonic         string
	derivationPath   string
	password         string
	privateKey       string
	port             uint16
	blockchainClient *network.Neighbor
	templatesPath    string
	logger           *log.Logger
}

func NewController(mnemonic string, derivationPath string, password string, privateKey string, port uint16, hostIp string, hostPort uint16, templatesPath string, level log.Level) *Controller {
	logger := log.NewLogger(level)
	target := network.NewTarget(hostIp, hostPort)
	peering := p2p.NewSenderFactory()
	blockchainClient, err := network.NewNeighbor(target, peering, logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to find blockchain client: %w", err).Error())
	}
	return &Controller{mnemonic, derivationPath, password, privateKey, port, blockchainClient, templatesPath, logger}
}

func (controller *Controller) Port() uint16 {
	return controller.port
}

func (controller *Controller) BlockchainClient() *network.Neighbor {
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
		wallet, err := encryption.DecodeWallet(controller.mnemonic, controller.derivationPath, controller.password, controller.privateKey)
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
		controller.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
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
			writer.WriteHeader(http.StatusBadRequest)
			controller.write(writer, "invalid transaction request")
			return
		}
		if transactionRequest.IsInvalid() {
			errorMessage := "field(s) are missing in transaction request"
			controller.logger.Error(errorMessage)
			writer.WriteHeader(http.StatusBadRequest)
			controller.write(writer, errorMessage)
			return
		}
		privateKey, err := encryption.DecodePrivateKey(*transactionRequest.SenderPrivateKey)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to decode transaction private key: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			controller.write(writer, "invalid private key")
			return
		}
		value, err := controller.atomsToParticles(*transactionRequest.Value)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to parse transaction value: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			controller.write(writer, "invalid transaction value")
			return
		}
		senderPublicKey := encryption.NewPublicKey(privateKey)
		transaction := NewTransaction(*transactionRequest.RecipientAddress, *transactionRequest.SenderAddress, senderPublicKey, time.Now().UnixNano(), value)
		err = transaction.Sign(privateKey)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to generate signature: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			controller.write(writer, "invalid signature")
			return
		}
		blockchainTransactionRequest := transaction.GetRequest()
		err = controller.blockchainClient.AddTransaction(blockchainTransactionRequest)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		controller.write(writer, "success")
	default:
		controller.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (controller *Controller) GetTransactions(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		transactions, err := controller.blockchainClient.GetTransactions()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to get transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		marshaledTransactions, err := json.Marshal(transactions)
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to marshal transactions: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
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
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		controller.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (controller *Controller) StartMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := controller.blockchainClient.StartMining()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to start mining: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		controller.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (controller *Controller) StopMining(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		err := controller.blockchainClient.StopMining()
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to stop mining: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
		}
	default:
		controller.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func (controller *Controller) WalletAmount(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		amountRequest := AmountRequest{
			Address: &address,
		}
		if amountRequest.IsInvalid() {
			controller.logger.Error("field(s) are missing in amount request")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		amount := NewAmount(*amountRequest.Address)
		amountResponse, err := controller.blockchainClient.GetAmount(*amount.GetRequest())
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to get amountResponse: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var marshaledAmount []byte
		marshaledAmount, err = json.Marshal(&AmountResponse{
			Amount: float64(amountResponse.Amount) / network.ParticlesCount,
		})
		if err != nil {
			controller.logger.Error(fmt.Errorf("failed to marshal amountResponse: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
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
	controller.logger.Info("user interface server is running...")
	controller.logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(controller.Port())), nil).Error())
}

func (controller *Controller) write(writer http.ResponseWriter, message string) {
	i, err := io.WriteString(writer, message)
	if err != nil || i == 0 {
		controller.logger.Error(fmt.Sprintf("failed to write message: %s", message))
	}
}

func (controller *Controller) atomsToParticles(atoms string) (particles uint64, err error) {
	const decimalSeparator = "."
	i := strings.Index(atoms, decimalSeparator)
	if i > 12 || (i == -1 && len(atoms) > 12) {
		err = fmt.Errorf("transaction value is too big")
		return
	}
	if i > -1 {
		unitsString := atoms[:i]
		var units uint64
		units, err = parseUint64(unitsString)
		if err != nil {
			return
		}
		decimalsString := atoms[i+1:]
		trailingZerosCount := len(strconv.Itoa(network.ParticlesCount)) - 1 - len(decimalsString)
		trailedDecimalsString := fmt.Sprintf("%s%s", decimalsString, strings.Repeat("0", trailingZerosCount))
		var decimals uint64
		decimals, err = parseUint64(trailedDecimalsString)
		if err != nil {
			return
		}
		particles = units*network.ParticlesCount + decimals
	} else {
		var units uint64
		units, err = parseUint64(atoms)
		if err != nil {
			return
		}
		particles = units * network.ParticlesCount
	}
	return
}

func parseUint64(valueString string) (value uint64, err error) {
	return strconv.ParseUint(valueString, 10, 64)
}
