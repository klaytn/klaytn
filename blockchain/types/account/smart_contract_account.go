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

package account

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

// SmartContractAccount represents a smart contract account containing
// storage root and code hash.
type SmartContractAccount struct {
	*AccountCommon
	storageRoot common.ExtHash // merkle root of the storage trie
	codeHash    common.ExtHash
	codeInfo    params.CodeInfo // consists of two information, vmVersion and codeFormat
}

// smartContractAccountSerializable is an internal data structure for RLP serialization.
// This structure inherits accountCommonSerializable.
type smartContractAccountSerializable struct {
	CommonSerializable *accountCommonSerializable
	StorageRoot        common.ExtHash
	CodeHash           common.ExtHash
	CodeInfo           params.CodeInfo
}

type smartContractAccountSerializableJSON struct {
	Nonce         uint64                           `json:"nonce"`
	Balance       *hexutil.Big                     `json:"balance"`
	HumanReadable bool                             `json:"humanReadable"`
	Key           *accountkey.AccountKeySerializer `json:"key"`
	StorageRoot   common.ExtHash                   `json:"storageRoot"`
	CodeHash      common.ExtHash                   `json:"codeHash"`
	CodeFormat    params.CodeFormat                `json:"codeFormat"`
	VmVersion     params.VmVersion                 `json:"vmVersion"`
}

func newSmartContractAccount() *SmartContractAccount {
	return &SmartContractAccount{
		newAccountCommon(),
		common.InitExtHash(),
		common.BytesToRootExtHash(emptyCodeHash),
		params.CodeInfo(0),
	}
}

func newSmartContractAccountWithMap(values map[AccountValueKeyType]interface{}) *SmartContractAccount {
	sca := &SmartContractAccount{
		newAccountCommonWithMap(values),
		common.InitExtHash(),
		common.BytesToRootExtHash(emptyCodeHash),
		params.CodeInfo(0),
	}

	if v, ok := values[AccountValueKeyStorageRoot].(common.ExtHash); ok {
		sca.storageRoot = v
	}

	if v, ok := values[AccountValueKeyCodeHash].([]byte); ok {
		sca.codeHash = common.BytesToRootExtHash(v)
	}

	if v, ok := values[AccountValueKeyCodeInfo].(params.CodeInfo); ok {
		sca.codeInfo = v
	}

	return sca
}

func newSmartContractAccountSerializable() *smartContractAccountSerializable {
	return &smartContractAccountSerializable{
		CommonSerializable: newAccountCommonSerializable(),
	}
}

func (sca *SmartContractAccount) toSerializable() *smartContractAccountSerializable {
	return &smartContractAccountSerializable{
		CommonSerializable: sca.AccountCommon.toSerializable(),
		StorageRoot:        sca.storageRoot,
		CodeHash:           sca.codeHash,
		CodeInfo:           sca.codeInfo,
	}
}

func (sca *SmartContractAccount) fromSerializable(o *smartContractAccountSerializable) {
	sca.AccountCommon.fromSerializable(o.CommonSerializable)
	sca.storageRoot = o.StorageRoot
	sca.codeHash = o.CodeHash
	sca.codeInfo = o.CodeInfo
}

func (sca *SmartContractAccount) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, sca.toSerializable())
}

func (sca *SmartContractAccount) DecodeRLP(s *rlp.Stream) error {
	serialized := &smartContractAccountSerializable{
		newAccountCommonSerializable(),
		common.InitExtHash(),
		common.InitExtHash(),
		params.CodeInfo(0),
	}

	if err := s.Decode(serialized); err != nil {
		return err
	}

	sca.fromSerializable(serialized)

	return nil
}

func (sca *SmartContractAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(&smartContractAccountSerializableJSON{
		Nonce:         sca.nonce,
		Balance:       (*hexutil.Big)(sca.balance),
		HumanReadable: sca.humanReadable,
		Key:           accountkey.NewAccountKeySerializerWithAccountKey(sca.key),
		StorageRoot:   sca.storageRoot,
		CodeHash:      sca.codeHash,
		CodeFormat:    sca.codeInfo.GetCodeFormat(),
		VmVersion:     sca.codeInfo.GetVmVersion(),
	})
}

func (sca *SmartContractAccount) UnmarshalJSON(b []byte) error {
	serialized := &smartContractAccountSerializableJSON{}

	if err := json.Unmarshal(b, serialized); err != nil {
		return err
	}

	sca.nonce = serialized.Nonce
	sca.balance = (*big.Int)(serialized.Balance)
	sca.humanReadable = serialized.HumanReadable
	sca.key = serialized.Key.GetKey()
	sca.storageRoot = serialized.StorageRoot
	sca.codeHash = serialized.CodeHash
	sca.codeInfo = params.NewCodeInfo(serialized.CodeFormat, serialized.VmVersion)

	return nil
}

func (sca *SmartContractAccount) Type() AccountType {
	return SmartContractAccountType
}

func (sca *SmartContractAccount) GetStorageRoot() common.ExtHash {
	return sca.storageRoot
}

func (sca *SmartContractAccount) GetCodeHash() common.ExtHash {
	return sca.codeHash
}

func (sca *SmartContractAccount) GetCodeFormat() params.CodeFormat {
	return sca.codeInfo.GetCodeFormat()
}

func (sca *SmartContractAccount) GetVmVersion() params.VmVersion {
	return sca.codeInfo.GetVmVersion()
}

func (sca *SmartContractAccount) SetStorageRoot(h common.ExtHash) {
	sca.storageRoot = h
}

func (sca *SmartContractAccount) SetCodeHash(h common.ExtHash) {
	sca.codeHash = h
}

func (sca *SmartContractAccount) SetCodeInfo(ci params.CodeInfo) {
	sca.codeInfo = ci
}

func (sca *SmartContractAccount) Empty() bool {
	return sca.nonce == 0 && sca.balance.Sign() == 0 && bytes.Equal(sca.codeHash.ToHash().Bytes(), emptyCodeHash)
}

func (sca *SmartContractAccount) UpdateKey(newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	return ErrAccountKeyNotModifiable
}

func (sca *SmartContractAccount) Equal(a Account) bool {
	sca2, ok := a.(*SmartContractAccount)
	if !ok {
		return false
	}

	return sca.AccountCommon.Equal(sca2.AccountCommon) &&
		sca.storageRoot == sca2.storageRoot &&
		bytes.Equal(sca.codeHash.ToHash().Bytes(), sca2.codeHash.ToHash().Bytes()) &&
		sca.codeInfo == sca2.codeInfo
}

func (sca *SmartContractAccount) DeepCopy() Account {
	return &SmartContractAccount{
		AccountCommon: sca.AccountCommon.DeepCopy(),
		storageRoot:   sca.storageRoot,
		codeHash:      sca.codeHash,
		codeInfo:      sca.codeInfo,
	}
}

func (sca *SmartContractAccount) TransCopy() AccountLH {
	return &SmartContractAccountLH{
		AccountCommon: sca.AccountCommon.DeepCopy(),
		storageRoot:   sca.storageRoot.ToHash(),
		codeHash:      sca.codeHash.ToHash().Bytes(),
		codeInfo:      sca.codeInfo,
	}
}

func (sca *SmartContractAccount) String() string {
	return fmt.Sprintf(`Common:%s
	StorageRoot: %s
	CodeHash: %s
	CodeInfo: %s`,
		sca.AccountCommon.String(),
		sca.storageRoot.String(),
		sca.codeHash.String(),
		sca.codeInfo.String())
}
