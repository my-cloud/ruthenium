package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/config"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/net"
	"github.com/my-cloud/ruthenium/src/ui/server/index"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction"
	"github.com/my-cloud/ruthenium/src/ui/server/transactions"
	"github.com/my-cloud/ruthenium/src/ui/server/validation/start"
	"github.com/my-cloud/ruthenium/src/ui/server/validation/stop"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet/address"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet/amount"
	"net/http"
	"strconv"
)

const (
	defaultPort     = 8080
	defaultHostPort = 8106
	transactionFee  = 1000
)

func main() {
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(defaultPort), "The TCP port number of the UI server")
	hostIp := flag.String("host-ip", environment.NewVariable("HOST_IP").GetStringValue(""), "The node host IP address")
	hostPort := flag.Int("host-port", environment.NewVariable("HOST_PORT").GetIntValue(defaultHostPort), "The TCP port number of the host node")
	templatesPath := flag.String("templates-path", environment.NewVariable("TEMPLATES_PATH").GetStringValue("templates"), "The UI templates path")
	configurationPath := flag.String("configuration-path", environment.NewVariable("CONFIGURATION_PATH").GetStringValue("config"), "The configuration files path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	target := p2p.NewTarget(*hostIp, strconv.Itoa(*hostPort))
	ipFinder := net.NewIpFinder(logger)
	clientFactory := gp2p.NewClientFactory(ipFinder)
	host, err := p2p.NewNeighbor(target, clientFactory)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to find blockchain client: %w", err).Error())
	}
	settings, err := config.NewSettings(*configurationPath)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to instantiate settings: %w", err).Error())
	}
	particlesCount := settings.ParticlesCount
	http.Handle("/", index.NewHandler(*templatesPath, logger))
	http.Handle("/transaction", transaction.NewHandler(host, particlesCount, transactionFee, logger))
	http.Handle("/transactions", transactions.NewHandler(host, logger))
	http.Handle("/validation/start", start.NewHandler(host, logger))
	http.Handle("/validation/stop", stop.NewHandler(host, logger))
	http.Handle("/wallet/address", address.NewHandler(logger))
	http.Handle("/wallet/amount", amount.NewHandler(host, particlesCount, logger))
	logger.Info("user interface server is running...")
	logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), nil).Error())
}
