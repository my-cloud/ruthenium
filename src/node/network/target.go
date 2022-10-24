package network

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

type Target struct {
	ip   string
	port uint16
}

func NewTarget(ip string, port uint16) *Target {
	return &Target{ip, port}
}

func (target *Target) Ip() string {
	return target.ip
}

func (target *Target) Port() uint16 {
	return target.port
}

func (target *Target) Value() string {
	return net.JoinHostPort(target.ip, strconv.Itoa(int(target.port)))
}

func (target *Target) Reach() error {
	lookedUpIps, err := net.LookupIP(target.ip)
	if err != nil {
		return fmt.Errorf("DNS discovery failed on addresse %s: %w", target.ip, err)
	}

	ipsCount := len(lookedUpIps)
	if ipsCount != 1 {
		return errors.New(fmt.Sprintf("DNS discovery did not find a single address (%d addresses found) for the given IP %s", ipsCount, target.ip))
	}
	_, err = net.DialTimeout("tcp", target.Value(), NeighborFindingTimeoutSecond*time.Second)
	return err
}
