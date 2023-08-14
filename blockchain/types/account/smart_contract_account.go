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

// type SmartContractAccount:                 Primary in-memory representation. storageRoot is ExtHash
// type smartContractAccountSerializable:     RLP encoding spec. StorageRoot is Hash
// type smartContractAccountSerializableExt:  RLP encoding spec. StorageRoot is ExtHash
// type smartContractAccountSerializableJSON: JSON encoding spec. StorageRoot is Hash

// SmartContractAccount represents a smart contract account containing
// storage root and code hash.
type SmartContractAccount struct {
	*AccountCommon
	storageRoot common.ExtHash // merkle root plus optional sequence of the storage trie
	codeHash    []byte
	codeInfo    params.CodeInfo // consists of two information, vmVersion and codeFormat
}

// smartContractAccountSerializable is an internal data structure for RLP serialization.
// This structure inherits accountCommonSerializable.
type smartContractAccountSerializable struct {
	CommonSerializable *accountCommonSerializable
	StorageRoot        common.Hash
	CodeHash           []byte
	CodeInfo           params.CodeInfo
}

// smartContractAccountSerializableExt is an internal data structure for RLP serialization.
// This structure inherits accountCommonSerializable.
// nolint: maligned  // Because it is a temporary struct, memory footprint is not important.
type smartContractAccountSerializableExt struct {
	CommonSerializable *accountCommonSerializable
	StorageRoot        common.ExtHash
	CodeHash           []byte
	CodeInfo           params.CodeInfo
}

type smartContractAccountSerializableJSON struct {
	Nonce         uint64                           `json:"nonce"`
	Balance       *hexutil.Big                     `json:"balance"`
	HumanReadable bool                             `json:"humanReadable"`
	Key           *accountkey.AccountKeySerializer `json:"key"`
	StorageRoot   common.Hash                      `json:"storageRoot"`
	CodeHash      []byte                           `json:"codeHash"`
	CodeFormat    params.CodeFormat                `json:"codeFormat"`
	VmVersion     params.VmVersion                 `json:"vmVersion"`
}

func newSmartContractAccount() *SmartContractAccount {
	return &SmartContractAccount{
		newAccountCommon(),
		common.ExtHash{},
		emptyCodeHash,
		params.CodeInfo(0),
	}
}

func newSmartContractAccountWithMap(values map[AccountValueKeyType]interface{}) *SmartContractAccount {
	sca := &SmartContractAccount{
		newAccountCommonWithMap(values),
		common.ExtHash{},
		emptyCodeHash,
		params.CodeInfo(0),
	}

	if v, ok := values[AccountValueKeyStorageRoot].(common.Hash); ok {
		sca.storageRoot = v.ExtendLegacy()
	}
	if v, ok := values[AccountValueKeyStorageRoot].(common.ExtHash); ok {
		sca.storageRoot = v
	}

	if v, ok := values[AccountValueKeyCodeHash].([]byte); ok {
		sca.codeHash = v
	}

	if v, ok := values[AccountValueKeyCodeInfo].(params.CodeInfo); ok {
		sca.codeInfo = v
	}

	return sca
}

func (sca *SmartContractAccount) toSerializable() *smartContractAccountSerializable {
	return &smartContractAccountSerializable{
		CommonSerializable: sca.AccountCommon.toSerializable(),
		StorageRoot:        sca.storageRoot.Unextend(),
		CodeHash:           sca.codeHash,
		CodeInfo:           sca.codeInfo,
	}
}

func (sca *SmartContractAccount) toSerializableExt() *smartContractAccountSerializableExt {
	return &smartContractAccountSerializableExt{
		CommonSerializable: sca.AccountCommon.toSerializable(),
		StorageRoot:        sca.storageRoot,
		CodeHash:           sca.codeHash,
		CodeInfo:           sca.codeInfo,
	}
}

func (sca *SmartContractAccount) fromSerializable(o *smartContractAccountSerializable) {
	sca.AccountCommon.fromSerializable(o.CommonSerializable)
	sca.storageRoot = o.StorageRoot.ExtendLegacy()
	sca.codeHash = o.CodeHash
	sca.codeInfo = o.CodeInfo
}

func (sca *SmartContractAccount) fromSerializableExt(o *smartContractAccountSerializableExt) {
	sca.AccountCommon.fromSerializable(o.CommonSerializable)
	sca.storageRoot = o.StorageRoot
	sca.codeHash = o.CodeHash
	sca.codeInfo = o.CodeInfo
}

func (sca *SmartContractAccount) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, sca.toSerializable())
}

func (sca *SmartContractAccount) EncodeRLPExt(w io.Writer) error {
	if sca.storageRoot.IsLegacy() {
		return rlp.Encode(w, sca.toSerializable())
	} else {
		return rlp.Encode(w, sca.toSerializableExt())
	}
}

func (sca *SmartContractAccount) DecodeRLP(s *rlp.Stream) error {
	savedStream, err := s.Raw()
	if err != nil {
		return err
	}

	// s.Raw() has consumed the stream. Refill with original data.
	s.Reset(bytes.NewReader(savedStream), 0)
	// Try decode into smartContractAccountSerializableExt
	serializedExt := &smartContractAccountSerializableExt{
		CommonSerializable: newAccountCommonSerializable(),
	}
	if err := s.Decode(serializedExt); err == nil {
		sca.fromSerializableExt(serializedExt)
		return nil
	}

	// s.Decode() may have consumed the stream. Refill with original data.
	s.Reset(bytes.NewReader(savedStream), 0)
	// Retry with smartContractAccountSerializable
	serialized := &smartContractAccountSerializable{
		CommonSerializable: newAccountCommonSerializable(),
	}
	if err := s.Decode(serialized); err == nil {
		sca.fromSerializable(serialized)
		return nil
	} else {
		return err
	}
}

func (sca *SmartContractAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(&smartContractAccountSerializableJSON{
		Nonce:         sca.nonce,
		Balance:       (*hexutil.Big)(sca.balance),
		HumanReadable: sca.humanReadable,
		Key:           accountkey.NewAccountKeySerializerWithAccountKey(sca.key),
		StorageRoot:   sca.storageRoot.Unextend(), // Unextend for API compatibility
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
	sca.storageRoot = serialized.StorageRoot.ExtendLegacy() // API inputs should contain merkle hash
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

func (sca *SmartContractAccount) GetCodeHash() []byte {
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

func (sca *SmartContractAccount) SetCodeHash(h []byte) {
	sca.codeHash = h
}

func (sca *SmartContractAccount) SetCodeInfo(ci params.CodeInfo) {
	sca.codeInfo = ci
}

func (sca *SmartContractAccount) Empty() bool {
	return sca.nonce == 0 && sca.balance.Sign() == 0 && bytes.Equal(sca.codeHash, emptyCodeHash)
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
		bytes.Equal(sca.codeHash, sca2.codeHash) &&
		sca.codeInfo == sca2.codeInfo
}

func (sca *SmartContractAccount) DeepCopy() Account {
	return &SmartContractAccount{
		AccountCommon: sca.AccountCommon.DeepCopy(),
		storageRoot:   sca.storageRoot,
		codeHash:      common.CopyBytes(sca.codeHash),
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
		common.Bytes2Hex(sca.codeHash),
		sca.codeInfo.String())
}
