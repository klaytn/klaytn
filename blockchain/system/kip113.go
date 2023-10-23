// Copyright 2023 The klaytn Authors
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

package system

import (
	"math/big"
	"sort"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/common"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto/bls"
)

type BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
}

type (
	BlsPublicKeyInfos map[common.Address]BlsPublicKeyInfo
	AllocKip113Init   = BlsPublicKeyInfos
)

func AllocKip113(init AllocKip113Init) map[common.Hash]common.Hash {
	if init == nil {
		return nil
	}
	storage := make(map[common.Hash]common.Hash)

	// slot[0]: mapping(address => BlsPublicKeyInfo) infos;
	// - infos[x].publicKey.length @ Hash(x, 0)      = el
	// - infos[x].publicKey        @ Hash(el) + 0..1
	// - infos[x].pop.length       @ Hash(x, 0) + 1  = el
	// - infos[x].pop              @ Hash(el) + 0..2
	for addr, info := range init {
		// The below slot calculation assumes 48-byte and 96-byte Solidity `bytes` values.
		if len(info.PublicKey) != 48 || len(info.Pop) != 96 {
			logger.Crit("Invalid AllocKip113Init")
		}

		pubSlot := calcMappingSlot(0, addr, 0)
		popSlot := calcMappingSlot(0, addr, 1)

		for k, v := range allocDynamicData(pubSlot, info.PublicKey) {
			storage[k] = v
		}
		for k, v := range allocDynamicData(popSlot, info.Pop) {
			storage[k] = v
		}
	}

	addrs := make([]common.Address, 0)
	for addr := range init {
		addrs = append(addrs, addr)
	}
	sort.Slice(addrs, func(i, j int) bool {
		return strings.Compare(addrs[i].Hex(), addrs[j].Hex()) > 0
	})

	// slot[1]: address[] addrs;
	// - addrs.length @ 1
	// - addrs[i] @ Hash(1) + i
	storage[lpad32(1)] = lpad32(len(addrs))
	for i, addr := range addrs {
		addrSlot := calcArraySlot(1, 1, i, 0)
		storage[addrSlot] = lpad32(addr)
	}

	return storage
}

func (info BlsPublicKeyInfo) Verify() error {
	pk, err := bls.PublicKeyFromBytes(info.PublicKey)
	if err != nil {
		return err
	}

	sig, err := bls.SignatureFromBytes(info.Pop)
	if err != nil {
		return err
	}

	if !bls.PopVerify(pk, sig) {
		return ErrKip113BadPop
	}

	return nil
}

func ReadKip113All(backend bind.ContractCaller, contractAddr common.Address, num *big.Int) (BlsPublicKeyInfos, error) {
	caller, err := contracts.NewIKIP113Caller(contractAddr, backend)
	if err != nil {
		return nil, err
	}

	opts := &bind.CallOpts{BlockNumber: num}
	ret, err := caller.GetAllBlsInfo(opts)
	if err != nil {
		return nil, err
	}

	if len(ret.NodeIdList) != len(ret.PubkeyList) {
		return nil, ErrKip113BadResult
	}

	infos := make(BlsPublicKeyInfos)
	for i := 0; i < len(ret.NodeIdList); i++ {
		addr := ret.NodeIdList[i]
		info := BlsPublicKeyInfo{
			PublicKey: ret.PubkeyList[i].PublicKey,
			Pop:       ret.PubkeyList[i].Pop,
		}

		if info.Verify() == nil {
			infos[addr] = info
		}
	}

	return infos, err
}
