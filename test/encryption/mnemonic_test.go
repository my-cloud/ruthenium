package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_PrivateKeyFromMnemonic(t *testing.T) {
	// Arrange
	mnemonic := encryption.NewMnemonic(test.Mnemonic1)

	// Act
	privateKey, _ := mnemonic.PrivateKey(test.DerivationPath, "")

	// Assert
	expectedPrivateKey := "0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd"
	actualPrivateKey := privateKey.String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}
