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
	backend := backends.NewSimulatedBackend(alloc)

	chargeAmount := big.NewInt(10000000)
	acc.Value = chargeAmount
	_, tx, b, err := bridge.DeployBridge(acc, backend, false)
	if err != nil {
		t.Fatalf("fail to DeployBridge %v", err)
	}
	backend.Commit()
	WaitMined(tx, backend, t)
	return &bridgeTestInfo{acc, b, backend}
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
	WaitMined(tx, info.sim, t)

	isOperator, err := info.b.Operators(nil, info.acc.From)
	assert.Equal(t, isOperator, true)
	if err != nil {
		t.Fatal("failed to get Operators.", "err", err)
	}

	opts = &bind.TransactOpts{From: info.acc.From, Signer: info.acc.Signer, GasLimit: gasLimit}
	tx, err = info.b.DeregisterOperator(opts, info.acc.From)
	assert.NoError(t, err)
	info.sim.Commit()
	WaitMined(tx, info.sim, t)

	isOperator, err = info.b.Operators(nil, info.acc.From)
	assert.Equal(t, isOperator, false)
	if err != nil {
		t.Fatal("failed to get Operators.", "err", err)
	}
}
