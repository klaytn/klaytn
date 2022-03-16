// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from ethclient/signer.go (2018/06/04).
// Modified and improved for the klaytn development.

package client

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
)

// senderFromServer is a types.Signer that remembers the sender address returned by the RPC
// server. It is stored in the transaction's sender address cache to avoid an additional
// request in TransactionSender.
type senderFromServer struct {
	addr      common.Address
	blockhash common.Hash
}

var errNotCached = errors.New("sender not cached")

func setSenderFromServer(tx *types.Transaction, addr common.Address, block common.Hash) {
	// Use types.Sender for side-effect to store our signer into the cache.
	types.Sender(&senderFromServer{addr, block}, tx)
}

func (s *senderFromServer) Equal(other types.Signer) bool {
	os, ok := other.(*senderFromServer)
	return ok && os.blockhash == s.blockhash
}

func (s *senderFromServer) Sender(tx *types.Transaction) (common.Address, error) {
	if s.blockhash == (common.Hash{}) {
		return common.Address{}, errNotCached
	}
	return s.addr, nil
}

func (s *senderFromServer) ChainID() *big.Int {
	// TODO-Klaytn: need to check this routine is never called or not.
	// `senderFromServer` is only used in klay_client.go.
	panic("ChainID should not be called!")
}

func (s *senderFromServer) SenderPubkey(tx *types.Transaction) ([]*ecdsa.PublicKey, error) {
	// TODO-Klaytn: need to check this routine is never called or not.
	// `senderFromServer` is only used in klay_client.go.
	panic("SenderPubkey should not be called!")
}

func (s *senderFromServer) SenderFeePayer(tx *types.Transaction) ([]*ecdsa.PublicKey, error) {
	// TODO-Klaytn: need to check this routine is never called or not.
	// `senderFromServer` is only used in klay_client.go.
	panic("SenderFeePayer should not be called!")
}

func (s *senderFromServer) Hash(tx *types.Transaction) common.Hash {
	// TODO-Klaytn: need to check this routine is never called or not.
	// `senderFromServer` is only used in klay_client.go.
	panic("can't sign with senderFromServer")
}

func (s *senderFromServer) HashFeePayer(tx *types.Transaction) (common.Hash, error) {
	// TODO-Klaytn: need to check this routine is never called or not.
	// `senderFromServer` is only used in klay_client.go.
	panic("can't sign with senderFromServer")
}

func (s *senderFromServer) SignatureValues(tx *types.Transaction, sig []byte) (R, S, V *big.Int, err error) {
	// TODO-Klaytn: need to check this routine is never called or not.
	// `senderFromServer` is only used in klay_client.go.
	panic("can't sign with senderFromServer")
}
