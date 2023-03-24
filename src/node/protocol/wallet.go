package protocol

type Wallet interface {
	Amount(currentTimestamp int64) uint64
}
