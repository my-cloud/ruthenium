package p2p

import (
	"github.com/my-cloud/ruthenium/src/p2p"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_CreateServer(t *testing.T) {
	// Arrange
	serverFactory := p2p.NewServerFactory()

	// Act
	server, _ := serverFactory.CreateServer(0)

	// Assert
	test.Assert(t, server != nil, "server is nil whereas it should not")
}
