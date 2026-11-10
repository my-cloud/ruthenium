package ledger

import (
	"fmt"
	"math"
	"testing"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
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
func Test_Value_ValueIsMaxUint64AndIsYielding_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = math.MaxUint64 // 18446744073709551615
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(float64(math.MaxUint64-limit))/2 + float64(limit)) // 9223377036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsTwiceTheLimitAndIsYielding_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	value := 2 * limit
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := (value + limit) / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsLimitAndIsYielding_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	value := limit
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterOneDay := utxo.Value(int64(genesisTimestamp+oneDay), halfLife, base, limit)
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterOneDay := limit
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	expectedValueAfterHalfLife := limit
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndIsYielding_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = 1
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterNoTimeElapsed := utxo.Value(int64(genesisTimestamp), halfLife, base, limit)
	actualValueAfter1Minute := utxo.Value(genesisTimestamp+oneMinute, halfLife, base, limit)

	// Assert
	expectedValueAfterNoTimeElapsed := value
	test.Assert(t, actualValueAfterNoTimeElapsed == expectedValueAfterNoTimeElapsed, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterNoTimeElapsed, actualValueAfterNoTimeElapsed))
	test.Assert(t, actualValueAfter1Minute >= value, fmt.Sprintf("Wrong value. Expected a value >= %d - Actual: %d", value, actualValueAfter1Minute))
}

func Test_Value_ValueIs0AndIsYielding_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := base
	isFloatingPointErrorAtMostOneSmallestUnit := actualValueAfterHalfLife <= expectedValueAfterHalfLife && expectedValueAfterHalfLife-actualValueAfterHalfLife <= 1
	test.Assert(t, isFloatingPointErrorAtMostOneSmallestUnit, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsBelowLimitAndIsYielding_ReturnsExpectedValues(t *testing.T) {
	// Arrange
	var incomeBase uint64 = 1000 * particlesCount
	var incomeLimit uint64 = 20000 * particlesCount
	emptyUtxo := NewUtxo(&InputInfo{}, NewOutput("", true, 0), 0)
	halfLimitUtxo := NewUtxo(&InputInfo{}, NewOutput("", true, incomeLimit/2), 0)
	nearLimitUtxo := NewUtxo(&InputInfo{}, NewOutput("", true, incomeLimit-particlesCount), 0)

	// Act
	actualEmptyValueAfter30Days := emptyUtxo.Value(int64(genesisTimestamp+30*oneDay), halfLife, incomeBase, incomeLimit)
	actualHalfLimitValueAfter30Days := halfLimitUtxo.Value(int64(genesisTimestamp+30*oneDay), halfLife, incomeBase, incomeLimit)
	actualNearLimitValueAfterOneDay := nearLimitUtxo.Value(int64(genesisTimestamp+oneDay), halfLife, incomeBase, incomeLimit)

	// Assert
	var expectedEmptyValueAfter30Days uint64 = 1827172467
	test.Assert(t, actualEmptyValueAfter30Days == expectedEmptyValueAfter30Days, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedEmptyValueAfter30Days, actualEmptyValueAfter30Days))
	var expectedHalfLimitValueAfter30Days uint64 = 1017338511889
	test.Assert(t, actualHalfLimitValueAfter30Days == expectedHalfLimitValueAfter30Days, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedHalfLimitValueAfter30Days, actualHalfLimitValueAfter30Days))
	var expectedNearLimitValueAfterOneDay uint64 = 1999900156698
	test.Assert(t, actualNearLimitValueAfterOneDay == expectedNearLimitValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedNearLimitValueAfterOneDay, actualNearLimitValueAfterOneDay))
	test.Assert(t, actualNearLimitValueAfterOneDay < incomeLimit, fmt.Sprintf("Wrong value. Expected a value < %d - Actual: %d", incomeLimit, actualNearLimitValueAfterOneDay))
}

func Test_Value_ValueIsBelowLimitAndIsYielding_IsConsistentOverTime(t *testing.T) {
	// Arrange
	var value uint64 = 10000 * particlesCount
	firstPeriod := int64(halfLife)
	secondPeriod := int64(30 * oneDay)
	utxo := NewUtxo(&InputInfo{}, NewOutput("", true, value), 0)
	valueAfterFirstPeriod := utxo.Value(firstPeriod, halfLife, base, limit)
	intermediateUtxo := NewUtxo(&InputInfo{}, NewOutput("", true, valueAfterFirstPeriod), firstPeriod)

	// Act
	valueAfterBothPeriodsInOneStep := utxo.Value(firstPeriod+secondPeriod, halfLife, base, limit)
	valueAfterBothPeriodsInTwoSteps := intermediateUtxo.Value(firstPeriod+secondPeriod, halfLife, base, limit)

	// Assert
	difference := int64(valueAfterBothPeriodsInOneStep) - int64(valueAfterBothPeriodsInTwoSteps)
	test.Assert(t, difference >= -1 && difference <= 1, fmt.Sprintf("Values differ by more than one smallest unit. One step: %d - Two steps: %d", valueAfterBothPeriodsInOneStep, valueAfterBothPeriodsInTwoSteps))
}

func Test_Value_IncomeBaseIs0AndIsYielding_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	value := limit / 2
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, 0, limit)

	// Assert
	expectedValueAfterHalfLife := value
	isFloatingPointErrorAtMostOneSmallestUnit := actualValueAfterHalfLife <= expectedValueAfterHalfLife && expectedValueAfterHalfLife-actualValueAfterHalfLife <= 1
	test.Assert(t, isFloatingPointErrorAtMostOneSmallestUnit, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_IncomeBaseIsGreaterThanLimitAndIsYielding_ReturnsLimit(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterOneDay := utxo.Value(int64(genesisTimestamp+oneDay), halfLife, 2*limit, limit)

	// Assert
	expectedValueAfterOneDay := limit
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
}

func Test_Value_SamePortionOfHalfLifeElapsedAndIsYielding_ReturnsSameValue(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := NewOutput("", true, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)
	portion := 1.5

	// Act
	value1 := utxo.Value(int64(genesisTimestamp+portion*halfLife), halfLife, base, limit)
	value2 := utxo.Value(int64(genesisTimestamp+portion*oneDay), oneDay, base, limit)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}

func Test_Value_SameElapsedTimeAndIsYielding_ReturnsSameValue(t *testing.T) {
	// Arrange
	var value uint64 = 0
	var outputTimestamp1 int64 = 0
	var outputTimestamp2 int64 = 10
	output := NewOutput("", true, value)
	utxo1 := NewUtxo(&InputInfo{}, output, outputTimestamp1)
	utxo2 := NewUtxo(&InputInfo{}, output, outputTimestamp2)
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
func Test_Value_ValueIsMaxUint64AndIsNotYielding_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = math.MaxUint64 // 18446744073709551615
	output := NewOutput("", false, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := uint64(math.Round(math.MaxUint64 / 2)) // 9223372036854775808
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsTwiceTheLimitAndIsNotYielding_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	value := 2 * limit
	output := NewOutput("", false, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := value / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIsLimitAndIsNotYielding_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	value := limit
	output := NewOutput("", false, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterHalfLife := utxo.Value(int64(genesisTimestamp+halfLife), halfLife, base, limit)

	// Assert
	expectedValueAfterHalfLife := limit / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndIsNotYielding_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 1
	output := NewOutput("", false, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

	// Act
	actualValueAfterNoTimeElapsed := utxo.Value(int64(genesisTimestamp), halfLife, base, limit)
	actualValueAfter1Minute := utxo.Value(genesisTimestamp+oneMinute, halfLife, base, limit)

	// Assert
	expectedValueAfterNoTimeElapsed := value
	test.Assert(t, actualValueAfterNoTimeElapsed == expectedValueAfterNoTimeElapsed, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterNoTimeElapsed, actualValueAfterNoTimeElapsed))
	var expectedValueAfter1Minute uint64 = 0
	test.Assert(t, actualValueAfter1Minute == expectedValueAfter1Minute, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfter1Minute, actualValueAfter1Minute))
}

func Test_Value_ValueIs0AndIsNotYielding_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 0
	output := NewOutput("", false, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)

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
	output := NewOutput("", false, value)
	utxo := NewUtxo(&InputInfo{}, output, 0)
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
	output := NewOutput("", false, value)
	utxo1 := NewUtxo(&InputInfo{}, output, outputTimestamp1)
	utxo2 := NewUtxo(&InputInfo{}, output, outputTimestamp2)
	var elapsedTimestamp int64 = oneDay
	utxo1Timestamp := outputTimestamp1
	utxo2Timestamp := outputTimestamp2

	// Act
	value1 := utxo1.Value(utxo1Timestamp+elapsedTimestamp, halfLife, base, limit)
	value2 := utxo2.Value(utxo2Timestamp+elapsedTimestamp, halfLife, base, limit)

	// Assert
	test.Assert(t, value1 == value2, "Values are not equals whereas it should be")
}
