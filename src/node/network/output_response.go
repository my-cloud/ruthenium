package network

type OutputResponse struct {
	Address     string
	BlockHeight int
	HasReward   bool
	HasIncome   bool
	Value       uint64
}
