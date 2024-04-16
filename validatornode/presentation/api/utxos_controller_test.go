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

func Test_HandleUtxosRequest_ValidUtxosRequest_UtxosByAddressCalled(t *testing.T) {
	// Arrange
	utxosManagerMock := new(ledger.UtxosManagerMock)
	utxosManagerMock.UtxosFunc = func(string) []*protocol.Utxo { return nil }
	controller := NewUtxosController(utxosManagerMock)
	address := "address"
	marshalledAddress, _ := json.Marshal(&address)
	req := gp2p.Data{Bytes: marshalledAddress}

	// Act
	_, _ = controller.HandleUtxosRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(utxosManagerMock.UtxosCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}
