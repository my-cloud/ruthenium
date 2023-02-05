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

func (finder *IpFinder) LookupIP(ip string) ([]net.IP, error) {
	return net.LookupIP(ip)
}

func (finder *IpFinder) FindHostPublicIp() (ip string, err error) {
	resp, err := http.Get("https://ifconfig.me")
	if err != nil {
		return
	}
	defer func() {
		if bodyCloseError := resp.Body.Close(); bodyCloseError != nil {
			finder.logger.Error(fmt.Errorf("failed to close public IP request body: %w", bodyCloseError).Error())
		}
	}()
	var body []byte
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	ip = string(body)
	return
}
