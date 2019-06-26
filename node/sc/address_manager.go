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

package sc

import (
	"errors"
	"github.com/klaytn/klaytn/common"
)

var (
	ErrAlreadyExistentBridgePair = errors.New("bridge already exists")
)

// AddressManager manages mapping addresses for bridge,contract,user
// to exchange value between parent and child chain
type AddressManager struct {
	bridgeContracts map[common.Address]common.Address
	tokenContracts  map[common.Address]common.Address
}

func NewAddressManager() (*AddressManager, error) {
	return &AddressManager{
		bridgeContracts: make(map[common.Address]common.Address),
		tokenContracts:  make(map[common.Address]common.Address),
	}, nil
}

func (am *AddressManager) AddBridge(bridge1 common.Address, bridge2 common.Address) error {
	_, ok1 := am.bridgeContracts[bridge1]
	_, ok2 := am.bridgeContracts[bridge2]

	if ok1 || ok2 {
		return ErrAlreadyExistentBridgePair
	}

	am.bridgeContracts[bridge1] = bridge2
	am.bridgeContracts[bridge2] = bridge1

	logger.Info("succeeded to AddBridge", "bridge1", bridge1.String(), "bridge2", bridge2.String())
	return nil
}

func (am *AddressManager) DeleteBridge(bridge1 common.Address) (common.Address, common.Address, error) {
	bridge2, ok1 := am.bridgeContracts[bridge1]
	if !ok1 {
		return common.Address{}, common.Address{}, errors.New("bridge does not exist")
	}

	delete(am.bridgeContracts, bridge1)
	delete(am.bridgeContracts, bridge2)

	logger.Info("succeeded to DeleteBridge", "bridge1", bridge1.String(), "bridge2", bridge2.String())
	return bridge1, bridge2, nil
}

func (am *AddressManager) AddToken(token1 common.Address, token2 common.Address) error {
	_, ok1 := am.tokenContracts[token1]
	_, ok2 := am.tokenContracts[token2]

	if ok1 || ok2 {
		return errors.New("token already exists")
	}

	am.tokenContracts[token2] = token1
	am.tokenContracts[token1] = token2
	return nil
}

func (am *AddressManager) DeleteToken(token1 common.Address) (common.Address, common.Address, error) {
	token2, ok1 := am.tokenContracts[token1]
	if !ok1 {
		return common.Address{}, common.Address{}, errors.New("token does not exist")
	}

	delete(am.tokenContracts, token1)
	delete(am.tokenContracts, token2)

	return token1, token2, nil
}

func (am *AddressManager) GetCounterPartBridge(addr common.Address) common.Address {
	bridge, ok := am.bridgeContracts[addr]
	if !ok {
		return common.Address{}
	}
	return bridge
}

func (am *AddressManager) GetCounterPartToken(addr common.Address) common.Address {
	token, ok := am.tokenContracts[addr]
	if !ok {
		return common.Address{}
	}
	return token
}
