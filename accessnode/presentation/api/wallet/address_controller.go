package wallet

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/my-cloud/ruthenium/accessnode/infrastructure/io"
	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type AddressController struct {
	logger log.Logger
}

func NewAddressController(logger log.Logger) *AddressController {
	return &AddressController{logger}
}

func (controller *AddressController) GetWalletAddress(writer http.ResponseWriter, req *http.Request) {
	jsonWriter := io.NewIoWriter(writer, controller.logger)
	publicKeyString := req.URL.Query().Get("publicKey")
	publicKey, err := encryption.NewPublicKeyFromHex(publicKeyString)
	if err != nil {
		errorMessage := "failed to decode public key"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		writer.WriteHeader(http.StatusBadRequest)
		jsonWriter.Write(errorMessage)
		return
	}
	address := publicKey.Address()
	marshaledAddress, err := json.Marshal(address)
	if err != nil {
		controller.logger.Error(fmt.Errorf("failed to marshal address: %w", err).Error())
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	io.NewIoWriter(writer, controller.logger).Write(string(marshaledAddress[:]))
}
