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
	"time"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
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

// ProcessStats includes the time statistics regarding StateProcessor.Process.
type ProcessStats struct {
	BeforeApplyTxs time.Time
	AfterApplyTxs  time.Time
	AfterFinalize  time.Time
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
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, []*vm.InternalTxTrace, ProcessStats, error) {
	var (
		receipts         types.Receipts
		usedGas          = new(uint64)
		header           = block.Header()
		allLogs          []*types.Log
		internalTxTraces []*vm.InternalTxTrace
		processStats     ProcessStats
	)

	// Enable the opcode computation cost limit
	cfg.UseOpcodeComputationCost = true

	// Extract author from the header
	author, _ := p.bc.Engine().Author(header) // Ignore error, we're past header validation

	processStats.BeforeApplyTxs = time.Now()
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, _, internalTxTrace, err := p.bc.ApplyTransaction(p.config, &author, statedb, header, tx, usedGas, &cfg)
		if err != nil {
			return nil, nil, 0, nil, processStats, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)
		internalTxTraces = append(internalTxTraces, internalTxTrace)
	}
	processStats.AfterApplyTxs = time.Now()

	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	if _, err := p.engine.Finalize(p.bc, header, statedb, block.Transactions(), receipts); err != nil {
		return nil, nil, 0, nil, processStats, err
	}
	processStats.AfterFinalize = time.Now()

	return receipts, allLogs, *usedGas, internalTxTraces, processStats, nil
}
