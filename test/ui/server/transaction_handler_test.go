package server

import (
	"bytes"
	"encoding/json"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"github.com/my-cloud/ruthenium/test"
	"github.com/my-cloud/ruthenium/test/mock"
	"net/http/httptest"
	"testing"
)

func Test_ServeHTTP_ValidTransaction_NeighborMethodCalled(t *testing.T) {
	// Arrange
	neighborMock := new(mock.NeighborMock)
	neighborMock.AddTransactionFunc = func(network.TransactionRequest) error { return nil }
	logger := log.NewLogger(log.Fatal)
	handler := server.NewTransactionHandler(neighborMock, 1, logger)

	transactionRequest := newTransactionRequest(
		test.PrivateKey,
		test.Address,
		"A",
		test.PublicKey,
		"1",
	)
	b, _ := json.Marshal(transactionRequest)
	body := bytes.NewReader(b)

	// Act
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("POST", "/transaction", body))

	// Assert
	isNeighborMethodCalled := len(neighborMock.AddTransactionCalls()) == 1
	test.Assert(t, isNeighborMethodCalled, "Neighbor method is not called whereas it should be.")
}

func newTransactionRequest(senderPrivateKey, senderAddress, recipientAddress, senderPublicKey, value string) *server.TransactionRequest {
	return &server.TransactionRequest{
		SenderPrivateKey: &senderPrivateKey,
		SenderAddress:    &senderAddress,
		RecipientAddress: &recipientAddress,
		SenderPublicKey:  &senderPublicKey,
		Value:            &value,
	}
}
