package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_AddressFromPublicKey(t *testing.T) {
	// Arrange
	publicKey, _ := encryption.DecodePublicKey(test.PublicKey)

	// Act
	address := publicKey.Address()

	// Assert
	expectedAddress := test.Address
	test.Assert(t, address == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, address))
}
