// Modifications Copyright 2018 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
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
// This file is derived from core/state_transition.go (2018/06/04).
// Modified and improved for the klaytn development.

package blockchain

import (
	"errors"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/kerrors"
)

var (
	errInsufficientBalanceForGas         = errors.New("insufficient balance of the sender to pay for gas")
	errInsufficientBalanceForGasFeePayer = errors.New("insufficient balance of the fee payer to pay for gas")
	errNotProgramAccount                 = errors.New("not a program account")
	errAccountAlreadyExists              = errors.New("account already exists")
	errMsgToNil                          = errors.New("msg.To() is nil")
	errInvalidCodeFormat                 = errors.New("smart contract code format is invalid")
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==
  4a) Attempt to run transaction data
  4b) If valid, use result as code for the new state object
== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	msg        Message
	gas        uint64
	gasPrice   *big.Int
	gasTipCap  *big.Int
	gasFeeCap  *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      vm.StateDB
	evm        *vm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	// ValidatedSender returns the sender of the transaction.
	// The returned sender should be derived by calling AsMessageAccountKeyPicker().
	ValidatedSender() common.Address

	// ValidatedFeePayer returns the fee payer of the transaction.
	// The returned fee payer should be derived by calling AsMessageAccountKeyPicker().
	ValidatedFeePayer() common.Address

	// ValidatedIntrinsicGas returns the intrinsic gas of the transaction.
	// The returned intrinsic gas should be derived by calling AsMessageAccountKeyPicker().
	ValidatedIntrinsicGas() uint64

	// FeeRatio returns a ratio of tx fee paid by the fee payer in percentage.
	// For example, if it is 30, 30% of tx fee will be paid by the fee payer.
	// 70% will be paid by the sender.
	FeeRatio() (types.FeeRatio, bool)

	//FromFrontier() (common.Address, error)
	To() *common.Address

	Hash() common.Hash

	GasPrice() *big.Int

	// For TxTypeEthereumDynamicFee
	GasTipCap() *big.Int
	GasFeeCap() *big.Int
	EffectiveGasTip(baseFee *big.Int) *big.Int
	EffectiveGasPrice(baseFee *big.Int) *big.Int

	Gas() uint64
	Value() *big.Int

	Nonce() uint64
	CheckNonce() bool
	Data() []byte

	// IntrinsicGas returns `intrinsic gas` based on the tx type.
	// This value is used to differentiate tx fee based on the tx type.
	IntrinsicGas(currentBlockNumber uint64) (uint64, error)

	// Type returns the transaction type of the message.
	Type() types.TxType

	// Validate performs additional validation for each transaction type
	Validate(stateDB types.StateDB, currentBlockNumber uint64) error

	// Execute performs execution of the transaction according to the transaction type.
	Execute(vm types.VM, stateDB types.StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) ([]byte, uint64, error)
}

// TODO-Klaytn Later we can merge Err and Status into one uniform error.
//         This might require changing overall error handling mechanism in Klaytn.
// Klaytn error type
// - Status: Indicate status of transaction after execution.
//           This value will be stored in Receipt if Receipt is available.
//           Please see getReceiptStatusFromErrTxFailed() how this value is calculated.
type kerror struct {
	ErrTxInvalid error
	Status       uint
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message) *StateTransition {
	effectiveGasPrice := msg.EffectiveGasPrice(evm.BaseFee)

	return &StateTransition{
		evm:       evm,
		msg:       msg,
		gasPrice:  effectiveGasPrice,
		gasFeeCap: msg.GasFeeCap(),
		gasTipCap: msg.GasTipCap(),
		value:     msg.Value(),
		data:      msg.Data(),
		state:     evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg Message) ([]byte, uint64, kerror) {
	return NewStateTransition(evm, msg).TransitionDb()
}

// to returns the recipient of the message.
func (st *StateTransition) to() common.Address {
	if st.msg == nil || st.msg.To() == nil /* contract creation */ {
		return common.Address{}
	}
	return *st.msg.To()
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gas < amount {
		return kerrors.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.msg.Gas()), st.gasPrice)

	validatedFeePayer := st.msg.ValidatedFeePayer()
	validatedSender := st.msg.ValidatedSender()
	feeRatio, isRatioTx := st.msg.FeeRatio()
	if isRatioTx {
		feePayerFee, senderFee := types.CalcFeeWithRatio(feeRatio, mgval)

		if st.state.GetBalance(validatedFeePayer).Cmp(feePayerFee) < 0 {
			logger.Debug(errInsufficientBalanceForGasFeePayer.Error(), "feePayer", validatedFeePayer.String(),
				"feePayerBalance", st.state.GetBalance(validatedFeePayer).Uint64(), "feePayerFee", feePayerFee.Uint64(),
				"txHash", st.msg.Hash().String())
			return errInsufficientBalanceForGasFeePayer
		}

		if st.state.GetBalance(validatedSender).Cmp(senderFee) < 0 {
			logger.Debug(errInsufficientBalanceForGas.Error(), "sender", validatedSender.String(),
				"senderBalance", st.state.GetBalance(validatedSender).Uint64(), "senderFee", senderFee.Uint64(),
				"txHash", st.msg.Hash().String())
			return errInsufficientBalanceForGas
		}

		st.state.SubBalance(validatedFeePayer, feePayerFee)
		st.state.SubBalance(validatedSender, senderFee)
	} else {
		// to make a short circuit, process the special case feeRatio == MaxFeeRatio
		if st.state.GetBalance(validatedFeePayer).Cmp(mgval) < 0 {
			logger.Debug(errInsufficientBalanceForGasFeePayer.Error(), "feePayer", validatedFeePayer.String(),
				"feePayerBalance", st.state.GetBalance(validatedFeePayer).Uint64(), "feePayerFee", mgval.Uint64(),
				"txHash", st.msg.Hash().String())
			return errInsufficientBalanceForGasFeePayer
		}

		st.state.SubBalance(validatedFeePayer, mgval)
	}

	st.gas += st.msg.Gas()

	st.initialGas = st.msg.Gas()
	return nil
}

func (st *StateTransition) preCheck() error {
	// Make sure this transaction's nonce is correct.
	if st.msg.CheckNonce() {
		nonce := st.state.GetNonce(st.msg.ValidatedSender())
		if nonce < st.msg.Nonce() {
			logger.Debug(ErrNonceTooHigh.Error(), "account", st.msg.ValidatedSender().String(),
				"accountNonce", nonce, "txNonce", st.msg.Nonce(), "txHash", st.msg.Hash().String())
			return ErrNonceTooHigh
		} else if nonce > st.msg.Nonce() {
			logger.Debug(ErrNonceTooLow.Error(), "account", st.msg.ValidatedSender().String(),
				"accountNonce", nonce, "txNonce", st.msg.Nonce(), "txHash", st.msg.Hash().String())
			return ErrNonceTooLow
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the result including the used gas. It returns an error if failed.
// An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, usedGas uint64, kerr kerror) {
	if st.evm.IsPrefetching() {
		st.gas = st.msg.Gas()
	} else {
		if kerr.ErrTxInvalid = st.preCheck(); kerr.ErrTxInvalid != nil {
			return
		}
	}

	msg := st.msg

	// Pay intrinsic gas.
	if kerr.ErrTxInvalid = st.useGas(msg.ValidatedIntrinsicGas()); kerr.ErrTxInvalid != nil {
		kerr.Status = getReceiptStatusFromErrTxFailed(nil)
		return nil, 0, kerr
	}

	var (
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error and total time limit reached error.
		errTxFailed error
	)

	ret, st.gas, errTxFailed = msg.Execute(st.evm, st.state, st.evm.BlockNumber.Uint64(), st.gas, st.value)

	if errTxFailed != nil {
		logger.Debug("VM returned with error", "err", errTxFailed, "txHash", st.msg.Hash().String())
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		// Another possible errTxFailed could be a time-limit error that happens
		// when the EVM is still running while the block proposer's total
		// execution time of txs for a candidate block reached the predefined
		// limit.
		if errTxFailed == vm.ErrInsufficientBalance || errTxFailed == vm.ErrTotalTimeLimitReached {
			kerr.ErrTxInvalid = errTxFailed
			kerr.Status = getReceiptStatusFromErrTxFailed(nil)
			return nil, 0, kerr
		}
	}
	st.refundGas()

	// Defer transferring Tx fee when DeferredTxFee is true
	if st.evm.ChainConfig().Governance == nil || !st.evm.ChainConfig().Governance.DeferredTxFee() {
		effectiveTip := msg.EffectiveGasTip(st.evm.BaseFee)
		st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), effectiveTip))
	}

	kerr.ErrTxInvalid = nil
	kerr.Status = getReceiptStatusFromErrTxFailed(errTxFailed)
	return ret, st.gasUsed(), kerr
}

var errTxFailed2receiptstatus = map[error]uint{
	nil:                                             types.ReceiptStatusSuccessful,
	vm.ErrDepth:                                     types.ReceiptStatusErrDepth,
	vm.ErrContractAddressCollision:                  types.ReceiptStatusErrContractAddressCollision,
	vm.ErrCodeStoreOutOfGas:                         types.ReceiptStatusErrCodeStoreOutOfGas,
	vm.ErrMaxCodeSizeExceeded:                       types.ReceiptStatuserrMaxCodeSizeExceed,
	kerrors.ErrOutOfGas:                             types.ReceiptStatusErrOutOfGas,
	vm.ErrWriteProtection:                           types.ReceiptStatusErrWriteProtection,
	vm.ErrExecutionReverted:                         types.ReceiptStatusErrExecutionReverted,
	vm.ErrOpcodeComputationCostLimitReached:         types.ReceiptStatusErrOpcodeComputationCostLimitReached,
	kerrors.ErrAccountAlreadyExists:                 types.ReceiptStatusErrAddressAlreadyExists,
	kerrors.ErrNotProgramAccount:                    types.ReceiptStatusErrNotAProgramAccount,
	kerrors.ErrNotHumanReadableAddress:              types.ReceiptStatusErrNotHumanReadableAddress,
	kerrors.ErrFeeRatioOutOfRange:                   types.ReceiptStatusErrFeeRatioOutOfRange,
	kerrors.ErrAccountKeyFailNotUpdatable:           types.ReceiptStatusErrAccountKeyFailNotUpdatable,
	kerrors.ErrDifferentAccountKeyType:              types.ReceiptStatusErrDifferentAccountKeyType,
	kerrors.ErrAccountKeyNilUninitializable:         types.ReceiptStatusErrAccountKeyNilUninitializable,
	kerrors.ErrNotOnCurve:                           types.ReceiptStatusErrNotOnCurve,
	kerrors.ErrZeroKeyWeight:                        types.ReceiptStatusErrZeroKeyWeight,
	kerrors.ErrUnserializableKey:                    types.ReceiptStatusErrUnserializableKey,
	kerrors.ErrDuplicatedKey:                        types.ReceiptStatusErrDuplicatedKey,
	kerrors.ErrWeightedSumOverflow:                  types.ReceiptStatusErrWeightedSumOverflow,
	kerrors.ErrUnsatisfiableThreshold:               types.ReceiptStatusErrUnsatisfiableThreshold,
	kerrors.ErrZeroLength:                           types.ReceiptStatusErrZeroLength,
	kerrors.ErrLengthTooLong:                        types.ReceiptStatusErrLengthTooLong,
	kerrors.ErrNestedCompositeType:                  types.ReceiptStatusErrNestedRoleBasedKey,
	kerrors.ErrLegacyTransactionMustBeWithLegacyKey: types.ReceiptStatusErrLegacyTransactionMustBeWithLegacyKey,
	kerrors.ErrDeprecated:                           types.ReceiptStatusErrDeprecated,
	kerrors.ErrNotSupported:                         types.ReceiptStatusErrNotSupported,
	kerrors.ErrInvalidCodeFormat:                    types.ReceiptStatusErrInvalidCodeFormat,
}

var receiptstatus2errTxFailed = map[uint]error{
	types.ReceiptStatusSuccessful:                              nil,
	types.ReceiptStatusErrDefault:                              ErrVMDefault,
	types.ReceiptStatusErrDepth:                                vm.ErrDepth,
	types.ReceiptStatusErrContractAddressCollision:             vm.ErrContractAddressCollision,
	types.ReceiptStatusErrCodeStoreOutOfGas:                    vm.ErrCodeStoreOutOfGas,
	types.ReceiptStatuserrMaxCodeSizeExceed:                    vm.ErrMaxCodeSizeExceeded,
	types.ReceiptStatusErrOutOfGas:                             kerrors.ErrOutOfGas,
	types.ReceiptStatusErrWriteProtection:                      vm.ErrWriteProtection,
	types.ReceiptStatusErrExecutionReverted:                    vm.ErrExecutionReverted,
	types.ReceiptStatusErrOpcodeComputationCostLimitReached:    vm.ErrOpcodeComputationCostLimitReached,
	types.ReceiptStatusErrAddressAlreadyExists:                 kerrors.ErrAccountAlreadyExists,
	types.ReceiptStatusErrNotAProgramAccount:                   kerrors.ErrNotProgramAccount,
	types.ReceiptStatusErrNotHumanReadableAddress:              kerrors.ErrNotHumanReadableAddress,
	types.ReceiptStatusErrFeeRatioOutOfRange:                   kerrors.ErrFeeRatioOutOfRange,
	types.ReceiptStatusErrAccountKeyFailNotUpdatable:           kerrors.ErrAccountKeyFailNotUpdatable,
	types.ReceiptStatusErrDifferentAccountKeyType:              kerrors.ErrDifferentAccountKeyType,
	types.ReceiptStatusErrAccountKeyNilUninitializable:         kerrors.ErrAccountKeyNilUninitializable,
	types.ReceiptStatusErrNotOnCurve:                           kerrors.ErrNotOnCurve,
	types.ReceiptStatusErrZeroKeyWeight:                        kerrors.ErrZeroKeyWeight,
	types.ReceiptStatusErrUnserializableKey:                    kerrors.ErrUnserializableKey,
	types.ReceiptStatusErrDuplicatedKey:                        kerrors.ErrDuplicatedKey,
	types.ReceiptStatusErrWeightedSumOverflow:                  kerrors.ErrWeightedSumOverflow,
	types.ReceiptStatusErrUnsatisfiableThreshold:               kerrors.ErrUnsatisfiableThreshold,
	types.ReceiptStatusErrZeroLength:                           kerrors.ErrZeroLength,
	types.ReceiptStatusErrLengthTooLong:                        kerrors.ErrLengthTooLong,
	types.ReceiptStatusErrNestedRoleBasedKey:                   kerrors.ErrNestedCompositeType,
	types.ReceiptStatusErrLegacyTransactionMustBeWithLegacyKey: kerrors.ErrLegacyTransactionMustBeWithLegacyKey,
	types.ReceiptStatusErrDeprecated:                           kerrors.ErrDeprecated,
	types.ReceiptStatusErrNotSupported:                         kerrors.ErrNotSupported,
	types.ReceiptStatusErrInvalidCodeFormat:                    kerrors.ErrInvalidCodeFormat,
}

// getReceiptStatusFromErrTxFailed returns corresponding ReceiptStatus for VM error.
func getReceiptStatusFromErrTxFailed(errTxFailed error) (status uint) {
	// TODO-Klaytn Add more VM error to ReceiptStatus
	status, ok := errTxFailed2receiptstatus[errTxFailed]
	if !ok {
		// No corresponding receiptStatus available for errTxFailed
		status = types.ReceiptStatusErrDefault
	}

	return
}

// GetVMerrFromReceiptStatus returns VM error according to status of receipt.
func GetVMerrFromReceiptStatus(status uint) (errTxFailed error) {
	errTxFailed, ok := receiptstatus2errTxFailed[status]
	if !ok {
		return ErrInvalidReceiptStatus
	}

	return
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	// Return KLAY for remaining gas, exchanged at the original rate.
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)

	validatedFeePayer := st.msg.ValidatedFeePayer()
	validatedSender := st.msg.ValidatedSender()
	feeRatio, isRatioTx := st.msg.FeeRatio()
	if isRatioTx {
		feePayer, feeSender := types.CalcFeeWithRatio(feeRatio, remaining)

		st.state.AddBalance(validatedFeePayer, feePayer)
		st.state.AddBalance(validatedSender, feeSender)
	} else {
		// To make a short circuit, the below routine processes when feeRatio == 100.
		st.state.AddBalance(validatedFeePayer, remaining)
	}
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gas
}
