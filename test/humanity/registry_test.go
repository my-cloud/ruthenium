package humanity

import (
	"github.com/my-cloud/ruthenium/src/humanity"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_IsRegistered_NotRegistered_ReturnsFalse(t *testing.T) {
	// Arrange
	address := "0x0000000000000000000000000000000000000001"
	registry := humanity.NewRegistry()

	// Act
	isRegistered, _ := registry.IsRegistered(address)

	// Assert
	test.Assert(t, !isRegistered, "proof of humanity is valid whereas it should not")
}

func Test_IsRegistered_Registered_ReturnsTrue(t *testing.T) {
	// Arrange
	address := "0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a"
	registry := humanity.NewRegistry()

	// Act
	isRegistered, _ := registry.IsRegistered(address)

	// Assert
	test.Assert(t, isRegistered, "proof of humanity is invalid whereas it should be")
}
