package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/my-cloud/ruthenium/accessnode/application/index"
	"github.com/my-cloud/ruthenium/accessnode/application/transaction"
	"github.com/my-cloud/ruthenium/accessnode/application/transaction/info"
	"github.com/my-cloud/ruthenium/accessnode/application/transaction/output/progress"
	"github.com/my-cloud/ruthenium/accessnode/application/transactions"
	"github.com/my-cloud/ruthenium/accessnode/application/wallet/address"
	"github.com/my-cloud/ruthenium/accessnode/application/wallet/amount"
	"github.com/my-cloud/ruthenium/validatornode/application/p2p"
	"github.com/my-cloud/ruthenium/validatornode/domain/clock"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/config"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/environment"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/gp2p"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log/console"
	"github.com/my-cloud/ruthenium/validatornode/infrastructure/net"
)

func main() {
	port := flag.Int("port", environment.NewVariable("PORT").GetIntValue(8080), "The TCP port number of the access node")
	hostIp := flag.String("host-ip", environment.NewVariable("HOST_IP").GetStringValue("127.0.0.1"), "The validator node IP or DNS address")
	hostPort := flag.Int("host-port", environment.NewVariable("HOST_PORT").GetIntValue(10600), "The TCP port number of the validator node")
	templatePath := flag.String("template-path", environment.NewVariable("TEMPLATE_PATH").GetStringValue("accessnode/template.html"), "The UI template path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level (possible values: 'debug', 'info', 'warn', 'error', 'fatal')")

	flag.Parse()
	logger := console.NewLogger(console.ParseLevel(*logLevel))
	target := p2p.NewTarget(*hostIp, strconv.Itoa(*hostPort))
	ipFinder := net.NewIpFinderImplementation(logger)
	clientFactory := gp2p.NewSenderFactory(ipFinder, time.Minute)
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
	watch := clock.NewWatch()
	http.Handle("/", index.NewHandler(*templatePath, logger))
	http.Handle("/transaction", transaction.NewHandler(host, logger))
	http.Handle("/transactions", transactions.NewHandler(host, logger))
	http.Handle("/transaction/info", info.NewHandler(host, settings, watch, logger))
	http.Handle("/transaction/output/progress", progress.NewHandler(host, settings, watch, logger))
	http.Handle("/wallet/address", address.NewHandler(logger))
	http.Handle("/wallet/amount", amount.NewHandler(host, settings, watch, logger))
	logger.Info("user interface server is running...")
	logger.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*port), nil).Error())
}
