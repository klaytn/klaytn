// Copyright 2022 The klaytn Authors
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
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

//go:generate mockgen -package governance -destination=interface_mock_test.go github.com/klaytn/klaytn/governance Engine
type Engine interface {
	HeaderEngine
	ReaderEngine
	HeaderGov() HeaderEngine
	ContractGov() ReaderEngine
}

type ReaderEngine interface {
	// CurrentParams returns the params at the current block. The returned params shall be
	// used to build the upcoming (head+1) block. Block processing codes
	// should use this method.
	CurrentParams() *params.GovParamSet

	// EffectiveParams returns the params at given block number. The returned params
	// were used to build the block at given number.
	// The number must be equal or less than current block height (head).
	EffectiveParams(num uint64) (*params.GovParamSet, error)

	// UpdateParams updates the current params (the ones returned by CurrentParams()).
	// by reading the latest blockchain states.
	// This function must be called after every block is mined to
	// guarantee that CurrentParams() works correctly.
	UpdateParams(num uint64) error
}

type HeaderEngine interface {
	// AddVote casts votes from API
	AddVote(key string, val interface{}) bool
	ValidateVote(vote *GovernanceVote) (*GovernanceVote, bool)

	// Access database for voting states
	CanWriteGovernanceState(num uint64) bool
	WriteGovernanceState(num uint64, isCheckpoint bool) error

	// Access database for network params
	ReadGovernance(num uint64) (uint64, map[string]interface{}, error)
	WriteGovernance(num uint64, data GovernanceSet, delta GovernanceSet) error

	// Compose header.Vote and header.Governance
	GetEncodedVote(addr common.Address, number uint64) []byte
	GetGovernanceChange() map[string]interface{}

	// Intake header.Vote and header.Governance
	VerifyGovernance(received []byte) error
	ClearVotes(num uint64)
	WriteGovernanceForNextEpoch(number uint64, governance []byte)
	UpdateCurrentSet(num uint64)
	HandleGovernanceVote(
		valset istanbul.ValidatorSet, votes []GovernanceVote, tally []GovernanceTallyItem,
		header *types.Header, proposer common.Address, self common.Address, writable bool) (
		istanbul.ValidatorSet, []GovernanceVote, []GovernanceTallyItem)

	// Get internal fields
	GetVoteMapCopy() map[string]VoteStatus
	GetGovernanceTalliesCopy() []GovernanceTallyItem
	CurrentSetCopy() map[string]interface{}
	PendingChanges() map[string]interface{}
	Votes() []GovernanceVote
	IdxCache() []uint64
	IdxCacheFromDb() []uint64

	NodeAddress() common.Address
	TotalVotingPower() uint64
	MyVotingPower() uint64
	BlockChain() blockChain
	DB() database.DBManager

	// Set internal fields
	SetNodeAddress(addr common.Address)
	SetTotalVotingPower(t uint64)
	SetMyVotingPower(t uint64)
	SetBlockchain(chain blockChain)
	SetTxPool(txpool txPool)
	GetTxPool() txPool
}

// blockChain is an interface for blockchain.Blockchain used in governance package.
type blockChain interface {
	blockchain.ChainContext

	CurrentBlock() *types.Block
	GetHeaderByNumber(val uint64) *types.Header
	StateAt(root common.Hash) (*state.StateDB, error)
	Config() *params.ChainConfig
}
