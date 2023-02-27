package net

import (
	"fmt"
	"github.com/my-cloud/ruthenium/src/log"
	"io"
	"net"
	"net/http"
)

type IpFinder struct {
	logger log.Logger
}

func NewIpFinder(logger log.Logger) *IpFinder {
	return &IpFinder{logger}
}

func (finder *IpFinder) LookupIP(ip string) (string, error) {
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

func (finder *IpFinder) FindHostPublicIp() (string, error) {
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
