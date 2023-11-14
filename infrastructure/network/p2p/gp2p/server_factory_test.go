package gp2p

import (
	"github.com/my-cloud/ruthenium/infrastructure/network/p2p"
	"github.com/my-cloud/ruthenium/infrastructure/test"
	"testing"
	"time"
)

func Test_CreateServer(t *testing.T) {
	// Arrange
	settings := new(p2p.SettingsMock)
	settings.ValidationTimeoutFunc = func() time.Duration { return 0 }
	serverFactory := NewServerFactory(nil, settings)

	// Act
	server, _ := serverFactory.CreateServer(0)

	// Assert
	test.Assert(t, server != nil, "server is nil whereas it should not")
}
