package gp2p

import (
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/node/network/p2p/p2ptest"
	"testing"
	"time"
)

func Test_CreateServer(t *testing.T) {
	// Arrange
	settings := new(p2ptest.SettingsMock)
	settings.ValidationTimeoutFunc = func() time.Duration { return 0 }
	serverFactory := gp2p.NewServerFactory(nil, settings)

	// Act
	server, _ := serverFactory.CreateServer(0)

	// Assert
	test.Assert(t, server != nil, "server is nil whereas it should not")
}
