// Copyright 2019 The klaytn Authors
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
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

// TxInternalDataFeeDelegatedAccountUpdateWithRatio represents a fee-delegated transaction updating a key of an account
// with a specified fee ratio between the sender and the fee payer.
// The ratio is a fee payer's ratio in percentage.
// For example, if it is 20, 20% of tx fee will be paid by the fee payer.
// 80% of tx fee will be paid by the sender.
type TxInternalDataFeeDelegatedAccountUpdateWithRatio struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	From         common.Address
	Key          accountkey.AccountKey
	FeeRatio     FeeRatio

	TxSignatures

	FeePayer           common.Address
	FeePayerSignatures TxSignatures

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type txInternalDataFeeDelegatedAccountUpdateWithRatioSerializable struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	From         common.Address
	Key          []byte
	FeeRatio     FeeRatio

	TxSignatures

	FeePayer           common.Address
	FeePayerSignatures TxSignatures
}

type TxInternalDataFeeDelegatedAccountUpdateWithRatioJSON struct {
	Type               TxType           `json:"typeInt"`
	TypeStr            string           `json:"type"`
	AccountNonce       hexutil.Uint64   `json:"nonce"`
	Price              *hexutil.Big     `json:"gasPrice"`
	GasLimit           hexutil.Uint64   `json:"gas"`
	From               common.Address   `json:"from"`
	Key                hexutil.Bytes    `json:"key"`
	FeeRatio           hexutil.Uint     `json:"feeRatio"`
	TxSignatures       TxSignaturesJSON `json:"signatures"`
	FeePayer           common.Address   `json:"feePayer"`
	FeePayerSignatures TxSignaturesJSON `json:"feePayerSignatures"`
	Hash               *common.Hash     `json:"hash"`
}

func newTxInternalDataFeeDelegatedAccountUpdateWithRatio() *TxInternalDataFeeDelegatedAccountUpdateWithRatio {
	return &TxInternalDataFeeDelegatedAccountUpdateWithRatio{
		Price:        new(big.Int),
		Key:          accountkey.NewAccountKeyLegacy(),
		TxSignatures: NewTxSignatures(),
	}
}

func newTxInternalDataFeeDelegatedAccountUpdateWithRatioWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataFeeDelegatedAccountUpdateWithRatio, error) {
	d := newTxInternalDataFeeDelegatedAccountUpdateWithRatio()

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
		return nil, errValueKeyGasLimitMustUint64
	}

	if v, ok := values[TxValueKeyGasPrice].(*big.Int); ok {
		d.Price.Set(v)
		delete(values, TxValueKeyGasPrice)
	} else {
		return nil, errValueKeyGasPriceMustBigInt
	}

	if v, ok := values[TxValueKeyFrom].(common.Address); ok {
		d.From = v
		delete(values, TxValueKeyFrom)
	} else {
		return nil, errValueKeyFromMustAddress
	}

	if v, ok := values[TxValueKeyAccountKey].(accountkey.AccountKey); ok {
		d.Key = v
		delete(values, TxValueKeyAccountKey)
	} else {
		return nil, errValueKeyAccountKeyMustAccountKey
	}

	if v, ok := values[TxValueKeyFeePayer].(common.Address); ok {
		d.FeePayer = v
		delete(values, TxValueKeyFeePayer)
	} else {
		return nil, errValueKeyFeePayerMustAddress
	}

	if v, ok := values[TxValueKeyFeeRatioOfFeePayer].(FeeRatio); ok {
		d.FeeRatio = v
		delete(values, TxValueKeyFeeRatioOfFeePayer)
	} else {
		return nil, errValueKeyFeeRatioMustUint8
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return d, nil
}

func newTxInternalDataFeeDelegatedAccountUpdateWithRatioSerializable() *txInternalDataFeeDelegatedAccountUpdateWithRatioSerializable {
	return &txInternalDataFeeDelegatedAccountUpdateWithRatioSerializable{}
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) toSerializable() *txInternalDataFeeDelegatedAccountUpdateWithRatioSerializable {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return &txInternalDataFeeDelegatedAccountUpdateWithRatioSerializable{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		keyEnc,
		t.FeeRatio,
		t.TxSignatures,
		t.FeePayer,
		t.FeePayerSignatures,
	}
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) fromSerializable(serialized *txInternalDataFeeDelegatedAccountUpdateWithRatioSerializable) error {
	t.AccountNonce = serialized.AccountNonce
	t.Price = serialized.Price
	t.GasLimit = serialized.GasLimit
	t.From = serialized.From
	t.TxSignatures = serialized.TxSignatures
	t.FeePayer = serialized.FeePayer
	t.FeePayerSignatures = serialized.FeePayerSignatures
	t.FeeRatio = serialized.FeeRatio

	serializer := accountkey.NewAccountKeySerializer()
	if err := rlp.DecodeBytes(serialized.Key, serializer); err != nil {
		return err
	}
	t.Key = serializer.GetKey()

	return nil
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, t.toSerializable())
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) DecodeRLP(s *rlp.Stream) error {
	dec := newTxInternalDataFeeDelegatedAccountUpdateWithRatioSerializable()

	if err := s.Decode(dec); err != nil {
		return err
	}
	if err := t.fromSerializable(dec); err != nil {
		logger.Warn("DecodeRLP failed", "err", err)
		return kerrors.ErrUnserializableKey
	}

	return nil
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) MarshalJSON() ([]byte, error) {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return json.Marshal(TxInternalDataFeeDelegatedAccountUpdateWithRatioJSON{
		t.Type(),
		t.Type().String(),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.From,
		(hexutil.Bytes)(keyEnc),
		(hexutil.Uint)(t.FeeRatio),
		t.TxSignatures.ToJSON(),
		t.FeePayer,
		t.FeePayerSignatures.ToJSON(),
		t.Hash,
	})
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) UnmarshalJSON(b []byte) error {
	js := &TxInternalDataFeeDelegatedAccountUpdateWithRatioJSON{}
	if err := json.Unmarshal(b, js); err != nil {
		return err
	}

	ser := accountkey.NewAccountKeySerializer()
	if err := rlp.DecodeBytes(js.Key, ser); err != nil {
		logger.Warn("UnmarshalJSON failed", "err", err)
		return kerrors.ErrUnserializableKey
	}

	t.AccountNonce = uint64(js.AccountNonce)
	t.Price = (*big.Int)(js.Price)
	t.GasLimit = uint64(js.GasLimit)
	t.From = js.From
	t.Key = ser.GetKey()
	t.FeeRatio = FeeRatio(js.FeeRatio)
	t.TxSignatures = js.TxSignatures.ToTxSignatures()
	t.FeePayer = js.FeePayer
	t.FeePayerSignatures = js.FeePayerSignatures.ToTxSignatures()
	t.Hash = js.Hash

	return nil
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) Type() TxType {
	return TxTypeFeeDelegatedAccountUpdateWithRatio
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleAccountUpdate
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetPrice() *big.Int {
	return t.Price
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetRecipient() *common.Address {
	return nil
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetAmount() *big.Int {
	return common.Big0
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetFrom() common.Address {
	return t.From
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetFeePayer() common.Address {
	return t.FeePayer
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetFeePayerRawSignatureValues() TxSignatures {
	return t.FeePayerSignatures.RawSignatureValues()
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) GetFeeRatio() FeeRatio {
	return t.FeeRatio
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) SetHash(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) Equal(a TxInternalData) bool {
	ta, ok := a.(*TxInternalDataFeeDelegatedAccountUpdateWithRatio)
	if !ok {
		return false
	}

	return t.AccountNonce == ta.AccountNonce &&
		t.Price.Cmp(ta.Price) == 0 &&
		t.GasLimit == ta.GasLimit &&
		t.From == ta.From &&
		t.FeeRatio == ta.FeeRatio &&
		t.Key.Equal(ta.Key) &&
		t.TxSignatures.equal(ta.TxSignatures) &&
		t.FeePayer == ta.FeePayer &&
		t.FeePayerSignatures.equal(ta.FeePayerSignatures)
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) String() string {
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
	Key:           %s
	Signature:     %s
	FeePayer:      %s
	FeeRatio:      %d
	FeePayerSig:   %s
	Hex:           %x
`,
		tx.Hash(),
		t.Type().String(),
		t.From.String(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Key.String(),
		t.TxSignatures.string(),
		t.FeePayer.String(),
		t.FeeRatio,
		t.FeePayerSignatures.string(),
		enc)
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) SetSignature(s TxSignatures) {
	t.TxSignatures = s
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) SetFeePayerSignatures(s TxSignatures) {
	t.FeePayerSignatures = s
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) RecoverFeePayerPubkey(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) ([]*ecdsa.PublicKey, error) {
	return t.FeePayerSignatures.RecoverPubkey(txhash, homestead, vfunc)
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	gasKey, err := t.Key.AccountCreationGas(currentBlockNumber)
	if err != nil {
		return 0, err
	}

	return params.TxGasAccountUpdate + gasKey + params.TxGasFeeDelegatedWithRatio, nil
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) SerializeForSignToBytes() []byte {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	b, _ := rlp.EncodeToBytes(struct {
		Txtype       TxType
		AccountNonce uint64
		Price        *big.Int
		GasLimit     uint64
		From         common.Address
		Key          []byte
		FeeRatio     FeeRatio
	}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		keyEnc,
		t.FeeRatio,
	})

	return b
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) SerializeForSign() []interface{} {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return []interface{}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		keyEnc,
		t.FeeRatio,
	}
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) SenderTxHash() common.Hash {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	hw := sha3.NewKeccak256()
	rlp.Encode(hw, t.Type())
	rlp.Encode(hw, []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		keyEnc,
		t.FeeRatio,
		t.TxSignatures,
	})

	h := common.Hash{}

	hw.Sum(h[:0])

	return h
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	oldKey := stateDB.GetKey(t.From)
	if err := accountkey.CheckReplacable(oldKey, t.Key, currentBlockNumber); err != nil {
		return err
	}
	return nil
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	stateDB.IncNonce(sender.Address())
	err = stateDB.UpdateKey(sender.Address(), t.Key, currentBlockNumber)
	return nil, gas, err
}

func (t *TxInternalDataFeeDelegatedAccountUpdateWithRatio) MakeRPCOutput() map[string]interface{} {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return map[string]interface{}{
		"typeInt":            t.Type(),
		"type":               t.Type().String(),
		"gas":                hexutil.Uint64(t.GasLimit),
		"gasPrice":           (*hexutil.Big)(t.Price),
		"nonce":              hexutil.Uint64(t.AccountNonce),
		"key":                hexutil.Bytes(keyEnc),
		"feeRatio":           hexutil.Uint(t.FeeRatio),
		"signatures":         t.TxSignatures.ToJSON(),
		"feePayer":           t.FeePayer,
		"feePayerSignatures": t.FeePayerSignatures.ToJSON(),
	}
}
