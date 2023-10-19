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

type AllocKIP113Init struct {
	// Owner common.Address
	Infos map[common.Address]BlsPublicKeyInfo
	Addrs []common.Address
}

type BlsPublicKeyInfos map[common.Address]BlsPublicKeyInfo

func AllocKIP113(init *AllocKIP113Init) map[common.Hash]common.Hash {
	if init == nil {
		return nil
	}
	storage := make(map[common.Hash]common.Hash)

	// TODO-klaytn: add storage slot calculation for owner address

	// slot[0]: mapping(address => BlsPublicKeyInfo) infos;
	// - infos[x].publicKey.length @ Hash(x, 0)
	//   - infos[x].publicKey[:32] @ Hash(Hash(x, 0)) + 0
	//   - infos[x].publicKey[32:] @ Hash(Hash(x, 0)) + 1
	// - infos[x].pop.length @ Hash(x, 1) + 1
	//   - infos[x].pop[:32] @ Hash(Hash(x, 1) + 1) + 0
	//   - infos[x].pop[32:64] @ Hash(Hash(x, 1) + 1) + 1
	//   - infos[x].pop[64:] @ Hash(Hash(x, 1) + 1) + 2
	for addr, info := range init.Infos {
		if len(info.PublicKey) != 48 || len(info.Pop) != 96 {
			logger.Crit("Invalid BLS key", "PublicKey", info.PublicKey, "Pop", info.Pop)
			continue
		}
		pupKeySlot := calcMappingSlot(0, addr)
		popKeySlot := addIntToHash(pupKeySlot, 1)

		// set publicKey.length and pop.length
		storage[pupKeySlot] = lpad32(len(info.PublicKey)*2 + 1)
		storage[popKeySlot] = lpad32(len(info.Pop)*2 + 1)

		// publicKey[:32]
		storage[calcStructSlot(pupKeySlot, 0)] = common.BytesToHash(info.PublicKey[:32])
		// publicKey[32:]
		// Note: solidity padded from the right for the byte type.
		var h common.Hash
		copy(h[:len(info.PublicKey[32:])], info.PublicKey[32:])
		storage[calcStructSlot(pupKeySlot, 1)] = h

		// pop[:32]
		storage[calcStructSlot(popKeySlot, 0)] = common.BytesToHash(info.Pop[:32])
		// pop[32:64]
		storage[calcStructSlot(popKeySlot, 1)] = common.BytesToHash(info.Pop[32:64])
		// pop[64:]
		storage[calcStructSlot(popKeySlot, 2)] = common.BytesToHash(info.Pop[64:])
	}

	// slot[1]: address[] addrs;
	// - addrs.length @ 1
	// - addrs[i] @ Hash(1) + i
	storage[lpad32(1)] = lpad32(len(init.Addrs))
	for i, addr := range init.Addrs {
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
