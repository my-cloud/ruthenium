package poh

import (
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
	"testing"
)

func Test_NewHumanityRegistry_InfuraKeyNotProvided_WarnIsLogged(t *testing.T) {
	// Arrange
	logger := log.NewLoggerMock()

	// Act
	_ = NewHumanityRegistry("", logger)

	// Assert
	test.AssertThatMessageIsLogged(t, logger.WarnCalls(), "infura key not provided")
}

func Test_IsRegistered_InfuraKeyNotProvided_ReturnsTrue(t *testing.T) {
	// Arrange
	registry := NewHumanityRegistry("", log.NewLoggerMock())

	// Act
	isRegistered, err := registry.IsRegistered("0x0000000000000000000000000000000000000000")

	// Assert
	test.Assert(t, err == nil, fmt.Sprintf("an error has been returned whereas it should not: %v", err))
	test.Assert(t, isRegistered, "address is not registered whereas it should be")
}

func Test_Url_InfuraKeyProvided_ReturnsWebsocketUrl(t *testing.T) {
	// Arrange
	registry := NewHumanityRegistry("key", log.NewLoggerMock())

	// Act
	actualUrl := registry.url()

	// Assert
	expectedUrl := "wss://mainnet.infura.io/ws/v3/key"
	test.Assert(t, actualUrl == expectedUrl, fmt.Sprintf("url is %s whereas it should be %s", actualUrl, expectedUrl))
}
