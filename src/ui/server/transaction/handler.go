package transaction

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	host               network.Neighbor
	particlesInOneAtom uint64
	logger             *log.Logger
}

func NewHandler(host network.Neighbor, particlesInOneAtom uint64, logger *log.Logger) *Handler {
	return &Handler{host, particlesInOneAtom, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		var transactionRequest server.TransactionRequest
		jsonWriter := server.NewIoWriter(writer, handler.logger)
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&transactionRequest)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to decode transaction request: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid transaction request")
			return
		}
		if transactionRequest.IsInvalid() {
			errorMessage := "field(s) are missing in transaction request"
			handler.logger.Error(errorMessage)
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write(errorMessage)
			return
		}
		privateKey, err := encryption.DecodePrivateKey(*transactionRequest.SenderPrivateKey)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to decode transaction private key: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid private key")
			return
		}
		value, err := atomsToParticles(*transactionRequest.Value, handler.particlesInOneAtom)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to parse transaction value: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid transaction value")
			return
		}
		senderPublicKey := encryption.NewPublicKey(privateKey)
		transaction := server.NewTransaction(*transactionRequest.RecipientAddress, *transactionRequest.SenderAddress, senderPublicKey, time.Now().UnixNano(), value)
		err = transaction.Sign(privateKey)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to generate signature: %w", err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write("invalid signature")
			return
		}
		blockchainTransactionRequest := transaction.GetRequest()
		err = handler.host.AddTransaction(blockchainTransactionRequest)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to create transaction: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		jsonWriter.Write("success")
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}

func atomsToParticles(atoms string, particlesInOneAtom uint64) (particles uint64, err error) {
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
		trailingZerosCount := len(strconv.Itoa(int(particlesInOneAtom))) - 1 - len(decimalsString)
		trailedDecimalsString := fmt.Sprintf("%s%s", decimalsString, strings.Repeat("0", trailingZerosCount))
		var decimals uint64
		decimals, err = parseUint64(trailedDecimalsString)
		if err != nil {
			return
		}
		particles = units*particlesInOneAtom + decimals
	} else {
		var units uint64
		units, err = parseUint64(atoms)
		if err != nil {
			return
		}
		particles = units * particlesInOneAtom
	}
	return
}

func parseUint64(valueString string) (value uint64, err error) {
	return strconv.ParseUint(valueString, 10, 64)
}
