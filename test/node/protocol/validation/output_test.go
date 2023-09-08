package validation

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/test"
	"testing"
	"time"
)

const (
	limit               uint64 = 10000000000000
	genesisTimestamp           = 0
	halfLife                   = 373.59 * 24 * 60 * 60 * 1e9
	k                          = 9.790310290581342
	validationTimestamp        = 60 * 1e9
)

//////////////////////////////////// INCOME ////////////////////////////////////
func Test_Value_ValueIsLimitAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     limit,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	expectedValueAfterOneDay := limit
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	expectedValueAfterHalfLife := limit
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     1,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 360342
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 50000445605
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs50kAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     5000000000000,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 5001577710768
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 5578244486849
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs200kAndHasIncome_ReturnsValueWithIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: true,
		Value:     20000000000000,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 19981463514647
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 15000000000000
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

//////////////////////////////////// NO INCOME ////////////////////////////////////
func Test_Value_ValueIsLimitAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     limit,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 9981463514647
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	expectedValueAfterHalfLife := limit / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs1AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     1,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 0
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	var expectedValueAfterHalfLife uint64 = 0
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs50kAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 5000000000000
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     value,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 4990731757323
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	expectedValueAfterHalfLife := value / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

func Test_Value_ValueIs200kAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
	// Arrange
	var value uint64 = 20000000000000
	utxo := &network.UtxoResponse{
		HasIncome: false,
		Value:     20000000000000,
	}
	output := validation.NewOutputFromUtxoResponse(utxo, genesisTimestamp, halfLife, limit, k, validationTimestamp)
	oneDay := 24 * float64(time.Hour.Nanoseconds())

	// Act
	actualValueAfterOneDay := output.Value(int64(oneDay))
	actualValueAfterHalfLife := output.Value(int64(halfLife))

	// Assert
	var expectedValueAfterOneDay uint64 = 19962927029295
	test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
	expectedValueAfterHalfLife := value / 2
	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
}

//
//// FIXME result is approximated
//func Test_Value_ValueIsMaxUint64AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
//	// Arrange
//	var value uint64 = 18000000000000000000 // Should be 18446744073709551615
//	utxo := &network.UtxoResponse{
//		HasIncome:     false,
//		Value:         value,
//	}
//	halfLife := 373.59 * 24 * float64(time.Hour.Nanoseconds())
//	validationTimestamp := 60 * time.Second.Nanoseconds()
//	output := validation.NewOutputFromUtxoResponse(utxo, halfLife, validationTimestamp, 0)
//	//oneDay := 24 * float64(time.Hour.Nanoseconds())
//
//	// Act
//	//actualValueAfterOneDay := output.Value(int64(oneDay))
//	actualValueAfterHalfLife := output.Value(int64(halfLife+5))
//
//	// Assert
//	//var expectedValueAfterOneDay uint64 = 19962927029295
//	//test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
//	expectedValueAfterHalfLife := value/2
//	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
//}
//
//func Test_Value_ValueIsMaxFloat64AndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
//	// Arrange
//	var value uint64 = 179769313486231583 // Should be 1,797,693,134.86231584
//	utxo := &network.UtxoResponse{
//		HasIncome:     false,
//		Value:         value,
//	}
//	halfLife := 373.59 * 24 * float64(time.Hour.Nanoseconds())
//	validationTimestamp := 60 * time.Second.Nanoseconds()
//	output := validation.NewOutputFromUtxoResponse(utxo, halfLife, validationTimestamp, 0)
//
//	//oneDay := 24 * float64(time.Hour.Nanoseconds())
//
//	// Act
//	//actualValueAfterOneDay := output.Value(int64(oneDay))
//	actualValueAfterHalfLife := output.Value(int64(halfLife))
//
//	// Assert
//	//var expectedValueAfterOneDay uint64 = 19962927029295
//	//test.Assert(t, actualValueAfterOneDay == expectedValueAfterOneDay, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterOneDay, actualValueAfterOneDay))
//	expectedValueAfterHalfLife := value/2
//	test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife))
//}
//
//func Test_Value_ValueIsOddAndHasNoIncome_ReturnsValueWithoutIncome(t *testing.T) {
//	// Arrange
//	//fmt.Println("MaxInt64: ", math.MaxInt64)
//	//fmt.Println("MaxInt32: ", math.MaxInt32)
//	//fmt.Println("MaxInt16: ", math.MaxInt16)
//	//fmt.Println("MaxUint32: ", math.MaxUint32)
//	//fmt.Println("MaxUint16: ", math.MaxUint16)
//	//var value uint64 = 3628651111137335
//	halfLife := 373.59 * 24 * float64(time.Hour.Nanoseconds())
//	validationTimestamp := 60 * time.Second.Nanoseconds()
//
//	for i:=10000000000000000; i < math.MaxInt64; i=i+1000 {
//		if i == 1000000 ||i == 10000000 ||i == 100000000 ||i == 1000000000 ||i == 10000000000 ||i == 100000000000 ||i == 1000000000000 ||i == 10000000000000 ||i == 100000000000000 ||i == 1000000000000000 ||i == 10000000000000000 ||i == 100000000000000000{
//			fmt.Println("loop index: ", i)
//		}
//		var value uint64 = uint64(i+3)
//		utxo := &network.UtxoResponse{
//			HasIncome: false,
//			Value:     value,
//		}
//		output := validation.NewOutputFromUtxoResponse(utxo, halfLife, validationTimestamp, 0)
//
//		// Act
//		actualValueAfterHalfLife := output.Value(int64(halfLife))
//
//		// Assert
//		expectedValueAfterHalfLife := value / 2
//		test.Assert(t, actualValueAfterHalfLife == expectedValueAfterHalfLife, fmt.Sprintf("Wrong value. Expected: %d - Actual: %d - Value: %d", expectedValueAfterHalfLife, actualValueAfterHalfLife, value))
//	}
//}
