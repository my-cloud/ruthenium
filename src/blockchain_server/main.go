package main

import (
	"flag"
	"fmt"
	"ruthenium/src/chain"
	"ruthenium/src/log"
)

func main() {
	publicKey := flag.String("public-key", "", "The public key (will be generated if not provided")
	privateKey := flag.String("private-key", "", "The private key (will be generated if not provided")
	port := flag.Uint("port", chain.DefaultPort, "TCP port number for blockchain server")
	logLevel := flag.String("log-level", "warn", "The log level")

	flag.Parse()
	app := chain.NewHost(*publicKey, *privateKey, uint16(*port), log.ParseLevel(*logLevel))
	fmt.Println("Running...")
	app.Run()
}
