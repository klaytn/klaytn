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
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"path"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

// TestBridgeAccountLockUnlock checks the lock/unlock functionality.
func TestBridgeAccountLockUnlock(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bAcc, err := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	assert.NoError(t, err)
	assert.Equal(t, true, bAcc.cAccount.IsUnlockedAccount())
	assert.Equal(t, true, bAcc.pAccount.IsUnlockedAccount())

	pPwdFilePath := path.Join(tempDir, "parent_bridge_account", bAcc.pAccount.address.String())
	pPwd, err := ioutil.ReadFile(pPwdFilePath)
	pPwdStr := string(pPwd)
	assert.NoError(t, err)

	cPwdFilePath := path.Join(tempDir, "child_bridge_account", bAcc.cAccount.address.String())
	cPwd, err := ioutil.ReadFile(cPwdFilePath)
	cPwdStr := string(cPwd)
	assert.NoError(t, err)

	lockAccountsWithCheck := func(t *testing.T, bAcc *BridgeAccounts) {
		{
			err := bAcc.cAccount.LockAccount()
			assert.NoError(t, err)
			assert.Equal(t, false, bAcc.cAccount.IsUnlockedAccount())
		}
		{
			err := bAcc.pAccount.LockAccount()
			assert.NoError(t, err)
			assert.Equal(t, false, bAcc.pAccount.IsUnlockedAccount())
		}
	}

	unlockParentAccountWithCheck := func(t *testing.T, bAcc *BridgeAccounts, passwd string, duration *uint64, expectedErr error) {
		expectedIsUnLock := false
		if expectedErr == nil {
			expectedIsUnLock = true
		}

		err := bAcc.pAccount.UnLockAccount(passwd, duration)
		assert.Equal(t, expectedErr, err)

		time.Sleep(time.Second)

		if duration == nil || *duration == 0 {
			assert.Equal(t, expectedIsUnLock, bAcc.pAccount.IsUnlockedAccount())
			return
		}

		time.Sleep(time.Duration(*duration) * time.Second)
		assert.Equal(t, false, bAcc.pAccount.IsUnlockedAccount())
	}

	unlockChildAccountWithCheck := func(t *testing.T, bAcc *BridgeAccounts, passwd string, duration *uint64, expectedErr error) {
		expectedIsUnLock := false
		if expectedErr == nil {
			expectedIsUnLock = true
		}

		err := bAcc.cAccount.UnLockAccount(passwd, duration)
		assert.Equal(t, expectedErr, err)

		time.Sleep(time.Second)

		if duration == nil || *duration == 0 {
			assert.Equal(t, expectedIsUnLock, bAcc.cAccount.IsUnlockedAccount())
			return
		}

		time.Sleep(time.Duration(*duration) * time.Second)
		assert.Equal(t, false, bAcc.cAccount.IsUnlockedAccount())
	}

	// Double time Lock Account
	lockAccountsWithCheck(t, bAcc)
	lockAccountsWithCheck(t, bAcc)

	// Fail to UnLock Account with invalid timeout
	lockAccountsWithCheck(t, bAcc)
	duration := uint64(time.Duration(math.MaxInt64)/time.Second) + 1
	unlockParentAccountWithCheck(t, bAcc, pPwdStr, &duration, errUnlockDurationTooLarge)
	unlockChildAccountWithCheck(t, bAcc, cPwdStr, &duration, errUnlockDurationTooLarge)

	// Fail to UnLock Account with wrong password
	lockAccountsWithCheck(t, bAcc)
	duration = uint64(0)
	unlockParentAccountWithCheck(t, bAcc, pPwdStr[:3], &duration, keystore.ErrDecrypt)
	unlockChildAccountWithCheck(t, bAcc, cPwdStr[:3], &duration, keystore.ErrDecrypt)

	// Succeed to UnLock Account
	lockAccountsWithCheck(t, bAcc)
	duration = uint64(0)
	unlockParentAccountWithCheck(t, bAcc, pPwdStr, &duration, nil)
	unlockChildAccountWithCheck(t, bAcc, cPwdStr, &duration, nil)

	// Succeed to UnLock Account with timeout
	lockAccountsWithCheck(t, bAcc)
	duration = uint64(5)
	unlockParentAccountWithCheck(t, bAcc, pPwdStr, &duration, nil)
	unlockChildAccountWithCheck(t, bAcc, cPwdStr, &duration, nil)

	// Fail to UnLock Account with wrong password
	lockAccountsWithCheck(t, bAcc)
	duration = uint64(5)
	unlockParentAccountWithCheck(t, bAcc, pPwdStr[:3], &duration, keystore.ErrDecrypt)
	unlockChildAccountWithCheck(t, bAcc, cPwdStr[:3], &duration, keystore.ErrDecrypt)

	// Succeed to UnLock Account with nil timeout
	lockAccountsWithCheck(t, bAcc)
	unlockParentAccountWithCheck(t, bAcc, pPwdStr, nil, nil)
	unlockChildAccountWithCheck(t, bAcc, cPwdStr, nil, nil)
}

// TestBridgeAccountInformation checks if the information result is right or not.
func TestBridgeAccountInformation(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), "sc")
	assert.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bAcc, err := NewBridgeAccounts(nil, config.DataDir, database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB}), DefaultBridgeTxGasLimit, DefaultBridgeTxGasLimit)
	assert.NoError(t, err)
	assert.Equal(t, true, bAcc.cAccount.IsUnlockedAccount())
	assert.Equal(t, true, bAcc.pAccount.IsUnlockedAccount())

	bAcc.pAccount.gasPrice = big.NewInt(100)
	bAcc.pAccount.nonce = 10
	bAcc.pAccount.chainID = big.NewInt(200)
	bAcc.pAccount.isNonceSynced = true
	err = bAcc.pAccount.LockAccount()
	assert.NoError(t, err)

	res := bAcc.GetBridgeOperators()
	assert.Equal(t, 2, len(res))

	pRes := res["parentOperator"].(map[string]interface{})
	cRes := res["childOperator"].(map[string]interface{})
	assert.Equal(t, 6, len(pRes))
	assert.Equal(t, 6, len(cRes))

	assert.Equal(t, pRes["address"], bAcc.pAccount.address)
	assert.Equal(t, pRes["nonce"], bAcc.pAccount.nonce)
	assert.Equal(t, pRes["chainID"].(*big.Int).String(), bAcc.pAccount.chainID.String())
	assert.Equal(t, pRes["gasPrice"].(*big.Int).String(), bAcc.pAccount.gasPrice.String())
	assert.Equal(t, pRes["isNonceSynced"], bAcc.pAccount.isNonceSynced)
	assert.Equal(t, pRes["isUnlocked"], bAcc.pAccount.IsUnlockedAccount())

	assert.Equal(t, cRes["address"], bAcc.cAccount.address)
	assert.Equal(t, cRes["nonce"], bAcc.cAccount.nonce)
	assert.Equal(t, cRes["chainID"].(*big.Int).String(), bAcc.cAccount.chainID.String())
	assert.Equal(t, cRes["gasPrice"].(*big.Int).String(), bAcc.cAccount.gasPrice.String())
	assert.Equal(t, cRes["isNonceSynced"], bAcc.cAccount.isNonceSynced)
	assert.Equal(t, cRes["isUnlocked"], bAcc.cAccount.IsUnlockedAccount())
}
