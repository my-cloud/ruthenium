package main

import (
	"flag"
	"fmt"
	"gitlab.com/coinsmaster/ruthenium/src/log"
	"gitlab.com/coinsmaster/ruthenium/src/node/blockchain"
)

func main() {
	publicKey := flag.String("public-key", "", "The public key (will be generated if not provided")
	privateKey := flag.String("private-key", "", "The private key (will be generated if not provided")
	port := flag.Uint("port", blockchain.DefaultPort, "TCP port number for blockchain server")
	logLevel := flag.String("log-level", "warn", "The log level")

	flag.Parse()
	app := blockchain.NewHost(*publicKey, *privateKey, uint16(*port), log.ParseLevel(*logLevel))
	fmt.Println("Running...")
	app.Run()
}
