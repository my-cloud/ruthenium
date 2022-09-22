package main

import (
	"flag"
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/blockchain"
)

func main() {
	mnemonic := flag.String("mnemonic", "", "The mnemonic (optional)")
	derivationPath := flag.String("derivation-path", "m/44'/60'/0'/0/0", "The derivation path (unused if the mnemonic is omitted)")
	password := flag.String("password", "", "The mnemonic password (unused if the mnemonic is omitted)")
	privateKey := flag.String("private-key", "", "The private key (will be generated if not provided)")
	port := flag.Uint("port", blockchain.DefaultPort, "TCP port number for blockchain server")
	logLevel := flag.String("log-level", "info", "The log level")

	flag.Parse()
	app := blockchain.NewHost(*mnemonic, *derivationPath, *password, *privateKey, uint16(*port), log.ParseLevel(*logLevel))
	app.Run()
}
