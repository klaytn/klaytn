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
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/stretchr/testify/assert"
	"math/big"
	"reflect"
	"testing"
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
	{k: "istanbul.committeesize", v: uint64(7), e: true},
	{k: "istanbul.committeesize", v: float64(7.0), e: true},
	{k: "istanbul.committeesize", v: float64(7.1), e: false},
	{k: "istanbul.committeesize", v: "7", e: false},
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
	{k: "governance.unitprice", v: float64(0.0), e: true},
	{k: "governance.unitprice", v: uint64(25000000000), e: true},
	{k: "governance.unitprice", v: float64(-10), e: false},
	{k: "governance.unitprice", v: "25000000000", e: false},
	{k: "reward.useginicoeff", v: false, e: true},
	{k: "reward.useginicoeff", v: true, e: true},
	{k: "reward.useginicoeff", v: "true", e: false},
	{k: "reward.useginicoeff", v: 0, e: false},
	{k: "reward.useginicoeff", v: 1, e: false},
	{k: "reward.mintingamount", v: "9600000000000000000", e: true},
	{k: "reward.mintingamount", v: "0", e: true},
	{k: "reward.mintingamount", v: 96000, e: false},
	{k: "reward.mintingamount", v: "many", e: false},
	{k: "reward.ratio", v: "30/40/30", e: true},
	{k: "reward.ratio", v: "10/10/80", e: true},
	{k: "reward.ratio", v: "30/70", e: false},
	{k: "reward.ratio", v: "30.5/40/29.5", e: false},
	{k: "reward.ratio", v: "30.5/40/30.5", e: false},
	{k: "reward.deferredtxfee", v: true, e: true},
	{k: "reward.deferredtxfee", v: false, e: true},
	{k: "reward.deferredtxfee", v: 0, e: false},
	{k: "reward.deferredtxfee", v: 1, e: false},
	{k: "reward.deferredtxfee", v: "true", e: false},
	{k: "reward.minimumstake", v: "2000000000000000000000000", e: true},
	{k: "reward.minimumstake", v: 200000000000000, e: false},
	{k: "reward.stakingupdateinterval", v: uint64(20), e: false},
	{k: "reward.proposerupdateinterval", v: uint64(20), e: false},
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
}

func getTestConfig() *params.ChainConfig {
	config := params.TestChainConfig
	config.Governance = GetDefaultGovernanceConfig(params.UseIstanbul)
	config.Istanbul = GetDefaultIstanbulConfig()
	return config
}

func getGovernance() *Governance {
	dbm := database.NewDBManager(&database.DBConfig{DBType: database.MemoryDB})
	config := getTestConfig()
	return NewGovernance(config, dbm)
}

func TestNewGovernance(t *testing.T) {
	config := getTestConfig()
	tstGovernance := NewGovernance(config, nil)

	if !reflect.DeepEqual(tstGovernance.ChainConfig, config) {
		t.Errorf("New governance's config is not same as the given one")
	}
}

func TestGetDefaultGovernanceConfig(t *testing.T) {
	tstGovernance := GetDefaultGovernanceConfig(params.UseIstanbul)

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

	if tstGovernance.Reward.MintingAmount.Cmp(big.NewInt(params.DefaultMintingAmount)) != 0 {
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
		if ret != val.e {
			t.Errorf("Want %v, got %v for %v and %v", val.e, ret, val.k, val.v)
		}
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

func countUncastedVote(data VoteSet) int {
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
	}
	gov.ClearVotes(0)
	if gov.voteMap.Size() != 0 {
		t.Errorf("Want 0, got %v after clearing votes", gov.voteMap.Size())
	}
}

func TestGovernance_GetEncodedVote(t *testing.T) {
	var err error
	gov := getGovernance()

	for _, val := range goodVotes {
		_ = gov.AddVote(val.k, val.v)
	}

	l := gov.voteMap.Size()
	for i := 0; i > l; i++ {
		voteData := gov.GetEncodedVote(common.HexToAddress("0x1234567890123456789012345678901234567890"), 1000)
		v := new(GovernanceVote)
		rlp.DecodeBytes(voteData, &v)

		if v, err = gov.ParseVoteValue(v); err != nil {
			assert.Equal(t, nil, err)
		}

		if v.Value != gov.voteMap.GetValue(v.Key).Value {
			t.Errorf("Encoded vote and Decoded vote are different! Encoded: %v, Decoded: %v\n", gov.voteMap.GetValue(v.Key).Value, v.Value)
		}
		gov.RemoveVote(v.Key, v.Value, 1000)
	}
}

func TestGovernance_ParseVoteValue(t *testing.T) {
	var err error
	gov := getGovernance()

	addr := common.HexToAddress("0x1234567890123456789012345678901234567890")
	for _, val := range goodVotes {
		v := &GovernanceVote{
			Key:       val.k,
			Value:     val.v,
			Validator: addr,
		}

		b, _ := rlp.EncodeToBytes(v)

		d := new(GovernanceVote)
		rlp.DecodeBytes(b, d)

		if d, err = gov.ParseVoteValue(d); err != nil {
			assert.Equal(t, nil, err)
		}

		if !reflect.DeepEqual(v, d) {
			t.Errorf("Parse was not successful! %v %v \n", v, d)
		}
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
		assert.Equal(t, nil, err)
	}
}

func TestBaoBabGenesisHash(t *testing.T) {
	baobabHash := common.HexToHash("0xe33ff05ceec2581ca9496f38a2bf9baad5d4eed629e896ccb33d1dc991bc4b4a")
	genesis := blockchain.DefaultBaobabGenesisBlock()
	genesis.Governance = blockchain.SetGenesisGovernance(genesis)
	blockchain.InitDeriveSha(genesis.Config.DeriveShaImpl)

	db := database.NewMemoryDBManager()
	block, _ := genesis.Commit(db)
	if block.Hash() != baobabHash {
		t.Errorf("Generated hash is not equal to Baobab's hash. Want %v, Have %v", baobabHash.String(), block.Hash().String())
	}
}

func TestCypressGenesisHash(t *testing.T) {
	cypressHash := common.HexToHash("0xc72e5293c3c3ba38ed8ae910f780e4caaa9fb95e79784f7ab74c3c262ea7137e")
	genesis := blockchain.DefaultGenesisBlock()
	genesis.Governance = blockchain.SetGenesisGovernance(genesis)
	blockchain.InitDeriveSha(genesis.Config.DeriveShaImpl)

	db := database.NewMemoryDBManager()
	block, _ := genesis.Commit(db)
	if block.Hash() != cypressHash {
		t.Errorf("Generated hash is not equal to Baobab's hash. Want %v, Have %v", cypressHash.String(), block.Hash().String())
	}
}
