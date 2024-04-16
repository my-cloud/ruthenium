package api

import (
	"context"
	"encoding/json"
	"testing"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/application/ledger"
	"github.com/my-cloud/ruthenium/validatornode/domain/protocol"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_HandleFirstBlockTimestampRequest_ValidRequest_FirstBlockTimestampCalled(t *testing.T) {
	// Arrange
	blocksManagerMock := new(ledger.BlocksManagerMock)
	blocksManagerMock.FirstBlockTimestampFunc = func() int64 { return 0 }
	controller := NewBlocksController(blocksManagerMock)
	req := gp2p.Data{}

	// Act
	_, _ = controller.HandleFirstBlockTimestampRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(blocksManagerMock.FirstBlockTimestampCalls()) != 0
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}

func Test_HandleBlocksRequest_ValidBlocksRequest_LastBlocksCalled(t *testing.T) {
	// Arrange
	blocksManagerMock := new(ledger.BlocksManagerMock)
	blocksManagerMock.BlocksFunc = func(uint64) []*protocol.Block { return nil }
	controller := NewBlocksController(blocksManagerMock)
	var height uint64 = 0
	marshalledHeight, _ := json.Marshal(&height)
	req := gp2p.Data{Bytes: marshalledHeight}

	// Act
	_, _ = controller.HandleBlocksRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(blocksManagerMock.BlocksCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}
