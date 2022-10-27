package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/clock"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/src/p2p"
	"github.com/my-cloud/ruthenium/src/poh"
	"time"
)

const validationIntervalInSeconds = 60

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (optional)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (will be generated if not provided)")
	port := flag.Uint64("port", environment.NewVariable("PORT").GetUint64Value(network.DefaultPort), "The TCP port number for the protocol host node")
	configurationPath := flag.String("configuration-path", environment.NewVariable("CONFIGURATION_PATH").GetStringValue("config"), "The configuration files path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level")

	flag.Parse()
	logger := log.NewLogger(log.ParseLevel(*logLevel))
	wallet, err := encryption.DecodeWallet(*mnemonic, *derivationPath, *password, *privateKey)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create wallet: %w", err).Error())
	}
	registry := poh.NewRegistry()
	validationTimer := validationIntervalInSeconds * time.Second
	watch := clock.NewWatch()
	peering := p2p.NewSenderFactory()
	synchronizer, err := network.NewSynchronizer(uint16(*port), watch, peering, *configurationPath, logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create synchronizer: %w", err).Error())
	}
	blockchain := protocol.NewBlockchain(registry, validationTimer, watch, synchronizer, logger)
	pool := protocol.NewTransactionsPool(registry, watch, logger)
	validation := protocol.NewValidation(wallet.Address(), blockchain, pool, watch, validationTimer, logger)
	server, err := p2p.NewServerFactory().CreateServer(int(*port))
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create server: %w", err).Error())
	}
	host := network.NewHost(server, blockchain, pool, validation, synchronizer, watch, logger)
	err = host.Run()
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to run host: %w", err).Error())
	}
}
