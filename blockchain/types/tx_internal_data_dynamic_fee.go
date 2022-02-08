// Copyright 2021 The klaytn Authors
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
	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"math/big"
)

type TxInternalDataDynamicFee struct {
	ChainID           *big.Int
	AccountNonce      uint64
	GasTipCap         *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap         *big.Int // a.k.a. maxFeePerGas
	GasLimit          uint64
	Recipient         *common.Address `rlp:"nil"` // nil means contract creation
	Amount            *big.Int
	Payload           []byte
	AccessList        AccessList

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataDynamicFeeJSON struct {
	Type                 TxType           `json:"typeInt"`
	TypeStr              string           `json:"type"`
	AccountNonce         hexutil.Uint64   `json:"nonce"`
	MaxPriorityFeePerGas *hexutil.Big     `json:"maxPriorityFeePerGas"`
	MaxFeePerGas         *hexutil.Big     `json:"maxFeePerGas"`
	GasLimit             hexutil.Uint64   `json:"gas"`
	Recipient            *common.Address  `json:"to"`
	Amount               *hexutil.Big     `json:"value"`
	Payload              hexutil.Bytes    `json:"input"`
	TxSignatures         TxSignaturesJSON `json:"signatures"`

	AccessList           AccessList       `json:"accessList"`
	ChainID              *hexutil.Big     `json:"chainId"`

	Hash                 *common.Hash     `json:"hash"`
}


func newEmptyTxInternalDataDynamicFee() *TxInternalDataDynamicFee {
	return &TxInternalDataDynamicFee{}
}

func newTxInternalDataDynamicFee() *TxInternalDataDynamicFee {
	return &TxInternalDataDynamicFee{
		ChainID: new(big.Int),
		AccountNonce: 0,
		GasTipCap:         new(big.Int),
		GasFeeCap:         new(big.Int),
		GasLimit:          0,
		Recipient:         nil,
		Amount:            new(big.Int),
		Payload:           []byte{},
		AccessList:        AccessList{},

		// Signature values
		V: new(big.Int),
		R: new(big.Int),
		S: new(big.Int),
	}
}

func newTxInternalDataDynamicFeeWithValues(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasTipCap *big.Int, gasFeeCap *big.Int, data []byte, accessList AccessList, chainID *big.Int) *TxInternalDataDynamicFee {
	d := newTxInternalDataDynamicFee()

	d.AccountNonce = nonce
	d.Recipient = to
	d.GasLimit = gasLimit

	if len(data) > 0 {
		d.Payload = common.CopyBytes(data)
	}

	if amount != nil {
		d.Amount.Set(amount)
	}

	if gasTipCap != nil {
		d.GasTipCap.Set(gasTipCap)
	}

	if gasFeeCap != nil {
		d.GasFeeCap.Set(gasFeeCap)
	}

	if accessList != nil {
		copy(d.AccessList, accessList)
	}

	if chainID != nil {
		d.ChainID.Set(chainID)
	}

	return d
}

func newTxInternalDataDynamicFeeWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataDynamicFee, error) {
	return nil, nil
}

func (t *TxInternalDataDynamicFee) Type() TxType {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetAccountNonce() uint64 {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetPrice() *big.Int {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetGasLimit() uint64 {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetRecipient() *common.Address {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetAmount() *big.Int {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetHash() *common.Hash {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) SetHash(hash *common.Hash) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) SetSignature(signatures TxSignatures) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) RawSignatureValues() TxSignatures {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) ValidateSignature() bool {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) RecoverAddress(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) (common.Address, error) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) RecoverPubkey(txhash common.Hash, homestead bool, vfunc func(*big.Int) *big.Int) ([]*ecdsa.PublicKey, error) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) ChainId() *big.Int {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) Equal(a TxInternalData) bool {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) SerializeForSign() []interface{} {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) SenderTxHash() common.Hash {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) IsLegacyTransaction() bool {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetRoleTypeForValidation() accountkey.RoleType {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) String() string {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) MakeRPCOutput() map[string]interface{} {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) MarshalJSON() ([]byte, error) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) UnmarshalJSON(bytes []byte) error {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) setSignatureValues(chainID, v, r, s *big.Int) {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) GetAccessList() AccessList {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) gasTipCap() *big.Int {
	panic("implement me")
}

func (t *TxInternalDataDynamicFee) gasFeeCap() *big.Int {
	panic("implement me")
}



