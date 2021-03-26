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
// This file is derived from quorum/consensus/istanbul/core/core.go (2018/06/04).
// Modified and improved for the klaytn development.

package core

import (
	"bytes"
	"math"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/prque"
	"github.com/klaytn/klaytn/consensus/istanbul"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/log"
	"github.com/rcrowley/go-metrics"
)

var logger = log.NewModuleLogger(log.ConsensusIstanbulCore)

// New creates an Istanbul consensus core
func New(backend istanbul.Backend, config *istanbul.Config) Engine {
	c := &core{
		config:             config,
		address:            backend.Address(),
		state:              StateAcceptRequest,
		handlerWg:          new(sync.WaitGroup),
		logger:             logger.NewWith("address", backend.Address()),
		backend:            backend,
		backlogs:           make(map[common.Address]*prque.Prque),
		backlogsMu:         new(sync.Mutex),
		pendingRequests:    prque.New(),
		pendingRequestsMu:  new(sync.Mutex),
		consensusTimestamp: time.Time{},

		roundMeter:         metrics.NewRegisteredMeter("consensus/istanbul/core/round", nil),
		currentRoundGauge:  metrics.NewRegisteredGauge("consensus/istanbul/core/currentRound", nil),
		sequenceMeter:      metrics.NewRegisteredMeter("consensus/istanbul/core/sequence", nil),
		consensusTimeGauge: metrics.NewRegisteredGauge("consensus/istanbul/core/timer", nil),
		councilSizeGauge:   metrics.NewRegisteredGauge("consensus/istanbul/core/councilSize", nil),
		committeeSizeGauge: metrics.NewRegisteredGauge("consensus/istanbul/core/committeeSize", nil),
		hashLockGauge:      metrics.NewRegisteredGauge("consensus/istanbul/core/hashLock", nil),
	}
	c.validateFn = c.checkValidatorSignature
	return c
}

// ----------------------------------------------------------------------------

type core struct {
	config  *istanbul.Config
	address common.Address
	state   State
	logger  log.Logger

	backend               istanbul.Backend
	events                *event.TypeMuxSubscription
	finalCommittedSub     *event.TypeMuxSubscription
	timeoutSub            *event.TypeMuxSubscription
	futurePreprepareTimer *time.Timer

	valSet                istanbul.ValidatorSet
	waitingForRoundChange bool
	validateFn            func([]byte, []byte) (common.Address, error)

	backlogs   map[common.Address]*prque.Prque
	backlogsMu *sync.Mutex

	current   *roundState
	handlerWg *sync.WaitGroup

	roundChangeSet    *roundChangeSet
	roundChangeTimer  atomic.Value //*time.Timer
	pendingRequests   *prque.Prque
	pendingRequestsMu *sync.Mutex

	consensusTimestamp time.Time
	// the meter to record the round change rate
	roundMeter metrics.Meter
	// the gauge to record the current round
	currentRoundGauge metrics.Gauge
	// the meter to record the sequence update rate
	sequenceMeter metrics.Meter
	// the gauge to record consensus duration (from accepting a preprepare to final committed stage)
	consensusTimeGauge metrics.Gauge
	// the gauge to record hashLock status (1 if hash-locked. 0 otherwise)
	hashLockGauge metrics.Gauge

	councilSizeGauge   metrics.Gauge
	committeeSizeGauge metrics.Gauge
}

func (c *core) finalizeMessage(msg *message) ([]byte, error) {
	var err error
	// Add sender address
	msg.Address = c.Address()

	// Add proof of consensus
	msg.CommittedSeal = []byte{}
	// Assign the CommittedSeal if it's a COMMIT message and proposal is not nil
	if msg.Code == msgCommit && c.current.Proposal() != nil {
		seal := PrepareCommittedSeal(c.current.Proposal().Hash())
		msg.CommittedSeal, err = c.backend.Sign(seal)
		if err != nil {
			return nil, err
		}
	}

	// Sign message
	data, err := msg.PayloadNoSig()
	if err != nil {
		return nil, err
	}
	msg.Signature, err = c.backend.Sign(data)
	if err != nil {
		return nil, err
	}

	// Convert to payload
	payload, err := msg.Payload()
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func (c *core) broadcast(msg *message) {
	logger := c.logger.NewWith("state", c.state)

	payload, err := c.finalizeMessage(msg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", msg, "err", err)
		return
	}

	// Broadcast payload
	if err = c.backend.Broadcast(msg.Hash, c.valSet, payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", msg, "err", err)
		return
	}
}

func (c *core) currentView() *istanbul.View {
	return &istanbul.View{
		Sequence: new(big.Int).Set(c.current.Sequence()),
		Round:    new(big.Int).Set(c.current.Round()),
	}
}

func (c *core) isProposer() bool {
	v := c.valSet
	if v == nil {
		return false
	}
	return v.IsProposer(c.backend.Address())
}

func (c *core) commit() {
	c.setState(StateCommitted)

	proposal := c.current.Proposal()
	if proposal != nil {
		committedSeals := make([][]byte, c.current.Commits.Size())
		for i, v := range c.current.Commits.Values() {
			committedSeals[i] = make([]byte, types.IstanbulExtraSeal)
			copy(committedSeals[i][:], v.CommittedSeal[:])
		}

		if err := c.backend.Commit(proposal, committedSeals); err != nil {
			c.current.UnlockHash() //Unlock block when insertion fails
			c.sendNextRoundChange("commit failure")
			return
		}
	} else {
		// TODO-Klaytn never happen, but if proposal is nil, mining is not working.
		logger.Error("istanbul.core current.Proposal is NULL")
		c.current.UnlockHash() //Unlock block when insertion fails
		c.sendNextRoundChange("commit failure. proposal is nil")
		return
	}
}

// startNewRound starts a new round. if round equals to 0, it means to starts a new sequence
func (c *core) startNewRound(round *big.Int) {
	var logger log.Logger
	if c.current == nil {
		logger = c.logger.NewWith("old_round", -1, "old_seq", 0)
	} else {
		logger = c.logger.NewWith("old_round", c.current.Round(), "old_seq", c.current.Sequence())
	}

	roundChange := false
	// Try to get last proposal
	lastProposal, lastProposer := c.backend.LastProposal()
	//if c.valSet != nil && c.valSet.IsSubSet() {
	//	c.current = nil
	//} else {
	if c.current == nil {
		logger.Trace("Start to the initial round")
	} else if lastProposal.Number().Cmp(c.current.Sequence()) >= 0 {
		diff := new(big.Int).Sub(lastProposal.Number(), c.current.Sequence())
		c.sequenceMeter.Mark(new(big.Int).Add(diff, common.Big1).Int64())

		if !c.consensusTimestamp.IsZero() {
			c.consensusTimeGauge.Update(int64(time.Since(c.consensusTimestamp)))
			c.consensusTimestamp = time.Time{}
		}
		logger.Trace("Catch up latest proposal", "number", lastProposal.Number().Uint64(), "hash", lastProposal.Hash())
	} else if lastProposal.Number().Cmp(big.NewInt(c.current.Sequence().Int64()-1)) == 0 {
		if round.Cmp(common.Big0) == 0 {
			// same seq and round, don't need to start new round
			return
		} else if round.Cmp(c.current.Round()) < 0 {
			logger.Warn("New round should not be smaller than current round", "seq", lastProposal.Number().Int64(), "new_round", round, "old_round", c.current.Round())
			return
		}
		roundChange = true
	} else {
		logger.Warn("New sequence should be larger than current sequence", "new_seq", lastProposal.Number().Int64())
		return
	}
	//}

	var newView *istanbul.View
	if roundChange {
		newView = &istanbul.View{
			Sequence: new(big.Int).Set(c.current.Sequence()),
			Round:    new(big.Int).Set(round),
		}
	} else {
		newView = &istanbul.View{
			Sequence: new(big.Int).Add(lastProposal.Number(), common.Big1),
			Round:    new(big.Int),
		}
		c.valSet = c.backend.Validators(lastProposal)

		councilSize := int64(c.valSet.Size())
		committeeSize := int64(c.valSet.SubGroupSize())
		if committeeSize > councilSize {
			committeeSize = councilSize
		}
		c.councilSizeGauge.Update(councilSize)
		c.committeeSizeGauge.Update(committeeSize)
	}
	c.backend.SetCurrentView(newView)

	// Update logger
	logger = logger.NewWith("old_proposer", c.valSet.GetProposer())
	// Clear invalid ROUND CHANGE messages
	c.roundChangeSet = newRoundChangeSet(c.valSet)
	// New snapshot for new round
	c.updateRoundState(newView, c.valSet, roundChange)
	// Calculate new proposer
	c.valSet.CalcProposer(lastProposer, newView.Round.Uint64())
	c.waitingForRoundChange = false
	c.setState(StateAcceptRequest)
	if roundChange && c.isProposer() && c.current != nil {
		// If it is locked, propose the old proposal
		// If we have pending request, propose pending request
		if c.current.IsHashLocked() {
			r := &istanbul.Request{
				Proposal: c.current.Proposal(), //c.current.Proposal would be the locked proposal by previous proposer, see updateRoundState
			}
			c.sendPreprepare(r)
		} else if c.current.pendingRequest != nil {
			c.sendPreprepare(c.current.pendingRequest)
		}
	}
	c.newRoundChangeTimer()

	logger.Debug("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "new_proposer", c.valSet.GetProposer(), "isProposer", c.isProposer())
	logger.Trace("New round", "new_round", newView.Round, "new_seq", newView.Sequence, "size", c.valSet.Size(), "valSet", c.valSet.List())
}

func (c *core) catchUpRound(view *istanbul.View) {
	logger := c.logger.NewWith("old_round", c.current.Round(), "old_seq", c.current.Sequence(), "old_proposer", c.valSet.GetProposer())

	if view.Round.Cmp(c.current.Round()) > 0 {
		c.roundMeter.Mark(new(big.Int).Sub(view.Round, c.current.Round()).Int64())
	}
	c.waitingForRoundChange = true

	// Need to keep block locked for round catching up
	c.updateRoundState(view, c.valSet, true)
	c.roundChangeSet.Clear(view.Round)

	c.newRoundChangeTimer()
	logger.Warn("[RC] Catch up round", "new_round", view.Round, "new_seq", view.Sequence, "new_proposer", c.valSet.GetProposer())
}

// updateRoundState updates round state by checking if locking block is necessary
func (c *core) updateRoundState(view *istanbul.View, validatorSet istanbul.ValidatorSet, roundChange bool) {
	// Lock only if both roundChange is true and it is locked
	if roundChange && c.current != nil {
		if c.current.IsHashLocked() {
			c.current = newRoundState(view, validatorSet, c.current.GetLockedHash(), c.current.Preprepare, c.current.pendingRequest, c.backend.HasBadProposal)
		} else {
			c.current = newRoundState(view, validatorSet, common.Hash{}, nil, c.current.pendingRequest, c.backend.HasBadProposal)
		}
	} else {
		c.current = newRoundState(view, validatorSet, common.Hash{}, nil, nil, c.backend.HasBadProposal)
	}
	c.currentRoundGauge.Update(c.current.round.Int64())
	if c.current.IsHashLocked() {
		c.hashLockGauge.Update(1)
	} else {
		c.hashLockGauge.Update(0)
	}
}

func (c *core) setState(state State) {
	if c.state != state {
		c.state = state
	}
	if state == StateAcceptRequest {
		c.processPendingRequests()
	}
	c.processBacklog()
}

func (c *core) Address() common.Address {
	return c.address
}

func (c *core) stopFuturePreprepareTimer() {
	if c.futurePreprepareTimer != nil {
		c.futurePreprepareTimer.Stop()
	}
}

func (c *core) stopTimer() {
	c.stopFuturePreprepareTimer()

	if c.roundChangeTimer.Load() != nil {
		c.roundChangeTimer.Load().(*time.Timer).Stop()
	}
}

func (c *core) newRoundChangeTimer() {
	c.stopTimer()

	// TODO-Klaytn-Istanbul: Replace &istanbul.DefaultConfig.Timeout to c.config.Timeout
	// set timeout based on the round number
	timeout := time.Duration(atomic.LoadUint64(&istanbul.DefaultConfig.Timeout)) * time.Millisecond
	round := c.current.Round().Uint64()
	if round > 0 {
		timeout += time.Duration(math.Pow(2, float64(round))) * time.Second
	}

	current := c.current
	proposer := c.valSet.GetProposer()

	c.roundChangeTimer.Store(time.AfterFunc(timeout, func() {
		var loc, proposerStr string

		if round == 0 {
			loc = "startNewRound"
		} else {
			loc = "catchUpRound"
		}
		if proposer == nil {
			proposerStr = ""
		} else {
			proposerStr = proposer.String()
		}

		if c.backend.NodeType() == common.CONSENSUSNODE {
			// Write log messages for validator activities analysis
			preparesSize := current.Prepares.Size()
			commitsSize := current.Commits.Size()
			logger.Warn("[RC] timeoutEvent Sent!", "set by", loc, "sequence",
				current.sequence, "round", current.round, "proposer", proposerStr, "preprepare is nil?",
				current.Preprepare == nil, "len(prepares)", preparesSize, "len(commits)", commitsSize)

			if preparesSize > 0 {
				logger.Warn("[RC] Prepares:", "messages", current.Prepares.GetMessages())
			}
			if commitsSize > 0 {
				logger.Warn("[RC] Commits:", "messages", current.Commits.GetMessages())
			}
		}

		c.sendEvent(timeoutEvent{&istanbul.View{
			Sequence: current.sequence,
			Round:    new(big.Int).Add(current.round, common.Big1),
		}})
	}))

	logger.Debug("New RoundChangeTimer Set", "seq", c.current.Sequence(), "round", round, "timeout", timeout)
}

func (c *core) checkValidatorSignature(data []byte, sig []byte) (common.Address, error) {
	return istanbul.CheckValidatorSignature(c.valSet, data, sig)
}

// PrepareCommittedSeal returns a committed seal for the given hash
func PrepareCommittedSeal(hash common.Hash) []byte {
	var buf bytes.Buffer
	buf.Write(hash.Bytes())
	buf.Write([]byte{byte(msgCommit)})
	return buf.Bytes()
}
