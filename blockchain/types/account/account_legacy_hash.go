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
	"encoding/json"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

// all code decode for legacy hash ( 32 byte hash )

//account.go
//
// NewAccountWithType creates an Account object with the given type.
func NewAccountLHWithType(t AccountType) (Account, error) {
	switch t {
	case LegacyAccountType:
		return newLegacyAccountLH(), nil
	case ExternallyOwnedAccountType:
		//return newExternallyOwnedAccountLH(), nil
		return newExternallyOwnedAccount(), nil
		// or panic()?
	case SmartContractAccountType:
		return newSmartContractAccountLH(), nil
	}

	return nil, ErrUndefinedAccountType
}

// Account is the Klaytn consensus representation of accounts.
// These objects are stored in the main account trie.
type AccountLH interface {
	Type() AccountType

	GetNonce() uint64
	GetBalance() *big.Int
	GetHumanReadable() bool

	SetNonce(n uint64)
	SetBalance(b *big.Int)
	SetHumanReadable(b bool)

	// UpdateKey updates the account's key with the given key.
	UpdateKey(newKey accountkey.AccountKey, currentBlockNumber uint64) error

	// Empty returns whether the account is considered empty.
	// The "empty" account may be defined differently depending on the actual account type.
	// An example of an empty account could be described as the one that satisfies the following conditions:
	// - nonce is zero
	// - balance is zero
	// - codeHash is the same as emptyCodeHash
	Empty() bool

	// Equal returns true if all the attributes are exactly same. Otherwise, returns false.
	Equal(Account) bool

	// DeepCopy copies all the attributes.
	DeepCopy() Account

	// String returns all attributes of this object as a string.
	String() string
}

//account_serializer.go
//
// AccountSerializer serializes an Account object using RLP/JSON.
type AccountLHSerializer struct {
	accType AccountType
	account AccountLH
}

// accountJSON is an internal data structure for JSON serialization.
type accountLHJSON struct {
	AccType AccountType     `json:"accType"`
	Account json.RawMessage `json:"account"`
}

// NewAccountSerializer creates a new AccountSerializer object with default values.
// This returned object will be used for decoding.
func NewAccountLHSerializer() *AccountLHSerializer {
	return &AccountLHSerializer{}
}

func (ser *AccountLHSerializer) DecodeRLP(s *rlp.Stream) error {
	if err := s.Decode(&ser.accType); err != nil {
		// fallback to decoding a LegacyAccount object.
		acc := newLegacyAccountLH()
		if err := s.Decode(acc); err != nil {
			return err
		}
		ser.accType = LegacyAccountType
		ser.account = acc
		return nil
	}

	var err error
	ser.account, err = NewAccountLHWithType(ser.accType)
	if err != nil {
		return err
	}

	return s.Decode(ser.account)
}

func (ser *AccountLHSerializer) Copy() (eser *AccountSerializer) {
	return &AccountSerializer{
		accType: ser.accType,
		account: ser.account.DeepCopy(),
	}
	/*
		        switch eser.accType {
		        case LegacyAccountType:
				eser.account = ser.account.DeepCopy()
		        case ExternallyOwnedAccountType:
				eser.account = ser.account.DeepCopy()
		        case SmartContractAccountType:
				eser.account = ser.account.DeepCopy()
		        }*/
}

/*
// account_common.go
//
// AccountCommon represents the common data structure of a Klaytn account.
type AccountLHCommon struct {
        nonce         uint64
        balance       *big.Int
        humanReadable bool
        key           accountkey.AccountKey
}

// accountCommonSerializable is an internal data structure for RLP serialization.
// This object is required due to AccountKey.
// AccountKey is an interface and it requires to use AccountKeySerializer for serialization.
type accountLHCommonSerializable struct {
        Nonce         uint64
        Balance       *big.Int
        HumanReadable bool
        Key           *accountkey.AccountKeySerializer
}

// newAccountCommon creates an AccountCommon object with default values.
func newAccountLHCommon() *AccountLHCommon {
        return &AccountLHCommon{
                nonce:         0,
                balance:       new(big.Int),
                humanReadable: false,
                key:           accountkey.NewAccountKeyLegacy(),
        }
}

func (e *AccountLHCommon) DecodeRLP(s *rlp.Stream) error {
        serialized := newAccountLHCommonSerializable()

        if err := s.Decode(serialized); err != nil {
                return err
        }
        e.fromSerializable(serialized)
        return nil
}

func newAccountLHCommonSerializable() *accountLHCommonSerializable {
        return &accountLHCommonSerializable{
                Balance: new(big.Int),
                Key:     accountkey.NewAccountKeySerializer(),
        }
}

// fromSerializable updates its values from the given accountCommonSerializable object.
func (e *AccountLHCommon) fromSerializable(o *accountLHCommonSerializable) {
        e.nonce = o.Nonce
        e.balance = o.Balance
        e.humanReadable = o.HumanReadable
        e.key = o.Key.GetKey()
}
*/

// legacy_account.go
//
// LegacyAccount is the Klaytn consensus representation of legacy accounts.
// These objects are stored in the main account trie.
type LegacyAccountLH struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
}

// newLegacyAccount returns a LegacyAccount object whose all
// attributes are initialized.
// This object is used when an account is created.
// Refer to StateDB.createObject().
func newLegacyAccountLH() *LegacyAccountLH {
	logger.CritWithStack("Legacy account is deprecated.")
	return &LegacyAccountLH{
		0, new(big.Int), common.Hash{}, emptyCodeHash,
	}
}

func (a *LegacyAccountLH) DeepCopy() Account {
	return &LegacyAccount{
		Nonce:    a.Nonce,
		Balance:  a.Balance,
		Root:     common.BytesLegacyToExtHash(a.Root.Bytes()),
		CodeHash: a.CodeHash,
	}
}

func (a *LegacyAccountLH) Type() AccountType       { return 0 }
func (a *LegacyAccountLH) GetNonce() uint64        { return 0 }
func (a *LegacyAccountLH) GetBalance() *big.Int    { return big.NewInt(0) }
func (a *LegacyAccountLH) GetHumanReadable() bool  { return false }
func (a *LegacyAccountLH) SetNonce(n uint64)       { return }
func (a *LegacyAccountLH) SetBalance(b *big.Int)   { return }
func (a *LegacyAccountLH) SetHumanReadable(b bool) { return }
func (a *LegacyAccountLH) UpdateKey(newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	return nil
}
func (a *LegacyAccountLH) Empty() bool        { return false }
func (a *LegacyAccountLH) Equal(Account) bool { return false }

//func (a *LegacyAccountLH) DeepCopy() Account { return nil }
func (a *LegacyAccountLH) String() string { return "Not Implemented" }

// externally_owned_account.go
//
// ExternallyOwnedAccount represents a Klaytn account used by a user.
/*type ExternallyOwnedAccountLH struct {
        *AccountLHCommon
}

// newExternallyOwnedAccount creates an ExternallyOwnedAccount object with default values.
func newExternallyOwnedAccountLH() *ExternallyOwnedAccountLH {
        return &ExternallyOwnedAccountLH{
                newAccountLHCommon(),
        }
}*/

// smart_contract_account.go
//
// SmartContractAccount represents a smart contract account containing
// storage root and code hash.
type SmartContractAccountLH struct {
	//aaa *AccountLHCommon
	*AccountCommon
	storageRoot common.Hash // merkle root of the storage trie
	codeHash    []byte
	codeInfo    params.CodeInfo // consists of two information, vmVersion and codeFormat
}

// smartContractAccountSerializable is an internal data structure for RLP serialization.
// This structure inherits accountCommonSerializable.
type smartContractAccountLHSerializable struct {
	//aaa CommonSerializable *accountLHCommonSerializable
	CommonSerializable *accountCommonSerializable
	StorageRoot        common.Hash
	CodeHash           []byte
	CodeInfo           params.CodeInfo
}

func newSmartContractAccountLH() *SmartContractAccountLH {
	return &SmartContractAccountLH{
		//aaa newAccountLHCommon(),
		newAccountCommon(),
		//common.InitExtHash(),
		common.Hash{},
		emptyCodeHash,
		params.CodeInfo(0),
	}
}
func (sca *SmartContractAccountLH) fromSerializable(o *smartContractAccountLHSerializable) {
	//aaa sca.AccountLHCommon.fromSerializable(o.CommonSerializable)
	sca.AccountCommon.fromSerializable(o.CommonSerializable)
	sca.storageRoot = o.StorageRoot
	sca.codeHash = o.CodeHash
	sca.codeInfo = o.CodeInfo
}

func (sca *SmartContractAccountLH) DecodeRLP(s *rlp.Stream) error {
	serialized := &smartContractAccountLHSerializable{
		//aaa newAccountLHCommonSerializable(),
		newAccountCommonSerializable(),
		//common.InitExtHash(),
		common.Hash{},
		[]byte{},
		params.CodeInfo(0),
	}

	if err := s.Decode(serialized); err != nil {
		return err
	}

	sca.fromSerializable(serialized)

	return nil
}

func (a *SmartContractAccountLH) DeepCopy() Account {
	return &SmartContractAccount{
		AccountCommon: a.AccountCommon.DeepCopy(),
		storageRoot:   common.BytesLegacyToExtHash(a.storageRoot.Bytes()),
		codeHash:      a.codeHash,
		codeInfo:      a.codeInfo,
	}
}

func (a *SmartContractAccountLH) Type() AccountType       { return 0 }
func (a *SmartContractAccountLH) GetNonce() uint64        { return 0 }
func (a *SmartContractAccountLH) GetBalance() *big.Int    { return big.NewInt(0) }
func (a *SmartContractAccountLH) GetHumanReadable() bool  { return false }
func (a *SmartContractAccountLH) SetNonce(n uint64)       { return }
func (a *SmartContractAccountLH) SetBalance(b *big.Int)   { return }
func (a *SmartContractAccountLH) SetHumanReadable(b bool) { return }
func (a *SmartContractAccountLH) UpdateKey(newKey accountkey.AccountKey, currentBlockNumber uint64) error {
	return nil
}
func (a *SmartContractAccountLH) Empty() bool        { return false }
func (a *SmartContractAccountLH) Equal(Account) bool { return false }

//func (a *SmartContractAccountLH) DeepCopy() Account { return nil }
func (a *SmartContractAccountLH) String() string { return "Not implemented" }
