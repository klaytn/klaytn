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

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/common"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto/bls"
)

type BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
}

type BlsPublicKeyInfos map[common.Address]BlsPublicKeyInfo

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
	ret, err := caller.GetAllInfo(opts)
	if err != nil {
		return nil, err
	}

	if len(ret.AddrList) != len(ret.PubkeyList) {
		return nil, ErrKip113BadResult
	}

	infos := make(BlsPublicKeyInfos)
	for i := 0; i < len(ret.AddrList); i++ {
		addr := ret.AddrList[i]
		info := BlsPublicKeyInfo{
			PublicKey: ret.PubkeyList[i].PublicKey,
			Pop:       ret.PubkeyList[i].Pop,
		}

		if err := info.Verify(); err == nil {
			infos[addr] = info
		}
	}

	return infos, err
}
