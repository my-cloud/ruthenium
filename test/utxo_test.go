package test

import (
	"github.com/my-cloud/ruthenium/src/node/protocol/validation"
	"testing"
)

func Test_Utxo(t *testing.T) {
	var addressByUtxo = make(map[[32]byte]string)
	tx1R := validation.NewRewardTransaction("A", 0, 1)
	tx1, _ := validation.NewTransactionFromResponse(tx1R)
	id1 := tx1.Id()
	addressByUtxo[id1] = "A"
	tx2R := validation.NewRewardTransaction("B", 0, 1)
	tx2, _ := validation.NewTransactionFromResponse(tx2R)
	id2 := tx2.Id()
	addressByUtxo[id2] = "B"
	print(addressByUtxo)
}
