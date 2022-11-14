package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_AddressFromPublicKey(t *testing.T) {
	// Arrange
	publicKey, _ := encryption.DecodePublicKey("0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782")

	// Act
	address := publicKey.Address()

	// Assert
	expectedAddress := "0x9C69443c3Ec0D660e257934ffc1754EB9aD039CB"
	test.Assert(t, address == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, address))
}
