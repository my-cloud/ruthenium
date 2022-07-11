package chain

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"
)

const (
	StartPort uint16 = 5000
	EndPort   uint16 = 5003
	StartIp   uint8  = 0
	EndIp     uint8  = 0
)

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

type Node struct {
	host string
	port uint16
}

func NewHostNode(port uint16) *Node {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "127.0.0.1"
	}
	host, err := net.LookupHost(hostname)
	if err != nil {
		host[0] = "127.0.0.1"
	}
	// FIXME host[3] = 192.168.1.90
	return &Node{"127.0.0.1", port}
}

func newNode(host string, port uint16) *Node {
	return &Node{host, port}
}

func (node *Node) FindNeighbors() []string {
	address := fmt.Sprintf("%s:%d", node.host, node.port)

	m := PATTERN.FindStringSubmatch(node.host)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, _ := strconv.Atoi(m[len(m)-1])
	neighbors := make([]string, 0)

	for port := StartPort; port <= EndPort; port += 1 {
		for ip := StartIp; ip <= EndIp; ip += 1 {
			guessHost := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ip))
			guessTarget := fmt.Sprintf("%s:%d", guessHost, port)
			neighbor := newNode(guessHost, port)
			if guessTarget != address && neighbor.isFound() {
				neighbors = append(neighbors, guessTarget)
			}
		}
	}
	return neighbors
}

func (node *Node) isFound() bool {
	target := fmt.Sprintf("%s:%d", node.host, node.port)

	_, err := net.DialTimeout("tcp", target, time.Millisecond)
	if err != nil {
		fmt.Printf("%s not found, err:%v\n", target, err)
		return false
	}
	fmt.Printf("%s found\n", target)
	return true
}
