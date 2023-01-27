package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/my-cloud/ruthenium/src/config"
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

const (
	synchronizationIntervalInSeconds = 10
	validationIntervalInSeconds      = 60
	verificationsCountPerValidation  = 6
	defaultPort                      = 8106
)

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (required if the private key is not provided)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)")
	port := flag.Uint64("port", environment.NewVariable("PORT").GetUint64Value(defaultPort), "The TCP port number of the host node")
	configurationPath := flag.String("configuration-path", environment.NewVariable("CONFIGURATION_PATH").GetStringValue("config"), "The configuration files path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	settings, err := config.NewSettings(*configurationPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to instantiate settings: %w", err).Error())
	}
	wallet, err := encryption.DecodeWallet(*mnemonic, *derivationPath, *password, *privateKey)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create wallet: %w", err).Error())
	}
	registry := poh.NewRegistry()
	validationTimer := validationIntervalInSeconds * time.Second
	watch := tick.NewWatch()
	clientFactory := gp2p.NewClientFactory(net.NewIpFinder())
	synchronizer, err := p2p.NewSynchronizer(uint16(*port), watch, clientFactory, *configurationPath, logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create synchronizer: %w", err).Error())
	}
	synchronizationTimer := time.Second * synchronizationIntervalInSeconds
	synchronizationEngine := tick.NewEngine(synchronizer.Synchronize, watch, synchronizationTimer, 1, 0)
	now := watch.Now()
	initialTimestamp := now.Truncate(validationTimer).Add(validationTimer).UnixNano()
	genesisTransaction := validation.NewRewardTransaction(wallet.Address(), initialTimestamp, settings.GenesisAmount)
	blockchain := verification.NewBlockchain(genesisTransaction, registry, validationTimer, synchronizer, logger)
	pool := validation.NewTransactionsPool(blockchain, registry, wallet.Address(), validationTimer, watch, logger)
	validationEngine := tick.NewEngine(pool.Validate, watch, validationTimer, 1, 0)
	verificationEngine := tick.NewEngine(blockchain.Update, watch, validationTimer, verificationsCountPerValidation, 1)
	serverFactory := gp2p.NewServerFactory()
	server, err := serverFactory.CreateServer(int(*port))
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create server: %w", err).Error())
	}
	host := p2p.NewHost(server, synchronizer, blockchain, pool, synchronizationEngine, validationEngine, verificationEngine, watch, logger)
	err = host.Run()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to run host: %w", err).Error())
	}
}
