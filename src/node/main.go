package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/config"
	"github.com/my-cloud/ruthenium/src/encryption"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/node/protocol/poh"
	"github.com/my-cloud/ruthenium/src/p2p"
	"github.com/my-cloud/ruthenium/src/p2p/net"
	"time"
)

const (
	validationIntervalInSeconds     = 60
	verificationsCountPerValidation = 6
)

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (required if the private key is not provided)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)")
	port := flag.Uint64("port", environment.NewVariable("PORT").GetUint64Value(network.DefaultPort), "The TCP port number of the host node")
	configurationPath := flag.String("configuration-path", environment.NewVariable("CONFIGURATION_PATH").GetStringValue("config"), "The configuration files path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level")

	flag.Parse()
	logger := log.NewLogger(log.ParseLevel(*logLevel))
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
	watch := clock.NewWatch()
	clientFactory := p2p.NewClientFactory(net.NewIpFinder())
	synchronizer, err := network.NewSynchronizer(uint16(*port), watch, clientFactory, *configurationPath, logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create synchronizer: %w", err).Error())
	}
	blockchain := protocol.NewBlockchain(registry, validationTimer, synchronizer, logger)
	pool := protocol.NewTransactionsPool(blockchain, registry, wallet.Address(), settings.GenesisAmount, validationTimer, watch, logger)
	validation := protocol.NewEngine(pool.Validate, watch, validationTimer, 1, 0, logger)
	verification := protocol.NewEngine(blockchain.Verify, watch, validationTimer, verificationsCountPerValidation, 1, logger)
	serverFactory := p2p.NewServerFactory()
	server, err := serverFactory.CreateServer(int(*port))
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create server: %w", err).Error())
	}
	host := network.NewHost(server, blockchain, pool, validation, verification, synchronizer, watch, logger)
	err = host.Run()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to run host: %w", err).Error())
	}
}
