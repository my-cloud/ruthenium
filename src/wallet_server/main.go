package main

import (
	"flag"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Wallet: ")
}

func main() {
	port := flag.Uint("port", 8080, "TCP port number for wallet server")
	hostIp := flag.String("host-ip", "", "Blockchain host IP address")
	hostPort := flag.Uint("host-port", 5000, "Blockchain host port")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *hostIp, uint16(*hostPort))
	fmt.Println("Running...")
	app.Run()
}
