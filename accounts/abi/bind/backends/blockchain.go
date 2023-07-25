// Copyright 2023 The klaytn Authors
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

package backends

import (
	"context"
	"math/big"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
)

// Maintain separate minimal interfaces of blockchain.BlockChain because ContractBackend are used
// in various situations. BlockChain instances are often passed down as different interfaces such as
// consensus.ChainReader, governance.blockChain, work.BlockChain.
type BlockChainForCaller interface {
	// Required by NewEVMContext
	blockchain.ChainContext

	// Below is a subset of consensus.ChainReader
	// Only using the vocabulary of consensus.ChainReader for potential
	// usability within consensus package.
	Config() *params.ChainConfig
	CurrentHeader() *types.Header
	GetHeaderByNumber(number uint64) *types.Header
	GetBlock(hash common.Hash, number uint64) *types.Block
	State() (*state.StateDB, error)
	StateAt(root common.Hash) (*state.StateDB, error)
	CurrentBlock() *types.Block
}

// BlockchainContractCaller implements bind.ContractCaller, based on
// a user-supplied blockchain.BlockChain instance.
// Its intended purpose is reading system contracts during block processing.
//
// Note that SimulatedBackend creates a new temporary BlockChain for testing,
// whereas BlockchainContractCaller uses an existing BlockChain with existing database.
type BlockchainContractCaller struct {
	bc BlockChainForCaller
}

// This nil assignment ensures at compile time that BlockchainContractCaller implements bind.ContractCaller.
var _ bind.ContractCaller = (*BlockchainContractCaller)(nil)

func NewBlockchainContractCaller(bc BlockChainForCaller) *BlockchainContractCaller {
	return &BlockchainContractCaller{
		bc: bc,
	}
}

func (b *BlockchainContractCaller) CodeAt(ctx context.Context, account common.Address, blockNumber *big.Int) ([]byte, error) {
	if _, state, err := b.getBlockAndState(blockNumber); err != nil {
		return nil, err
	} else {
		return state.GetCode(account), nil
	}
}

// Executes a read-only function call with respect to the specified block's state, or latest state if not specified.
//
// Returns call result in []byte.
// Returns error when:
// - cannot find the corresponding block or stateDB
// - VM revert error
// - VM other errors (e.g. NotProgramAccount, OutOfGas)
// - Error outside VM
func (b *BlockchainContractCaller) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	block, state, err := b.getBlockAndState(blockNumber)
	if err != nil {
		return nil, err
	}

	res, err := b.callContract(call, block, state)
	if err != nil {
		return nil, err
	}
	if len(res.Revert()) > 0 {
		return nil, newRevertError(res)
	}
	return res.Return(), res.Unwrap()
}

func (b *BlockchainContractCaller) callContract(call klaytn.CallMsg, block *types.Block, state *state.StateDB) (*blockchain.ExecutionResult, error) {
	if call.Gas == 0 {
		call.Gas = uint64(3e8) // enough gas for ordinary contract calls
	}

	intrinsicGas, err := types.IntrinsicGas(call.Data, nil, call.To == nil, b.bc.Config().Rules(block.Number()))
	if err != nil {
		return nil, err
	}

	msg := types.NewMessage(call.From, call.To, 0, call.Value, call.Gas, call.GasPrice, call.Data,
		false, intrinsicGas)

	evmContext := blockchain.NewEVMContext(msg, block.Header(), b.bc, nil)
	// EVM demands the sender to have enough KLAY balance (gasPrice * gasLimit) in buyGas()
	// After KIP-71, gasPrice is nonzero baseFee, regardless of the msg.gasPrice (usually 0)
	// But our sender (usually 0x0) won't have enough balance. Instead we override gasPrice = 0 here
	evmContext.GasPrice = big.NewInt(0)
	evm := vm.NewEVM(evmContext, state, b.bc.Config(), &vm.Config{})

	return blockchain.ApplyMessage(evm, msg)
}

func (b *BlockchainContractCaller) getBlockAndState(num *big.Int) (*types.Block, *state.StateDB, error) {
	var block *types.Block
	if num == nil {
		block = b.bc.CurrentBlock()
	} else {
		header := b.bc.GetHeaderByNumber(num.Uint64())
		if header == nil {
			return nil, nil, errBlockDoesNotExist
		}
		block = b.bc.GetBlock(header.Hash(), header.Number.Uint64())
	}

	state, err := b.bc.StateAt(block.Root())
	return block, state, err
}
