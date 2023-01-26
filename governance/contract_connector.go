// Copyright 2022 The klaytn Authors
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

package governance

import (
	"errors"
	"math/big"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	govcontract "github.com/klaytn/klaytn/contracts/gov"
	"github.com/klaytn/klaytn/params"
)

type contractCaller struct {
	chain        blockChain
	contractAddr common.Address
}

var govParamAbi, _ = abi.JSON(strings.NewReader(govcontract.GovParamABI))

const govparamfunc = "getAllParamsAt"

func (c *contractCaller) getAllParamsAt(num *big.Int) (*params.GovParamSet, error) {
	// respect error; some nodes can succeed without error while others do not
	tx, evm, err := c.prepareCall(govParamAbi, govparamfunc, num)
	if err != nil {
		return nil, err
	}

	// ignore error; if error, all nodes will have the same error
	res, err := c.callTx(tx, evm)
	if err != nil {
		return params.NewGovParamSet(), nil
	}

	// ignore error; if error, all nodes will have the same error
	pset, err := c.parseGetAllParamsAt(res)
	if err != nil {
		return params.NewGovParamSet(), nil
	}
	return pset, nil
}

func (c *contractCaller) prepareCall(contractAbi abi.ABI, fn string, args ...interface{}) (*types.Transaction, *vm.EVM, error) {
	tx, err := c.makeTx(contractAbi, fn, args...)
	if err != nil {
		return nil, nil, err
	}

	evm, err := c.makeEVM(tx)
	if err != nil {
		return nil, nil, err
	}

	return tx, evm, nil
}

// Make contract execution transaction
func (c *contractCaller) makeTx(contractAbi abi.ABI, fn string, args ...interface{},
) (*types.Transaction, error) {
	calldata, err := contractAbi.Pack(fn, args...)
	if err != nil {
		logger.Error("Could not pack ABI", "err", err)
		return nil, err
	}

	rules := c.chain.Config().Rules(c.chain.CurrentBlock().Number())
	intrinsicGas, err := types.IntrinsicGas(calldata, nil, false, rules)
	if err != nil {
		logger.Error("Could not fetch intrinsicGas", "err", err)
		return nil, err
	}

	var (
		from       = common.Address{}
		to         = &c.contractAddr
		nonce      = uint64(0)
		amount     = big.NewInt(0)
		gasLimit   = uint64(1e8)
		gasPrice   = big.NewInt(0)
		checkNonce = false
	)

	tx := types.NewMessage(from, to, nonce, amount, gasLimit, gasPrice, calldata,
		checkNonce, intrinsicGas)
	return tx, nil
}

// Make contract execution transaction
func (c *contractCaller) makeEVM(tx *types.Transaction) (*vm.EVM, error) {
	// Load the latest state
	block := c.chain.GetBlockByNumber(c.chain.CurrentBlock().NumberU64())
	if block == nil {
		logger.Error("Could not find the latest block", "num", c.chain.CurrentBlock().NumberU64())
		return nil, errors.New("no block")
	}

	statedb, err := c.chain.StateAt(block.Root())
	if err != nil {
		logger.Error("Could not find the state", "err", err, "num", c.chain.CurrentBlock().NumberU64())
		return nil, err
	}

	// Run EVM at given states
	evmCtx := blockchain.NewEVMContext(tx, block.Header(), c.chain, nil)
	// EVM demands the sender to have enough KLAY balance (gasPrice * gasLimit) in buyGas()
	// After KIP-71, gasPrice is baseFee (=nonzero), regardless of the msg.gasPrice (=zero)
	// But our sender (0x0) won't have enough balance. Instead we override gasPrice = 0 here
	evmCtx.GasPrice = big.NewInt(0)
	evm := vm.NewEVM(evmCtx, statedb, c.chain.Config(), &vm.Config{})
	return evm, nil
}

// Execute contract call at the latest block context
func (c *contractCaller) callTx(tx *types.Transaction, evm *vm.EVM) ([]byte, error) {
	res, _, kerr := blockchain.ApplyMessage(evm, tx)
	if kerr.ErrTxInvalid != nil {
		logger.Warn("Invalid tx")
		return nil, kerr.ErrTxInvalid
	}

	return res, nil
}

func (c *contractCaller) parseGetAllParamsAt(b []byte) (*params.GovParamSet, error) {
	if len(b) == 0 {
		return params.NewGovParamSet(), nil
	}

	var ( // c.f. contracts/gov/GovParam.go:GetAllParamsAt()
		pNames  = new([]string)                   // *[]string = nil
		pValues = new([][]byte)                   // *[][]byte = nil
		out     = &[]interface{}{pNames, pValues} // array of pointers
	)
	if err := govParamAbi.Unpack(out, govparamfunc, b); err != nil {
		return nil, err
	}
	var ( // Retrieve the slices allocated inside Unpack().
		names  = *pNames
		values = *pValues
	)

	// verification
	if len(names) != len(values) {
		logger.Warn("Malformed contract.getAllParams result",
			"len(names)", len(names), "len(values)", len(values))
		return nil, errors.New("malformed contract.getAllParams result")
	}

	bytesMap := make(map[string][]byte)
	for i := 0; i < len(names); i++ {
		bytesMap[names[i]] = values[i]
	}
	return params.NewGovParamSetBytesMap(bytesMap)
}
