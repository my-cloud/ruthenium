package wallet

import (
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
	response := io.NewResponse(writer, controller.logger)
	publicKeyString := req.URL.Query().Get("publicKey")
	publicKey, err := encryption.NewPublicKeyFromHex(publicKeyString)
	if err != nil {
		errorMessage := "failed to decode public key"
		controller.logger.Error(fmt.Errorf("%s: %w", errorMessage, err).Error())
		response.Write(http.StatusBadRequest, errorMessage)
		return
	}
	address := publicKey.Address()
	response.WriteJson(http.StatusOK, address)
}
