package protocol

import (
	"github.com/my-cloud/ruthenium/src/log"
	"github.com/my-cloud/ruthenium/src/node/neighborhood"
	"github.com/my-cloud/ruthenium/src/node/protocol"
	"github.com/my-cloud/ruthenium/test"
	"testing"
	"time"
)

func Test_Verify_NeighborBlockchainIsBetter_IsReplaced(t *testing.T) {
	// Arrange
	registry := new(RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*neighborhood.BlockResponse, error) {
		blockResponse1 := NewRewardedBlockResponse([32]byte{}, 0)
		block1, _ := protocol.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := NewRewardedBlockResponse(hash, 1)
		return []*neighborhood.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(SynchronizerMock)
	synchronizer.NeighborsFunc = func() []neighborhood.Neighbor {
		return []neighborhood.Neighbor{neighborMock}
	}
	blockchain := protocol.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

	// Act
	blockchain.Verify()

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, isReplaced, "blockchain is not replaced whereas it should be")
}

func Test_Verify_NeighborNewBlockTimestampIsInvalid_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 2) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(NeighborMock)
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(SynchronizerMock)
	synchronizer.NeighborsFunc = func() []neighborhood.Neighbor {
		return []neighborhood.Neighbor{neighborMock}
	}
	blockchain := protocol.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

	type args struct {
		firstBlockTimestamp  int64
		secondBlockTimestamp int64
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "SecondTimestampBeforeTheFirstOne",
			args: args{
				firstBlockTimestamp:  1,
				secondBlockTimestamp: 0,
			},
		},
		{
			name: "BlockMissing",
			args: args{
				firstBlockTimestamp:  0,
				secondBlockTimestamp: 2,
			},
		},
		{
			name: "SameZeroedTimestamp",
			args: args{
				firstBlockTimestamp:  0,
				secondBlockTimestamp: 0,
			},
		},
		{
			name: "SameNonZeroTimestamp",
			args: args{
				firstBlockTimestamp:  1,
				secondBlockTimestamp: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			neighborMock.GetBlocksFunc = func() ([]*neighborhood.BlockResponse, error) {
				blockResponse1 := NewRewardedBlockResponse([32]byte{}, tt.args.firstBlockTimestamp)
				block1, _ := protocol.NewBlockFromResponse(blockResponse1)
				hash, _ := block1.Hash()
				blockResponse2 := NewRewardedBlockResponse(hash, tt.args.secondBlockTimestamp)
				return []*neighborhood.BlockResponse{blockResponse1, blockResponse2}, nil
			}

			// Act
			blockchain.Verify()

			// Assert
			if blockchain.IsReplaced() {
				t.Errorf("blockchain is replaced whereas it should not")
			}
		})
	}
}

func Test_Verify_NeighborNewBlockTimestampIsInTheFuture_IsNotReplaced(t *testing.T) {
	// Arrange
	registry := new(RegistryMock)
	registry.IsRegisteredFunc = func(address string) (bool, error) { return true, nil }
	timeMock := new(TimeMock)
	timeMock.NowFunc = func() time.Time { return time.Unix(0, 1) }
	logger := log.NewLogger(log.Fatal)
	neighborMock := new(NeighborMock)
	neighborMock.GetBlocksFunc = func() ([]*neighborhood.BlockResponse, error) {
		blockResponse1 := NewRewardedBlockResponse([32]byte{}, 1)
		block1, _ := protocol.NewBlockFromResponse(blockResponse1)
		hash, _ := block1.Hash()
		blockResponse2 := NewRewardedBlockResponse(hash, 2)
		return []*neighborhood.BlockResponse{blockResponse1, blockResponse2}, nil
	}
	neighborMock.TargetFunc = func() string {
		return "neighbor"
	}
	synchronizer := new(SynchronizerMock)
	synchronizer.NeighborsFunc = func() []neighborhood.Neighbor {
		return []neighborhood.Neighbor{neighborMock}
	}
	blockchain := protocol.NewBlockchain(registry, 1, timeMock, synchronizer, logger)

	// Act
	blockchain.Verify()

	// Assert
	isReplaced := blockchain.IsReplaced()
	test.Assert(t, !isReplaced, "blockchain is replaced whereas it should not")
}
