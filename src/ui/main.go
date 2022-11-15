package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/config"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/p2p"
	"github.com/my-cloud/ruthenium/src/p2p/net"
	"github.com/my-cloud/ruthenium/src/ui/server"
	"net/http"
	"strconv"
)

const defaultPort = 8080

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (required if the private key is not provided)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)")
	port := flag.Uint64("port", environment.NewVariable("PORT").GetUint64Value(defaultPort), "The TCP port number of the UI server")
	hostIp := flag.String("host-ip", environment.NewVariable("HOST_IP").GetStringValue(""), "The node host IP address")
	hostPort := flag.Uint64("host-port", environment.NewVariable("HOST_PORT").GetUint64Value(network.DefaultPort), "The TCP port number of the host node")
	templatesPath := flag.String("templates-path", environment.NewVariable("TEMPLATES_PATH").GetStringValue("src/ui/templates"), "The UI templates path")
	configurationPath := flag.String("configuration-path", environment.NewVariable("CONFIGURATION_PATH").GetStringValue("config"), "The configuration files path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level")

	flag.Parse()
	logger := log.NewLogger(log.ParseLevel(*logLevel))
	target := network.NewTarget(*hostIp, uint16(*hostPort))
	clientFactory := p2p.NewClientFactory(net.NewIpFinder())
	host, err := network.NewNeighbor(target, clientFactory, logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to find blockchain client: %w", err).Error())
	}
	settings, err := config.NewSettings(*configurationPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to instantiate settings: %w", err).Error())
	}
	particlesCount := settings.ParticlesCount
	http.Handle("/", server.NewIndexHandler(*templatesPath, logger))
	http.Handle("/wallet", server.NewWalletHandler(*mnemonic, *derivationPath, *password, *privateKey, logger))
	http.Handle("/transaction", server.NewTransactionHandler(host, particlesCount, logger))
	http.Handle("/transactions", server.NewTransactionsHandler(host, logger))
	http.Handle("/wallet/amount", server.NewWalletAmountHandler(host, particlesCount, logger))
	http.Handle("/mine", server.NewValidationHandler(host, logger))
	http.Handle("/mine/start", server.NewValidationStartHandler(host, logger))
	http.Handle("/mine/stop", server.NewValidationStopHandler(host, logger))
	logger.Info("user interface server is running...")
	logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(int(*port)), nil).Error())
}
