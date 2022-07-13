package main

import (
	"fmt"
	"ruthenium/src/chain"
)

func main() {
	fmt.Println(chain.NewHostNode(5000).FindNeighbors())
}
