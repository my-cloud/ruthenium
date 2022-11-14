package server

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/log"
	"net/http"
)

type WalletHandler struct {
	mnemonic       string
	derivationPath string
	password       string
	privateKey     string
	logger         *log.Logger
}

func NewWalletHandler(mnemonic string, derivationPath string, password string, privateKey string, logger *log.Logger) *WalletHandler {
	return &WalletHandler{mnemonic, derivationPath, password, privateKey, logger}
}

func (handler *WalletHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
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
		NewIoWriter(writer, handler.logger).Write(string(marshaledWallet[:]))
	default:
		handler.logger.Error("invalid HTTP method")
		writer.WriteHeader(http.StatusBadRequest)
	}
}
