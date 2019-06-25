// Copyright 2018 The klaytn Authors
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
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common/hexutil"
)

func Decode(extraData string) ([]byte, *types.IstanbulExtra, error) {
	extra, err := hexutil.Decode(extraData)
	if err != nil {
		return nil, nil, err
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(&types.Header{Extra: extra})
	if err != nil {
		return nil, nil, err
	}
	return extra[:types.IstanbulExtraVanity], istanbulExtra, nil
}
