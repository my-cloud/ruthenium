package main

import (
	"flag"
	"fmt"
	"log"
	"ruthenium/src/chain"
)

func init() {
	log.SetPrefix("Blockchain: ")
}

func main() {
	port := flag.Uint("port", chain.DefaultPort, "TCP port number for blockchain server")
	flag.Parse()
	app := chain.NewHost(uint16(*port))
	fmt.Println("Running...")
	app.Run()
}
