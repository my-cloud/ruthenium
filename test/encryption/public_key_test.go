package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_NewPublicKey(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)

	// Act
	publicKey := encryption.NewPublicKey(privateKey)

	// Assert
	expectedPublicKey := test.PublicKey
	actualPublicKey := publicKey.String()
	test.Assert(t, actualPublicKey == expectedPublicKey, fmt.Sprintf("Wrong public key. Expected: %s - Actual: %s", expectedPublicKey, actualPublicKey))
}

func Test_Address(t *testing.T) {
	// Arrange
	publicKey, _ := encryption.NewPublicKeyFromHex(test.PublicKey)

	// Act
	address := publicKey.Address()

	// Assert
	expectedAddress := test.Address
	test.Assert(t, address == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, address))
}
