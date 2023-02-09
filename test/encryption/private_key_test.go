package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_PublicKeyFromPrivateKey(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.DecodePrivateKey(test.PrivateKey)

	// Act
	publicKey := encryption.NewPublicKey(privateKey)

	// Assert
	expectedPublicKey := test.PublicKey
	actualPublicKey := publicKey.String()
	test.Assert(t, actualPublicKey == expectedPublicKey, fmt.Sprintf("Wrong public key. Expected: %s - Actual: %s", expectedPublicKey, actualPublicKey))
}
