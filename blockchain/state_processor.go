// Modifications Copyright 2018 The klaytn Authors
// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/state_processor.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/params"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Klaytn rules by running
// the transaction messages using the statedb and applying any rewards to the processor.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
	)

	// Enable the opcode computation cost limit
	cfg.UseOpcodeComputationCost = true

	// Extract author from the header
	author, _ := p.bc.Engine().Author(header) // Ignore error, we're past header validation

	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, _, err := ApplyTransaction(p.config, p.bc, &author, statedb, header, tx, usedGas, &cfg)
		if err != nil {
			return nil, nil, 0, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
	}

	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	if _, err := p.engine.Finalize(p.bc, header, statedb, block.Transactions(), receipts); err != nil {
		return nil, nil, 0, err
	}

	return receipts, allLogs, *usedGas, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc *BlockChain, author *common.Address, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg *vm.Config) (*types.Receipt, uint64, error) {

	// TODO-Klaytn We reject transactions with unexpected gasPrice and do not put the transaction into TxPool.
	//         And we run transactions regardless of gasPrice if we push transactions in the TxPool.
	/*
		// istanbul BFT
		if tx.GasPrice() != nil && tx.GasPrice().Cmp(common.Big0) > 0 {
			return nil, uint64(0), ErrInvalidGasPrice
		}
	*/

	blockNumber := header.Number.Uint64()

	// validation for each transaction before execution
	if err := tx.Validate(statedb, blockNumber); err != nil {
		return nil, 0, err
	}

	msg, err := tx.AsMessageWithAccountKeyPicker(types.MakeSigner(config, header.Number), statedb, blockNumber)
	if err != nil {
		return nil, 0, err
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)
	// Apply the transaction to the current state (included in the env)
	_, gas, kerr := ApplyMessage(vmenv, msg)
	err = kerr.ErrTxInvalid
	if err != nil {
		return nil, 0, err
	}
	// Update the state with pending changes
	statedb.Finalise(true)
	*usedGas += gas

	receipt := types.NewReceipt(kerr.Status, tx.Hash(), gas)
	// if the transaction created a contract, store the creation address in the receipt.
	msg.FillContractAddress(vmenv.Context.Origin, receipt)
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	return receipt, gas, err
}
