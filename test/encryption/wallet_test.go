package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_DecodeWallet_PrivateKeyProvided_ReturnsWalletForPrivateKey(t *testing.T) {
	// Arrange
	expectedPrivateKey := "0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd"

	// Act
	wallet, _ := encryption.DecodeWallet("", "", "", expectedPrivateKey)

	// Assert
	actualPrivateKey := wallet.PrivateKey().String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}

func Test_DecodeWallet_BothPrivateKeyAndMnemonicAreEmpty_ReturnsEmptyWallet(t *testing.T) {
	// Act
	wallet, _ := encryption.DecodeWallet("", "", "", "")

	// Assert
	test.Assert(t, wallet.PrivateKey() == nil, "Private key is not nil whereas it should be.")
}

func Test_MarshalJSON_ValidPrivateKey_ReturnsMarshaledJsonWithoutError(t *testing.T) {
	// Arrange
	privateKey := "0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd"
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
