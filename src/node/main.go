package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/config"
	"github.com/my-cloud/ruthenium/src/file"
	"strconv"
	"time"

	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/clock/tick"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/net"
	"github.com/my-cloud/ruthenium/src/node/protocol/poh"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
)

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (required if the private key is not provided)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKeyString := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)")
	infuraKey := flag.String("infura-key", environment.NewVariable("INFURA_KEY").GetStringValue(""), "The infura key (required to check the proof of humanity)")
	ip := flag.String("ip", environment.NewVariable("IP").GetStringValue(""), "The node IP or DNS address (detected if not provided)")
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(10600), "The TCP port number of the host node")
	settingsPath := flag.String("settings-path", environment.NewVariable("SETTINGS_PATH").GetStringValue("config/settings.json"), "The settings file path")
	seedsPath := flag.String("seeds-path", environment.NewVariable("SEEDS_PATH").GetStringValue("config/seeds.json"), "The seeds file path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level (possible values: 'debug', 'info', 'warn', 'error', 'fatal')")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	address := decodeAddress(mnemonic, derivationPath, password, privateKeyString, logger)
	host := createHost(settingsPath, infuraKey, seedsPath, ip, port, address, logger)
	logger.Info(fmt.Sprintf("host node starting for address: %s", address))
	err := host.Run()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to run host: %w", err).Error())
	}
}

func decodeAddress(mnemonic *string, derivationPath *string, password *string, privateKeyString *string, logger *console.Logger) string {
	var privateKey *encryption.PrivateKey
	var err error
	if *mnemonic != "" {
		privateKey, err = encryption.NewPrivateKeyFromMnemonic(*mnemonic, *derivationPath, *password)
	} else if *privateKeyString != "" {
		privateKey, err = encryption.NewPrivateKeyFromHex(*privateKeyString)
	} else {
		logger.Fatal(fmt.Errorf("nor the mnemonic neither the private key have been provided").Error())
	}
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create private key: %w", err).Error())
	}
	publicKey := encryption.NewPublicKey(privateKey)
	return publicKey.Address()
}

func createHost(settingsPath *string, infuraKey *string, seedsPath *string, ip *string, port *int, address string, logger *console.Logger) *p2p.Host {
	parser := file.NewJsonParser()
	var settings *config.Settings
	err := parser.Parse(*settingsPath, &settings)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to parse settings: %w", err).Error())
	}
	registry := poh.NewRegistry(*infuraKey, logger)
	watch := tick.NewWatch()
	synchronizer := createSynchronizer(parser, *seedsPath, *ip, *port, settings.MaxOutboundsCount(), watch, logger)
	blockchain := verification.NewBlockchain(registry, settings, synchronizer, logger)
	transactionsPool := validation.NewTransactionsPool(blockchain, settings, synchronizer, address, settings.ValidationTimer(), logger)
	synchronizationEngine := tick.NewEngine(synchronizer.Synchronize, watch, settings.SynchronizationTimer(), 1, 0)
	validationEngine := tick.NewEngine(transactionsPool.Validate, watch, settings.ValidationTimer(), 1, 0)
	verificationEngine := tick.NewEngine(blockchain.Update, watch, settings.ValidationTimer(), settings.VerificationsCountPerValidation(), 1)
	serverFactory := gp2p.NewServerFactory()
	handler := gp2p.NewHandler(blockchain, synchronizer, transactionsPool, watch, logger)
	server, err := serverFactory.CreateServer(*port, handler)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create server: %w", err).Error())
	}
	return p2p.NewHost(server, synchronizationEngine, validationEngine, verificationEngine, logger)
}

func createSynchronizer(parser *file.JsonParser, seedsPath string, hostIp string, port int, maxOutboundsCount int, watch *tick.Watch, logger *console.Logger) *p2p.Synchronizer {
	var seedsStringTargets []string
	err := parser.Parse(seedsPath, &seedsStringTargets)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to parse seeds: %w", err).Error())
	}
	scoresBySeedTarget := map[string]int{}
	for _, seedStringTarget := range seedsStringTargets {
		scoresBySeedTarget[seedStringTarget] = 0
	}
	ipFinder := net.NewIpFinder(logger)
	if hostIp == "" {
		hostIp, err = ipFinder.FindHostPublicIp()
		if err != nil {
			logger.Fatal(fmt.Errorf("failed to find the public IP: %w", err).Error())
		}
	}
	clientFactory := gp2p.NewClientFactory(ipFinder)
	return p2p.NewSynchronizer(clientFactory, hostIp, strconv.Itoa(port), maxOutboundsCount, scoresBySeedTarget, watch)
}
