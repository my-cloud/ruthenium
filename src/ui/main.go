package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/file"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/net"
	"github.com/my-cloud/ruthenium/src/ui/config"
	"github.com/my-cloud/ruthenium/src/ui/server/index"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction"
	"github.com/my-cloud/ruthenium/src/ui/server/transactions"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet/address"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet/amount"
	"net/http"
	"strconv"
)

func main() {
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(8080), "The TCP port number of the UI server")
	hostIp := flag.String("host-ip", environment.NewVariable("HOST_IP").GetStringValue("127.0.0.1"), "The node host IP or DNS address")
	hostPort := flag.Int("host-port", environment.NewVariable("HOST_PORT").GetIntValue(10600), "The TCP port number of the host node")
	templatesPath := flag.String("templates-path", environment.NewVariable("TEMPLATES_PATH").GetStringValue("templates"), "The UI templates path")
	settingsPath := flag.String("settings-path", environment.NewVariable("SETTINGS_PATH").GetStringValue("config/settings.json"), "The settings file path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level (possible values: 'debug', 'info', 'warn', 'error', 'fatal')")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	target := p2p.NewTarget(*hostIp, strconv.Itoa(*hostPort))
	ipFinder := net.NewIpFinder(logger)
	clientFactory := gp2p.NewClientFactory(ipFinder)
	host, err := p2p.NewNeighbor(target, clientFactory)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to find blockchain client: %w", err).Error())
	}
	parser := file.NewJsonParser()
	var settings config.Settings
	err = parser.Parse(*settingsPath, &settings)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to parse settings: %w", err).Error())
	}
	particlesCount := settings.ParticlesPerToken
	http.Handle("/", index.NewHandler(*templatesPath, logger))
	http.Handle("/transaction", transaction.NewHandler(host, logger))
	http.Handle("/transactions", transactions.NewHandler(host, logger))
	http.Handle("/wallet/address", address.NewHandler(logger))
	http.Handle("/wallet/amount", amount.NewHandler(host, particlesCount, logger))
	logger.Info("user interface server is running...")
	logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), nil).Error())
}
