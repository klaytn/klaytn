// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
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
// This file is derived from accounts/abi/bind/auth.go (2018/06/04).
// Modified and improved for the klaytn development.

package bind

import (
	"crypto/ecdsa"
	"errors"
	"io"
	"io/ioutil"
	"math/big"

	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
)

// NewTransactor is a utility method to easily create a transaction signer from
// an encrypted json key stream and the associated passphrase.
func NewTransactor(keyin io.Reader, passphrase string) (*TransactOpts, error) {
	json, err := ioutil.ReadAll(keyin)
	if err != nil {
		return nil, err
	}
	key, err := keystore.DecryptKey(json, passphrase)
	if err != nil {
		return nil, err
	}
	return NewKeyedTransactor(key.GetPrivateKey()), nil
}

// NewKeyStoreTransactor is a utility method to easily create a transaction signer from
// an decrypted key from a keystore
func NewKeyStoreTransactor(keystore *keystore.KeyStore, account accounts.Account) (*TransactOpts, error) {
	return &TransactOpts{
		From: account.Address,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}
			signature, err := keystore.SignHash(account, signer.Hash(tx).Bytes())
			if err != nil {
				return nil, err
			}
			return tx.WithSignature(signer, signature)
		},
	}, nil
}

// NewKeyedTransactor is a utility method to easily create a transaction signer
// from a single private key.
func NewKeyedTransactor(key *ecdsa.PrivateKey) *TransactOpts {
	keyAddr := crypto.PubkeyToAddress(key.PublicKey)
	return &TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			return types.SignTx(tx, signer, key)
		},
	}
}

// TODO-klaytn: clef related code
/*
// NewClefTransactor is a utility method to easily create a transaction signer
// with a clef backend.
func NewClefTransactor(clef *external.ExternalSigner, account accounts.Account) *TransactOpts {
	return &TransactOpts{
		From: account.Address,
		Signer: func(signer types.Signer, address common.Address, transaction *types.Transaction) (*types.Transaction, error) {
			if address != account.Address {
				return nil, errors.New("not authorized to sign this account")
			}
			return clef.SignTx(account, transaction, nil) // Clef enforces its own chain id
		},
	}
}
*/

// NewKeyedTransactorWithKeystore is a utility method to easily create a transaction signer
// from a keystore wallet.
func NewKeyedTransactorWithKeystore(address common.Address, ks *keystore.KeyStore, chainID *big.Int) *TransactOpts {
	keyAddr := address
	return &TransactOpts{
		From: keyAddr,
		Signer: func(signer types.Signer, address common.Address, tx *types.Transaction) (*types.Transaction, error) {
			if address != keyAddr {
				return nil, errors.New("not authorized to sign this account")
			}
			account := accounts.Account{Address: address}
			return ks.SignTx(account, tx, chainID)
		},
	}
}

// MakeTransactOpts creates a transaction signer with nonce, gasLimit, and gasPrice from a single private key.
func MakeTransactOpts(accountKey *ecdsa.PrivateKey, nonce *big.Int, gasLimit uint64, gasPrice *big.Int) *TransactOpts {
	if accountKey == nil {
		return nil
	}
	auth := NewKeyedTransactor(accountKey)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	auth.Nonce = nonce
	return auth
}

// MakeTransactOptsWithKeystore creates a transaction signer with nonce, gasLimit, and gasPrice from a keystore wallet.
func MakeTransactOptsWithKeystore(ks *keystore.KeyStore, from common.Address, nonce *big.Int, chainID *big.Int, gasLimit uint64, gasPrice *big.Int) *TransactOpts {
	if ks == nil {
		return nil
	}

	auth := NewKeyedTransactorWithKeystore(from, ks, chainID)
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	auth.Nonce = nonce
	return auth
}
