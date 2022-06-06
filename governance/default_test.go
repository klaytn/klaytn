// Copyright 2019 The klaytn Authors
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

package governance

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	gotest_assert "gotest.tools/assert"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/consensus/istanbul/validator"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
)

type voteValue struct {
	k string
	v interface{}
	e bool
}

var tstData = []voteValue{
	{k: "istanbul.epoch", v: uint64(30000), e: true},
	{k: "istanbul.epoch", v: "bad", e: false},
	{k: "istanbul.epoch", v: float64(30000.00), e: true},
	{k: "istanbul.Epoch", v: float64(30000.10), e: false},
	{k: "istanbul.epoch", v: true, e: false},
	{k: "istanbul.committeesize", v: uint64(7), e: true},
	{k: "istanbul.committeesize", v: float64(7.0), e: true},
	{k: "istanbul.committeesize", v: float64(7.1), e: false},
	{k: "istanbul.committeesize", v: "7", e: false},
	{k: "istanbul.committeesize", v: true, e: false},
	{k: "istanbul.committeesize", v: float64(-7), e: false},
	{k: "istanbul.committeesize", v: uint64(0), e: false},
	{k: "istanbul.policy", v: "roundrobin", e: false},
	{k: "istanbul.policy", v: "RoundRobin", e: false},
	{k: "istanbul.policy", v: "sticky", e: false},
	{k: "istanbul.policy", v: "weightedrandom", e: false},
	{k: "istanbul.policy", v: "WeightedRandom", e: false},
	{k: "istanbul.policy", v: uint64(0), e: false},
	{k: "istanbul.policy", v: uint64(1), e: false},
	{k: "istanbul.policy", v: uint64(2), e: false},
	{k: "istanbul.policy", v: float64(1.2), e: false},
	{k: "istanbul.policy", v: float64(1.0), e: false},
	{k: "istanbul.policy", v: true, e: false},
	{k: "governance.governancemode", v: "none", e: true},
	{k: "governance.governancemode", v: "single", e: true},
	{k: "governance.governancemode", v: "ballot", e: true},
	{k: "governance.governancemode", v: 0, e: false},
	{k: "governance.governancemode", v: 1, e: false},
	{k: "governance.governancemode", v: 2, e: false},
	{k: "governance.governancemode", v: "unexpected", e: false},
	{k: "governance.governingnode", v: "0x00000000000000000000", e: false},
	{k: "governance.governingnode", v: "0x0000000000000000000000000000000000000000", e: true},
	{k: "governance.governingnode", v: "0x000000000000000000000000000abcd000000000", e: true},
	{k: "governance.governingnode", v: "000000000000000000000000000abcd000000000", e: true},
	{k: "governance.governingnode", v: common.HexToAddress("000000000000000000000000000abcd000000000"), e: true},
	{k: "governance.governingnode", v: "0x000000000000000000000000000xxxx000000000", e: false},
	{k: "governance.governingnode", v: "address", e: false},
	{k: "governance.governingnode", v: 0, e: false},
	{k: "governance.governingnode", v: true, e: false},
	{k: "governance.unitprice", v: float64(0.0), e: true},
	{k: "governance.unitprice", v: float64(0.1), e: false},
	{k: "governance.unitprice", v: uint64(25000000000), e: true},
	{k: "governance.unitprice", v: float64(-10), e: false},
	{k: "governance.unitprice", v: "25000000000", e: false},
	{k: "governance.unitprice", v: true, e: false},
	{k: "reward.useginicoeff", v: false, e: true},
	{k: "reward.useginicoeff", v: true, e: true},
	{k: "reward.useginicoeff", v: "true", e: false},
	{k: "reward.useginicoeff", v: 0, e: false},
	{k: "reward.useginicoeff", v: 1, e: false},
	{k: "reward.mintingamount", v: "9600000000000000000", e: true},
	{k: "reward.mintingamount", v: "0", e: true},
	{k: "reward.mintingamount", v: 96000, e: false},
	{k: "reward.mintingamount", v: "many", e: false},
	{k: "reward.mintingamount", v: true, e: false},
	{k: "reward.ratio", v: "30/40/30", e: true},
	{k: "reward.ratio", v: "10/10/80", e: true},
	{k: "reward.ratio", v: "30/70", e: false},
	{k: "reward.ratio", v: "30/40/31", e: false},
	{k: "reward.ratio", v: "30/40/29", e: false},
	{k: "reward.ratio", v: 30 / 40 / 30, e: false},
	{k: "reward.ratio", v: "0/0/100", e: true},
	{k: "reward.ratio", v: "0/100/0", e: true},
	{k: "reward.ratio", v: "100/0/0", e: true},
	{k: "reward.ratio", v: "0/0/0", e: false},
	{k: "reward.ratio", v: "30.5/40/29.5", e: false},
	{k: "reward.ratio", v: "30.5/40/30.5", e: false},
	{k: "reward.deferredtxfee", v: true, e: true},
	{k: "reward.deferredtxfee", v: false, e: true},
	{k: "reward.deferredtxfee", v: 0, e: false},
	{k: "reward.deferredtxfee", v: 1, e: false},
	{k: "reward.deferredtxfee", v: "true", e: false},
	{k: "reward.minimumstake", v: "2000000000000000000000000", e: true},
	{k: "reward.minimumstake", v: 200000000000000, e: false},
	{k: "reward.minimumstake", v: "-1", e: false},
	{k: "reward.minimumstake", v: "0", e: true},
	{k: "reward.minimumstake", v: 0, e: false},
	{k: "reward.minimumstake", v: 1.1, e: false},
	{k: "reward.stakingupdateinterval", v: uint64(20), e: false},
	{k: "reward.stakingupdateinterval", v: float64(20.0), e: false},
	{k: "reward.stakingupdateinterval", v: float64(20.2), e: false},
	{k: "reward.stakingupdateinterval", v: "20", e: false},
	{k: "reward.proposerupdateinterval", v: uint64(20), e: false},
	{k: "reward.proposerupdateinterval", v: float64(20.0), e: false},
	{k: "reward.proposerupdateinterval", v: float64(20.2), e: false},
	{k: "reward.proposerupdateinterval", v: "20", e: false},
	{k: "istanbul.timeout", v: uint64(10000), e: true},
	{k: "istanbul.timeout", v: uint64(5000), e: true},
	{k: "istanbul.timeout", v: float64(-1000), e: false},
	{k: "istanbul.timeout", v: true, e: false},
	{k: "istanbul.timeout", v: "10", e: false},
	{k: "istanbul.timeout", v: 5.3, e: false},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c", e: true},
	{k: "governance.addvalidator", v: " 0x639e5ebfc483716fbac9810b230ff6ad487f366c ", e: true},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c,0x828880c5f09cc1cc6a58715e3fe2b4c4cf3c5869", e: true},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c,0x828880c5f09cc1cc6a58715e3fe2b4c4cf3c58", e: false},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c,0x639e5ebfc483716fbac9810b230ff6ad487f366c", e: false},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c, 0x828880c5f09cc1cc6a58715e3fe2b4c4cf3c5869, ", e: false},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c,, 0x828880c5f09cc1cc6a58715e3fe2b4c4cf3c5869, ", e: false},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c, 0x828880c5f09cc1cc6a58715e3fe2b4c4cf3c5869 ", e: true},
}

var goodVotes = []voteValue{
	{k: "istanbul.epoch", v: uint64(20000), e: true},
	{k: "istanbul.committeesize", v: uint64(7), e: true},
	{k: "governance.governancemode", v: "single", e: true},
	{k: "governance.governingnode", v: common.HexToAddress("0x0000000000000000000000000000000000000000"), e: true},
	{k: "governance.unitprice", v: uint64(25000000000), e: true},
	{k: "reward.useginicoeff", v: false, e: true},
	{k: "reward.mintingamount", v: "9600000000000000000", e: true},
	{k: "reward.ratio", v: "10/10/80", e: true},
	{k: "istanbul.timeout", v: uint64(5000), e: true},
	{k: "governance.addvalidator", v: "0x639e5ebfc483716fbac9810b230ff6ad487f366c,0x828880c5f09cc1cc6a58715e3fe2b4c4cf3c5869", e: true},
}

func getTestConfig() *params.ChainConfig {
	config := params.TestChainConfig
	config.Governance = params.GetDefaultGovernanceConfig(params.UseIstanbul)
	config.Istanbul = params.GetDefaultIstanbulConfig()
	return config
}

func getGovernance() *Governance {
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	config := getTestConfig()
	return NewGovernanceInitialize(config, dbm)
}

func TestGetDefaultGovernanceConfig(t *testing.T) {
	tstGovernance := params.GetDefaultGovernanceConfig(params.UseIstanbul)

	want := []interface{}{
		params.DefaultUseGiniCoeff,
		params.DefaultRatio,
		common.HexToAddress(params.DefaultGoverningNode),
		params.DefaultGovernanceMode,
		params.DefaultDefferedTxFee,
	}

	got := []interface{}{
		tstGovernance.Reward.UseGiniCoeff,
		tstGovernance.Reward.Ratio,
		tstGovernance.GoverningNode,
		tstGovernance.GovernanceMode,
		tstGovernance.DeferredTxFee(),
	}

	if !reflect.DeepEqual(want, got) {
		t.Fatalf("Want %v, got %v", want, got)
	}

	if tstGovernance.Reward.MintingAmount.Cmp(params.DefaultMintingAmount) != 0 {
		t.Errorf("Default minting amount is not equal")
	}
}

func TestGovernance_ValidateVote(t *testing.T) {
	gov := getGovernance()

	for _, val := range tstData {
		vote := &GovernanceVote{
			Key:   val.k,
			Value: val.v,
		}
		_, ret := gov.ValidateVote(vote)
		if ret != val.e {
			if _, ok := GovernanceForbiddenKeyMap[val.k]; !ok && !ret {
				t.Errorf("Want %v, got %v for %v and %v", val.e, ret, val.k, val.v)
			}
		}
	}
}

func TestGovernance_AddVote(t *testing.T) {
	gov := getGovernance()

	for _, val := range tstData {
		ret := gov.AddVote(val.k, val.v)
		assert.Equal(t, val.e, ret, fmt.Sprintf("key %v, value %v", val.k, val.v))
	}

}

func TestGovernance_RemoveVote(t *testing.T) {
	gov := getGovernance()

	for _, val := range goodVotes {
		ret := gov.AddVote(val.k, val.v)
		if ret != val.e {
			t.Errorf("Want %v, got %v for %v and %v", val.e, ret, val.k, val.v)
		}
	}

	// Length check. Because []votes has all valid votes, length of voteMap and votes should be equal
	if countUncastedVote(gov.voteMap) != len(goodVotes) {
		t.Errorf("Length of voteMap should be %d, but %d\n", len(goodVotes), gov.voteMap.Size())
	}

	// Remove unvoted vote. Length should still be same
	gov.RemoveVote("istanbul.Epoch", uint64(10000), 0)
	if countUncastedVote(gov.voteMap) != len(goodVotes) {
		t.Errorf("Length of voteMap should be %d, but %d\n", len(goodVotes), gov.voteMap.Size())
	}

	// Remove vote with wrong key. Length should still be same
	gov.RemoveVote("istanbul.EpochEpoch", uint64(10000), 0)
	if countUncastedVote(gov.voteMap) != len(goodVotes) {
		t.Errorf("Length of voteMap should be %d, but %d\n", len(goodVotes), gov.voteMap.Size())
	}

	// Removed a vote. Length should be len(goodVotes) -1
	gov.RemoveVote("istanbul.epoch", uint64(20000), 0)
	if countUncastedVote(gov.voteMap) != (len(goodVotes) - 1) {
		t.Errorf("Length of voteMap should be %d, but %d\n", len(goodVotes)-1, gov.voteMap.Size())
	}
}

func countUncastedVote(data VoteMap) int {
	size := 0

	for _, v := range data.Copy() {
		if v.Casted == false {
			size++
		}
	}
	return size
}

func TestGovernance_ClearVotes(t *testing.T) {
	gov := getGovernance()

	for _, val := range tstData {
		ret := gov.AddVote(val.k, val.v)
		if ret != val.e {
			t.Errorf("Want %v, got %v for %v and %v", val.e, ret, val.k, val.v)
		}
		avt := gov.adjustValueType(val.k, val.v)
		gov.RemoveVote(val.k, avt, 0)
	}
	gov.ClearVotes(0)
	if gov.voteMap.Size() != 0 {
		t.Errorf("Want 0, got %v after clearing votes", gov.voteMap.Size())
	}
}

func TestGovernance_GetEncodedVote(t *testing.T) {
	gov := getGovernance()

	for _, val := range goodVotes {
		_ = gov.AddVote(val.k, val.v)
	}

	l := gov.voteMap.Size()
	for i := 0; i < l; i++ {
		voteData := gov.GetEncodedVote(common.HexToAddress("0x1234567890123456789012345678901234567890"), 1000)
		v := new(GovernanceVote)
		rlp.DecodeBytes(voteData, &v)

		v, err := gov.ParseVoteValue(v)
		assert.Equal(t, nil, err)
		gotest_assert.DeepEqual(t, gov.voteMap.GetValue(v.Key).Value, v.Value)
	}
}

func TestGovernance_ParseVoteValue(t *testing.T) {
	gov := getGovernance()

	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	for _, val := range goodVotes {
		v := &GovernanceVote{
			Key:       val.k,
			Value:     gov.adjustValueType(val.k, val.v),
			Validator: addr,
		}

		b, _ := rlp.EncodeToBytes(v)

		d := new(GovernanceVote)
		rlp.DecodeBytes(b, d)

		d, err := gov.ParseVoteValue(d)
		assert.Equal(t, nil, err)
		gotest_assert.DeepEqual(t, v, d)
	}
}

var testGovernanceMap = map[string]interface{}{
	"governance.governancemode": "none",
	"governance.governingnode":  common.HexToAddress("0x1234567890123456789012345678901234567890"),
	"governance.unitprice":      uint64(25000000000),
	"reward.ratio":              "30/40/30",
	"reward.useginicoeff":       true,
	"reward.deferredtxfee":      true,
	"reward.minimumstake":       2000000,
}

func copyMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{})
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func TestGovernancePersistence(t *testing.T) {
	gov := getGovernance()

	var MAXITEMS = int(10)

	// Write Test
	// WriteGovernance() and WriteGovernanceIdx()
	for i := 1; i < MAXITEMS; i++ {
		blockNum := params.DefaultEpoch * uint64(i)
		tstMap := copyMap(testGovernanceMap)

		// Make every stored governance map has a difference
		tstMap["governance.unitprice"] = tstMap["governance.unitprice"].(uint64) + blockNum
		if gov.CanWriteGovernanceState(blockNum) {
			if err := gov.db.WriteGovernance(tstMap, blockNum); err != nil {
				t.Errorf("Write governance failed: %v", err)
			}
		}
	}

	// Read Test
	// ReadRecentGovernanceIdx() ReadGovernance()
	tstIdx, _ := gov.db.ReadRecentGovernanceIdx(MAXITEMS)
	length := len(tstIdx)
	for i := 1; i < length; i++ {
		num := tstIdx[i]
		compMap, _ := gov.db.ReadGovernance(num)

		expected := testGovernanceMap["governance.unitprice"].(uint64) + uint64(i)*params.DefaultEpoch
		if uint64(compMap["governance.unitprice"].(float64)) != expected {
			t.Errorf("Retrieved %v, Expected %v", compMap["governance.unitprice"], expected)
		}
	}

	tstIdx2, _ := gov.db.ReadRecentGovernanceIdx(0)

	if len(tstIdx2) != MAXITEMS {
		t.Errorf("ReadRecentGovernanceIdx with 0 failure. want %v have %v", MAXITEMS, len(tstIdx2))
	}

	// ReadGovernanceAtNumber

	for i := 0; i < MAXITEMS; i++ {
		num := params.DefaultEpoch*uint64(i) + 123
		idx, _, err := gov.db.ReadGovernanceAtNumber(num, params.DefaultEpoch)
		if err != nil {
			t.Errorf("Failed to get the governance information for block %d", num)
		}
		tValue := num - (num % params.DefaultEpoch)
		if tValue >= params.DefaultEpoch {
			tValue -= params.DefaultEpoch
		}
		if idx != tValue {
			t.Errorf("Wrong block number, Want %v, have %v", tValue, idx)
		}
	}
}

type governanceData struct {
	n uint64
	e uint64
}

var tstGovernanceInfo = []governanceData{
	{n: 1, e: 25000000000},
	{n: 1209600, e: 25001209600},
	{n: 2419200, e: 25002419200},
	{n: 3628800, e: 25003628800},
	{n: 4838400, e: 25004838400},
}

var tstGovernanceData = []governanceData{
	{n: 123, e: 1}, // 1 is set at params.TestChainConfig
	{n: 604923, e: 1},
	{n: 1209723, e: 25000000000},
	{n: 1814523, e: 25001209600},
	{n: 2419323, e: 25001209600},
	{n: 3024123, e: 25002419200},
	{n: 3628923, e: 25002419200},
	{n: 4233723, e: 25003628800},
	{n: 4838523, e: 25003628800},
	{n: 5443323, e: 25004838400},
}

func TestSaveGovernance(t *testing.T) {
	gov := getGovernance()

	var MAXITEMS = int(10)

	// Set Data
	for i := 0; i < len(tstGovernanceInfo); i++ {
		blockNum := tstGovernanceInfo[i].n
		tstMap := copyMap(testGovernanceMap)

		// Make every stored governance map has a difference
		tstMap["governance.unitprice"] = tstGovernanceInfo[i].e
		src := NewGovernanceSet()
		delta := NewGovernanceSet()
		src.Import(tstMap)
		if err := gov.WriteGovernance(blockNum, src, delta); err != nil {
			t.Errorf("Error in storing governance: %v", err)
		}
	}

	// retrieve governance information. some will come from cache, others will be searched
	for i := 0; i < MAXITEMS; i++ {
		blockNum := tstGovernanceData[i].n
		_, data, err := gov.ReadGovernance(blockNum)
		if err == nil {
			if data["governance.unitprice"] != tstGovernanceData[i].e {
				t.Errorf("Data mismatch want %v, have %v for block %d", tstGovernanceData[i].e, data["governance.unitprice"], tstGovernanceData[i].n)
			}
		}
	}
}

type epochTest struct {
	v uint64
	e uint64
}

// Assume epoch is 30
var epochTestData = []epochTest{
	{0, 0},
	{30, 0},
	{60, 30},
	{90, 60},
	{120, 90},
	{150, 120},
	{180, 150},
	{210, 180},
	{240, 210},
	{270, 240},
	{300, 270},
	{330, 300},
	{360, 330},
}

func TestCalcGovernanceInfoBlock(t *testing.T) {
	for _, v := range epochTestData {
		res := CalcGovernanceInfoBlock(v.v, 30)
		if res != v.e {
			t.Errorf("Governance Block Number Mismatch: want %v, have %v", v.e, res)
		}
	}
}

func TestVoteValueNilInterface(t *testing.T) {
	gov := getGovernance()
	gVote := new(GovernanceVote)
	var test []byte

	gVote.Key = "istanbul.policy"
	// gVote.Value is an interface list
	{
		gVote.Value = []interface{}{[]byte{1, 2}}
		test, _ = rlp.EncodeToBytes(gVote)
		if err := rlp.DecodeBytes(test, gVote); err != nil {
			t.Fatal("RLP decode error")
		}

		// Parse vote.Value and make it has appropriate type
		_, err := gov.ParseVoteValue(gVote)
		assert.Equal(t, ErrValueTypeMismatch, err)
	}

	// gVote.Value is an empty interface list
	{
		gVote.Value = []interface{}{[]byte{}}
		test, _ = rlp.EncodeToBytes(gVote)
		if err := rlp.DecodeBytes(test, gVote); err != nil {
			t.Fatal("RLP decode error")
		}

		// Parse vote.Value and make it has appropriate type
		_, err := gov.ParseVoteValue(gVote)
		assert.Equal(t, ErrValueTypeMismatch, err)
	}

	// gVote.Value is nil
	{
		gVote.Value = nil
		test, _ = rlp.EncodeToBytes(gVote)
		if err := rlp.DecodeBytes(test, gVote); err != nil {
			t.Fatal("RLP decode error")
		}

		// Parse vote.Value and make it has appropriate type
		_, err := gov.ParseVoteValue(gVote)
		assert.Equal(t, ErrValueTypeMismatch, err)
	}

	// gVote.Value is uint8 list
	{
		gVote.Value = []uint8{0x11}
		test, _ = rlp.EncodeToBytes(gVote)
		if err := rlp.DecodeBytes(test, gVote); err != nil {
			t.Fatal("RLP decode error")
		}

		// Parse vote.Value and make it has appropriate type
		_, err := gov.ParseVoteValue(gVote)
		assert.NoError(t, err)
	}

	gVote.Key = "governance.addvalidator"
	{
		gVote.Value = [][]uint8{{0x3, 0x4}, {0x5, 0x6}}
		test, _ = rlp.EncodeToBytes(gVote)
		if err := rlp.DecodeBytes(test, gVote); err != nil {
			t.Fatal("RLP decode error")
		}
		// Parse vote.Value and make it has appropriate type
		_, err := gov.ParseVoteValue(gVote)
		assert.Equal(t, ErrValueTypeMismatch, err)
	}

	{
		gVote.Value = []uint8{0x1, 0x2, 0x3}
		test, _ = rlp.EncodeToBytes(gVote)
		if err := rlp.DecodeBytes(test, gVote); err != nil {
			t.Fatal("RLP decode error")
		}

		// Parse vote.Value and make it has appropriate type
		_, err := gov.ParseVoteValue(gVote)
		assert.NoError(t, err)
	}
}

func TestBaoBabGenesisHash(t *testing.T) {
	baobabHash := params.BaobabGenesisHash
	genesis := blockchain.DefaultBaobabGenesisBlock()
	genesis.Governance = blockchain.SetGenesisGovernance(genesis)
	blockchain.InitDeriveSha(genesis.Config.DeriveShaImpl)

	db := database.NewMemoryDBManager()
	block, _ := genesis.Commit(common.Hash{}, db)
	if block.Hash() != baobabHash {
		t.Errorf("Generated hash is not equal to Baobab's hash. Want %v, Have %v", baobabHash.String(), block.Hash().String())
	}
}

func TestCypressGenesisHash(t *testing.T) {
	cypressHash := params.CypressGenesisHash
	genesis := blockchain.DefaultGenesisBlock()
	genesis.Governance = blockchain.SetGenesisGovernance(genesis)
	blockchain.InitDeriveSha(genesis.Config.DeriveShaImpl)

	db := database.NewMemoryDBManager()
	block, _ := genesis.Commit(common.Hash{}, db)
	if block.Hash() != cypressHash {
		t.Errorf("Generated hash is not equal to Cypress's hash. Want %v, Have %v", cypressHash.String(), block.Hash().String())
	}
}

func TestGovernance_initializeCache(t *testing.T) {
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	config := getTestConfig()
	config.Istanbul.Epoch = 30

	gov := NewGovernanceInitialize(config, dbm)

	testData := []struct {
		// test input
		governanceUpdateNum int
		blockNums           []uint64
		unitPrices          []uint64
		currentBlockNum     int64
		// expected result
		unitPriceInCurrentSet uint64
		actualGovernanceBlock uint64
	}{
		{0, []uint64{0}, []uint64{1}, 20, 1, 0},
		{0, []uint64{0}, []uint64{1}, 50, 1, 0},
		{0, []uint64{0}, []uint64{1}, 80, 1, 0},
		{3, []uint64{0, 30, 60, 90}, []uint64{1, 2, 3, 4}, 90, 3, 60},
		{3, []uint64{0, 30, 60, 90}, []uint64{1, 2, 3, 4}, 110, 3, 60},
		{3, []uint64{0, 30, 60, 90}, []uint64{1, 2, 3, 4}, 130, 4, 90},
	}

	for _, tc := range testData {
		// 1. initializing

		// store governance items to the governance db for 'tc.governanceUpdateNum' times
		for idx := 1; idx <= tc.governanceUpdateNum; idx++ {
			config.UnitPrice = tc.unitPrices[idx]

			src := GetGovernanceItemsFromChainConfig(config)
			delta := NewGovernanceSet()

			if ret := gov.WriteGovernance(tc.blockNums[idx], src, delta); ret != nil {
				t.Errorf("Error in testing WriteGovernance, %v", ret)
			}
		}
		// store head block (fake block, it only contains block number) to the headerDB
		header := &types.Header{Number: big.NewInt(tc.currentBlockNum)}
		gov.db.WriteHeadBlockHash(header.Hash())
		gov.db.WriteHeader(header)

		// reset idxCache and itemCache. purpose - reproduce the environment of the restarted node
		gov.idxCache = nil
		gov.itemCache.Purge()

		// 2. call initializeCache
		err := gov.initializeCache()

		// 3. check the affected values with expected results
		assert.NoError(t, err)

		v, ok := gov.currentSet.GetValue(GovernanceKeyMap["governance.unitprice"])
		assert.True(t, ok)
		assert.Equal(t, tc.unitPriceInCurrentSet, v)

		assert.Equal(t, tc.blockNums, gov.idxCache) // the order should be same
		assert.True(t, gov.itemCache.Contains(getGovernanceCacheKey(tc.blockNums[tc.governanceUpdateNum])))
		assert.Equal(t, tc.actualGovernanceBlock, gov.actualGovernanceBlock.Load().(uint64))
	}
}

func TestWriteGovernance_idxCache(t *testing.T) {
	gov := getGovernance()

	tstMap := copyMap(testGovernanceMap)

	src := NewGovernanceSet()
	delta := NewGovernanceSet()
	src.Import(tstMap)

	blockNum := []uint64{30, 30, 60, 60, 50}

	for _, num := range blockNum {
		if ret := gov.WriteGovernance(num, src, delta); ret != nil {
			t.Errorf("Error in testing WriteGovernance, %v", ret)
		}
	}

	// idxCache should have 0, 30 and 60
	if len(gov.idxCache) != 3 || gov.idxCache[len(gov.idxCache)-1] != 60 {
		t.Errorf("idxCache has wrong value")
	}
}

func getTestValidators() []common.Address {
	return []common.Address{
		common.HexToAddress("0x414790CA82C14A8B975cEBd66098c3dA590bf969"), // Node Address for test
		common.HexToAddress("0x604973C51f6389dF2782E018000c3AC1257dee90"),
		common.HexToAddress("0x5Ac1689ae5F521B05145C5Cd15a3E8F6ab39Af19"),
		common.HexToAddress("0x0688CaC68bbF7c1a0faedA109c668a868BEd855E"),
	}
}

func getTestDemotedValidators() []common.Address {
	return []common.Address{
		common.HexToAddress("0x3BB17a8A4f915cC9A8CAAcdC062Ef9b903511Ffa"),
		common.HexToAddress("0x82588D33A48e6Bda012714f1C680d254ff607472"),
		common.HexToAddress("0xceB7ADDFBa9665d8767173D47dE4453D7b7B900D"),
		common.HexToAddress("0x38Ea854792EB956620E53090E8bc4e5C5C917123"),
	}
}

func getTestRewards() []common.Address {
	return []common.Address{
		common.HexToAddress("0x2A35FE72F847aa0B509e4055883aE90c87558AaD"),
		common.HexToAddress("0xF91B8EBa583C7fa603B400BE17fBaB7629568A4a"),
		common.HexToAddress("0x240ed27c8bDc9Bb6cA08fa3D239699Fba525d05a"),
		common.HexToAddress("0x3B980293396Fb0e827929D573e3e42d2EA902502"),
	}
}

func getTestVotingPowers(num int) []uint64 {
	vps := make([]uint64, 0, num)
	for i := 0; i < num; i++ {
		vps = append(vps, 1000)
	}
	return vps
}

const (
	GovernanceModeBallot = "ballot"
)

func TestGovernance_HandleGovernanceVote_None_mode(t *testing.T) {
	// Create ValidatorSet
	validators := getTestValidators()
	demotedValidators := getTestDemotedValidators()
	rewards := getTestRewards()

	blockCounter := common.Big0
	valSet := validator.NewWeightedCouncil(validators, demotedValidators, rewards, getTestVotingPowers(len(validators)), nil, istanbul.WeightedRandom, 21, 0, 0, nil)
	gov := getGovernance()
	gov.nodeAddress.Store(validators[len(validators)-1])

	votes := make([]GovernanceVote, 0)
	tally := make([]GovernanceTallyItem, 0)

	proposer := validators[0]
	self := validators[len(validators)-1]
	header := &types.Header{}

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test for "none" mode
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.BlockScore = common.Big1
	gov.AddVote("governance.unitprice", uint64(22000))
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, self)
	gov.RemoveVote("governance.unitprice", uint64(22000), 0)

	if _, ok := gov.changeSet.items["governance.unitprice"]; !ok {
		t.Errorf("Vote had to be applied but it wasn't")
	}
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test for "istanbul.timeout" in "none" mode
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.BlockScore = common.Big1

	newValue := istanbul.DefaultConfig.Timeout + uint64(10000)
	gov.AddVote("istanbul.timeout", newValue)
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, self)
	gov.RemoveVote("istanbul.timeout", newValue, 0)
	assert.Equal(t, istanbul.DefaultConfig.Timeout, newValue, "Vote had to be applied but it wasn't")

	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test removing a validator
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	gov.AddVote("governance.removevalidator", validators[1].String())
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, self)
	gov.RemoveVote("governance.removevalidator", validators[1], 0)
	if i, _ := valSet.GetByAddress(validators[1]); i != -1 {
		t.Errorf("Validator removal failed, %d validators remains", valSet.Size())
	}
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test removing a non-existing validator
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	gov.AddVote("governance.removevalidator", validators[1].String())
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, proposer) // self = proposer
	// check if casted
	if !gov.voteMap.items["governance.removevalidator"].Casted {
		t.Errorf("Removing a non-existing validator failed")
	}
	gov.RemoveVote("governance.removevalidator", validators[1], 0)
	if i, _ := valSet.GetByAddress(validators[1]); i != -1 {
		t.Errorf("Removing a non-existing validator failed, %d validators remains", valSet.Size())
	}
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test adding a validator
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	gov.AddVote("governance.addvalidator", validators[1].String())
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, self)
	gov.RemoveVote("governance.addvalidator", validators[1], 0)
	if i, _ := valSet.GetByAddress(validators[1]); i == -1 {
		t.Errorf("Validator addition failed, %d validators remains", valSet.Size())
	}
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test adding an existing validator
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	gov.AddVote("governance.addvalidator", validators[1].String())
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, proposer) // self = proposer
	// check if casted
	if !gov.voteMap.items["governance.addvalidator"].Casted {
		t.Errorf("Adding an existing validator failed")
	}
	gov.RemoveVote("governance.addvalidator", validators[1], 0)
	if i, _ := valSet.GetByAddress(validators[1]); i == -1 {
		t.Errorf("Validator addition failed, %d validators remains", valSet.Size())
	}
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test removing a demoted validator
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	gov.AddVote("governance.removevalidator", demotedValidators[1].String())
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, self)
	gov.RemoveVote("governance.removevalidator", demotedValidators[1], 0)
	if i, _ := valSet.GetDemotedByAddress(demotedValidators[1]); i != -1 {
		t.Errorf("Demoted validator removal failed, %d demoted validators remains", len(valSet.DemotedList()))
	}
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test adding a demoted validator
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	gov.AddVote("governance.addvalidator", demotedValidators[1].String())
	header.Vote = gov.GetEncodedVote(proposer, blockCounter.Uint64())

	gov.HandleGovernanceVote(valSet, votes, tally, header, proposer, self)
	gov.RemoveVote("governance.addvalidator", demotedValidators[1], 0)
	// At first, demoted validator is added to the validators, but it will be refreshed right after
	// So, we here check only if the adding demoted validator to validators
	if i, _ := valSet.GetByAddress(demotedValidators[1]); i == -1 {
		t.Errorf("Demoted validator addition failed, %d demoted validators remains", len(valSet.DemotedList()))
	}
	gov.voteMap.Clear()
}

func TestGovernance_HandleGovernanceVote_Ballot_mode(t *testing.T) {
	// Create ValidatorSet
	validators := getTestValidators()
	demotedValidators := getTestDemotedValidators()
	rewards := getTestRewards()

	blockCounter := common.Big0
	var valSet istanbul.ValidatorSet
	valSet = validator.NewWeightedCouncil(validators, demotedValidators, rewards, getTestVotingPowers(len(validators)), nil, istanbul.WeightedRandom, 21, 0, 0, nil)

	config := getTestConfig()
	config.Governance.GovernanceMode = GovernanceModeBallot
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	gov := NewGovernanceInitialize(config, dbm)
	gov.nodeAddress.Store(validators[len(validators)-1])

	votes := make([]GovernanceVote, 0)
	tally := make([]GovernanceTallyItem, 0)

	self := validators[len(validators)-1]
	header := &types.Header{}

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test for "ballot" mode
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.BlockScore = common.Big1
	gov.AddVote("governance.unitprice", uint64(22000))

	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], self)

	header.Vote = gov.GetEncodedVote(validators[1], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[1], self)

	if _, ok := gov.changeSet.items["governance.unitprice"]; ok {
		t.Errorf("Vote shouldn't be applied yet but it was applied")
	}

	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], self)
	if _, ok := gov.changeSet.items["governance.unitprice"]; !ok {
		t.Errorf("Vote should be applied but it was not")
	}

	gov.RemoveVote("governance.unitprice", uint64(22000), blockCounter.Uint64())
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test for "istanbul.timeout" in "ballot" mode
	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.BlockScore = common.Big1
	newValue := istanbul.DefaultConfig.Timeout + uint64(10000)
	gov.AddVote("istanbul.timeout", newValue)

	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], self)

	header.Vote = gov.GetEncodedVote(validators[1], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[1], self)

	assert.NotEqual(t, istanbul.DefaultConfig.Timeout, newValue, "Vote shouldn't be applied yet but it was applied")

	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], self)

	assert.Equal(t, istanbul.DefaultConfig.Timeout, newValue, "Vote should be applied but it was not")
	gov.RemoveVote("istanbul.timeout", newValue, blockCounter.Uint64())
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test removing a validator, because there are 4 nodes 3 votes are required to remove a validator
	gov.AddVote("governance.removevalidator", validators[1].String())

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], self)

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], self)
	if i, _ := valSet.GetByAddress(validators[1]); i == -1 {
		t.Errorf("Validator removal shouldn't be done yet, %d validators remains", valSet.Size())
	}

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[3], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[3], self)

	if i, _ := valSet.GetByAddress(validators[1]); i != -1 {
		t.Errorf("Validator removal failed, %d validators remains", valSet.Size())
	}
	gov.RemoveVote("governance.removevalidator", validators[1], blockCounter.Uint64())
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test removing a non-existing validator. Because there are 3 nodes, 2 votes are required to remove a validator
	gov.AddVote("governance.removevalidator", validators[1].String())

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], validators[0])
	// check if casted
	if !gov.voteMap.items["governance.removevalidator"].Casted {
		t.Errorf("Removing a non-existing validator failed")
	}

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], validators[2])
	// check if casted
	if !gov.voteMap.items["governance.removevalidator"].Casted {
		t.Errorf("Removing a non-existing validator failed")
	}

	gov.RemoveVote("governance.removevalidator", validators[1], blockCounter.Uint64())
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test adding a validator, because there are 3 nodes 2 plus votes are required to add a new validator
	gov.AddVote("governance.addvalidator", validators[1].String())

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], self)
	if i, _ := valSet.GetByAddress(validators[1]); i != -1 {
		t.Errorf("Validator addition shouldn't be done yet, %d validators remains", valSet.Size())
	}

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], self)

	if i, _ := valSet.GetByAddress(validators[1]); i == -1 {
		t.Errorf("Validator addition failed, %d validators remains", valSet.Size())
	}
	gov.RemoveVote("governance.addvalidator", validators[1], blockCounter.Uint64())
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test adding an existing validator, because there are 3 nodes 2 plus votes are required to add a new validator
	gov.AddVote("governance.addvalidator", validators[1].String())

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], validators[0])
	// check if casted
	if !gov.voteMap.items["governance.addvalidator"].Casted {
		t.Errorf("Adding an existing validator failed")
	}

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], validators[2])
	// check if casted
	if !gov.voteMap.items["governance.addvalidator"].Casted {
		t.Errorf("Adding an existing validator failed")
	}

	gov.RemoveVote("governance.addvalidator", validators[1], blockCounter.Uint64())
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test removing a demoted validator, because there are 4 nodes 3 votes are required to remove a demoted validator
	gov.AddVote("governance.removevalidator", demotedValidators[1].String())

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], self)

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], self)
	if i, _ := valSet.GetDemotedByAddress(demotedValidators[1]); i == -1 {
		t.Errorf("Demoted validator removal shouldn't be done yet, %d validators remains", len(valSet.DemotedList()))
	}

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[3], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[3], self)

	if i, _ := valSet.GetDemotedByAddress(demotedValidators[1]); i != -1 {
		t.Errorf("Demoted validator removal failed, %d validators remains", len(valSet.DemotedList()))
	}
	gov.voteMap.Clear()

	//////////////////////////////////////////////////////////////////////////////////////////////////////////
	// Test adding a demoted validator, because there are 4 nodes 3 votes are required to add a demoted validator
	gov.AddVote("governance.addvalidator", demotedValidators[1].String())

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[0], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[0], self)
	if i, _ := valSet.GetByAddress(demotedValidators[1]); i != -1 {
		t.Errorf("Validator addition shouldn't be done yet, %d validators remains", len(valSet.DemotedList()))
	}

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[2], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[2], self)

	header.Number = blockCounter.Add(blockCounter, common.Big1)
	header.Vote = gov.GetEncodedVote(validators[3], blockCounter.Uint64())
	valSet, votes, tally = gov.HandleGovernanceVote(valSet, votes, tally, header, validators[3], self)

	// At first, demoted validator is added to the validators, but it will be refreshed right after
	// So, we here check only if the adding demoted validator to validators
	if i, _ := valSet.GetByAddress(demotedValidators[1]); i == -1 {
		t.Errorf("Demoted validator addition failed, %d validators remains", len(valSet.DemotedList()))
	}
	gov.voteMap.Clear()
}

func TestGovernance_checkVote(t *testing.T) {
	// Create ValidatorSet
	council := getTestValidators()
	validators := []common.Address{council[0], council[1]}
	demotedValidators := []common.Address{council[2], council[3]}

	valSet := validator.NewWeightedCouncil(validators, demotedValidators, nil, getTestVotingPowers(len(validators)), nil, istanbul.WeightedRandom, 21, 0, 0, nil)

	config := getTestConfig()
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	gov := NewGovernanceInitialize(config, dbm)

	unknown := common.HexToAddress("0xa")

	// for adding validator
	assert.True(t, gov.checkVote(unknown, true, valSet))
	assert.True(t, gov.checkVote(validators[0], false, valSet))
	assert.True(t, gov.checkVote(demotedValidators[0], false, valSet))

	// for removing validator
	assert.False(t, gov.checkVote(unknown, false, valSet))
	assert.False(t, gov.checkVote(validators[1], true, valSet))
	assert.False(t, gov.checkVote(demotedValidators[1], true, valSet))
}

func TestGovernance_VerifyGovernance(t *testing.T) {
	gov := getGovernance()
	vote := GovernanceVote{
		Key:   "governance.governingnode",
		Value: common.HexToAddress("000000000000000000000000000abcd000000000"),
	}
	gov.updateChangeSet(vote)

	// consensus/istanbul/backend/engine.go:Prepare()
	// Correct case
	g := gov.GetGovernanceChange()
	j, err := json.Marshal(g)
	assert.Nil(t, err)
	r, err := rlp.EncodeToBytes(j)
	assert.Nil(t, err)
	err = gov.VerifyGovernance(r)
	assert.Nil(t, err)

	// Value mismatch
	g = gov.GetGovernanceChange()
	g["governance.governingnode"] = "000000000000000000000000000abcd000001111"
	j, err = json.Marshal(g)
	assert.Nil(t, err)
	r, err = rlp.EncodeToBytes(j)
	assert.Nil(t, err)
	err = gov.VerifyGovernance(r)
	assert.Equal(t, ErrVoteValueMismatch, err)

	// Type mismatch
	g = gov.GetGovernanceChange()
	g["governance.governingnode"] = 123
	j, err = json.Marshal(g)
	assert.Nil(t, err)
	r, err = rlp.EncodeToBytes(j)
	assert.Nil(t, err)
	err = gov.VerifyGovernance(r)
	assert.Equal(t, ErrVoteValueMismatch, err)

	// Length mismatch
	g = gov.GetGovernanceChange()
	g["governance.governingnode"] = 123
	g["istanbul.epoch"] = uint64(10000)
	j, err = json.Marshal(g)
	assert.Nil(t, err)
	r, err = rlp.EncodeToBytes(j)
	assert.Nil(t, err)
	err = gov.VerifyGovernance(r)
	assert.Equal(t, ErrVoteValueMismatch, err)
}
