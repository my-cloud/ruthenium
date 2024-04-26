package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"testing"
)

// TODO
// func Test_UtxosByAddress_UnknownAddress_ReturnsEmptyArray(t *testing.T) {
//	// Arrange
//	logger := log.NewLoggerMock()
//	genesisValidatorAddress := ""
//	var genesisAmount uint64 = 0
//	settings := new(validatornode.SettingsProviderMock)
//	settings.GenesisAmountFunc = func() uint64 { return genesisAmount }
//	blockchain := NewBlockchain(nil, settings, nil, logger)
//
//	// Act
//	utxosBytes := blockchain.Utxos(genesisValidatorAddress)
//
//	// Assert
//	var utxos []*ledger.Utxo
//	_ = json.Unmarshal(utxosBytes, &utxos)
//	test.Assert(t, len(utxos) == 0, "utxos should be empty")
// }
//
// func Test_Utxos_UtxoExists_ReturnsUtxo(t *testing.T) {
//	// Arrange
//	registry := new(validatornode.AddressesManagerMock)
//	registry.IsRegisteredFunc = func(string) (bool, error) { return true, nil }
//	logger := log.NewLoggerMock()
//	sendersManagerMock := new(network.SendersManagerMock)
//	var validationInterval int64 = 1
//	settings := new(validatornode.SettingsProviderMock)
//	settings.GenesisAmountFunc = func() uint64 { return 1 }
//	blockchain := NewBlockchain(registry, settings, sendersManagerMock, logger)
//	registeredAddress := ""
//	var expectedValue uint64 = 1
//	var genesisTimestamp int64 = 0
//	transaction, _ := ledger.NewRewardTransaction(registeredAddress, true, genesisTimestamp+validationInterval, expectedValue)
//	transactions := []*ledger.Transaction{transaction}
//	transactionsBytes, _ := json.Marshal(transactions)
//	_ = blockchain.AddBlock(genesisTimestamp, nil, nil)
//	_ = blockchain.AddBlock(genesisTimestamp+validationInterval, transactionsBytes, []string{registeredAddress})
//	_ = blockchain.AddBlock(genesisTimestamp+2*validationInterval, nil, nil)
//
//	// Act
//	utxosBytes := blockchain.Utxos(registeredAddress)
//
//	// Assert
//	var utxos []*ledger.Utxo
//	_ = json.Unmarshal(utxosBytes, &utxos)
//	actualValue := utxos[0].Value(genesisTimestamp+2*validationInterval, 1, 1, 1)
//	test.Assert(t, actualValue == expectedValue, fmt.Sprintf("utxo amount is %d whereas it should be %d", actualValue, expectedValue))
// }

func Test_Utxos_UnknownAddress_ReturnsEmptyArray(t *testing.T) {
	// Arrange
	registry := NewUtxosRegistry(new(application.ProtocolSettingsProviderMock))

	// Act
	actualUtxos := registry.Utxos("")

	// Assert
	actualUtxosLength := len(actualUtxos)
	expectedUtxosLength := 0
	test.Assert(t, actualUtxosLength == expectedUtxosLength, fmt.Sprintf("utxos length is %d whereas it should be %d", actualUtxosLength, expectedUtxosLength))
}

func Test_Utxos_OneCorrespondingUtxo_ReturnsArrayWithOneUtxo(t *testing.T) {
	// Arrange
	address := ""
	initialUtxos := utxosRegistrationInfo{
		address,
		"",
		[]*ledger.Utxo{ledger.NewUtxo(nil, ledger.NewOutput(address, false, 1), 0)},
	}
	registry := NewUtxosRegistry(new(application.ProtocolSettingsProviderMock), initialUtxos)

	// Act
	actualUtxos := registry.Utxos("")

	// Assert
	actualUtxosLength := len(actualUtxos)
	expectedUtxosLength := 1
	test.Assert(t, actualUtxosLength == expectedUtxosLength, fmt.Sprintf("utxo length is %d whereas it should be %d", actualUtxosLength, expectedUtxosLength))
}
