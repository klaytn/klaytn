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
	"fmt"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"log"
	"math/big"
	"testing"
)

func TestSmartContract(t *testing.T) {
	// Generate a new random account and a funded simulator
	key, _ := crypto.GenerateKey()
	auth := bind.NewKeyedTransactor(key)

	key2, _ := crypto.GenerateKey()
	auth2 := bind.NewKeyedTransactor(key2)

	alloc := blockchain.GenesisAlloc{auth.From: {Balance: big.NewInt(params.KLAY)}, auth2.From: {Balance: big.NewInt(params.KLAY)}}
	sim := backends.NewSimulatedBackend(alloc)

	// Deploy a token contract on the simulated blockchain
	_, _, reward, err := DeployKlaytnReward(auth, sim)
	if err != nil {
		log.Fatalf("Failed to deploy new token contract: %v", err)
	}
	// Print the current (non existent) and pending name of the contract
	tx, err := reward.Reward(&bind.TransactOpts{From: auth.From, Signer: auth.Signer, Value: big.NewInt(500000000)}, auth2.From)
	if err != nil {
		log.Fatalf("Failed to call reward : %v", err)
	}
	fmt.Println("reward tx.hash :", tx.Hash().Hex())

	// Commit all pending transactions in the simulator and print the names again
	sim.Commit()

	balance, _ := reward.BalanceOf((&bind.CallOpts{Pending: true}), auth2.From)
	fmt.Println("balance :", balance)

	amount, _ := reward.TotalAmount((&bind.CallOpts{Pending: true}))
	fmt.Println("total amount :", amount)

	balance1, _ := sim.BalanceAt(context.Background(), auth2.From, big.NewInt(1))
	fmt.Println("before reward, balance :", balance1)

	tx2, err2 := reward.SafeWithdrawal(&bind.TransactOpts{From: auth2.From, Signer: auth2.Signer, Value: big.NewInt(0)})
	if err2 != nil {
		log.Fatalf("Failed to call reward : %v", err2)
	}
	fmt.Println("reward tx.hash :", tx2.Hash().Hex())

	sim.Commit()

	balance2, _ := sim.BalanceAt(context.Background(), auth2.From, big.NewInt(2))
	fmt.Println("after reward, balance :", balance2)

	if balance1.Cmp(balance2) >= 0 {
		log.Fatalf("Failed to withdraw safely")
	}
}
