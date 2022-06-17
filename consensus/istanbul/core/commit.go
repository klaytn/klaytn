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
// This file is derived from quorum/consensus/istanbul/core/commit.go (2018/06/04).
// Modified and improved for the klaytn development.

package core

import (
	"reflect"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/consensus/istanbul"
)

func (c *core) sendCommit() {
	logger := c.logger.NewWith("state", c.state)
	if c.current.Preprepare == nil {
		logger.Error("Failed to get parentHash from roundState in sendCommit")
		return
	}

	sub := c.current.Subject()
	prevHash := c.current.Proposal().ParentHash()

	// Do not send message if the owner of the core is not a member of the committee for the `sub.View`
	if !c.valSet.CheckInSubList(prevHash, sub.View, c.Address()) {
		return
	}

	// TODO-Klaytn-Istanbul: generalize broadcastCommit for all istanbul message types
	c.broadcastCommit(sub)
}

func (c *core) sendCommitForOldBlock(view *istanbul.View, digest common.Hash, prevHash common.Hash) {
	sub := &istanbul.Subject{
		View:     view,
		Digest:   digest,
		PrevHash: prevHash,
	}
	c.broadcastCommit(sub)
}

func (c *core) broadcastCommit(sub *istanbul.Subject) {
	logger := c.logger.NewWith("state", c.state)

	encodedSubject, err := Encode(sub)
	if err != nil {
		logger.Error("Failed to encode", "subject", sub)
		return
	}

	c.broadcast(&message{
		Hash: sub.PrevHash,
		Code: msgCommit,
		Msg:  encodedSubject,
	})
}

func (c *core) handleCommit(msg *message, src istanbul.Validator) error {
	// Decode COMMIT message
	var commit *istanbul.Subject
	err := msg.Decode(&commit)
	if err != nil {
		return errFailedDecodeCommit
	}

	//logger.Error("receive handle commit","num", commit.View.Sequence)
	if err := c.checkMessage(msgCommit, commit.View); err != nil {
		//logger.Error("### istanbul/commit.go checkMessage","num",commit.View.Sequence,"err",err)
		return err
	}

	if err := c.verifyCommit(commit, src); err != nil {
		return err
	}

	if !c.valSet.CheckInSubList(msg.Hash, commit.View, src.Address()) {
		logger.Warn("received an istanbul commit message from non-committee",
			"currentSequence", c.current.sequence.Uint64(), "sender", src.Address().String(), "msgView", commit.View.String())
		return errNotFromCommittee
	}

	c.acceptCommit(msg, src)

	// Change to Prepared state if we've received enough PREPARE/COMMIT messages or it is locked
	// and we are in earlier state before Prepared state.
	// Both of PREPARE and COMMIT messages are counted since the nodes which is hashlocked in
	// the previous round skip sending PREPARE messages.
	if c.state.Cmp(StatePrepared) < 0 {
		if c.current.IsHashLocked() && commit.Digest == c.current.GetLockedHash() {
			logger.Warn("received commit of the hash locked proposal and change state to prepared", "msgType", msgCommit)
			c.setState(StatePrepared)
			c.sendCommit()
		} else if c.current.GetPrepareOrCommitSize() > 2*c.valSet.F() {
			logger.Info("received more than 2f agreements and change state to prepared", "msgType", msgCommit)
			c.current.LockHash()
			c.setState(StatePrepared)
			c.sendCommit()
		}
	}

	// Commit the proposal once we have enough COMMIT messages and we are not in the Committed state.
	//
	// If we already have a proposal, we may have chance to speed up the consensus process
	// by committing the proposal without PREPARE messages.
	//logger.Error("### consensus check","len(commits)",c.current.Commits.Size(),"f(2/3)",2*c.valSet.F(),"state",c.state.Cmp(StateCommitted))
	if c.state.Cmp(StateCommitted) < 0 && c.current.Commits.Size() > 2*c.valSet.F() {
		// Still need to call LockHash here since state can skip Prepared state and jump directly to the Committed state.
		c.current.LockHash()
		c.commit()
	}

	return nil
}

// verifyCommit verifies if the received COMMIT message is equivalent to our subject
func (c *core) verifyCommit(commit *istanbul.Subject, src istanbul.Validator) error {
	logger := c.logger.NewWith("from", src, "state", c.state)

	sub := c.current.Subject()
	if !reflect.DeepEqual(commit, sub) {
		logger.Warn("Inconsistent subjects between commit and proposal", "expected", sub, "got", commit)
		return errInconsistentSubject
	}

	return nil
}

func (c *core) acceptCommit(msg *message, src istanbul.Validator) error {
	logger := c.logger.NewWith("from", src, "state", c.state)

	// Add the COMMIT message to current round state
	if err := c.current.Commits.Add(msg); err != nil {
		logger.Error("Failed to record commit message", "msg", msg, "err", err)
		return err
	}

	return nil
}
