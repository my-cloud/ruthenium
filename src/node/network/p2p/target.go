package p2p

import (
	"fmt"
	"net"
)

type Target struct {
	ip    string
	port  string
	value string
}

func NewTarget(ip string, port string) *Target {
	value := net.JoinHostPort(ip, port)
	return &Target{ip, port, value}
}

func NewTargetFromValue(value string) (*Target, error) {
	ip, port, err := net.SplitHostPort(value)
	if err != nil {
		return nil, fmt.Errorf("seed target format is invalid: %w", err)
	}
	return &Target{ip, port, value}, nil
}

func (target *Target) IsSameNetworkId(other *Target) bool {
	return target.networkId() == other.networkId()
}

func (target *Target) Ip() string {
	return target.ip
}

func (target *Target) Port() string {
	return target.port
}

func (target *Target) Value() string {
	return target.value
}

func (target *Target) networkId() string {
	if target.port == "10600" {
		return "mainnet"
	} else if len(target.port) == 5 && target.port[:3] == "106" {
		return "testnet"
	} else {
		return "unknown"
	}
}
