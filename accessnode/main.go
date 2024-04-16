package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strconv"
	"time"

	"github.com/my-cloud/ruthenium/accessnode/presentation"
	"github.com/my-cloud/ruthenium/validatornode/domain/clock"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/config"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/environment"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/p2p"
)

func main() {
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(8080), "The TCP port number of the access node")
	validatorIp := flag.String("validator-ip", environment.NewVariable("VALIDATOR_IP").GetStringValue("127.0.0.1"), "The validator node IP or DNS address")
	validatorPort := flag.Int("validator-port", environment.NewVariable("VALIDATOR_PORT").GetIntValue(10600), "The TCP port number of the validator node")
	templatePath := flag.String("template-path", environment.NewVariable("TEMPLATE_PATH").GetStringValue("accessnode/template.html"), "The UI template path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level (possible values: 'debug', 'info', 'warn', 'error', 'fatal')")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	// ipFinder := net.NewIpFinderImplementation(logger)
	// neighborFactory := p2p.NewNeighborFactory(ipFinder, time.Minute)
	// validatorNeighbor, err := neighborFactory.CreateSender(*validatorIp, strconv.Itoa(*validatorPort))
	validatorNeighbor, err := p2p.NewNeighbor(*validatorIp, strconv.Itoa(*validatorPort), time.Minute, console.NewLogger(console.Fatal))
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to find blockchain client: %w", err).Error())
	}
	settingsBytes, err := validatorNeighbor.GetSettings()
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to get settings: %w", err).Error())
	}
	var settings *config.Settings
	err = json.Unmarshal(settingsBytes, &settings)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to unmarshal settings: %w", err).Error())
	}
	watch := clock.NewWatch()
	node := presentation.NewNode(strconv.Itoa(*port), validatorNeighbor, settings, *templatePath, watch, logger)
	logger.Info("host access node is running...")
	logger.Fatal(node.Run().Error())
}
