// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/types/derive_sha.go (2018/06/04).
// Modified and improved for the klaytn development.

package types

import (
	"math/big"

	"github.com/klaytn/klaytn/common"
)

type DerivableList interface {
	Len() int
	GetRlp(i int) []byte
}

const (
	ImplDeriveShaOriginal int = iota
	ImplDeriveShaSimple
	ImplDeriveShaConcat
)

var (
	// EmptyRootHashOriginal is the empty root hash of a state trie,
	// which is equal to EmptyRootHash with ImplDeriveShaOriginal.
	EmptyRootHashOriginal = common.HexToHash("0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// DeriveSha and EmptyRootHash are populated by derivesha.InitDeriveSha().
	// DeriveSha is used to calculate TransactionsRoot and ReceiptsRoot.
	// EmptyRootHash is a transaction/receipt root hash when there is no transaction.
	DeriveSha     func(list DerivableList, num *big.Int) common.Hash = DeriveShaNone
	EmptyRootHash func(num *big.Int) common.Hash                     = EmptyRootHashNone
)

func DeriveShaNone(list DerivableList, num *big.Int) common.Hash {
	logger.Crit("DeriveSha not initialized")
	return common.Hash{}
}

func EmptyRootHashNone(num *big.Int) common.Hash {
	logger.Crit("DeriveSha not initialized")
	return common.Hash{}
}
