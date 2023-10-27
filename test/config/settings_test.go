package config

import (
	"encoding/json"
	"fmt"
	"github.com/my-cloud/ruthenium/src/config"
	"github.com/my-cloud/ruthenium/test"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func Test_NewSettings_UnableToOpenFile_ReturnsError(t *testing.T) {
	// Arrange
	// Act
	_, err := config.NewSettings("")

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
	_, err := config.NewSettings(jsonFileName)

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
	var genesisAmountInParticles uint64 = 2
	var halfLifeInDays float64 = 3
	var incomeBaseInParticles uint64 = 4
	var incomeLimitInParticles uint64 = 5
	var maxOutboundsCount int = 6
	var minimalTransactionFee uint64 = 7
	var particlesPerToken uint64 = 8
	var synchronizationIntervalInSeconds int = 9
	var validationIntervalInSeconds int64 = 10
	var validationTimeoutInSeconds int64 = 11
	var verificationsCountPerValidation int64 = 12
	bytes, _ := json.Marshal(struct {
		BlocksCountLimit                 uint64
		GenesisAmountInParticles         uint64
		HalfLifeInDays                   float64
		IncomeBaseInParticles            uint64
		IncomeLimitInParticles           uint64
		MaxOutboundsCount                int
		MinimalTransactionFee            uint64
		ParticlesPerToken                uint64
		SynchronizationIntervalInSeconds int
		ValidationIntervalInSeconds      int64
		ValidationTimeoutInSeconds       int64
		VerificationsCountPerValidation  int64
	}{
		BlocksCountLimit:                 blocksCountLimit,
		GenesisAmountInParticles:         genesisAmountInParticles,
		HalfLifeInDays:                   halfLifeInDays,
		IncomeBaseInParticles:            incomeBaseInParticles,
		IncomeLimitInParticles:           incomeLimitInParticles,
		MaxOutboundsCount:                maxOutboundsCount,
		MinimalTransactionFee:            minimalTransactionFee,
		ParticlesPerToken:                particlesPerToken,
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
	settings, _ := config.NewSettings(jsonFileName)

	// Assert
	test.Assert(t, settings != nil, "Settings are nil whereas it should not.")
	actualBytes := settings.Bytes()
	actualBlocksCountLimit := settings.BlocksCountLimit()
	actualGenesisAmountInParticles := settings.GenesisAmountInParticles()
	actualHalfLifeInNanoseconds := settings.HalfLifeInNanoseconds()
	actualIncomeBaseInParticles := settings.IncomeBaseInParticles()
	actualIncomeLimitInParticles := settings.IncomeLimitInParticles()
	actualMaxOutboundsCount := settings.MaxOutboundsCount()
	actualMinimalTransactionFee := settings.MinimalTransactionFee()
	actualParticlesPerToken := settings.ParticlesPerToken()
	actualSynchronizationTimer := settings.SynchronizationTimer()
	actualValidationTimer := settings.ValidationTimer()
	actualValidationTimestamp := settings.ValidationTimestamp()
	actualValidationTimeout := settings.ValidationTimeout()
	actualVerificationsCountPerValidation := settings.VerificationsCountPerValidation()
	test.Assert(t, reflect.DeepEqual(actualBytes, bytes), fmt.Sprintf("wrong bytes. expected: %v, actual: %v", bytes, actualBytes))
	test.Assert(t, actualBlocksCountLimit == blocksCountLimit, fmt.Sprintf("wrong blocksCountLimit. expected: %v, actual: %v", blocksCountLimit, actualBlocksCountLimit))
	test.Assert(t, actualGenesisAmountInParticles == genesisAmountInParticles, fmt.Sprintf("wrong genesisAmountInParticles. expected: %v, actual: %v", genesisAmountInParticles, actualGenesisAmountInParticles))
	expectedHalfLifeInNanoseconds := halfLifeInDays * 24 * float64(time.Hour.Nanoseconds())
	test.Assert(t, actualHalfLifeInNanoseconds == expectedHalfLifeInNanoseconds, fmt.Sprintf("wrong halfLifeInNanoseconds. expected: %v, actual: %v", expectedHalfLifeInNanoseconds, actualHalfLifeInNanoseconds))
	test.Assert(t, actualIncomeBaseInParticles == incomeBaseInParticles, fmt.Sprintf("wrong incomeBaseInParticles. expected: %v, actual: %v", incomeBaseInParticles, actualIncomeBaseInParticles))
	test.Assert(t, actualIncomeLimitInParticles == incomeLimitInParticles, fmt.Sprintf("wrong incomeLimitInParticles. expected: %v, actual: %v", incomeLimitInParticles, actualIncomeLimitInParticles))
	test.Assert(t, actualMaxOutboundsCount == maxOutboundsCount, fmt.Sprintf("wrong maxOutboundsCount. expected: %v, actual: %v", maxOutboundsCount, actualMaxOutboundsCount))
	test.Assert(t, actualMinimalTransactionFee == minimalTransactionFee, fmt.Sprintf("wrong minimalTransactionFee. expected: %v, actual: %v", minimalTransactionFee, actualMinimalTransactionFee))
	test.Assert(t, actualParticlesPerToken == particlesPerToken, fmt.Sprintf("wrong particlesPerToken. expected: %v, actual: %v", particlesPerToken, actualParticlesPerToken))
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
