package encryption

import (
	"fmt"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_NewPrivateKeyFromHex(t *testing.T) {
	// Arrange
	// Act
	privateKey, _ := NewPrivateKeyFromHex(test.PrivateKey)

	// Assert
	expectedPrivateKey := test.PrivateKey
	actualPrivateKey := privateKey.String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}

func Test_NewPrivateKeyFromMnemonic(t *testing.T) {
	// Arrange
	// Act
	privateKey, _ := NewPrivateKeyFromMnemonic(test.Mnemonic, test.DerivationPath, "")

	// Assert
	expectedPrivateKey := test.PrivateKey
	actualPrivateKey := privateKey.String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}
