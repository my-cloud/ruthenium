package main

import (
	"flag"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/blockchain"
	"github.com/my-cloud/ruthenium/src/ui/server"
)

func main() {
	mnemonic := flag.String("mnemonic", "", "The mnemonic (optional)")
	derivationPath := flag.String("derivation-path", "m/44'/60'/0'/0/0", "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", "", "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", "", "The private key (will be generated if not provided")
	port := flag.Uint("port", 8080, "The UI server TCP port number")
	hostIp := flag.String("host-ip", "", "The blockchain host IP address")
	hostPort := flag.Uint("host-port", blockchain.DefaultPort, "The blockchain host port")
	templatesPath := flag.String("templates-path", "src/ui/templates", "The UI templates path")
	logLevel := flag.String("log-level", "info", "The log level")
	flag.Parse()

	app := server.NewController(*mnemonic, *derivationPath, *password, *privateKey, uint16(*port), *hostIp, uint16(*hostPort), *templatesPath, log.ParseLevel(*logLevel))
	app.Run()
}
