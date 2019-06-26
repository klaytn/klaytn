// Modifications Copyright 2019 The klaytn Authors
// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from core/state_transition.go (2018/06/04).
// Modified and improved for the klaytn development.

package types

import "github.com/klaytn/klaytn/common"

// ContractRef is a reference to the contract's backing object
type ContractRef interface {
	Address() common.Address
	FeePayer() common.Address
}

// AccountRefWithFeePayer implements ContractRef.
// This structure has an additional field `feePayer` compared to `AccountRef`.
type AccountRefWithFeePayer struct {
	SenderAddress   common.Address
	FeePayerAddress common.Address
}

func NewAccountRefWithFeePayer(sender common.Address, feePayer common.Address) *AccountRefWithFeePayer {
	return &AccountRefWithFeePayer{sender, feePayer}
}

func (a *AccountRefWithFeePayer) Address() common.Address  { return a.SenderAddress }
func (a *AccountRefWithFeePayer) FeePayer() common.Address { return a.FeePayerAddress }
