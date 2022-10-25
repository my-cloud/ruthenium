package encryption

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_PrivateKeyFromMnemonic(t *testing.T) {
	// Arrange
	mnemonic := encryption.NewMnemonic("artist silver basket insane canvas top drill social reflect park fruit bless")

	// Act
	privateKey, _ := mnemonic.PrivateKey("m/44'/60'/0'/0/0", "")

	// Assert
	expectedPrivateKey := "0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd"
	actualPrivateKey := privateKey.String()
	test.Assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}

func Test_PublicKeyFromPrivateKey(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.DecodePrivateKey("0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd")

	// Act
	publicKey := encryption.NewPublicKey(privateKey)

	// Assert
	expectedPublicKey := "0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782"
	actualPublicKey := publicKey.String()
	test.Assert(t, actualPublicKey == expectedPublicKey, fmt.Sprintf("Wrong public key. Expected: %s - Actual: %s", expectedPublicKey, actualPublicKey))
}

func Test_AddressFromPublicKey(t *testing.T) {
	// Arrange
	publicKey, _ := encryption.DecodePublicKey("0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782")

	// Act
	address := publicKey.Address()

	// Assert
	expectedAddress := "0x9C69443c3Ec0D660e257934ffc1754EB9aD039CB"
	test.Assert(t, address == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, address))
}
