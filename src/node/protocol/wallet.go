package protocol

type Wallet interface {
	Balance(currentTimestamp int64) uint64
}
