package poh

import (
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

// func Test_IsRegistered_InfuraKey_ReturnsTrue(t *testing.T) {
// 	// Arrange
// 	address := "0x0000000000000000000000000000000000000001"
// 	logger := log.NewLoggerMock()
// 	registry := NewRegistry("", logger)
//
// 	// Act
// 	isRegistered, err := registry.IsRegistered(address)
// 	fmt.Println(err)
//
// 	// Assert
// 	test.Assert(t, isRegistered, "proof of humanity is invalid whereas it should be")
// }

func Test_TakeAddedAddresses(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()
	registry := NewRegistry("", logger)

	// Act
	addedAddresses := registry.Filter([]string{"qesf"})
	addedAddresses2 := registry.Filter([]string{"qesf"})

	// Assert
	test.Assert(t, len(addedAddresses) == 1 && len(addedAddresses2) == 1, "oh!")
}
