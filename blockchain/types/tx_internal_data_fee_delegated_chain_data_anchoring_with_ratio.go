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
	"crypto/ecdsa"
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

// TxInternalDataFeeDelegatedChainDataAnchoringWithRatio represents the fee-delegated transaction
// anchoring child chain data with a specified fee ratio between the sender and the fee payer.
// The ratio is a fee payer's ratio in percentage.
// For example, if it is 20, 20% of tx fee will be paid by the fee payer.
// 80% of tx fee will be paid by the sender.
type TxInternalDataFeeDelegatedChainDataAnchoringWithRatio struct {
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	From         common.Address
	Payload      []byte
	FeeRatio     FeeRatio

	TxSignatures

	FeePayer           common.Address
	FeePayerSignatures TxSignatures

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataFeeDelegatedChainDataAnchoringWithRatioJSON struct {
	Type               TxType           `json:"typeInt"`
	TypeStr            string           `json:"type"`
	AccountNonce       hexutil.Uint64   `json:"nonce"`
	Price              *hexutil.Big     `json:"gasPrice"`
	GasLimit           hexutil.Uint64   `json:"gas"`
	From               common.Address   `json:"from"`
	Payload            hexutil.Bytes    `json:"input"`
	InputJson          interface{}      `json:"inputJSON"`
	FeeRatio           hexutil.Uint     `json:"feeRatio"`
	TxSignatures       TxSignaturesJSON `json:"signatures"`
	FeePayer           common.Address   `json:"feePayer"`
	FeePayerSignatures TxSignaturesJSON `json:"feePayerSignatures"`
	Hash               *common.Hash     `json:"hash"`
}

func newTxInternalDataFeeDelegatedChainDataAnchoringWithRatio() *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio {
	h := common.Hash{}

	return &TxInternalDataFeeDelegatedChainDataAnchoringWithRatio{
		Price: new(big.Int),
		Hash:  &h,
	}
}

func newTxInternalDataFeeDelegatedChainDataAnchoringWithRatioWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataFeeDelegatedChainDataAnchoringWithRatio, error) {
	d := newTxInternalDataFeeDelegatedChainDataAnchoringWithRatio()

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

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) Type() TxType {
	return TxTypeFeeDelegatedChainDataAnchoringWithRatio
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) Equal(b TxInternalData) bool {
	tb, ok := b.(*TxInternalDataFeeDelegatedChainDataAnchoringWithRatio)
	if !ok {
		return false
	}

	return t.AccountNonce == tb.AccountNonce &&
		t.Price.Cmp(tb.Price) == 0 &&
		t.GasLimit == tb.GasLimit &&
		t.From == tb.From &&
		t.FeeRatio == tb.FeeRatio &&
		bytes.Equal(t.Payload, tb.Payload) &&
		t.TxSignatures.equal(tb.TxSignatures) &&
		t.FeePayer == tb.FeePayer &&
		t.FeePayerSignatures.equal(tb.FeePayerSignatures)
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) String() string {
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
	AnchoredData:  %s
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
		common.Bytes2Hex(t.Payload),
		t.TxSignatures.string(),
		t.FeePayer.String(),
		t.FeeRatio,
		t.FeePayerSignatures.string(),
		enc)
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) SerializeForSignToBytes() []byte {
	b, _ := rlp.EncodeToBytes(struct {
		Txtype       TxType
		AccountNonce uint64
		Price        *big.Int
		GasLimit     uint64
		From         common.Address
		AnchoredData []byte
		FeeRatio     FeeRatio
	}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.Payload,
		t.FeeRatio,
	})

	return b
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) SerializeForSign() []interface{} {
	return []interface{}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.Payload,
		t.FeeRatio,
	}
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) SenderTxHash() common.Hash {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, t.Type())
	rlp.Encode(hw, []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.From,
		t.Payload,
		t.FeeRatio,
		t.TxSignatures,
	})

	h := common.Hash{}

	hw.Sum(h[:0])

	return h
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetPrice() *big.Int {
	return new(big.Int).Set(t.Price)
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetRecipient() *common.Address {
	return nil
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetAmount() *big.Int {
	return common.Big0
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetFrom() common.Address {
	return t.From
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetPayload() []byte {
	return t.Payload
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetFeePayer() common.Address {
	return t.FeePayer
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetFeePayerRawSignatureValues() TxSignatures {
	return t.FeePayerSignatures.RawSignatureValues()
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) GetFeeRatio() FeeRatio {
	return t.FeeRatio
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) SetHash(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) SetSignature(s TxSignatures) {
	t.TxSignatures = s
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) SetFeePayerSignatures(s TxSignatures) {
	t.FeePayerSignatures = s
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) RecoverFeePayerPubkey(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) ([]*ecdsa.PublicKey, error) {
	return t.FeePayerSignatures.RecoverPubkey(txhash, homestead, vfunc)
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	gas := params.TxChainDataAnchoringGas + params.TxGasFeeDelegatedWithRatio

	gasPayloadWithGas, err := IntrinsicGasPayload(gas, t.Payload)
	if err != nil {
		return 0, err
	}

	return gasPayloadWithGas, nil
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	return nil
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	stateDB.IncNonce(sender.Address())
	return nil, gas, nil
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) MakeRPCOutput() map[string]interface{} {
	decoded, err := DecodeAnchoringDataToJSON(t.Payload)
	if err != nil {
		logger.Trace("decode anchor payload", "err", err)
	}

	return map[string]interface{}{
		"typeInt":            t.Type(),
		"type":               t.Type().String(),
		"gas":                hexutil.Uint64(t.GasLimit),
		"gasPrice":           (*hexutil.Big)(t.Price),
		"nonce":              hexutil.Uint64(t.AccountNonce),
		"input":              hexutil.Bytes(t.Payload),
		"inputJSON":          decoded,
		"feeRatio":           hexutil.Uint(t.FeeRatio),
		"signatures":         t.TxSignatures.ToJSON(),
		"feePayer":           t.FeePayer,
		"feePayerSignatures": t.FeePayerSignatures.ToJSON(),
	}
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) MarshalJSON() ([]byte, error) {
	decoded, err := DecodeAnchoringDataToJSON(t.Payload)
	if err != nil {
		logger.Trace("decode anchor payload", "err", err)
	}
	return json.Marshal(TxInternalDataFeeDelegatedChainDataAnchoringWithRatioJSON{
		t.Type(),
		t.Type().String(),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.From,
		t.Payload,
		decoded,
		(hexutil.Uint)(t.FeeRatio),
		t.TxSignatures.ToJSON(),
		t.FeePayer,
		t.FeePayerSignatures.ToJSON(),
		t.Hash,
	})
}

func (t *TxInternalDataFeeDelegatedChainDataAnchoringWithRatio) UnmarshalJSON(b []byte) error {
	js := &TxInternalDataFeeDelegatedChainDataAnchoringWithRatioJSON{}
	if err := json.Unmarshal(b, js); err != nil {
		return err
	}

	t.AccountNonce = uint64(js.AccountNonce)
	t.Price = (*big.Int)(js.Price)
	t.GasLimit = uint64(js.GasLimit)
	t.From = js.From
	t.Payload = js.Payload
	t.FeeRatio = FeeRatio(js.FeeRatio)
	t.TxSignatures = js.TxSignatures.ToTxSignatures()
	t.FeePayer = js.FeePayer
	t.FeePayerSignatures = js.FeePayerSignatures.ToTxSignatures()
	t.Hash = js.Hash

	return nil
}
