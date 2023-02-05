package network

import (
	"net"
)

type IpFinder interface {
	LookupIP(ip string) ([]net.IP, error)
	FindHostPublicIp() (ip string, err error)
}
