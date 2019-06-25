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
	"fmt"
)

// ExternallyOwnedAccount represents a Klaytn account used by a user.
type ExternallyOwnedAccount struct {
	*AccountCommon
}

// newExternallyOwnedAccount creates an ExternallyOwnedAccount object with default values.
func newExternallyOwnedAccount() *ExternallyOwnedAccount {
	return &ExternallyOwnedAccount{
		newAccountCommon(),
	}
}

// newExternallyOwnedAccountWithMap creates an ExternallyOwnedAccount object initialized with the given values.
func newExternallyOwnedAccountWithMap(values map[AccountValueKeyType]interface{}) *ExternallyOwnedAccount {
	return &ExternallyOwnedAccount{
		newAccountCommonWithMap(values),
	}
}

func (e *ExternallyOwnedAccount) Type() AccountType {
	return ExternallyOwnedAccountType
}

func (e *ExternallyOwnedAccount) Dump() {
	fmt.Println(e.String())
}

func (e *ExternallyOwnedAccount) String() string {
	return fmt.Sprintf("EOA: %s", e.AccountCommon.String())
}

func (e *ExternallyOwnedAccount) DeepCopy() Account {
	return &ExternallyOwnedAccount{
		AccountCommon: e.AccountCommon.DeepCopy(),
	}
}

func (e *ExternallyOwnedAccount) Equal(a Account) bool {
	e2, ok := a.(*ExternallyOwnedAccount)
	if !ok {
		return false
	}

	return e.AccountCommon.Equal(e2.AccountCommon)
}
