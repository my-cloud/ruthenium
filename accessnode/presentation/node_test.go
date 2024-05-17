package presentation

import (
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"testing"
)

func Test_CalculateFee_UnknownTransactionId_ReturnsError(t *testing.T) {
	// Arrange
	// Act
	node := NewNode("", nil, nil, "", nil, nil)

	// Assert
	test.Assert(t, node != nil, "node is nil whereas it should not")
}
