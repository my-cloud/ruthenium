package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"testing"
)

func Test_Clear_ListsAreNotEmpty_ListsAreCleared(t *testing.T) {
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
		removedAddresses:    []string{"test"},
	}
	registry.Clear()
	test.Assert(t, !registry.registeredAddresses["test"], "registeredAddresses has not been correctly cleared")
	test.Assert(t, registry.removedAddresses == nil, "removedAddresses has not been correctly cleared")
}

func Test_Filter_OneOfTwoAddressesIsAlreadyRegistered_ReturnsOneAddress(t *testing.T) {
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
	}
	expectedAddress := "new"
	newAddresses := registry.Filter([]string{"test", expectedAddress})
	actualLength := len(newAddresses)
	expectedLength := 1
	test.Assert(t, actualLength == expectedLength, fmt.Sprintf("newAddresses should contain %d item whereas it contains %d", expectedLength, actualLength))
	actualAddress := newAddresses[0]
	test.Assert(t, actualAddress == expectedAddress, fmt.Sprintf("newAddresses should contain %s whereas it contains %s", expectedAddress, actualAddress))
}

func Test_IsRegistered_Registered_ReturnsTrue(t *testing.T) {
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
	}
	isRegistered := registry.IsRegistered("test")
	test.Assert(t, isRegistered, "address is not registered whereas it should be")
}

func Test_IsRegistered_NotRegistered_ReturnsFalse(t *testing.T) {
	registry := NewAddressesRegistry(nil, nil)
	isRegistered := !registry.IsRegistered("new")
	test.Assert(t, isRegistered, "address is registered whereas it should not")
}

func Test_RemovedAddresses_OneAddressRemoved_ReturnsOneAddress(t *testing.T) {
	registry := &AddressesRegistry{
		removedAddresses: []string{"test"},
	}
	removedAddresses := registry.RemovedAddresses()
	test.Assert(t, removedAddresses[0] == "test", "address is not removed whereas it should be")
}

func Test_Update_AddingAndRemovingAddresses_AddressesCorrectlyUpdated(t *testing.T) {
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
		removedAddresses:    []string{"test"},
	}
	registry.Update([]string{"new"}, []string{"test"})
	test.Assert(t, registry.registeredAddresses["new"], "registeredAddresses has not been correctly updated")
	test.Assert(t, !registry.registeredAddresses["test"], "registeredAddresses has not been correctly updated")
	test.Assert(t, len(registry.removedAddresses) == 0, "removedAddresses has not been correctly updated")
}
