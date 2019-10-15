// Modifications Copyright 2018 The klaytn Authors
// Copyright 2016 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from cmd/geth/genesis_test.go (2018/06/04).
// Modified and improved for the klaytn development.

package nodecmd

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

var customGenesisTests = []struct {
	genesis string
	query   []string
	result  []string
}{
	// Plain genesis file without anything extra
	{
		genesis: `{
			"alloc"      : {},
			"blockScore" : "0x20000",
			"extraData"  : "0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0",
			"gasLimit"   : "0x2fefd8",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00"
		}`,
		query:  []string{"klay.getBlock(0).parentHash"},
		result: []string{"0x0000000000000000000000000000000000000000000000000000000000000000"},
	},
	// Genesis file with an empty chain configuration (ensure missing fields work)
	{
		genesis: `{
			"alloc"      : {},
			"blockScore" : "0x20000",
			"extraData"  : "0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0",
			"gasLimit"   : "0x2fefd8",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {}
		}`,
		query:  []string{"klay.getBlock(0).parentHash"},
		result: []string{"0x0000000000000000000000000000000000000000000000000000000000000000"},
	},
	// Genesis file with specific chain configurations
	{
		genesis: `{
			"alloc"      : {},
			"blockScore" : "0x20000",
			"extraData"  : "0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0",
			"gasLimit"   : "0x2fefd8",
			"parentHash" : "0x0000000000000000000000000000000000000000000000000000000000000000",
			"timestamp"  : "0x00",
			"config"     : {
				"homesteadBlock" : 314,
				"daoForkBlock"   : 141,
				"daoForkSupport" : true
			}
		}`,
		query:  []string{"klay.getBlock(0).parentHash"},
		result: []string{"0x0000000000000000000000000000000000000000000000000000000000000000"},
	},
	{
		genesis: `{
		    "config": {
		        "chainId": 1000,
		        "istanbul": {
		            "epoch": 30,
		            "policy": 0,
		            "sub": 22
		        },
		        "unitPrice": 0,
		        "deriveShaImpl": 2,
		        "governance": null
		    },
		    "timestamp": "0x5da3fdfd",
		    "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000f89af8549413c9a496fb7b84ea6bf39f3602658c41f0dd7a51947cf8c6a6a6a4fbb9b846527b5d762b278adfe7369438b99937e557a21d9ab9b6d9cfc07498b50660d09442b1d1fccee0e7b51a5091c532eb2b92a1e49296b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0",
		    "governanceData": null,
		    "blockScore": "0x1",
		    "alloc": {
		        "04d67289f30351fc46c161eb3c588d07c2995367": {
		            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
		        },
		        "13c9a496fb7b84ea6bf39f3602658c41f0dd7a51": {
		            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
		        },
		        "2913ecaec7da798466611271a27b53836f20b108": {
		            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
		        },
		        "30208f32c70e8b53a67ea171c8720cbfe32888ff": {
		            "balance": "0x446c3b15f9926687d2c40534fdb564000000000000"
		        }
		    },
		    "number": "0x0",
		    "gasUsed": "0x0",
		    "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"
		}`,
		query:  []string{"klay.getBlock(0).parentHash", "klay.getBlock(0).blockscore", "klay.getBlock(0).extraData", "governance.chainConfig.chainId", "governance.chainConfig.deriveShaImpl", "governance.chainConfig.unitPrice", "governance.chainConfig.governance.governanceMode", "governance.chainConfig.governance.governingNode", "governance.chainConfig.governance.reward.deferredTxFee", "governance.chainConfig.governance.reward.minimumStake", "governance.chainConfig.governance.reward.mintingAmount", "governance.chainConfig.governance.reward.proposerUpdateInterval", "governance.chainConfig.governance.reward.ratio", "governance.chainConfig.governance.reward.stakingUpdateInterval", "governance.chainConfig.governance.reward.useGiniCoeff", "governance.chainConfig.istanbul.epoch", "governance.chainConfig.istanbul.policy", "governance.chainConfig.istanbul.sub"},
		result: []string{"0x0000000000000000000000000000000000000000000000000000000000000000", "0x1", "0x0000000000000000000000000000000000000000000000000000000000000000f89af8549413c9a496fb7b84ea6bf39f3602658c41f0dd7a51947cf8c6a6a6a4fbb9b846527b5d762b278adfe7369438b99937e557a21d9ab9b6d9cfc07498b50660d09442b1d1fccee0e7b51a5091c532eb2b92a1e49296b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0", "1000", "2", "0", "none", "0x0000000000000000000000000000000000000000", "false", "2000000", "0", "3600", "100/0/0", "86400", "false", "30", "0", "22"},
	},
	{
		genesis: `{"config":{"chainId":2019,"istanbul":{"epoch":30,"policy":2,"sub":13},"unitPrice":25000000000,"deriveShaImpl":2,"governance":{"governingNode":"0xdddfb991127b43e209c2f8ed08b8b3d0b5843d36","governanceMode":"single","reward":{"mintingAmount":9600000000000000000,"ratio":"34/54/12","useGiniCoeff":false,"deferredTxFee":true,"stakingUpdateInterval":60,"proposerUpdateInterval":30,"minimumStake":5000000}}},"timestamp":"0x5ce33d6e","extraData":"0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0","governanceData":null,"blockScore":"0x1","alloc":{"195ba9cc787b00796a7ae6356e5b656d43603537":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"777fd033b5e3bcaad6006bc9f481ffed6b83cf5a":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"d473284239f704adccd24647c7ca132992a28973":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"dddfb991127b43e209c2f8ed08b8b3d0b5843d36":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"},"f4316f69d9522667c0674afcd8638288489f0333":{"balance":"0x446c3b15f9926687d2c40534fdb564000000000000"}},"number":"0x0","gasUsed":"0x0","parentHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`,
		query:   []string{"klay.getBlock(0).parentHash", "klay.getBlock(0).blockscore", "klay.getBlock(0).extraData", "governance.chainConfig.chainId", "governance.chainConfig.deriveShaImpl", "governance.chainConfig.unitPrice", "governance.chainConfig.governance.governanceMode", "governance.chainConfig.governance.governingNode", "governance.chainConfig.governance.reward.deferredTxFee", "governance.chainConfig.governance.reward.minimumStake", "governance.chainConfig.governance.reward.mintingAmount", "governance.chainConfig.governance.reward.proposerUpdateInterval", "governance.chainConfig.governance.reward.ratio", "governance.chainConfig.governance.reward.stakingUpdateInterval", "governance.chainConfig.governance.reward.useGiniCoeff", "governance.chainConfig.istanbul.epoch", "governance.chainConfig.istanbul.policy", "governance.chainConfig.istanbul.sub"},
		result:  []string{"0x0000000000000000000000000000000000000000000000000000000000000000", "0x1", "0x0000000000000000000000000000000000000000000000000000000000000000f89af85494dddfb991127b43e209c2f8ed08b8b3d0b5843d3694195ba9cc787b00796a7ae6356e5b656d4360353794777fd033b5e3bcaad6006bc9f481ffed6b83cf5a94d473284239f704adccd24647c7ca132992a28973b8410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0", "2019", "2", "25000000000", "single", "0xdddfb991127b43e209c2f8ed08b8b3d0b5843d36", "true", "5000000", "9600000000000000000", "30", "34/54/12", "60", "false", "30", "2", "13"},
	},
}

// Tests that initializing Klay with a custom genesis block and chain definitions
// work properly.
func TestCustomGenesis(t *testing.T) {
	for i, tt := range customGenesisTests {
		// Create a temporary data directory to use and inspect later
		datadir := tmpdir(t)
		defer os.RemoveAll(datadir)

		// Initialize the data directory with the custom genesis block
		json := filepath.Join(datadir, "genesis.json")
		if err := ioutil.WriteFile(json, []byte(tt.genesis), 0600); err != nil {
			t.Fatalf("test %d: failed to write genesis file: %v", i, err)
		}
		runKlay(t, "klay-test", "--datadir", datadir, "--verbosity", "0", "init", json).WaitExit()

		// Query the custom genesis block
		if len(tt.query) != len(tt.result) {
			t.Errorf("Test cases are wrong, #query: %v, #result, %v", len(tt.query), len(tt.result))
		}
		for idx, query := range tt.query {
			klay := runKlay(t,
				"klay-test", "--datadir", datadir, "--maxconnections", "0", "--port", "0",
				"--nodiscover", "--nat", "none", "--ipcdisable",
				"--exec", query, "--verbosity", "0", "console")
			klay.ExpectRegexp(tt.result[idx])
			klay.ExpectExit()
		}
	}
}
