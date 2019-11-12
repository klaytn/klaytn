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
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"math/big"
	"reflect"
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
	var mapOfFieldMap = map[types.TxType]map[string]bool{}
	var internalDataTypes = map[types.TxType]interface{}{
		// since legacy tx has optional fields, some fields can be omitted
		//types.TxTypeLegacyTransaction:                           types.TxInternalDataLegacy{},
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
}

// SendTxArgs represents the arguments to submit a new transaction into the transaction pool.
type SendTxArgs struct {
	TypeInt      *types.TxType   `json:"typeInt"`
	From         common.Address  `json:"from"`
	Recipient    *common.Address `json:"to"`
	GasLimit     *hexutil.Uint64 `json:"gas"`
	Price        *hexutil.Big    `json:"gasPrice"`
	Amount       *hexutil.Big    `json:"value"`
	AccountNonce *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data    *hexutil.Bytes `json:"data"`
	Payload *hexutil.Bytes `json:"input"`

	CodeFormat    *params.CodeFormat `json:"codeFormat"`
	HumanReadable *bool              `json:"humanReadable"`

	Key *hexutil.Bytes `json:"key"`

	FeePayer *common.Address `json:"feePayer"`
	FeeRatio *types.FeeRatio `json:"feeRatio"`

	TxSignatures types.TxSignaturesJSON `json:"signatures"`
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
	if args.Price == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.Price = (*hexutil.Big)(price)
	}
	if args.AccountNonce == nil {
		nonce := b.GetPoolNonce(ctx, args.From)
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
	// Skip legacy transaction type since it has optional fields
	if *args.TypeInt == types.TxTypeLegacyTransaction {
		return nil
	}

	argsType := reflect.TypeOf(*args)
	argsValue := reflect.ValueOf(*args)

	for i := 0; i < argsType.NumField(); i++ {
		// Skip From since it is an essential field and a non-pointer value
		// Skip TxSignatures since the value is not considered by all APIs
		if argsType.Field(i).Name == "From" || argsType.Field(i).Name == "TxSignatures" {
			continue
		}

		// An args field doesn't have a value but the field name exist on the tx type
		if argsValue.Field(i).IsNil() && isTxField[*args.TypeInt][argsType.Field(i).Name] {
			return errors.New((string)(argsType.Field(i).Tag) + " is required for " + (*args.TypeInt).String())
		}

		// An args field has a value but the field name doesn't exist on the tx type
		if !argsValue.Field(i).IsNil() && !isTxField[*args.TypeInt][argsType.Field(i).Name] {
			return errors.New((string)(argsType.Field(i).Tag) + " is not a field of " + (*args.TypeInt).String())
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
	if args.TypeInt == nil || args.AccountNonce == nil || args.GasLimit == nil || args.Price == nil {
		return values
	}
	values[types.TxValueKeyFrom] = args.From
	values[types.TxValueKeyNonce] = uint64(*args.AccountNonce)
	values[types.TxValueKeyGasLimit] = uint64(*args.GasLimit)
	values[types.TxValueKeyGasPrice] = (*big.Int)(args.Price)

	// optional tx fields
	if args.TypeInt.IsContractDeploy() {
		// contract deploy type allows nil as TxValueKeyTo value
		values[types.TxValueKeyTo] = (*common.Address)(args.Recipient)
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
	}
	if args.Payload != nil {
		values[types.TxValueKeyData] = ([]byte)(*args.Payload)
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
	From     common.Address  `json:"from"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	To       common.Address  `json:"to"`
	Value    *hexutil.Big    `json:"value"`
}

func (args *ValueTransferTxArgs) from() common.Address {
	return args.From
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *ValueTransferTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.Gas == nil {
		args.Gas = new(hexutil.Uint64)
		*(*uint64)(args.Gas) = 90000
	}
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Nonce == nil {
		nonce := b.GetPoolNonce(ctx, args.From)
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	return nil
}

func (args *ValueTransferTxArgs) toTransaction() (*types.Transaction, error) {
	tx, err := types.NewTransactionWithMap(types.TxTypeValueTransfer, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:    (uint64)(*args.Nonce),
		types.TxValueKeyGasLimit: (uint64)(*args.Gas),
		types.TxValueKeyGasPrice: (*big.Int)(args.GasPrice),
		types.TxValueKeyFrom:     args.From,
		types.TxValueKeyTo:       args.To,
		types.TxValueKeyAmount:   (*big.Int)(args.Value),
	})

	if err != nil {
		return nil, err
	}

	return tx, nil
}

type AccountUpdateTxArgs struct {
	From     common.Address  `json:"from"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	Key      *hexutil.Bytes  `json:"key"`
}

func (args *AccountUpdateTxArgs) from() common.Address {
	return args.From
}

// setDefaults is a helper function that fills in default values for unspecified tx fields.
func (args *AccountUpdateTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.Gas == nil {
		args.Gas = new(hexutil.Uint64)
		*(*uint64)(args.Gas) = 90000
	}
	if args.GasPrice == nil {
		price, err := b.SuggestPrice(ctx)
		if err != nil {
			return err
		}
		args.GasPrice = (*hexutil.Big)(price)
	}
	if args.Nonce == nil {
		nonce := b.GetPoolNonce(ctx, args.From)
		args.Nonce = (*hexutil.Uint64)(&nonce)
	}
	return nil
}

func (args *AccountUpdateTxArgs) toTransaction() (*types.Transaction, error) {
	serializer := accountkey.NewAccountKeySerializer()

	if err := rlp.DecodeBytes(*args.Key, &serializer); err != nil {
		return nil, err
	}
	tx, err := types.NewTransactionWithMap(types.TxTypeAccountUpdate, map[types.TxValueKeyType]interface{}{
		types.TxValueKeyNonce:      (uint64)(*args.Nonce),
		types.TxValueKeyGasLimit:   (uint64)(*args.Gas),
		types.TxValueKeyGasPrice:   (*big.Int)(args.GasPrice),
		types.TxValueKeyFrom:       args.From,
		types.TxValueKeyAccountKey: serializer.GetKey(),
	})

	if err != nil {
		return nil, err
	}

	return tx, nil
}
