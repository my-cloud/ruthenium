package verification

import (
	"github.com/my-cloud/ruthenium/validatornode/application"
	"sync"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type AddressesRegistry struct {
	humansManager       HumansManager
	registeredMutex     sync.RWMutex
	temporaryMutex      sync.RWMutex
	removedMutex        sync.RWMutex
	registeredAddresses map[string]bool
	removedAddresses    []string
	logger              log.Logger
}

func NewAddressesRegistry(humansManager HumansManager, logger log.Logger) *AddressesRegistry {
	registry := &AddressesRegistry{}
	registry.humansManager = humansManager
	registry.registeredAddresses = make(map[string]bool)
	registry.logger = logger
	return registry
}

func (registry *AddressesRegistry) Clear() {
	registry.registeredMutex.Lock()
	defer registry.registeredMutex.Unlock()
	registry.temporaryMutex.Lock()
	defer registry.temporaryMutex.Unlock()
	registry.removedMutex.Lock()
	defer registry.removedMutex.Unlock()
	registry.registeredAddresses = make(map[string]bool)
	registry.removedAddresses = nil
}

func (registry *AddressesRegistry) Copy() application.AddressesManager {
	registry.registeredMutex.RLock()
	defer registry.registeredMutex.RUnlock()
	registry.temporaryMutex.RLock()
	defer registry.temporaryMutex.RUnlock()
	registry.removedMutex.RLock()
	defer registry.removedMutex.RUnlock()
	registryCopy := &AddressesRegistry{}
	registryCopy.humansManager = registry.humansManager
	registryCopy.registeredAddresses = copyAddressesMap(registry.registeredAddresses)
	registryCopy.removedAddresses = registry.removedAddresses
	registryCopy.logger = registry.logger
	return registryCopy
}

func (registry *AddressesRegistry) Filter(addresses []string) []string {
	var newAddresses []string
	for _, address := range addresses {
		if !registry.registeredAddresses[address] {
			newAddresses = append(newAddresses, address)
		}
	}
	return newAddresses
}

func (registry *AddressesRegistry) IsRegistered(address string) bool {
	return registry.registeredAddresses[address]
}

func (registry *AddressesRegistry) RemovedAddresses() []string {
	return registry.removedAddresses
}

func (registry *AddressesRegistry) Synchronize(_ int64) {
	registry.registeredMutex.RLock()
	defer registry.registeredMutex.RUnlock()
	registry.removedMutex.Lock()
	defer registry.removedMutex.Unlock()
	for address := range registry.registeredAddresses {
		isPohValid, err := registry.humansManager.IsRegistered(address)
		if err != nil {
			registry.logger.Debug(err.Error())
		} else if !isPohValid {
			registry.removedAddresses = append(registry.removedAddresses, address)
		}
	}
	registry.temporaryMutex.Lock()
	defer registry.temporaryMutex.Unlock()
}

func (registry *AddressesRegistry) Update(addedAddresses []string, removedAddresses []string) {
	registry.registeredMutex.Lock()
	defer registry.registeredMutex.Unlock()
	registry.removedMutex.Lock()
	defer registry.removedMutex.Unlock()
	for _, address := range removedAddresses {
		registry.removedAddresses = removeAddress(registry.removedAddresses, address)
		delete(registry.registeredAddresses, address)
	}
	for _, address := range addedAddresses {
		registry.registeredAddresses[address] = true
	}
}

func copyAddressesMap(addresses map[string]bool) map[string]bool {
	addressesCopy := make(map[string]bool, len(addresses))
	for address := range addresses {
		addressesCopy[address] = true
	}
	return addressesCopy
}

func removeAddress(addresses []string, address string) []string {
	for i := 0; i < len(addresses); i++ {
		if address == addresses[i] {
			return append(addresses[:i], addresses[i+1:]...)
		}
	}
	return addresses
}
