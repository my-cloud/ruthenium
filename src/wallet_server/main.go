package main

import (
	"flag"
	"fmt"
	"ruthenium/src/log"
	"ruthenium/src/wallet_server/wallet"
)

func main() {
	port := flag.Uint("port", 8080, "TCP port number for wallet server")
	hostIp := flag.String("host-ip", "", "Blockchain host IP address")
	hostPort := flag.Uint("host-port", 5000, "Blockchain host port")
	logLevel := flag.String("log-level", "warn", "The log level")
	flag.Parse()

	app := wallet.NewWalletServer(uint16(*port), *hostIp, uint16(*hostPort), log.ParseLevel(*logLevel))
	fmt.Println("Running...")
	app.Run()
}
