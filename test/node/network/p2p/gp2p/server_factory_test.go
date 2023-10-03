package gp2p

import (
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_CreateServer(t *testing.T) {
	// Arrange
	serverFactory := gp2p.NewServerFactory()

	// Act
	server, _ := serverFactory.CreateServer(0, nil)

	// Assert
	test.Assert(t, server != nil, "server is nil whereas it should not")
}
