package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"testing"
)

func Test_Clear_ListsAreNotEmpty_ListsAreCleared(t *testing.T) {
	// Arrange
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
		removedAddresses:    []string{"test"},
	}

	// Act
	registry.Clear()

	// Assert
	test.Assert(t, !registry.registeredAddresses["test"], "registeredAddresses has not been correctly cleared")
	test.Assert(t, registry.removedAddresses == nil, "removedAddresses has not been correctly cleared")
}

func Test_Filter_OneOfTwoAddressesIsAlreadyRegistered_ReturnsOneAddress(t *testing.T) {
	// Arrange
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
	}
	expectedAddress := "new"

	// Act
	newAddresses := registry.Filter([]string{"test", expectedAddress})

	// Assert
	actualLength := len(newAddresses)
	expectedLength := 1
	test.Assert(t, actualLength == expectedLength, fmt.Sprintf("newAddresses should contain %d item whereas it contains %d", expectedLength, actualLength))
	actualAddress := newAddresses[0]
	test.Assert(t, actualAddress == expectedAddress, fmt.Sprintf("newAddresses should contain %s whereas it contains %s", expectedAddress, actualAddress))
}

func Test_IsRegistered_Registered_ReturnsTrue(t *testing.T) {
	// Arrange
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
	}

	// Act
	isRegistered := registry.IsRegistered("test")

	// Assert
	test.Assert(t, isRegistered, "address is not registered whereas it should be")
}

func Test_IsRegistered_NotRegistered_ReturnsFalse(t *testing.T) {
	// Arrange
	registry := NewAddressesRegistry(nil, nil)

	// Act
	isRegistered := !registry.IsRegistered("new")

	// Assert
	test.Assert(t, isRegistered, "address is registered whereas it should not")
}

func Test_RemovedAddresses_OneAddressRemoved_ReturnsOneAddress(t *testing.T) {
	// Arrange
	registry := &AddressesRegistry{
		removedAddresses: []string{"test"},
	}

	// Act
	removedAddresses := registry.RemovedAddresses()

	// Assert
	test.Assert(t, removedAddresses[0] == "test", "address is not removed whereas it should be")
}

func Test_Update_AddingAndRemovingAddresses_AddressesCorrectlyUpdated(t *testing.T) {
	// Arrange
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
		removedAddresses:    []string{"test"},
	}

	// Act
	registry.Update([]string{"new"}, []string{"test"})

	// Assert
	test.Assert(t, registry.registeredAddresses["new"], "registeredAddresses has not been correctly updated")
	test.Assert(t, !registry.registeredAddresses["test"], "registeredAddresses has not been correctly updated")
	test.Assert(t, len(registry.removedAddresses) == 0, "removedAddresses has not been correctly updated")
}
