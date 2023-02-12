package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_DecodeWallet_BothMnemonicAndDerivationPathProvided_ReturnsCorrespondingWallet(t *testing.T) {
	// Act
	wallet, _ := encryption.DecodeWallet(test.Mnemonic, test.DerivationPath, "", "")

	// Assert
	actualAddress := wallet.Address()
	expectedAddress := test.Address
	test.Assert(t, actualAddress == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, actualAddress))
}

func Test_DecodeWallet_PrivateKeyProvided_ReturnsCorrespondingWallet(t *testing.T) {
	// Act
	wallet, _ := encryption.DecodeWallet("", "", "", test.PrivateKey)

	// Assert
	actualAddress := wallet.Address()
	expectedAddress := test.Address
	test.Assert(t, actualAddress == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, actualAddress))
}

func Test_DecodeWallet_BothPrivateKeyAndMnemonicAreEmpty_ReturnsEmptyWallet(t *testing.T) {
	// Act
	wallet, _ := encryption.DecodeWallet("", "", "", "")

	// Assert
	test.Assert(t, wallet.Address() == "", "Address is not empty whereas it should be.")
}

func Test_MarshalJSON_ValidPrivateKey_ReturnsMarshaledJsonWithoutError(t *testing.T) {
	// Arrange
	privateKey := test.PrivateKey
	wallet, _ := encryption.DecodeWallet("", "", "", privateKey)

	// Act
	marshaledWallet, err := wallet.MarshalJSON()

	// Assert
	test.Assert(t, marshaledWallet != nil, "Marshaled wallet is nil.")
	test.Assert(t, err == nil, "Marshal wallet returned an error.")
}

func Test_MarshalJSON_EmptyWallet_ReturnsMarshaledJsonWithoutError(t *testing.T) {
	// Arrange
	wallet := encryption.NewEmptyWallet()

	// Act
	marshaledWallet, err := wallet.MarshalJSON()

	// Assert
	test.Assert(t, marshaledWallet != nil, "Marshaled wallet is nil.")
	test.Assert(t, err == nil, "Marshal wallet returned an error.")
}
