package net

type IpFinder interface {
	LookupIP(ip string) (string, error)
	FindHostPublicIp() (string, error)
}
