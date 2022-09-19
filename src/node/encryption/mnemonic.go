package encryption

import (
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
	"strconv"
	"strings"
)

const (
	purposeId      = "purpose"
	coinTypeId     = "coin type"
	accountId      = "account"
	changeId       = "change"
	addressIndexId = "address index"
)

type Mnemonic struct {
	phrase string
}

func NewMnemonic(phrase string) *Mnemonic {
	return &Mnemonic{phrase}
}

func (mnemonic *Mnemonic) PrivateKey(derivationPath string, password string) (*PrivateKey, error) {
	indexes, err := parsePath(derivationPath)
	if err != nil {
		return nil, err
	}

	// Generate a new seed for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic.phrase, password)

	// Generate a new master node using the seed.
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	// Derive the master node to extend it with the purpose
	purposeExtendedKey, err := masterKey.Derive(hdkeychain.HardenedKeyStart + indexes[purposeId])
	if err != nil {
		return nil, err
	}

	// Derive the master node to extend it with the coin type
	coinTypeExtendedKey, err := purposeExtendedKey.Derive(hdkeychain.HardenedKeyStart + indexes[coinTypeId])
	if err != nil {
		return nil, err
	}

	// Derive the master node to extend it with the account
	accountExtendedKey, err := coinTypeExtendedKey.Derive(hdkeychain.HardenedKeyStart + indexes[accountId])
	if err != nil {
		return nil, err
	}

	// Derive the master node to extend it with the change
	changeExtendedKey, err := accountExtendedKey.Derive(indexes[changeId])
	if err != nil {
		return nil, err
	}

	// Derive the master node to extend it with the addressIndex
	addressIndexExtendedKey, err := changeExtendedKey.Derive(indexes[addressIndexId])
	if err != nil {
		return nil, err
	}

	btcecPrivateKey, err := addressIndexExtendedKey.ECPrivKey()
	if err != nil {
		return nil, err
	}

	ecdsaPrivateKey := btcecPrivateKey.ToECDSA()

	return &PrivateKey{ecdsaPrivateKey}, nil
}

func parsePath(path string) (map[string]uint32, error) {
	indexes := make(map[string]uint32)
	const derivationStartString = "m/"
	const derivationSeparator1 = "'/"
	const derivationSeparator2 = "/"
	purposeString := path[len(derivationStartString):strings.Index(path, derivationSeparator1)]
	purpose, err := strconv.ParseUint(purposeString, 10, 32)
	if err != nil {
		return nil, err
	}
	indexes[purposeId] = uint32(purpose)
	pathWithoutPurpose := path[strings.Index(path, purposeString)+len(purposeString)+len(derivationSeparator1):]
	coinTypeString := pathWithoutPurpose[:strings.Index(pathWithoutPurpose, derivationSeparator1)]
	coinType, err := strconv.ParseUint(coinTypeString, 10, 32)
	if err != nil {
		return nil, err
	}
	indexes[coinTypeId] = uint32(coinType)
	pathWithoutCoinType := pathWithoutPurpose[strings.Index(pathWithoutPurpose, coinTypeString)+len(coinTypeString)+len(derivationSeparator1):]
	accountString := pathWithoutCoinType[:strings.Index(pathWithoutCoinType, derivationSeparator1)]
	account, err := strconv.ParseUint(accountString, 10, 32)
	if err != nil {
		return nil, err
	}
	indexes[accountId] = uint32(account)
	pathWithoutAccount := pathWithoutCoinType[strings.Index(pathWithoutCoinType, accountString)+len(accountString)+len(derivationSeparator1):]
	changeString := pathWithoutAccount[:strings.Index(pathWithoutAccount, derivationSeparator2)]
	change, err := strconv.ParseUint(changeString, 10, 32)
	if err != nil {
		return nil, err
	}
	indexes[changeId] = uint32(change)
	pathWithoutChange := pathWithoutAccount[strings.Index(pathWithoutAccount, changeString)+len(changeString)+len(derivationSeparator2):]
	addressIndex, err := strconv.ParseUint(pathWithoutChange, 10, 32)
	if err != nil {
		return nil, err
	}
	indexes[addressIndexId] = uint32(addressIndex)
	return indexes, nil
}
