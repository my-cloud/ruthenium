package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/my-cloud/ruthenium/validatornode/application/p2p"
	"github.com/my-cloud/ruthenium/validatornode/application/p2p/gp2p"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol/validation"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol/verification"
	"github.com/my-cloud/ruthenium/validatornode/application/protocol/verification/poh"
	"github.com/my-cloud/ruthenium/validatornode/domain/clock"
	"github.com/my-cloud/ruthenium/validatornode/domain/encryption"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/config"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/environment"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/file"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/net"
)

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (required if the private key is not provided)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKeyString := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)")
	infuraKey := flag.String("infura-key", environment.NewVariable("INFURA_KEY").GetStringValue(""), "The infura key (required to check the proof of humanity)")
	ip := flag.String("ip", environment.NewVariable("IP").GetStringValue(""), "The validator node IP or DNS address (detected if not provided)")
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(10600), "The TCP port number of the validator node")
	settingsPath := flag.String("settings-path", environment.NewVariable("SETTINGS_PATH").GetStringValue("config/settings.json"), "The settings file path")
	seedsPath := flag.String("seeds-path", environment.NewVariable("SEEDS_PATH").GetStringValue("config/seeds.json"), "The seeds file path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level (possible values: 'debug', 'info', 'warn', 'error', 'fatal')")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	address, err := decodeAddress(*mnemonic, *derivationPath, *password, *privateKeyString)
	if err != nil {
		logger.Fatal(err.Error())
	}
	node, err := createNode(*settingsPath, *infuraKey, *seedsPath, *ip, *port, address, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Info(fmt.Sprintf("validator node running for address: %s", address))
	if err = node.Run(); err != nil {
		logger.Fatal(fmt.Errorf("failed to run node: %w", err).Error())
	}
}

func decodeAddress(mnemonic string, derivationPath string, password string, privateKeyString string) (string, error) {
	var privateKey *encryption.PrivateKey
	var err error
	if mnemonic != "" {
		privateKey, err = encryption.NewPrivateKeyFromMnemonic(mnemonic, derivationPath, password)
	} else if privateKeyString != "" {
		privateKey, err = encryption.NewPrivateKeyFromHex(privateKeyString)
	} else {
		return "", fmt.Errorf("nor the mnemonic neither the private key have been provided")
	}
	if err != nil {
		return "", fmt.Errorf("failed to create private key: %w", err)
	}
	publicKey := encryption.NewPublicKey(privateKey)
	return publicKey.Address(), nil
}

func createNode(settingsPath string, infuraKey string, seedsPath string, ip string, port int, address string, logger *console.Logger) (*p2p.Node, error) {
	settings, err := config.NewSettings(settingsPath)
	if err != nil {
		return nil, err
	}
	registry := poh.NewRegistry(infuraKey, logger)
	watch := clock.NewWatch()
	neighborhood, err := createNeighborhood(seedsPath, ip, port, settings, watch, logger)
	if err != nil {
		return nil, err
	}
	utxosPool := verification.NewUtxosPool()
	blockchain := verification.NewBlockchain(registry, settings, neighborhood, utxosPool, logger)
	transactionsPool := validation.NewTransactionsPool(blockchain, settings, neighborhood, utxosPool, address, logger)
	neighborhoodSynchronizationEngine := clock.NewEngine(neighborhood.Synchronize, watch, settings.SynchronizationTimer(), 1, 0)
	validationEngine := clock.NewEngine(transactionsPool.Validate, watch, settings.ValidationTimer(), 1, 0)
	verificationEngine := clock.NewEngine(blockchain.Update, watch, settings.ValidationTimer(), settings.VerificationsCountPerValidation(), 1)
	registrySynchronizationEngine := clock.NewEngine(registry.Synchronize, watch, 600*settings.ValidationTimer() /* TODO extract*/, 1, 0)
	handler := gp2p.NewHandler(blockchain, settings.Bytes(), neighborhood, transactionsPool, utxosPool, watch, logger)
	host, err := gp2p.NewHost(port, handler, settings.ValidationTimeout())
	if err != nil {
		return nil, err
	}
	return p2p.NewNode(host, logger, neighborhoodSynchronizationEngine, validationEngine, verificationEngine, registrySynchronizationEngine), nil
}

func createNeighborhood(seedsPath string, hostIp string, port int, settings *config.Settings, watch *clock.Watch, logger *console.Logger) (*p2p.Neighborhood, error) {
	var seedsStringTargets []string
	parser := file.NewJsonParser()
	err := parser.Parse(seedsPath, &seedsStringTargets)
	if err != nil {
		return nil, fmt.Errorf("unable to parse seeds: %w", err)
	}
	scoresBySeedTargetValue := map[string]int{}
	for _, seedStringTargetValue := range seedsStringTargets {
		scoresBySeedTargetValue[seedStringTargetValue] = 0
	}
	ipFinder := net.NewIpFinderImplementation(logger)
	if hostIp == "" {
		hostIp, err = ipFinder.FindHostPublicIp()
		if err != nil {
			return nil, fmt.Errorf("failed to find the public IP: %w", err)
		}
	}
	clientFactory := gp2p.NewSenderFactory(ipFinder, settings.ValidationTimeout())
	return p2p.NewNeighborhood(clientFactory, hostIp, strconv.Itoa(port), settings.MaxOutboundsCount(), scoresBySeedTargetValue, watch), nil
}
