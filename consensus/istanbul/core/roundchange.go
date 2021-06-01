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
// This file is derived from quorum/consensus/istanbul/core/roundchange.go (2018/06/04).
// Modified and improved for the klaytn development.

package core

import (
	"math/big"
	"sync"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
)

// sendNextRoundChange sends the ROUND CHANGE message with current round + 1
func (c *core) sendNextRoundChange(loc string) {
	if c.backend.NodeType() != common.CONSENSUSNODE {
		return
	}
	logger.Warn("[RC] sendNextRoundChange happened", "where", loc)
	c.sendRoundChange(new(big.Int).Add(c.currentView().Round, common.Big1))
}

// sendRoundChange sends the ROUND CHANGE message with the given round
func (c *core) sendRoundChange(round *big.Int) {
	logger := c.logger.NewWith("state", c.state)

	cv := c.currentView()
	if cv.Round.Cmp(round) >= 0 {
		logger.Warn("[RC] Skip sending out the round change message", "current round", cv.Round,
			"target round", round)
		return
	}

	logger.Warn("[RC] Prepare messages received before catchUpRound",
		"len(prepares)", c.current.Prepares.Size(), "messages", c.current.Prepares.GetMessages())
	logger.Warn("[RC] Commit messages received before catchUpRound",
		"len(commits)", c.current.Commits.Size(), "messages", c.current.Commits.GetMessages())

	c.catchUpRound(&istanbul.View{
		// The round number we'd like to transfer to.
		Round:    new(big.Int).Set(round),
		Sequence: new(big.Int).Set(cv.Sequence),
	})

	lastProposal, _ := c.backend.LastProposal()

	// Now we have the new round number and sequence number
	cv = c.currentView()
	rc := &istanbul.Subject{
		View:     cv,
		Digest:   common.Hash{},
		PrevHash: lastProposal.Hash(),
	}

	payload, err := Encode(rc)
	if err != nil {
		logger.Error("Failed to encode ROUND CHANGE", "rc", rc, "err", err)
		return
	}

	c.broadcast(&message{
		Hash: rc.PrevHash,
		Code: msgRoundChange,
		Msg:  payload,
	})
}

func (c *core) handleRoundChange(msg *message, src istanbul.Validator) error {
	logger := c.logger.NewWith("state", c.state, "from", src.Address().Hex())

	// Decode ROUND CHANGE message
	var rc *istanbul.Subject
	if err := msg.Decode(&rc); err != nil {
		logger.Error("Failed to decode ROUND CHANGE", "err", err)
		return errInvalidMessage
	}

	// TODO-Klaytn-Istanbul: establish round change messaging policy and then apply it
	//if !c.valSet.CheckInSubList(msg.Hash, rc.View, src.Address()) {
	//	return errNotFromCommittee
	//}

	if err := c.checkMessage(msgRoundChange, rc.View); err != nil {
		return err
	}

	cv := c.currentView()
	roundView := rc.View

	// Add the ROUND CHANGE message to its message set and return how many
	// messages we've got with the same round number and sequence number.
	num, err := c.roundChangeSet.Add(roundView.Round, msg)
	if err != nil {
		logger.Warn("Failed to add round change message", "from", src, "msg", msg, "err", err)
		return err
	}

	var numCatchUp, numStartNewRound int
	if c.valSet.Size() < 4 {
		n := int(c.valSet.Size())
		// N ROUND CHANGE messages can start new round.
		numStartNewRound = n
		// N - 1 ROUND CHANGE messages can catch up the round.
		numCatchUp = n - 1
	} else {
		f := int(c.valSet.F())
		// 2*F + 1 ROUND CHANGE messages can start new round.
		numStartNewRound = 2*f + 1
		// F + 1 ROUND CHANGE messages can start catch up the round.
		numCatchUp = f + 1
	}

	if num == numStartNewRound && (c.waitingForRoundChange || cv.Round.Cmp(roundView.Round) < 0) {
		// We've received enough ROUND CHANGE messages, start a new round immediately.
		logger.Warn("[RC] Prepare messages received before startNewRound", "round", cv.Round.String(),
			"len(prepares)", c.current.Prepares.Size(), "messages", c.current.Prepares.GetMessages())
		logger.Warn("[RC] Commit messages received before startNewRound", "round", cv.Round.String(),
			"len(commits)", c.current.Commits.Size(), "messages", c.current.Commits.GetMessages())
		logger.Warn("[RC] Received 2f+1 Round Change Messages. Starting new round",
			"currentRound", cv.Round.String(), "newRound", roundView.Round.String())
		c.startNewRound(roundView.Round)
		return nil
	} else if c.waitingForRoundChange && num == numCatchUp {
		// Once we received enough ROUND CHANGE messages, those messages form a weak certificate.
		// If our round number is smaller than the certificate's round number, we would
		// try to catch up the round number.
		if cv.Round.Cmp(roundView.Round) < 0 {
			logger.Warn("[RC] Send round change because we have f+1 round change messages",
				"currentRound", cv.Round.String(), "newRound", roundView.Round.String())
			c.sendRoundChange(roundView.Round)
		}
		return nil
	} else if cv.Round.Cmp(roundView.Round) < 0 {
		// Only gossip the message with current round to other validators.
		logger.Warn("[RC] Received round is bigger but not enough number of messages. Message ignored",
			"currentRound", cv.Round.String(), "newRound", roundView.Round.String(), "numRC", num)
		return errIgnored
	}
	return nil
}

// ----------------------------------------------------------------------------

func newRoundChangeSet(valSet istanbul.ValidatorSet) *roundChangeSet {
	return &roundChangeSet{
		validatorSet: valSet,
		roundChanges: make(map[uint64]*messageSet),
		mu:           new(sync.Mutex),
	}
}

type roundChangeSet struct {
	validatorSet istanbul.ValidatorSet
	roundChanges map[uint64]*messageSet
	mu           *sync.Mutex
}

// Add adds the round and message into round change set
func (rcs *roundChangeSet) Add(r *big.Int, msg *message) (int, error) {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	round := r.Uint64()
	if rcs.roundChanges[round] == nil {
		rcs.roundChanges[round] = newMessageSet(rcs.validatorSet)
	}
	err := rcs.roundChanges[round].Add(msg)
	if err != nil {
		return 0, err
	}
	return rcs.roundChanges[round].Size(), nil
}

// Clear deletes the messages with smaller round
func (rcs *roundChangeSet) Clear(round *big.Int) {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	for k, rms := range rcs.roundChanges {
		if len(rms.Values()) == 0 || k < round.Uint64() {
			delete(rcs.roundChanges, k)
		}
	}
}

// MaxRound returns the max round which the number of messages is equal or larger than num
func (rcs *roundChangeSet) MaxRound(num int) *big.Int {
	rcs.mu.Lock()
	defer rcs.mu.Unlock()

	var maxRound *big.Int
	for k, rms := range rcs.roundChanges {
		if rms.Size() < num {
			continue
		}
		r := big.NewInt(int64(k))
		if maxRound == nil || maxRound.Cmp(r) < 0 {
			maxRound = r
		}
	}
	return maxRound
}
