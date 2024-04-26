package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/accessnode/infrastructure/configuration"
	"time"

	"github.com/my-cloud/ruthenium/accessnode/presentation"
	"github.com/my-cloud/ruthenium/validatornode/domain/clock"
	validatorconfiguration "github.com/my-cloud/ruthenium/validatornode/infrastructure/configuration"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/environment"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/p2p"
)

func main() {
	settingsPath := flag.String("settings-path", environment.NewVariable("SETTINGS_PATH").GetStringValue("accessnode/settings.json"), "The settings file path")
	flag.Parse()
	settings, err := configuration.NewSettings(*settingsPath)
	if err != nil {
		panic(err.Error())
	}
	logger := console.NewLogger(settings.Log().Level())
	validatorNeighbor, err := p2p.NewNeighbor(settings.Validator().Ip(), settings.Validator().Port(), time.Minute, console.NewFatalLogger())
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to find blockchain client: %w", err).Error())
	}
	settingsBytes, err := validatorNeighbor.GetSettings()
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to get protocol settings: %w", err).Error())
	}
	var protocolSettings *validatorconfiguration.ProtocolSettings
	err = json.Unmarshal(settingsBytes, &protocolSettings)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to unmarshal protocol settings: %w", err).Error())
	}
	watch := clock.NewWatch()
	node := presentation.NewNode(settings.Host().Port(), validatorNeighbor, protocolSettings, settings.Template().Path(), watch, logger)
	logger.Info("host access node is running...")
	logger.Fatal(node.Run().Error())
}
