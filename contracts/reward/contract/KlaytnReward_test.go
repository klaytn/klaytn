// Copyright 2018 The klaytn Authors
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

package contract

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/stretchr/testify/assert"
)

func TestSmartContract(t *testing.T) {
	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	initialValue := big.NewInt(params.KLAY)
	withdrawValue := big.NewInt(500000000)

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: initialValue}, auth2.From: {Balance: initialValue}}
	sim := backends.NewSimulatedBackend(alloc)
	defer sim.Close()

	// Deploy a token contract on the simulated blockchain
	_, _, reward, err := DeployKlaytnReward(auth, sim)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	sim.Commit()

	// Set reward
	tx, err := reward.Reward(&bind.TransactOpts{From: auth.From, Signer: auth.Signer,
		Value: withdrawValue}, auth2.From)
	if err != nil {
		log.Fatalf("Failed to call reward : %v", err)
	}
	sim.Commit()
	assert.Equal(t, uint(1), sim.BlockChain().GetReceiptByTxHash(tx.Hash()).Status)

	// Check the balance before withdrawal
	balance1, err := sim.BalanceAt(context.Background(), auth2.From, nil)
	assert.Nil(t, err)
	assert.Equal(t, initialValue, balance1)

	// Withdraw reward
	tx2, err := reward.SafeWithdrawal(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer,
		Value: big.NewInt(0)})
	if err != nil {
		log.Fatalf("Failed to call reward : %v", err)
	}
	sim.Commit()
	assert.Equal(t, uint(1), sim.BlockChain().GetReceiptByTxHash(tx2.Hash()).Status)

	// Check the balance after withdrawal
	balance2, err := sim.BalanceAt(context.Background(), auth2.From, nil)
	assert.Nil(t, err)
	assert.Equal(t, new(big.Int).Add(balance1, withdrawValue), balance2)
}
