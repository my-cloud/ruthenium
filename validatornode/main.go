package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/validatornode/application/validation"
	"github.com/my-cloud/ruthenium/validatornode/application/verification"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/net"
	"github.com/my-cloud/ruthenium/validatornode/presentation/api"
	"io"
	"net/http"

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
	logger := console.NewLogger(settings.Log().Level())
	node, err := createHostNode(settings, logger)
	if err != nil {
		logger.Fatal(err.Error())
	} else if err = node.Run(); err != nil {
		logger.Fatal(fmt.Errorf("failed to run host validator node: %w", err).Error())
	}
}

func createHostNode(settings *configuration.Settings, logger *console.Logger) (*presentation.Node, error) {
	humanityRegistry := poh.NewHumanityRegistry(settings.Validator().InfuraKey(), logger)
	addressesRegistry := verification.NewAddressesRegistry(humanityRegistry, logger)
	watch := clock.NewWatch()
	neighborhood, err := createNeighborhood(settings, watch, logger)
	if err != nil {
		return nil, err
	}
	utxosRegistry := verification.NewUtxosRegistry(settings.Protocol())
	blockchain := verification.NewBlockchain(addressesRegistry, settings.Protocol(), neighborhood, utxosRegistry, logger)
	transactionsPool := validation.NewTransactionsPool(blockchain, settings.Protocol(), neighborhood, utxosRegistry, settings.Validator().Address(), logger)
	neighborhoodSynchronizationEngine := clock.NewEngine(neighborhood.Synchronize, watch, settings.Network().SynchronizationTimer(), 1, 0)
	validationEngine := clock.NewEngine(transactionsPool.Validate, watch, settings.Protocol().ValidationTimer(), 1, 0)
	verificationEngine := clock.NewEngine(blockchain.Update, watch, settings.Protocol().ValidationTimer(), settings.Protocol().VerificationsCountPerValidation(), 1)
	registrySynchronizationEngine := clock.NewEngine(addressesRegistry.Synchronize, watch, settings.Registry().SynchronizationTimer(), 1, 0)
	host, err := api.NewHost(blockchain, neighborhood, transactionsPool, utxosRegistry, settings.Host().Port(), settings.ProtocolBytes(), settings.Protocol().ValidationTimeout())
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("host validator node running for address: %s", settings.Validator().Address()))
	return presentation.NewNode(host, neighborhoodSynchronizationEngine, validationEngine, verificationEngine, registrySynchronizationEngine), nil
}

func createNeighborhood(settings *configuration.Settings, watch *clock.Watch, logger *console.Logger) (*network.Neighborhood, error) {
	seedsStringTargets := settings.Network().Seeds()
	scoresBySeedTargetValue := map[string]int{}
	for _, seedStringTargetValue := range seedsStringTargets {
		scoresBySeedTargetValue[seedStringTargetValue] = 0
	}
	ipFinder := net.NewIpFinderImplementation(logger)
	neighborFactory := p2p.NewNeighborFactory(ipFinder, settings.Network().ConnectionTimeout(), console.NewFatalLogger())
	hostIp, err := findHostPublicIp(settings.Host().Ip(), logger)
	if err != nil {
		return nil, err
	}
	return network.NewNeighborhood(neighborFactory, hostIp, settings.Host().Port(), settings.Network().MaxOutboundsCount(), scoresBySeedTargetValue, watch), nil
}

func findHostPublicIp(ip string, logger *console.Logger) (string, error) {
	if ip != "" {
		return ip, nil
	}
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", fmt.Errorf("failed to find the public IP: %w", err)
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
