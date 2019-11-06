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

type NewTxArgs interface {
	setDefaults(context.Context, Backend) error
	toTransaction() (*types.Transaction, error)
	from() common.Address
}

// SendTxArgs represents the arguments to submit a new transaction into the transaction pool.
type SendTxArgs struct {
	TypeInt  *types.TxType   `json:"typeInt"`
	From     common.Address  `json:"from"`
	To       *common.Address `json:"to"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Value    *hexutil.Big    `json:"value"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	// We accept "data" and "input" for backwards-compatibility reasons. "input" is the
	// newer name and should be preferred by clients.
	Data  *hexutil.Bytes `json:"data"`
	Input *hexutil.Bytes `json:"input"`

	CodeFormat    *params.CodeFormat `json:"codeFormat"`
	HumanReadable *bool              `json:"humanReadable"`

	AccountKey *hexutil.Bytes `json:"Key"`

	FeePayer *common.Address `json:"feePayer"`
	FeeRatio *types.FeeRatio `json:"feeRatio"`

	Signatures types.TxSignaturesJSON `json:"Signatures"`
}

// setDefaults is a helper function that fills in default values for unspecified common tx fields.
func (args *SendTxArgs) setDefaults(ctx context.Context, b Backend) error {
	if args.TypeInt == nil {
		args.TypeInt = new(types.TxType)
		*args.TypeInt = types.TxTypeLegacyTransaction
	}
	if args.Gas == nil {
		args.Gas = new(hexutil.Uint64)
		*args.Gas = hexutil.Uint64(90000)
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

// checkArgs checks the validity of SendTxArgs values.
// The each tx types has its own validation logic to give detailed errors to users.
func (args *SendTxArgs) checkArgs() error {
	if args.TypeInt == nil {
		return errTxArgNilTxType
	}

	// TODO-Klaytn-TxType Arguments validation will be implemented for each tx type
	//switch *args.TypeInt {
	//case types.TxTypeLegacyTransaction:
	//case types.TxTypeValueTransfer:
	//case types.TxTypeFeeDelegatedValueTransfer:
	//}

	return nil
}

// genTxValuesMap generates a value map used used in "NewTransactionWithMap" function.
// This function assigned all non-nil values regardless of the tx type.
// Invalid values in the map will be validated in "NewTransactionWithMap" function.
func (args *SendTxArgs) genTxValuesMap() map[types.TxValueKeyType]interface{} {
	values := make(map[types.TxValueKeyType]interface{})

	// common tx fields. They should have values after executing "setDefaults" function.
	if args.TypeInt == nil || args.Nonce == nil || args.Gas == nil || args.GasPrice == nil {
		return values
	}
	values[types.TxValueKeyFrom] = args.From
	values[types.TxValueKeyNonce] = uint64(*args.Nonce)
	values[types.TxValueKeyGasLimit] = uint64(*args.Gas)
	values[types.TxValueKeyGasPrice] = (*big.Int)(args.GasPrice)

	// optional tx fields
	if args.TypeInt.IsContractDeploy() {
		// contract deploy type allows nil as TxValueKeyTo value
		values[types.TxValueKeyTo] = (*common.Address)(args.To)
	} else if args.To != nil {
		values[types.TxValueKeyTo] = *args.To
	}
	if args.FeePayer != nil {
		values[types.TxValueKeyFeePayer] = *args.FeePayer
	}
	if args.FeeRatio != nil {
		values[types.TxValueKeyFeeRatioOfFeePayer] = *args.FeeRatio
	}
	if args.Value != nil {
		values[types.TxValueKeyAmount] = (*big.Int)(args.Value)
	}
	if args.Data != nil {
		values[types.TxValueKeyData] = ([]byte)(*args.Data)
	}
	if args.CodeFormat != nil {
		values[types.TxValueKeyCodeFormat] = *args.CodeFormat
	}
	if args.HumanReadable != nil {
		values[types.TxValueKeyHumanReadable] = *args.HumanReadable
	}
	if args.AccountKey != nil {
		serializer := accountkey.NewAccountKeySerializer()
		if err := rlp.DecodeBytes(*args.AccountKey, &serializer); err == nil {
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
		if args.Data != nil && args.Input != nil && !bytes.Equal(*args.Data, *args.Input) {
			return nil, errTxArgInvalidInputData
		}

		if args.Data != nil {
			input = *args.Data
		} else if args.Input != nil {
			input = *args.Input
		}

		if args.To == nil {
			if len(input) == 0 {
				return nil, errTxArgNilContractData
			}
			return types.NewContractCreation(uint64(*args.Nonce), (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input), nil
		}
		return types.NewTransaction(uint64(*args.Nonce), *args.To, (*big.Int)(args.Value), uint64(*args.Gas), (*big.Int)(args.GasPrice), input), nil
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
