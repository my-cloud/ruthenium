package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/config"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log/console"
	"github.com/my-cloud/ruthenium/src/node/clock/tick"
	"github.com/my-cloud/ruthenium/src/node/network/p2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/gp2p"
	"github.com/my-cloud/ruthenium/src/node/network/p2p/net"
	"github.com/my-cloud/ruthenium/src/ui/server/index"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction/info"
	"github.com/my-cloud/ruthenium/src/ui/server/transaction/output/progress"
	"github.com/my-cloud/ruthenium/src/ui/server/transactions"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet/address"
	"github.com/my-cloud/ruthenium/src/ui/server/wallet/amount"
	"net/http"
	"strconv"
	"time"
)

func main() {
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(8080), "The TCP port number of the UI server")
	hostIp := flag.String("host-ip", environment.NewVariable("HOST_IP").GetStringValue("127.0.0.1"), "The node host IP or DNS address")
	hostPort := flag.Int("host-port", environment.NewVariable("HOST_PORT").GetIntValue(10600), "The TCP port number of the host node")
	templatesPath := flag.String("templates-path", environment.NewVariable("TEMPLATES_PATH").GetStringValue("templates"), "The UI templates path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level (possible values: 'debug', 'info', 'warn', 'error', 'fatal')")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	target := p2p.NewTarget(*hostIp, strconv.Itoa(*hostPort))
	ipFinder := net.NewIpFinder(logger)
	clientFactory := gp2p.NewClientFactory(ipFinder, time.Minute)
	host, err := p2p.NewNeighbor(target, clientFactory)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to find blockchain client: %w", err).Error())
	}
	settingsBytes, err := host.GetSettings()
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to get settings: %w", err).Error())
	}
	var settings *config.Settings
	err = json.Unmarshal(settingsBytes, &settings)
	if err != nil {
		logger.Fatal(fmt.Errorf("unable to unmarshal settings: %w", err).Error())
	}
	watch := tick.NewWatch()
	http.Handle("/", index.NewHandler(*templatesPath, logger))
	http.Handle("/transaction", transaction.NewHandler(host, logger))
	http.Handle("/transactions", transactions.NewHandler(host, logger))
	http.Handle("/transaction/info", info.NewHandler(host, settings, watch, logger))
	http.Handle("/transaction/output/progress", progress.NewHandler(host, settings, watch, logger))
	http.Handle("/wallet/address", address.NewHandler(logger))
	http.Handle("/wallet/amount", amount.NewHandler(host, settings, watch, logger))
	logger.Info("user interface server is running...")
	logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), nil).Error())
}
