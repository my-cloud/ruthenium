package wallet

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
)

type Handler struct {
	mnemonic       string
	derivationPath string
	password       string
	privateKey     string
	logger         *log.Logger
}

func NewHandler(mnemonic string, derivationPath string, password string, privateKey string, logger *log.Logger) *Handler {
	return &Handler{mnemonic, derivationPath, password, privateKey, logger}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		wallet, err := encryption.DecodeWallet(handler.mnemonic, handler.derivationPath, handler.password, handler.privateKey)
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to create wallet: %w", err).Error())
			return
		}
		marshaledWallet, err := wallet.MarshalJSON()
		if err != nil {
			handler.logger.Error(fmt.Errorf("failed to marshal wallet: %w", err).Error())
			return
		}
		writer.Header().Add("Content-Type", "application/json")
		server.NewIoWriter(writer, handler.logger).Write(string(marshaledWallet[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
