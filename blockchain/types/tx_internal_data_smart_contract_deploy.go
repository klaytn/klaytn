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
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/kerrors"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

// TxInternalDataSmartContractDeploy represents a transaction creating a smart contract.
type TxInternalDataSmartContractDeploy struct {
	AccountNonce  uint64
	Price         *big.Int
	GasLimit      uint64
	Recipient     *common.Address `rlp:"nil"`
	Amount        *big.Int
	From          common.Address
	Payload       []byte
	HumanReadable bool
	CodeFormat    params.CodeFormat

	TxSignatures

	// This is only used when marshaling to JSON.
	Hash *common.Hash `json:"hash" rlp:"-"`
}

type TxInternalDataSmartContractDeployJSON struct {
	Type          TxType           `json:"typeInt"`
	TypeStr       string           `json:"type"`
	AccountNonce  hexutil.Uint64   `json:"nonce"`
	Price         *hexutil.Big     `json:"gasPrice"`
	GasLimit      hexutil.Uint64   `json:"gas"`
	Recipient     *common.Address  `json:"to"`
	Amount        *hexutil.Big     `json:"value"`
	From          common.Address   `json:"from"`
	Payload       hexutil.Bytes    `json:"input"`
	HumanReadable bool             `json:"humanReadable"`
	CodeFormat    hexutil.Uint     `json:"codeFormat"`
	TxSignatures  TxSignaturesJSON `json:"signatures"`
	Hash          *common.Hash     `json:"hash"`
}

func newTxInternalDataSmartContractDeploy() *TxInternalDataSmartContractDeploy {
	h := common.Hash{}
	return &TxInternalDataSmartContractDeploy{
		Price:  new(big.Int),
		Amount: new(big.Int),
		Hash:   &h,
	}
}

func newTxInternalDataSmartContractDeployWithMap(values map[TxValueKeyType]interface{}) (*TxInternalDataSmartContractDeploy, error) {
	t := newTxInternalDataSmartContractDeploy()

	if v, ok := values[TxValueKeyNonce].(uint64); ok {
		t.AccountNonce = v
		delete(values, TxValueKeyNonce)
	} else {
		return nil, errValueKeyNonceMustUint64
	}

	if v, ok := values[TxValueKeyTo].(*common.Address); ok {
		t.Recipient = v
		delete(values, TxValueKeyTo)
	} else {
		return nil, errValueKeyToMustAddressPointer
	}

	if v, ok := values[TxValueKeyAmount].(*big.Int); ok {
		t.Amount.Set(v)
		delete(values, TxValueKeyAmount)
	} else {
		return nil, errValueKeyAmountMustBigInt
	}

	if v, ok := values[TxValueKeyGasLimit].(uint64); ok {
		t.GasLimit = v
		delete(values, TxValueKeyGasLimit)
	} else {
		return nil, errValueKeyGasLimitMustUint64
	}

	if v, ok := values[TxValueKeyGasPrice].(*big.Int); ok {
		t.Price.Set(v)
		delete(values, TxValueKeyGasPrice)
	} else {
		return nil, errValueKeyGasPriceMustBigInt
	}

	if v, ok := values[TxValueKeyFrom].(common.Address); ok {
		t.From = v
		delete(values, TxValueKeyFrom)
	} else {
		return nil, errValueKeyFromMustAddress
	}

	if v, ok := values[TxValueKeyData].([]byte); ok {
		t.Payload = v
		delete(values, TxValueKeyData)
	} else {
		return nil, errValueKeyDataMustByteSlice
	}

	if v, ok := values[TxValueKeyHumanReadable].(bool); ok {
		t.HumanReadable = v
		delete(values, TxValueKeyHumanReadable)
	} else {
		return nil, errValueKeyHumanReadableMustBool
	}

	if v, ok := values[TxValueKeyCodeFormat].(params.CodeFormat); ok {
		t.CodeFormat = v
		delete(values, TxValueKeyCodeFormat)
	} else {
		return nil, errValueKeyCodeFormatInvalid
	}

	if len(values) != 0 {
		for k := range values {
			logger.Warn("unnecessary key", k.String())
		}
		return nil, errUndefinedKeyRemains
	}

	return t, nil
}

func (t *TxInternalDataSmartContractDeploy) Type() TxType {
	return TxTypeSmartContractDeploy
}

func (t *TxInternalDataSmartContractDeploy) GetRoleTypeForValidation() accountkey.RoleType {
	return accountkey.RoleTransaction
}

func (t *TxInternalDataSmartContractDeploy) GetPayload() []byte {
	return t.Payload
}

func (t *TxInternalDataSmartContractDeploy) Equal(a TxInternalData) bool {
	ta, ok := a.(*TxInternalDataSmartContractDeploy)
	if !ok {
		return false
	}

	return t.AccountNonce == ta.AccountNonce &&
		t.Price.Cmp(ta.Price) == 0 &&
		t.GasLimit == ta.GasLimit &&
		equalRecipient(t.Recipient, ta.Recipient) &&
		t.Amount.Cmp(ta.Amount) == 0 &&
		t.From == ta.From &&
		bytes.Equal(t.Payload, ta.Payload) &&
		t.HumanReadable == ta.HumanReadable &&
		t.TxSignatures.equal(ta.TxSignatures) &&
		t.CodeFormat == ta.CodeFormat
}

func (t *TxInternalDataSmartContractDeploy) IsLegacyTransaction() bool {
	return false
}

func (t *TxInternalDataSmartContractDeploy) GetAccountNonce() uint64 {
	return t.AccountNonce
}

func (t *TxInternalDataSmartContractDeploy) GetPrice() *big.Int {
	return new(big.Int).Set(t.Price)
}

func (t *TxInternalDataSmartContractDeploy) GetGasLimit() uint64 {
	return t.GasLimit
}

func (t *TxInternalDataSmartContractDeploy) GetRecipient() *common.Address {
	return t.Recipient
}

func (t *TxInternalDataSmartContractDeploy) GetAmount() *big.Int {
	return new(big.Int).Set(t.Amount)
}

func (t *TxInternalDataSmartContractDeploy) GetFrom() common.Address {
	return t.From
}

func (t *TxInternalDataSmartContractDeploy) GetHash() *common.Hash {
	return t.Hash
}

func (t *TxInternalDataSmartContractDeploy) GetCodeFormat() params.CodeFormat {
	return t.CodeFormat
}

func (t *TxInternalDataSmartContractDeploy) SetHash(h *common.Hash) {
	t.Hash = h
}

func (t *TxInternalDataSmartContractDeploy) SetSignature(s TxSignatures) {
	t.TxSignatures = s
}

func (t *TxInternalDataSmartContractDeploy) String() string {
	var to common.Address
	if t.Recipient != nil {
		to = *t.Recipient
	} else {
		to = crypto.CreateAddress(t.From, t.AccountNonce)
	}
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
	Signature:     %s
	Data:          %x
	HumanReadable: %v
	CodeFormat:    %s
	Hex:           %x
`,
		tx.Hash(),
		t.Type().String(),
		t.From.String(),
		to.String(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Amount,
		t.TxSignatures.string(),
		common.Bytes2Hex(t.Payload),
		t.HumanReadable,
		t.CodeFormat.String(),
		enc)

}

func (t *TxInternalDataSmartContractDeploy) IntrinsicGas(currentBlockNumber uint64) (uint64, error) {
	gas := params.TxGasContractCreation

	if t.HumanReadable {
		gas += params.TxGasHumanReadable
	}

	gasPayloadWithGas, err := IntrinsicGasPayload(gas, t.Payload)
	if err != nil {
		return 0, err
	}

	return gasPayloadWithGas, nil
}

func (t *TxInternalDataSmartContractDeploy) SerializeForSignToBytes() []byte {
	b, _ := rlp.EncodeToBytes(struct {
		Txtype        TxType
		AccountNonce  uint64
		Price         *big.Int
		GasLimit      uint64
		Recipient     *common.Address
		Amount        *big.Int
		From          common.Address
		Payload       []byte
		HumanReadable bool
		CodeFormat    params.CodeFormat
	}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.From,
		t.Payload,
		t.HumanReadable,
		t.CodeFormat,
	})

	return b
}

func (t *TxInternalDataSmartContractDeploy) SerializeForSign() []interface{} {
	return []interface{}{
		t.Type(),
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.From,
		t.Payload,
		t.HumanReadable,
		t.CodeFormat,
	}
}

func (t *TxInternalDataSmartContractDeploy) SenderTxHash() common.Hash {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, t.Type())
	rlp.Encode(hw, []interface{}{
		t.AccountNonce,
		t.Price,
		t.GasLimit,
		t.Recipient,
		t.Amount,
		t.From,
		t.Payload,
		t.HumanReadable,
		t.CodeFormat,
		t.TxSignatures,
	})

	h := common.Hash{}

	hw.Sum(h[:0])

	return h
}

func (t *TxInternalDataSmartContractDeploy) Validate(stateDB StateDB, currentBlockNumber uint64) error {
	var to common.Address
	if t.Recipient != nil {
		return kerrors.ErrInvalidContractAddress
	} else {
		to = crypto.CreateAddress(t.From, t.AccountNonce)
	}
	if common.IsPrecompiledContractAddress(to) {
		return kerrors.ErrPrecompiledContractAddress
	}
	if t.HumanReadable {
		return kerrors.ErrHumanReadableNotSupported
	}
	// Fail if the codeFormat is invalid.
	if !t.CodeFormat.Validate() {
		return kerrors.ErrInvalidCodeFormat
	}
	return t.ValidateMutableValue(stateDB, currentBlockNumber)
}

func (t *TxInternalDataSmartContractDeploy) ValidateMutableValue(stateDB StateDB, currentBlockNumber uint64) error {
	// Fail if the address is already created.
	if t.Recipient != nil && stateDB.Exist(*t.Recipient) {
		return kerrors.ErrAccountAlreadyExists
	}
	return nil
}

func (t *TxInternalDataSmartContractDeploy) FillContractAddress(from common.Address, r *Receipt) {
	if t.Recipient == nil {
		r.ContractAddress = crypto.CreateAddress(from, t.AccountNonce)
	} else {
		r.ContractAddress = *t.Recipient
	}
}

func (t *TxInternalDataSmartContractDeploy) Execute(sender ContractRef, vm VM, stateDB StateDB, currentBlockNumber uint64, gas uint64, value *big.Int) (ret []byte, usedGas uint64, err error) {
	///////////////////////////////////////////////////////
	// OpcodeComputationCostLimit: The below code is commented and will be usd for debugging purposes.
	//start := time.Now()
	//defer func() {
	//	elapsed := time.Since(start)
	//	logger.Debug("[TxInternalDataSmartContractDeploy] EVM execution done", "elapsed", elapsed)
	//}()
	///////////////////////////////////////////////////////
	// Sender's nonce will be increased in '`vm.Create()` or `vm.CreateWithAddress()`
	if t.Recipient == nil {
		ret, _, usedGas, err = vm.Create(sender, t.Payload, gas, value, t.CodeFormat)
	} else {
		ret, _, usedGas, err = vm.CreateWithAddress(sender, t.Payload, gas, value, *t.Recipient, t.HumanReadable, t.CodeFormat)
	}
	return ret, usedGas, err
}

func (t *TxInternalDataSmartContractDeploy) MakeRPCOutput() map[string]interface{} {
	return map[string]interface{}{
		"typeInt":       t.Type(),
		"type":          t.Type().String(),
		"gas":           hexutil.Uint64(t.GasLimit),
		"gasPrice":      (*hexutil.Big)(t.Price),
		"nonce":         hexutil.Uint64(t.AccountNonce),
		"to":            t.Recipient,
		"value":         (*hexutil.Big)(t.Amount),
		"input":         hexutil.Bytes(t.Payload),
		"humanReadable": t.HumanReadable,
		"codeFormat":    hexutil.Uint(t.CodeFormat),
		"signatures":    t.TxSignatures.ToJSON(),
	}
}

func (t *TxInternalDataSmartContractDeploy) MarshalJSON() ([]byte, error) {
	return json.Marshal(TxInternalDataSmartContractDeployJSON{
		t.Type(),
		t.Type().String(),
		(hexutil.Uint64)(t.AccountNonce),
		(*hexutil.Big)(t.Price),
		(hexutil.Uint64)(t.GasLimit),
		t.Recipient,
		(*hexutil.Big)(t.Amount),
		t.From,
		t.Payload,
		t.HumanReadable,
		(hexutil.Uint)(t.CodeFormat),
		t.TxSignatures.ToJSON(),
		t.Hash,
	})
}

func (t *TxInternalDataSmartContractDeploy) UnmarshalJSON(b []byte) error {
	js := &TxInternalDataSmartContractDeployJSON{}
	if err := json.Unmarshal(b, js); err != nil {
		return err
	}

	t.AccountNonce = uint64(js.AccountNonce)
	t.Price = (*big.Int)(js.Price)
	t.GasLimit = uint64(js.GasLimit)
	t.Recipient = js.Recipient
	t.Amount = (*big.Int)(js.Amount)
	t.From = js.From
	t.Payload = js.Payload
	t.HumanReadable = js.HumanReadable
	t.CodeFormat = params.CodeFormat(js.CodeFormat)
	t.TxSignatures = js.TxSignatures.ToTxSignatures()
	t.Hash = js.Hash

	return nil
}
