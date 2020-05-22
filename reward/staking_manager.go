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
	"errors"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/state"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
	"sync"
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

type migrationCheck interface {
	InMigration() bool
}

type StakingManager struct {
	addressBookConnector *addressBookConnector
	stakingInfoCache     *stakingInfoCache
	stakingInfoDB        stakingInfoDB
	governanceHelper     governanceHelper
	blockchain           blockChain
	chainHeadChan        chan blockchain.ChainHeadEvent
	chainHeadSub         event.Subscription
	migrationCheck       migrationCheck
}

var (
	// variables for sole StakingManager
	once           sync.Once
	stakingManager *StakingManager

	// errors for staking manager
	ErrStakingManagerNotSet = errors.New("staking manager is not set")
	ErrChainHeadChanNotSet  = errors.New("chain head channel is not set")
)

// NewStakingManager creates and returns StakingManager.
//
// On the first call, a StakingManager is created with given parameters.
// From next calls, the existing StakingManager is returned. (Parameters
// from the next calls will not affect.)
func NewStakingManager(bc blockChain, gh governanceHelper, db database.DBManager) *StakingManager {
	if bc != nil && gh != nil {
		// this is only called once
		once.Do(func() {
			stakingManager = &StakingManager{
				addressBookConnector: newAddressBookConnector(bc, gh),
				stakingInfoCache:     newStakingInfoCache(),
				stakingInfoDB:        db,
				governanceHelper:     gh,
				blockchain:           bc,
				chainHeadChan:        make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
				migrationCheck:       migrationCheck(db),
			}

			// Before migration, staking information of current and before should be stored in DB.
			//
			// Staking information from block of StakingUpdateInterval ahead is needed to create a block.
			// If there is no staking info in either cache, db or state trie, the node cannot make a block.
			// The information in state trie is deleted after state trie migration.
			blockchain.RegisterMigrationPrerequisites(func(blockNum uint64) error {
				if _, err := updateStakingInfo(blockNum); err != nil {
					return err
				}
				_, err := updateStakingInfo(blockNum + params.StakingUpdateInterval())
				return err
			})
		})
	} else {
		logger.Error("unable to set StakingManager", "blockchain", bc, "governanceHelper", gh)
	}

	return stakingManager
}

func GetStakingManager() *StakingManager {
	return stakingManager
}

// GetStakingInfo returns a corresponding stakingInfo for a blockNum.
func GetStakingInfo(blockNum uint64) *StakingInfo {
	if stakingManager == nil {
		logger.Error("unable to GetStakingInfo", "err", ErrStakingManagerNotSet)
		return nil
	}

	stakingBlockNumber := params.CalcStakingBlockNumber(blockNum)

	// Get staking info from cache
	if cachedStakingInfo := stakingManager.stakingInfoCache.get(stakingBlockNumber); cachedStakingInfo != nil {
		logger.Debug("StakingInfoCache hit.", "blockNum", blockNum, "staking block number", stakingBlockNumber, "stakingInfo", cachedStakingInfo)
		return cachedStakingInfo
	}

	// Get staking info from DB
	if storedStakingInfo, err := getStakingInfoFromDB(stakingBlockNumber); storedStakingInfo != nil && err == nil {
		logger.Debug("StakingInfoDB hit.", "blockNum", blockNum, "staking block number", stakingBlockNumber, "stakingInfo", storedStakingInfo)
		stakingManager.stakingInfoCache.add(storedStakingInfo)
		return storedStakingInfo
	} else {
		logger.Debug("failed to get stakingInfo from DB", "err", err, "blockNum", blockNum)
	}

	// Calculate staking info from block header and updates it to cache and db
	calcStakingInfo, _ := updateStakingInfo(stakingBlockNumber)
	if calcStakingInfo == nil {
		if stakingManager.migrationCheck.InMigration() {
			logger.Error("YOU MUST NOT RESTART NODE", "action", "if you want to restart the node, you should wait a day after the migration completes", "reason", "failed to store staking info while migration")
		}
		logger.Error("failed to update stakingInfo", "blockNum", blockNum, "staking block number", stakingBlockNumber)

		return nil
	}

	logger.Debug("Get stakingInfo from header.", "blockNum", blockNum, "staking block number", stakingBlockNumber, "stakingInfo", calcStakingInfo)
	return calcStakingInfo
}

// updateStakingInfo updates staking info in cache and db created from given block number.
func updateStakingInfo(blockNum uint64) (*StakingInfo, error) {
	if stakingManager == nil {
		return nil, ErrStakingManagerNotSet
	}

	stakingBlockNumber := params.CalcStakingBlockNumber(blockNum)

	stakingInfo, err := stakingManager.addressBookConnector.getStakingInfoFromAddressBook(stakingBlockNumber)
	if err != nil {
		return nil, err
	}

	stakingManager.stakingInfoCache.add(stakingInfo)

	if err := addStakingInfoToDB(stakingInfo); err != nil {
		logger.Debug("failed to write staking info to db", "err", err, "stakingInfo", stakingInfo)
		return stakingInfo, err
	}

	logger.Info("Add a new stakingInfo to stakingInfoCache and stakingInfoDB", "blockNum", stakingBlockNumber)
	logger.Debug("Added stakingInfo", "stakingInfo", stakingInfo)
	return stakingInfo, nil
}

// StakingManagerSubscribe setups a channel to listen chain head event and starts a goroutine to update staking cache.
func StakingManagerSubscribe() {
	if stakingManager == nil {
		logger.Warn("unable to subscribe; this can slow down node", "err", ErrStakingManagerNotSet)
		return
	}

	stakingManager.chainHeadSub = stakingManager.blockchain.SubscribeChainHeadEvent(stakingManager.chainHeadChan)

	go handleChainHeadEvent()
}

func handleChainHeadEvent() {
	if stakingManager == nil {
		logger.Warn("unable to start chain head event", "err", ErrStakingManagerNotSet)
		return
	} else if stakingManager.chainHeadSub == nil {
		logger.Info("unable to start chain head event", "err", ErrChainHeadChanNotSet)
		return
	}

	defer StakingManagerUnsubscribe()

	logger.Info("Start listening chain head event to update stakingInfoCache.")

	for {
		// A real event arrived, process interesting content
		select {
		// Handle ChainHeadEvent
		case ev := <-stakingManager.chainHeadChan:
			if stakingManager.governanceHelper.ProposerPolicy() == params.WeightedRandom {
				// check and update if staking info is not valid before for the next update interval blocks
				stakingInfo := GetStakingInfo(ev.Block.NumberU64() + params.StakingUpdateInterval())
				if stakingInfo == nil {
					logger.Error("unable to fetch staking info", "blockNum", ev.Block.NumberU64())
				}
			}
		case <-stakingManager.chainHeadSub.Err():
			return
		}
	}
}

// StakingManagerUnsubscribe can unsubscribe a subscription on chain head event.
func StakingManagerUnsubscribe() {
	if stakingManager == nil {
		logger.Warn("unable to start chain head event", "err", ErrStakingManagerNotSet)
		return
	} else if stakingManager.chainHeadSub == nil {
		logger.Info("unable to start chain head event", "err", ErrChainHeadChanNotSet)
		return
	}

	stakingManager.chainHeadSub.Unsubscribe()
}
