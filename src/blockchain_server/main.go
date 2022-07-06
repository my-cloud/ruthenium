package main

import (
	"flag"
	"fmt"
	"log"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", 5000, "TCP port number for blockchain server")
	flag.Parse()
	app := NewBlockchainServer(uint16(*port))
	fmt.Println("Running...")
	app.Run()
}
