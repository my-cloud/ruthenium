package protocol

import (
	"context"
	"reflect"
	"testing"

	gp2p "github.com/leprosus/golang-p2p"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/test"
)

func Test_HandleSettingsRequest_ValidRequest_SettingsCalled(t *testing.T) {
	// Arrange
	expectedSettings := []byte{0}
	controller := NewSettingsController(expectedSettings)
	req := gp2p.Data{}

	// Act
	data, _ := controller.HandleSettingsRequest(context.TODO(), req)

	// Assert
	actualSettings := data.GetBytes()
	test.Assert(t, reflect.DeepEqual(expectedSettings, actualSettings), "Settings are not the expected ones.")
}
