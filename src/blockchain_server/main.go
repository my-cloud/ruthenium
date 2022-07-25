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
	flag.Parse()
	app := chain.NewHost()
	fmt.Println("Running...")
	app.Run()
}
