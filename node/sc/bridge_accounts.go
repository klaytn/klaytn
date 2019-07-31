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
	"github.com/klaytn/klaytn/accounts"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/cmd/homi/setup"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"sync"
)

const (
	DefaultBridgeTxGasLimit = 5000000
)

// accountInfo has bridge account's information to make and sign a transaction.
type accountInfo struct {
	wallet   accounts.Wallet
	address  common.Address
	nonce    uint64
	chainID  *big.Int
	gasPrice *big.Int

	isNonceSynced bool
	mu            sync.RWMutex
}

// BridgeAccounts manages bridge account for main/service chain.
type BridgeAccounts struct {
	pAccount *accountInfo
	cAccount *accountInfo
}

// NewBridgeAccounts returns bridgeAccounts created by main/service bridge account keys.
func NewBridgeAccounts(dataDir string) (*BridgeAccounts, error) {
	pWallet, pAccAddr, isLock, err := InitializeBridgeAccountKeystore(path.Join(dataDir, "parent_bridge_account"))
	if err != nil {
		return nil, err
	}

	if isLock {
		logger.Warn("parent_bridge_account is locked. Please unlock the account manually for Service Chain")
	}

	cWallet, cAccAddr, isLock, err := InitializeBridgeAccountKeystore(path.Join(dataDir, "child_bridge_account"))
	if err != nil {
		return nil, err
	}

	if isLock {
		logger.Warn("child_bridge_account is locked. Please unlock the account manually for Service Chain")
	}

	logger.Info("bridge account is loaded", "parent", pAccAddr.String(), "child", cAccAddr.String())

	pAccInfo := &accountInfo{
		wallet:   pWallet,
		address:  pAccAddr,
		nonce:    0,
		chainID:  nil,
		gasPrice: nil,
	}

	cAccInfo := &accountInfo{
		wallet:   cWallet,
		address:  cAccAddr,
		nonce:    0,
		chainID:  nil,
		gasPrice: nil,
	}

	return &BridgeAccounts{
		pAccount: pAccInfo,
		cAccount: cAccInfo,
	}, nil
}

// InitializeBridgeAccountKeystore initializes a keystore, imports existing keys, and tries to unlock the bridge account.
// This returns the 1st account of the wallet, its address, the lock status and the error.
func InitializeBridgeAccountKeystore(keystorePath string) (accounts.Wallet, common.Address, bool, error) {
	ks := keystore.NewKeyStore(keystorePath, keystore.StandardScryptN, keystore.StandardScryptP)

	// If there is no keystore file, this creates a random account and the corresponded password file.
	// TODO-Klaytn-Servicechain A test-option will be added and this routine will be only executed with it.
	if len(ks.Accounts()) == 0 {
		password := setup.RandStringRunes(params.PasswordLength)
		acc, err := ks.NewAccount(password)
		if err != nil {
			return nil, common.Address{}, true, err
		}
		setup.WriteFile([]byte(password), keystorePath, acc.Address.String())

		if err := ks.Unlock(acc, password); err != nil {
			logger.Error("bridge account wallet unlock is failed by created password file.", "address", acc.Address, "err", err)
			os.RemoveAll(keystorePath)
			return nil, common.Address{}, true, err
		}

		return ks.Wallets()[0], acc.Address, false, nil
	}

	// Try to unlock 1st account if valid password file exist. (optional behavior)
	// If unlocking failed, user should unlock it through API.
	acc := ks.Accounts()[0]
	pwdFilePath := path.Join(keystorePath, acc.Address.String())
	pwdStr, err := ioutil.ReadFile(pwdFilePath)
	if err == nil {
		if err := ks.Unlock(acc, string(pwdStr)); err != nil {
			logger.Warn("bridge account wallet unlock is failed by exist password file.", "address", acc.Address, "err", err)
			return ks.Wallets()[0], acc.Address, true, nil
		}
		return ks.Wallets()[0], acc.Address, false, nil
	}

	return ks.Wallets()[0], acc.Address, true, nil
}

// GetTransactOpts returns a transactOpts for transact on local/remote backend.
func (acc *accountInfo) GetTransactOpts() *bind.TransactOpts {
	var nonce *big.Int

	// Only for unit test, if the nonce is not synced yet, return transaction option with nil nonce.
	// Backend will use state nonce.
	if acc.isNonceSynced {
		nonce = new(big.Int).SetUint64(acc.nonce)
	}
	return bind.MakeTransactOptsWithKeystore(acc.wallet, acc.address, nonce, acc.chainID, DefaultBridgeTxGasLimit, acc.gasPrice)
}

// SignTx signs a transaction with the accountInfo.
func (acc *accountInfo) SignTx(tx *types.Transaction) (*types.Transaction, error) {
	return acc.wallet.SignTx(accounts.Account{Address: acc.address}, tx, acc.chainID)
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
