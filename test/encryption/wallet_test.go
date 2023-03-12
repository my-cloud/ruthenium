package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_NewWallet_BothMnemonicAndDerivationPathProvided_ReturnsCorrespondingWallet(t *testing.T) {
	// Act
	publicKey, _ := encryption.NewPublicKeyFromHex(test.PublicKey)
	wallet := encryption.NewWallet(publicKey)

	// Assert
	actualAddress := wallet.Address()
	expectedAddress := test.Address
	test.Assert(t, actualAddress == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, actualAddress))
}

func Test_MarshalJSON_ValidPrivateKey_ReturnsMarshaledJsonWithoutError(t *testing.T) {
	// Arrange
	publicKey, _ := encryption.NewPublicKeyFromHex(test.PublicKey)
	wallet := encryption.NewWallet(publicKey)

	// Act
	marshaledWallet, err := wallet.MarshalJSON()

	// Assert
	test.Assert(t, marshaledWallet != nil, "Marshaled wallet is nil.")
	test.Assert(t, err == nil, "Marshal wallet returned an error.")
}
