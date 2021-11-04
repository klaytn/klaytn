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
	"errors"
	"fmt"
	"io"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/types/accountkey"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

var (
	ErrNotEOA     = errors.New("Not a EOA")
	ErrNilAccount = errors.New("Account not set")
)

func GetInitialBalanceLimit() *big.Int {
	return new(big.Int).Mul(big.NewInt(params.KLAY), big.NewInt(params.KLAY))
}

// ExternallyOwnedAccount represents a Klaytn account used by a user.
type ExternallyOwnedAccount struct {
	*AccountCommon
	balanceLimit *big.Int
}

type externallyOwnedAccountSerializable struct {
	CommonSerializable *accountCommonSerializable
	BalanceLimit       *big.Int
}

type externallyOwnedAccountSerializableJSON struct {
	Nonce         uint64                           `json:"nonce"`
	Balance       *hexutil.Big                     `json:"balance"`
	HumanReadable bool                             `json:"humanReadable"`
	Key           *accountkey.AccountKeySerializer `json:"key"`
	BalanceLimit  *hexutil.Big                     `json:"balanceLimit"`
}

// newExternallyOwnedAccount creates an ExternallyOwnedAccount object with default values.
func newExternallyOwnedAccount() *ExternallyOwnedAccount {
	return &ExternallyOwnedAccount{
		AccountCommon: newAccountCommon(),
		balanceLimit:  GetInitialBalanceLimit(),
	}
}

// newExternallyOwnedAccountWithMap creates an ExternallyOwnedAccount object initialized with the given values.
func newExternallyOwnedAccountWithMap(values map[AccountValueKeyType]interface{}) *ExternallyOwnedAccount {
	balanceLimit := new(big.Int)
	if v, ok := values[AccountValueBalanceLimit].(*big.Int); ok {
		balanceLimit.Set(v)
	} else {
		balanceLimit.Set(GetInitialBalanceLimit())
	}

	return &ExternallyOwnedAccount{
		AccountCommon: newAccountCommonWithMap(values),
		balanceLimit:  balanceLimit,
	}
}

func (e *ExternallyOwnedAccount) toSerializable() *externallyOwnedAccountSerializable {
	return &externallyOwnedAccountSerializable{
		CommonSerializable: e.AccountCommon.toSerializable(),
		BalanceLimit:       e.balanceLimit,
	}
}

func (e *ExternallyOwnedAccount) fromSerializable(o *externallyOwnedAccountSerializable) {
	e.AccountCommon.fromSerializable(o.CommonSerializable)
	e.balanceLimit = o.BalanceLimit
}

func (e *ExternallyOwnedAccount) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, e.toSerializable())
}

func (e *ExternallyOwnedAccount) DecodeRLP(s *rlp.Stream) error {
	serialized := &externallyOwnedAccountSerializable{
		newAccountCommonSerializable(),
		new(big.Int),
	}
	if err := s.Decode(serialized); err != nil {
		return err
	}

	e.fromSerializable(serialized)

	return nil
}

func (e *ExternallyOwnedAccount) MarshalJSON() ([]byte, error) {
	return json.Marshal(&externallyOwnedAccountSerializableJSON{
		Nonce:         e.nonce,
		Balance:       (*hexutil.Big)(e.balance),
		HumanReadable: e.humanReadable,
		Key:           accountkey.NewAccountKeySerializerWithAccountKey(e.key),
		BalanceLimit:  (*hexutil.Big)(e.balanceLimit),
	})
}

func (e *ExternallyOwnedAccount) UnmarshalJSON(b []byte) error {
	serialized := &externallyOwnedAccountSerializableJSON{}

	if err := json.Unmarshal(b, serialized); err != nil {
		return err
	}

	e.nonce = serialized.Nonce
	e.balance = (*big.Int)(serialized.Balance)
	e.humanReadable = serialized.HumanReadable
	e.key = serialized.Key.GetKey()
	e.balanceLimit = (*big.Int)(serialized.BalanceLimit)

	return nil
}

func (e *ExternallyOwnedAccount) Type() AccountType {
	return ExternallyOwnedAccountType
}

func (e *ExternallyOwnedAccount) GetBalanceLimit() *big.Int {
	return new(big.Int).Set(e.balanceLimit)
}

func (e *ExternallyOwnedAccount) SetBalanceLimit(balanceLimit *big.Int) {
	e.balanceLimit.Set(balanceLimit)
}

func (e *ExternallyOwnedAccount) Dump() {
	fmt.Println(e.String())
}

func (e *ExternallyOwnedAccount) String() string {
	return fmt.Sprintf(`Common: %s
	BalanceLimit: %s`,
		e.AccountCommon.String(), e.balanceLimit.String())
}

func (e *ExternallyOwnedAccount) DeepCopy() Account {
	return &ExternallyOwnedAccount{
		AccountCommon: e.AccountCommon.DeepCopy(),
		balanceLimit:  e.balanceLimit,
	}
}

func (e *ExternallyOwnedAccount) Equal(a Account) bool {
	e2, ok := a.(*ExternallyOwnedAccount)
	if !ok {
		return false
	}

	if e.balanceLimit.Cmp(e2.balanceLimit) != 0 {
		return false
	}

	return e.AccountCommon.Equal(e2.AccountCommon)
}
