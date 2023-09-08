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

//////////////////////////////////// WITH INCOME ////////////////////////////////////
func Test_Value_ValueIsMaxUint64AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     math.MaxUint64, // 18446744073709551615
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(float64(math.MaxUint64-limit))/2 + float64(limit)) // 9223377036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsTwiceTheLimitAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	value := 2 * limit
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

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
	actualValueAfterOneDay := output.Value(int64(genesisTimestamp + oneDay))
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

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
	actualValueAfterOneDay := output.Value(int64(genesisTimestamp + oneDay))
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 5001577710768
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 5578244486849
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = 1
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterNoTimeElapsed := output.Value(int64(genesisTimestamp))
	actualValueAfter1Minute := output.Value(genesisTimestamp + validationTimestamp)

	// Assert
	expectedValueAfterNoTimeElapsed := value
	test.Assert(t, actualValueAfterNoTimeElapsed == expectedValueAfterNoTimeElapsed, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterNoTimeElapsed, actualValueAfterNoTimeElapsed))
	test.Assert(t, actualValueAfter1Minute >= value, fmt.Sprintf("Wrong value. Expected a value >= %d - Actual: %d", value, actualValueAfter1Minute))
}

func Test_Value_ValueIs0AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     0,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

	// Assert
	expectedValueAfterHalfLife := base
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_SamePortionOfHalfLifeElapsedWithIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     0,
	}
	output1 := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)
	output2 := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, oneDay, base, limit, validationTimestamp)
	portion := 1.5

	// Act
	value1 := output1.Value(int64(genesisTimestamp + portion*halfLife))
	value2 := output2.Value(int64(genesisTimestamp + portion*oneDay))

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}

func Test_Value_SameElapsedTimeWithIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	blockHeight1 := 0
	utxo1 := &network.UtxoResponse{
		BlockHeight: blockHeight1,
		HasIncome:   true,
		Value:       0,
	}
	blockHeight2 := 10
	utxo2 := &network.UtxoResponse{
		BlockHeight: blockHeight2,
		HasIncome:   true,
		Value:       0,
	}
	output1 := validation.NewOutputFromUtxoResponse(utxo1, genesisTimestamp, halfLife, base, limit, validationTimestamp)
	output2 := validation.NewOutputFromUtxoResponse(utxo2, genesisTimestamp, halfLife, base, limit, validationTimestamp)
	var elapsedTimestamp int64 = oneDay
	utxo1Timestamp := int64(blockHeight1 * validationTimestamp)
	utxo2Timestamp := int64(blockHeight2 * validationTimestamp)

	// Act
	value1 := output1.Value(utxo1Timestamp + elapsedTimestamp)
	value2 := output2.Value(utxo2Timestamp + elapsedTimestamp)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}

//////////////////////////////////// WITHOUT INCOME ////////////////////////////////////
func Test_Value_ValueIsMaxUint64AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     math.MaxUint64, // 18446744073709551615
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(math.MaxUint64 / 2)) // 9223372036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsTwiceTheLimitAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	value := 2 * limit
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

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
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

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
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

	// Assert
	expectedValueAfterHalfLife := value / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 1
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterNoTimeElapsed := output.Value(int64(genesisTimestamp))
	actualValueAfter1Minute := output.Value(genesisTimestamp + validationTimestamp)

	// Assert
	expectedValueAfterNoTimeElapsed := value
	test.Assert(t, actualValueAfterNoTimeElapsed == expectedValueAfterNoTimeElapsed, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterNoTimeElapsed, actualValueAfterNoTimeElapsed))
	var expectedValueAfter1Minute uint64 = 0
	test.Assert(t, actualValueAfter1Minute == expectedValueAfter1Minute, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfter1Minute, actualValueAfter1Minute))
}

func Test_Value_ValueIs0AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     0,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(genesisTimestamp + halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 0
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 0
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_SamePortionOfHalfLifeElapsedWithoutIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     0,
	}
	output1 := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, base, limit, validationTimestamp)
	output2 := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, oneDay, base, limit, validationTimestamp)
	portion := 1.5

	// Act
	value1 := output1.Value(int64(genesisTimestamp + portion*halfLife))
	value2 := output2.Value(int64(genesisTimestamp + portion*oneDay))

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}

func Test_Value_SameElapsedTimeWithoutIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	blockHeight1 := 0
	utxo1 := &network.UtxoResponse{
		BlockHeight: blockHeight1,
		HasIncome:   false,
		Value:       0,
	}
	blockHeight2 := 10
	utxo2 := &network.UtxoResponse{
		BlockHeight: blockHeight2,
		HasIncome:   false,
		Value:       0,
	}
	output1 := validation.NewOutputFromUtxoResponse(utxo1, genesisTimestamp, halfLife, base, limit, validationTimestamp)
	output2 := validation.NewOutputFromUtxoResponse(utxo2, genesisTimestamp, halfLife, base, limit, validationTimestamp)
	var elapsedTimestamp int64 = oneDay
	utxo1Timestamp := int64(blockHeight1 * validationTimestamp)
	utxo2Timestamp := int64(blockHeight2 * validationTimestamp)

	// Act
	value1 := output1.Value(utxo1Timestamp + elapsedTimestamp)
	value2 := output2.Value(utxo2Timestamp + elapsedTimestamp)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}
