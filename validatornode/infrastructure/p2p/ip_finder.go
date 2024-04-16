package p2p

type IpFinder interface {
	LookupIP(ip string) (string, error)
}
