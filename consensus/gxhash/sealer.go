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
// This file is derived from consensus/ethash/sealer.go (2018/06/04).
// Modified and improved for the klaytn development.

package gxhash

import (
	crand "crypto/rand"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"sync"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/consensus"
)

// Seal implements consensus.Engine, attempting to find a nonce that satisfies
// the block's blockscore requirements.
func (gxhash *Gxhash) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	// If we're running a fake PoW, simply return a 0 nonce immediately
	if gxhash.config.PowMode == ModeFake || gxhash.config.PowMode == ModeFullFake {
		header := block.Header()
		return block.WithSeal(header), nil
	}
	// If we're running a shared PoW, delegate sealing to it
	if gxhash.shared != nil {
		return gxhash.shared.Seal(chain, block, stop)
	}
	// Create a runner and the multiple search threads it directs
	abort := make(chan struct{})
	found := make(chan *types.Block)

	gxhash.lock.Lock()
	threads := gxhash.threads
	if gxhash.rand == nil {
		seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			gxhash.lock.Unlock()
			return nil, err
		}
		gxhash.rand = rand.New(rand.NewSource(seed.Int64()))
	}
	gxhash.lock.Unlock()
	if threads == 0 {
		threads = runtime.NumCPU()
	}
	if threads < 0 {
		threads = 0 // Allows disabling local mining without extra logic around local/remote
	}
	var pend sync.WaitGroup
	for i := 0; i < threads; i++ {
		pend.Add(1)
		go func(id int, nonce uint64) {
			defer pend.Done()
			gxhash.mine(block, id, nonce, abort, found)
		}(i, uint64(gxhash.rand.Int63()))
	}
	// Wait until sealing is terminated or a nonce is found
	var result *types.Block
	select {
	case <-stop:
		// Outside abort, stop all miner threads
		close(abort)
	case result = <-found:
		// One of the threads found a block, abort all others
		close(abort)
	case <-gxhash.update:
		// Thread count was changed on user request, restart
		close(abort)
		pend.Wait()
		return gxhash.Seal(chain, block, stop)
	}
	// Wait for all miners to terminate and return the block
	pend.Wait()
	return result, nil
}

// mine is the actual proof-of-work miner that searches for a nonce starting from
// seed that results in correct final block blockscore.
func (gxhash *Gxhash) mine(block *types.Block, id int, seed uint64, abort chan struct{}, found chan *types.Block) {
	// Extract some data from the header
	var (
		header  = block.Header()
		hash    = header.HashNoNonce().Bytes()
		target  = new(big.Int).Div(maxUint256, header.BlockScore)
		number  = header.Number.Uint64()
		dataset = gxhash.dataset(number)
	)
	// Start generating random nonces until we abort or find a good one
	var (
		attempts = int64(0)
		nonce    = seed
	)
	localLogger := logger.NewWith("miner", id)
	localLogger.Trace("Started gxhash search for new nonces", "seed", seed)
search:
	for {
		select {
		case <-abort:
			// Mining terminated, update stats and abort
			localLogger.Trace("Gxhash nonce search aborted", "attempts", nonce-seed)
			gxhash.hashrate.Mark(attempts)
			break search

		default:
			// We don't have to update hash rate on every nonce, so update after after 2^X nonces
			attempts++
			if (attempts % (1 << 15)) == 0 {
				gxhash.hashrate.Mark(attempts)
				attempts = 0
			}
			// Compute the PoW value of this nonce
			_, result := hashimotoFull(dataset.dataset, hash, nonce)
			if new(big.Int).SetBytes(result).Cmp(target) <= 0 {
				// Correct nonce found, create a new header with it
				header = types.CopyHeader(header)

				// Seal and return a block (if still needed)
				select {
				case found <- block.WithSeal(header):
					localLogger.Trace("Gxhash nonce found and reported", "attempts", nonce-seed, "nonce", nonce)
				case <-abort:
					localLogger.Trace("Gxhash nonce found but discarded", "attempts", nonce-seed, "nonce", nonce)
				}
				break search
			}
			nonce++
		}
	}
	// Datasets are unmapped in a finalizer. Ensure that the dataset stays live
	// during sealing so it's not unmapped while being read.
	runtime.KeepAlive(dataset)
}
