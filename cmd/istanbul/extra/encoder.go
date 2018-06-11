// Copyright 2017 AMIS Technologies
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

package extra

import (
	"bytes"
	"ground-x/go-gxplatform/common"
	"ground-x/go-gxplatform/common/hexutil"
	"ground-x/go-gxplatform/rlp"
	"ground-x/go-gxplatform/core/types"
)

func Encode(vanity string, validators []common.Address) (string, error) {
	newVanity, err := hexutil.Decode(vanity)
	if err != nil {
		return "", err
	}

	if len(newVanity) < types.IstanbulExtraVanity {
		newVanity = append(newVanity, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(newVanity))...)
	}
	newVanity = newVanity[:types.IstanbulExtraVanity]

	ist := &types.IstanbulExtra{
		Validators:    validators,
		Seal:          make([]byte, types.IstanbulExtraSeal),
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return "", err
	}

	return "0x" + common.Bytes2Hex(append(newVanity, payload...)), nil
}