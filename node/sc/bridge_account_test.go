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
	"github.com/klaytn/klaytn/accounts/keystore"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math"
	"math/big"
	"os"
	"path"
	"testing"
	"time"
)

// TestBridgeAccountLockUnlock checks the lock/unlock functionality.
func TestBridgeAccountLockUnlock(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bAcc, err := NewBridgeAccounts(config.DataDir)
	assert.NoError(t, err)
	assert.Equal(t, true, bAcc.cAccount.IsUnlockedAccount())
	assert.Equal(t, true, bAcc.pAccount.IsUnlockedAccount())

	pPwdFilePath := path.Join(tempDir, "parent_bridge_account", bAcc.pAccount.address.String())
	pPwdStr, err := ioutil.ReadFile(pPwdFilePath)
	assert.NoError(t, err)

	cPwdFilePath := path.Join(tempDir, "child_bridge_account", bAcc.cAccount.address.String())
	cPwdStr, err := ioutil.ReadFile(cPwdFilePath)
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

	unlockAccountsWithCheck := func(t *testing.T, bAcc *BridgeAccounts, duration *uint64, expectedErr error, expectedIsUnLock bool) {
		// Because the duration error has higher priority than wrong password error.
		// In normal case, this wrong password case will be checked.
		if expectedErr == nil {
			// Wrong password
			{
				err := bAcc.cAccount.UnLockAccount(duration, string(cPwdStr)[3:])
				assert.Equal(t, keystore.ErrDecrypt, err)
			}
			{
				err := bAcc.pAccount.UnLockAccount(duration, string(pPwdStr)[3:])
				assert.Equal(t, keystore.ErrDecrypt, err)
			}
		}

		// Right password
		{
			err := bAcc.cAccount.UnLockAccount(duration, string(cPwdStr))
			assert.Equal(t, expectedErr, err)
		}
		{
			err := bAcc.pAccount.UnLockAccount(duration, string(pPwdStr))
			assert.Equal(t, expectedErr, err)
		}

		time.Sleep(time.Second)
		if duration == nil || *duration == 0 {
			assert.Equal(t, expectedIsUnLock, bAcc.cAccount.IsUnlockedAccount())
			assert.Equal(t, expectedIsUnLock, bAcc.pAccount.IsUnlockedAccount())
			return
		}

		time.Sleep(time.Duration(*duration) * time.Second)
		assert.Equal(t, false, bAcc.cAccount.IsUnlockedAccount())
		assert.Equal(t, false, bAcc.pAccount.IsUnlockedAccount())
	}

	// Double time Lock Account
	lockAccountsWithCheck(t, bAcc)
	lockAccountsWithCheck(t, bAcc)

	// Fail to UnLock Account with invalid timeout
	lockAccountsWithCheck(t, bAcc)
	duration := uint64(time.Duration(math.MaxInt64)/time.Second) + 1
	unlockAccountsWithCheck(t, bAcc, &duration, errUnlockDurationTooLarge, false)

	// Succeed to UnLock Account
	lockAccountsWithCheck(t, bAcc)
	duration = uint64(0)
	unlockAccountsWithCheck(t, bAcc, &duration, nil, true)

	// Succeed to UnLock Account with timeout
	lockAccountsWithCheck(t, bAcc)
	duration = uint64(5)
	unlockAccountsWithCheck(t, bAcc, &duration, nil, true)

	// Succeed to UnLock Account with default timeout
	lockAccountsWithCheck(t, bAcc)
	unlockAccountsWithCheck(t, bAcc, nil, nil, true)
}

// TestBridgeAccountInformation checks if the information result is right or not.
func TestBridgeAccountInformation(t *testing.T) {
	tempDir := os.TempDir() + "sc"
	os.MkdirAll(tempDir, os.ModePerm)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Fatalf("fail to delete file %v", err)
		}
	}()

	// Config Bridge Account Manager
	config := &SCConfig{}
	config.DataDir = tempDir
	bAcc, err := NewBridgeAccounts(config.DataDir)
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
