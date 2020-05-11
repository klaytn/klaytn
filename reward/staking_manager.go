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

package reward

import (
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

const (
	chainHeadChanSize = 100
)

// blockChain is an interface for blockchain.Blockchain used in reward package.
type blockChain interface {
	SubscribeChainHeadEvent(ch chan<- blockchain.ChainHeadEvent) event.Subscription
	GetBlockByNumber(number uint64) *types.Block
	StateAt(root common.Hash) (*state.StateDB, error)
	Config() *params.ChainConfig

	blockchain.ChainContext
}

type StakingManager struct {
	addressBookConnector *addressBookConnector
	stakingInfoCache     *stakingInfoCache
	dbManager            database.DBManager // to store staking info into db
	governanceHelper     governanceHelper
	blockchain           blockChain
	chainHeadChan        chan blockchain.ChainHeadEvent
	chainHeadSub         event.Subscription
	isActivated          bool
}

func NewStakingManager(bc blockChain, gh governanceHelper, dbm database.DBManager) *StakingManager {
	return &StakingManager{
		addressBookConnector: newAddressBookConnector(bc, gh),
		stakingInfoCache:     newStakingInfoCache(),
		dbManager:            dbm,
		governanceHelper:     gh,
		blockchain:           bc,
		chainHeadChan:        make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
		isActivated:          false,
	}
}

// GetStakingInfo returns a corresponding stakingInfo for a blockNum.
func (sm *StakingManager) GetStakingInfo(blockNum uint64) *StakingInfo {
	stakingBlockNumber := params.CalcStakingBlockNumber(blockNum)

	// Get staking info from cache
	if cachedStakingInfo := sm.stakingInfoCache.get(stakingBlockNumber); cachedStakingInfo != nil {
		logger.Debug("StakingInfoCache hit.", "Block number", blockNum, "staking block number", stakingBlockNumber, "stakingInfo", cachedStakingInfo)
		return cachedStakingInfo
	}

	// Get staking info from db
	if sm.dbManager != nil {
		if storedStakingInfo, err := sm.dbManager.ReadStakingInfo(stakingBlockNumber); storedStakingInfo != nil && err == nil {
			logger.Debug("StakingInfoDB hit.", "Block number", blockNum, "staking block number", stakingBlockNumber, "stakingInfo", storedStakingInfo)
			return storedStakingInfo.(*StakingInfo)
		} else {
			logger.Warn("Failed to get staking info from db.", "blockNum", blockNum, "err", err, "value", storedStakingInfo)
		}
	}

	// Get staking info from block header and updates it to cache and db
	calcStakingInfo, err := sm.updateStakingInfo(stakingBlockNumber)
	if err != nil {
		logger.Error("Failed to get stakingInfo", "Block number", blockNum, "staking block number", stakingBlockNumber, "err", err, "stakingInfo", calcStakingInfo)
		return nil
	}

	logger.Debug("Get stakingInfo from header.", "Block number", blockNum, "staking block number", stakingBlockNumber, "stakingInfo", calcStakingInfo)
	return calcStakingInfo
}

func (sm *StakingManager) IsActivated() bool {
	return sm.isActivated
}

// updateStakingInfo updates staking info in cache and db created from given block number.
func (sm *StakingManager) updateStakingInfo(blockNum uint64) (*StakingInfo, error) {
	stakingInfo, err := sm.addressBookConnector.getStakingInfoFromAddressBook(blockNum)
	if err != nil {
		return nil, err
	}
	sm.isActivated = true

	// update cache
	sm.stakingInfoCache.add(stakingInfo)

	// update db
	if sm.dbManager != nil {
		err := sm.dbManager.WriteStakingInfo(stakingInfo)
		if err != nil {
			logger.Warn("Failed to write staking info to db.", "err", err, "stakingInfo", stakingInfo)
		}
	}

	logger.Info("Add a new stakingInfo to stakingInfoCache and stakingInfoDB", "stakingInfo", stakingInfo)
	return stakingInfo, nil
}

// Subscribe setups a channel to listen chain head event and starts a goroutine to update staking cache.
func (sm *StakingManager) Subscribe() {
	sm.chainHeadSub = sm.blockchain.SubscribeChainHeadEvent(sm.chainHeadChan)

	go sm.handleChainHeadEvent()
}

func (sm *StakingManager) handleChainHeadEvent() {
	defer sm.Unsubscribe()

	logger.Info("Start listening chain head event to update stakingInfoCache.")

	for {
		// A real event arrived, process interesting content
		select {
		// Handle ChainHeadEvent
		case ev := <-sm.chainHeadChan:
			if sm.governanceHelper.ProposerPolicy() == params.WeightedRandom {
				stakingBlockNum := ev.Block.NumberU64() - ev.Block.NumberU64()%sm.governanceHelper.StakingUpdateInterval()
				if cachedStakingInfo := sm.stakingInfoCache.get(stakingBlockNum); cachedStakingInfo == nil {
					_, err := sm.updateStakingInfo(stakingBlockNum)
					if err != nil {
						logger.Error("Failed to update stakingInfoCache", "blockNumber", ev.Block.NumberU64(), "stakingNumber", stakingBlockNum, "err", err)
					}
				}
			}
		case <-sm.chainHeadSub.Err():
			return
		}
	}
}

// Unsubscribe can unsubscribe a subscription to listen chain head event.
func (sm *StakingManager) Unsubscribe() {
	sm.chainHeadSub.Unsubscribe()
}
