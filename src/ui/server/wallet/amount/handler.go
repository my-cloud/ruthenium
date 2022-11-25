package amount

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet"
	"net/http"
)

type Handler struct {
	host           network.Neighbor
	particlesCount uint64
	logger         *log.Logger
}

func NewHandler(host network.Neighbor, particlesCount uint64, logger *log.Logger) *Handler {
	return &Handler{host, particlesCount, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		address := req.URL.Query().Get("address")
		amountRequest := network.AmountRequest{
			Address: &address,
		}
		if amountRequest.IsInvalid() {
			handler.logger.Error("field(s) are missing in amount request")
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		amount := wallet.NewAmount(*amountRequest.Address)
		amountResponse, err := handler.host.GetAmount(*amount.GetRequest())
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get amountResponse: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		var marshaledAmount []byte
		marshaledAmount, err = json.Marshal(&wallet.AmountResponse{
			Amount: float64(amountResponse.Amount) / float64(handler.particlesCount),
		})
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal amountResponse: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		server.NewIoWriter(writer, handler.logger).Write(string(marshaledAmount[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
