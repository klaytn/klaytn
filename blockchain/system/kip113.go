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
	"encoding/hex"
	"encoding/json"
	"math/big"
	"sort"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/klaytn/klaytn/common"
	contracts "github.com/klaytn/klaytn/contracts/system_contracts"
	"github.com/klaytn/klaytn/crypto/bls"
	"github.com/klaytn/klaytn/params"
)

type BlsPublicKeyInfo struct {
	PublicKey []byte
	Pop       []byte
	VerifyErr error // Nil if valid. Must check before use.
}

func newBlsPublicKeyInfo(publicKey []byte, pop []byte) BlsPublicKeyInfo {
	return BlsPublicKeyInfo{
		PublicKey: publicKey,
		Pop:       pop,
		VerifyErr: verifyBlsPublicKeyInfo(publicKey, pop),
	}
}

func verifyBlsPublicKeyInfo(publicKey []byte, pop []byte) error {
	pk, err := bls.PublicKeyFromBytes(publicKey)
	if err != nil {
		return err
	}

	sig, err := bls.SignatureFromBytes(pop)
	if err != nil {
		return err
	}

	if !bls.PopVerify(pk, sig) {
		return ErrKip113BadPop
	}

	return nil
}

type BlsPublicKeyInfos map[common.Address]BlsPublicKeyInfo

func (infos BlsPublicKeyInfos) String() string {
	obj := make(map[string]string)
	for addr, info := range infos {
		obj[addr.Hex()] = hex.EncodeToString(info.PublicKey)
	}
	j, _ := json.Marshal(obj)
	return string(j)
}

type AllocKip113Init struct {
	Infos BlsPublicKeyInfos
	Owner common.Address
}

func AllocKip113Proxy(init AllocKip113Init) map[common.Hash]common.Hash {
	if init.Infos == nil {
		return nil
	}
	storage := make(map[common.Hash]common.Hash)

	// Overall storage layout for KIP113 contract:
	//
	// | Name          | Type                                                | Slot | Offset | Bytes |
	// |---------------|-----------------------------------------------------|------|--------|-------|
	// | _initialized  | uint8                                               | 0    | 0      | 1     |
	// | _initializing | bool                                                | 0    | 1      | 1     |
	// | __gap         | uint256[50]                                         | 1    | 0      | 1600  |
	// | __gap         | uint256[50]                                         | 51   | 0      | 1600  |
	// | __gap         | uint256[50]                                         | 101  | 0      | 1600  |
	// | _owner        | address                                             | 151  | 0      | 20    |
	// | __gap         | uint256[49]                                         | 152  | 0      | 1568  |
	// | allNodeIds    | address[]                                           | 201  | 0      | 32    |
	// | record        | mapping(address => struct IKIP113.BlsPublicKeyInfo) | 202  | 0      | 32    |
	//
	// We need to consider the following:
	storage[lpad32(0)] = lpad32([]byte{0, 1}) // false, 1
	storage[lpad32(151)] = lpad32(init.Owner)

	addrs := make([]common.Address, 0)
	for addr := range init.Infos {
		addrs = append(addrs, addr)
	}
	sort.Slice(addrs, func(i, j int) bool {
		return strings.Compare(addrs[i].Hex(), addrs[j].Hex()) > 0
	})

	// slot[201]: address[] allNodeIds;
	// - addrs.length @ 0
	// - addrs[i] @ Hash(201) + i
	storage[lpad32(201)] = lpad32(len(addrs))
	for i, addr := range addrs {
		addrSlot := calcArraySlot(201, 1, i, 0)
		storage[addrSlot] = lpad32(addr)
	}

	// slot[202]: mapping(address => BlsPublicKeyInfo) record;
	// - infos[x].publicKey.length @ Hash(x, 202)      = el
	// - infos[x].publicKey        @ Hash(el) + 0..1
	// - infos[x].pop.length       @ Hash(x, 202) + 1  = el
	// - infos[x].pop              @ Hash(el) + 0..2
	for addr, info := range init.Infos {
		// The below slot calculation assumes 48-byte and 96-byte Solidity `bytes` values.
		if len(info.PublicKey) != 48 || len(info.Pop) != 96 {
			logger.Crit("Invalid AllocKip113Init")
		}

		pubSlot := calcMappingSlot(202, addr, 0)
		popSlot := calcMappingSlot(202, addr, 1)

		for k, v := range allocDynamicData(pubSlot, info.PublicKey) {
			storage[k] = v
		}
		for k, v := range allocDynamicData(popSlot, info.Pop) {
			storage[k] = v
		}
	}

	return storage
}

func AllocKip113Logic() map[common.Hash]common.Hash {
	storage := make(map[common.Hash]common.Hash)

	// We only need to case about _initialized, which is max(uint8).
	storage[lpad32(0)] = lpad32([]byte{0xff})

	return storage
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
		infos[addr] = newBlsPublicKeyInfo(
			ret.PubkeyList[i].PublicKey,
			ret.PubkeyList[i].Pop,
		)
	}

	return infos, err
}

func ReadKip113FromConfig(config *params.ChainConfig) (common.Address, error) {
	if (config.RandaoRegistry == nil) || (config.RandaoRegistry.Records == nil) {
		logger.Error("RandaoRegistry not correctly set in ChainConfig")
		return common.Address{}, ErrKip113NotConfigured
	}
	kip113Addr, ok := config.RandaoRegistry.Records[Kip113Name]
	if !ok {
		logger.Error("KIP113 address not set in ChainConfig")
		return common.Address{}, ErrKip113NotConfigured
	}
	return kip113Addr, nil
}
