package backend

import (
	"context"
	"errors"
	"math/big"

	"github.com/klaytn/klaytn"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/contracts/kip103"
)

// kip103ContractCaller is an implementation of contractCaller only for KIP-103.
// The caller interacts with a KIP-103 contract on a read only basis.
type kip103ContractCaller struct {
	state  *state.StateDB        // the state that is under process
	chain  consensus.ChainReader // chain containing the blockchain information
	header *types.Header         // the header of a new block that is under process
}

func (caller *kip103ContractCaller) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return caller.state.GetCode(contract), nil
}

func (caller *kip103ContractCaller) CallContract(ctx context.Context, call klaytn.CallMsg, blockNumber *big.Int) ([]byte, error) {
	gasLimit := uint64(1e8)   // enough gas limit to execute kip103 contract functions
	intrinsicGas := uint64(0) // read operation doesn't require intrinsicGas

	// call.From: zero address will be assigned if nothing is specified
	// call.To: the target contract address will be assigned by `BoundContract`
	// call.Value: nil value is acceptable for `types.NewMessage`
	// call.Data: a proper value will be assigned by `BoundContract`
	msg := types.NewMessage(call.From, call.To, caller.state.GetNonce(call.From),
		call.Value, gasLimit, caller.header.BaseFee, call.Data, false, intrinsicGas)

	context := blockchain.NewEVMContext(msg, caller.header, caller.chain, nil)
	evm := vm.NewEVM(context, caller.state, caller.chain.Config(), &vm.Config{}) // no additional vm config required

	res, _, kerr := blockchain.ApplyMessage(evm, msg)
	return res, kerr.ErrTxInvalid
}

type kip103Receipt struct {
	Retired map[common.Address]*big.Int `json:"retired"`
	Newbie  map[common.Address]*big.Int `json:"newbie"`
	Burnt   *big.Int                    `json:"burnt"`
	Success bool                        `json:"success"`
}

func newKip103Receipt() *kip103Receipt {
	return &kip103Receipt{
		Retired: make(map[common.Address]*big.Int),
		Newbie:  make(map[common.Address]*big.Int),
		Burnt:   big.NewInt(0),
		Success: false,
	}
}

func (receipt *kip103Receipt) fillRetired(contract *kip103.TreasuryRebalanceCaller, state *state.StateDB) error {
	numRetiredBigInt, err := contract.GetRetiredCount(nil)
	if err != nil {
		logger.Error("Failed to get RetiredCount from TreasuryRebalance contract", "err", err)
		return err
	}

	for i := 0; i < int(numRetiredBigInt.Int64()); i++ {
		ret, err := contract.Retirees(nil, big.NewInt(int64(i)))
		if err != nil {
			logger.Error("Failed to get Retirees from TreasuryRebalance contract", "err", err)
			return err
		}
		receipt.Retired[ret] = state.GetBalance(ret)
	}
	return nil
}

func (receipt *kip103Receipt) fillNewbie(contract *kip103.TreasuryRebalanceCaller) error {
	numNewbieBigInt, err := contract.GetNewbieCount(nil)
	if err != nil {
		logger.Error("Failed to get NewbieCount from TreasuryRebalance contract", "err", err)
		return nil
	}

	for i := 0; i < int(numNewbieBigInt.Int64()); i++ {
		ret, err := contract.Newbies(nil, big.NewInt(int64(i)))
		if err != nil {
			logger.Error("Failed to get Newbies from TreasuryRebalance contract", "err", err)
			return err
		}
		receipt.Newbie[ret.Newbie] = ret.Amount
	}
	return nil
}

func (receipt *kip103Receipt) totalRetriedBalance() *big.Int {
	total := big.NewInt(0)
	for _, bal := range receipt.Retired {
		total.Add(total, bal)
	}
	return total
}

func (receipt *kip103Receipt) totalNewbieBalance() *big.Int {
	total := big.NewInt(0)
	for _, bal := range receipt.Newbie {
		total.Add(total, bal)
	}
	return total
}

// rebalanceTreasury reads data from a contract, validates stored values, and executes treasury rebalancing (KIP-103).
// It can change the global state by removing old treasury balances and allocating new treasury balances.
// The new allocation can be larger than the removed amount, and the difference between two amounts will be burnt.
func rebalanceTreasury(state *state.StateDB, chain consensus.ChainReader, header *types.Header) (*kip103Receipt, error) {
	receipt := newKip103Receipt()
	c := &kip103ContractCaller{state, chain, header}

	// Inside check to avoid a panic case
	if chain.Config().KIP103 == nil {
		return receipt, errors.New("cannot find KIP103 configuration")
	}

	caller, err := kip103.NewTreasuryRebalanceCaller(chain.Config().KIP103.Kip103ContractAddress, c)
	if err != nil {
		return receipt, err
	}

	// Retrieve 1) Get Retired
	if err := receipt.fillRetired(caller, state); err != nil {
		return receipt, err
	}

	// Retrieve 2) Get Newbie
	if err := receipt.fillNewbie(caller); err != nil {
		return receipt, err
	}

	// Validation 1) Check the target block number
	if blockNum, err := caller.RebalanceBlockNumber(nil); err != nil || blockNum.Cmp(header.Number) != 0 {
		return receipt, errors.New("cannot find a proper target block number")
	}

	// Validation 2) Check whether status is approved. It should be 3 meaning approved
	if status, err := caller.Status(nil); err != nil || status != 3 {
		return receipt, errors.New("cannot read a proper status value")
	}

	// Validation 3) Check approvals from retirees
	if err := caller.CheckRetiredsApproved(nil); err != nil {
		return receipt, err
	}

	// Validation 4) Check the total balance of retirees are bigger than the distributing amount
	totalRetiredAmount := receipt.totalRetriedBalance()
	totalNewbieAmount := receipt.totalNewbieBalance()
	if totalRetiredAmount.Cmp(totalNewbieAmount) < 0 {
		return receipt, errors.New("the sum of retired accounts' balance is smaller than the distributing amount")
	}

	// Execution 1) Clear all balances of retirees
	for addr := range receipt.Retired {
		state.SetBalance(addr, big.NewInt(0))
	}
	// Execution 2) Distribute KLAY to all newbies
	for addr, balance := range receipt.Newbie {
		state.SetBalance(addr, balance)
	}

	// Fill the remaining fields of the receipt
	receipt.Burnt.Sub(totalRetiredAmount, totalNewbieAmount)
	receipt.Success = true

	return receipt, nil
}
