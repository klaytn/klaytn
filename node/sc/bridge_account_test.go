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
	"math/big"
	"os"
	"path"
	"testing"
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
	assert.Equal(t, false, bAcc.cAccount.IsLockedAccount())
	assert.Equal(t, false, bAcc.pAccount.IsLockedAccount())

	pPwdFilePath := path.Join(tempDir, "parent_bridge_account", bAcc.pAccount.address.String())
	pPwdStr, err := ioutil.ReadFile(pPwdFilePath)
	assert.NoError(t, err)

	cPwdFilePath := path.Join(tempDir, "child_bridge_account", bAcc.cAccount.address.String())
	cPwdStr, err := ioutil.ReadFile(cPwdFilePath)
	assert.NoError(t, err)

	// Lock Account
	{
		err := bAcc.cAccount.LockAccount()
		assert.NoError(t, err)
		assert.Equal(t, true, bAcc.cAccount.IsLockedAccount())
	}
	{
		err := bAcc.pAccount.LockAccount()
		assert.NoError(t, err)
		assert.Equal(t, true, bAcc.pAccount.IsLockedAccount())
	}

	// Fail to UnLock Account
	{
		err := bAcc.cAccount.UnLockAccount(string(cPwdStr)[3:])
		assert.EqualError(t, err, keystore.ErrDecrypt.Error())
		assert.Equal(t, true, bAcc.cAccount.IsLockedAccount())
	}
	{
		err := bAcc.pAccount.UnLockAccount(string(pPwdStr)[3:])
		assert.EqualError(t, err, keystore.ErrDecrypt.Error())
		assert.Equal(t, true, bAcc.pAccount.IsLockedAccount())
	}

	// Succeed to UnLock Account
	{
		err := bAcc.cAccount.UnLockAccount(string(cPwdStr))
		assert.NoError(t, err)
		assert.Equal(t, false, bAcc.cAccount.IsLockedAccount())
	}
	{
		err := bAcc.pAccount.UnLockAccount(string(pPwdStr))
		assert.NoError(t, err)
		assert.Equal(t, false, bAcc.pAccount.IsLockedAccount())
	}
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
	assert.Equal(t, false, bAcc.cAccount.IsLockedAccount())
	assert.Equal(t, false, bAcc.pAccount.IsLockedAccount())

	bAcc.pAccount.gasPrice = big.NewInt(100)
	bAcc.pAccount.nonce = 10
	bAcc.pAccount.chainID = big.NewInt(200)
	bAcc.pAccount.isNonceSynced = true
	err = bAcc.pAccount.LockAccount()
	assert.NoError(t, err)

	res := bAcc.GetBridgeOperators()
	assert.Equal(t, 2, len(res))

	pRes := (res["parentOperator"]).(map[string]interface{})
	cRes := res["childOperator"].(map[string]interface{})
	assert.Equal(t, 6, len(pRes))
	assert.Equal(t, 6, len(cRes))

	assert.Equal(t, pRes["address"], bAcc.pAccount.address)
	assert.Equal(t, pRes["nonce"], bAcc.pAccount.nonce)
	assert.Equal(t, pRes["chainID"].(*big.Int).String(), bAcc.pAccount.chainID.String())
	assert.Equal(t, pRes["gasPrice"].(*big.Int).String(), bAcc.pAccount.gasPrice.String())
	assert.Equal(t, pRes["isNonceSynced"], bAcc.pAccount.isNonceSynced)
	assert.Equal(t, pRes["isLocked"], bAcc.pAccount.IsLockedAccount())

	assert.Equal(t, cRes["address"], bAcc.cAccount.address)
	assert.Equal(t, cRes["nonce"], bAcc.cAccount.nonce)
	assert.Equal(t, cRes["chainID"].(*big.Int).String(), bAcc.cAccount.chainID.String())
	assert.Equal(t, cRes["gasPrice"].(*big.Int).String(), bAcc.cAccount.gasPrice.String())
	assert.Equal(t, cRes["isNonceSynced"], bAcc.cAccount.isNonceSynced)
	assert.Equal(t, cRes["isLocked"], bAcc.cAccount.IsLockedAccount())
}
