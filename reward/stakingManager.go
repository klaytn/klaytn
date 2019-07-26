package reward

import (
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/event"
	"github.com/klaytn/klaytn/params"
)

const (
	chainHeadChanSize = 100
)

type stakingManager struct {
	abm          *addressBookManager
	sic          *stakingInfoCache
	gh           governanceHelper
	bc           *blockchain.BlockChain
	chainHeadCh  chan blockchain.ChainHeadEvent
	chainHeadSub event.Subscription
}

func newStakingManager(bc *blockchain.BlockChain, gh governanceHelper) *stakingManager {
	abm := newAddressBookManager(bc, gh)
	sic := newStakingInfoCache()

	return &stakingManager{
		abm:         abm,
		sic:         sic,
		gh:          gh,
		bc:          bc,
		chainHeadCh: make(chan blockchain.ChainHeadEvent, chainHeadChanSize),
	}
}

// GetStakingInfoFromStakingCache returns a corresponding stakingInfo for a blockNum.
func (sm *stakingManager) getStakingInfo(blockNum uint64) *StakingInfo {
	stakingBlockNumber := params.CalcStakingBlockNumber(blockNum)

	if cachedStakingInfo := sm.sic.get(stakingBlockNumber); cachedStakingInfo != nil {
		logger.Debug("StakingInfoCache hit.", "Block number", blockNum, "number of staking block", stakingBlockNumber)
		return cachedStakingInfo
	}

	stakingInfo, err := sm.updateStakingCache(stakingBlockNumber)
	if err != nil {
		logger.Error("Failed to get stakingInfo", "Block number", blockNum, "number of staking block", stakingBlockNumber, "err", err)
		return nil
	}

	logger.Debug("Complete StakingInfoCache update.", "Block number", blockNum, "number of staking block", stakingBlockNumber)
	return stakingInfo
}

// updateStakingCache updates staking cache with a stakingInfo of a given block number.
func (sm *stakingManager) updateStakingCache(blockNum uint64) (*StakingInfo, error) {
	stakingInfo, err := sm.abm.getStakingInfoFromAddressBook(blockNum)
	if err != nil {
		return nil, err
	}

	sm.sic.add(stakingInfo)

	logger.Info("Add a new stakingInfo to the stakingInfoCache", "stakingInfo", stakingInfo)
	return stakingInfo, nil
}

// Subscribe setups a channel to listen chain head event and starts a goroutine to update staking cache.
func (sm *stakingManager) subscribe() {
	sm.chainHeadSub = sm.bc.SubscribeChainHeadEvent(sm.chainHeadCh)

	go sm.handleChainHeadEvent()
}

func (sm *stakingManager) handleChainHeadEvent() {
	defer sm.unsubscribe()

	logger.Info("Start listening chain head event to update stakingInfoCache.")

	for {
		// A real event arrived, process interesting content
		select {
		// Handle ChainHeadEvent
		case ev := <-sm.chainHeadCh:
			if sm.gh.ProposerPolicy() == params.WeightedRandom {
				stakingBlockNum := ev.Block.NumberU64() - ev.Block.NumberU64()%sm.gh.StakingUpdateInterval()
				if cachedStakingInfo := sm.sic.get(stakingBlockNum); cachedStakingInfo == nil {
					_, err := sm.updateStakingCache(stakingBlockNum)
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

func (sm *stakingManager) unsubscribe() {
	sm.chainHeadSub.Unsubscribe()
}
