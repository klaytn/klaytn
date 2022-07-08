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

package benchmarks

import (
	"math"
	"math/big"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

type BenchConfig struct {
	ChainConfig *params.ChainConfig
	BlockScore  *big.Int
	Origin      common.Address
	BlockNumber *big.Int
	Time        *big.Int
	GasLimit    uint64
	GasPrice    *big.Int
	Value       *big.Int
	Debug       bool
	EVMConfig   vm.Config

	State     *state.StateDB
	GetHashFn func(n uint64) common.Hash
}

func makeBenchConfig() *BenchConfig {
	cfg := &BenchConfig{}

	cfg.ChainConfig = &params.ChainConfig{ChainID: big.NewInt(1)}
	cfg.BlockScore = new(big.Int)
	// Origin      common.Address
	cfg.BlockNumber = new(big.Int)
	cfg.Time = big.NewInt(time.Now().Unix())
	cfg.GasLimit = math.MaxUint64
	cfg.GasPrice = new(big.Int)
	cfg.Value = new(big.Int)
	// Debug       bool
	// EVMConfig   vm.Config

	memDBManager := database.NewMemoryDBManager()
	cfg.State, _ = state.New(common.Hash{}, state.NewDatabase(memDBManager), nil)
	cfg.GetHashFn = func(n uint64) common.Hash {
		return common.BytesToHash(crypto.Keccak256([]byte(new(big.Int).SetUint64(n).String())))
	}

	return cfg
}

func prepareInterpreterAndContract(code []byte) (*vm.Interpreter, *vm.Contract) {
	// runtime.go:Execute()
	cfg := makeBenchConfig()
	context := vm.Context{
		CanTransfer: blockchain.CanTransfer,
		Transfer:    blockchain.Transfer,
		GetHash:     func(uint64) common.Hash { return common.Hash{} },

		Origin:      cfg.Origin,
		BlockNumber: cfg.BlockNumber,
		Time:        cfg.Time,
		BlockScore:  cfg.BlockScore,
		GasLimit:    cfg.GasLimit,
		GasPrice:    cfg.GasPrice,
	}

	evm := vm.NewEVM(context, cfg.State, cfg.ChainConfig, &cfg.EVMConfig)

	address := common.BytesToAddress([]byte("contract"))
	sender := vm.AccountRef(cfg.Origin)

	cfg.State.CreateSmartContractAccount(address, params.CodeFormatEVM, cfg.ChainConfig.Rules(cfg.BlockNumber))
	cfg.State.SetCode(address, code)

	// Parameters for NewContract()
	caller := sender
	to := vm.AccountRef(address)
	value := cfg.Value
	gas := cfg.GasLimit

	contract := vm.NewContract(caller, to, value, gas)

	contract.SetCallCode(&address, evm.StateDB.GetCodeHash(address), evm.StateDB.GetCode(address))

	return evm.Interpreter(), contract
}
