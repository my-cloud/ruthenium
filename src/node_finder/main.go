package main

import (
	"fmt"
)

func main() {
	//fmt.Println(IsFoundHost("127.0.0.1", 5000))
	fmt.Println(FindNeighbors("127.0.0.1", 5000, 0, 3, 5000, 5003))
}
