package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_NewWallet_BothMnemonicAndDerivationPathProvided_ReturnsCorrespondingWallet(t *testing.T) {
	// Act
	wallet, _ := encryption.NewWallet(test.Mnemonic, test.DerivationPath, "", "")

	// Assert
	actualAddress := wallet.Address()
	expectedAddress := test.Address
	test.Assert(t, actualAddress == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, actualAddress))
}

func Test_NewWallet_PrivateKeyProvided_ReturnsCorrespondingWallet(t *testing.T) {
	// Act
	wallet, _ := encryption.NewWallet("", "", "", test.PrivateKey)

	// Assert
	actualAddress := wallet.Address()
	expectedAddress := test.Address
	test.Assert(t, actualAddress == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, actualAddress))
}

func Test_NewWallet_BothPrivateKeyAndMnemonicAreEmpty_ReturnsError(t *testing.T) {
	// Act
	wallet, err := encryption.NewWallet("", "", "", "")

	// Assert
	test.Assert(t, wallet == nil, "Wallet is not nil whereas it should be.")
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
}

func Test_MarshalJSON_ValidPrivateKey_ReturnsMarshaledJsonWithoutError(t *testing.T) {
	// Arrange
	privateKey := test.PrivateKey
	wallet, _ := encryption.NewWallet("", "", "", privateKey)

	// Act
	marshaledWallet, err := wallet.MarshalJSON()

	// Assert
	test.Assert(t, marshaledWallet != nil, "Marshaled wallet is nil.")
	test.Assert(t, err == nil, "Marshal wallet returned an error.")
}
