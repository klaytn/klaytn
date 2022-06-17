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
	"math/big"
	"testing"
	"time"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/contracts/bridge"
	sc_token "github.com/klaytn/klaytn/contracts/sc_erc20"
	sc_nft "github.com/klaytn/klaytn/contracts/sc_erc721"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
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
	const maxAccounts = 4
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
	bAddr, tx, b, err := bridge.DeployBridge(res.accounts[0], res.sim, true)
	res.b = b
	res.bAddr = bAddr
	assert.NoError(t, err)
	res.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(res.sim, tx))

	owner := res.accounts[0]
	for i := 0; i < maxAccounts; i++ {
		acc := res.accounts[i]
		opts := &bind.TransactOpts{From: owner.From, Signer: owner.Signer, GasLimit: DefaultBridgeTxGasLimit}
		_, _ = b.RegisterOperator(opts, acc.From)
		res.sim.Commit()
	}

	res.requestCh = make(chan *bridge.BridgeRequestValueTransfer, 100)
	res.handleCh = make(chan *bridge.BridgeHandleValueTransfer, 100)
	res.requestSub, err = b.WatchRequestValueTransfer(nil, res.requestCh, nil, nil, nil)
	assert.NoError(t, err)
	res.handleSub, err = b.WatchHandleValueTransfer(nil, res.handleCh, nil, nil, nil)
	assert.NoError(t, err)

	return &res
}

// TestRegisterDeregisterOperator checks the following:
// - the specified operator is registered by the contract method RegisterOperator.
// - the specified operator is deregistered by the contract method DeregisterOperator.
func TestRegisterDeregisterOperator(t *testing.T) {
	info := prepareMultiBridgeTest(t)
	defer info.sim.Close()

	testAddrs := []common.Address{{10}, {20}}

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}

	// info.acc.From is duplicated because it is an owner. ignored.
	operatorList, err := info.b.GetOperatorList(nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(operatorList))

	tx, err := info.b.RegisterOperator(opts, testAddrs[0])
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	tx, err = info.b.RegisterOperator(opts, testAddrs[1])
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	// Check operators
	isOperator, err := info.b.Operators(nil, info.acc.From)
	assert.NoError(t, err)
	assert.Equal(t, true, isOperator)

	isOperator, err = info.b.Operators(nil, testAddrs[0])
	assert.NoError(t, err)
	assert.Equal(t, true, isOperator)

	isOperator, err = info.b.Operators(nil, testAddrs[1])
	assert.NoError(t, err)
	assert.Equal(t, true, isOperator)

	operatorList, err = info.b.GetOperatorList(nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(operatorList))

	// Deregister an operator
	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.DeregisterOperator(opts, testAddrs[0])
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	// Check operators
	isOperator, err = info.b.Operators(nil, testAddrs[0])
	assert.NoError(t, err)
	assert.Equal(t, false, isOperator)

	operatorList, err = info.b.GetOperatorList(nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(operatorList))
	assert.Equal(t, info.acc.From, operatorList[0])
	assert.Equal(t, testAddrs[1], operatorList[1])
}

// TestStartStop checks the following:
// - the bridge contract method Start.
// - the bridge contract method Stop.
func TestStartStop(t *testing.T) {
	info := prepareMultiBridgeTest(t)
	defer info.sim.Close()

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.Start(opts, true)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	isRunning, err := info.b.IsRunning(nil)
	assert.NoError(t, err)
	assert.Equal(t, true, isRunning)

	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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
	defer info.sim.Close()

	dummy := common.Address{10}

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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
	defer info.sim.Close()

	dummy1 := common.Address{10}
	dummy2 := common.Address{20}

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.RegisterToken(opts, dummy1, dummy2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	// Try duplicated insertion.
	tx, err = info.b.RegisterToken(opts, dummy1, dummy2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))

	res, err := info.b.RegisteredTokens(nil, dummy1)
	assert.NoError(t, err)
	assert.Equal(t, dummy2, res)

	allowedTokenList, err := info.b.GetRegisteredTokenList(nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(allowedTokenList))
	assert.Equal(t, dummy1, allowedTokenList[0])

	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.DeregisterToken(opts, dummy1)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	res, err = info.b.RegisteredTokens(nil, dummy1)
	assert.NoError(t, err)
	assert.Equal(t, common.Address{0}, res)

	allowedTokenList, err = info.b.GetRegisteredTokenList(nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(allowedTokenList))
}

// TestMultiBridgeKLAYTransfer1 checks the following:
// - successful value transfer with proper transaction counts.
func TestMultiBridgeKLAYTransfer1(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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
			balance, err := info.sim.BalanceAt(nil, to, nil)
			assert.NoError(t, err)
			assert.Equal(t, big.NewInt(int64(transferAmount)).String(), balance.String())
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
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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

// TestMultiBridgeERC20Transfer checks the following:
// - successful value transfer with proper transaction counts.
func TestMultiBridgeERC20Transfer(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	// Deploy token
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	erc20Addr, tx, erc20, err := sc_token.DeployServiceChainToken(opts, info.sim, info.bAddr)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	// Register token
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.RegisterToken(opts, erc20Addr, erc20Addr)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	// Give minter role to bridge contract
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = erc20.AddMinter(opts, info.bAddr)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	amount := int64(100)
	sentBlockNumber := uint64(100000)

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.HandleERC20Transfer(opts, common.Hash{10}, to, to, erc20Addr, big.NewInt(amount), sentNonce, sentBlockNumber, nil)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.HandleERC20Transfer(opts, common.Hash{10}, to, to, erc20Addr, big.NewInt(amount), sentNonce, sentBlockNumber, nil)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	for {
		select {
		case ev := <-info.handleCh:
			assert.Equal(t, nonceOffset, ev.HandleNonce)
			balance, err := erc20.BalanceOf(nil, to)
			assert.NoError(t, err)
			assert.Equal(t, big.NewInt(amount).String(), balance.String())
			return
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			t.Fatal("Contract Event Loop Running Stop by timeout")
		}
	}
}

// TestMultiBridgeERC721Transfer checks the following:
// - successful value transfer with proper transaction counts.
func TestMultiBridgeERC721Transfer(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 1)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	// Deploy token
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	erc721Addr, tx, erc721, err := sc_nft.DeployServiceChainNFT(opts, info.sim, info.bAddr)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	// Register token
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.RegisterToken(opts, erc721Addr, erc721Addr)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	// Give minter role to bridge contract
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = erc721.AddMinter(opts, info.bAddr)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	amount := int64(100)
	sentBlockNumber := uint64(100000)

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.HandleERC721Transfer(opts, common.Hash{10}, to, to, erc721Addr, big.NewInt(amount), sentNonce, sentBlockNumber, "", nil)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	for {
		select {
		case ev := <-info.handleCh:
			assert.Equal(t, nonceOffset, ev.HandleNonce)
			owner, err := erc721.OwnerOf(nil, big.NewInt(amount))
			assert.NoError(t, err)
			assert.Equal(t, to, owner)
			return
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			t.Fatal("Contract Event Loop Running Stop by timeout")
		}
	}
}

// TestMultiBridgeSetKLAYFee checks the following:
// - successfully setting KLAY fee
func TestMultiBridgeSetKLAYFee(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	const confThreshold = uint8(2)
	const fee = 1000
	requestNonce := uint64(0)

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeConfiguration, confThreshold)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetKLAYFee(opts, big.NewInt(fee), requestNonce)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	ret, err := info.b.FeeOfKLAY(nil)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0).String(), ret.String())

	acc = info.accounts[1]
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetKLAYFee(opts, big.NewInt(fee), requestNonce)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	ret, err = info.b.FeeOfKLAY(nil)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(fee).String(), ret.String())
}

// TestMultiBridgeSetERC20Fee checks the following:
// - successfully setting ERC20 fee
func TestMultiBridgeSetERC20Fee(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	const confThreshold = uint8(2)
	const fee = 1000
	tokenAddr := common.Address{10}
	requestNonce := uint64(0)

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeConfiguration, confThreshold)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetERC20Fee(opts, tokenAddr, big.NewInt(fee), requestNonce)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	ret, err := info.b.FeeOfERC20(nil, tokenAddr)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0).String(), ret.String())

	acc = info.accounts[1]
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetERC20Fee(opts, tokenAddr, big.NewInt(fee), requestNonce)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	ret, err = info.b.FeeOfERC20(nil, tokenAddr)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(fee).String(), ret.String())
}

// TestSetFeeReceiver checks the following:
// - the bridge contract method SetFeeReceiver.
func TestSetFeeReceiver(t *testing.T) {
	info := prepareMultiBridgeTest(t)
	defer info.sim.Close()

	dummy := common.Address{10}

	opts := &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetFeeReceiver(opts, dummy)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	cBridge, err := info.b.FeeReceiver(nil)
	assert.NoError(t, err)
	assert.Equal(t, dummy, cBridge)
}

// TestMultiBridgeErrNotOperator1 checks the following:
// - set threshold to 1.
// - non-operator failed to handle value transfer.
func TestMultiBridgeErrNotOperator1(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	to := common.Address{100}
	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	accKey, _ := crypto.GenerateKey()
	acc := bind.NewKeyedTransactor(accKey)

	tx := SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))
}

// TestMultiBridgeErrNotOperator2 checks the following:
// - set threshold to 2.
// - operator succeeds to handle value transfer.
// - non-operator fails to handle value transfer.
func TestMultiBridgeErrNotOperator2(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	to := common.Address{100}
	acc := info.accounts[0]

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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

	accKey, _ := crypto.GenerateKey()
	acc = bind.NewKeyedTransactor(accKey)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))
}

// TestMultiBridgeErrInvalTx checks the following:
// - set threshold to 2.
// - first operator succeeds to handle value transfer.
// - second operator fails to handle value transfer with invalid tx.
func TestMultiBridgeErrInvalTx(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	to := common.Address{100}
	acc := info.accounts[0]

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce+1, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	for {
		select {
		case _ = <-info.handleCh:
			t.Fatal("Failure is expected")
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			return
		}
	}
}

// TestMultiBridgeErrOverSign checks the following:
// - set threshold to 2.
// - two operators succeed to handle value transfer.
// - the last operator fails to handle value transfer because of vote closing.
func TestMultiBridgeErrOverSign(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	to := common.Address{100}
	acc := info.accounts[0]

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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

	acc = info.accounts[2]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))
}

// TestMultiOperatorKLAYTransferDup checks the following:
// - set threshold to 1.
// - an operator succeed to handle value transfer.
// - the operator (same auth) fails to handle value transfer because of vote closing (duplicated).
// - another operator (different auth) fails to handle value transfer because of vote closing (duplicated).
func TestMultiOperatorKLAYTransferDup(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	to := common.Address{100}
	acc := info.accounts[0]

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 1)
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

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))
}

// TestMultiBridgeSetKLAYFeeErrNonce checks the following:
// - failed to set KLAY fee because of the wrong nonce.
func TestMultiBridgeSetKLAYFeeErrNonce(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	const confThreshold = uint8(2)
	const fee = 1000
	requestNonce := uint64(0)

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeConfiguration, confThreshold)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetKLAYFee(opts, big.NewInt(fee), requestNonce)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	ret, err := info.b.FeeOfKLAY(nil)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0).String(), ret.String())

	acc = info.accounts[1]
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetKLAYFee(opts, big.NewInt(fee), requestNonce+1)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetKLAYFee(opts, big.NewInt(fee), requestNonce-1)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Error(t, bind.CheckWaitMined(info.sim, tx))
}

// TestMultiBridgeKLAYTransferNonceJump checks the following:
// - set threshold to 2.
// - first request is not executed yet since one operator does not vote.
// - jump 100 nonce and successfully handle value transfer.
func TestMultiBridgeKLAYTransferNonceJump(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
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

	sentNonce = sentNonce + 100
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
			assert.Equal(t, sentNonce, ev.HandleNonce)
			return
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			t.Fatal("Contract Event Loop Running Stop by timeout")
		}
	}
}

// TestMultiBridgeKLAYTransferParallel checks the following:
// - set threshold to 2.
// - two different value transfers succeed to handle the value transfer.
func TestMultiBridgeKLAYTransferParallel(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	to := common.Address{100}
	acc := info.accounts[0]

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.Nil(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	go func() {
		acc := info.accounts[0]
		tx := SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
		info.sim.Commit()
		assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

		acc = info.accounts[1]
		tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
		info.sim.Commit()
		assert.NoError(t, bind.CheckWaitMined(info.sim, tx))
	}()

	go func() {
		acc := info.accounts[2]
		sentNonce := sentNonce + 20
		tx := SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
		info.sim.Commit()
		assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

		acc = info.accounts[3]
		tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
		info.sim.Commit()
		assert.NoError(t, bind.CheckWaitMined(info.sim, tx))
	}()

	handleCounter := 0
	for {
		select {
		case _ = <-info.handleCh:
			handleCounter++
			if handleCounter == 2 {
				return
			}
		case err := <-info.handleSub.Err():
			t.Fatal("Contract Event Loop Running Stop by sub.Err()", "err", err)
		case <-time.After(timeOut):
			t.Fatal("Contract Event Loop Running Stop by timeout")
		}
	}
}

// TestMultiBridgeKLAYTransferMixConfig1 checks the following:
// - set threshold to 2.
// - the first tx is done.
// - set threshold to 1.
// - the second operator successfully handles the value transfer.
func TestMultiBridgeKLAYTransferMixConfig1(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 1)
	assert.NoError(t, err)
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

// TestMultiBridgeKLAYTransferMixConfig2 checks the following:
// - set threshold to 2.
// - The first tx is done.
// - set threshold to 3.
// - remain operators successfully handle the value transfer.
func TestMultiBridgeKLAYTransferMixConfig2(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 3)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[2]
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

// TestMultiBridgeKLAYTransferMixConfig1 checks the following:
// - set threshold to 2.
// - the first tx is done.
// - set threshold to 1.
// - the first operator successfully handles the value transfer if retry.
func TestMultiBridgeKLAYTransferMixConfig3(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	nonceOffset := uint64(17)
	sentNonce := nonceOffset
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	opts = &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err = info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 1)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

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

// TestNoncesAndBlockNumber checks the following:
// - set threshold to 2.
// - the first transfer is done and check nonces and block number (req 0, blk 100000)
// - the second transfer is done and check nonces and block number (req 1, blk 100100)
func TestNoncesAndBlockNumber(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	opts := &bind.TransactOpts{From: acc.From, Signer: acc.Signer, GasLimit: DefaultBridgeTxGasLimit}
	tx, err := info.b.SetOperatorThreshold(opts, voteTypeValueTransfer, 2)
	assert.NoError(t, err)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	// First value transfer.
	sentNonce := uint64(0)
	transferAmount := uint64(100)
	sentBlockNumber := uint64(100000)

	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	hMaxNonce, err := info.b.UpperHandleNonce(nil)
	assert.NoError(t, err)
	assert.Equal(t, sentNonce, hMaxNonce)

	lowerHandleNonce, err := info.b.LowerHandleNonce(nil)
	assert.NoError(t, err)
	assert.Equal(t, sentNonce+1, lowerHandleNonce)

	hBlkNum, err := info.b.RecoveryBlockNumber(nil)
	assert.NoError(t, err)
	assert.Equal(t, sentBlockNumber, hBlkNum)

	// Second value transfer.
	sentNonce++
	sentBlockNumber = sentBlockNumber + 100
	acc = info.accounts[0]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	acc = info.accounts[1]
	tx = SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
	info.sim.Commit()
	assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

	hMaxNonce, err = info.b.UpperHandleNonce(nil)
	assert.NoError(t, err)
	assert.Equal(t, sentNonce, hMaxNonce)

	lowerHandleNonce, err = info.b.LowerHandleNonce(nil)
	assert.NoError(t, err)
	assert.Equal(t, sentNonce+1, lowerHandleNonce)

	hBlkNum, err = info.b.RecoveryBlockNumber(nil)
	assert.NoError(t, err)
	assert.Equal(t, sentBlockNumber, hBlkNum)
}

// TestNoncesAndBlockNumberUnordered checks the following:
// - default threshold (1).
// - check if reversed request nonce from 2 to 0 success.
// - check if reversed request nonce from 2 to 0 fail.
func TestNoncesAndBlockNumberUnordered(t *testing.T) {
	info := prepareMultiBridgeEventTest(t)
	defer info.sim.Close()

	acc := info.accounts[0]
	to := common.Address{100}

	transferAmount := uint64(100)

	type testParams struct {
		// test input
		requestNonce  uint64
		requestBlkNum uint64
		// expected result
		maxHandledRequestedNonce uint64
		minUnhandledRequestNonce uint64
		recoveryBlockNumber      uint64
	}

	testCases := []testParams{
		{2, 300, 2, 0, 1},
		{1, 200, 2, 0, 1},
		{0, 100, 2, 3, 300},
		{3, 400, 3, 4, 400},
	}

	for i := 0; i < len(testCases); i++ {
		sentNonce := testCases[i].requestNonce
		sentBlockNumber := testCases[i].requestBlkNum
		t.Log("test round", "i", i, "nonce", sentNonce, "blk", sentBlockNumber)
		tx := SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
		info.sim.Commit()
		assert.NoError(t, bind.CheckWaitMined(info.sim, tx))

		hMaxNonce, err := info.b.UpperHandleNonce(nil)
		assert.NoError(t, err)
		assert.Equal(t, testCases[i].maxHandledRequestedNonce, hMaxNonce)

		minUnhandledNonce, err := info.b.LowerHandleNonce(nil)
		assert.NoError(t, err)
		assert.Equal(t, testCases[i].minUnhandledRequestNonce, minUnhandledNonce)

		blkNum, err := info.b.RecoveryBlockNumber(nil)
		assert.NoError(t, err)
		assert.Equal(t, testCases[i].recoveryBlockNumber, blkNum)
	}

	lowerHandleNonce, _ := info.b.LowerHandleNonce(nil)
	assert.Equal(t, uint64(4), lowerHandleNonce)

	for i := 0; i < len(testCases); i++ {
		sentNonce := testCases[i].requestNonce
		sentBlockNumber := testCases[i].requestBlkNum
		t.Log("test round", "i", i, "nonce", sentNonce, "blk", sentBlockNumber)
		tx := SendHandleKLAYTransfer(info.b, acc, to, transferAmount, sentNonce, sentBlockNumber, t)
		info.sim.Commit()
		assert.Error(t, bind.CheckWaitMined(info.sim, tx))
	}
}
