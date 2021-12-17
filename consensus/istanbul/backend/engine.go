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
// This file is derived from quorum/consensus/istanbul/backend/engine.go (2018/06/04).
// Modified and improved for the klaytn development.

package backend

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/klaytn/klaytn/reward"

	lru "github.com/hashicorp/golang-lru"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/consensus"
	"github.com/klaytn/klaytn/consensus/istanbul"
	istanbulCore "github.com/klaytn/klaytn/consensus/istanbul/core"
	"github.com/klaytn/klaytn/consensus/istanbul/validator"
	"github.com/klaytn/klaytn/crypto/sha3"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
)

const (
	//checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	//inmemorySnapshots  = 128  // Number of recent vote snapshots to keep in memory
	//inmemoryPeers      = 40
	//inmemoryMessages   = 1024

	checkpointInterval = 1024 // Number of blocks after which to save the vote snapshot to the database
	inmemorySnapshots  = 496  // Number of recent vote snapshots to keep in memory
	inmemoryPeers      = 200
	inmemoryMessages   = 4096

	allowedFutureBlockTime = 1 * time.Second // Max time from current time allowed for blocks, before they're considered future blocks
)

var (
	// errInvalidProposal is returned when a prposal is malformed.
	errInvalidProposal = errors.New("invalid proposal")
	// errInvalidSignature is returned when given signature is not signed by given
	// address.
	errInvalidSignature = errors.New("invalid signature")
	// errUnknownBlock is returned when the list of validators is requested for a block
	// that is not part of the local blockchain.
	errUnknownBlock = errors.New("unknown block")
	// errUnauthorized is returned if a header is signed by a non authorized entity.
	errUnauthorized = errors.New("unauthorized")
	// errInvalidBlockScore is returned if the BlockScore of a block is not 1
	errInvalidBlockScore = errors.New("invalid blockscore")
	// errInvalidExtraDataFormat is returned when the extra data format is incorrect
	errInvalidExtraDataFormat = errors.New("invalid extra data format")
	// errInvalidTimestamp is returned if the timestamp of a block is lower than the previous block's timestamp + the minimum block period.
	errInvalidTimestamp = errors.New("invalid timestamp")
	// errInvalidVotingChain is returned if an authorization list is attempted to
	// be modified via out-of-range or non-contiguous headers.
	errInvalidVotingChain = errors.New("invalid voting chain")
	// errInvalidVote is returned if a nonce value is something else that the two
	// allowed constants of 0x00..0 or 0xff..f.
	errInvalidVote = errors.New("vote nonce not 0x00..0 or 0xff..f")
	// errInvalidCommittedSeals is returned if the committed seal is not signed by any of parent validators.
	errInvalidCommittedSeals = errors.New("invalid committed seals")
	// errEmptyCommittedSeals is returned if the field of committed seals is zero.
	errEmptyCommittedSeals = errors.New("zero committed seals")
	// errMismatchTxhashes is returned if the TxHash in header is mismatch.
	errMismatchTxhashes = errors.New("mismatch transactions hashes")
)
var (
	defaultBlockScore = big.NewInt(1)
	now               = time.Now

	nonceAuthVote = hexutil.MustDecode("0xffffffffffffffff") // Magic nonce number to vote on adding a new validator
	nonceDropVote = hexutil.MustDecode("0x0000000000000000") // Magic nonce number to vote on removing a validator.

	inmemoryBlocks             = 2048 // Number of blocks to precompute validators' addresses
	inmemoryValidatorsPerBlock = 30   // Approximate number of validators' addresses from ecrecover
	signatureAddresses, _      = lru.NewARC(inmemoryBlocks * inmemoryValidatorsPerBlock)
)

// cacheSignatureAddresses extracts the address from the given data and signature and cache them for later usage.
func cacheSignatureAddresses(data []byte, sig []byte) (common.Address, error) {
	sigStr := hex.EncodeToString(sig)
	if addr, ok := signatureAddresses.Get(sigStr); ok {
		return addr.(common.Address), nil
	}
	addr, err := istanbul.GetSignatureAddress(data, sig)
	if err != nil {
		return common.Address{}, err
	}
	signatureAddresses.Add(sigStr, addr)
	return addr, err
}

// Author retrieves the Klaytn address of the account that minted the given block.
func (sb *backend) Author(header *types.Header) (common.Address, error) {
	return ecrecover(header)
}

// CanVerifyHeadersConcurrently returns true if concurrent header verification possible, otherwise returns false.
func (sb *backend) CanVerifyHeadersConcurrently() bool {
	return false
}

// PreprocessHeaderVerification prepares header verification for heavy computation before synchronous header verification such as ecrecover.
func (sb *backend) PreprocessHeaderVerification(headers []*types.Header) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, inmemoryBlocks)
	go func() {
		for _, header := range headers {
			err := sb.computeSignatureAddrs(header)

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// computeSignatureAddrs computes the addresses of signer and validators and caches them.
func (sb *backend) computeSignatureAddrs(header *types.Header) error {
	_, err := ecrecover(header)
	if err != nil {
		return err
	}

	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}

	proposalSeal := istanbulCore.PrepareCommittedSeal(header.Hash())
	for _, seal := range istanbulExtra.CommittedSeal {
		_, err := cacheSignatureAddresses(proposalSeal, seal)
		if err != nil {
			return errInvalidSignature
		}
	}
	return nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (sb *backend) VerifyHeader(chain consensus.ChainReader, header *types.Header, seal bool) error {
	return sb.verifyHeader(chain, header, nil)
}

// verifyHeader checks whether a header conforms to the consensus rules.The
// caller may optionally pass in a batch of parents (ascending order) to avoid
// looking those up from the database. This is useful for concurrently verifying
// a batch of new headers.
func (sb *backend) verifyHeader(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	if header.Number == nil {
		return errUnknownBlock
	}

	// Don't waste time checking blocks from the future
	if header.Time.Cmp(big.NewInt(now().Add(allowedFutureBlockTime).Unix())) > 0 {
		return consensus.ErrFutureBlock
	}

	// Ensure that the extra data format is satisfied
	if _, err := types.ExtractIstanbulExtra(header); err != nil {
		return errInvalidExtraDataFormat
	}
	// Ensure that the block's blockscore is meaningful (may not be correct at this point)
	if header.BlockScore == nil || header.BlockScore.Cmp(defaultBlockScore) != 0 {
		return errInvalidBlockScore
	}

	return sb.verifyCascadingFields(chain, header, parents)
}

// verifyCascadingFields verifies all the header fields that are not standalone,
// rather depend on a batch of previous headers. The caller may optionally pass
// in a batch of parents (ascending order) to avoid looking those up from the
// database. This is useful for concurrently verifying a batch of new headers.
func (sb *backend) verifyCascadingFields(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// The genesis block is the always valid dead-end
	number := header.Number.Uint64()
	if number == 0 {
		return nil
	}
	// Ensure that the block's timestamp isn't too close to it's parent
	var parent *types.Header
	if len(parents) > 0 {
		parent = parents[len(parents)-1]
	} else {
		parent = chain.GetHeader(header.ParentHash, number-1)
	}
	if parent == nil || parent.Number.Uint64() != number-1 || parent.Hash() != header.ParentHash {
		return consensus.ErrUnknownAncestor
	}
	if parent.Time.Uint64()+sb.config.BlockPeriod > header.Time.Uint64() {
		return errInvalidTimestamp
	}
	if err := sb.verifySigner(chain, header, parents); err != nil {
		return err
	}

	// At every epoch governance data will come in block header. Verify it.
	if number%sb.governance.Epoch() == 0 && len(header.Governance) > 0 {
		return sb.governance.VerifyGovernance(header.Governance)
	}
	return sb.verifyCommittedSeals(chain, header, parents)
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (sb *backend) VerifyHeaders(chain consensus.ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))
	go func() {
		for i, header := range headers {
			err := sb.verifyHeader(chain, header, headers[:i])

			select {
			case <-abort:
				return
			case results <- err:
			}
		}
	}()
	return abort, results
}

// verifySigner checks whether the signer is in parent's validator set
func (sb *backend) verifySigner(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	// Verifying the genesis block is not supported
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	// resolve the authorization key and check against signers
	signer, err := ecrecover(header)
	if err != nil {
		return err
	}

	// Signer should be in the validator set of previous block's extraData.
	if _, v := snap.ValSet.GetByAddress(signer); v == nil {
		return errUnauthorized
	}
	return nil
}

// verifyCommittedSeals checks whether every committed seal is signed by one of the parent's validators
func (sb *backend) verifyCommittedSeals(chain consensus.ChainReader, header *types.Header, parents []*types.Header) error {
	number := header.Number.Uint64()
	// We don't need to verify committed seals in the genesis block
	if number == 0 {
		return nil
	}

	// Retrieve the snapshot needed to verify this header and cache it
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, parents)
	if err != nil {
		return err
	}

	extra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return err
	}
	// The length of Committed seals should be larger than 0
	if len(extra.CommittedSeal) == 0 {
		return errEmptyCommittedSeals
	}

	validators := snap.ValSet.Copy()
	// Check whether the committed seals are generated by parent's validators
	validSeal := 0
	proposalSeal := istanbulCore.PrepareCommittedSeal(header.Hash())
	// 1. Get committed seals from current header
	for _, seal := range extra.CommittedSeal {
		// 2. Get the original address by seal and parent block hash
		addr, err := cacheSignatureAddresses(proposalSeal, seal)
		if err != nil {
			return errInvalidSignature
		}
		// Every validator can have only one seal. If more than one seals are signed by a
		// validator, the validator cannot be found and errInvalidCommittedSeals is returned.
		if validators.RemoveValidator(addr) {
			validSeal += 1
		} else {
			return errInvalidCommittedSeals
		}
	}

	// The length of validSeal should be larger than number of faulty node + 1
	if validSeal <= 2*snap.ValSet.F() {
		return errInvalidCommittedSeals
	}

	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (sb *backend) VerifySeal(chain consensus.ChainReader, header *types.Header) error {
	// get parent header and ensure the signer is in parent's validator set
	number := header.Number.Uint64()
	if number == 0 {
		return errUnknownBlock
	}

	// ensure that the blockscore equals to defaultBlockScore
	if header.BlockScore.Cmp(defaultBlockScore) != 0 {
		return errInvalidBlockScore
	}
	return sb.verifySigner(chain, header, nil)
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (sb *backend) Prepare(chain consensus.ChainReader, header *types.Header) error {
	// unused fields, force to set to empty
	header.Rewardbase = sb.rewardbase

	// copy the parent extra data as the header extra data
	number := header.Number.Uint64()
	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	// use the same blockscore for all blocks
	header.BlockScore = defaultBlockScore

	// Assemble the voting snapshot
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return err
	}

	// If it reaches the Epoch, governance config will be added to block header
	if number%sb.governance.Epoch() == 0 {
		if g := sb.governance.GetGovernanceChange(); g != nil {
			if data, err := json.Marshal(g); err != nil {
				logger.Error("Failed to encode governance changes!! Possible configuration mismatch!! ")
			} else {
				if header.Governance, err = rlp.EncodeToBytes(data); err != nil {
					logger.Error("Failed to encode governance data for the header", "num", number)
				}
			}
		}
	}

	// if there is a vote to attach, attach it to the header
	header.Vote = sb.governance.GetEncodedVote(sb.address, number)

	// add validators (council list) in snapshot to extraData's validators section
	extra, err := prepareExtra(header, snap.validators())
	if err != nil {
		return err
	}
	header.Extra = extra

	// set header's timestamp
	header.Time = new(big.Int).Add(parent.Time, new(big.Int).SetUint64(sb.config.BlockPeriod))
	header.TimeFoS = parent.TimeFoS
	if header.Time.Int64() < time.Now().Unix() {
		t := time.Now()
		header.Time = big.NewInt(t.Unix())
		header.TimeFoS = uint8((t.UnixNano() / 1000 / 1000 / 10) % 100)
	}
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
//
// Note, the block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (sb *backend) Finalize(chain consensus.ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	receipts []*types.Receipt) (*types.Block, error) {

	// If sb.chain is nil, it means backend is not initialized yet.
	if sb.chain != nil && sb.governance.ProposerPolicy() == uint64(istanbul.WeightedRandom) {
		// TODO-Klaytn Let's redesign below logic and remove dependency between block reward and istanbul consensus.

		pocAddr := common.Address{}
		kirAddr := common.Address{}
		lastHeader := chain.CurrentHeader()
		valSet := sb.getValidators(lastHeader.Number.Uint64(), lastHeader.Hash())

		// Determine and update Rewardbase when mining. When mining, state root is not yet determined and will be determined at the end of this Finalize below.
		if common.EmptyHash(header.Root) {
			var logMsg string
			_, nodeValidator := valSet.GetByAddress(sb.address)
			if nodeValidator == nil || (nodeValidator.RewardAddress() == common.Address{}) {
				logMsg = "No reward address for nodeValidator. Use node's rewardbase."
			} else {
				// use reward address of current node.
				// only a block made by proposer will be accepted. However, due to round change any node can be the proposer of a block.
				// so need to write reward address of current node to receive reward when it becomes proposer.
				// if current node does not become proposer, the block will be abandoned
				header.Rewardbase = nodeValidator.RewardAddress()
				logMsg = "Use reward address for nodeValidator."
			}
			logger.Trace(logMsg, "header.Number", header.Number.Uint64(), "node address", sb.address, "rewardbase", header.Rewardbase)
		}

		if stakingInfo := reward.GetStakingInfo(header.Number.Uint64()); stakingInfo != nil {
			kirAddr = stakingInfo.KIRAddr
			pocAddr = stakingInfo.PoCAddr
		}

		if err := sb.rewardDistributor.DistributeBlockReward(state, header, pocAddr, kirAddr); err != nil {
			return nil, err
		}
	} else {
		if err := sb.rewardDistributor.MintKLAY(state, header); err != nil {
			return nil, err
		}
	}

	header.Root = state.IntermediateRoot(true)

	// Assemble and return the final block for sealing
	return types.NewBlock(header, txs, receipts), nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (sb *backend) Seal(chain consensus.ChainReader, block *types.Block, stop <-chan struct{}) (*types.Block, error) {
	// update the block header timestamp and signature and propose the block to core engine
	header := block.Header()
	number := header.Number.Uint64()

	// Bail out if we're unauthorized to sign a block
	snap, err := sb.snapshot(chain, number-1, header.ParentHash, nil)
	if err != nil {
		return nil, err
	}
	if _, v := snap.ValSet.GetByAddress(sb.address); v == nil {
		return nil, errUnauthorized
	}

	parent := chain.GetHeader(header.ParentHash, number-1)
	if parent == nil {
		return nil, consensus.ErrUnknownAncestor
	}
	block, err = sb.updateBlock(parent, block)
	if err != nil {
		return nil, err
	}

	// wait for the timestamp of header, use this to adjust the block period
	delay := time.Unix(block.Header().Time.Int64(), 0).Sub(now())
	select {
	case <-time.After(delay):
	case <-stop:
		return nil, nil
	}

	// get the proposed block hash and clear it if the seal() is completed.
	sb.sealMu.Lock()
	sb.proposedBlockHash = block.Hash()
	clear := func() {
		sb.proposedBlockHash = common.Hash{}
		sb.sealMu.Unlock()
	}
	defer clear()

	// post block into Istanbul engine
	go sb.EventMux().Post(istanbul.RequestEvent{
		Proposal: block,
	})

	for {
		select {
		case result := <-sb.commitCh:
			if result == nil {
				return nil, nil
			}
			// if the block hash and the hash from channel are the same,
			// return the result. Otherwise, keep waiting the next hash.
			block = types.SetRoundToBlock(block, result.Round)
			if block.Hash() == result.Block.Hash() {
				return result.Block, nil
			}
		case <-stop:
			return nil, nil
		}
	}
}

// update timestamp and signature of the block based on its number of transactions
func (sb *backend) updateBlock(parent *types.Header, block *types.Block) (*types.Block, error) {
	header := block.Header()
	// sign the hash
	seal, err := sb.Sign(sigHash(header).Bytes())
	if err != nil {
		return nil, err
	}

	err = writeSeal(header, seal)
	if err != nil {
		return nil, err
	}

	return block.WithSeal(header), nil
}

func (sb *backend) CalcBlockScore(chain consensus.ChainReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(0)
}

// APIs returns the RPC APIs this consensus engine provides.
func (sb *backend) APIs(chain consensus.ChainReader) []rpc.API {
	return []rpc.API{
		{
			Namespace: "istanbul",
			Version:   "1.0",
			Service:   &API{chain: chain, istanbul: sb},
			Public:    true,
		}, {
			Namespace: "klay",
			Version:   "1.0",
			Service:   &APIExtension{chain: chain, istanbul: sb},
			Public:    true,
		},
	}
}

// SetChain sets chain of the Istanbul backend
func (sb *backend) SetChain(chain consensus.ChainReader) {
	sb.chain = chain
}

// Start implements consensus.Istanbul.Start
func (sb *backend) Start(chain consensus.ChainReader, currentBlock func() *types.Block, hasBadBlock func(hash common.Hash) bool) error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if sb.coreStarted {
		return istanbul.ErrStartedEngine
	}

	// clear previous data
	sb.proposedBlockHash = common.Hash{}
	if sb.commitCh != nil {
		close(sb.commitCh)
	}
	sb.commitCh = make(chan *types.Result, 1)

	sb.SetChain(chain)
	sb.currentBlock = currentBlock
	sb.hasBadBlock = hasBadBlock

	if err := sb.core.Start(); err != nil {
		return err
	}

	sb.coreStarted = true
	return nil
}

// Stop implements consensus.Istanbul.Stop
func (sb *backend) Stop() error {
	sb.coreMu.Lock()
	defer sb.coreMu.Unlock()
	if !sb.coreStarted {
		return istanbul.ErrStoppedEngine
	}
	if err := sb.core.Stop(); err != nil {
		return err
	}
	sb.coreStarted = false
	return nil
}

// initSnapshot initializes and stores a new Snapshot.
func (sb *backend) initSnapshot(chain consensus.ChainReader) (*Snapshot, error) {
	genesis := chain.GetHeaderByNumber(0)
	if err := sb.VerifyHeader(chain, genesis, false); err != nil {
		return nil, err
	}
	istanbulExtra, err := types.ExtractIstanbulExtra(genesis)
	if err != nil {
		return nil, err
	}

	proposerPolicy, ok := sb.governance.GetGovernanceValue(params.Policy).(uint64)
	if !ok {
		proposerPolicy = params.DefaultProposerPolicy
	}
	committeeSize, ok := sb.governance.GetGovernanceValue(params.CommitteeSize).(uint64)
	if !ok {
		committeeSize = params.DefaultSubGroupSize
	}
	snap := newSnapshot(sb.governance, 0, genesis.Hash(), validator.NewValidatorSet(istanbulExtra.Validators, nil, istanbul.ProposerPolicy(proposerPolicy), committeeSize, chain), chain.Config())

	if err := snap.store(sb.db); err != nil {
		return nil, err
	}
	logger.Trace("Stored genesis voting snapshot to disk")
	return snap, nil
}

// getPrevHeaderAndUpdateParents returns previous header to find stored Snapshot object and drops the last element of the parents parameter.
func getPrevHeaderAndUpdateParents(chain consensus.ChainReader, number uint64, hash common.Hash, parents *[]*types.Header) *types.Header {
	var header *types.Header
	if len(*parents) > 0 {
		// If we have explicit parents, pick from there (enforced)
		header = (*parents)[len(*parents)-1]
		if header.Hash() != hash || header.Number.Uint64() != number {
			return nil
		}
		*parents = (*parents)[:len(*parents)-1]
	} else {
		// No explicit parents (or no more left), reach out to the database
		header = chain.GetHeader(hash, number)
		if header == nil {
			return nil
		}
	}
	return header
}

// CreateSnapshot does not return a snapshot but creates a new snapshot at a given point in time.
func (sb *backend) CreateSnapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) error {
	_, err := sb.snapshot(chain, number, hash, parents)
	return err
}

// snapshot retrieves the authorization snapshot at a given point in time.
func (sb *backend) snapshot(chain consensus.ChainReader, number uint64, hash common.Hash, parents []*types.Header) (*Snapshot, error) {
	// Search for a snapshot in memory or on disk for checkpoints
	var (
		headers []*types.Header
		snap    *Snapshot
	)

	for snap == nil {
		// If an in-memory snapshot was found, use that
		if s, ok := sb.recents.Get(hash); ok {
			snap = s.(*Snapshot)
			break
		}
		// If an on-disk checkpoint snapshot can be found, use that
		if number%checkpointInterval == 0 {
			if s, err := loadSnapshot(sb.db, hash); err == nil {
				logger.Trace("Loaded voting snapshot form disk", "number", number, "hash", hash)
				snap = s
				break
			}
		}
		// If we're at block zero, make a snapshot
		if number == 0 {
			var err error
			if snap, err = sb.initSnapshot(chain); err != nil {
				return nil, err
			}
			break
		}
		// No snapshot for this header, gather the header and move backward
		if header := getPrevHeaderAndUpdateParents(chain, number, hash, &parents); header == nil {
			return nil, consensus.ErrUnknownAncestor
		} else {
			headers = append(headers, header)
			number, hash = number-1, header.ParentHash
		}
	}
	// Previous snapshot found, apply any pending headers on top of it
	for i := 0; i < len(headers)/2; i++ {
		headers[i], headers[len(headers)-1-i] = headers[len(headers)-1-i], headers[i]
	}
	snap, err := snap.apply(headers, sb.governance, sb.address, sb.governance.ProposerPolicy(), chain)
	if err != nil {
		return nil, err
	}

	// If we've generated a new checkpoint snapshot, save to disk
	if snap.Number%checkpointInterval == 0 && len(headers) > 0 {
		if sb.governance.CanWriteGovernanceState(snap.Number) {
			sb.governance.WriteGovernanceState(snap.Number, true)
		}
		if err = snap.store(sb.db); err != nil {
			return nil, err
		}
		logger.Trace("Stored voting snapshot to disk", "number", snap.Number, "hash", snap.Hash)
	}

	sb.recents.Add(snap.Hash, snap)
	return snap, err
}

// FIXME: Need to update this for Istanbul
// sigHash returns the hash which is used as input for the Istanbul
// signing. It is the hash of the entire header apart from the 65 byte signature
// contained at the end of the extra data.
//
// Note, the method requires the extra data to be at least 65 bytes, otherwise it
// panics. This is done to avoid accidentally using both forms (signature present
// or not), which could be abused to produce different hashes for the same header.
func sigHash(header *types.Header) (hash common.Hash) {
	hasher := sha3.NewKeccak256()

	// Clean seal is required for calculating proposer seal.
	rlp.Encode(hasher, types.IstanbulFilteredHeader(header, false))
	hasher.Sum(hash[:0])
	return hash
}

// ecrecover extracts the Klaytn account address from a signed header.
func ecrecover(header *types.Header) (common.Address, error) {
	// Retrieve the signature from the header extra-data
	istanbulExtra, err := types.ExtractIstanbulExtra(header)
	if err != nil {
		return common.Address{}, err
	}
	addr, err := cacheSignatureAddresses(sigHash(header).Bytes(), istanbulExtra.Seal)
	if err != nil {
		return addr, err
	}

	return addr, nil
}

// prepareExtra returns a extra-data of the given header and validators
func prepareExtra(header *types.Header, vals []common.Address) ([]byte, error) {
	var buf bytes.Buffer

	// compensate the lack bytes if header.Extra is not enough IstanbulExtraVanity bytes.
	if len(header.Extra) < types.IstanbulExtraVanity {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, types.IstanbulExtraVanity-len(header.Extra))...)
	}
	buf.Write(header.Extra[:types.IstanbulExtraVanity])

	ist := &types.IstanbulExtra{
		Validators:    vals,
		Seal:          []byte{},
		CommittedSeal: [][]byte{},
	}

	payload, err := rlp.EncodeToBytes(&ist)
	if err != nil {
		return nil, err
	}

	return append(buf.Bytes(), payload...), nil
}

// writeSeal writes the extra-data field of the given header with the given seals.
// suggest to rename to writeSeal.
func writeSeal(h *types.Header, seal []byte) error {
	if len(seal)%types.IstanbulExtraSeal != 0 {
		return errInvalidSignature
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.Seal = seal
	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}

// writeCommittedSeals writes the extra-data field of a block header with given committed seals.
func writeCommittedSeals(h *types.Header, committedSeals [][]byte) error {
	if len(committedSeals) == 0 {
		return errInvalidCommittedSeals
	}

	for _, seal := range committedSeals {
		if len(seal) != types.IstanbulExtraSeal {
			return errInvalidCommittedSeals
		}
	}

	istanbulExtra, err := types.ExtractIstanbulExtra(h)
	if err != nil {
		return err
	}

	istanbulExtra.CommittedSeal = make([][]byte, len(committedSeals))
	copy(istanbulExtra.CommittedSeal, committedSeals)

	payload, err := rlp.EncodeToBytes(&istanbulExtra)
	if err != nil {
		return err
	}

	h.Extra = append(h.Extra[:types.IstanbulExtraVanity], payload...)
	return nil
}
