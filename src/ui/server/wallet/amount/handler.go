package amount

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
)

type Handler struct {
	host           network.Neighbor
	particlesCount uint64
	logger         log.Logger
}

func NewHandler(host network.Neighbor, particlesCount uint64, logger log.Logger) *Handler {
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
		amount, err := handler.host.GetAmount(*amountRequest.Address)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to get amount: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		marshaledAmount, err := json.Marshal(float64(amount) / float64(handler.particlesCount))
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal amount: %w", err).Error())
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
