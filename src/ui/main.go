package main

import (
	"coinsmaster/ruthenium/src/log"
	"coinsmaster/ruthenium/src/node/blockchain"
	"coinsmaster/ruthenium/src/ui/server"
	"flag"
	"fmt"
)

func main() {
	publicKey := flag.String("public-key", "", "The public key (will be generated if not provided")
	privateKey := flag.String("private-key", "", "The private key (will be generated if not provided")
	port := flag.Uint("port", 8080, "The UI server TCP port number")
	hostIp := flag.String("host-ip", "", "The blockchain host IP address")
	hostPort := flag.Uint("host-port", blockchain.DefaultPort, "The blockchain host port")
	templatesPath := flag.String("templates-path", "src/ui/templates", "The UI templates path")
	logLevel := flag.String("log-level", "warn", "The log level")
	flag.Parse()

	app := server.NewController(*publicKey, *privateKey, uint16(*port), *hostIp, uint16(*hostPort), *templatesPath, log.ParseLevel(*logLevel))
	fmt.Println("Running...")
	app.Run()
}
