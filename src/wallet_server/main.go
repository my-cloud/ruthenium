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
	gateway := flag.String("gateway", "http://127.0.0.1:5000", "Blockchain gateway")
	flag.Parse()

	app := NewWalletServer(uint16(*port), *gateway)
	fmt.Println("Running...")
	app.Run()
}
