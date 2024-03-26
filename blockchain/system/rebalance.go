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

package system

import (
	"context"
	"errors"
	"math/big"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/contracts/system_contracts"
)

type RebalanceCaller interface {
	RebalanceBlockNumber(opts *bind.CallOpts) (*big.Int, error)
	Status(opts *bind.CallOpts) (uint8, error)

	GetZeroedCount(opts *bind.CallOpts) (*big.Int, error)
	Zeroeds(opts *bind.CallOpts, index *big.Int) (common.Address, error)
	GetAllocatedCount(opts *bind.CallOpts) (*big.Int, error)
	Allocateds(opts *bind.CallOpts, arg0 *big.Int) (struct {
		Addr   common.Address
		Amount *big.Int
	}, error)

	CheckZeroedsApproved(opts *bind.CallOpts) error
}

// Kip103ContractCaller implements RebalanceCaller interface for Treasury rebalance contract.
type Kip103ContractCaller struct {
	*system_contracts.TreasuryRebalanceCaller // Assuming TreasuryRebalance is the contract ABI wrapper

	state  *state.StateDB        // the state that is under process
	chain  consensus.ChainReader // chain containing the blockchain information
	header *types.Header         // the header of a new block that is under process
}

// NewKip103ContractCaller creates a new instance of TreasuryRebalanceCaller.
func NewKip103ContractCaller(state *state.StateDB, chain consensus.ChainReader, header *types.Header,
) (*Kip103ContractCaller, error) {
	caller, err := system_contracts.NewTreasuryRebalanceCaller(chain.Config().Kip103ContractAddress,
		&Kip103ContractCaller{state: state, chain: chain, header: header},
	)
	if err != nil {
		return nil, err
	}
	return &Kip103ContractCaller{caller, state, chain, header}, nil
}

func (caller *Kip103ContractCaller) Zeroeds(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	return caller.Retirees(opts, arg0)
}

func (caller *Kip103ContractCaller) GetZeroedCount(opts *bind.CallOpts) (*big.Int, error) {
	return caller.GetRetiredCount(opts)
}

func (caller *Kip103ContractCaller) Allocateds(opts *bind.CallOpts, arg0 *big.Int) (
	struct {
		Addr   common.Address
		Amount *big.Int
	}, error,
) {
	newbie, err := caller.Newbies(opts, arg0)
	return struct {
		Addr   common.Address
		Amount *big.Int
	}{Addr: newbie.Newbie, Amount: newbie.Amount}, err
}

func (caller *Kip103ContractCaller) GetAllocatedCount(opts *bind.CallOpts) (*big.Int, error) {
	return caller.GetNewbieCount(opts)
}

func (caller *Kip103ContractCaller) CheckZeroedsApproved(opts *bind.CallOpts) error {
	return caller.CheckRetiredsApproved(opts)
}

func (caller *Kip103ContractCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return caller.state.GetCode(contract), nil
}

func (caller *Kip103ContractCaller) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	gasPrice := big.NewInt(0) // execute call regardless of the balance of the sender
	gasLimit := uint64(1e8)   // enough gas limit to execute kip103 contract functions
	intrinsicGas := uint64(0) // read operation doesn't require intrinsicGas

	// call.From: zero address will be assigned if nothing is specified
	// call.To: the target contract address will be assigned by `BoundContract`
	// call.Value: nil value is acceptable for `types.NewMessage`
	// call.Data: a proper value will be assigned by `BoundContract`
	// No need to handle acccess list here

	// To fix 0x0 nonce increase, uncomment next line to generate state instance whenever it is called
	// But for backward compatiblity after hardfork, remain this code as commented
	//state, err := caller.chain.State()
	//if err != nil {
	//	return nil, err
	//}
	msg := types.NewMessage(call.From, call.To, caller.state.GetNonce(call.From),
		call.Value, gasLimit, gasPrice, call.Data, false, intrinsicGas, nil)

	blockContext := blockchain.NewEVMBlockContext(caller.header, caller.chain, nil)
	txContext := blockchain.NewEVMTxContext(msg, caller.header)
	txContext.GasPrice = gasPrice                                                                // set gasPrice again if baseFee is assigned
	evm := vm.NewEVM(blockContext, txContext, caller.state, caller.chain.Config(), &vm.Config{}) // no additional vm config required

	result, err := blockchain.ApplyMessage(evm, msg)
	return result.Return(), err
}

type rebalanceResult struct {
	Zeroed    map[common.Address]*big.Int `json:"zeroed"`
	Allocated map[common.Address]*big.Int `json:"allocated"`
	Burnt     *big.Int                    `json:"burnt"`
	Success   bool                        `json:"success"`
}

func newRebalanceReceipt() *rebalanceResult {
	return &rebalanceResult{
		Zeroed:    make(map[common.Address]*big.Int),
		Allocated: make(map[common.Address]*big.Int),
		Burnt:     big.NewInt(0),
		Success:   false,
	}
}

func (result *rebalanceResult) fillZeroed(contract RebalanceCaller, state *state.StateDB) error {
	numRetiredBigInt, err := contract.GetZeroedCount(nil)
	if err != nil {
		logger.Error("Failed to get ZeroedCount from TreasuryRebalance contract", "err", err)
		return err
	}

	for i := 0; i < int(numRetiredBigInt.Int64()); i++ {
		ret, err := contract.Zeroeds(nil, big.NewInt(int64(i)))
		if err != nil {
			logger.Error("Failed to get Zeroeds from TreasuryRebalance contract", "err", err)
			return err
		}
		result.Zeroed[ret] = state.GetBalance(ret)
	}
	return nil
}

func (result *rebalanceResult) fillAllocated(contract RebalanceCaller) error {
	numNewbieBigInt, err := contract.GetAllocatedCount(nil)
	if err != nil {
		logger.Error("Failed to get AllocatedCount from TreasuryRebalance contract", "err", err)
		return nil
	}

	for i := 0; i < int(numNewbieBigInt.Int64()); i++ {
		ret, err := contract.Allocateds(nil, big.NewInt(int64(i)))
		if err != nil {
			logger.Error("Failed to get Allocateds from TreasuryRebalance contract", "err", err)
			return err
		}

		result.Allocated[ret.Addr] = ret.Amount
	}
	return nil
}

func (result *rebalanceResult) totalZeroedBalance() *big.Int {
	total := big.NewInt(0)
	for _, bal := range result.Zeroed {
		total.Add(total, bal)
	}
	return total
}

func (result *rebalanceResult) totalAllocatedBalance() *big.Int {
	total := big.NewInt(0)
	for _, bal := range result.Allocated {
		total.Add(total, bal)
	}
	return total
}

// RebalanceTreasury reads data from a contract, validates stored values, and executes treasury rebalancing (KIP-103, KIP-160).
// It can change the global state by removing old treasury balances and allocating new treasury balances.
// The new allocation can be larger than the removed amount, and the difference between two amounts will be burnt.
func RebalanceTreasury(state *state.StateDB, chain consensus.ChainReader, header *types.Header) (*rebalanceResult, error) {
	var (
		result = newRebalanceReceipt()
		c      = backends.NewBlockchainContractBackend(chain, nil, nil)

		caller RebalanceCaller
		err    error
	)

	if chain.Config().IsDragonForkBlock(header.Number) {
		caller, err = system_contracts.NewTreasuryRebalanceV2Caller(chain.Config().RandaoRegistry.Records[Kip160Name], c)
		if err != nil {
			return nil, err
		}
	} else if chain.Config().IsKIP103ForkBlock(header.Number) {
		caller, err = NewKip103ContractCaller(state, chain, header)
		if err != nil {
			return result, err
		}
	} else {
		return nil, errors.New("unsupported block number for rebalance treasury")
	}

	// Retrieve 1) Get Zeroed
	if err = result.fillZeroed(caller, state); err != nil {
		return result, err
	}

	// Retrieve 2) Get Allocated
	if err = result.fillAllocated(caller); err != nil {
		return result, err
	}

	// Validation 1) Check the target block number
	if blockNum, err := caller.RebalanceBlockNumber(nil); err != nil || blockNum.Cmp(header.Number) != 0 {
		return result, ErrRebalanceIncorrectBlock
	}

	// Validation 2) Check whether status is approved. It should be 2 meaning approved
	if status, err := caller.Status(nil); err != nil || status != 2 {
		return result, ErrRebalanceBadStatus
	}

	// Validation  3) Check approvals from zeroeds
	if err := caller.CheckZeroedsApproved(nil); err != nil {
		return result, err
	}

	// Validation 4) Check the total balance of retirees are bigger than the distributing amount
	totalZeroedAmount := result.totalZeroedBalance()
	totalAllocatedAmount := result.totalAllocatedBalance()
	if chain.Config().IsKIP103ForkBlock(header.Number) && totalZeroedAmount.Cmp(totalAllocatedAmount) < 0 {
		return result, ErrRebalanceNotEnoughBalance
	}

	// Execution 1) Clear all balances of retirees
	for addr := range result.Zeroed {
		state.SetBalance(addr, big.NewInt(0))
	}
	// Execution 2) Distribute KLAY to all newbies
	for addr, balance := range result.Allocated {
		// if newbie has KLAY before the allocation, it will be burnt
		currentBalance := state.GetBalance(addr)
		result.Burnt.Add(result.Burnt, currentBalance)

		state.SetBalance(addr, balance)
	}

	// Fill the remaining fields of the result
	remainder := new(big.Int).Sub(totalZeroedAmount, totalAllocatedAmount)
	result.Burnt.Add(result.Burnt, remainder)
	result.Success = true

	return result, nil
}
