package verification

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
	"github.com/my-cloud/ruthenium/test"
	"math"
	"testing"
)

const (
	particlesCount          = 100000000
	base             uint64 = 500 * particlesCount
	limit            uint64 = 100000 * particlesCount
	genesisTimestamp        = 0
	oneMinute               = 60 * 1e9
	oneDay                  = 24 * 60 * oneMinute
	halfLife                = 373.59 * oneDay
)

// ////////////////////////////////// WITH INCOME ////////////////////////////////////
func Test_Value_ValueIsMaxUint64AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = math.MaxUint64 // 18446744073709551615
	output := verification.NewOutput("", true, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(float64(math.MaxUint64-limit))/2 + float64(limit)) // 9223377036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsTwiceTheLimitAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	value := 2 * limit
	output := verification.NewOutput("", true, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := (value + limit) / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsLimitAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	value := limit
	output := verification.NewOutput("", true, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterOneDay := utxo.Value(int64(genesisTimestamp+oneDay), halfLife, base, limit)
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterOneDay := limit
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	expectedValueAfterHalfLife := limit
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = 1
	output := verification.NewOutput("", true, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterNoTimeElapsed := utxo.Value(int64(genesisTimestamp), halfLife, base, limit)
	actualValueAfter1Minute := utxo.Value(genesisTimestamp+oneMinute, halfLife, base, limit)

	// Assert
	expectedValueAfterNoTimeElapsed := value
	test.Assert(t, actualValueAfterNoTimeElapsed == expectedValueAfterNoTimeElapsed, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterNoTimeElapsed, actualValueAfterNoTimeElapsed))
	test.Assert(t, actualValueAfter1Minute >= value, fmt.Sprintf("Wrong value. Expected a value >= %d - Actual: %d", value, actualValueAfter1Minute))
}

func Test_Value_ValueIs0AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := verification.NewOutput("", true, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := base
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_SamePortionOfHalfLifeElapsedWithIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := verification.NewOutput("", true, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)
	portion := 1.5

	// Act
	value1 := utxo.Value(int64(genesisTimestamp+portion*halfLife), halfLife, base, limit)
	value2 := utxo.Value(int64(genesisTimestamp+portion*oneDay), oneDay, base, limit)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}

func Test_Value_SameElapsedTimeWithIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	var value uint64 = 0
	var outputTimestamp1 int64 = 0
	var outputTimestamp2 int64 = 10
	output := verification.NewOutput("", true, value)
	utxo1 := verification.NewUtxo(&verification.InputInfo{}, output, outputTimestamp1)
	utxo2 := verification.NewUtxo(&verification.InputInfo{}, output, outputTimestamp2)
	var elapsedTimestamp int64 = oneDay
	utxo1Timestamp := outputTimestamp1
	utxo2Timestamp := outputTimestamp2

	// Act
	value1 := utxo1.Value(utxo1Timestamp+elapsedTimestamp, halfLife, base, limit)
	value2 := utxo2.Value(utxo2Timestamp+elapsedTimestamp, halfLife, base, limit)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}

// ////////////////////////////////// WITHOUT INCOME ////////////////////////////////////
func Test_Value_ValueIsMaxUint64AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = math.MaxUint64 // 18446744073709551615
	output := verification.NewOutput("", false, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(math.MaxUint64 / 2)) // 9223372036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsTwiceTheLimitAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	value := 2 * limit
	output := verification.NewOutput("", false, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := value / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsLimitAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	value := limit
	output := verification.NewOutput("", false, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := limit / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 1
	output := verification.NewOutput("", false, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterNoTimeElapsed := utxo.Value(int64(genesisTimestamp), halfLife, base, limit)
	actualValueAfter1Minute := utxo.Value(genesisTimestamp+oneMinute, halfLife, base, limit)

	// Assert
	expectedValueAfterNoTimeElapsed := value
	test.Assert(t, actualValueAfterNoTimeElapsed == expectedValueAfterNoTimeElapsed, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterNoTimeElapsed, actualValueAfterNoTimeElapsed))
	var expectedValueAfter1Minute uint64 = 0
	test.Assert(t, actualValueAfter1Minute == expectedValueAfter1Minute, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfter1Minute, actualValueAfter1Minute))
}

func Test_Value_ValueIs0AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := verification.NewOutput("", false, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)

	// Act
	actualValueAfterOneDay := utxo.Value(int64(oneDay), halfLife, base, limit)
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	var expectedValueAfterOneDay uint64 = 0
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 0
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_SamePortionOfHalfLifeElapsedWithoutIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := verification.NewOutput("", false, value)
	utxo := verification.NewUtxo(&verification.InputInfo{}, output, 0)
	portion := 1.5

	// Act
	value1 := utxo.Value(int64(genesisTimestamp+portion*halfLife), halfLife, base, limit)
	value2 := utxo.Value(int64(genesisTimestamp+portion*oneDay), oneDay, base, limit)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}

func Test_Value_SameElapsedTimeWithoutIncome_ReturnsSameValue(t *testing.T) {
	// Arrange
	var value uint64 = 0
	var outputTimestamp1 int64 = 0
	var outputTimestamp2 int64 = 10
	output := verification.NewOutput("", false, value)
	utxo1 := verification.NewUtxo(&verification.InputInfo{}, output, outputTimestamp1)
	utxo2 := verification.NewUtxo(&verification.InputInfo{}, output, outputTimestamp2)
	var elapsedTimestamp int64 = oneDay
	utxo1Timestamp := outputTimestamp1
	utxo2Timestamp := outputTimestamp2

	// Act
	value1 := utxo1.Value(utxo1Timestamp+elapsedTimestamp, halfLife, base, limit)
	value2 := utxo2.Value(utxo2Timestamp+elapsedTimestamp, halfLife, base, limit)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}
