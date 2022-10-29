package main

import (
	"flag"
	"github.com/my-cloud/ruthenium/src/environment"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/network"
	"github.com/my-cloud/ruthenium/src/ui/server"
)

func main() {
	mnemonic := flag.String("mnemonic", environment.NewVariable("MNEMONIC").GetStringValue(""), "The mnemonic (required if the private key is not provided)")
	derivationPath := flag.String("derivation-path", environment.NewVariable("DERIVATION_PATH").GetStringValue("m/44'/60'/0'/0/0"), "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", environment.NewVariable("PASSWORD").GetStringValue(""), "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", environment.NewVariable("PRIVATE_KEY").GetStringValue(""), "The private key (required if the mnemonic is not provided, unused if the mnemonic is provided)")
	port := flag.Uint64("port", environment.NewVariable("PORT").GetUint64Value(server.DefaultPort), "The TCP port number for the UI server")
	hostIp := flag.String("host-ip", environment.NewVariable("HOST_IP").GetStringValue(""), "The blockchain host IP address")
	hostPort := flag.Uint64("host-port", environment.NewVariable("HOST_PORT").GetUint64Value(network.DefaultPort), "The TCP port number for the protocol host node")
	templatesPath := flag.String("templates-path", environment.NewVariable("TEMPLATES_PATH").GetStringValue("src/ui/templates"), "The UI templates path")
	logLevel := flag.String("log-level", environment.NewVariable("LOG_LEVEL").GetStringValue("info"), "The log level")

	flag.Parse()
	app := server.NewController(*mnemonic, *derivationPath, *password, *privateKey, uint16(*port), *hostIp, uint16(*hostPort), *templatesPath, log.ParseLevel(*logLevel))
	app.Run()
}
