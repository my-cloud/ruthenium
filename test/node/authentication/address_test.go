package authentication

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
)

const (
	purposeId      = "purpose"
	coinTypeId     = "coin type"
	accountId      = "account"
	changeId       = "change"
	addressIndexId = "address index"
)

func Test_PrivateKeyFromMnemonic(t *testing.T) {
	// Arrange
	// Act
	privateKey, _ := privateKeyFromMnemonic("artist silver basket insane canvas top drill social reflect park fruit bless", "m/44'/60'/0'/0/0", "")

	// Assert
	expectedPrivateKey := "0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd"
	privateKeyBytes := crypto.FromECDSA(privateKey)
	actualPrivateKey := hexutil.Encode(privateKeyBytes)
	assert(t, actualPrivateKey == expectedPrivateKey, fmt.Sprintf("Wrong private key. Expected: %s - Actual: %s", expectedPrivateKey, actualPrivateKey))
}

func Test_PublicKeyFromPrivateKey(t *testing.T) {
	bytes, _ := hexutil.Decode("0x48913790c2bebc48417491f96a7e07ec94c76ccd0fe1562dc1749479d9715afd")
	privateKey := crypto.ToECDSAUnsafe(bytes)

	// Act
	publicKey := privateKey.Public()

	// Assert
	expectedPublicKey := "0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782"
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	actualPublicKey := hexutil.Encode(publicKeyBytes)
	assert(t, actualPublicKey == expectedPublicKey, fmt.Sprintf("Wrong public key. Expected: %s - Actual: %s", expectedPublicKey, actualPublicKey))
}

func Test_AddressFromPublicKey(t *testing.T) {
	// Arrange
	bytes, _ := hexutil.Decode("0x046bd857ce80ff5238d6561f3a775802453c570b6ea2cbf93a35a8a6542b2edbe5f625f9e3fbd2a5df62adebc27391332a265fb94340fb11b69cf569605a5df782")
	publicKey, _ := crypto.UnmarshalPubkey(bytes)

	// Act
	address := crypto.PubkeyToAddress(*publicKey).Hex()

	// Assert
	expectedAddress := "0x9C69443c3Ec0D660e257934ffc1754EB9aD039CB"
	assert(t, address == expectedAddress, fmt.Sprintf("Wrong address. Expected: %s - Actual: %s", expectedAddress, address))
}

func assert(t testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		t.FailNow()
	}
}

func privateKeyFromMnemonic(mnemonic string, path string, password string) (*ecdsa.PrivateKey, error) {
	indexes := parsePath(path)
	// Generate a Bip32 HD wallet for the mnemonic and a user supplied password
	seed := bip39.NewSeed(mnemonic, password)

	// Generate a new master node using the seed.
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H
	acc44H, err := masterKey.Derive(hdkeychain.HardenedKeyStart + indexes[purposeId])
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H
	acc44H60H, err := acc44H.Derive(hdkeychain.HardenedKeyStart + indexes[coinTypeId])
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H/0H
	acc44H60H0H, err := acc44H60H.Derive(hdkeychain.HardenedKeyStart + indexes[accountId])
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H/0H/0
	acc44H60H0H0, err := acc44H60H0H.Derive(indexes[changeId])
	if err != nil {
		return nil, err
	}

	// This gives the path: m/44H/60H/0H/0/0
	acc44H60H0H00, err := acc44H60H0H0.Derive(indexes[addressIndexId])
	if err != nil {
		return nil, err
	}

	btcecPrivateKey, err := acc44H60H0H00.ECPrivKey()
	if err != nil {
		return nil, err
	}

	privateKey := btcecPrivateKey.ToECDSA()

	return privateKey, nil
}

func parsePath(path string) map[string]uint32 {
	indexes := make(map[string]uint32)
	const derivationStartString = "m/"
	const derivationSeparator1 = "'/"
	const derivationSeparator2 = "/"
	purposeString := path[len(derivationStartString):strings.Index(path, derivationSeparator1)]
	purpose, _ := strconv.Atoi(purposeString)
	indexes[purposeId] = uint32(purpose)
	pathWithoutPurpose := path[strings.Index(path, purposeString)+len(purposeString)+len(derivationSeparator1):]
	coinTypeString := pathWithoutPurpose[:strings.Index(pathWithoutPurpose, derivationSeparator1)]
	coinType, _ := strconv.Atoi(coinTypeString)
	indexes[coinTypeId] = uint32(coinType)
	pathWithoutCoinType := pathWithoutPurpose[strings.Index(pathWithoutPurpose, coinTypeString)+len(coinTypeString)+len(derivationSeparator1):]
	accountString := pathWithoutCoinType[:strings.Index(pathWithoutCoinType, derivationSeparator1)]
	account, _ := strconv.Atoi(accountString)
	indexes[accountId] = uint32(account)
	pathWithoutAccount := pathWithoutCoinType[strings.Index(pathWithoutCoinType, accountString)+len(accountString)+len(derivationSeparator1):]
	changeString := pathWithoutAccount[:strings.Index(pathWithoutAccount, derivationSeparator2)]
	change, _ := strconv.Atoi(changeString)
	indexes[changeId] = uint32(change)
	pathWithoutChange := pathWithoutAccount[strings.Index(pathWithoutAccount, changeString)+len(changeString)+len(derivationSeparator2):]
	addressIndex, _ := strconv.Atoi(pathWithoutChange)
	indexes[addressIndexId] = uint32(addressIndex)
	return indexes
}
