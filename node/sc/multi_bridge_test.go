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
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

type bridgeTestInfo struct {
	acc *bind.TransactOpts
	b   *bridge.Bridge
	sim *backends.SimulatedBackend
}

func prepareMultiBridgeTest(t *testing.T) *bridgeTestInfo {
	accKey, _ := crypto.GenerateKey()
	acc := bind.NewKeyedTransactor(accKey)

	alloc := blockchain.GenesisAlloc{acc.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	chargeAmount := big.NewInt(10000000)
	acc.Value = chargeAmount
	_, tx, b, err := bridge.DeployBridge(acc, sim, false)
	assert.NoError(t, err)
	sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(sim, tx))
	return &bridgeTestInfo{acc, b, sim}
}

// TestRegisterDeregisterOperator checks the following:
// - the specified operator is registered by the contract method RegisterOperator.
// - the specified operator is deregistered by the contract method DeregisterOperator.
func TestRegisterDeregisterOperator(t *testing.T) {
	info := prepareMultiBridgeTest(t)

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.RegisterOperator(opts, info.acc.From)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	isOperator, err := info.b.Operators(nil, info.acc.From)
	assert.NoError(t, err)
	assert.Equal(t, true, isOperator)

	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err = info.b.DeregisterOperator(opts, info.acc.From)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	isOperator, err = info.b.Operators(nil, info.acc.From)
	assert.NoError(t, err)
	assert.Equal(t, false, isOperator)
}

// TestStartStop checks the following:
// - the bridge contract method Start.
// - the bridge contract method Stop.
func TestStartStop(t *testing.T) {
	info := prepareMultiBridgeTest(t)

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.Start(opts, true)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	isRunning, err := info.b.IsRunning(nil)
	assert.NoError(t, err)
	assert.Equal(t, true, isRunning)

	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err = info.b.Start(opts, false)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	isRunning, err = info.b.IsRunning(nil)
	assert.NoError(t, err)
	assert.Equal(t, false, isRunning)
}

// TestSetCounterPartBridge checks the following:
// - the bridge contract method TestSetCounterPartBridge.
func TestSetCounterPartBridge(t *testing.T) {
	info := prepareMultiBridgeTest(t)
	dummy := common.Address{10}

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.SetCounterPartBridge(opts, dummy)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	cBridge, err := info.b.CounterpartBridge(nil)
	assert.NoError(t, err)
	assert.Equal(t, dummy, cBridge)
}
