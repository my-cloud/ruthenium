package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
	"github.com/my-cloud/ruthenium/validatornode/domain/ledger"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"testing"
)

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

func Test_UpdateUtxos_ValidTransactions_ReturnsNil(t *testing.T) {
	// Arrange
	privateKey, _ := encryption.NewPrivateKeyFromHex(test.PrivateKey)
	publicKey := encryption.NewPublicKey(privateKey)
	address := publicKey.Address()
	transactionId := ""
	initialUtxos := utxosRegistrationInfo{
		address,
		transactionId,
		[]*ledger.Utxo{ledger.NewUtxo(ledger.NewInputInfo(0, ""), ledger.NewOutput(address, false, 1), 0)},
	}
	registry := NewUtxosRegistry(new(application.ProtocolSettingsProviderMock), initialUtxos)
	transaction := ledger.NewSignedTransaction(1, 1, 0, address, privateKey, publicKey, 0, transactionId, 0, false)

	// Act
	err := registry.UpdateUtxos([]*ledger.Transaction{transaction}, 0)

	// Assert
	test.Assert(t, err == nil, fmt.Errorf("error should be nil but was: %w", err).Error())
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
