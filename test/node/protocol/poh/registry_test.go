package poh

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/node/protocol/poh"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/log/logtest"
	"testing"
)

func Test_IsRegistered_InfuraKey_ReturnsTrue(t *testing.T) {
	// Arrange
	address := "0x0000000000000000000000000000000000000001"
	logger := logtest.NewLoggerMock()
	registry := poh.NewRegistry("", logger)

	// Act
	isRegistered, err := registry.IsRegistered(address)
	fmt.Println(err)

	// Assert
	test.Assert(t, isRegistered, "proof of humanity is invalid whereas it should be")
}
