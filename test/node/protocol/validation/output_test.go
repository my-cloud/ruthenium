package validation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/test"
	"math"
	"testing"
)

const (
	base                uint64 = 50000000000
	limit               uint64 = 10000000000000
	genesisTimestamp           = 0
	oneDay                     = 24 * 60 * 60 * 1e9
	halfLife                   = 373.59 * oneDay
	validationTimestamp        = 60 * 1e9
)

//////////////////////////////////// INCOME ////////////////////////////////////
func Test_Value_ValueIsMaxUint64AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     math.MaxUint64, // 18446744073709551615
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(float64(math.MaxUint64-limit))/2 + float64(limit)) // 9223377036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs200kAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = 20000000000000
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterHalfLife := (value + limit) / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsLimitAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     limit,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterOneDay := limit
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	expectedValueAfterHalfLife := limit
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs50kAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     5000000000000,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 5001577710768 // Can be verified graphically at https://www.desmos.com/calculator/utdwq2cdh9
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 5578244486849 // Can be verified graphically at https://www.desmos.com/calculator/utdwq2cdh9
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     1,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 360342 // Can be verified graphically at https://www.desmos.com/calculator/utdwq2cdh9
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 50000445605 // Can be verified graphically at https://www.desmos.com/calculator/utdwq2cdh9
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs0AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     0,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterHalfLife := base
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

//////////////////////////////////// NO INCOME ////////////////////////////////////
func Test_Value_ValueIsMaxUint64AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     math.MaxUint64, // 18446744073709551615
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(math.MaxUint64 / 2)) // 9223372036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs200kAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 20000000000000
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterHalfLife := value / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsLimitAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     limit,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterHalfLife := limit / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs50kAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 5000000000000
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterHalfLife := value / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     1,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 0
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 0
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs0AndHasNoIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     0,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 0
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 0
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}
