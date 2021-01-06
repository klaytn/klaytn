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

// TxInternalDataAccountCreation represents a transaction creating an account.
type TxInternalDataAccountCreation struct {
	AccountNonce  uint64
	Price         *big.Int
	GasLimit      uint64
	Recipient     common.Address
	Amount        *big.Int
	From          common.Address
	HumanReadable bool
	Key           accountkey.AccountKey

	TxSignatures

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataAccountCreationJSON struct {
	Type          TxType           `json:"typeInt"`
	TypeStr       string           `json:"type"`
	AccountNonce  hexutil.Uint64   `json:"nonce"`
	Price         *hexutil.Big     `json:"gasPrice"`
	GasLimit      hexutil.Uint64   `json:"gas"`
	Recipient     common.Address   `json:"to"`
	Amount        *hexutil.Big     `json:"value"`
	From          common.Address   `json:"from"`
	HumanReadable bool             `json:"humanReadable"`
	Key           hexutil.Bytes    `json:"key"`
	TxSignatures  TxSignaturesJSON `json:"signatures"`
	Hash          *common.Hash     `json:"hash"`
}

// txInternalDataAccountCreationSerializable for RLP serialization.
type txInternalDataAccountCreationSerializable struct {
	AccountNonce  uint64
	Price         *big.Int
	GasLimit      uint64
	Recipient     common.Address
	Amount        *big.Int
	From          common.Address
	HumanReadable bool
	KeyData       []byte

	TxSignatures
}

func newTxInternalDataAccountCreation() *TxInternalDataAccountCreation {
	h := common.Hash{}
	return &TxInternalDataAccountCreation{
		Price:  new(big.Int),
		Amount: new(big.Int),
		Key:    accountkey.NewAccountKeyLegacy(),
		Hash:   &h,
	}
}

func newTxInternalDataAccountCreationWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataAccountCreation, error) {
	b := newTxInternalDataAccountCreation()

	if v, ok := values[TxValueKeyNonce].(uint64); ok {
		b.AccountNonce = v
		delete(values, TxValueKeyNonce)
	} else {
		return nil, errValueKeyNonceMustUint64
	}

	if v, ok := values[TxValueKeyGasPrice].(*big.Int); ok {
		b.Price.Set(v)
		delete(values, TxValueKeyGasPrice)
	} else {
		return nil, errValueKeyGasPriceMustBigInt
	}

	if v, ok := values[TxValueKeyGasLimit].(uint64); ok {
		b.GasLimit = v
		delete(values, TxValueKeyGasLimit)
	} else {
		return nil, errValueKeyGasLimitMustUint64
	}

	if v, ok := values[TxValueKeyTo].(common.Address); ok {
		b.Recipient = v
		delete(values, TxValueKeyTo)
	} else {
		return nil, errValueKeyToMustAddress
	}

	if v, ok := values[TxValueKeyAmount].(*big.Int); ok {
		b.Amount.Set(v)
		delete(values, TxValueKeyAmount)
	} else {
		return nil, errValueKeyAmountMustBigInt
	}

	if v, ok := values[TxValueKeyFrom].(common.Address); ok {
		b.From = v
		delete(values, TxValueKeyFrom)
	} else {
		return nil, errValueKeyFromMustAddress
	}

	if v, ok := values[TxValueKeyHumanReadable].(bool); ok {
		b.HumanReadable = v
		delete(values, TxValueKeyHumanReadable)
	} else {
		return nil, errValueKeyHumanReadableMustBool
	}

	if v, ok := values[TxValueKeyAccountKey].(accountkey.AccountKey); ok {
		b.Key = v
		delete(values, TxValueKeyAccountKey)
	} else {
		return nil, errValueKeyAccountKeyMustAccountKey
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return b, nil
}

func newTxInternalDataAccountCreationSerializable() *txInternalDataAccountCreationSerializable {
	return &txInternalDataAccountCreationSerializable{}
}

func (t *TxInternalDataAccountCreation) toSerializable() *txInternalDataAccountCreationSerializable {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return &txInternalDataAccountCreationSerializable{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.From,
		t.HumanReadable,
		keyEnc,
		t.TxSignatures,
	}
}

func (t *TxInternalDataAccountCreation) fromSerializable(serialized *txInternalDataAccountCreationSerializable) error {
	t.AccountNonce = serialized.AccountNonce
	t.Price = serialized.Price
	t.GasLimit = serialized.GasLimit
	t.Recipient = serialized.Recipient
	t.Amount = serialized.Amount
	t.From = serialized.From
	t.HumanReadable = serialized.HumanReadable
	t.TxSignatures = serialized.TxSignatures

	serializer := accountkey.NewAccountKeySerializer()
	if err := rlp.DecodeBytes(serialized.KeyData, serializer); err != nil {
		return err
	}
	t.Key = serializer.GetKey()

	return nil
}

func (t *TxInternalDataAccountCreation) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, t.toSerializable())
}

func (t *TxInternalDataAccountCreation) DecodeRLP(s *rlp.Stream) error {
	dec := newTxInternalDataAccountCreationSerializable()

	if err := s.Decode(dec); err != nil {
		return err
	}
	if err := t.fromSerializable(dec); err != nil {
		logger.Warn("DecodeRLP failed", "err", err)
		return kerrors.ErrUnserializableKey
	}

	return nil
}

func (t *TxInternalDataAccountCreation) MarshalJSON() ([]byte, error) {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return json.Marshal(TxInternalDataAccountCreationJSON{
		t.Type(),
		t.Type().String(),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.Recipient,
		(*hexutil.Big)(t.Amount),
		t.From,
		t.HumanReadable,
		(hexutil.Bytes)(keyEnc),
		t.TxSignatures.ToJSON(),
		t.Hash,
	})
}

func (t *TxInternalDataAccountCreation) UnmarshalJSON(b []byte) error {
	js := &TxInternalDataAccountCreationJSON{}
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
	t.Recipient = js.Recipient
	t.Amount = (*big.Int)(js.Amount)
	t.From = js.From
	t.HumanReadable = js.HumanReadable
	t.Key = ser.GetKey()
	t.TxSignatures = js.TxSignatures.ToTxSignatures()
	t.Hash = js.Hash

	return nil
}

func (t *TxInternalDataAccountCreation) Type() TxType {
	return TxTypeAccountCreation
}

func (t *TxInternalDataAccountCreation) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataAccountCreation) Equal(a TxInternalData) bool {
	ta, ok := a.(*TxInternalDataAccountCreation)
	if !ok {
		return false
	}

	return t.AccountNonce == ta.AccountNonce &&
		t.Price.Cmp(ta.Price) == 0 &&
		t.GasLimit == ta.GasLimit &&
		t.Recipient == ta.Recipient &&
		t.Amount.Cmp(ta.Amount) == 0 &&
		t.From == ta.From &&
		t.HumanReadable == ta.HumanReadable &&
		t.Key.Equal(ta.Key) &&
		t.TxSignatures.equal(ta.TxSignatures)
}

func (t *TxInternalDataAccountCreation) String() string {
	ser := newTxInternalDataSerializerWithValues(t)
	tx := Transaction{data: t}
	enc, _ := rlp.EncodeToBytes(ser)
	return fmt.Sprintf(`
	TX(%x)
	Type:          %s
	From:          %s
	To:            %s
	Nonce:         %v
	GasPrice:      %#x
	GasLimit:      %#x
	Value:         %#x
	HumanReadable: %t
	Key:           %s
	Signature:     %s
	Hex:           %x
`,
		tx.Hash(),
		t.Type().String(),
		t.From.String(),
		t.Recipient.String(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Amount,
		t.HumanReadable,
		t.Key.String(),
		t.TxSignatures.string(),
		enc)
}

func (t *TxInternalDataAccountCreation) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataAccountCreation) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataAccountCreation) GetPrice() *big.Int {
	return new(big.Int).Set(t.Price)
}

func (t *TxInternalDataAccountCreation) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataAccountCreation) GetRecipient() *common.Address {
	if t.Recipient == (common.Address{}) {
		return nil
	}

	to := common.Address(t.Recipient)
	return &to
}

func (t *TxInternalDataAccountCreation) GetAmount() *big.Int {
	return new(big.Int).Set(t.Amount)
}

func (t *TxInternalDataAccountCreation) GetFrom() common.Address {
	return t.From
}

func (t *TxInternalDataAccountCreation) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataAccountCreation) SetHash(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataAccountCreation) SetSignature(s TxSignatures) {
	t.TxSignatures = s
}

func (t *TxInternalDataAccountCreation) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	gasKey, err := t.Key.AccountCreationGas(currentBlockNumber)
	if err != nil {
		return 0, err
	}

	gas := params.TxGasAccountCreation + gasKey
	if t.HumanReadable {
		gas += params.TxGasHumanReadable
	}

	return gas, nil
}

func (t *TxInternalDataAccountCreation) SerializeForSignToBytes() []byte {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	b, _ := rlp.EncodeToBytes(struct {
		Txtype        TxType
		AccountNonce  uint64
		Price         *big.Int
		GasLimit      uint64
		Recipient     common.Address
		Amount        *big.Int
		From          common.Address
		HumanReadable bool
		Key           []byte
	}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.From,
		t.HumanReadable,
		keyEnc,
	})

	return b
}

func (t *TxInternalDataAccountCreation) SerializeForSign() []interface{} {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return []interface{}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.From,
		t.HumanReadable,
		keyEnc,
	}
}

func (t *TxInternalDataAccountCreation) SenderTxHash() common.Hash {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	hw := sha3.NewKeccak256()
	rlp.Encode(hw, t.Type())
	rlp.Encode(hw, []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.From,
		t.HumanReadable,
		keyEnc,
		t.TxSignatures,
	})

	h := common.Hash{}

	hw.Sum(h[:0])

	return h
}

func (t *TxInternalDataAccountCreation) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	return errUndefinedTxType
	//if t.HumanReadable {
	//	return kerrors.ErrHumanReadableNotSupported
	//}
	//if common.IsPrecompiledContractAddress(t.Recipient) {
	//	return kerrors.ErrPrecompiledContractAddress
	//}
	//if err := t.Key.CheckInstallable(currentBlockNumber); err != nil {
	//	return err
	//}
	//
	//return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataAccountCreation) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	return errUndefinedTxType
	//// Fail if the address is already created.
	//if stateDB.Exist(t.Recipient) {
	//	return kerrors.ErrAccountAlreadyExists
	//}
	//return nil
}

func (t *TxInternalDataAccountCreation) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	to := t.Recipient
	stateDB.IncNonce(sender.Address())
	stateDB.CreateEOA(to, t.HumanReadable, t.Key)
	return vm.Call(sender, to, []byte{}, gas, value)
}

func (t *TxInternalDataAccountCreation) MakeRPCOutput() map[string]interface{} {
	serializer := accountkey.NewAccountKeySerializerWithAccountKey(t.Key)
	keyEnc, _ := rlp.EncodeToBytes(serializer)

	return map[string]interface{}{
		"typeInt":       t.Type(),
		"type":          t.Type().String(),
		"gas":           hexutil.Uint64(t.GasLimit),
		"gasPrice":      (*hexutil.Big)(t.Price),
		"nonce":         hexutil.Uint64(t.AccountNonce),
		"to":            t.Recipient,
		"value":         (*hexutil.Big)(t.Amount),
		"humanReadable": t.HumanReadable,
		"key":           hexutil.Bytes(keyEnc),
		"signatures":    t.TxSignatures.ToJSON(),
	}
}
