package server

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/network"
	"net/http"
)

type WalletAmountHandler struct {
	host   neighborhood.Neighbor
	logger *log.Logger
}

func NewWalletAmountHandler(host neighborhood.Neighbor, logger *log.Logger) *WalletAmountHandler {
	return &WalletAmountHandler{host, logger}
}

func (handler *WalletAmountHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		amountRequest := neighborhood.AmountRequest{
			Address: &address,
		}
		if amountRequest.IsInvalid() {
			handler.logger.Error("field(s) are missing in amount request")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		amount := NewAmount(*amountRequest.Address)
		amountResponse, err := handler.host.GetAmount(*amount.GetRequest())
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get amountResponse: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var marshaledAmount []byte
		marshaledAmount, err = json.Marshal(&AmountResponse{
			Amount: float64(amountResponse.Amount) / network.ParticlesCount,
		})
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal amountResponse: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		NewIoWriter(writer, handler.logger).Write(string(marshaledAmount[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
