package configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_NewSettings_UnableToOpenFile_ReturnsError(t *testing.T) {
	// Arrange
	// Act
	_, err := NewSettings("")

	// Assert
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
	if err != nil {
		expectedErrorMessage := "unable to open file"
		actualErrorMessage := err.Error()
		test.Assert(t, strings.Contains(actualErrorMessage, expectedErrorMessage), fmt.Sprintf("Wrong error message.\nExpected: %s\nActual:   %s", expectedErrorMessage, actualErrorMessage))
	}
}

func Test_NewSettings_UnableToUnmarshalBytes_ReturnsError(t *testing.T) {
	// Arrange
	jsonFile, _ := os.CreateTemp("", "Test_NewSettings_UnableToUnmarshalBytes_ReturnsError.json")
	jsonFileName := jsonFile.Name()
	defer func() { _ = os.Remove(jsonFileName) }()
	jsonData := []byte(`{`)
	_, _ = jsonFile.Write(jsonData)
	_ = jsonFile.Close()

	// Act
	_, err := NewSettings(jsonFileName)

	// Assert
	test.Assert(t, err != nil, "Error is nil whereas it should not.")
	if err != nil {
		expectedErrorMessage := "unable to unmarshal"
		actualErrorMessage := err.Error()
		test.Assert(t, strings.Contains(actualErrorMessage, expectedErrorMessage), fmt.Sprintf("Wrong error message.\nExpected: %s\nActual:   %s", expectedErrorMessage, actualErrorMessage))
	}
}

func Test_NewSettings_ValidBytes_NoError(t *testing.T) {
	// Arrange
	var blocksCountLimit uint64 = 1
	var genesisAmount uint64 = 2
	var halfLifeInDays float64 = 3
	var incomeBase uint64 = 4
	var incomeLimit uint64 = 5
	var maxOutboundsCount int = 6
	var minimalTransactionFee uint64 = 7
	var smallestUnitsPerCoin uint64 = 8
	var synchronizationIntervalInSeconds int = 9
	var validationIntervalInSeconds int64 = 10
	var validationTimeoutInSeconds int64 = 11
	var verificationsCountPerValidation int64 = 12
	bytes, _ := json.Marshal(struct {
		BlocksCountLimit                 uint64
		GenesisAmount                    uint64
		HalfLifeInDays                   float64
		IncomeBase                       uint64
		IncomeLimit                      uint64
		MaxOutboundsCount                int
		MinimalTransactionFee            uint64
		SmallestUnitsPerCoin             uint64
		SynchronizationIntervalInSeconds int
		ValidationIntervalInSeconds      int64
		ValidationTimeoutInSeconds       int64
		VerificationsCountPerValidation  int64
	}{
		BlocksCountLimit:                 blocksCountLimit,
		GenesisAmount:                    genesisAmount,
		HalfLifeInDays:                   halfLifeInDays,
		IncomeBase:                       incomeBase,
		IncomeLimit:                      incomeLimit,
		MaxOutboundsCount:                maxOutboundsCount,
		MinimalTransactionFee:            minimalTransactionFee,
		SmallestUnitsPerCoin:             smallestUnitsPerCoin,
		SynchronizationIntervalInSeconds: synchronizationIntervalInSeconds,
		ValidationIntervalInSeconds:      validationIntervalInSeconds,
		ValidationTimeoutInSeconds:       validationTimeoutInSeconds,
		VerificationsCountPerValidation:  verificationsCountPerValidation,
	})
	jsonFile, _ := os.CreateTemp("", "Test_NewSettings_ValidBytes_NoError.json")
	jsonFileName := jsonFile.Name()
	defer func() { _ = os.Remove(jsonFileName) }()
	_, _ = jsonFile.Write(bytes)
	_ = jsonFile.Close()

	// Act
	settings, _ := NewSettings(jsonFileName)

	// Assert
	test.Assert(t, settings != nil, "Settings are nil whereas it should not.")
	actualBytes := settings.Bytes()
	actualBlocksCountLimit := settings.BlocksCountLimit()
	actualGenesisAmount := settings.GenesisAmount()
	actualHalfLifeInNanoseconds := settings.HalfLifeInNanoseconds()
	actualIncomeBase := settings.IncomeBase()
	actualIncomeLimit := settings.IncomeLimit()
	actualMaxOutboundsCount := settings.MaxOutboundsCount()
	actualMinimalTransactionFee := settings.MinimalTransactionFee()
	actualSmallestUnitsPerCoin := settings.SmallestUnitsPerCoin()
	actualSynchronizationTimer := settings.SynchronizationTimer()
	actualValidationTimer := settings.ValidationTimer()
	actualValidationTimestamp := settings.ValidationTimestamp()
	actualValidationTimeout := settings.ValidationTimeout()
	actualVerificationsCountPerValidation := settings.VerificationsCountPerValidation()
	test.Assert(t, reflect.DeepEqual(actualBytes, bytes), fmt.Sprintf("wrong bytes. expected: %v, actual: %v", bytes, actualBytes))
	test.Assert(t, actualBlocksCountLimit == blocksCountLimit, fmt.Sprintf("wrong blocksCountLimit. expected: %v, actual: %v", blocksCountLimit, actualBlocksCountLimit))
	test.Assert(t, actualGenesisAmount == genesisAmount, fmt.Sprintf("wrong genesisAmount. expected: %v, actual: %v", genesisAmount, actualGenesisAmount))
	expectedHalfLifeInNanoseconds := halfLifeInDays * 24 * float64(time.Hour.Nanoseconds())
	test.Assert(t, actualHalfLifeInNanoseconds == expectedHalfLifeInNanoseconds, fmt.Sprintf("wrong halfLifeInNanoseconds. expected: %v, actual: %v", expectedHalfLifeInNanoseconds, actualHalfLifeInNanoseconds))
	test.Assert(t, actualIncomeBase == incomeBase, fmt.Sprintf("wrong incomeBase. expected: %v, actual: %v", incomeBase, actualIncomeBase))
	test.Assert(t, actualIncomeLimit == incomeLimit, fmt.Sprintf("wrong incomeLimit. expected: %v, actual: %v", incomeLimit, actualIncomeLimit))
	test.Assert(t, actualMaxOutboundsCount == maxOutboundsCount, fmt.Sprintf("wrong maxOutboundsCount. expected: %v, actual: %v", maxOutboundsCount, actualMaxOutboundsCount))
	test.Assert(t, actualMinimalTransactionFee == minimalTransactionFee, fmt.Sprintf("wrong minimalTransactionFee. expected: %v, actual: %v", minimalTransactionFee, actualMinimalTransactionFee))
	test.Assert(t, actualSmallestUnitsPerCoin == smallestUnitsPerCoin, fmt.Sprintf("wrong smallestUnitsPerCoin. expected: %v, actual: %v", smallestUnitsPerCoin, actualSmallestUnitsPerCoin))
	expectedSynchronizationTimer := time.Duration(synchronizationIntervalInSeconds) * time.Second
	test.Assert(t, actualSynchronizationTimer == expectedSynchronizationTimer, fmt.Sprintf("wrong synchronizationTimer. expected: %v, actual: %v", expectedSynchronizationTimer, actualSynchronizationTimer))
	expectedValidationTimer := time.Duration(validationIntervalInSeconds) * time.Second
	test.Assert(t, actualValidationTimer == expectedValidationTimer, fmt.Sprintf("wrong validationTimer. expected: %v, actual: %v", expectedValidationTimer, actualValidationTimer))
	expectedValidationTimestamp := validationIntervalInSeconds * time.Second.Nanoseconds()
	test.Assert(t, actualValidationTimestamp == expectedValidationTimestamp, fmt.Sprintf("wrong validationTimestamp. expected: %v, actual: %v", expectedValidationTimestamp, actualValidationTimestamp))
	expectedValidationTimeout := time.Duration(validationTimeoutInSeconds) * time.Second
	test.Assert(t, actualValidationTimeout == expectedValidationTimeout, fmt.Sprintf("wrong validationTimeout. expected: %v, actual: %v", expectedValidationTimeout, actualValidationTimeout))
	test.Assert(t, actualVerificationsCountPerValidation == verificationsCountPerValidation, fmt.Sprintf("wrong verificationsCountPerValidation. expected: %v, actual: %v", verificationsCountPerValidation, actualVerificationsCountPerValidation))
}
