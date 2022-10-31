package net

import (
	"net"
	"time"
)

type IpFinder struct{}

func NewIpFinder() *IpFinder {
	return &IpFinder{}
}

func (finder *IpFinder) LookupIP(ip string) ([]net.IP, error) {
	return net.LookupIP(ip)
}

func (finder *IpFinder) DialTimeout(network, address string, timeout time.Duration) (net.Conn, error) {
	return net.DialTimeout(network, address, timeout)
}
