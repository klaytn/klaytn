// Copyright 2021 The klaytn CBDC Authors
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

package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types/account"
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

// TxInternalDataBalanceLimitUpdate represents a transaction updating the balance limit of an account.
type TxInternalDataBalanceLimitUpdate struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	From         common.Address
	BalanceLimit *big.Int

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`

	TxSignatures
}

type TxInternalDataBalanceLimitUpdateJSON struct {
	Type         TxType           `json:"typeInt"`
	TypeStr      string           `json:"type"`
	AccountNonce hexutil.Uint64   `json:"nonce"`
	Price        *hexutil.Big     `json:"gasPrice"`
	GasLimit     hexutil.Uint64   `json:"gas"`
	From         common.Address   `json:"from"`
	BalanceLimit *hexutil.Big     `json:"balanceLimit"`
	TxSignatures TxSignaturesJSON `json:"signatures"`
	Hash         *common.Hash     `json:"hash"`
}

func newTxInternalDataBalanceLimitUpdate() *TxInternalDataBalanceLimitUpdate {
	return &TxInternalDataBalanceLimitUpdate{
		Price:        new(big.Int),
		BalanceLimit: new(big.Int),
		TxSignatures: NewTxSignatures(),
	}
}

func newTxInternalDataBalanceLimitUpdateWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataBalanceLimitUpdate, error) {
	d := newTxInternalDataBalanceLimitUpdate()

	if v, ok := values[TxValueKeyNonce].(uint64); ok {
		d.AccountNonce = v
		delete(values, TxValueKeyNonce)
	} else {
		return nil, errValueKeyNonceMustUint64
	}

	if v, ok := values[TxValueKeyGasLimit].(uint64); ok {
		d.GasLimit = v
		delete(values, TxValueKeyGasLimit)
	} else {
		d.GasLimit = 0
		delete(values, TxValueKeyGasLimit)
	}

	if v, ok := values[TxValueKeyGasPrice].(*big.Int); ok {
		d.Price.Set(v)
		delete(values, TxValueKeyGasPrice)
	} else {
		d.Price.Set(big.NewInt(0))
		delete(values, TxValueKeyGasPrice)
	}

	if v, ok := values[TxValueKeyFrom].(common.Address); ok {
		d.From = v
		delete(values, TxValueKeyFrom)
	} else {
		return nil, errValueKeyFromMustAddress
	}

	if v, ok := values[TxValueKeyBalanceLimit].(*big.Int); ok {
		d.BalanceLimit.Set(v)
		delete(values, TxValueKeyBalanceLimit)
	} else {
		return nil, errValueKeyBalanceLimitMustBigInt
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return d, nil
}

func (t *TxInternalDataBalanceLimitUpdate) MarshalJSON() ([]byte, error) {
	return json.Marshal(TxInternalDataBalanceLimitUpdateJSON{
		t.Type(),
		t.Type().String(),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.From,
		(*hexutil.Big)(t.BalanceLimit),
		t.TxSignatures.ToJSON(),
		t.Hash,
	})
}

func (t *TxInternalDataBalanceLimitUpdate) UnmarshalJSON(b []byte) error {
	js := &TxInternalDataBalanceLimitUpdateJSON{}
	if err := json.Unmarshal(b, js); err != nil {
		return err
	}

	t.AccountNonce = uint64(js.AccountNonce)
	t.Price = (*big.Int)(js.Price)
	t.GasLimit = uint64(js.GasLimit)
	t.From = js.From
	t.BalanceLimit = (*big.Int)(js.BalanceLimit)
	t.TxSignatures = js.TxSignatures.ToTxSignatures()
	t.Hash = js.Hash

	return nil
}

func (t *TxInternalDataBalanceLimitUpdate) Type() TxType {
	return TxTypeBalanceLimitUpdate
}

func (t *TxInternalDataBalanceLimitUpdate) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleAccountUpdate
}

func (t *TxInternalDataBalanceLimitUpdate) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataBalanceLimitUpdate) GetPrice() *big.Int {
	return t.Price
}

func (t *TxInternalDataBalanceLimitUpdate) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataBalanceLimitUpdate) GetRecipient() *common.Address {
	return nil
}

func (t *TxInternalDataBalanceLimitUpdate) GetAmount() *big.Int {
	return common.Big0
}

func (t *TxInternalDataBalanceLimitUpdate) GetFrom() common.Address {
	return t.From
}

func (t *TxInternalDataBalanceLimitUpdate) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataBalanceLimitUpdate) SetHash(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataBalanceLimitUpdate) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataBalanceLimitUpdate) Equal(a TxInternalData) bool {
	ta, ok := a.(*TxInternalDataBalanceLimitUpdate)
	if !ok {
		return false
	}

	return t.AccountNonce == ta.AccountNonce &&
		t.Price.Cmp(ta.Price) == 0 &&
		t.GasLimit == ta.GasLimit &&
		t.From == ta.From &&
		t.BalanceLimit.Cmp(ta.BalanceLimit) == 0 &&
		t.TxSignatures.equal(ta.TxSignatures)
}

func (t *TxInternalDataBalanceLimitUpdate) String() string {
	ser := newTxInternalDataSerializerWithValues(t)
	tx := Transaction{data: t}
	enc, _ := rlp.EncodeToBytes(ser)
	return fmt.Sprintf(`
	TX(%x)
	Type:          %s
	From:          %s
	Nonce:         %v
	GasPrice:      %#x
	GasLimit:      %#x
	BalanceLimit:  %#x
	Signature:     %s
	Hex:           %x
`,
		tx.Hash(),
		t.Type().String(),
		t.From.String(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.BalanceLimit,
		t.TxSignatures.string(),
		enc)
}

func (t *TxInternalDataBalanceLimitUpdate) SetSignature(s TxSignatures) {
	t.TxSignatures = s
}

func (t *TxInternalDataBalanceLimitUpdate) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	return params.TxGasBalanceLimitUpdate, nil
}

func (t *TxInternalDataBalanceLimitUpdate) SerializeForSignToBytes() []byte {
	b, _ := rlp.EncodeToBytes(struct {
		Txtype       TxType
		AccountNonce uint64
		Price        *big.Int
		GasLimit     uint64
		From         common.Address
		BalanceLimit *big.Int
	}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.BalanceLimit,
	})

	return b
}

func (t *TxInternalDataBalanceLimitUpdate) SerializeForSign() []interface{} {
	b, _ := rlp.EncodeToBytes(struct {
		Txtype       TxType
		AccountNonce uint64
		Price        *big.Int
		GasLimit     uint64
		From         common.Address
		BalanceLimit *big.Int
	}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.BalanceLimit,
	})

	return []interface{}{b}
}

func (t *TxInternalDataBalanceLimitUpdate) SenderTxHash() common.Hash {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, t.Type())
	rlp.Encode(hw, []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.BalanceLimit,
		t.TxSignatures,
	})

	h := common.Hash{}

	hw.Sum(h[:0])

	return h
}

func (t *TxInternalDataBalanceLimitUpdate) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataBalanceLimitUpdate) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	if stateDB.IsProgramAccount(t.From) || common.IsPrecompiledContractAddress(t.From) {
		return account.ErrNotEOA
	}
	return nil
}

func (t *TxInternalDataBalanceLimitUpdate) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	stateDB.IncNonce(sender.Address())
	stateDB.SetBalanceLimit(sender.Address(), t.BalanceLimit)
	return nil, gas, nil
}

func (t *TxInternalDataBalanceLimitUpdate) MakeRPCOutput() map[string]interface{} {
	return map[string]interface{}{
		"type":         t.Type().String(),
		"typeInt":      t.Type(),
		"gas":          hexutil.Uint64(t.GasLimit),
		"gasPrice":     (*hexutil.Big)(t.Price),
		"nonce":        hexutil.Uint64(t.AccountNonce),
		"balanceLimit": (*hexutil.Big)(t.BalanceLimit),
		"signatures":   t.TxSignatures.ToJSON(),
	}
}
