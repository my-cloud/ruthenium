package wallet

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
)

type Handler struct {
	hostWallet *encryption.Wallet
	logger     log.Logger
}

func NewHandler(hostWallet *encryption.Wallet, logger log.Logger) *Handler {
	return &Handler{hostWallet, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		marshaledWallet, err := handler.hostWallet.MarshalJSON()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal host wallet: %w", err).Error())
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		server.NewIoWriter(writer, handler.logger).Write(string(marshaledWallet[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
