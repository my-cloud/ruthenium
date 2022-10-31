package p2p

import (
	"net"
	"time"
)

type IpFinder interface {
	LookupIP(ip string) ([]net.IP, error)
	DialTimeout(network, address string, timeout time.Duration) (net.Conn, error)
}
