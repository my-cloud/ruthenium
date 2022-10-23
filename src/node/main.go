package main

import (
	"flag"
	"fmt"
	"github.com/my-cloud/ruthenium/src/clock"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/humanity"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/encryption"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"io/ioutil"
	"net/http"
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
	ip, err := findPublicIp(logger)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to find the public IP: %w", err).Error())
	}
	wallet, err := encryption.DecodeWallet(*mnemonic, *derivationPath, *password, *privateKey)
	if err != nil {
		logger.Fatal(fmt.Errorf("failed to create wallet: %w", err).Error())
	} else {
		registry := humanity.NewRegistry()
		validationTimer := validationIntervalInSeconds * time.Second
		watch := clock.NewWatch()
		blockchain := protocol.NewBlockchain(registry, validationTimer, watch, logger)
		pool := protocol.NewPool(registry, watch, logger)
		validation := protocol.NewValidation(wallet.Address(), blockchain, pool, watch, validationTimer, logger)
		neighborhood := network.NewNeighborhood(ip, uint16(*port), watch, *configurationPath, logger)
		app := network.NewHost(ip, uint16(*port), blockchain, pool, validation, neighborhood, logger)
		app.Run()
	}
}

func findPublicIp(logger *log.Logger) (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return
	}
	defer func() {
		if bodyCloseError := resp.Body.Close(); bodyCloseError != nil {
			logger.Error(fmt.Errorf("failed to close public IP request body: %w", bodyCloseError).Error())
		}
	}()
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ip = string(body)
	return
}
