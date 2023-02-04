package p2p

import (
	"fmt"
	"net"
	"strings"
)

type Target struct {
	ip    string
	port  string
	value string
}

func NewTarget(target string) (*Target, error) {
	separator := ":"
	separatorIndex := strings.Index(target, separator)
	if separatorIndex == -1 {
		return nil, fmt.Errorf("seed target format is invalid (should be of the form 89.82.76.241:8106")
	}
	ip := target[:separatorIndex]
	port := target[separatorIndex+1:]
	value := net.JoinHostPort(ip, port)
	return &Target{ip, port, value}, nil
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
