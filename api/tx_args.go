// Modifications Copyright 2019 The klaytn Authors
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
// This file is derived from internal/ethapi/api.go (2018/06/04).
// Modified and improved for the klaytn development.

package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

var (
	errTxArgInvalidInputData = errors.New(`Both "data" and "input" are set and not equal. Please use "input" to pass transaction call data.`)
	errTxArgInvalidFeePayer  = errors.New("invalid fee payer is set")
	errTxArgNilTxType        = errors.New("tx should have a type value")
	errTxArgNilContractData  = errors.New(`contract creation without any data provided`)
	errTxArgNilSenderSig     = errors.New("sender signature is not set")
	errTxArgNilNonce         = errors.New("nonce of the sender is not set")
	errTxArgNilGas           = errors.New("gas limit is not set")
	errTxArgNilGasPrice      = errors.New("gas price is not set")
	errNotForFeeDelegationTx = errors.New("fee-delegation type transactions are not allowed to use this API")
)

// isTxField checks whether the string is a field name of the specific txType.
// isTxField[txType][txFieldName] has true/false.
var isTxField = func() map[types.TxType]map[string]bool {
	mapOfFieldMap := map[types.TxType]map[string]bool{}
	internalDataTypes := map[types.TxType]interface{}{
		// since legacy tx has optional fields, some fields can be omitted
		// types.TxTypeLegacyTransaction:                           types.TxInternalDataLegacy{},
		types.TxTypeValueTransfer:                               types.TxInternalDataValueTransfer{},
		types.TxTypeFeeDelegatedValueTransfer:                   types.TxInternalDataFeeDelegatedValueTransfer{},
		types.TxTypeFeeDelegatedValueTransferWithRatio:          types.TxInternalDataFeeDelegatedValueTransferWithRatio{},
		types.TxTypeValueTransferMemo:                           types.TxInternalDataValueTransferMemo{},
		types.TxTypeFeeDelegatedValueTransferMemo:               types.TxInternalDataFeeDelegatedValueTransferMemo{},
		types.TxTypeFeeDelegatedValueTransferMemoWithRatio:      types.TxInternalDataFeeDelegatedValueTransferMemoWithRatio{},
		types.TxTypeAccountUpdate:                               types.TxInternalDataAccountUpdate{},
		types.TxTypeFeeDelegatedAccountUpdate:                   types.TxInternalDataFeeDelegatedAccountUpdate{},
		types.TxTypeFeeDelegatedAccountUpdateWithRatio:          types.TxInternalDataFeeDelegatedAccountUpdateWithRatio{},
		types.TxTypeSmartContractDeploy:                         types.TxInternalDataSmartContractDeploy{},
		types.TxTypeFeeDelegatedSmartContractDeploy:             types.TxInternalDataFeeDelegatedSmartContractDeploy{},
		types.TxTypeFeeDelegatedSmartContractDeployWithRatio:    types.TxInternalDataFeeDelegatedSmartContractDeployWithRatio{},
		types.TxTypeSmartContractExecution:                      types.TxInternalDataSmartContractExecution{},
		types.TxTypeFeeDelegatedSmartContractExecution:          types.TxInternalDataFeeDelegatedSmartContractExecution{},
		types.TxTypeFeeDelegatedSmartContractExecutionWithRatio: types.TxInternalDataFeeDelegatedSmartContractExecutionWithRatio{},
		types.TxTypeCancel:                                      types.TxInternalDataCancel{},
		types.TxTypeFeeDelegatedCancel:                          types.TxInternalDataFeeDelegatedCancel{},
		types.TxTypeFeeDelegatedCancelWithRatio:                 types.TxInternalDataFeeDelegatedCancelWithRatio{},
		types.TxTypeChainDataAnchoring:                          types.TxInternalDataChainDataAnchoring{},
		types.TxTypeFeeDelegatedChainDataAnchoring:              types.TxInternalDataFeeDelegatedChainDataAnchoring{},
		types.TxTypeFeeDelegatedChainDataAnchoringWithRatio:     types.TxInternalDataFeeDelegatedChainDataAnchoringWithRatio{},
	}

	// generate field maps for each tx type
	for txType, internalData := range internalDataTypes {
		fieldMap := map[string]bool{}
		internalDataType := reflect.TypeOf(internalData)

		// key of filedMap is tx field name and value of fieldMap means the existence of field name
		for i := 0; i < internalDataType.NumField(); i++ {
			fieldMap[internalDataType.Field(i).Name] = true
		}

		// additional field of SendTxArgs to support various tx types
		fieldMap["TypeInt"] = true
		// additional field of SendTxArgs to support a legacy tx field (skip checking)
		fieldMap["Data"] = false

		mapOfFieldMap[txType] = fieldMap
	}
	return mapOfFieldMap
}()

type NewTxArgs interface {
	setDefaults(context.Context, Backend) error
	toTransaction() (*types.Transaction, error)
	from() common.Address
	gas() *hexutil.Uint64
	gasPrice() *hexutil.Big
	nonce() *hexutil.Uint64
	setGas(*hexutil.Uint64)
	setGasPrice(*hexutil.Big)
}

// TransactionArgs represents the arguments to construct a new transaction
// into the txpool for eth namespace. For Kaia namespace, use SendTxArgs instead.
// For the message call, both Kaia and Eth namespace use TransactionArgs.
type TransactionArgs struct {
	From         *common.Address `json:"from"`
	Recipient    *common.Address `json:"to"`
	GasLimit     *hexutil.Uint64 `json:"gas"`
	Price        *hexutil.Big    `json:"gasPrice"`
	Amount       *hexutil.Big    `json:"value"`
	AccountNonce *hexutil.Uint64 `json:"nonce"`

	// We accept "data" and "input" for backwards-compatibility reasons.
	// "input" is the newer name and should be preferred by clients.
	Data    *hexutil.Bytes `json:"data"`
	Payload *hexutil.Bytes `json:"input"`

	// For EthereumDynamicFee TxType
	MaxPriorityFeePerGas *hexutil.Big `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big `json:"maxFeePerGas"`

	// For EthereumAccessList TxType
	AccessList *types.AccessList `json:"accessList,omitempty"`
	ChainID    *hexutil.Big      `json:"chainId,omitempty"`
}

func (args *TransactionArgs) from() common.Address {
	if args.From == nil {
		return common.Address{}
	}
	return *args.From
}

func (args *TransactionArgs) gas() *hexutil.Uint64 {
	return args.GasLimit
}

func (args *TransactionArgs) gasPrice() *hexutil.Big {
	return args.Price
}

func (args *TransactionArgs) nonce() *hexutil.Uint64 {
	return args.AccountNonce
}

func (args *TransactionArgs) InputData() []byte {
	if args.Payload != nil {
		return *args.Payload
	}
	if args.Data != nil {
		return *args.Data
	}
	return nil
}

func (args *TransactionArgs) setGas(gas *hexutil.Uint64) {
	args.GasLimit = gas
}

func (args *TransactionArgs) setGasPrice(gasPrice *hexutil.Big) {
	args.Price = gasPrice
}

func (args *TransactionArgs) ToMessage(globalGasCap uint64, baseFee *big.Int, intrinsicGas uint64) (*types.Transaction, error) {
	// Reject invalid combinations of pre- and post-1559 fee styles
	if args.Price != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return nil, errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	} else if args.MaxFeePerGas != nil && args.MaxPriorityFeePerGas != nil {
		if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
			return nil, errors.New("MaxPriorityFeePerGas is greater than MaxFeePerGas")
		}
	}

	// Set sender address or use zero address if none specified.
	addr := args.from()

	// Set default gas & gas price if none were set
	gas := globalGasCap
	if gas == 0 {
		gas = params.UpperGasLimit
	}
	if args.GasLimit != nil {
		gas = uint64(*args.GasLimit)
	}
	if globalGasCap != 0 && globalGasCap < gas {
		logger.Warn("Caller gas above allowance, capping", "requested", gas, "cap", globalGasCap)
		gas = globalGasCap
	}

	// Do not update gasPrice unless any of args.Price and args.MaxFeePerGas is specified.
	gasPrice := new(big.Int)
	if args.Price != nil {
		gasPrice = args.Price.ToInt()
	} else if args.MaxFeePerGas != nil {
		gasPrice = args.MaxFeePerGas.ToInt()
	} else if baseFee.Cmp(new(big.Int).SetUint64(params.ZeroBaseFee)) != 0 {
		// User specified neither GasPrice nor MaxFeePerGas, use baseFee
		gasPrice = new(big.Int).Mul(baseFee, common.Big2)
	}

	value := new(big.Int)
	if &args.Amount != nil {
		value = args.Amount.ToInt()
	}

	var accessList types.AccessList
	if args.AccessList != nil {
		accessList = *args.AccessList
	}
	return types.NewMessage(addr, args.Recipient, 0, value, gas, gasPrice, args.InputData(), false, intrinsicGas, accessList), nil
}

// setDefaults fills in default values for unspecified tx fields.
func (args *TransactionArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.Price != nil && (args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil) {
		return errors.New("both gasPrice and (maxFeePerGas or maxPriorityFeePerGas) specified")
	}
	// After london, default to 1559 uncles gasPrice is set
	head := b.CurrentBlock().Header()
	isMagma := head.BaseFee != nil

	// b.SuggestPrice = unitPrice, for before Magma
	//                = baseFee,   for after Magma
	gasPrice, err := b.SuggestPrice(ctx)
	if err != nil {
		return err
	}

	// If user specifies both maxPriorityFee and maxFee, then we do not
	// need to consult the chain for defaults. It's definitely a London tx.
	if args.MaxPriorityFeePerGas == nil || args.MaxFeePerGas == nil {
		if b.ChainConfig().IsEthTxTypeForkEnabled(head.Number) && args.Price == nil {
			if args.MaxPriorityFeePerGas == nil {
				args.MaxPriorityFeePerGas = (*hexutil.Big)(gasPrice)
			}
			if args.MaxFeePerGas == nil {
				// Before Magma hard fork, `gasFeeCap` was set to `maxPriorityFeePerGas` by default.
				args.MaxFeePerGas = args.MaxPriorityFeePerGas
				if isMagma {
					// After Magma hard fork, `gasFeeCap` was set to `baseFee*2` by default.
					args.MaxFeePerGas = (*hexutil.Big)(gasPrice)
				}
			}
			if isMagma {
				if args.MaxFeePerGas.ToInt().Cmp(new(big.Int).Div(gasPrice, common.Big2)) < 0 {
					return fmt.Errorf("maxFeePerGas (%v) < BaseFee (%v)", args.MaxFeePerGas, gasPrice)
				}
			} else if args.MaxPriorityFeePerGas.ToInt().Cmp(gasPrice) != 0 || args.MaxFeePerGas.ToInt().Cmp(gasPrice) != 0 {
				return fmt.Errorf("only %s is allowed to be used as maxFeePerGas and maxPriorityPerGas", gasPrice.Text(16))
			}
			if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
				return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
			}
		} else {
			if args.MaxFeePerGas != nil || args.MaxPriorityFeePerGas != nil {
				return errors.New("maxFeePerGas or maxPriorityFeePerGas specified but london is not active yet")
			}
			if args.Price == nil {
				// TODO-Kaia: Original logic of Ethereum uses b.SuggestTipCap which suggests TipCap, not a GasPrice.
				// But Kaia currently uses fixed unit price determined by Governance, so using b.SuggestPrice
				// is fine as now.
				if b.ChainConfig().IsEthTxTypeForkEnabled(head.Number) {
					// TODO-Kaia: Kaia is using fixed BaseFee(0) as now but
					// if we apply dynamic BaseFee, we should add calculated BaseFee instead of params.ZeroBaseFee.
					gasPrice.Add(gasPrice, new(big.Int).SetUint64(params.ZeroBaseFee))
				}
				args.Price = (*hexutil.Big)(gasPrice)
			}
		}
	} else {
		// Both maxPriorityFee and maxFee set by caller. Sanity-check their internal relation
		if isMagma {
			if args.MaxFeePerGas.ToInt().Cmp(new(big.Int).Div(gasPrice, common.Big2)) < 0 {
				return fmt.Errorf("maxFeePerGas (%v) < BaseFee (%v)", args.MaxFeePerGas, gasPrice)
			}
		} else {
			if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
				return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
			}
		}
	}
	if args.Amount == nil {
		args.Amount = new(hexutil.Big)
	}
	if args.AccountNonce == nil {
		nonce := b.GetPoolNonce(ctx, args.from())
		args.AccountNonce = (*hexutil.Uint64)(&nonce)
	}
	if args.Data != nil && args.Payload != nil && !bytes.Equal(*args.Data, *args.Payload) {
		return errors.New(`both "data" and "input" are set and not equal. Please use "input" to pass transaction call data`)
	}
	if args.Recipient == nil && len(args.InputData()) == 0 {
		return errors.New(`contract creation without any data provided`)
	}
	// Estimate the gas usage if necessary.
	if args.GasLimit == nil {
		// These fields are immutable during the estimation, safe to
		// pass the pointer directly.
		data := args.InputData()
		callArgs := TransactionArgs{}
		callArgs.From = args.From
		callArgs.Recipient = args.Recipient
		callArgs.Price = args.Price
		callArgs.MaxFeePerGas = args.MaxFeePerGas
		callArgs.MaxPriorityFeePerGas = args.MaxPriorityFeePerGas
		callArgs.Amount = args.Amount
		callArgs.Data = (*hexutil.Bytes)(&data)
		callArgs.AccessList = args.AccessList

		pendingBlockNr := rpc.NewBlockNumberOrHashWithNumber(rpc.PendingBlockNumber)
		gasCap := uint64(0)
		if rpcGasCap := b.RPCGasCap(); rpcGasCap != nil {
			gasCap = rpcGasCap.Uint64()
		}
		estimated, err := EthDoEstimateGas(ctx, b, callArgs, pendingBlockNr, gasCap)
		if err != nil {
			return err
		}
		args.GasLimit = &estimated
		logger.Trace("Estimate gas usage automatically", "gas", args.GasLimit)
	}
	if args.ChainID == nil {
		id := (*hexutil.Big)(b.ChainConfig().ChainID)
		args.ChainID = id
	}
	return nil
}

// toTransaction converts the arguments to a transaction.
// This assumes that setDefaults has been called.
func (args *TransactionArgs) toTransaction() (*types.Transaction, error) {
	var tx *types.Transaction
	switch {
	case args.MaxFeePerGas != nil:
		al := types.AccessList{}
		if args.AccessList != nil {
			al = *args.AccessList
		}
		tx = types.NewTx(&types.TxInternalDataEthereumDynamicFee{
			ChainID:      (*big.Int)(args.ChainID),
			AccountNonce: uint64(*args.AccountNonce),
			GasTipCap:    (*big.Int)(args.MaxPriorityFeePerGas),
			GasFeeCap:    (*big.Int)(args.MaxFeePerGas),
			GasLimit:     uint64(*args.GasLimit),
			Recipient:    args.Recipient,
			Amount:       (*big.Int)(args.Amount),
			Payload:      args.InputData(),
			AccessList:   al,
		})
	case args.AccessList != nil:
		tx = types.NewTx(&types.TxInternalDataEthereumAccessList{
			ChainID:      (*big.Int)(args.ChainID),
			AccountNonce: uint64(*args.AccountNonce),
			Recipient:    args.Recipient,
			GasLimit:     uint64(*args.GasLimit),
			Price:        (*big.Int)(args.Price),
			Amount:       (*big.Int)(args.Amount),
			Payload:      args.InputData(),
			AccessList:   *args.AccessList,
		})
	default:
		tx = types.NewTx(&types.TxInternalDataLegacy{
			AccountNonce: uint64(*args.AccountNonce),
			Price:        (*big.Int)(args.Price),
			GasLimit:     uint64(*args.GasLimit),
			Recipient:    args.Recipient,
			Amount:       (*big.Int)(args.Amount),
			Payload:      args.InputData(),
		})
	}
	return tx, nil
}

// SendTxArgs represents the arguments to submit a new transaction into the transaction pool
// for a kaia namespace.
type SendTxArgs struct {
	TransactionArgs
	TypeInt       *types.TxType          `json:"typeInt,omitempty"`
	FeePayer      *common.Address        `json:"feePayer"`
	FeeRatio      *types.FeeRatio        `json:"feeRatio"`
	CodeFormat    *params.CodeFormat     `json:"codeFormat"`
	HumanReadable *bool                  `json:"humanReadable"`
	Key           *hexutil.Bytes         `json:"key"`
	TxSignatures  types.TxSignaturesJSON `json:"signatures"`
}

// setDefaults is a helper function that fills in default values for unspecified common tx fields.
func (args *SendTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.TypeInt == nil {
		args.TypeInt = new(types.TxType)
		*args.TypeInt = types.TxTypeLegacyTransaction
	}
	if args.GasLimit == nil {
		args.GasLimit = new(hexutil.Uint64)
		*args.GasLimit = hexutil.Uint64(90000)
	}
	// Eth typed transactions requires chainId.
	if args.TypeInt.IsEthTypedTransaction() {
		if args.ChainID == nil {
			args.ChainID = (*hexutil.Big)(b.ChainConfig().ChainID)
		}
	}
	// After london, default to 1559 uncles gasPrice is set
	head := b.CurrentBlock().Header()
	isMagma := head.BaseFee != nil

	// b.SuggestPrice = unitPrice, for before Magma
	//                = baseFee * 2,   for after Magma
	gasPrice, err := b.SuggestPrice(ctx)
	if err != nil {
		return err
	}

	// For the transaction that do not use the gasPrice field, the default value of gasPrice is not set.
	if args.Price == nil && *args.TypeInt != types.TxTypeEthereumDynamicFee {
		args.Price = (*hexutil.Big)(gasPrice)
	}

	if *args.TypeInt == types.TxTypeEthereumDynamicFee {
		if args.MaxPriorityFeePerGas == nil {
			args.MaxPriorityFeePerGas = (*hexutil.Big)(gasPrice)
		}
		if args.MaxFeePerGas == nil {
			// After EthTxtype, `gasFeeCap` was set to `maxPriorityFeePerGas` by default.
			// Anyway, if it's different from unitPrice(gasPrice), rejected.
			args.MaxFeePerGas = args.MaxPriorityFeePerGas

			// After Magma hard fork, `gasFeeCap` was set to `baseFee*2` by default.
			if isMagma {
				args.MaxFeePerGas = (*hexutil.Big)(gasPrice)
			}
		}
		if isMagma {
			if args.MaxFeePerGas.ToInt().Cmp(new(big.Int).Div(gasPrice, common.Big2)) < 0 {
				return fmt.Errorf("maxFeePerGas (%v) < BaseFee (%v)", args.MaxFeePerGas, gasPrice)
			}
		} else if args.MaxPriorityFeePerGas.ToInt().Cmp(gasPrice) != 0 || args.MaxFeePerGas.ToInt().Cmp(gasPrice) != 0 {
			return fmt.Errorf("only %s is allowed to be used as maxFeePerGas and maxPriorityPerGas", gasPrice.Text(16))
		}
		if args.MaxFeePerGas.ToInt().Cmp(args.MaxPriorityFeePerGas.ToInt()) < 0 {
			return fmt.Errorf("maxFeePerGas (%v) < maxPriorityFeePerGas (%v)", args.MaxFeePerGas, args.MaxPriorityFeePerGas)
		}
	}
	if args.AccountNonce == nil {
		nonce := b.GetPoolNonce(ctx, args.from())
		args.AccountNonce = (*hexutil.Uint64)(&nonce)
	}

	return nil
}

// checkArgs checks the validity of SendTxArgs values.
// The each tx types has its own validation logic to give detailed errors to users.
func (args *SendTxArgs) checkArgs() error {
	if args.TypeInt == nil {
		return errTxArgNilTxType
	}
	// Skip ethereum transaction type since it has optional fields
	if args.TypeInt.IsEthereumTransaction() {
		return nil
	}

	checkArg := func(field reflect.StructField, value reflect.Value) error {
		// Skip From since it is an essential field
		// Skip TxSignatures since the value is not considered by all APIs
		if field.Name == "From" || field.Name == "TxSignatures" {
			return nil
		}
		// An args field doesn't have a value but the field name exist on the tx type
		if value.IsNil() && isTxField[*args.TypeInt][field.Name] {
			// if argsValue.Field(i).IsNil() && isTxField[*args.TypeInt][argsType.Field(i).Name] {
			// Allow only contract deploying txs to set the recipient as nil
			if (*args.TypeInt).IsContractDeploy() && field.Name == "To" {
				return nil
			}
			return errors.New((string)(field.Tag) + " is required for " + (*args.TypeInt).String())
		}

		// An args field has a value but the field name doesn't exist on the tx type
		if !value.IsNil() && !isTxField[*args.TypeInt][field.Name] {
			return errors.New((string)(field.Tag) + " is not a field of " + (*args.TypeInt).String())
		}
		return nil
	}

	// check common fields first
	argsType := reflect.TypeOf(args.TransactionArgs)
	argsValue := reflect.ValueOf(args.TransactionArgs)
	for i := 0; i < argsType.NumField(); i++ {
		if err := checkArg(argsType.Field(i), argsValue.Field(i)); err != nil {
			return err
		}
	}

	// then, check kaia-specific fields
	argsType = reflect.TypeOf(*args)
	argsValue = reflect.ValueOf(*args)
	for i := 0; i < argsType.NumField(); i++ {
		if argsType.Field(i).Name == "TransactionArgs" {
			continue
		}
		if err := checkArg(argsType.Field(i), argsValue.Field(i)); err != nil {
			return err
		}
	}

	return nil
}

// genTxValuesMap generates a value map used used in "NewTransactionWithMap" function.
// This function assigned all non-nil values regardless of the tx type.
// Invalid values in the map will be validated in "NewTransactionWithMap" function.
func (args *SendTxArgs) genTxValuesMap() map[types.TxValueKeyType]interface{} {
	values := make(map[types.TxValueKeyType]interface{})

	// common tx fields. They should have values after executing "setDefaults" function.
	if args.TypeInt == nil || args.AccountNonce == nil || args.GasLimit == nil {
		return values
	}
	// GasPrice can be an optional tx filed for TxTypeEthereumDynamicFee
	if args.Price == nil && *args.TypeInt != types.TxTypeEthereumDynamicFee {
		return values
	}

	if !args.TypeInt.IsEthereumTransaction() {
		values[types.TxValueKeyFrom] = args.from()
	}
	values[types.TxValueKeyNonce] = uint64(*args.AccountNonce)
	values[types.TxValueKeyGasLimit] = uint64(*args.GasLimit)

	// optional tx fields
	if args.Price != nil {
		values[types.TxValueKeyGasPrice] = (*big.Int)(args.Price)
	}
	if args.TypeInt.IsContractDeploy() || args.TypeInt.IsEthereumTransaction() {
		// contract deploy type and ethereum tx types allow nil as TxValueKeyTo value
		values[types.TxValueKeyTo] = args.Recipient
	} else if args.Recipient != nil {
		values[types.TxValueKeyTo] = *args.Recipient
	}
	if args.FeePayer != nil {
		values[types.TxValueKeyFeePayer] = *args.FeePayer
	}
	if args.FeeRatio != nil {
		values[types.TxValueKeyFeeRatioOfFeePayer] = *args.FeeRatio
	}
	if args.Amount != nil {
		values[types.TxValueKeyAmount] = (*big.Int)(args.Amount)
	} else if args.TypeInt.IsEthereumTransaction() {
		values[types.TxValueKeyAmount] = common.Big0
	}
	if args.Payload != nil {
		// chain data anchoring type uses the TxValueKeyAnchoredData field
		if args.TypeInt.IsChainDataAnchoring() {
			values[types.TxValueKeyAnchoredData] = ([]byte)(*args.Payload)
		} else {
			values[types.TxValueKeyData] = ([]byte)(*args.Payload)
		}
	} else if args.TypeInt.IsEthereumTransaction() {
		// For Ethereum transactions, Input is an optional field.
		values[types.TxValueKeyData] = []byte{}
	}
	if args.CodeFormat != nil {
		values[types.TxValueKeyCodeFormat] = *args.CodeFormat
	}
	if args.HumanReadable != nil {
		values[types.TxValueKeyHumanReadable] = *args.HumanReadable
	}
	if args.Key != nil {
		serializer := accountkey.NewAccountKeySerializer()
		if err := rlp.DecodeBytes(*args.Key, &serializer); err == nil {
			values[types.TxValueKeyAccountKey] = serializer.GetKey()
		}
	}
	if args.ChainID != nil {
		values[types.TxValueKeyChainID] = (*big.Int)(args.ChainID)
	}
	if args.AccessList != nil {
		values[types.TxValueKeyAccessList] = *args.AccessList
	}
	if args.MaxPriorityFeePerGas != nil {
		values[types.TxValueKeyGasTipCap] = (*big.Int)(args.MaxPriorityFeePerGas)
	}
	if args.MaxFeePerGas != nil {
		values[types.TxValueKeyGasFeeCap] = (*big.Int)(args.MaxFeePerGas)
	}

	return values
}

// toTransaction returns an unsigned transaction filled with values in SendTxArgs.
func (args *SendTxArgs) toTransaction() (*types.Transaction, error) {
	var input []byte

	// provide detailed error messages to users (optional)
	if err := args.checkArgs(); err != nil {
		return nil, err
	}

	// for TxTypeLegacyTransaction
	if *args.TypeInt == types.TxTypeLegacyTransaction {
		if args.Data != nil && args.Payload != nil && !bytes.Equal(*args.Data, *args.Payload) {
			return nil, errTxArgInvalidInputData
		}

		if args.Data != nil {
			input = *args.Data
		} else if args.Payload != nil {
			input = *args.Payload
		}

		if args.Recipient == nil {
			if len(input) == 0 {
				return nil, errTxArgNilContractData
			}
			return types.NewContractCreation(uint64(*args.AccountNonce), (*big.Int)(args.Amount), uint64(*args.GasLimit), (*big.Int)(args.Price), input), nil
		}
		return types.NewTransaction(uint64(*args.AccountNonce), *args.Recipient, (*big.Int)(args.Amount), uint64(*args.GasLimit), (*big.Int)(args.Price), input), nil
	}

	// for other tx types except TxTypeLegacyTransaction
	values := args.genTxValuesMap()
	return types.NewTransactionWithMap(*args.TypeInt, values)
}

type ValueTransferTxArgs struct {
	TransactionArgs
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *ValueTransferTxArgs) setDefaults(ctx context.Context, b Backend) error {
	// Check if invalid arguments exist
	if args.InputData() != nil || args.MaxPriorityFeePerGas != nil || args.MaxFeePerGas != nil || args.AccessList != nil || args.ChainID != nil {
		return errors.New(fmt.Sprintln(
			"deny next arguments when it's valueTransferTx", "args.Data", args.Data, "args.Payload", args.Payload,
			"args.MaxPriorityFeePerGas", args.MaxPriorityFeePerGas, "args.MaxFeePerGas", args.MaxFeePerGas,
			"args.AccessList", args.AccessList, "args.ChainID", args.ChainID))
	}
	if args.GasLimit == nil {
		args.GasLimit = new(hexutil.Uint64)
		*(*uint64)(args.GasLimit) = 90000
	}
	if args.Price == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.Price = (*hexutil.Big)(price)
	}
	if args.AccountNonce == nil {
		nonce := b.GetPoolNonce(ctx, args.from())
		args.AccountNonce = (*hexutil.Uint64)(&nonce)
	}
	return nil
}

func (args *ValueTransferTxArgs) toTransaction() (*types.Transaction, error) {
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    (uint64)(*args.AccountNonce),
		types.TxValueKeyGasLimit: (uint64)(*args.GasLimit),
		types.TxValueKeyGasPrice: (*big.Int)(args.Price),
		types.TxValueKeyFrom:     args.From,
		types.TxValueKeyTo:       args.Recipient,
		types.TxValueKeyAmount:   (*big.Int)(args.Amount),
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}

type AccountUpdateTxArgs struct {
	TransactionArgs
	Key *hexutil.Bytes `json:"key"`
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *AccountUpdateTxArgs) setDefaults(ctx context.Context, b Backend) error {
	// Check if invalid arguments exist
	if args.Recipient != nil || args.Amount != nil || args.InputData() != nil || args.MaxPriorityFeePerGas != nil || args.MaxFeePerGas != nil || args.AccessList != nil || args.ChainID != nil {
		return errors.New(fmt.Sprintln(
			"deny next arguments when it's valueTransferTx", "args.Recipient", args.Recipient, "args.Amount", args.Amount,
			"args.Data", args.Data, "args.Payload", args.Payload,
			"args.MaxPriorityFeePerGas", args.MaxPriorityFeePerGas, "args.MaxFeePerGas", args.MaxFeePerGas,
			"args.AccessList", args.AccessList, "args.ChainID", args.ChainID))
	}

	if args.GasLimit == nil {
		args.GasLimit = new(hexutil.Uint64)
		*(*uint64)(args.GasLimit) = 90000
	}
	if args.Price == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.Price = (*hexutil.Big)(price)
	}
	if args.AccountNonce == nil {
		nonce := b.GetPoolNonce(ctx, args.from())
		args.AccountNonce = (*hexutil.Uint64)(&nonce)
	}
	return nil
}

func (args *AccountUpdateTxArgs) toTransaction() (*types.Transaction, error) {
	serializer := accountkey.NewAccountKeySerializer()

	if err := rlp.DecodeBytes(*args.Key, &serializer); err != nil {
		return nil, err
	}
	tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      (uint64)(*args.AccountNonce),
		types.TxValueKeyGasLimit:   (uint64)(*args.GasLimit),
		types.TxValueKeyGasPrice:   (*big.Int)(args.Price),
		types.TxValueKeyFrom:       args.From,
		types.TxValueKeyAccountKey: serializer.GetKey(),
	})
	if err != nil {
		return nil, err
	}

	return tx, nil
}
