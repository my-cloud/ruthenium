package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
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

func Test_CalculateFee_UnknownTransactionId_ReturnsError(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	transactionId := ""
	initialUtxos := utxosRegistrationInfo{
		"",
		transactionId,
		[]*ledger.Utxo{ledger.NewUtxo(nil, ledger.NewOutput("", false, 1), 0)},
	}
	settingsMock := new(application.ProtocolSettingsProviderMock)
	registry := NewUtxosRegistry(settingsMock, initialUtxos)
	transaction := ledger.NewSignedTransaction(1, 0, 0, "", privateKey, publicKey, 0, "unknown", 1, false)

	// Act
	_, err := registry.CalculateFee(transaction, 0)

	// Assert
	if err == nil {
		test.Assert(t, false, "error was nil whereas it should not")
		return
	} else {
		test.AssertThatMessageIsLogged(t, []struct{ Msg string }{{Msg: err.Error()}}, "failed to find transaction ID")
	}
}

func Test_CalculateFee_UnknownOutputIndex_ReturnsError(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	transactionId := ""
	initialUtxos := utxosRegistrationInfo{
		"",
		transactionId,
		[]*ledger.Utxo{ledger.NewUtxo(nil, ledger.NewOutput("", false, 1), 0)},
	}
	settingsMock := new(application.ProtocolSettingsProviderMock)
	registry := NewUtxosRegistry(settingsMock, initialUtxos)
	transaction := ledger.NewSignedTransaction(1, 0, 1, "", privateKey, publicKey, 0, transactionId, 1, false)

	// Act
	_, err := registry.CalculateFee(transaction, 0)

	// Assert
	if err == nil {
		test.Assert(t, false, "error was nil whereas it should not")
		return
	} else {
		test.AssertThatMessageIsLogged(t, []struct{ Msg string }{{Msg: err.Error()}}, "failed to find output index")
	}
}

func Test_CalculateFee_WrongRecipientAddress_ReturnsError(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	transactionId := ""
	initialUtxos := utxosRegistrationInfo{
		"",
		transactionId,
		[]*ledger.Utxo{ledger.NewUtxo(nil, ledger.NewOutput("", false, 1), 0)},
	}
	settingsMock := new(application.ProtocolSettingsProviderMock)
	registry := NewUtxosRegistry(settingsMock, initialUtxos)
	transaction := ledger.NewSignedTransaction(1, 0, 0, "", privateKey, publicKey, 0, transactionId, 1, false)

	// Act
	_, err := registry.CalculateFee(transaction, 0)

	// Assert
	if err == nil {
		test.Assert(t, false, "error was nil whereas it should not")
		return
	} else {
		test.AssertThatMessageIsLogged(t, []struct{ Msg string }{{Msg: err.Error()}}, "failed to verify input recipient address")
	}
}

func Test_CalculateFee_FeeIsNegative_ReturnsError(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	address := publicKey.Address()
	transactionId := ""
	initialUtxos := utxosRegistrationInfo{
		address,
		transactionId,
		[]*ledger.Utxo{ledger.NewUtxo(nil, ledger.NewOutput(address, false, 1), 0)},
	}
	settingsMock := new(application.ProtocolSettingsProviderMock)
	settingsMock.HalfLifeInNanosecondsFunc = func() float64 { return 1 }
	settingsMock.IncomeBaseFunc = func() uint64 { return 1 }
	settingsMock.IncomeLimitFunc = func() uint64 { return 1 }
	settingsMock.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	registry := NewUtxosRegistry(settingsMock, initialUtxos)
	transaction := ledger.NewSignedTransaction(1, -1, 0, "", privateKey, publicKey, 0, transactionId, 1, false)

	// Act
	_, err := registry.CalculateFee(transaction, 0)

	// Assert
	if err == nil {
		test.Assert(t, false, "error was nil whereas it should not")
		return
	} else {
		test.AssertThatMessageIsLogged(t, []struct{ Msg string }{{Msg: err.Error()}}, "fee is negative")
	}
}

func Test_CalculateFee_FeeIsTooLow_ReturnsError(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	address := publicKey.Address()
	transactionId := ""
	initialUtxos := utxosRegistrationInfo{
		address,
		transactionId,
		[]*ledger.Utxo{ledger.NewUtxo(nil, ledger.NewOutput(address, false, 1), 0)},
	}
	settingsMock := new(application.ProtocolSettingsProviderMock)
	settingsMock.HalfLifeInNanosecondsFunc = func() float64 { return 1 }
	settingsMock.IncomeBaseFunc = func() uint64 { return 1 }
	settingsMock.IncomeLimitFunc = func() uint64 { return 1 }
	settingsMock.MinimalTransactionFeeFunc = func() uint64 { return 1 }
	registry := NewUtxosRegistry(settingsMock, initialUtxos)
	transaction := ledger.NewSignedTransaction(1, 0, 0, "", privateKey, publicKey, 0, transactionId, 1, false)

	// Act
	_, err := registry.CalculateFee(transaction, 0)

	// Assert
	if err == nil {
		test.Assert(t, false, "error was nil whereas it should not")
		return
	} else {
		test.AssertThatMessageIsLogged(t, []struct{ Msg string }{{Msg: err.Error()}}, "fee is too low")
	}
}

func Test_CalculateFee_ValidTransaction_ReturnsFee(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	address := publicKey.Address()
	transactionId := ""
	initialUtxos := utxosRegistrationInfo{
		address,
		transactionId,
		[]*ledger.Utxo{ledger.NewUtxo(nil, ledger.NewOutput(address, false, 1), 0)},
	}
	settingsMock := new(application.ProtocolSettingsProviderMock)
	settingsMock.HalfLifeInNanosecondsFunc = func() float64 { return 1 }
	settingsMock.IncomeBaseFunc = func() uint64 { return 1 }
	settingsMock.IncomeLimitFunc = func() uint64 { return 1 }
	settingsMock.MinimalTransactionFeeFunc = func() uint64 { return 0 }
	registry := NewUtxosRegistry(settingsMock, initialUtxos)
	transaction := ledger.NewSignedTransaction(1, 1, 0, "", privateKey, publicKey, 0, transactionId, 0, false)

	// Act
	actualFee, _ := registry.CalculateFee(transaction, 0)

	// Assert
	var expectedFee uint64 = 1
	test.Assert(t, actualFee == expectedFee, fmt.Sprintf("fee is %d whereas it should be %d", actualFee, expectedFee))
}

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
