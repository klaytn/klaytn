package derivesha

import (
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/statedb"
)

type DeriveShaOrig struct{}

func (d DeriveShaOrig) DeriveSha(list types.DerivableList) common.Hash {
	trie := statedb.NewStackTrie(nil)
	trie.Reset()
	var buf []byte

	// StackTrie requires values to be inserted in increasing
	// hash order, which is not the order that `list` provides
	// hashes in. This insertion sequence ensures that the
	// order is correct.
	for i := 1; i < list.Len() && i <= 0x7f; i++ {
		buf = rlp.AppendUint64(buf[:0], uint64(i))
		trie.Update(buf, list.GetRlp(i))
	}
	if list.Len() > 0 {
		buf = rlp.AppendUint64(buf[:0], 0)
		trie.Update(buf, list.GetRlp(0))
	}
	for i := 0x80; i < list.Len(); i++ {
		buf = rlp.AppendUint64(buf[:0], uint64(i))
		trie.Update(buf, list.GetRlp(i))
	}
	return trie.Hash()
}
