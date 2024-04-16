package net

import (
	"fmt"
	"net"

	"github.com/my-cloud/ruthenium/validatornode/infrastructure/log"
)

type IpFinderImplementation struct {
	logger log.Logger
}

func NewIpFinderImplementation(logger log.Logger) *IpFinderImplementation {
	return &IpFinderImplementation{logger}
}

func (finder *IpFinderImplementation) LookupIP(ip string) (string, error) {
	ips, err := net.LookupIP(ip)
	if err != nil {
		return "", fmt.Errorf("DNS discovery failed on addresse %s: %w", ip, err)
	}
	ipsCount := len(ips)
	if ipsCount != 1 {
		return "", fmt.Errorf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", ipsCount, ip)
	}
	return ips[0].String(), nil
}
