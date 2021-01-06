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
	"io"

	"github.com/klaytn/klaytn/rlp"
)

// AccountSerializer serializes an Account object using RLP/JSON.
type AccountSerializer struct {
	accType AccountType
	account Account
}

// accountJSON is an internal data structure for JSON serialization.
type accountJSON struct {
	AccType AccountType     `json:"accType"`
	Account json.RawMessage `json:"account"`
}

// NewAccountSerializer creates a new AccountSerializer object with default values.
// This returned object will be used for decoding.
func NewAccountSerializer() *AccountSerializer {
	return &AccountSerializer{}
}

// NewAccountSerializerWithAccount creates a new AccountSerializer object with the given account.
func NewAccountSerializerWithAccount(a Account) *AccountSerializer {
	return &AccountSerializer{a.Type(), a}
}

func (ser *AccountSerializer) EncodeRLP(w io.Writer) error {
	// If it is a LegacyAccount object, do not encode the account type.
	if ser.accType == LegacyAccountType {
		return rlp.Encode(w, ser.account.(*LegacyAccount))
	}

	if err := rlp.Encode(w, ser.accType); err != nil {
		return err
	}

	return rlp.Encode(w, ser.account)
}

func (ser *AccountSerializer) GetAccount() Account {
	return ser.account
}

func (ser *AccountSerializer) DecodeRLP(s *rlp.Stream) error {
	if err := s.Decode(&ser.accType); err != nil {
		// fallback to decoding a LegacyAccount object.
		acc := newLegacyAccount()
		if err := s.Decode(acc); err != nil {
			return err
		}
		ser.accType = LegacyAccountType
		ser.account = acc
		return nil
	}

	var err error
	ser.account, err = NewAccountWithType(ser.accType)
	if err != nil {
		return err
	}

	return s.Decode(ser.account)
}

func (ser *AccountSerializer) MarshalJSON() ([]byte, error) {
	// if it is a legacyAccount object, do not marshal the account type.
	if ser.accType == LegacyAccountType {
		return json.Marshal(ser.account)
	}
	b, err := json.Marshal(ser.account)
	if err != nil {
		return nil, err
	}

	return json.Marshal(&accountJSON{ser.accType, b})
}

func (ser *AccountSerializer) UnmarshalJSON(b []byte) error {
	dec := &accountJSON{}

	if err := json.Unmarshal(b, dec); err != nil {
		return err
	}

	if len(dec.Account) == 0 {
		// fallback to unmarshal a LegacyAccount object.
		acc := newLegacyAccount()
		if err := json.Unmarshal(b, acc); err != nil {
			return err
		}
		ser.accType = LegacyAccountType
		ser.account = acc

		return nil

	}

	ser.accType = dec.AccType

	var err error
	ser.account, err = NewAccountWithType(ser.accType)
	if err != nil {
		return err
	}

	return json.Unmarshal(dec.Account, ser.account)
}
