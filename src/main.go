package main

import (
	"flag"
	"fmt"
	"ruthenium/src/server"
)

func main() {
	port := flag.Uint("port", 5000, "TCP port number for blockchain server")
	flag.Parse()
	app := server.NewServer(uint16(*port))
	fmt.Println("Running...")
	app.Run()
}
