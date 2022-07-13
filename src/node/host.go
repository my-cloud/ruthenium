package node

import (
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

const (
	StartPort     uint16 = 5000
	EndPort       uint16 = 5002
	StartIpSuffix uint8  = 0
	EndIpSuffix   uint8  = 0
)

var PATTERN = regexp.MustCompile(`((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?\.){3})(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)

type Host struct {
	ip   string
	port uint16
}

func NewHost(port uint16) *Host {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "127.0.0.1"
	}
	ips, err := net.LookupHost(hostname)
	if err != nil {
		ips[0] = "127.0.0.1"
	}
	// FIXME host[3] = 192.168.1.90
	return &Host{"127.0.0.1", port}
}

func (host *Host) FindNeighbors() []*Neighbor {
	address := fmt.Sprintf("%s:%d", host.ip, host.port)

	m := PATTERN.FindStringSubmatch(host.ip)
	if m == nil {
		return nil
	}
	prefixHost := m[1]
	lastIp, err := strconv.Atoi(m[len(m)-1])
	if err != nil {
		fmt.Printf("ERROR: Failed to parse IP %s, err:%v\n", m[len(m)-1], err)
	}
	neighbors := make([]*Neighbor, 0)

	for port := StartPort; port <= EndPort; port += 1 {
		for ipSuffix := StartIpSuffix; ipSuffix <= EndIpSuffix; ipSuffix += 1 {
			guessIp := fmt.Sprintf("%s%d", prefixHost, lastIp+int(ipSuffix))
			neighbor := NewNeighbor(guessIp, port)
			guessTarget := neighbor.IpAndPort()
			if guessTarget != address && neighbor.isFound() {
				// FIXME set only one read-writer by couple host-neighbor
				neighbor.SetReadWriter(port)
				neighbors = append(neighbors, neighbor)
			}
		}
	}
	return neighbors
}
