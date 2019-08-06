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
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
	"time"
)

type bridgeTestInfo struct {
	acc *bind.TransactOpts
	b   *bridge.Bridge
	sim *backends.SimulatedBackend
}

type multiBridgeTestInfo struct {
	accounts   []*bind.TransactOpts
	b          *bridge.Bridge
	bAddr      common.Address
	sim        *backends.SimulatedBackend
	requestCh  chan *bridge.BridgeRequestValueTransfer
	requestSub event.Subscription
	handleCh   chan *bridge.BridgeHandleValueTransfer
	handleSub  event.Subscription
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

func prepareMultiBridgeEventTest(t *testing.T) *multiBridgeTestInfo {
	const maxAccounts = 3
	var res = multiBridgeTestInfo{}

	accountMap := make(map[common.Address]blockchain.GenesisAccount)
	res.accounts = make([]*bind.TransactOpts, maxAccounts)

	for i := 0; i < maxAccounts; i++ {
		accKey, _ := crypto.GenerateKey()
		res.accounts[i] = bind.NewKeyedTransactor(accKey)
		accountMap[res.accounts[i].From] = blockchain.GenesisAccount{Balance: big.NewInt(params.KLAY)}
	}

	res.sim = backends.NewSimulatedBackend(accountMap)

	chargeAmount := big.NewInt(10000000)
	res.accounts[0].Value = chargeAmount
	bAddr, tx, b, err := bridge.DeployBridge(res.accounts[0], res.sim, false)
	res.b = b
	res.bAddr = bAddr
	assert.NoError(t, err)
	res.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(res.sim, tx))

	owner := res.accounts[0]
	for i := 0; i < maxAccounts; i++ {
		acc := res.accounts[i]
		opts := &bind.TransactOpts{From: owner.From, Signer: owner.Signer, GasLimit: gasLimit}
		tx, err := b.RegisterOperator(opts, acc.From)
		assert.NoError(t, err)
		res.sim.Commit()
		assert.Nil(t, bind.CheckWaitMined(res.sim, tx))
	}

	res.requestCh = make(chan *bridge.BridgeRequestValueTransfer, 100)
	res.handleCh = make(chan *bridge.BridgeHandleValueTransfer, 100)
	res.requestSub, err = b.WatchRequestValueTransfer(nil, res.requestCh)
	assert.NoError(t, err)
	res.handleSub, err = b.WatchHandleValueTransfer(nil, res.handleCh)
	assert.NoError(t, err)

	return &res
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

// TestRegisterDeregisterToken checks the following:
// - the bridge contract method RegisterToken and DeregisterToken.
func TestRegisterDeregisterToken(t *testing.T) {
	info := prepareMultiBridgeTest(t)
	dummy1 := common.Address{10}
	dummy2 := common.Address{20}

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.RegisterToken(opts, dummy1, dummy2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	res, err := info.b.AllowedTokens(nil, dummy1)
	assert.NoError(t, err)
	assert.Equal(t, dummy2, res)

	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err = info.b.DeregisterToken(opts, dummy1)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	res, err = info.b.AllowedTokens(nil, dummy1)
	assert.NoError(t, err)
	assert.Equal(t, common.Address{0}, res)
}

// TestOperatorThreshold checks the following:
// - the bridge contract method SetOperatorThreshold
func TestOperatorThreshold(t *testing.T) {
	info := prepareMultiBridgeTest(t)
	const vtThreshold = uint8(128)
	const confThreshold = uint8(255)

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, vtThreshold)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err = info.b.SetOperatorThreshold(opts, voteTypeConfiguration, confThreshold)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	res, err := info.b.OperatorThresholds(nil, voteTypeValueTransfer)
	assert.NoError(t, err)
	assert.Equal(t, vtThreshold, res)

	res, err = info.b.OperatorThresholds(nil, voteTypeConfiguration)
	assert.NoError(t, err)
	assert.Equal(t, confThreshold, res)
}

// TestMultiBridgeKLAYTransfer1 checks the following:
// - successful value transfer with proper transaction counts.
func TestMultiBridgeKLAYTransfer1(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	for {
		select {
		case ev := <-info.handleCh:
			assert.Equal(t, nonceOffset, ev.HandleNonce)
			return
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			t.Fatal("Contract Event Loop Running Stop by timeout")
		}
	}
}

// TestMultiBridgeKLAYTransfer2 checks the following:
// - failed value transfer without proper transaction counts.
// - timeout is expected since operator threshold will not be satisfied.
func TestMultiBridgeKLAYTransfer2(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	for {
		select {
		case _ = <-info.handleCh:
			t.Fatal("unexpected handling of value transfer")
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			// expected timeout
			return
		}
	}
}

// TestMultiBridgeKLAYTransfer3 checks the following:
// - tx is actually handled when operator threshold is satisfied.
// - no double spending is made due to additional tx.
func TestMultiBridgeKLAYTransfer3(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: gasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	for {
		select {
		case _ = <-info.handleCh:
			acc = info.accounts[2]
			tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
			info.sim.Commit()
			assert.Error(t, bind.CheckWaitMined(info.sim, tx))
			return
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			t.Fatal("Contract Event Loop Running Stop by timeout")
		}
	}
}
