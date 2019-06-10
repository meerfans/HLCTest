package blockchain

import (
	"Nox-DAG-test/script/tool/types"
	"math"
	"Nox-DAG-test/script/tool/hash"
)
var zeroHash = &hash.ZeroHash
//
// This function only differs from IsCoinBase in that it works with a raw wire
// transaction as opposed to a higher level util transaction.
func IsCoinBaseTx(tx *types.Transaction) bool {
	// A coin base must only have one transaction input.
	if len(tx.TxIn) != 1 {
		return false
	}
	// The previous output of a coin base must have a max value index and a
	// zero hash.
	prevOut := &tx.TxIn[0].PreviousOut
	if prevOut.OutIndex != math.MaxUint32 || !prevOut.Hash.IsEqual(zeroHash) {
		return false
	}
	return true
}
