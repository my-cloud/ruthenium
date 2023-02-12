package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_PrivateKeyFromMnemonic(t *testing.T) {
	// Arrange
	mnemonic := encryption.NewMnemonic(test.Mnemonic)

	// Act
	privateKey, _ := mnemonic.PrivateKey(test.DerivationPath, "")

	// Assert
	expectedPrivateKey := test.PrivateKey
	actualPrivateKey := privateKey.String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}
