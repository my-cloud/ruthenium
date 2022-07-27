package main

import (
	"flag"
	"fmt"
	"ruthenium/src/chain"
	"ruthenium/src/log"
)

func main() {
	port := flag.Uint("port", chain.DefaultPort, "TCP port number for blockchain server")
	logLevel := flag.String("log-level", "warn", "The log level")
	flag.Parse()
	app := chain.NewHost(uint16(*port), log.ParseLevel(*logLevel))
	fmt.Println("Running...")
	app.Run()
}
