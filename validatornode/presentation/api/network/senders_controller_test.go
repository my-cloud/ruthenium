package network

import (
	"context"
	"encoding/json"
	"github.com/my-cloud/ruthenium/validatornode/application"
	"sync"
	"testing"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_HandleTargetsRequest_AddInvalidTargets_AddTargetsNotCalled(t *testing.T) {
	// Arrange
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.AddTargetsFunc = func([]string) {}
	controller := NewSendersController(sendersManagerMock)
	targets := []string{"target"}
	marshalledTargets, _ := json.Marshal(targets)
	req := gp2p.Data{Bytes: marshalledTargets}

	// Act
	_, _ = controller.HandleTargetsRequest(context.TODO(), req)

	// Assert
	isMethodCalled := len(sendersManagerMock.AddTargetsCalls()) != 0
	test.Assert(t, !isMethodCalled, "Method is called whereas it should not.")
}

func Test_HandleTargetsRequest_AddValidTargets_AddTargetsCalled(t *testing.T) {
	// Arrange
	waitGroup := sync.WaitGroup{}
	sendersManagerMock := new(application.SendersManagerMock)
	sendersManagerMock.AddTargetsFunc = func([]string) { waitGroup.Done() }
	controller := NewSendersController(sendersManagerMock)
	targets := []string{"target"}
	marshalledTargets, _ := json.Marshal(targets)
	req := gp2p.Data{Bytes: marshalledTargets}
	waitGroup.Add(1)

	// Act
	_, _ = controller.HandleTargetsRequest(context.TODO(), req)

	// Assert
	waitGroup.Wait()
	isMethodCalled := len(sendersManagerMock.AddTargetsCalls()) == 1
	test.Assert(t, isMethodCalled, "Method is not called whereas it should be.")
}
