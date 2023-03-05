package p2p

import (
	"fmt"
	"net"
)

type Target struct {
	ip        string
	port      string
	value     string
	networkId string
}

func NewTarget(ip string, port string) *Target {
	value := net.JoinHostPort(ip, port)
	return &Target{ip, port, value, networkId(port)}
}

func NewTargetFromValue(value string) (*Target, error) {
	ip, port, err := net.SplitHostPort(value)
	if err != nil {
		return nil, fmt.Errorf("seed target format is invalid: %w", err)
	}
	return &Target{ip, port, value, networkId(port)}, nil
}

func (target *Target) Ip() string {
	return target.ip
}

func (target *Target) IsSameNetworkId(other *Target) bool {
	return target.networkId == other.networkId
}

func (target *Target) Port() string {
	return target.port
}

func (target *Target) Value() string {
	return target.value
}

func networkId(port string) string {
	if port == "10600" {
		return "mainnet"
	} else if len(port) == 5 && port[:3] == "106" {
		return "testnet"
	} else {
		return "unknown"
	}
}
