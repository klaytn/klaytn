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
// This file is derived from quorum/consensus/istanbul/backend/snapshot.go (2018/06/04).
// Modified and improved for the klaytn development.

package backend

import (
	"bytes"
	"encoding/json"

	"github.com/klaytn/klaytn/consensus"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/consensus/istanbul/validator"
	"github.com/klaytn/klaytn/governance"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

const (
	dbKeySnapshotPrefix = "istanbul-snapshot"
)

// Snapshot is the state of the authorization voting at a given point in time.
type Snapshot struct {
	Epoch         uint64                // The number of blocks after which to checkpoint and reset the pending votes
	Number        uint64                // Block number where the snapshot was created
	Hash          common.Hash           // Block hash where the snapshot was created
	ValSet        istanbul.ValidatorSet // Set of authorized validators at this moment
	Policy        uint64
	CommitteeSize uint64
	Votes         []governance.GovernanceVote      // List of votes cast in chronological order
	Tally         []governance.GovernanceTallyItem // Current vote tally to avoid recalculating
}

func getGovernanceValue(gov governance.Engine, number uint64) (epoch uint64, policy uint64, committeeSize uint64) {
	if r, err := gov.GetGovernanceItemAtNumber(number, governance.GovernanceKeyMapReverse[params.Epoch]); err == nil && r != nil {
		epoch = r.(uint64)
	} else {
		logger.Error("Couldn't get governance value istanbul.epoch", "err", err, "received", r)
		epoch = params.DefaultEpoch
	}

	if r, err := gov.GetGovernanceItemAtNumber(number, governance.GovernanceKeyMapReverse[params.Policy]); err == nil && r != nil {
		policy = r.(uint64)
	} else {
		logger.Error("Couldn't get governance value istanbul.policy", "err", err, "received", r)
		policy = params.DefaultProposerPolicy
	}

	if r, err := gov.GetGovernanceItemAtNumber(number, governance.GovernanceKeyMapReverse[params.CommitteeSize]); err == nil && r != nil {
		committeeSize = r.(uint64)
	} else {
		logger.Error("Couldn't get governance value istanbul.committeesize", "err", err, "received", r)
		committeeSize = params.DefaultSubGroupSize
	}
	return
}

// newSnapshot create a new snapshot with the specified startup parameters. This
// method does not initialize the set of recent validators, so only ever use if for
// the genesis block.
func newSnapshot(gov governance.Engine, number uint64, hash common.Hash, valSet istanbul.ValidatorSet, chainConfig *params.ChainConfig) *Snapshot {
	epoch, policy, committeeSize := getGovernanceValue(gov, number)

	snap := &Snapshot{
		Epoch:         epoch,
		Number:        number,
		Hash:          hash,
		ValSet:        valSet,
		Policy:        policy,
		CommitteeSize: committeeSize,
		Votes:         make([]governance.GovernanceVote, 0),
		Tally:         make([]governance.GovernanceTallyItem, 0),
	}
	return snap
}

// loadSnapshot loads an existing snapshot from the database.
func loadSnapshot(db database.DBManager, hash common.Hash) (*Snapshot, error) {
	blob, err := db.ReadIstanbulSnapshot(hash)
	if err != nil {
		return nil, err
	}
	snap := new(Snapshot)
	if err := json.Unmarshal(blob, snap); err != nil {
		return nil, err
	}
	return snap, nil
}

// store inserts the snapshot into the database.
func (s *Snapshot) store(db database.DBManager) error {
	blob, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return db.WriteIstanbulSnapshot(s.Hash, blob)
}

// copy creates a deep copy of the snapshot, though not the individual votes.
func (s *Snapshot) copy() *Snapshot {
	cpy := &Snapshot{
		Epoch:         s.Epoch,
		Number:        s.Number,
		Hash:          s.Hash,
		ValSet:        s.ValSet.Copy(),
		Policy:        s.Policy,
		CommitteeSize: s.CommitteeSize,
		Votes:         make([]governance.GovernanceVote, len(s.Votes)),
		Tally:         make([]governance.GovernanceTallyItem, len(s.Tally)),
	}

	copy(cpy.Votes, s.Votes)
	copy(cpy.Tally, s.Tally)

	return cpy
}

// checkVote return whether it's a valid vote
func (s *Snapshot) checkVote(address common.Address, authorize bool) bool {
	_, validator := s.ValSet.GetByAddress(address)
	return (validator != nil && !authorize) || (validator == nil && authorize)
}

// apply creates a new authorization snapshot by applying the given headers to
// the original one.
func (s *Snapshot) apply(headers []*types.Header, gov governance.Engine, addr common.Address, policy uint64, chain consensus.ChainReader) (*Snapshot, error) {
	// Allow passing in no headers for cleaner code
	if len(headers) == 0 {
		return s, nil
	}
	// Sanity check that the headers can be applied
	for i := 0; i < len(headers)-1; i++ {
		if headers[i+1].Number.Uint64() != headers[i].Number.Uint64()+1 {
			return nil, errInvalidVotingChain
		}
	}
	if headers[0].Number.Uint64() != s.Number+1 {
		return nil, errInvalidVotingChain
	}

	// Iterate through the headers and create a new snapshot
	snap := s.copy()

	// Copy values which might be changed by governance vote
	snap.Epoch, snap.Policy, snap.CommitteeSize = getGovernanceValue(gov, snap.Number)

	for _, header := range headers {
		// Remove any votes on checkpoint blocks
		number := header.Number.Uint64()

		// Resolve the authorization key and check against validators
		validator, err := ecrecover(header)
		if err != nil {
			return nil, err
		}
		if _, v := snap.ValSet.GetByAddress(validator); v == nil {
			return nil, errUnauthorized
		}

		if number%snap.Epoch == 0 {
			gov.UpdateCurrentSet(number)
			if len(header.Governance) > 0 {
				gov.WriteGovernanceForNextEpoch(number, header.Governance)
			}
			gov.ClearVotes(number)

			// Reload governance values because epoch changed
			snap.Epoch, snap.Policy, snap.CommitteeSize = getGovernanceValue(gov, number)
			snap.Votes = make([]governance.GovernanceVote, 0)
			snap.Tally = make([]governance.GovernanceTallyItem, 0)
		}

		snap.ValSet, snap.Votes, snap.Tally = gov.HandleGovernanceVote(snap.ValSet, snap.Votes, snap.Tally, header, validator, addr)
		if policy == uint64(params.WeightedRandom) {
			// Snapshot of block N (Snapshot_N) should contain proposers for N+1 and following blocks.
			// Validators for Block N+1 can be calculated based on the staking information from the previous stakingUpdateInterval block.
			// If the governance mode is single, the governing node is added to validator all the time.
			//
			// Proposers for Block N+1 can be calculated from the nearest previous proposersUpdateInterval block.
			// Refresh proposers in Snapshot_N using previous proposersUpdateInterval block for N+1, if not updated yet.
			isSingle, govNode, err := gov.GetGoverningInfoAtNumber(number)
			if err != nil {
				return nil, err
			}

			minStaking, err := gov.GetMinimumStakingAtNumber(number)
			if err != nil {
				return nil, err
			}

			pHeader := chain.GetHeaderByNumber(params.CalcProposerBlockNumber(number + 1))
			if pHeader != nil {
				if err := snap.ValSet.Refresh(pHeader.Hash(), pHeader.Number.Uint64(), chain.Config(), isSingle, govNode, minStaking); err != nil {
					// There are three error cases and they just don't refresh proposers
					// (1) no validator at all
					// (2) invalid formatted hash
					// (3) no staking info available
					logger.Trace("Skip refreshing proposers while creating snapshot", "snap.Number", snap.Number, "pHeader.Number", pHeader.Number.Uint64(), "err", err)
				}
			} else {
				logger.Trace("Can't refreshing proposers while creating snapshot due to lack of required header", "snap.Number", snap.Number)
			}
		}
	}
	snap.Number += uint64(len(headers))
	snap.Hash = headers[len(headers)-1].Hash()

	if snap.ValSet.Policy() == istanbul.WeightedRandom {
		// TODO-Klaytn-Issue1166 We have to update block number of ValSet too.
		snap.ValSet.SetBlockNum(snap.Number)
	}
	snap.ValSet.SetSubGroupSize(snap.CommitteeSize)

	gov.SetTotalVotingPower(snap.ValSet.TotalVotingPower())
	gov.SetMyVotingPower(snap.getMyVotingPower(addr))

	return snap, nil
}

func (s *Snapshot) getMyVotingPower(addr common.Address) uint64 {
	for _, a := range s.ValSet.List() {
		if a.Address() == addr {
			return a.VotingPower()
		}
	}
	return 0
}

// validators retrieves the list of authorized validators in ascending order.
func (s *Snapshot) validators() []common.Address {
	validators := make([]common.Address, 0, s.ValSet.Size())
	for _, validator := range s.ValSet.List() {
		validators = append(validators, validator.Address())
	}
	return sortValidatorArray(validators)
}

// demotedValidators retrieves the list of authorized, but demoted validators in ascending order.
func (s *Snapshot) demotedValidators() []common.Address {
	demotedValidators := make([]common.Address, 0, len(s.ValSet.DemotedList()))
	for _, demotedValidator := range s.ValSet.DemotedList() {
		demotedValidators = append(demotedValidators, demotedValidator.Address())
	}
	return sortValidatorArray(demotedValidators)
}

func (s *Snapshot) committee(prevHash common.Hash, view *istanbul.View) []common.Address {
	committeeList := s.ValSet.SubList(prevHash, view)

	committee := make([]common.Address, 0, len(committeeList))
	for _, v := range committeeList {
		committee = append(committee, v.Address())
	}
	return committee
}

func sortValidatorArray(validators []common.Address) []common.Address {
	for i := 0; i < len(validators); i++ {
		for j := i + 1; j < len(validators); j++ {
			if bytes.Compare(validators[i][:], validators[j][:]) > 0 {
				validators[i], validators[j] = validators[j], validators[i]
			}
		}
	}
	return validators
}

type snapshotJSON struct {
	Epoch  uint64                           `json:"epoch"`
	Number uint64                           `json:"number"`
	Hash   common.Hash                      `json:"hash"`
	Votes  []governance.GovernanceVote      `json:"votes"`
	Tally  []governance.GovernanceTallyItem `json:"tally"`

	// for validator set
	Validators   []common.Address        `json:"validators"`
	Policy       istanbul.ProposerPolicy `json:"policy"`
	SubGroupSize uint64                  `json:"subgroupsize"`

	// for weighted validator
	RewardAddrs       []common.Address `json:"rewardAddrs"`
	VotingPowers      []uint64         `json:"votingPower"`
	Weights           []uint64         `json:"weight"`
	Proposers         []common.Address `json:"proposers"`
	ProposersBlockNum uint64           `json:"proposersBlockNum"`
	DemotedValidators []common.Address `json:"demotedValidators"`
}

func (s *Snapshot) toJSONStruct() *snapshotJSON {
	var rewardAddrs []common.Address
	var votingPowers []uint64
	var weights []uint64
	var proposers []common.Address
	var proposersBlockNum uint64
	var validators []common.Address
	var demotedValidators []common.Address

	// TODO-Klaytn-Issue1166 For weightedCouncil
	if s.ValSet.Policy() == istanbul.WeightedRandom {
		validators, demotedValidators, rewardAddrs, votingPowers, weights, proposers, proposersBlockNum = validator.GetWeightedCouncilData(s.ValSet)
	} else {
		validators = s.validators()
	}

	return &snapshotJSON{
		Epoch:             s.Epoch,
		Number:            s.Number,
		Hash:              s.Hash,
		Votes:             s.Votes,
		Tally:             s.Tally,
		Validators:        validators,
		Policy:            istanbul.ProposerPolicy(s.Policy),
		SubGroupSize:      s.CommitteeSize,
		RewardAddrs:       rewardAddrs,
		VotingPowers:      votingPowers,
		Weights:           weights,
		Proposers:         proposers,
		ProposersBlockNum: proposersBlockNum,
		DemotedValidators: demotedValidators,
	}
}

// Unmarshal from a json byte array
func (s *Snapshot) UnmarshalJSON(b []byte) error {
	var j snapshotJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}

	s.Epoch = j.Epoch
	s.Number = j.Number
	s.Hash = j.Hash
	s.Votes = j.Votes
	s.Tally = j.Tally

	// TODO-Klaytn-Issue1166 For weightedCouncil
	if j.Policy == istanbul.WeightedRandom {
		s.ValSet = validator.NewWeightedCouncil(j.Validators, j.DemotedValidators, j.RewardAddrs, j.VotingPowers, j.Weights, j.Policy, j.SubGroupSize, j.Number, j.ProposersBlockNum, nil)
		validator.RecoverWeightedCouncilProposer(s.ValSet, j.Proposers)
	} else {
		s.ValSet = validator.NewSubSet(j.Validators, j.Policy, j.SubGroupSize)
	}
	return nil
}

// Marshal to a json byte array
func (s *Snapshot) MarshalJSON() ([]byte, error) {
	j := s.toJSONStruct()
	return json.Marshal(j)
}
