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
// This file is derived from consensus/ethash/consensus.go (2018/06/04).
// Modified and improved for the klaytn development.

package gxhash

import (
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"time"

	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/params"
)

var (
	ByzantiumBlockReward   *big.Int = big.NewInt(3e+18) // Block reward in peb for successfully mining a block upward from Byzantium
	allowedFutureBlockTime          = 15 * time.Second  // Max time from current time allowed for blocks, before they're considered future blocks
)

// Various error messages to mark blocks invalid. These should be private to
// prevent engine specific errors from being referenced in the remainder of the
// codebase, inherently breaking if the engine is swapped out. Please put common
// error types into the consensus package.
var (
	errLargeBlockTime    = errors.New("timestamp too big")
	errZeroBlockTime     = errors.New("timestamp equals parent's")
	errInvalidBlockScore = errors.New("non-positive blockScore")
	errInvalidPoW        = errors.New("invalid proof-of-work")
)

// Author implements consensus.Engine, returning the header's coinbase as the
// proof-of-work verified author of the block.
func (gxhash *Gxhash) Author(header *types.Header) (common.Address, error) {
	// Returns arbitrary address because gxhash is used just for testing
	return params.AuthorAddressForTesting, nil
}

// CanVerifyHeadersConcurrently returns true if concurrent header verification possible, otherwise returns false.
func (gxhash *Gxhash) CanVerifyHeadersConcurrently() bool {
	return true
}

// PreprocessHeaderVerification prepares header verification for heavy computation before synchronous header verification such as ecrecover.
func (gxhash *Gxhash) PreprocessHeaderVerification(headers []*types.Header) (chan<- struct{}, <-chan error) {
	panic("this method is not used for PoW engine")
}

// CreateSnapshot is not used for PoW engine.
func (gxhash *Gxhash) CreateSnapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) error {
	return nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of the
// stock Klaytn gxhash engine.
func (gxhash *Gxhash) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	// If we're running a full engine faking, accept any input as valid
	if gxhash.config.PowMode == ModeFullFake {
		return nil
	}
	// Short circuit if the header is known, or it's parent not
	number := header.Number.Uint64()
	if chain.GetHeader(header.Hash(), number) != nil {
		return nil
	}
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// Sanity checks passed, do a proper verification
	return gxhash.verifyHeader(chain, header, parent, seal)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications.
func (gxhash *Gxhash) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	// If we're running a full engine faking, accept any input as valid
	if gxhash.config.PowMode == ModeFullFake || len(headers) == 0 {
		abort, results := make(chan struct{}), make(chan error, len(headers))
		for i := 0; i < len(headers); i++ {
			results <- nil
		}
		return abort, results
	}

	// Spawn as many workers as allowed threads
	workers := runtime.GOMAXPROCS(0)
	if len(headers) < workers {
		workers = len(headers)
	}

	// Create a task channel and spawn the verifiers
	var (
		inputs = make(chan int)
		done   = make(chan int, workers)
		errors = make([]error, len(headers))
		abort  = make(chan struct{})
	)
	for i := 0; i < workers; i++ {
		go func() {
			for index := range inputs {
				errors[index] = gxhash.verifyHeaderWorker(chain, headers, seals, index)
				done <- index
			}
		}()
	}

	errorsOut := make(chan error, len(headers))
	go func() {
		defer close(inputs)
		var (
			in, out = 0, 0
			checked = make([]bool, len(headers))
			inputs  = inputs
		)
		for {
			select {
			case inputs <- in:
				if in++; in == len(headers) {
					// Reached end of headers. Stop sending to workers.
					inputs = nil
				}
			case index := <-done:
				for checked[index] = true; checked[out]; out++ {
					errorsOut <- errors[out]
					if out == len(headers)-1 {
						return
					}
				}
			case <-abort:
				return
			}
		}
	}()
	return abort, errorsOut
}

func (gxhash *Gxhash) verifyHeaderWorker(chain consensus.ChainReader, headers []*types.Header, seals []bool, index int) error {
	var parent *types.Header
	if index == 0 {
		parent = chain.GetHeader(headers[0].ParentHash, headers[0].Number.Uint64()-1)
	} else if headers[index-1].Hash() == headers[index].ParentHash {
		parent = headers[index-1]
	}
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	if chain.GetHeader(headers[index].Hash(), headers[index].Number.Uint64()) != nil {
		return nil // known block
	}
	return gxhash.verifyHeader(chain, headers[index], parent, seals[index])
}

// verifyHeader checks whether a header conforms to the consensus rules of the
// stock Klaytn gxhash engine.
// See YP section 4.3.4. "Block Header Validity"
func (gxhash *Gxhash) verifyHeader(chain consensus.ChainReader, header, parent *types.Header, seal bool) error {
	// Ensure that the header's extra-data section is of a reasonable size
	if uint64(len(header.Extra)) > params.MaximumExtraDataSize {
		return fmt.Errorf("extra-data too long: %d > %d", len(header.Extra), params.MaximumExtraDataSize)
	}
	// Verify the header's timestamp
	if header.Time.Cmp(big.NewInt(time.Now().Add(allowedFutureBlockTime).Unix())) > 0 {
		return consensus.ErrFutureBlock
	}
	if header.Time.Cmp(parent.Time) <= 0 {
		return errZeroBlockTime
	}
	// Verify the block's blockscore based in it's timestamp and parent's blockscore
	expected := gxhash.CalcBlockScore(chain, header.Time.Uint64(), parent)

	if expected.Cmp(header.BlockScore) != 0 {
		return fmt.Errorf("invalid blockscore: have %v, want %v", header.BlockScore, expected)
	}
	// Verify that the block number is parent's +1
	if diff := new(big.Int).Sub(header.Number, parent.Number); diff.Cmp(common.Big1) != 0 {
		return consensus.ErrInvalidNumber
	}
	// Verify the engine specific seal securing the block
	if seal {
		if err := gxhash.VerifySeal(chain, header); err != nil {
			return err
		}
	}
	return nil
}

// CalcBlockScore is the blockscore adjustment algorithm. It returns
// the blockscore that a new block should have when created at time
// given the parent block's time and blockscore.
func (gxhash *Gxhash) CalcBlockScore(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return CalcBlockScore(chain.Config(), time, parent)
}

// CalcBlockScore is the blockscore adjustment algorithm. It returns
// the blockscore that a new block should have when created at time
// given the parent block's time and blockscore.
func CalcBlockScore(config *params.ChainConfig, time uint64, parent *types.Header) *big.Int {
	return calcBlockScoreByzantium(time, parent)
}

// Some weird constants to avoid constant memory allocs for them.
var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	big10         = big.NewInt(10)
	bigMinus99    = big.NewInt(-99)
	big2999999    = big.NewInt(2999999)
)

// calcBlockScoreByzantium is the blockscore adjustment algorithm. It returns
// the blockscore that a new block should have when created at time given the
// parent block's time and blockscore. The calculation uses the Byzantium rules.
func calcBlockScoreByzantium(time uint64, parent *types.Header) *big.Int {
	// https://github.com/ethereum/EIPs/issues/100.
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max((1) - ((timestamp - parent.timestamp) // 9), -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// (1) - (block_timestamp - parent_timestamp) // 9
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big9)
	x.Sub(big1, x)

	// max((1) - (block_timestamp - parent_timestamp) // 9, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// parent_diff + (parent_diff / 2048 * max((1) - ((timestamp - parent.timestamp) // 9), -99))
	y.Div(parent.BlockScore, params.BlockScoreBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.BlockScore, x)

	// minimum blockscore can ever be (before exponential factor)
	if x.Cmp(params.MinimumBlockScore) < 0 {
		x.Set(params.MinimumBlockScore)
	}
	// calculate a fake block number for the ice-age delay:
	//   https://github.com/ethereum/EIPs/pull/669
	//   fake_block_number = min(0, block.number - 3_000_000
	fakeBlockNumber := new(big.Int)
	if parent.Number.Cmp(big2999999) >= 0 {
		fakeBlockNumber = fakeBlockNumber.Sub(parent.Number, big2999999) // Note, parent is 1 less than the actual block number
	}
	// for the exponential factor
	periodCount := fakeBlockNumber
	periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

// calcBlockScoreHomestead is the blockscore adjustment algorithm. It returns
// the blockscore that a new block should have when created at time given the
// parent block's time and blockscore. The calculation uses the Homestead rules.
func calcBlockScoreHomestead(time uint64, parent *types.Header) *big.Int {
	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-2.md
	// algorithm:
	// diff = (parent_diff +
	//         (parent_diff / 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	//        ) + 2^(periodCount - 2)

	bigTime := new(big.Int).SetUint64(time)
	bigParentTime := new(big.Int).Set(parent.Time)

	// holds intermediate values to make the algo easier to read & audit
	x := new(big.Int)
	y := new(big.Int)

	// 1 - (block_timestamp - parent_timestamp) // 10
	x.Sub(bigTime, bigParentTime)
	x.Div(x, big10)
	x.Sub(big1, x)

	// max(1 - (block_timestamp - parent_timestamp) // 10, -99)
	if x.Cmp(bigMinus99) < 0 {
		x.Set(bigMinus99)
	}
	// (parent_diff + parent_diff // 2048 * max(1 - (block_timestamp - parent_timestamp) // 10, -99))
	y.Div(parent.BlockScore, params.BlockScoreBoundDivisor)
	x.Mul(y, x)
	x.Add(parent.BlockScore, x)

	// minimum blockscore can ever be (before exponential factor)
	if x.Cmp(params.MinimumBlockScore) < 0 {
		x.Set(params.MinimumBlockScore)
	}
	// for the exponential factor
	periodCount := new(big.Int).Add(parent.Number, big1)
	periodCount.Div(periodCount, expDiffPeriod)

	// the exponential factor, commonly referred to as "the bomb"
	// diff = diff + 2^(periodCount - 2)
	if periodCount.Cmp(big1) > 0 {
		y.Sub(periodCount, big2)
		y.Exp(big2, y, nil)
		x.Add(x, y)
	}
	return x
}

// VerifySeal implements consensus.Engine, checking whether the given block satisfies
// the PoW blockscore requirements.
func (gxhash *Gxhash) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// If we're running a fake PoW, accept any seal as valid
	if gxhash.config.PowMode == ModeFake || gxhash.config.PowMode == ModeFullFake {
		time.Sleep(gxhash.fakeDelay)
		if gxhash.fakeFail == header.Number.Uint64() {
			return errInvalidPoW
		}
		return nil
	}
	// If we're running a shared PoW, delegate verification to it
	if gxhash.shared != nil {
		return gxhash.shared.VerifySeal(chain, header)
	}
	// Ensure that we have a valid blockscore for the block
	if header.BlockScore.Sign() <= 0 {
		return errInvalidBlockScore
	}
	// Recompute the digest and PoW value and verify against the header
	number := header.Number.Uint64()

	cache := gxhash.cache(number)
	//size := datasetSize(number)
	//if gxhash.config.PowMode == ModeTest {
	//	size = 32 * 1024
	//}
	//digest, result := hashimotoLight(size, cache.cache, header.HashNoNonce().Bytes(), 0)
	// Caches are unmapped in a finalizer. Ensure that the cache stays live
	// until after the call to hashimotoLight so it's not unmapped while being used.
	runtime.KeepAlive(cache)

	//target := new(big.Int).Div(maxUint256, header.blockscore)
	//if new(big.Int).SetBytes(result).Cmp(target) > 0 {
	//	return errInvalidPoW
	//}
	return nil
}

// Prepare implements consensus.Engine, initializing the blockscore field of a
// header to conform to the gxhash protocol. The changes are done inline.
func (gxhash *Gxhash) Prepare(chain consensus.ChainReader, header *types.Header) error {
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	header.BlockScore = gxhash.CalcBlockScore(chain, header.Time.Uint64(), parent)
	return nil
}

// Finalize implements consensus.Engine, accumulating the block rewards,
// setting the final state and assembling the block.
func (gxhash *Gxhash) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, receipts []*types.Receipt) (*types.Block, error) {
	// Accumulate any block rewards and commit the final state root
	accumulateRewards(chain.Config(), state, header)
	header.Root = state.IntermediateRoot(true)

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, receipts), nil
}

// Some weird constants to avoid constant memory allocs for them.
var (
	big8  = big.NewInt(8)
	big32 = big.NewInt(32)
)

// AccumulateRewards credits the coinbase of the given block with the mining
// reward. The total reward consists of the static block reward.
func accumulateRewards(config *params.ChainConfig, state *state.StateDB, header *types.Header) {
	// Select the correct block reward based on chain progression
	blockReward := ByzantiumBlockReward

	// Accumulate the rewards for the miner
	reward := new(big.Int).Set(blockReward)

	state.AddBalance(params.AuthorAddressForTesting, reward)
}
