package main

import (
	"flag"
	"fmt"
	"ruthenium/src/chain"
	"ruthenium/src/log"
	"ruthenium/src/wallet_server/wallet"
)

func main() {
	publicKey := flag.String("public-key", "", "The public key (will be generated if not provided")
	privateKey := flag.String("private-key", "", "The private key (will be generated if not provided")
	port := flag.Uint("port", 8080, "TCP port number for wallet server")
	hostIp := flag.String("host-ip", "", "Blockchain host IP address")
	hostPort := flag.Uint("host-port", chain.DefaultPort, "Blockchain host port")
	logLevel := flag.String("log-level", "warn", "The log level")
	flag.Parse()

	app := wallet.NewWalletServer(*publicKey, *privateKey, uint16(*port), *hostIp, uint16(*hostPort), log.ParseLevel(*logLevel))
	fmt.Println("Running...")
	app.Run()
}
