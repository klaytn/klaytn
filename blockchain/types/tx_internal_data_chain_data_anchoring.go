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
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

// TxInternalDataChainDataAnchoring represents the transaction anchoring child chain data.
type TxInternalDataChainDataAnchoring struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	From         common.Address
	Payload      []byte

	TxSignatures

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataChainDataAnchoringJSON struct {
	Type         TxType           `json:"typeInt"`
	TypeStr      string           `json:"type"`
	AccountNonce hexutil.Uint64   `json:"nonce"`
	Price        *hexutil.Big     `json:"gasPrice"`
	GasLimit     hexutil.Uint64   `json:"gas"`
	From         common.Address   `json:"from"`
	Payload      hexutil.Bytes    `json:"input"`
	InputJson    interface{}      `json:"inputJSON"`
	TxSignatures TxSignaturesJSON `json:"signatures"`
	Hash         *common.Hash     `json:"hash"`
}

func newTxInternalDataChainDataAnchoring() *TxInternalDataChainDataAnchoring {
	h := common.Hash{}

	return &TxInternalDataChainDataAnchoring{
		Price: new(big.Int),
		Hash:  &h,
	}
}

func newTxInternalDataChainDataAnchoringWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataChainDataAnchoring, error) {
	d := newTxInternalDataChainDataAnchoring()

	if v, ok := values[TxValueKeyNonce].(uint64); ok {
		d.AccountNonce = v
		delete(values, TxValueKeyNonce)
	} else {
		return nil, errValueKeyNonceMustUint64
	}

	if v, ok := values[TxValueKeyGasPrice].(*big.Int); ok {
		d.Price.Set(v)
		delete(values, TxValueKeyGasPrice)
	} else {
		return nil, errValueKeyGasPriceMustBigInt
	}

	if v, ok := values[TxValueKeyGasLimit].(uint64); ok {
		d.GasLimit = v
		delete(values, TxValueKeyGasLimit)
	} else {
		return nil, errValueKeyGasLimitMustUint64
	}

	if v, ok := values[TxValueKeyFrom].(common.Address); ok {
		d.From = v
		delete(values, TxValueKeyFrom)
	} else {
		return nil, errValueKeyFromMustAddress
	}

	if v, ok := values[TxValueKeyAnchoredData].([]byte); ok {
		d.Payload = v
		delete(values, TxValueKeyAnchoredData)
	} else {
		return nil, errValueKeyAnchoredDataMustByteSlice
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return d, nil
}

func (t *TxInternalDataChainDataAnchoring) Type() TxType {
	return TxTypeChainDataAnchoring
}

func (t *TxInternalDataChainDataAnchoring) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataChainDataAnchoring) Equal(b TxInternalData) bool {
	tb, ok := b.(*TxInternalDataChainDataAnchoring)
	if !ok {
		return false
	}

	return t.AccountNonce == tb.AccountNonce &&
		t.Price.Cmp(tb.Price) == 0 &&
		t.GasLimit == tb.GasLimit &&
		t.From == tb.From &&
		t.TxSignatures.equal(tb.TxSignatures) &&
		bytes.Equal(t.Payload, tb.Payload)
}

func (t *TxInternalDataChainDataAnchoring) String() string {
	ser := newTxInternalDataSerializerWithValues(t)
	enc, _ := rlp.EncodeToBytes(ser)
	tx := Transaction{data: t}

	return fmt.Sprintf(`
	TX(%x)
	Type:          %s
	From:          %s
	Nonce:         %v
	GasPrice:      %#x
	GasLimit:      %#x
	Signature:     %s
	Hex:           %x
	AnchoredData:  %s
`,
		tx.Hash(),
		t.Type().String(),
		t.From.String(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.TxSignatures.string(),
		enc,
		common.Bytes2Hex(t.Payload))
}

func (t *TxInternalDataChainDataAnchoring) SerializeForSignToBytes() []byte {
	b, _ := rlp.EncodeToBytes(struct {
		Txtype       TxType
		AccountNonce uint64
		Price        *big.Int
		GasLimit     uint64
		From         common.Address
		AnchoredData []byte
	}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.Payload,
	})

	return b
}

func (t *TxInternalDataChainDataAnchoring) SerializeForSign() []interface{} {
	return []interface{}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.Payload,
	}
}

func (t *TxInternalDataChainDataAnchoring) SenderTxHash() common.Hash {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, t.Type())
	rlp.Encode(hw, []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.Payload,
		t.TxSignatures,
	})

	h := common.Hash{}

	hw.Sum(h[:0])

	return h
}

func (t *TxInternalDataChainDataAnchoring) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataChainDataAnchoring) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataChainDataAnchoring) GetPrice() *big.Int {
	return new(big.Int).Set(t.Price)
}

func (t *TxInternalDataChainDataAnchoring) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataChainDataAnchoring) GetRecipient() *common.Address {
	return nil
}

func (t *TxInternalDataChainDataAnchoring) GetAmount() *big.Int {
	return common.Big0
}

func (t *TxInternalDataChainDataAnchoring) GetFrom() common.Address {
	return t.From
}

func (t *TxInternalDataChainDataAnchoring) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataChainDataAnchoring) GetPayload() []byte {
	return t.Payload
}

func (t *TxInternalDataChainDataAnchoring) SetHash(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataChainDataAnchoring) SetSignature(s TxSignatures) {
	t.TxSignatures = s
}

func (t *TxInternalDataChainDataAnchoring) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	gas := params.TxChainDataAnchoringGas

	gasPayloadWithGas, err := IntrinsicGasPayload(gas, t.Payload)
	if err != nil {
		return 0, err
	}

	return gasPayloadWithGas, nil
}

func (t *TxInternalDataChainDataAnchoring) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataChainDataAnchoring) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	return nil
}

func (t *TxInternalDataChainDataAnchoring) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	stateDB.IncNonce(sender.Address())
	return nil, gas, nil
}

func (t *TxInternalDataChainDataAnchoring) MakeRPCOutput() map[string]interface{} {
	decoded, err := DecodeAnchoringDataToJSON(t.Payload)
	if err != nil {
		logger.Trace("decode anchor payload", "err", err)
	}

	return map[string]interface{}{
		"typeInt":    t.Type(),
		"type":       t.Type().String(),
		"gas":        hexutil.Uint64(t.GasLimit),
		"gasPrice":   (*hexutil.Big)(t.Price),
		"nonce":      hexutil.Uint64(t.AccountNonce),
		"input":      hexutil.Bytes(t.Payload),
		"inputJSON":  decoded,
		"signatures": t.TxSignatures.ToJSON(),
	}
}

func (t *TxInternalDataChainDataAnchoring) MarshalJSON() ([]byte, error) {
	decoded, err := DecodeAnchoringDataToJSON(t.Payload)
	if err != nil {
		logger.Trace("decode anchor payload", "err", err)
	}
	return json.Marshal(TxInternalDataChainDataAnchoringJSON{
		t.Type(),
		t.Type().String(),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.From,
		t.Payload,
		decoded,
		t.TxSignatures.ToJSON(),
		t.Hash,
	})
}

func (t *TxInternalDataChainDataAnchoring) UnmarshalJSON(b []byte) error {
	js := &TxInternalDataChainDataAnchoringJSON{}
	if err := json.Unmarshal(b, js); err != nil {
		return err
	}

	t.AccountNonce = uint64(js.AccountNonce)
	t.Price = (*big.Int)(js.Price)
	t.GasLimit = uint64(js.GasLimit)
	t.From = js.From
	t.Payload = js.Payload
	t.TxSignatures = js.TxSignatures.ToTxSignatures()
	t.Hash = js.Hash

	return nil
}
