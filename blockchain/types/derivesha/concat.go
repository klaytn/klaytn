package derivesha

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto/sha3"
)

// An alternative implementation of DeriveSha()
// This function generates a hash of `DerivableList` as below:
// 1. make a byte slice by concatenating RLP-encoded items
// 2. make a hash of the byte slice.
type DeriveShaConcat struct{}

func (d DeriveShaConcat) DeriveSha(list types.DerivableList) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	for i := 0; i < list.Len(); i++ {
		hasher.Write(list.GetRlp(i))
	}
	hasher.Sum(hash[:0])

	return hash
}
