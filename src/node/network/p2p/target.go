package p2p

import (
	"net"
	"strconv"
)

type Target struct {
	ip    string
	port  uint16
	value string
}

func NewTarget(ip string, port uint16) *Target {
	value := net.JoinHostPort(ip, strconv.Itoa(int(port)))
	return &Target{ip, port, value}
}

func (target *Target) Ip() string {
	return target.ip
}

func (target *Target) Port() uint16 {
	return target.port
}

func (target *Target) Value() string {
	return target.value
}
