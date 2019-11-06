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

package sc

import (
	"github.com/klaytn/klaytn/accounts/abi/bind"
	"github.com/pkg/errors"
	"sync"
	"time"
)

var (
	filterLogsStride = uint64(100)
	maxPendingTxs    = 1000
)

// valueTransferHint stores the last handled block number and nonce (Request or Handle).
type valueTransferHint struct {
	blockNumber     uint64 // block number to start searching event logs
	requestNonce    uint64
	handleNonce     uint64
	prevHandleNonce uint64 // previous handleNonce between recovery interval
	candidate       bool   // to check recovery candidate between recovery interval
}

// valueTransferRecovery stores status information for the value transfer recovery.
type valueTransferRecovery struct {
	stopCh    chan interface{}
	isRunning bool           // to check duplicated start
	wg        sync.WaitGroup // wait group to handle the Stop() sync

	child2parentHint *valueTransferHint
	parent2childHint *valueTransferHint
	childEvents      []*RequestValueTransferEvent
	parentEvents     []*RequestValueTransferEvent

	config      *SCConfig
	cBridgeInfo *BridgeInfo
	pBridgeInfo *BridgeInfo
}

var (
	ErrVtrDisabled       = errors.New("VTR is disabled")
	ErrVtrAlreadyStarted = errors.New("VTR is already started")
)

// NewValueTransferRecovery creates a new value transfer recovery structure.
func NewValueTransferRecovery(config *SCConfig, cBridgeInfo, pBridgeInfo *BridgeInfo) *valueTransferRecovery {
	return &valueTransferRecovery{
		stopCh:           make(chan interface{}),
		isRunning:        false,
		wg:               sync.WaitGroup{},
		child2parentHint: &valueTransferHint{},
		parent2childHint: &valueTransferHint{},
		childEvents:      []*RequestValueTransferEvent{},
		parentEvents:     []*RequestValueTransferEvent{},
		config:           config,
		cBridgeInfo:      cBridgeInfo,
		pBridgeInfo:      pBridgeInfo,
	}
}

// Start implements starting all internal goroutines used by the value transfer recovery.
func (vtr *valueTransferRecovery) Start() error {
	if !vtr.config.VTRecovery {
		return ErrVtrDisabled
	}

	// TODO-Klaytn-Servicechain If there is no user API to start recovery, remove isRunning in Start/Stop.
	if vtr.isRunning {
		return ErrVtrAlreadyStarted
	}

	vtr.wg.Add(1)

	go func() {
		ticker := time.NewTicker(time.Duration(vtr.config.VTRecoveryInterval) * time.Second)
		defer func() {
			ticker.Stop()
			vtr.wg.Done()
		}()

		if err := vtr.Recover(); err != nil {
			logger.Warn("initial value transfer recovery is failed", "err", err)
		}

		vtr.isRunning = true

		for {
			select {
			case <-vtr.stopCh:
				logger.Info("value transfer recovery is stopped")
				return
			case <-ticker.C:
				if vtr.isRunning {
					if err := vtr.Recover(); err != nil {
						logger.Trace("value transfer recovery is failed", "err", err)
					}
				}
			}
		}
	}()

	return nil
}

// Stop implements terminating all internal goroutines used by the value transfer recovery.
func (vtr *valueTransferRecovery) Stop() error {
	if !vtr.isRunning {
		logger.Info("value transfer recovery is already stopped")
		return nil
	}
	close(vtr.stopCh)
	vtr.wg.Wait()
	vtr.isRunning = false
	return nil
}

// Recover implements the whole recovery process of the value transfer recovery.
func (vtr *valueTransferRecovery) Recover() error {
	logger.Trace("update value transfer hint")
	err := vtr.updateRecoveryHint()
	if err != nil {
		return err
	}

	logger.Trace("retrieve pending events")
	err = vtr.retrievePendingEvents()
	if err != nil {
		return err
	}

	logger.Trace("recover pending events")
	err = vtr.recoverPendingEvents()
	if err != nil {
		return err
	}

	return nil
}

// updateRecoveryHint updates hints for value transfers on the both side.
// One is from child chain to parent chain, the other is from parent chain to child chain value transfers.
// The hint includes a block number to begin search, request nonce and handle nonce.
func (vtr *valueTransferRecovery) updateRecoveryHint() error {
	if vtr.cBridgeInfo == nil {
		return errors.New("child chain bridge is nil")
	}
	if vtr.pBridgeInfo == nil {
		return errors.New("parent chain bridge is nil")
	}

	var err error
	vtr.child2parentHint, err = updateRecoveryHintFromTo(vtr.child2parentHint, vtr.cBridgeInfo, vtr.pBridgeInfo)
	if err != nil {
		return err
	}

	vtr.parent2childHint, err = updateRecoveryHintFromTo(vtr.parent2childHint, vtr.pBridgeInfo, vtr.cBridgeInfo)
	if err != nil {
		return err
	}

	// Update the hint for the initial status.
	if !vtr.isRunning {
		vtr.child2parentHint.prevHandleNonce = vtr.child2parentHint.handleNonce
		vtr.parent2childHint.prevHandleNonce = vtr.parent2childHint.handleNonce
		vtr.child2parentHint.candidate = true
		vtr.parent2childHint.candidate = true
	}

	return nil
}

// updateRecoveryHint updates a hint for the one-way value transfers.
func updateRecoveryHintFromTo(prevHint *valueTransferHint, from, to *BridgeInfo) (*valueTransferHint, error) {
	var err error
	var hint valueTransferHint

	logger.Trace("updateRecoveryHintFromTo start")
	if prevHint != nil {
		logger.Trace("recovery prevHint", "rnonce", prevHint.requestNonce, "hnonce", prevHint.handleNonce, "phnonce", prevHint.prevHandleNonce, "cand", prevHint.candidate)
	}

	hint.blockNumber, err = to.bridge.RecoveryBlockNumber(nil)
	if err != nil {
		return nil, err
	}

	requestNonce, err := from.bridge.RequestNonce(nil)
	if err != nil {
		return nil, err
	}
	hint.requestNonce = requestNonce

	handleNonce, err := to.bridge.LowerHandleNonce(nil)
	if err != nil {
		return nil, err
	}
	if prevHint != nil {
		hint.prevHandleNonce = prevHint.handleNonce
		hint.candidate = prevHint.candidate
	}
	hint.handleNonce = handleNonce

	logger.Trace("updateRecoveryHintFromTo finish", "rnonce", hint.requestNonce, "hnonce", hint.handleNonce, "phnonce", hint.prevHandleNonce, "cand", hint.candidate)

	return &hint, nil
}

// retrievePendingEvents retrieves pending events on the child chain or parent chain.
// The pending event is the value transfer without processing HandleValueTransfer.
func (vtr *valueTransferRecovery) retrievePendingEvents() error {
	if vtr.cBridgeInfo.bridge == nil {
		return errors.New("bridge is nil")
	}

	var err error
	vtr.childEvents, err = retrievePendingEventsFrom(vtr.child2parentHint, vtr.cBridgeInfo)
	if err != nil {
		return err
	}
	vtr.parentEvents, err = retrievePendingEventsFrom(vtr.parent2childHint, vtr.pBridgeInfo)
	if err != nil {
		return err
	}

	return nil
}

// retrievePendingEventsFrom retrieves pending events from the specified bridge by using the hint provided.
// The filter uses a hint as a search range. It returns a slice of events that has log details.
func retrievePendingEventsFrom(hint *valueTransferHint, bi *BridgeInfo) ([]*RequestValueTransferEvent, error) {
	if bi.bridge == nil {
		return nil, errors.New("bridge is nil")
	}
	if hint.requestNonce == hint.handleNonce {
		return nil, nil
	}
	if !checkRecoveryCondition(hint) {
		return nil, nil
	}

	var pendingEvents []*RequestValueTransferEvent

	curBlkNum, err := bi.GetCurrentBlockNumber()
	if err != nil {
		return nil, err
	}

	startBlkNum := hint.blockNumber
	endBlkNum := startBlkNum + filterLogsStride

pendingTxLoop:
	for startBlkNum <= curBlkNum {
		if endBlkNum > curBlkNum {
			endBlkNum = curBlkNum
		}
		opts := &bind.FilterOpts{Start: startBlkNum, End: &endBlkNum}
		it, err := bi.bridge.FilterRequestValueTransfer(opts, nil, nil, nil)
		if err != nil {
			return nil, err
		}
		for it.Next() {
			logger.Trace("pending nonce in the event", "requestNonce", it.Event.RequestNonce)
			if it.Event.RequestNonce >= hint.handleNonce {
				logger.Trace("filtered pending nonce", "requestNonce", it.Event.RequestNonce, "handledNonce", hint.handleNonce)
				pendingEvents = append(pendingEvents, &RequestValueTransferEvent{it.Event})
				if len(pendingEvents) >= maxPendingTxs {
					it.Close()
					break pendingTxLoop
				}
			}
		}
		startBlkNum = endBlkNum + 1
		endBlkNum = startBlkNum + filterLogsStride
		it.Close()
	}

	if len(pendingEvents) > 0 {
		logger.Info("retrieved pending events", "bridge", bi.address.String(), "len(pendingEvents)", len(pendingEvents), "1st nonce", pendingEvents[0].Nonce())
	}
	return pendingEvents, nil
}

// checkRecoveryCandidateCondition checks if vtr is recovery candidate or not.
// candidate is introduced to check any normal request just before checking start.
//
// For example,
//
// ======== ======== ======== ========
// Round    R Nonce  H Nonce  Result
// ======== ======== ======== ========
// 1        10       10       false
// <burst requests just before checking>
// 2        1000     10       ? (it can be normal but candidate)
// 3        2000     10       true
func checkRecoveryCandidateCondition(hint *valueTransferHint) bool {
	return hint.requestNonce != hint.handleNonce && hint.prevHandleNonce == hint.handleNonce
}

// checkRecoveryCondition checks if recovery for the handle value transfers is needed or not.
func checkRecoveryCondition(hint *valueTransferHint) bool {
	if checkRecoveryCandidateCondition(hint) && hint.candidate {
		hint.candidate = false
		return true
	}
	if checkRecoveryCandidateCondition(hint) && !hint.candidate {
		hint.candidate = true
		return false
	}
	hint.candidate = false
	return false
}

// recoverPendingEvents recovers all pending events by resending them.
func (vtr *valueTransferRecovery) recoverPendingEvents() error {
	defer func() {
		vtr.childEvents = []*RequestValueTransferEvent{}
		vtr.parentEvents = []*RequestValueTransferEvent{}
	}()

	if len(vtr.childEvents) > 0 {
		logger.Warn("VT Recovery : Child -> Parent Chain", "cBridge", vtr.cBridgeInfo.address.String(), "events", len(vtr.childEvents))
	}

	vtRequestEventMeter.Mark(int64(len(vtr.childEvents)))
	vtRecoveredRequestEventMeter.Mark(int64(len(vtr.childEvents)))

	vtr.pBridgeInfo.AddRequestValueTransferEvents(vtr.childEvents)

	if len(vtr.parentEvents) > 0 {
		logger.Warn("VT Recovery : Parent -> Child Chain", "pBridge", vtr.pBridgeInfo.address.String(), "events", len(vtr.parentEvents))
	}

	vtHandleEventMeter.Mark(int64(len(vtr.parentEvents)))
	vtr.cBridgeInfo.AddRequestValueTransferEvents(vtr.parentEvents)

	return nil
}

func (vtr *valueTransferRecovery) WaitRunningStatus(expected bool, timeout time.Duration) error {
	for i := 0; i < int(timeout/time.Second); i++ {
		if vtr.isRunning == expected {
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return errors.New("timeout to wait expect value")
}
