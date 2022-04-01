// Copyright 2022 The klaytn Authors
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
	"math/big"
	"reflect"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/fork"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

//go:generate gencodec -type AccessTuple -out gen_access_tuple.go

// AccessList is an EIP-2930 access list
type AccessList []AccessTuple

// AccessTuple is the element type of the access list
type AccessTuple struct {
	Address     common.Address `json:"address"        gencodec:"required"`
	StorageKeys []common.Hash  `json:"storageKeys"    gencodec:"required"`
}

// StorageKeys returns the total number of storage keys in the access list.
func (al AccessList) StorageKeys() int {
	sum := 0
	for _, tuple := range al {
		sum += len(tuple.StorageKeys)
	}
	return sum
}

// TxInternalDataEthereumAccessList is the data of EIP-2930 access list transactions.
type TxInternalDataEthereumAccessList struct {
	ChainID      *big.Int
	AccountNonce uint64
	Price        *big.Int
	GasLimit     uint64
	Recipient    *common.Address `rlp:"nil"` // nil means contract creation.
	Amount       *big.Int
	Payload      []byte
	AccessList   AccessList

	// Signature values
	V *big.Int
	R *big.Int
	S *big.Int

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataEthereumAccessListJSON struct {
	Type         TxType           `json:"typeInt"`
	TypeStr      string           `json:"type"`
	ChainID      *hexutil.Big     `json:"chainId"`
	AccountNonce hexutil.Uint64   `json:"nonce"`
	Price        *hexutil.Big     `json:"gasPrice"`
	GasLimit     hexutil.Uint64   `json:"gas"`
	Recipient    *common.Address  `json:"to"`
	Amount       *hexutil.Big     `json:"value"`
	Payload      hexutil.Bytes    `json:"input"`
	AccessList   AccessList       `json:"accessList"`
	TxSignatures TxSignaturesJSON `json:"signatures"`
	Hash         *common.Hash     `json:"hash"`
}

func newEmptyTxInternalDataEthereumAccessList() *TxInternalDataEthereumAccessList {
	return &TxInternalDataEthereumAccessList{}
}

func newTxInternalDataEthereumAccessList() *TxInternalDataEthereumAccessList {
	return &TxInternalDataEthereumAccessList{
		ChainID:      new(big.Int),
		AccountNonce: 0,
		Price:        new(big.Int),
		GasLimit:     0,
		Recipient:    nil,
		Amount:       new(big.Int),
		Payload:      []byte{},
		AccessList:   AccessList{},
		V:            new(big.Int),
		R:            new(big.Int),
		S:            new(big.Int),
	}
}

func newTxInternalDataEthereumAccessListWithValues(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte, accessList AccessList, chainID *big.Int) *TxInternalDataEthereumAccessList {
	d := newTxInternalDataEthereumAccessList()

	d.AccountNonce = nonce
	d.Recipient = to
	d.GasLimit = gasLimit

	if chainID != nil {
		d.ChainID.Set(chainID)
	}

	if gasPrice != nil {
		d.Price.Set(gasPrice)
	}

	if amount != nil {
		d.Amount.Set(amount)
	}

	if len(data) > 0 {
		d.Payload = common.CopyBytes(data)
	}

	if accessList != nil {
		copy(d.AccessList, accessList)
	}

	return d
}

func newTxInternalDataEthereumAccessListWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataEthereumAccessList, error) {
	d := newTxInternalDataEthereumAccessList()

	if v, ok := values[TxValueKeyChainID].(*big.Int); ok {
		d.ChainID.Set(v)
		delete(values, TxValueKeyChainID)
	} else {
		return nil, errValueKeyChainIDInvalid
	}

	if v, ok := values[TxValueKeyNonce].(uint64); ok {
		d.AccountNonce = v
		delete(values, TxValueKeyNonce)
	} else {
		return nil, errValueKeyNonceMustUint64
	}

	if v, ok := values[TxValueKeyTo].(*common.Address); ok {
		d.Recipient = v
		delete(values, TxValueKeyTo)
	} else {
		return nil, errValueKeyToMustAddressPointer
	}

	if v, ok := values[TxValueKeyAmount].(*big.Int); ok {
		d.Amount.Set(v)
		delete(values, TxValueKeyAmount)
	} else {
		return nil, errValueKeyAmountMustBigInt
	}

	if v, ok := values[TxValueKeyData].([]byte); ok {
		d.Payload = common.CopyBytes(v)
		delete(values, TxValueKeyData)
	} else {
		return nil, errValueKeyDataMustByteSlice
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

	if v, ok := values[TxValueKeyAccessList].(AccessList); ok {
		d.AccessList = make(AccessList, len(v))
		copy(d.AccessList, v)
		delete(values, TxValueKeyAccessList)
	} else {
		return nil, errValueKeyAccessListInvalid
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return d, nil
}

func (t *TxInternalDataEthereumAccessList) Type() TxType {
	return TxTypeEthereumAccessList
}

func (t *TxInternalDataEthereumAccessList) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataEthereumAccessList) ChainId() *big.Int {
	return t.ChainID
}

func (t *TxInternalDataEthereumAccessList) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataEthereumAccessList) GetPrice() *big.Int {
	return new(big.Int).Set(t.Price)
}

func (t *TxInternalDataEthereumAccessList) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataEthereumAccessList) GetRecipient() *common.Address {
	return t.Recipient
}

func (t *TxInternalDataEthereumAccessList) GetAmount() *big.Int {
	return new(big.Int).Set(t.Amount)
}

func (t *TxInternalDataEthereumAccessList) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataEthereumAccessList) GetPayload() []byte {
	return t.Payload
}

func (t *TxInternalDataEthereumAccessList) GetAccessList() AccessList {
	return t.AccessList
}

func (t *TxInternalDataEthereumAccessList) SetHash(hash *common.Hash) {
	t.Hash = hash
}

func (t *TxInternalDataEthereumAccessList) SetSignature(signatures TxSignatures) {
	if len(signatures) != 1 {
		logger.Crit("TxTypeEthereumAccessList can receive only single signature!")
	}

	t.V = signatures[0].V
	t.R = signatures[0].R
	t.S = signatures[0].S
}

func (t *TxInternalDataEthereumAccessList) RawSignatureValues() TxSignatures {
	return TxSignatures{&TxSignature{t.V, t.R, t.S}}
}

func (t *TxInternalDataEthereumAccessList) ValidateSignature() bool {
	v := byte(t.V.Uint64())
	return crypto.ValidateSignatureValues(v, t.R, t.S, false)
}

func (t *TxInternalDataEthereumAccessList) RecoverAddress(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) (common.Address, error) {
	V := vfunc(t.V)
	return recoverPlain(txhash, t.R, t.S, V, homestead)
}

func (t *TxInternalDataEthereumAccessList) RecoverPubkey(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) ([]*ecdsa.PublicKey, error) {
	V := vfunc(t.V)

	pk, err := recoverPlainPubkey(txhash, t.R, t.S, V, homestead)
	if err != nil {
		return nil, err
	}

	return []*ecdsa.PublicKey{pk}, nil
}

func (t *TxInternalDataEthereumAccessList) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	return IntrinsicGas(t.Payload, t.AccessList, t.Recipient == nil, *fork.Rules(big.NewInt(int64(currentBlockNumber))))
}

func (t *TxInternalDataEthereumAccessList) setSignatureValues(chainID, v, r, s *big.Int) {
	t.ChainID, t.V, t.R, t.S = chainID, v, r, s
}

func (t *TxInternalDataEthereumAccessList) SerializeForSign() []interface{} {
	// If the chainId has nil or empty value, It will be set signer's chainId.
	return []interface{}{
		t.ChainID,
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
		t.AccessList,
	}
}

func (t *TxInternalDataEthereumAccessList) TxHash() common.Hash {
	return prefixedRlpHash(byte(t.Type()), []interface{}{
		t.ChainID,
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
		t.AccessList,
		t.V,
		t.R,
		t.S,
	})
}

func (t *TxInternalDataEthereumAccessList) SenderTxHash() common.Hash {
	return prefixedRlpHash(byte(t.Type()), []interface{}{
		t.ChainID,
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.Payload,
		t.AccessList,
		t.V,
		t.R,
		t.S,
	})
}

func (t *TxInternalDataEthereumAccessList) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataEthereumAccessList) Equal(a TxInternalData) bool {
	ta, ok := a.(*TxInternalDataEthereumAccessList)
	if !ok {
		return false
	}

	return t.ChainID.Cmp(ta.ChainID) == 0 &&
		t.AccountNonce == ta.AccountNonce &&
		t.Price.Cmp(ta.Price) == 0 &&
		t.GasLimit == ta.GasLimit &&
		equalRecipient(t.Recipient, ta.Recipient) &&
		t.Amount.Cmp(ta.Amount) == 0 &&
		reflect.DeepEqual(t.AccessList, ta.AccessList) &&
		t.V.Cmp(ta.V) == 0 &&
		t.R.Cmp(ta.R) == 0 &&
		t.S.Cmp(ta.S) == 0
}

func (t *TxInternalDataEthereumAccessList) String() string {
	var from, to string
	tx := &Transaction{data: t}

	v, r, s := t.V, t.R, t.S

	signer := LatestSignerForChainID(t.ChainId())
	if f, err := Sender(signer, tx); err != nil { // derive but don't cache
		from = "[invalid sender: invalid sig]"
	} else {
		from = fmt.Sprintf("%x", f[:])
	}

	if t.GetRecipient() == nil {
		to = "[contract creation]"
	} else {
		to = fmt.Sprintf("%x", t.GetRecipient().Bytes())
	}
	enc, _ := rlp.EncodeToBytes(tx)
	return fmt.Sprintf(`
		TX(%x)
		Contract: %v
		Chaind:   %#x
		From:     %s
		To:       %s
		Nonce:    %v
		GasPrice: %#x
		GasLimit  %#x
		Value:    %#x
		Data:     0x%x
		AccessList: %x
		V:        %#x
		R:        %#x
		S:        %#x
		Hex:      %x
	`,
		tx.Hash(),
		t.GetRecipient() == nil,
		t.ChainId(),
		from,
		to,
		t.GetAccountNonce(),
		t.GetPrice(),
		t.GetGasLimit(),
		t.GetAmount(),
		t.GetPayload(),
		t.AccessList,
		v,
		r,
		s,
		enc,
	)
}

func (t *TxInternalDataEthereumAccessList) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	if t.Recipient != nil {
		if common.IsPrecompiledContractAddress(*t.Recipient) {
			return kerrors.ErrPrecompiledContractAddress
		}
	}
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataEthereumAccessList) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	return nil
}

func (t *TxInternalDataEthereumAccessList) FillContractAddress(from common.Address, r *Receipt) {
	if t.Recipient == nil {
		r.ContractAddress = crypto.CreateAddress(from, t.AccountNonce)
	}
}

func (t *TxInternalDataEthereumAccessList) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	///////////////////////////////////////////////////////
	// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
	//start := time.Now()
	//defer func() {
	//	elapsed := time.Since(start)
	//	logger.Debug("[TxInternalDataLegacy] EVM execution done", "elapsed", elapsed)
	//}()
	///////////////////////////////////////////////////////
	if t.Recipient == nil {
		// Sender's nonce will be increased in '`vm.Create()`
		ret, _, usedGas, err = vm.Create(sender, t.Payload, gas, value, params.CodeFormatEVM)
	} else {
		stateDB.IncNonce(sender.Address())
		ret, usedGas, err = vm.Call(sender, *t.Recipient, t.Payload, gas, value)
	}
	return ret, usedGas, err
}

func (t *TxInternalDataEthereumAccessList) MakeRPCOutput() map[string]interface{} {
	return map[string]interface{}{
		"typeInt":    t.Type(),
		"type":       t.Type().String(),
		"chainId":    (*hexutil.Big)(t.ChainId()),
		"nonce":      hexutil.Uint64(t.AccountNonce),
		"gasPrice":   (*hexutil.Big)(t.Price),
		"gas":        hexutil.Uint64(t.GasLimit),
		"to":         t.Recipient,
		"input":      hexutil.Bytes(t.Payload),
		"value":      (*hexutil.Big)(t.Amount),
		"accessList": t.AccessList,
		"signatures": TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
	}
}

func (t *TxInternalDataEthereumAccessList) MarshalJSON() ([]byte, error) {
	return json.Marshal(TxInternalDataEthereumAccessListJSON{
		t.Type(),
		t.Type().String(),
		(*hexutil.Big)(t.ChainID),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.Recipient,
		(*hexutil.Big)(t.Amount),
		t.Payload,
		t.AccessList,
		TxSignaturesJSON{&TxSignatureJSON{(*hexutil.Big)(t.V), (*hexutil.Big)(t.R), (*hexutil.Big)(t.S)}},
		t.Hash,
	})
}

func (t *TxInternalDataEthereumAccessList) UnmarshalJSON(bytes []byte) error {
	js := &TxInternalDataEthereumAccessListJSON{}
	if err := json.Unmarshal(bytes, js); err != nil {
		return err
	}

	t.ChainID = (*big.Int)(js.ChainID)
	t.AccountNonce = uint64(js.AccountNonce)
	t.Price = (*big.Int)(js.Price)
	t.GasLimit = uint64(js.GasLimit)
	t.Recipient = js.Recipient
	t.Amount = (*big.Int)(js.Amount)
	t.Payload = js.Payload
	t.AccessList = js.AccessList
	t.V = (*big.Int)(js.TxSignatures[0].V)
	t.R = (*big.Int)(js.TxSignatures[0].R)
	t.S = (*big.Int)(js.TxSignatures[0].S)
	t.Hash = js.Hash

	return nil
}
