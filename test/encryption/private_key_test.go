package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_NewPrivateKeyFromHex(t *testing.T) {
	// Arrange
	// Act
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)

	// Assert
	expectedPrivateKey := test.PrivateKey
	actualPrivateKey := privateKey.String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}

func Test_NewPrivateKeyFromMnemonic(t *testing.T) {
	// Arrange
	// Act
	privateKey, _ := encryption.NewPrivateKeyFromMnemonic(test.Mnemonic, test.DerivationPath, "")

	// Assert
	expectedPrivateKey := test.PrivateKey
	actualPrivateKey := privateKey.String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}
