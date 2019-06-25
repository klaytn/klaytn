// Copyright 2019 The klaytn Authors
// This file is part of the klaytn library.
//
// The klaytn library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The klaytn library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the klaytn library. If not, see <http://www.gnu.org/licenses/>.

package sc

import (
	"crypto/ecdsa"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"math/big"
	"sync"
)

type accountInfo struct {
	key      *ecdsa.PrivateKey
	address  common.Address
	nonce    uint64
	chainID  *big.Int
	gasPrice *big.Int

	isNonceSynced bool
	mu            sync.RWMutex
}

// BridgeAccountManager manages bridge account for main/service chain.
type BridgeAccountManager struct {
	// TODO-Klaytn need to consider multiple bridge accounts?
	mcAccount *accountInfo
	scAccount *accountInfo
}

// NewBridgeAccountManager returns bridgeAccountManager created by main/service bridge account keys.
func NewBridgeAccountManager(mcKey, scKey *ecdsa.PrivateKey) (*BridgeAccountManager, error) {
	mcAcc := &accountInfo{
		key:      mcKey,
		address:  crypto.PubkeyToAddress(mcKey.PublicKey),
		nonce:    0,
		chainID:  nil,
		gasPrice: nil,
	}

	scAcc := &accountInfo{
		key:      scKey,
		address:  crypto.PubkeyToAddress(scKey.PublicKey),
		nonce:    0,
		chainID:  nil,
		gasPrice: nil,
	}

	return &BridgeAccountManager{
		mcAccount: mcAcc,
		scAccount: scAcc,
	}, nil
}

// GetTransactOpts return a transactOpts for transact on local/remote backend.
func (acc *accountInfo) GetTransactOpts() *bind.TransactOpts {
	var nonce *big.Int

	// Only for unit test, if the nonce is not synced yet, return transaction option with nil nonce.
	// Backend will use state nonce.
	if acc.isNonceSynced {
		nonce = new(big.Int).SetUint64(acc.nonce)
	}
	return MakeTransactOpts(acc.key, nonce, acc.chainID, acc.gasPrice)
}

// SetChainID sets the chain ID of the chain of the account.
func (acc *accountInfo) SetChainID(cID *big.Int) {
	acc.chainID = cID
}

// SetGasPrice sets the gas price of the chain of the account.
func (acc *accountInfo) SetGasPrice(gp *big.Int) {
	acc.gasPrice = gp
}

// Lock can lock the account for nonce management.
func (acc *accountInfo) Lock() {
	acc.mu.Lock()
}

// UnLock can unlock the account for nonce management.
func (acc *accountInfo) UnLock() {
	acc.mu.Unlock()
}

// SetNonce can set the nonce of the account.
func (acc *accountInfo) SetNonce(n uint64) {
	acc.nonce = n
	acc.isNonceSynced = true
}

// GetNonce can return the nonce of the account.
func (acc *accountInfo) GetNonce() uint64 {
	return acc.nonce
}

// IncNonce can increase the nonce of the account.
func (acc *accountInfo) IncNonce() {
	acc.nonce++
}
