package poh

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/my-cloud/ruthenium/validatornode/application/protocol/verification"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

const (
	pohSmartContractAddressHex = "0xC5E9dDebb09Cd64DfaCab4011A0D5cEDaf7c9BDb"
	clientUrl                  = "https://mainnet.infura.io/v3/"
)

type Registry struct {
	infuraKey           string
	registeredMutex     sync.RWMutex
	temporaryMutex      sync.RWMutex
	removedMutex        sync.RWMutex
	registeredAddresses map[string]bool
	temporaryAddresses  map[string]bool
	removedAddresses    []string
	logger              log.Logger
}

func NewRegistry(infuraKey string, logger log.Logger) *Registry {
	if infuraKey == "" {
		logger.Warn("infura key not provided")
	}
	registry := &Registry{}
	registry.infuraKey = infuraKey
	registry.registeredAddresses = make(map[string]bool)
	registry.temporaryAddresses = make(map[string]bool)
	registry.logger = logger
	return registry
}

func (registry *Registry) Clear() {
	registry.registeredMutex.Lock()
	defer registry.registeredMutex.Unlock()
	registry.temporaryMutex.Lock()
	defer registry.temporaryMutex.Unlock()
	registry.removedMutex.Lock()
	defer registry.removedMutex.Unlock()
	registry.registeredAddresses = make(map[string]bool)
	registry.temporaryAddresses = make(map[string]bool)
	registry.removedAddresses = nil
}

func (registry *Registry) Copy() verification.RegistrationsManager {
	registry.registeredMutex.RLock()
	defer registry.registeredMutex.RUnlock()
	registry.temporaryMutex.RLock()
	defer registry.temporaryMutex.RUnlock()
	registry.removedMutex.RLock()
	defer registry.removedMutex.RUnlock()
	registryCopy := &Registry{}
	registryCopy.infuraKey = registry.infuraKey
	registryCopy.registeredAddresses = copyAddressesMap(registry.registeredAddresses)
	registryCopy.temporaryAddresses = copyAddressesMap(registry.temporaryAddresses)
	registryCopy.removedAddresses = registry.removedAddresses
	registryCopy.logger = registry.logger
	return registryCopy
}

func (registry *Registry) Filter(addresses []string) []string {
	var newAddresses []string
	for _, address := range addresses {
		if !registry.registeredAddresses[address] {
			newAddresses = append(newAddresses, address)
		}
	}
	return newAddresses
}

func (registry *Registry) IsRegistered(address string) bool {
	return registry.registeredAddresses[address]
}

func (registry *Registry) RemovedAddresses() []string {
	return registry.removedAddresses
}

func (registry *Registry) Synchronize(int64) {
	registry.registeredMutex.RLock()
	defer registry.registeredMutex.RUnlock()
	registry.removedMutex.Lock()
	defer registry.removedMutex.Unlock()
	for address, _ := range registry.registeredAddresses {
		isPohValid, err := registry.isRegistered(address)
		if err != nil {
			registry.logger.Debug(err.Error())
		} else if !isPohValid {
			registry.removedAddresses = append(registry.removedAddresses, address)
		}
	}
	registry.temporaryMutex.Lock()
	defer registry.temporaryMutex.Unlock()
	registry.temporaryAddresses = make(map[string]bool)
}

func (registry *Registry) Update(addedAddresses []string, removedAddresses []string) {
	registry.registeredMutex.Lock()
	defer registry.registeredMutex.Unlock()
	registry.removedMutex.Lock()
	defer registry.removedMutex.Unlock()
	for _, address := range removedAddresses {
		removeAddress(registry.removedAddresses, address)
		delete(registry.registeredAddresses, address)
	}
	for _, address := range addedAddresses {
		registry.registeredAddresses[address] = true
	}
}

func (registry *Registry) Verify(addedAddresses []string, removedAddresses []string) error {
	registry.registeredMutex.RLock()
	defer registry.registeredMutex.RUnlock()
	registry.temporaryMutex.Lock()
	defer registry.temporaryMutex.Unlock()
	for _, address := range removedAddresses {
		if registry.registeredAddresses[address] && registry.temporaryAddresses[address] {
			isPohValid, err := registry.isRegistered(address)
			if err != nil {
				registry.logger.Debug(err.Error())
			} else if isPohValid {
				return fmt.Errorf("a removed address is registered")
			} else {
				registry.temporaryAddresses[address] = false
			}
		}
	}
	for _, address := range addedAddresses {
		if !registry.registeredAddresses[address] && registry.temporaryAddresses[address] {
			isPohValid, err := registry.isRegistered(address)
			if err != nil {
				registry.logger.Debug(err.Error())
			} else if !isPohValid {
				return fmt.Errorf("an added address is not registered")
			} else {
				registry.temporaryAddresses[address] = true
			}
		}
	}
	return nil
}

func copyAddressesMap(addresses map[string]bool) map[string]bool {
	addressesCopy := make(map[string]bool, len(addresses))
	for address := range addresses {
		addressesCopy[address] = true
	}
	return addressesCopy
}

func (registry *Registry) isRegistered(address string) (isRegistered bool, err error) {
	if registry.infuraKey == "" {
		return true, nil
	}
	client, err := ethclient.Dial(fmt.Sprint(clientUrl, registry.infuraKey))
	if err != nil {
		return
	}
	pohSmartContractAddress := common.HexToAddress(pohSmartContractAddressHex)
	proofOfHumanity, err := NewPoh(pohSmartContractAddress, client)
	if err != nil {
		return false, fmt.Errorf("failed to get proof of humanity for address %s: %w", address, err)
	}
	isRegistered, err = proofOfHumanity.PohCaller.IsRegistered(nil, common.HexToAddress(address))
	if err != nil {
		return false, fmt.Errorf("failed to get proof of humanity for address %s: %w", address, err)
	}
	return isRegistered, nil
}

func removeAddress(addresses []string, address string) []string {
	for i := 0; i < len(addresses); i++ {
		if address == addresses[i] {
			return append(addresses[:i], addresses[i+1:]...)
		}
	}
	return addresses
}
