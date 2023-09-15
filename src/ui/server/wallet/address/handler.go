package address

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
)

type Handler struct {
	logger log.Logger
}

func NewHandler(logger log.Logger) *Handler {
	return &Handler{logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		jsonWriter := server.NewIoWriter(writer, handler.logger)
		publicKeyString := req.URL.Query().Get("publicKey")
		publicKey, err := encryption.NewPublicKeyFromHex(publicKeyString)
		if err != nil {
			errorMessage := "failed to decode public key"
			handler.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
			writer.WriteHeader(http.StatusBadRequest)
			jsonWriter.Write(errorMessage)
			return
		}
		address := publicKey.Address()
		marshaledAddress, err := json.Marshal(address)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal address: %w", err).Error())
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		server.NewIoWriter(writer, handler.logger).Write(string(marshaledAddress[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
