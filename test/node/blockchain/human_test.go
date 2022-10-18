package blockchain

import (
	"github.com/my-cloud/ruthenium/src/node/blockchain"
	"github.com/my-cloud/ruthenium/test"
	"testing"
)

func Test_IsRegistered_NotRegistered_ReturnsFalse(t *testing.T) {
	address := "0x0000000000000000000000000000000000000001"
	human := blockchain.NewHuman(address)
	isRegistered, _ := human.IsRegistered()
	test.Assert(t, !isRegistered, "proof of humanity is valid whereas it should not")
}

func Test_IsRegistered_Registered_ReturnsTrue(t *testing.T) {
	address := "0xf14DB86A3292ABaB1D4B912dbF55e8abc112593a"
	human := blockchain.NewHuman(address)
	isRegistered, _ := human.IsRegistered()
	test.Assert(t, isRegistered, "proof of humanity is invalid whereas it should be")
}
