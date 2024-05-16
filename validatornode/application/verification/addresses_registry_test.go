package verification

import (
	"errors"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"testing"
)

// TODO fix method names
func TestNewAddressesRegistry(t *testing.T) {
	humansManager := new(HumansManagerMock)
	logger := log.NewLoggerMock()
	registry := NewAddressesRegistry(humansManager, logger)
	test.Assert(t, registry != nil, "registry is nil whereas it should not")
}

func TestAddressesRegistry_Clear(t *testing.T) {
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
		temporaryAddresses:  map[string]bool{"test": true},
		removedAddresses:    []string{"test"},
	}
	registry.Clear()
	test.Assert(t, !registry.registeredAddresses["test"], "registeredAddresses has not been correctly cleared")
	test.Assert(t, !registry.temporaryAddresses["test"], "temporaryAddresses has not been correctly cleared")
	test.Assert(t, registry.removedAddresses == nil, "removedAddresses has not been correctly cleared")
}

// TODO
//func TestAddressesRegistry_Copy(t *testing.T) {
//	logger := log.NewLoggerMock()
//	registry := &AddressesRegistry{
//		registeredAddresses: map[string]bool{"test": true},
//		temporaryAddresses:  map[string]bool{"test": true},
//		removedAddresses:    []string{"test"},
//		logger:              logger,
//	}
//	_ = registry.Copy()
//	assert.Equal(t, registry.humansManager, registryCopy.humansManager)
//	assert.NotEqual(t, registry.registeredAddresses, registryCopy.registeredAddresses)
//	assert.NotEqual(t, registry.temporaryAddresses, registryCopy.temporaryAddresses)
//	assert.Equal(t, registry.removedAddresses, registryCopy.removedAddresses)
//	assert.Equal(t, registry.logger, registryCopy.logger)
//}

func TestAddressesRegistry_Filter(t *testing.T) {
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

func TestAddressesRegistry_IsRegistered(t *testing.T) {
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
	}
	isRegistered := registry.IsRegistered("test")
	test.Assert(t, isRegistered, "address is not registered whereas it should be")
	isRegistered = !registry.IsRegistered("new")
	test.Assert(t, isRegistered, "address is registered whereas it should not")
}

func TestAddressesRegistry_RemovedAddresses(t *testing.T) {
	registry := &AddressesRegistry{
		removedAddresses: []string{"test"},
	}
	removedAddresses := registry.RemovedAddresses()
	test.Assert(t, removedAddresses[0] == "test", "address is not removed whereas it should be")
}

func TestAddressesRegistry_Update(t *testing.T) {
	registry := &AddressesRegistry{
		registeredAddresses: map[string]bool{"test": true},
		removedAddresses:    []string{"test"},
	}
	registry.Update([]string{"new"}, []string{"test"})
	test.Assert(t, registry.registeredAddresses["new"], "registeredAddresses has not been correctly updated")
	test.Assert(t, !registry.registeredAddresses["test"], "registeredAddresses has not been correctly updated")
	test.Assert(t, len(registry.removedAddresses) == 0, "removedAddresses has not been correctly updated")
}

func TestAddressesRegistry_Verify(t *testing.T) {
	humansManager := new(HumansManagerMock)
	humansManager.IsRegisteredFunc = func(address string) (bool, error) {
		if address == "registered" {
			return true, nil
		} else if address == "notregistered" {
			return false, nil
		}
		return false, errors.New("error fetching registration status")
	}
	logger := log.NewLoggerMock()
	registry := &AddressesRegistry{
		humansManager:       humansManager,
		registeredAddresses: map[string]bool{"registered": true},
		temporaryAddresses:  map[string]bool{"registered": true, "notregistered": true},
		logger:              logger,
	}

	// Test adding a registered address
	err := registry.Verify([]string{"registered"}, nil)
	test.Assert(t, err == nil, "error is not nil whereas it should be")

	// Test removing a registered address
	err = registry.Verify(nil, []string{"registered"})
	expectedMessage := "a removed address is registered"
	actualMessage := err.Error()
	test.Assert(t, actualMessage == expectedMessage, fmt.Sprintf("expected message was not logged: %s\nlogged message: %s", expectedMessage, actualMessage))

	// Test adding a not registered address
	err = registry.Verify([]string{"notregistered"}, nil)
	expectedMessage = "an added address is not registered"
	actualMessage = err.Error()
	test.Assert(t, actualMessage == expectedMessage, fmt.Sprintf("expected message was not logged: %s\nlogged message: %s", expectedMessage, actualMessage))
}
