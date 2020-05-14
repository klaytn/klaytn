// Copyright 2020 The klaytn Authors
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

package database

import (
	"bytes"
	"testing"
)

func TestDatabaseManager_StakingInfo(t *testing.T) {
	dbm := dbManagers[0]

	key := uint64(1234)
	value := []byte("{\"BlockNum\":2880,\"CouncilNodeAddrs\":[\"0x159ae5ccda31b77475c64d88d4499c86f77b7ecc\",\"0x181deb121304b0430d99328ff1a9122df9f09d7f\",\"0x324ec8f2681cd73642cc55057970540a1f4393e0\",\"0x11191029025d3fcd21001746f949b25c6e8435cc\"],\"CouncilStakingAddrs\":[\"0x70e051c46ea76b9af9977407bb32192319907f9e\",\"0xe4a0c3821a2711758306ed57c2f4900aa9ddbb3d\",\"0xf3ba3a33b3bf7cf2085890315b41cc788770feb3\",\"0x9285a85777d0ae7e12bee3ffd7842908b2295f45\"],\"CouncilRewardAddrs\":[\"0xd155d4277c99fa837c54a37a40a383f71a3d082a\",\"0x2b8cc0ca62537fa5e49dce197acc8a15d3c5d4a8\",\"0x7d892f470ecde693f52588dd0cfe46c3d26b6219\",\"0xa0f7354a0cef878246820b6caa19d2bdef74a0cc\"],\"KIRAddr\":\"0x673003e5f9a852d3dc85b83d16ef62d45497fb96\",\"PoCAddr\":\"0x576dc0c2afeb1661da3cf53a60e76dd4e32c7ab1\",\"UseGini\":false,\"Gini\":-1,\"CouncilStakingAmounts\":[5000000,5000000,5000000,5000000]}")
	err := dbm.WriteStakingInfo(key, value)
	if err != nil {
		t.Fatal(err)
	}

	rValue, err := dbm.ReadStakingInfo(key)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(value, rValue) {
		t.Fatal(err)
	}
}
