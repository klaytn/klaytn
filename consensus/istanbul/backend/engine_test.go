// Modifications Copyright 2020 The klaytn Authors
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
// This file is derived from quorum/consensus/istanbul/backend/engine_test.go (2020/04/16).
// Modified and improved for the klaytn development.

package backend

import (
	"bytes"
	"crypto/ecdsa"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/consensus/istanbul/core"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/reward"
	"github.com/klaytn/klaytn/rlp"
	"github.com/stretchr/testify/assert"
)

// These variables are the global variables of the test blockchain.
var (
	nodeKeys []*ecdsa.PrivateKey
	addrs    []common.Address
)

// These are the types in order to add a custom configuration of the test chain.
// You may need to create a configuration type if necessary.
type istanbulCompatibleBlock *big.Int
type minimumStake *big.Int
type stakingUpdateInterval uint64
type proposerUpdateInterval uint64
type proposerPolicy uint64
type governanceMode string
type epoch uint64
type subGroupSize uint64
type blockPeriod uint64

// makeCommittedSeals returns a list of committed seals for the global variable nodeKeys.
func makeCommittedSeals(hash common.Hash) [][]byte {
	committedSeals := make([][]byte, len(nodeKeys))
	hashData := crypto.Keccak256(core.PrepareCommittedSeal(hash))
	for i, key := range nodeKeys {
		sig, _ := crypto.Sign(hashData, key)
		committedSeals[i] = make([]byte, types.IstanbulExtraSeal)
		copy(committedSeals[i][:], sig)
	}
	return committedSeals
}

// Include a node from the global nodeKeys and addrs
func includeNode(addr common.Address, key *ecdsa.PrivateKey) {
	for _, a := range addrs {
		if a.String() == addr.String() {
			// already exists
			return
		}
	}
	nodeKeys = append(nodeKeys, key)
	addrs = append(addrs, addr)
}

// Exclude a node from the global nodeKeys and addrs
func excludeNodeByAddr(target common.Address) {
	for i, a := range addrs {
		if a.String() == target.String() {
			nodeKeys = append(nodeKeys[:i], nodeKeys[i+1:]...)
			addrs = append(addrs[:i], addrs[i+1:]...)
			break
		}
	}
}

// in this test, we can set n to 1, and it means we can process Istanbul and commit a
// block by one node. Otherwise, if n is larger than 1, we have to generate
// other fake events to process Istanbul.
func newBlockChain(n int, items ...interface{}) (*blockchain.BlockChain, *backend) {
	// generate a genesis block
	genesis := blockchain.DefaultGenesisBlock()
	config := *params.TestChainConfig // copy test chain config which may be modified
	genesis.Config = &config
	genesis.Timestamp = uint64(time.Now().Unix())

	var (
		key    *ecdsa.PrivateKey
		period = istanbul.DefaultConfig.BlockPeriod
	)
	// force enable Istanbul engine and governance
	genesis.Config.Istanbul = params.GetDefaultIstanbulConfig()
	genesis.Config.Governance = params.GetDefaultGovernanceConfig(params.UseIstanbul)
	for _, item := range items {
		switch v := item.(type) {
		case istanbulCompatibleBlock:
			genesis.Config.IstanbulCompatibleBlock = v
		case proposerPolicy:
			genesis.Config.Istanbul.ProposerPolicy = uint64(v)
		case epoch:
			genesis.Config.Istanbul.Epoch = uint64(v)
		case subGroupSize:
			genesis.Config.Istanbul.SubGroupSize = uint64(v)
		case minimumStake:
			genesis.Config.Governance.Reward.MinimumStake = v
		case stakingUpdateInterval:
			genesis.Config.Governance.Reward.StakingUpdateInterval = uint64(v)
		case proposerUpdateInterval:
			genesis.Config.Governance.Reward.ProposerUpdateInterval = uint64(v)
		case governanceMode:
			genesis.Config.Governance.GovernanceMode = string(v)
		case *ecdsa.PrivateKey:
			key = v
		case blockPeriod:
			period = uint64(v)
		}
	}
	nodeKeys = make([]*ecdsa.PrivateKey, n)
	addrs = make([]common.Address, n)

	var b *backend
	if len(items) != 0 {
		b = newTestBackendWithConfig(genesis.Config, period, key)
	} else {
		b = newTestBackend()
	}

	nodeKeys[0] = b.privateKey
	addrs[0] = b.address // if governance mode is single, this address is the governing node address
	for i := 1; i < n; i++ {
		nodeKeys[i], _ = crypto.GenerateKey()
		addrs[i] = crypto.PubkeyToAddress(nodeKeys[i].PublicKey)
	}

	appendValidators(genesis, addrs)

	genesis.MustCommit(b.db)

	bc, err := blockchain.NewBlockChain(b.db, nil, genesis.Config, b, vm.Config{})
	if err != nil {
		panic(err)
	}
	if b.Start(bc, bc.CurrentBlock, bc.HasBadBlock) != nil {
		panic(err)
	}

	return bc, b
}

func appendValidators(genesis *blockchain.Genesis, addrs []common.Address) {
	if len(genesis.ExtraData) < types.IstanbulExtraVanity {
		genesis.ExtraData = append(genesis.ExtraData, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)...)
	}
	genesis.ExtraData = genesis.ExtraData[:types.IstanbulExtraVanity]

	ist := &types.IstanbulExtra{
		Validators:    addrs,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	istPayload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		panic("failed to encode istanbul extra")
	}
	genesis.ExtraData = append(genesis.ExtraData, istPayload...)
}

func makeHeader(parent *types.Block, config *istanbul.Config) *types.Header {
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     parent.Number().Add(parent.Number(), common.Big1),
		GasUsed:    0,
		Extra:      parent.Extra(),
		Time:       new(big.Int).Add(parent.Time(), new(big.Int).SetUint64(config.BlockPeriod)),
		BlockScore: defaultBlockScore,
	}
	return header
}

func makeBlock(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	block := makeBlockWithoutSeal(chain, engine, parent)
	stopCh := make(chan struct{})
	result, err := engine.Seal(chain, block, stopCh)
	if err != nil {
		panic(err)
	}
	return result
}

// makeBlockWithSeal creates a block with the proposer seal as well as all committed seals of validators.
func makeBlockWithSeal(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	blockWithoutSeal := makeBlockWithoutSeal(chain, engine, parent)

	// add proposer seal for the block
	block, err := engine.updateBlock(nil, blockWithoutSeal)
	if err != nil {
		panic(err)
	}

	// write validators committed seals to the block
	header := block.Header()
	committedSeals := makeCommittedSeals(block.Hash())
	err = writeCommittedSeals(header, committedSeals)
	if err != nil {
		panic(err)
	}
	block = block.WithSeal(header)

	return block
}

func makeBlockWithoutSeal(chain *blockchain.BlockChain, engine *backend, parent *types.Block) *types.Block {
	header := makeHeader(parent, engine.config)
	if err := engine.Prepare(chain, header); err != nil {
		panic(err)
	}
	state, _ := chain.StateAt(parent.Root())
	block, _ := engine.Finalize(chain, header, state, nil, nil)
	return block
}

func TestPrepare(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	header := makeHeader(chain.Genesis(), engine.config)
	err := engine.Prepare(chain, header)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	header.ParentHash = common.HexToHash("0x1234567890")
	err = engine.Prepare(chain, header)
	if err != consensus.ErrUnknownAncestor {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrUnknownAncestor)
	}
}

func TestSealStopChannel(t *testing.T) {
	chain, engine := newBlockChain(4)
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	stop := make(chan struct{}, 1)
	eventSub := engine.EventMux().Subscribe(istanbul.RequestEvent{})
	eventLoop := func() {
		select {
		case ev := <-eventSub.Chan():
			_, ok := ev.Data.(istanbul.RequestEvent)
			if !ok {
				t.Errorf("unexpected event comes: %v", reflect.TypeOf(ev.Data))
			}
			stop <- struct{}{}
		}
		eventSub.Unsubscribe()
	}
	go eventLoop()

	finalBlock, err := engine.Seal(chain, block, stop)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}

	if finalBlock != nil {
		t.Errorf("block mismatch: have %v, want nil", finalBlock)
	}
}

func TestSealCommitted(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	expectedBlock, _ := engine.updateBlock(engine.chain.GetHeader(block.ParentHash(), block.NumberU64()-1), block)

	actualBlock, err := engine.Seal(chain, block, make(chan struct{}))
	if err != nil {
		t.Errorf("error mismatch: have %v, want %v", err, expectedBlock)
	}

	if actualBlock.Hash() != expectedBlock.Hash() {
		t.Errorf("hash mismatch: have %v, want %v", actualBlock.Hash(), expectedBlock.Hash())
	}
}

func TestVerifyHeader(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	// errEmptyCommittedSeals case
	block := makeBlockWithoutSeal(chain, engine, chain.Genesis())
	block, _ = engine.updateBlock(chain.Genesis().Header(), block)
	err := engine.VerifyHeader(chain, block.Header(), false)
	if err != errEmptyCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, errEmptyCommittedSeals)
	}

	// short extra data
	header := block.Header()
	header.Extra = []byte{}
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}
	// incorrect extra format
	header.Extra = []byte("0000000000000000000000000000000012300000000000000000000000000000000000000000000000000000000000000000")
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidExtraDataFormat {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidExtraDataFormat)
	}

	// invalid difficulty
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.BlockScore = big.NewInt(2)
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidBlockScore {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidBlockScore)
	}

	// invalid timestamp
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = new(big.Int).Add(chain.Genesis().Time(), new(big.Int).SetUint64(engine.config.BlockPeriod-1))
	err = engine.VerifyHeader(chain, header, false)
	if err != errInvalidTimestamp {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidTimestamp)
	}

	// future block
	block = makeBlockWithoutSeal(chain, engine, chain.Genesis())
	header = block.Header()
	header.Time = new(big.Int).Add(big.NewInt(now().Unix()), new(big.Int).SetUint64(10))
	err = engine.VerifyHeader(chain, header, false)
	if err != consensus.ErrFutureBlock {
		t.Errorf("error mismatch: have %v, want %v", err, consensus.ErrFutureBlock)
	}

	// TODO-Klaytn: add more tests for header.Governance, header.Rewardbase, header.Vote
}

func TestVerifySeal(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	genesis := chain.Genesis()

	// cannot verify genesis
	err := engine.VerifySeal(chain, genesis.Header())
	if err != errUnknownBlock {
		t.Errorf("error mismatch: have %v, want %v", err, errUnknownBlock)
	}
	block := makeBlock(chain, engine, genesis)

	// clean cache before testing
	signatureAddresses.Purge()

	// change block content
	header := block.Header()
	header.Number = big.NewInt(4)
	block1 := block.WithSeal(header)
	err = engine.VerifySeal(chain, block1.Header())
	if err != errUnauthorized {
		t.Errorf("error mismatch: have %v, want %v", err, errUnauthorized)
	}

	// clean cache before testing
	signatureAddresses.Purge()

	// unauthorized users but still can get correct signer address
	engine.privateKey, _ = crypto.GenerateKey()
	err = engine.VerifySeal(chain, block.Header())
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
}

func TestVerifyHeaders(t *testing.T) {
	chain, engine := newBlockChain(1)
	defer engine.Stop()

	genesis := chain.Genesis()

	// success case
	headers := []*types.Header{}
	blocks := []*types.Block{}
	size := 100

	for i := 0; i < size; i++ {
		var b *types.Block
		if i == 0 {
			b = makeBlockWithoutSeal(chain, engine, genesis)
			b, _ = engine.updateBlock(genesis.Header(), b)
			engine.db.WriteHeader(b.Header())
		} else {
			b = makeBlockWithoutSeal(chain, engine, blocks[i-1])
			b, _ = engine.updateBlock(blocks[i-1].Header(), b)
			engine.db.WriteHeader(b.Header())
		}
		blocks = append(blocks, b)
		headers = append(headers, blocks[i].Header())
	}

	// proceed time to avoid future block errors
	now = func() time.Time {
		return time.Unix(headers[size-1].Time.Int64(), 0)
	}
	defer func() {
		now = time.Now
	}()

	_, results := engine.VerifyHeaders(chain, headers, nil)
	const timeoutDura = 2 * time.Second
	timeout := time.NewTimer(timeoutDura)
	index := 0
OUT1:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT1
				}
			}
			index++
			if index == size {
				break OUT1
			}
		case <-timeout.C:
			break OUT1
		}
	}
	// abort cases
	abort, results := engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDura)
	index = 0
OUT2:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					t.Errorf("error mismatch: have %v, want errEmptyCommittedSeals|errInvalidCommittedSeals", err)
					break OUT2
				}
			}
			index++
			if index == 5 {
				abort <- struct{}{}
			}
			if index >= size {
				t.Errorf("verifyheaders should be aborted")
				break OUT2
			}
		case <-timeout.C:
			break OUT2
		}
	}
	// error header cases
	headers[2].Number = big.NewInt(100)
	abort, results = engine.VerifyHeaders(chain, headers, nil)
	timeout = time.NewTimer(timeoutDura)
	index = 0
	errors := 0
	expectedErrors := 2
OUT3:
	for {
		select {
		case err := <-results:
			if err != nil {
				if err != errEmptyCommittedSeals && err != errInvalidCommittedSeals {
					errors++
				}
			}
			index++
			if index == size {
				if errors != expectedErrors {
					t.Errorf("error mismatch: have %v, want %v", err, expectedErrors)
				}
				break OUT3
			}
		case <-timeout.C:
			break OUT3
		}
	}
}

func TestPrepareExtra(t *testing.T) {
	validators := make([]common.Address, 4)
	validators[0] = common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a"))
	validators[1] = common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212"))
	validators[2] = common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6"))
	validators[3] = common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440"))

	vanity := make([]byte, types.IstanbulExtraVanity)
	expectedResult := append(vanity, hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")...)

	h := &types.Header{
		Extra: vanity,
	}

	payload, err := prepareExtra(h, validators)
	if err != nil {
		t.Errorf("error mismatch: have %v, want: nil", err)
	}
	if !reflect.DeepEqual(payload, expectedResult) {
		t.Errorf("payload mismatch: have %v, want %v", payload, expectedResult)
	}

	// append useless information to extra-data
	h.Extra = append(vanity, make([]byte, 15)...)

	payload, err = prepareExtra(h, validators)
	if !reflect.DeepEqual(payload, expectedResult) {
		t.Errorf("payload mismatch: have %v, want %v", payload, expectedResult)
	}
}

func TestWriteSeal(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)
	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
	expectedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-3)...)
	expectedIstExtra := &types.IstanbulExtra{
		Validators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
		},
		Seal:          expectedSeal,
		CommittedSeal: [][]byte{},
	}
	var expectedErr error

	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	err := writeSeal(h, expectedSeal)
	if err != expectedErr {
		t.Errorf("error mismatch: have %v, want %v", err, expectedErr)
	}

	// verify istanbul extra-data
	istExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if !reflect.DeepEqual(istExtra, expectedIstExtra) {
		t.Errorf("extra data mismatch: have %v, want %v", istExtra, expectedIstExtra)
	}

	// invalid seal
	unexpectedSeal := append(expectedSeal, make([]byte, 1)...)
	err = writeSeal(h, unexpectedSeal)
	if err != errInvalidSignature {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidSignature)
	}
}

func TestWriteCommittedSeals(t *testing.T) {
	vanity := bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity)
	istRawData := hexutil.MustDecode("0xf858f8549444add0ec310f115a0e603b2d7db9f067778eaf8a94294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212946beaaed781d2d2ab6350f5c4566a2c6eaac407a6948be76812f765c24641ec63dc2852b378aba2b44080c0")
	expectedCommittedSeal := append([]byte{1, 2, 3}, bytes.Repeat([]byte{0x00}, types.IstanbulExtraSeal-3)...)
	expectedIstExtra := &types.IstanbulExtra{
		Validators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
			common.BytesToAddress(hexutil.MustDecode("0x6beaaed781d2d2ab6350f5c4566a2c6eaac407a6")),
			common.BytesToAddress(hexutil.MustDecode("0x8be76812f765c24641ec63dc2852b378aba2b440")),
		},
		Seal:          []byte{},
		CommittedSeal: [][]byte{expectedCommittedSeal},
	}
	var expectedErr error

	h := &types.Header{
		Extra: append(vanity, istRawData...),
	}

	// normal case
	err := writeCommittedSeals(h, [][]byte{expectedCommittedSeal})
	if err != expectedErr {
		t.Errorf("error mismatch: have %v, want %v", err, expectedErr)
	}

	// verify istanbul extra-data
	istExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		t.Errorf("error mismatch: have %v, want nil", err)
	}
	if !reflect.DeepEqual(istExtra, expectedIstExtra) {
		t.Errorf("extra data mismatch: have %v, want %v", istExtra, expectedIstExtra)
	}

	// invalid seal
	unexpectedCommittedSeal := append(expectedCommittedSeal, make([]byte, 1)...)
	err = writeCommittedSeals(h, [][]byte{unexpectedCommittedSeal})
	if err != errInvalidCommittedSeals {
		t.Errorf("error mismatch: have %v, want %v", err, errInvalidCommittedSeals)
	}
}

func makeSnapshotTestConfigItems() []interface{} {
	return []interface{}{
		stakingUpdateInterval(1),
		proposerUpdateInterval(1),
		proposerPolicy(params.WeightedRandom),
	}
}

func makeFakeStakingInfo(blockNumber uint64, keys []*ecdsa.PrivateKey, amounts []uint64) *reward.StakingInfo {
	stakingInfo := &reward.StakingInfo{
		BlockNum: blockNumber,
	}
	for idx, key := range keys {
		addr := crypto.PubkeyToAddress(key.PublicKey)

		pk, _ := crypto.GenerateKey()
		rewardAddr := crypto.PubkeyToAddress(pk.PublicKey)

		stakingInfo.CouncilNodeAddrs = append(stakingInfo.CouncilNodeAddrs, addr)
		stakingInfo.CouncilStakingAmounts = append(stakingInfo.CouncilStakingAmounts, amounts[idx])
		stakingInfo.CouncilRewardAddrs = append(stakingInfo.CouncilRewardAddrs, rewardAddr)
	}
	return stakingInfo
}

func toAddressList(validators []istanbul.Validator) []common.Address {
	addresses := make([]common.Address, len(validators))
	for idx, val := range validators {
		addresses[idx] = val.Address()
	}
	return addresses
}

func copyAndSortAddrs(addrs []common.Address) []common.Address {
	copied := make([]common.Address, len(addrs))
	copy(copied, addrs)

	sort.Slice(copied, func(i, j int) bool {
		return strings.Compare(copied[i].String(), copied[j].String()) < 0
	})

	return copied
}

func makeExpectedResult(indices []int, candidate []common.Address) []common.Address {
	expected := make([]common.Address, len(indices))
	for eIdx, cIdx := range indices {
		expected[eIdx] = candidate[cIdx]
	}
	return copyAndSortAddrs(expected)
}

// Asserts taht if all (key,value) pairs of `subset` exists in `set`
func assertMapSubset(t *testing.T, subset, set map[string]interface{}) {
	for k, v := range subset {
		assert.Equal(t, set[k], v)
	}
}

func TestSnapshot_Validators_AfterMinimumStakingVotes(t *testing.T) {
	type vote struct {
		key   string
		value interface{}
	}
	type expected struct {
		blocks     []uint64
		validators []int
		demoted    []int
	}
	type testcase struct {
		stakingInfo []uint64
		votes       []vote
		expected    []expected
	}

	testcases := []testcase{
		{
			// test the validators are updated properly when minimum staking is changed in none mode
			[]uint64{8000000, 7000000, 6000000, 5000000},
			[]vote{
				{"governance.governancemode", "none"}, // voted on epoch 1, applied from 6-8
				{"reward.minimumstake", "5500000"},    // voted on epoch 2, applied from 9-11
				{"reward.minimumstake", "6500000"},    // voted on epoch 3, applied from 12-14
				{"reward.minimumstake", "7500000"},    // voted on epoch 4, applied from 15-17
				{"reward.minimumstake", "8500000"},    // voted on epoch 5, applied from 18-20
				{"reward.minimumstake", "7500000"},    // voted on epoch 6, applied from 21-23
				{"reward.minimumstake", "6500000"},    // voted on epoch 7, applied from 24-26
				{"reward.minimumstake", "5500000"},    // voted on epoch 8, applied from 27-29
				{"reward.minimumstake", "4500000"},    // voted on epoch 9, applied from 30-32
			},
			[]expected{
				{[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{9, 10, 11}, []int{0, 1, 2}, []int{3}},
				{[]uint64{12, 13, 14}, []int{0, 1}, []int{2, 3}},
				{[]uint64{15, 16, 17}, []int{0}, []int{1, 2, 3}},
				{[]uint64{18, 19, 20}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{21, 22, 23}, []int{0}, []int{1, 2, 3}},
				{[]uint64{24, 25, 26}, []int{0, 1}, []int{2, 3}},
				{[]uint64{27, 28, 29}, []int{0, 1, 2}, []int{3}},
				{[]uint64{30, 31, 32}, []int{0, 1, 2, 3}, []int{}},
			},
		},
		{
			// test the validators (including governing node) are updated properly when minimum staking is changed in single mode
			[]uint64{5000000, 6000000, 7000000, 8000000},
			[]vote{
				{"reward.minimumstake", "8500000"}, // voted on epoch 1, applied from 6-8
				{"reward.minimumstake", "7500000"}, // voted on epoch 2, applied from 9-11
				{"reward.minimumstake", "6500000"}, // voted on epoch 3, applied from 12-14
				{"reward.minimumstake", "5500000"}, // voted on epoch 4, applied from 15-17
				{"reward.minimumstake", "4500000"}, // voted on epoch 5, applied from 18-20
				{"reward.minimumstake", "5500000"}, // voted on epoch 6, applied from 21-23
				{"reward.minimumstake", "6500000"}, // voted on epoch 7, applied from 24-26
				{"reward.minimumstake", "7500000"}, // voted on epoch 8, applied from 27-29
				{"reward.minimumstake", "8500000"}, // voted on epoch 9, applied from 30-32
			},
			[]expected{
				// 0 is governing node, so it is included in the validators all the time
				{[]uint64{0, 1, 2, 3, 4, 5, 6, 7, 8}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{9, 10, 11}, []int{0, 3}, []int{1, 2}},
				{[]uint64{12, 13, 14}, []int{0, 2, 3}, []int{1}},
				{[]uint64{15, 16, 17, 18, 19, 20, 21, 22, 23}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{24, 25, 26}, []int{0, 2, 3}, []int{1}},
				{[]uint64{27, 28, 29}, []int{0, 3}, []int{1, 2}},
				{[]uint64{30, 31, 32}, []int{0, 1, 2, 3}, []int{}},
			},
		},
		{
			// test the validators are updated properly if governing node is changed
			[]uint64{6000000, 6000000, 5000000, 5000000},
			[]vote{
				{"reward.minimumstake", "5500000"}, // voted on epoch 1, applied from 6-8
				{"governance.governingnode", 2},    // voted on epoch 2, applied from 9-11
			},
			[]expected{
				// 0 is governing node, so it is included in the validators all the time
				{[]uint64{0, 1, 2, 3, 4, 5}, []int{0, 1, 2, 3}, []int{}},
				{[]uint64{6, 7, 8}, []int{0, 1}, []int{2, 3}},
				{[]uint64{9, 10, 11}, []int{0, 1, 2}, []int{3}},
			},
		},
	}

	testEpoch := 3

	var configItems []interface{}
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(testEpoch))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block

	for _, tc := range testcases {
		chain, engine := newBlockChain(4, configItems...)

		// set old staking manager after finishing this test.
		oldStakingManager := reward.GetStakingManager()

		// set new staking manager with the given staking information.
		stakingInfo := makeFakeStakingInfo(0, nodeKeys, tc.stakingInfo)
		reward.SetTestStakingManagerWithStakingInfoCache(stakingInfo)

		var (
			previousBlock, currentBlock *types.Block = nil, chain.Genesis()
		)

		for _, v := range tc.votes {
			// vote a vote in each epoch
			if v.key == "governance.governingnode" {
				idx := v.value.(int)
				v.value = addrs[idx].String()
			}
			engine.governance.AddVote(v.key, v.value)

			for i := 0; i < testEpoch; i++ {
				previousBlock = currentBlock
				currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
				_, err := chain.InsertChain(types.Blocks{currentBlock})
				assert.NoError(t, err)
			}
		}

		// insert blocks on extra epoch
		for i := 0; i < 2*testEpoch; i++ {
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err := chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)
		}

		for _, e := range tc.expected {
			for _, num := range e.blocks {
				block := chain.GetBlockByNumber(num)
				snap, err := engine.snapshot(chain, block.NumberU64(), block.Hash(), nil)
				assert.NoError(t, err)

				validators := toAddressList(snap.ValSet.List())
				demoted := toAddressList(snap.ValSet.DemotedList())

				expectedValidators := makeExpectedResult(e.validators, addrs)
				expectedDemoted := makeExpectedResult(e.demoted, addrs)

				assert.Equal(t, expectedValidators, validators)
				assert.Equal(t, expectedDemoted, demoted)
			}
		}

		reward.SetTestStakingManager(oldStakingManager)
		engine.Stop()
	}
}

func TestSnapshot_Validators_BasedOnStaking(t *testing.T) {
	type testcase struct {
		stakingAmounts       []uint64 // test staking amounts of each validator
		isIstanbulCompatible bool     // whether or not if the inserted block is istanbul compatible
		isSingleMode         bool     // whether or not if the governance mode is single
		expectedValidators   []int    // the indices of expected validators
		expectedDemoted      []int    // the indices of expected demoted validators
	}

	testcases := []testcase{
		// The following testcases are the ones before istanbul incompatible change
		{
			[]uint64{5000000, 5000000, 5000000, 5000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 5000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 6000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 6000000, 6000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{6000000, 6000000, 6000000, 6000000},
			false,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		// The following testcases are the ones after istanbul incompatible change
		{
			[]uint64{5000000, 5000000, 5000000, 5000000},
			true,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 5000000, 6000000},
			true,
			false,
			[]int{3},
			[]int{0, 1, 2},
		},
		{
			[]uint64{5000000, 5000000, 6000000, 6000000},
			true,
			false,
			[]int{2, 3},
			[]int{0, 1},
		},
		{
			[]uint64{5000000, 6000000, 6000000, 6000000},
			true,
			false,
			[]int{1, 2, 3},
			[]int{0},
		},
		{
			[]uint64{6000000, 6000000, 6000000, 6000000},
			true,
			false,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5500001, 5500000, 5499999, 0},
			true,
			false,
			[]int{0, 1},
			[]int{2, 3},
		},
		// The following testcases are the ones for testing governing node in single mode
		// The first staking amount is of the governing node
		{
			[]uint64{6000000, 6000000, 6000000, 6000000},
			true,
			true,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 6000000, 6000000, 6000000},
			true,
			true,
			[]int{0, 1, 2, 3},
			[]int{},
		},
		{
			[]uint64{5000000, 5000000, 6000000, 6000000},
			true,
			true,
			[]int{0, 2, 3},
			[]int{1},
		},
		{
			[]uint64{5000000, 5000000, 5000000, 6000000},
			true,
			true,
			[]int{0, 3},
			[]int{1, 2},
		},
		{
			[]uint64{5000000, 5000000, 5000000, 5000000},
			true,
			true,
			[]int{0, 1, 2, 3},
			[]int{},
		},
	}

	testNum := 4
	ms := uint64(5500000)
	configItems := makeSnapshotTestConfigItems()
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(ms)))
	for _, tc := range testcases {
		if tc.isIstanbulCompatible {
			configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(0)))
		}
		if tc.isSingleMode {
			configItems = append(configItems, governanceMode("single"))
		}
		chain, engine := newBlockChain(testNum, configItems...)

		// set old staking manager after finishing this test.
		oldStakingManager := reward.GetStakingManager()

		// set new staking manager with the given staking information.
		stakingInfo := makeFakeStakingInfo(0, nodeKeys, tc.stakingAmounts)
		reward.SetTestStakingManagerWithStakingInfoCache(stakingInfo)

		block := makeBlockWithSeal(chain, engine, chain.Genesis())
		_, err := chain.InsertChain(types.Blocks{block})
		assert.NoError(t, err)

		snap, err := engine.snapshot(chain, block.NumberU64(), block.Hash(), nil)
		assert.NoError(t, err)

		validators := toAddressList(snap.ValSet.List())
		demoted := toAddressList(snap.ValSet.DemotedList())

		expectedValidators := makeExpectedResult(tc.expectedValidators, addrs)
		expectedDemoted := makeExpectedResult(tc.expectedDemoted, addrs)

		assert.Equal(t, expectedValidators, validators)
		assert.Equal(t, expectedDemoted, demoted)

		reward.SetTestStakingManager(oldStakingManager)
		engine.Stop()
	}
}

func TestSnapshot_Validators_AddRemove(t *testing.T) {
	type vote struct {
		key   string
		value interface{}
	}
	type expected struct {
		validators []int // expected validator indexes at given block
	}
	type testcase struct {
		length   int // total number of blocks to simulate
		votes    map[int]vote
		expected map[int]expected
	}

	testcases := []testcase{
		{ // Singular change
			5,
			map[int]vote{
				1: {"governance.removevalidator", 3},
				3: {"governance.addvalidator", 3},
			},
			map[int]expected{
				0: {[]int{0, 1, 2, 3}},
				1: {[]int{0, 1, 2, 3}},
				2: {[]int{0, 1, 2}},
				3: {[]int{0, 1, 2}},
				4: {[]int{0, 1, 2, 3}},
			},
		},
		{ // Plural change
			5,
			map[int]vote{
				1: {"governance.removevalidator", []int{1, 2, 3}},
				3: {"governance.addvalidator", []int{1, 2}},
			},
			map[int]expected{
				0: {[]int{0, 1, 2, 3}},
				1: {[]int{0, 1, 2, 3}},
				2: {[]int{0}},
				3: {[]int{0}},
				4: {[]int{0, 1, 2}},
			},
		},
		{ // Around checkpoint interval (i.e. every 1024 block)
			checkpointInterval + 10,
			map[int]vote{
				checkpointInterval - 5: {"governance.removevalidator", 3},
				checkpointInterval - 1: {"governance.removevalidator", 2},
				checkpointInterval + 0: {"governance.removevalidator", 1},
				checkpointInterval + 1: {"governance.addvalidator", 1},
				checkpointInterval + 2: {"governance.addvalidator", 2},
				checkpointInterval + 3: {"governance.addvalidator", 3},
			},
			map[int]expected{
				0:                      {[]int{0, 1, 2, 3}},
				1:                      {[]int{0, 1, 2, 3}},
				checkpointInterval - 4: {[]int{0, 1, 2}},
				checkpointInterval + 0: {[]int{0, 1}},
				checkpointInterval + 1: {[]int{0}},
				checkpointInterval + 2: {[]int{0, 1}},
				checkpointInterval + 3: {[]int{0, 1, 2}},
				checkpointInterval + 4: {[]int{0, 1, 2, 3}},
				checkpointInterval + 9: {[]int{0, 1, 2, 3}},
			},
		},
		{ // multiple addvalidator & removevalidator
			10,
			map[int]vote{
				0: {"governance.removevalidator", 3},
				2: {"governance.addvalidator", 3},
				4: {"governance.addvalidator", 3},
				6: {"governance.removevalidator", 3},
				8: {"governance.removevalidator", 3},
			},
			map[int]expected{
				1: {[]int{0, 1, 2}},
				3: {[]int{0, 1, 2, 3}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1, 2}},
				9: {[]int{0, 1, 2}},
			},
		},
		{ // multiple removevalidator & addvalidator
			10,
			map[int]vote{
				0: {"governance.removevalidator", 3},
				2: {"governance.removevalidator", 3},
				4: {"governance.addvalidator", 3},
				6: {"governance.addvalidator", 3},
			},
			map[int]expected{
				1: {[]int{0, 1, 2}},
				3: {[]int{0, 1, 2}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1, 2, 3}},
			},
		},
		{ // multiple addvalidators & removevalidators
			10,
			map[int]vote{
				0: {"governance.removevalidator", []int{2, 3}},
				2: {"governance.addvalidator", []int{2, 3}},
				4: {"governance.addvalidator", []int{2, 3}},
				6: {"governance.removevalidator", []int{2, 3}},
				8: {"governance.removevalidator", []int{2, 3}},
			},
			map[int]expected{
				1: {[]int{0, 1}},
				3: {[]int{0, 1, 2, 3}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1}},
				9: {[]int{0, 1}},
			},
		},
		{ // multiple removevalidators & addvalidators
			10,
			map[int]vote{
				0: {"governance.removevalidator", []int{2, 3}},
				2: {"governance.removevalidator", []int{2, 3}},
				4: {"governance.addvalidator", []int{2, 3}},
				6: {"governance.addvalidator", []int{2, 3}},
			},
			map[int]expected{
				1: {[]int{0, 1}},
				3: {[]int{0, 1}},
				5: {[]int{0, 1, 2, 3}},
				7: {[]int{0, 1, 2, 3}},
			},
		},
	}

	var configItems []interface{}
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, subGroupSize(4))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	stakes := []uint64{4000000, 4000000, 4000000, 4000000}

	for _, tc := range testcases {
		// Create test blockchain
		chain, engine := newBlockChain(4, configItems...)

		oldStakingManager := reward.GetStakingManager()
		stakingInfo := makeFakeStakingInfo(0, nodeKeys, stakes)
		reward.SetTestStakingManagerWithStakingInfoCache(stakingInfo)

		// Backup the globals. The globals `nodeKeys` and `addrs` will be
		// modified according to validator change votes.
		allNodeKeys := make([]*ecdsa.PrivateKey, len(nodeKeys))
		allAddrs := make([]common.Address, len(addrs))
		copy(allNodeKeys, nodeKeys)
		copy(allAddrs, addrs)

		var previousBlock, currentBlock *types.Block = nil, chain.Genesis()

		// Create blocks with votes
		for i := 0; i < tc.length; i++ {
			if v, ok := tc.votes[i]; ok { // If a vote is scheduled in this block,
				if idx, ok := v.value.(int); ok {
					addr := allAddrs[idx]
					engine.governance.AddVote(v.key, addr)
				} else {
					addrList := makeExpectedResult(v.value.([]int), allAddrs)
					engine.governance.AddVote(v.key, addrList)
				}
				// t.Logf("Voting at block #%d for %s, %v", i, v.key, v.value)
			}

			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err := chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)

			// After a voting, reflect the validator change to the globals
			if v, ok := tc.votes[i]; ok {
				var indices []int
				if idx, ok := v.value.(int); ok {
					indices = []int{idx}
				} else {
					indices = v.value.([]int)
				}
				if v.key == "governance.addvalidator" {
					for _, i := range indices {
						includeNode(allAddrs[i], allNodeKeys[i])
					}
				}
				if v.key == "governance.removevalidator" {
					for _, i := range indices {
						excludeNodeByAddr(allAddrs[i])
					}
				}
			}
		}

		// Calculate historical validators using the snapshot.
		for i := 0; i < tc.length; i++ {
			if _, ok := tc.expected[i]; !ok {
				continue
			}
			block := chain.GetBlockByNumber(uint64(i))
			snap, err := engine.snapshot(chain, block.NumberU64(), block.Hash(), nil)
			assert.NoError(t, err)
			validators := copyAndSortAddrs(toAddressList(snap.ValSet.List()))

			expectedValidators := makeExpectedResult(tc.expected[i].validators, allAddrs)
			assert.Equal(t, expectedValidators, validators)
			// t.Logf("snap at block #%d: size %d", i, snap.ValSet.Size())
		}

		reward.SetTestStakingManager(oldStakingManager)
		engine.Stop()
	}
}

func TestGovernance_Votes(t *testing.T) {
	type vote struct {
		key   string
		value interface{}
	}
	type governanceItem struct {
		vote
		appliedBlockNumber uint64 // if applied block number is 0, then it checks the item on current block
	}
	type testcase struct {
		votes    []vote
		expected []governanceItem
	}

	testcases := []testcase{
		{
			votes: []vote{
				{"governance.governancemode", "none"},     // voted on block 1
				{"istanbul.committeesize", uint64(4)},     // voted on block 2
				{"governance.unitprice", uint64(2000000)}, // voted on block 3
				{"reward.mintingamount", "96000000000"},   // voted on block 4
				{"reward.ratio", "34/33/33"},              // voted on block 5
				{"reward.useginicoeff", true},             // voted on block 6
				{"reward.minimumstake", "5000000"},        // voted on block 7
			},
			expected: []governanceItem{
				{vote{"governance.governancemode", "none"}, 6},
				{vote{"istanbul.committeesize", uint64(4)}, 6},
				{vote{"governance.unitprice", uint64(2000000)}, 9},
				{vote{"reward.mintingamount", "96000000000"}, 9},
				{vote{"reward.ratio", "34/33/33"}, 9},
				{vote{"reward.useginicoeff", true}, 12},
				{vote{"reward.minimumstake", "5000000"}, 12},
				// check governance items on current block
				{vote{"governance.governancemode", "none"}, 0},
				{vote{"istanbul.committeesize", uint64(4)}, 0},
				{vote{"governance.unitprice", uint64(2000000)}, 0},
				{vote{"reward.mintingamount", "96000000000"}, 0},
				{vote{"reward.ratio", "34/33/33"}, 0},
				{vote{"reward.useginicoeff", true}, 0},
				{vote{"reward.minimumstake", "5000000"}, 0},
			},
		},
		{
			votes: []vote{
				{"governance.governancemode", "none"},   // voted on block 1
				{"governance.governancemode", "single"}, // voted on block 2
				{"governance.governancemode", "none"},   // voted on block 3
				{"governance.governancemode", "single"}, // voted on block 4
				{"governance.governancemode", "none"},   // voted on block 5
				{"governance.governancemode", "single"}, // voted on block 6
				{"governance.governancemode", "none"},   // voted on block 7
				{"governance.governancemode", "single"}, // voted on block 8
				{"governance.governancemode", "none"},   // voted on block 9
			},
			expected: []governanceItem{
				{vote{"governance.governancemode", "single"}, 6},
				{vote{"governance.governancemode", "none"}, 9},
				{vote{"governance.governancemode", "single"}, 12},
				{vote{"governance.governancemode", "none"}, 15},
			},
		},
		{
			votes: []vote{
				{"governance.governancemode", "none"},     // voted on block 1
				{"istanbul.committeesize", uint64(4)},     // voted on block 2
				{"governance.unitprice", uint64(2000000)}, // voted on block 3
				{"governance.governancemode", "single"},   // voted on block 4
				{"istanbul.committeesize", uint64(22)},    // voted on block 5
				{"governance.unitprice", uint64(2)},       // voted on block 6
				{"governance.governancemode", "none"},     // voted on block 7
			},
			expected: []governanceItem{
				// governance mode for all blocks
				{vote{"governance.governancemode", "single"}, 1},
				{vote{"governance.governancemode", "single"}, 2},
				{vote{"governance.governancemode", "single"}, 3},
				{vote{"governance.governancemode", "single"}, 4},
				{vote{"governance.governancemode", "single"}, 5},
				{vote{"governance.governancemode", "none"}, 6},
				{vote{"governance.governancemode", "none"}, 7},
				{vote{"governance.governancemode", "none"}, 8},
				{vote{"governance.governancemode", "single"}, 9},
				{vote{"governance.governancemode", "single"}, 10},
				{vote{"governance.governancemode", "single"}, 11},
				{vote{"governance.governancemode", "none"}, 12},
				{vote{"governance.governancemode", "none"}, 13},
				{vote{"governance.governancemode", "none"}, 14},
				{vote{"governance.governancemode", "none"}, 0}, // check on current

				// committee size for all blocks
				{vote{"istanbul.committeesize", uint64(21)}, 1},
				{vote{"istanbul.committeesize", uint64(21)}, 2},
				{vote{"istanbul.committeesize", uint64(21)}, 3},
				{vote{"istanbul.committeesize", uint64(21)}, 4},
				{vote{"istanbul.committeesize", uint64(21)}, 5},
				{vote{"istanbul.committeesize", uint64(4)}, 6},
				{vote{"istanbul.committeesize", uint64(4)}, 7},
				{vote{"istanbul.committeesize", uint64(4)}, 8},
				{vote{"istanbul.committeesize", uint64(22)}, 9},
				{vote{"istanbul.committeesize", uint64(22)}, 10},
				{vote{"istanbul.committeesize", uint64(22)}, 11},
				{vote{"istanbul.committeesize", uint64(22)}, 12},
				{vote{"istanbul.committeesize", uint64(22)}, 13},
				{vote{"istanbul.committeesize", uint64(22)}, 14},
				{vote{"istanbul.committeesize", uint64(22)}, 0}, // check on current

				// unitprice for all blocks
				{vote{"governance.unitprice", uint64(1)}, 1},
				{vote{"governance.unitprice", uint64(1)}, 2},
				{vote{"governance.unitprice", uint64(1)}, 3},
				{vote{"governance.unitprice", uint64(1)}, 4},
				{vote{"governance.unitprice", uint64(1)}, 5},
				{vote{"governance.unitprice", uint64(1)}, 6},
				{vote{"governance.unitprice", uint64(1)}, 7},
				{vote{"governance.unitprice", uint64(1)}, 8},
				{vote{"governance.unitprice", uint64(2000000)}, 9},
				{vote{"governance.unitprice", uint64(2000000)}, 10},
				{vote{"governance.unitprice", uint64(2000000)}, 11},
				{vote{"governance.unitprice", uint64(2)}, 12},
				{vote{"governance.unitprice", uint64(2)}, 13},
				{vote{"governance.unitprice", uint64(2)}, 14},
				{vote{"governance.unitprice", uint64(2)}, 0}, // check on current
			},
		},
	}

	var configItems []interface{}
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	for _, tc := range testcases {
		chain, engine := newBlockChain(1, configItems...)

		// test initial governance items
		assert.Equal(t, uint64(3), engine.governance.Epoch())
		assert.Equal(t, "single", engine.governance.GovernanceMode())
		assert.Equal(t, uint64(21), engine.governance.CommitteeSize())
		assert.Equal(t, uint64(1), engine.governance.UnitPrice())
		assert.Equal(t, "0", engine.governance.MintingAmount())
		assert.Equal(t, "100/0/0", engine.governance.Ratio())
		assert.Equal(t, false, engine.governance.UseGiniCoeff())
		assert.Equal(t, "2000000", engine.governance.MinimumStake())

		// add votes and insert voted blocks
		var (
			previousBlock, currentBlock *types.Block = nil, chain.Genesis()
			err                         error
		)

		for _, v := range tc.votes {
			engine.governance.AddVote(v.key, v.value)
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err = chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)
		}

		// insert blocks until the vote is applied
		for i := 0; i < 6; i++ {
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err = chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)
		}

		for _, item := range tc.expected {
			blockNumber := item.appliedBlockNumber
			if blockNumber == 0 {
				blockNumber = chain.CurrentBlock().NumberU64()
			}
			_, items, err := engine.governance.ReadGovernance(blockNumber)
			assert.NoError(t, err)
			assert.Equal(t, item.value, items[item.key])
		}

		engine.Stop()
	}
}

func TestGovernance_ReaderEngine(t *testing.T) {
	// Test that ReaderEngine (Params(), ParamsAt(), UpdateParams()) works.
	type vote = map[string]interface{}
	type expected = map[string]interface{} // expected (subset of) governance items
	type testcase struct {
		length   int // total number of blocks to simulate
		votes    map[int]vote
		expected map[int]expected
	}

	testcases := []testcase{
		{
			8,
			map[int]vote{
				1: {"governance.unitprice": uint64(17)},
			},
			map[int]expected{
				0: {"governance.unitprice": uint64(1)},
				1: {"governance.unitprice": uint64(1)},
				2: {"governance.unitprice": uint64(1)},
				3: {"governance.unitprice": uint64(1)},
				4: {"governance.unitprice": uint64(1)},
				5: {"governance.unitprice": uint64(1)},
				6: {"governance.unitprice": uint64(17)},
				7: {"governance.unitprice": uint64(17)},
				8: {"governance.unitprice": uint64(17)},
			},
		},
	}

	var configItems []interface{}
	configItems = append(configItems, proposerPolicy(params.WeightedRandom))
	configItems = append(configItems, proposerUpdateInterval(1))
	configItems = append(configItems, epoch(3))
	configItems = append(configItems, governanceMode("single"))
	configItems = append(configItems, minimumStake(new(big.Int).SetUint64(4000000)))
	configItems = append(configItems, istanbulCompatibleBlock(new(big.Int).SetUint64(0)))
	configItems = append(configItems, blockPeriod(0)) // set block period to 0 to prevent creating future block
	stakes := []uint64{4000000, 4000000, 4000000, 4000000}

	for _, tc := range testcases {
		// Create test blockchain
		chain, engine := newBlockChain(4, configItems...)

		oldStakingManager := reward.GetStakingManager()
		stakingInfo := makeFakeStakingInfo(0, nodeKeys, stakes)
		reward.SetTestStakingManagerWithStakingInfoCache(stakingInfo)

		var previousBlock, currentBlock *types.Block = nil, chain.Genesis()

		// Create blocks with votes
		for num := 0; num <= tc.length; num++ {
			// Validate current params with Params() and CurrentSetCopy().
			// Check that both returns the expected result.
			assertMapSubset(t, tc.expected[num], engine.governance.Params().StrMap())
			assertMapSubset(t, tc.expected[num], engine.governance.CurrentSetCopy())

			// Place a vote if a vote is scheduled in upcoming block
			// Note that we're building (head+1)'th block here.
			for k, v := range tc.votes[num+1] {
				ok := engine.governance.AddVote(k, v)
				assert.True(t, ok)
			}

			// Create a block
			previousBlock = currentBlock
			currentBlock = makeBlockWithSeal(chain, engine, previousBlock)
			_, err := chain.InsertChain(types.Blocks{currentBlock})
			assert.NoError(t, err)

			// Load parameters for the next block
			err = engine.governance.UpdateParams()
			assert.NoError(t, err)
		}

		// Validate historic parameters with ParamsAt() and ReadGovernance().
		// Check that both returns the expected result.
		for num := 0; num <= tc.length; num++ {
			pset, err := engine.governance.ParamsAt(uint64(num))
			assert.NoError(t, err)
			assertMapSubset(t, tc.expected[num], pset.StrMap())

			_, items, err := engine.governance.ReadGovernance(uint64(num))
			assert.NoError(t, err)
			assertMapSubset(t, tc.expected[num], items)
		}

		reward.SetTestStakingManager(oldStakingManager)
		engine.Stop()
	}
}
