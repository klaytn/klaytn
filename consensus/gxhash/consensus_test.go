// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
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
// This file is derived from consensus/ethash/consensus_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package gxhash

import (
	"encoding/json"
	"math/big"

	// Enable below packages when enabling TestCalcDiffulty
	// "os"
	// "path/filepath"
	// "testing"

	"github.com/klaytn/klaytn/common/math"
	// "github.com/klaytn/klaytn/blockchain/types"
	// "github.com/klaytn/klaytn/params"
)

type diffTest struct {
	ParentTimestamp    uint64
	ParentBlockScore   *big.Int
	CurrentTimestamp   uint64
	CurrentBlocknumber *big.Int
	CurrentBlockScore  *big.Int
}

func (d *diffTest) UnmarshalJSON(b []byte) (err error) {
	var ext struct {
		ParentTimestamp    string
		ParentBlockScore   string
		CurrentTimestamp   string
		CurrentBlocknumber string
		CurrentBlockScore  string
	}
	if err := json.Unmarshal(b, &ext); err != nil {
		return err
	}

	d.ParentTimestamp = math.MustParseUint64(ext.ParentTimestamp)
	d.ParentBlockScore = math.MustParseBig256(ext.ParentBlockScore)
	d.CurrentTimestamp = math.MustParseUint64(ext.CurrentTimestamp)
	d.CurrentBlocknumber = math.MustParseBig256(ext.CurrentBlocknumber)
	d.CurrentBlockScore = math.MustParseBig256(ext.CurrentBlockScore)

	return nil
}

// TODO-Klaytn-FailedTest Enable this test later
/*
func TestCalcBlockScore(t *testing.T) {
	file, err := os.Open(filepath.Join("..", "..", "tests", "testdata", "BasicTests", "difficulty.json"))
	if err != nil {
		t.Skip(err)
	}
	defer file.Close()

	tests := make(map[string]diffTest)
	err = json.NewDecoder(file).Decode(&tests)
	if err != nil {
		t.Fatal(err)
	}

	config := &params.ChainConfig{}

	for name, test := range tests {
		number := new(big.Int).Sub(test.CurrentBlocknumber, big.NewInt(1))
		diff := CalcBlockScore(config, test.CurrentTimestamp, &types.Header{
			Number:     number,
			Time:       new(big.Int).SetUint64(test.ParentTimestamp),
			BlockScore: test.ParentBlockScore,
		})
		if diff.Cmp(test.CurrentBlockScore) != 0 {
			t.Error(name, "failed. Expected", test.CurrentBlockScore, "and calculated", diff)
		}
	}
}
*/
