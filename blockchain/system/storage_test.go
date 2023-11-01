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
	"testing"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestStorageCalc(t *testing.T) {
	assert.Equal(t, "0x839613f731613c3a2f728362760f939c8004b5d9066154aab51d6dadf74733f3",
		calcMappingSlot(0, common.HexToAddress("0xaaaa"), 0).Hex())
	assert.Equal(t, "0x839613f731613c3a2f728362760f939c8004b5d9066154aab51d6dadf74733f4",
		calcMappingSlot(0, common.HexToAddress("0xaaaa"), 1).Hex())

	assert.Equal(t, "0xb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf6",
		calcArraySlot(1, 1, 0, 0).Hex())
	assert.Equal(t, "0xb10e2d527612073b26eecdfd717e6a320cf44b4afac2b0732d9fcbe2b7fa0cf7",
		calcArraySlot(1, 1, 1, 0).Hex())

	assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000001234",
		lpad32(0x1234).Hex())
}

func TestStorageDynamicData(t *testing.T) {
	testcases := []struct {
		baseSlot interface{}
		data     []byte
		expected map[string]string
	}{
		{ // short data
			baseSlot: 1,
			data:     []byte("AcmeContract"),
			expected: map[string]string{
				"0x0000000000000000000000000000000000000000000000000000000000000001": "0x41636d65436f6e74726163740000000000000000000000000000000000000018",
			},
		},
		{ // long 48-byte data (e.g. BLS public key)
			baseSlot: common.HexToHash("0x839613f731613c3a2f728362760f939c8004b5d9066154aab51d6dadf74733f3"),
			data:     hexutil.MustDecode("0x91edb62902264ada822e48f0beedf582fa8c5d1518c0b536b1347983f2fe9c80f283c1e5dfc7a37c3437ceaaa9a1f867"),
			expected: map[string]string{
				"0x839613f731613c3a2f728362760f939c8004b5d9066154aab51d6dadf74733f3": "0x0000000000000000000000000000000000000000000000000000000000000061",
				"0x89005672ae92e5773a2b83cb802d9bc9be4cf262b7ddcca0635550f0ea96e12d": "0x91edb62902264ada822e48f0beedf582fa8c5d1518c0b536b1347983f2fe9c80",
				"0x89005672ae92e5773a2b83cb802d9bc9be4cf262b7ddcca0635550f0ea96e12e": "0xf283c1e5dfc7a37c3437ceaaa9a1f86700000000000000000000000000000000",
			},
		},
		{ // long 95-byte data (e.g. BLS PoP)
			baseSlot: common.HexToHash("0x839613f731613c3a2f728362760f939c8004b5d9066154aab51d6dadf74733f4"),
			data:     hexutil.MustDecode("0xb72d443815683d633b5933366ac8aa6e2891aaff8ab95e69a7531aa2afe10f0740368128992a2a54cb24f0c93dad5c220692915ddb5013fb475726baf711590f7ea1eafaf396264cc2940c744ee6914f3734a4d404dffa91d15e3d138ee81464"),
			expected: map[string]string{
				"0x839613f731613c3a2f728362760f939c8004b5d9066154aab51d6dadf74733f4": "0x00000000000000000000000000000000000000000000000000000000000000c1",
				"0xccf7c844150261de21551c00e90c49fd101dcf60b58bfeb3f7c1ce662560a383": "0xb72d443815683d633b5933366ac8aa6e2891aaff8ab95e69a7531aa2afe10f07",
				"0xccf7c844150261de21551c00e90c49fd101dcf60b58bfeb3f7c1ce662560a384": "0x40368128992a2a54cb24f0c93dad5c220692915ddb5013fb475726baf711590f",
				"0xccf7c844150261de21551c00e90c49fd101dcf60b58bfeb3f7c1ce662560a385": "0x7ea1eafaf396264cc2940c744ee6914f3734a4d404dffa91d15e3d138ee81464",
			},
		},
	}

	for _, tc := range testcases {
		alloc := allocDynamicData(tc.baseSlot, tc.data)

		strAlloc := make(map[string]string)
		for k, v := range alloc {
			strAlloc[k.Hex()] = v.Hex()
		}
		assert.Equal(t, tc.expected, strAlloc)
	}
}
