package net

import (
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/my-cloud/ruthenium/common/infrastructure/log"
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

func (finder *IpFinderImplementation) FindHostPublicIp() (string, error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return "", err
	}
	defer func() {
		if bodyCloseError := resp.Body.Close(); bodyCloseError != nil {
			finder.logger.Error(fmt.Errorf("failed to close public IP request body: %w", bodyCloseError).Error())
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
