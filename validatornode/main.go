package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application/validation"
	"github.com/my-cloud/ruthenium/validatornode/application/verification"
	"github.com/my-cloud/ruthenium/validatornode/presentation/api"
	"io"
	"net/http"
	"time"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/net"

	"github.com/my-cloud/ruthenium/validatornode/application/network"
	"github.com/my-cloud/ruthenium/validatornode/domain/clock"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/configuration"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/environment"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/p2p"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/poh"
	"github.com/my-cloud/ruthenium/validatornode/presentation"
)

func main() {
	settingsPath := flag.String("settings-path", environment.NewVariable("SETTINGS_PATH").GetStringValue("validatornode/settings.json"), "The settings file path")

	flag.Parse()
	settings, err := configuration.NewSettings(*settingsPath)
	if err != nil {
		panic(err.Error())
	}
	logger := console.NewLogger(console.ParseLevel(settings.LogLevel()))
	node, err := createHostNode(settings, logger)
	if err != nil {
		logger.Fatal(err.Error())
	} else if err = node.Run(); err != nil {
		logger.Fatal(fmt.Errorf("failed to run host validator node: %w", err).Error())
	}
}

func createHostNode(settings *configuration.Settings, logger *console.Logger) (*presentation.Node, error) {
	humanityRegistry := poh.NewHumanityRegistry(settings.InfuraKey(), logger)
	addressesRegistry := verification.NewAddressesRegistry(humanityRegistry, logger)
	watch := clock.NewWatch()
	neighborhood, err := createNeighborhood(settings, watch, logger)
	if err != nil {
		return nil, err
	}
	utxosRegistry := verification.NewUtxosRegistry(settings)
	blockchain := verification.NewBlockchain(addressesRegistry, settings, neighborhood, utxosRegistry, logger)
	transactionsPool := validation.NewTransactionsPool(blockchain, settings, neighborhood, utxosRegistry, settings.Address(), logger)
	neighborhoodSynchronizationEngine := clock.NewEngine(neighborhood.Synchronize, watch, settings.SynchronizationTimer(), 1, 0)
	validationEngine := clock.NewEngine(transactionsPool.Validate, watch, settings.ValidationTimer(), 1, 0)
	verificationEngine := clock.NewEngine(blockchain.Update, watch, settings.ValidationTimer(), settings.VerificationsCountPerValidation(), 1)
	registrySynchronizationEngine := clock.NewEngine(addressesRegistry.Synchronize, watch, time.Hour /* TODO extract */, 1, 0)
	host, err := api.NewHost(settings, blockchain, neighborhood, transactionsPool, utxosRegistry)
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("host validator node running for address: %s", settings.Address()))
	return presentation.NewNode(host, neighborhoodSynchronizationEngine, validationEngine, verificationEngine, registrySynchronizationEngine), nil
}

func createNeighborhood(settings *configuration.Settings, watch *clock.Watch, logger *console.Logger) (*network.Neighborhood, error) {
	seedsStringTargets := settings.Seeds()
	scoresBySeedTargetValue := map[string]int{}
	for _, seedStringTargetValue := range seedsStringTargets {
		scoresBySeedTargetValue[seedStringTargetValue] = 0
	}
	ipFinder := net.NewIpFinderImplementation(logger)
	hostIp := settings.Ip()
	if hostIp == "" {
		var err error
		hostIp, err = findHostPublicIp(logger)
		if err != nil {
			return nil, fmt.Errorf("failed to find the public IP: %w", err)
		}
	}
	neighborFactory := p2p.NewNeighborFactory(ipFinder, settings.ValidationTimeout(), console.NewLogger(console.Fatal))
	return network.NewNeighborhood(neighborFactory, hostIp, settings.Port(), settings.MaxOutboundsCount(), scoresBySeedTargetValue, watch), nil
}

func findHostPublicIp(logger *console.Logger) (string, error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", err
	}
	defer func() {
		if bodyCloseError := resp.Body.Close(); bodyCloseError != nil {
			logger.Error(fmt.Errorf("failed to close public IP request body: %w", bodyCloseError).Error())
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
