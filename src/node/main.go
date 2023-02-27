package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/my-cloud/ruthenium/src/config"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/clock/tick"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/net"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/poh"
	"github.com/my-cloud/ruthenium/src/node/protocol/pohtest"
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"github.com/my-cloud/ruthenium/src/node/protocol/verification"
)

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (required if the private key is not provided)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)")
	ip := flag.String("ip", environment.NewVariable("IP").GetStringValue(""), "The node IP or DNS address (detected if not provided)")
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(8106), "The TCP port number of the host node")
	configurationPath := flag.String("configuration-path", environment.NewVariable("CONFIGURATION_PATH").GetStringValue("config"), "The configuration files path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level (possible values: 'debug', 'info', 'warn', 'error', 'fatal')")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	settings, err := config.NewSettings(*configurationPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to instantiate settings: %w", err).Error())
	}
	wallet, err := encryption.NewWallet(*mnemonic, *derivationPath, *password, *privateKey)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create wallet: %w", err).Error())
	}
	var registry protocol.Registry
	if settings.NetworkId == "mainnet" {
		registry = poh.NewRegistry()
	} else {
		registry = pohtest.NewRegistry()
	}
	validationTimer := time.Duration(settings.ValidationIntervalInSeconds) * time.Second
	watch := tick.NewWatch()
	ipFinder := net.NewIpFinder(logger)
	clientFactory := gp2p.NewClientFactory(ipFinder)
	hostIp := *ip
	if hostIp == "" {
		hostIp, err = ipFinder.FindHostPublicIp()
		if err != nil {
			logger.Fatal(fmt.Errorf("failed to find the public IP: %w", err).Error())
		}
	}
	scoresBySeedTarget, err := readSeedsTargets(*configurationPath, logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to read seeds targets: %w", err).Error())
	}
	synchronizer := p2p.NewSynchronizer(clientFactory, hostIp, strconv.Itoa(*port), settings.MaxOutboundsCount, scoresBySeedTarget, watch)
	synchronizationTimer := time.Second * time.Duration(settings.SynchronizationIntervalInSeconds)
	synchronizationEngine := tick.NewEngine(synchronizer.Synchronize, watch, synchronizationTimer, 1, 0)
	now := watch.Now()
	initialTimestamp := now.Truncate(validationTimer).Add(validationTimer).UnixNano()
	genesisTransaction := validation.NewRewardTransaction(wallet.Address(), initialTimestamp, settings.GenesisAmountInParticles)
	blockchain := verification.NewBlockchain(genesisTransaction, settings.MinimalTransactionFee, registry, validationTimer, synchronizer, logger)
	transactionsPool := validation.NewTransactionsPool(blockchain, settings.MinimalTransactionFee, registry, synchronizer, wallet.Address(), validationTimer, logger)
	validationEngine := tick.NewEngine(transactionsPool.Validate, watch, validationTimer, 1, 0)
	verificationEngine := tick.NewEngine(blockchain.Update, watch, validationTimer, settings.VerificationsCountPerValidation, 1)
	serverFactory := gp2p.NewServerFactory()
	server, err := serverFactory.CreateServer(*port)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create server: %w", err).Error())
	}
	handler := p2p.NewHandler(blockchain, synchronizer, transactionsPool, watch, logger)
	host := p2p.NewHost(handler, server, synchronizationEngine, validationEngine, verificationEngine, logger)
	logger.Info(fmt.Sprintf("host node starting for address: %s", wallet.Address()))
	err = host.Run()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to run host: %w", err).Error())
	}
}

func readSeedsTargets(configurationPath string, logger *console.Logger) (map[string]int, error) {
	jsonFile, err := os.Open(configurationPath + "/seeds.json")
	if err != nil {
		return nil, fmt.Errorf("unable to open seeds IPs configuration file: %w", err)
	}
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read seeds IPs configuration file: %w", err)
	}
	if err = jsonFile.Close(); err != nil {
		logger.Error(fmt.Errorf("unable to close seeds IPs configuration file: %w", err).Error())
	}
	var seedsStringTargets []string
	if err = json.Unmarshal(byteValue, &seedsStringTargets); err != nil {
		return nil, fmt.Errorf("unable to unmarshal seeds IPs: %w", err)
	}
	scoresBySeedTarget := map[string]int{}
	for _, seedStringTarget := range seedsStringTargets {
		scoresBySeedTarget[seedStringTarget] = 0
	}
	return scoresBySeedTarget, nil
}
