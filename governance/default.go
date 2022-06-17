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

package governance

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/pkg/errors"
)

var (
	ErrValueTypeMismatch  = errors.New("Value's type mismatch")
	ErrDecodeGovChange    = errors.New("Failed to decode received governance changes")
	ErrUnmarshalGovChange = errors.New("Failed to unmarshal received governance changes")
	ErrVoteValueMismatch  = errors.New("Received change mismatches with the value this node has!!")
	ErrNotInitialized     = errors.New("Cache not initialized")
	ErrItemNotFound       = errors.New("Failed to find governance item")
	ErrItemNil            = errors.New("Governance Item is nil")
	ErrUnknownKey         = errors.New("Governnace value of the given key not found")
)

var (
	GovernanceKeyMap = map[string]int{
		"governance.governancemode":     params.GovernanceMode,
		"governance.governingnode":      params.GoverningNode,
		"istanbul.epoch":                params.Epoch,
		"istanbul.policy":               params.Policy,
		"istanbul.committeesize":        params.CommitteeSize,
		"governance.unitprice":          params.UnitPrice,
		"reward.mintingamount":          params.MintingAmount,
		"reward.ratio":                  params.Ratio,
		"reward.useginicoeff":           params.UseGiniCoeff,
		"reward.deferredtxfee":          params.DeferredTxFee,
		"reward.minimumstake":           params.MinimumStake,
		"reward.stakingupdateinterval":  params.StakeUpdateInterval,
		"reward.proposerupdateinterval": params.ProposerRefreshInterval,
		"governance.addvalidator":       params.AddValidator,
		"governance.removevalidator":    params.RemoveValidator,
		"param.txgashumanreadable":      params.ConstTxGasHumanReadable,
		"istanbul.timeout":              params.Timeout,
	}

	GovernanceForbiddenKeyMap = map[string]int{
		"istanbul.policy":               params.Policy,
		"reward.stakingupdateinterval":  params.StakeUpdateInterval,
		"reward.proposerupdateinterval": params.ProposerRefreshInterval,
	}

	GovernanceKeyMapReverse = map[int]string{
		params.GovernanceMode:          "governance.governancemode",
		params.GoverningNode:           "governance.governingnode",
		params.Epoch:                   "istanbul.epoch",
		params.CliqueEpoch:             "clique.epoch",
		params.Policy:                  "istanbul.policy",
		params.CommitteeSize:           "istanbul.committeesize",
		params.UnitPrice:               "governance.unitprice",
		params.MintingAmount:           "reward.mintingamount",
		params.Ratio:                   "reward.ratio",
		params.UseGiniCoeff:            "reward.useginicoeff",
		params.DeferredTxFee:           "reward.deferredtxfee",
		params.MinimumStake:            "reward.minimumstake",
		params.StakeUpdateInterval:     "reward.stakingupdateinterval",
		params.ProposerRefreshInterval: "reward.proposerupdateinterval",
		params.AddValidator:            "governance.addvalidator",
		params.RemoveValidator:         "governance.removevalidator",
		params.ConstTxGasHumanReadable: "param.txgashumanreadable",
		params.Timeout:                 "istanbul.timeout",
	}

	ProposerPolicyMap = map[string]int{
		"roundrobin":     params.RoundRobin,
		"sticky":         params.Sticky,
		"weightedrandom": params.WeightedRandom,
	}

	ProposerPolicyMapReverse = map[int]string{
		params.RoundRobin:     "roundrobin",
		params.Sticky:         "sticky",
		params.WeightedRandom: "weightedrandom",
	}

	GovernanceModeMap = map[string]int{
		"none":   params.GovernanceMode_None,
		"single": params.GovernanceMode_Single,
		"ballot": params.GovernanceMode_Ballot,
	}
)

var logger = log.NewModuleLogger(log.Governance)

// Governance item set
type GovernanceSet struct {
	items map[string]interface{}
	mu    *sync.RWMutex
}

// Governance represents vote information given from istanbul.vote()
type GovernanceVote struct {
	Validator common.Address `json:"validator"`
	Key       string         `json:"key"`
	Value     interface{}    `json:"value"`
}

// GovernanceTallies represents a tally for each governance item
type GovernanceTallyItem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Votes uint64      `json:"votes"`
}

type GovernanceTallyList struct {
	items []GovernanceTallyItem
	mu    *sync.RWMutex
}
type GovernanceVotes struct {
	items []GovernanceVote
	mu    *sync.RWMutex
}

type VoteStatus struct {
	Value  interface{} `json:"value"`
	Casted bool        `json:"casted"`
	Num    uint64      `json:"num"`
}

type VoteMap struct {
	items map[string]VoteStatus
	mu    *sync.RWMutex
}

// txPool is an interface for blockchain.TxPool used in governance package.
type txPool interface {
	SetGasPrice(price *big.Int)
}

// blockChain is an interface for blockchain.Blockchain used in governance package.
type blockChain interface {
	CurrentHeader() *types.Header
	SetProposerPolicy(val uint64)
	SetUseGiniCoeff(val bool)
}

type Governance struct {
	ChainConfig *params.ChainConfig

	// Map used to keep multiple types of votes
	voteMap VoteMap

	nodeAddress      atomic.Value //common.Address
	totalVotingPower uint64
	votingPower      uint64

	GovernanceVotes   GovernanceVotes
	GovernanceTallies GovernanceTallyList

	db           database.DBManager
	itemCache    common.Cache
	idxCache     []uint64 // elements should be in ascending order
	idxCacheLock *sync.RWMutex

	// The block number when current governance information was changed
	actualGovernanceBlock atomic.Value //uint64

	// The last block number at governance state was stored (used not to replay old votes)
	lastGovernanceStateBlock uint64

	currentSet GovernanceSet
	changeSet  GovernanceSet

	TxPool txPool

	blockChain blockChain
}

func NewVoteMap() VoteMap {
	return VoteMap{
		items: make(map[string]VoteStatus),
		mu:    new(sync.RWMutex),
	}
}

func NewGovernanceTallies() GovernanceTallyList {
	return GovernanceTallyList{
		items: []GovernanceTallyItem{},
		mu:    new(sync.RWMutex),
	}
}

func NewGovernanceVotes() GovernanceVotes {
	return GovernanceVotes{
		items: []GovernanceVote{},
		mu:    new(sync.RWMutex),
	}
}

func (gt *GovernanceTallyList) Clear() {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	gt.items = make([]GovernanceTallyItem, 0)
}

func (vl *VoteMap) Copy() map[string]VoteStatus {
	vl.mu.RLock()
	defer vl.mu.RUnlock()

	ret := make(map[string]VoteStatus)
	for k, v := range vl.items {
		ret[k] = v
	}

	return ret
}

func (vl *VoteMap) GetValue(key string) VoteStatus {
	vl.mu.RLock()
	defer vl.mu.RUnlock()

	return vl.items[key]
}

func (vl *VoteMap) SetValue(key string, val VoteStatus) {
	vl.mu.Lock()
	defer vl.mu.Unlock()

	vl.items[key] = val
}

func (vl *VoteMap) Import(src map[string]VoteStatus) {
	vl.mu.Lock()
	defer vl.mu.Unlock()

	for k, v := range src {
		vl.items[k] = v
	}
}

func (vl *VoteMap) Clear() {
	vl.mu.Lock()
	defer vl.mu.Unlock()

	// TODO-Governance if vote is not casted, it can remain forever. So, it would be better to add expiration.
	newItems := make(map[string]VoteStatus)
	for k, v := range vl.items {
		if !v.Casted {
			newItems[k] = v
		}
	}
	vl.items = newItems
}

func (vl *VoteMap) Size() int {
	vl.mu.RLock()
	defer vl.mu.RUnlock()

	return len(vl.items)
}

func (gt *GovernanceTallyList) Copy() []GovernanceTallyItem {
	gt.mu.RLock()
	defer gt.mu.RUnlock()

	ret := make([]GovernanceTallyItem, len(gt.items))
	copy(ret, gt.items)

	return ret
}

func (gt *GovernanceTallyList) Import(src []GovernanceTallyItem) {
	gt.mu.Lock()
	defer gt.mu.Unlock()

	gt.items = make([]GovernanceTallyItem, len(src))
	copy(gt.items, src)
}

func (gv *GovernanceVotes) Clear() {
	gv.mu.Lock()
	defer gv.mu.Unlock()
	gv.items = make([]GovernanceVote, 0)
}

func (gv *GovernanceVotes) Copy() []GovernanceVote {
	gv.mu.RLock()
	defer gv.mu.RUnlock()

	ret := make([]GovernanceVote, len(gv.items))
	copy(ret, gv.items)

	return ret
}

func (gv *GovernanceVotes) Import(src []GovernanceVote) {
	gv.mu.Lock()
	defer gv.mu.Unlock()

	gv.items = make([]GovernanceVote, len(src))
	copy(gv.items, src)
}

func NewGovernanceSet() GovernanceSet {
	return GovernanceSet{
		items: map[string]interface{}{},
		mu:    new(sync.RWMutex),
	}
}

func (gs *GovernanceSet) Clear() {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.items = make(map[string]interface{})
}

func (gs *GovernanceSet) SetValue(itemType int, value interface{}) error {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	key := GovernanceKeyMapReverse[itemType]
	if !checkValueType(value, GovernanceItems[itemType].t) {
		return ErrValueTypeMismatch
	}
	gs.items[key] = value
	return nil
}

func (gs *GovernanceSet) GetValue(key int) (interface{}, bool) {
	sKey, ok := GovernanceKeyMapReverse[key]
	if !ok {
		return nil, false
	}

	gs.mu.RLock()
	defer gs.mu.RUnlock()
	ret, ok := gs.items[sKey]
	return ret, ok
}

func (gs *GovernanceSet) RemoveItem(key string) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	delete(gs.items, key)
}

func (gs *GovernanceSet) Size() int {
	gs.mu.RLock()
	defer gs.mu.RUnlock()

	return len(gs.items)
}

func (gs *GovernanceSet) Import(src map[string]interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	gs.items = make(map[string]interface{})
	for k, v := range src {
		gs.items[k] = v
	}
}

func (gs *GovernanceSet) Items() map[string]interface{} {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	ret := make(map[string]interface{})
	for k, v := range gs.items {
		ret[k] = v
	}
	return ret
}

func (gs *GovernanceSet) Merge(change map[string]interface{}) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	for k, v := range change {
		gs.items[k] = v
	}
}

// NewGovernance creates Governance with the given configuration.
func NewGovernance(chainConfig *params.ChainConfig, dbm database.DBManager) *Governance {
	return &Governance{
		ChainConfig:              chainConfig,
		voteMap:                  NewVoteMap(),
		db:                       dbm,
		itemCache:                newGovernanceCache(),
		currentSet:               NewGovernanceSet(),
		changeSet:                NewGovernanceSet(),
		lastGovernanceStateBlock: 0,
		GovernanceTallies:        NewGovernanceTallies(),
		GovernanceVotes:          NewGovernanceVotes(),
		idxCacheLock:             new(sync.RWMutex),
	}
}

// NewGovernanceInitialize creates Governance with the given configuration and read governance state from DB.
// If any items are not stored in DB, it stores governance items of the genesis block to DB.
func NewGovernanceInitialize(chainConfig *params.ChainConfig, dbm database.DBManager) *Governance {
	ret := NewGovernance(chainConfig, dbm)
	// nil is for testing or simple function usage
	if dbm != nil {
		ret.ReadGovernanceState()
		if err := ret.initializeCache(); err != nil {
			// If this is the first time to run, store governance information for genesis block on database
			cfg := GetGovernanceItemsFromChainConfig(chainConfig)
			if err := ret.WriteGovernance(0, cfg, NewGovernanceSet()); err != nil {
				logger.Crit("Error in writing governance information", "err", err)
			}
			// If failed again after writing governance, stop booting up
			if err = ret.initializeCache(); err != nil {
				logger.Crit("No governance cache index found in a database", "err", err)
			}
		}
	}
	return ret
}

func (g *Governance) updateGovernanceParams() {
	params.SetStakingUpdateInterval(g.StakingUpdateInterval())
	params.SetProposerUpdateInterval(g.ProposerUpdateInterval())

	// NOTE: HumanReadable related functions are inactivated now
	if txGasHumanReadable, ok := g.currentSet.GetValue(params.ConstTxGasHumanReadable); ok {
		params.TxGasHumanReadable = txGasHumanReadable.(uint64)
	}
}

func (g *Governance) SetNodeAddress(addr common.Address) {
	g.nodeAddress.Store(addr)
}

func (g *Governance) SetTotalVotingPower(t uint64) {
	atomic.StoreUint64(&g.totalVotingPower, t)
}

func (g *Governance) SetMyVotingPower(t uint64) {
	atomic.StoreUint64(&g.votingPower, t)
}

func (g *Governance) NodeAddress() common.Address {
	return g.nodeAddress.Load().(common.Address)
}

func (g *Governance) TotalVotingPower() uint64 {
	return atomic.LoadUint64(&g.totalVotingPower)
}

func (g *Governance) MyVotingPower() uint64 {
	return atomic.LoadUint64(&g.votingPower)
}

func (gov *Governance) BlockChain() blockChain {
	return gov.blockChain
}

func (gov *Governance) DB() database.DBManager {
	return gov.db
}

func (g *Governance) GetEncodedVote(addr common.Address, number uint64) []byte {
	// TODO-Klaytn-Governance Change this part to add all votes to the header at once
	for key, val := range g.voteMap.Copy() {
		if val.Casted == false {
			vote := new(GovernanceVote)
			vote.Validator = addr
			vote.Key = key
			vote.Value = val.Value
			encoded, err := rlp.EncodeToBytes(vote)
			if err != nil {
				logger.Error("Failed to RLP Encode a vote", "vote", vote)
				g.RemoveVote(key, val, number)
				continue
			}
			return encoded
		}
	}
	return nil
}

func (g *Governance) getKey(k string) string {
	return strings.Trim(strings.ToLower(k), " ")
}

// RemoveVote remove a vote from the voteMap to prevent repetitive addition of same vote
func (g *Governance) RemoveVote(key string, value interface{}, number uint64) {
	k := GovernanceKeyMap[key]
	if isEqualValue(k, g.voteMap.GetValue(key).Value, value) {
		g.voteMap.SetValue(key, VoteStatus{
			Value:  value,
			Casted: true,
			Num:    number,
		})
	}
	if g.CanWriteGovernanceState(number) {
		g.WriteGovernanceState(number, false)
	}
}

func (g *Governance) ClearVotes(num uint64) {
	g.GovernanceVotes.Clear()
	g.GovernanceTallies.Clear()
	g.changeSet.Clear()
	g.voteMap.Clear()
	logger.Info("Governance votes are cleared", "num", num)
}

// parseVoteValue parse vote.Value from []uint8, [][]uint8 to appropriate type
func (g *Governance) ParseVoteValue(gVote *GovernanceVote) (*GovernanceVote, error) {
	var val interface{}
	k, ok := GovernanceKeyMap[gVote.Key]
	if !ok {
		logger.Warn("Unknown key was given", "key", k)
		return nil, ErrUnknownKey
	}

	switch k {
	case params.GovernanceMode, params.MintingAmount, params.MinimumStake, params.Ratio:
		v, ok := gVote.Value.([]uint8)
		if !ok {
			return nil, ErrValueTypeMismatch
		}
		val = string(v)
	case params.GoverningNode:
		v, ok := gVote.Value.([]uint8)
		if !ok {
			return nil, ErrValueTypeMismatch
		}
		val = common.BytesToAddress(v)
	case params.AddValidator, params.RemoveValidator:
		if v, ok := gVote.Value.([]uint8); ok {
			// if value contains single address, gVote.Value type should be []uint8{}
			val = common.BytesToAddress(v)
		} else if addresses, ok := gVote.Value.([]interface{}); ok {
			// if value contains multiple addresses, gVote.Value type should be [][]uint8{}
			if len(addresses) == 0 {
				return nil, ErrValueTypeMismatch
			}
			var nodeAddresses []common.Address
			for _, item := range addresses {
				if in, ok := item.([]uint8); !ok || len(in) != common.AddressLength {
					return nil, ErrValueTypeMismatch
				}
				nodeAddresses = append(nodeAddresses, common.BytesToAddress(item.([]uint8)))
			}
			val = nodeAddresses
		} else {
			return nil, ErrValueTypeMismatch
		}
	case params.Epoch, params.CommitteeSize, params.UnitPrice, params.StakeUpdateInterval, params.ProposerRefreshInterval, params.ConstTxGasHumanReadable, params.Policy, params.Timeout:
		v, ok := gVote.Value.([]uint8)
		if !ok {
			return nil, ErrValueTypeMismatch
		}
		v = append(make([]byte, 8-len(v)), v...)
		val = binary.BigEndian.Uint64(v)
	case params.UseGiniCoeff, params.DeferredTxFee:
		v, ok := gVote.Value.([]uint8)
		if !ok {
			return nil, ErrValueTypeMismatch
		}
		v = append(make([]byte, 8-len(v)), v...)
		if binary.BigEndian.Uint64(v) != uint64(0) {
			val = true
		} else {
			val = false
		}
	}
	gVote.Value = val
	return gVote, nil
}

func (gov *Governance) ReflectVotes(vote GovernanceVote) {
	if ok := gov.updateChangeSet(vote); !ok {
		logger.Error("Failed to reflect Governance Config", "Key", vote.Key, "Value", vote.Value)
	}
}

func (gov *Governance) updateChangeSet(vote GovernanceVote) bool {
	switch GovernanceKeyMap[vote.Key] {
	case params.GoverningNode:
		gov.changeSet.SetValue(GovernanceKeyMap[vote.Key], vote.Value.(common.Address))
		return true
	case params.GovernanceMode, params.Ratio:
		gov.changeSet.SetValue(GovernanceKeyMap[vote.Key], vote.Value.(string))
		return true
	case params.Epoch, params.StakeUpdateInterval, params.ProposerRefreshInterval, params.CommitteeSize, params.UnitPrice, params.ConstTxGasHumanReadable, params.Policy, params.Timeout:
		gov.changeSet.SetValue(GovernanceKeyMap[vote.Key], vote.Value.(uint64))
		return true
	case params.MintingAmount, params.MinimumStake:
		gov.changeSet.SetValue(GovernanceKeyMap[vote.Key], vote.Value.(string))
		return true
	case params.UseGiniCoeff, params.DeferredTxFee:
		gov.changeSet.SetValue(GovernanceKeyMap[vote.Key], vote.Value.(bool))
		return true
	default:
		logger.Warn("Unknown key was given", "key", vote.Key)
	}
	return false
}

func CheckGenesisValues(c *params.ChainConfig) error {
	gov := NewGovernanceInitialize(c, nil)

	var tstMap = map[string]interface{}{
		"istanbul.epoch":                c.Istanbul.Epoch,
		"istanbul.committeesize":        c.Istanbul.SubGroupSize,
		"istanbul.policy":               uint64(c.Istanbul.ProposerPolicy),
		"governance.governancemode":     c.Governance.GovernanceMode,
		"governance.governingnode":      c.Governance.GoverningNode,
		"governance.unitprice":          c.UnitPrice,
		"reward.ratio":                  c.Governance.Reward.Ratio,
		"reward.useginicoeff":           c.Governance.Reward.UseGiniCoeff,
		"reward.deferredtxfee":          c.Governance.Reward.DeferredTxFee,
		"reward.mintingamount":          c.Governance.Reward.MintingAmount.String(),
		"reward.minimumstake":           c.Governance.Reward.MinimumStake.String(),
		"reward.stakingupdateinterval":  c.Governance.Reward.StakingUpdateInterval,
		"reward.proposerupdateinterval": c.Governance.Reward.ProposerUpdateInterval,
	}

	for k, v := range tstMap {
		if _, ok := gov.ValidateVote(&GovernanceVote{Key: k, Value: v}); !ok {
			return errors.New(k + " value is wrong")
		}
	}
	return nil
}

func newGovernanceCache() common.Cache {
	cache := common.NewCache(common.LRUConfig{CacheSize: params.GovernanceCacheLimit})
	return cache
}

// initializeCache reads governance item data from database and updates Governance.itemCache.
// It also initializes currentSet and actualGovernanceBlock according to head block number.
func (g *Governance) initializeCache() error {
	// get last n governance change block number
	indices, err := g.db.ReadRecentGovernanceIdx(params.GovernanceIdxCacheLimit)
	if err != nil {
		return ErrNotInitialized
	}
	g.idxCache = indices

	// Put governance items into the itemCache
	for _, v := range indices {
		if data, err := g.db.ReadGovernance(v); err == nil {
			data = adjustDecodedSet(data)
			g.itemCache.Add(getGovernanceCacheKey(v), data)
		} else {
			logger.Crit("Couldn't read governance cache from database. Check database consistency", "index", v, "err", err)
		}
	}

	// g.db.ReadGovernance(N) returns the governance data stored in block 'N'. If the data is not exists, it simply returns an error.
	// On the other hand, g.ReadGovernance(N) returns the governance data viewed by block 'N' and the block number at the time the governance data is stored.
	// Since g.currentSet means the governance data viewed by the head block number of the node
	// and g.actualGovernanceBlock means the block number at the time the governance data is stored,
	// So those two variables are initialized by using the return values of the g.ReadGovernance(headBlockNum).

	// head block number is used to get the appropriate g.currentSet and g.actualGovernanceBlock
	headBlockNumber := uint64(0)
	if headBlockHash := g.db.ReadHeadBlockHash(); !common.EmptyHash(headBlockHash) {
		if num := g.db.ReadHeaderNumber(headBlockHash); num != nil {
			headBlockNumber = *num
		}
	}
	newBlockNumber, newGovernanceSet, err := g.ReadGovernance(headBlockNumber)
	if err != nil {
		return err
	}
	// g.actualGovernanceBlock and currentSet is set
	g.actualGovernanceBlock.Store(newBlockNumber)
	g.currentSet.Import(newGovernanceSet)
	g.updateGovernanceParams()

	// g.lastGovernanceStateBlock contains the last block number when voting data is included.
	// we check the order between g.actualGovernanceBlock and g.lastGovernanceStateBlock, so make sure that voting is not missed.
	governanceBlock, governanceStateBlock := g.actualGovernanceBlock.Load().(uint64), atomic.LoadUint64(&g.lastGovernanceStateBlock)
	if governanceBlock >= governanceStateBlock {
		ret, ok := g.itemCache.Get(getGovernanceCacheKey(governanceBlock))
		if !ok || ret == nil {
			logger.Error("cannot get governance data at actualGovernanceBlock", "actualGovernanceBlock", governanceBlock)
			return errors.New("Currentset initialization failed")
		}
		g.currentSet.Import(ret.(map[string]interface{}))
		g.updateGovernanceParams()
	}

	return nil
}

// getGovernanceCache returns cached governance config as a byte slice
func (g *Governance) getGovernanceCache(num uint64) (map[string]interface{}, bool) {
	cKey := getGovernanceCacheKey(num)

	if ret, ok := g.itemCache.Get(cKey); ok && ret != nil {
		return ret.(map[string]interface{}), true
	}
	return nil, false
}

func (g *Governance) addGovernanceCache(num uint64, data GovernanceSet) {
	// Don't update cache if num (block number) is smaller than the biggest number of cached block number
	g.idxCacheLock.Lock()
	defer g.idxCacheLock.Unlock()

	if len(g.idxCache) > 0 && num <= g.idxCache[len(g.idxCache)-1] {
		logger.Error("The same or more recent governance index exist. Skip updating governance cache",
			"newIdx", num, "govIdxes", g.idxCache)
		return
	}
	cKey := getGovernanceCacheKey(num)
	g.itemCache.Add(cKey, data.Items())

	g.idxCache = append(g.idxCache, num)
	if len(g.idxCache) > params.GovernanceIdxCacheLimit {
		g.idxCache = g.idxCache[len(g.idxCache)-params.GovernanceIdxCacheLimit:]
	}
}

// getGovernanceCacheKey returns cache key of the given block number
func getGovernanceCacheKey(num uint64) common.GovernanceCacheKey {
	v := fmt.Sprintf("%v", num)
	return common.GovernanceCacheKey(params.GovernanceCachePrefix + "_" + v)
}

// Store new governance data on DB. This updates Governance cache too.
func (g *Governance) WriteGovernance(num uint64, data GovernanceSet, delta GovernanceSet) error {
	g.idxCacheLock.RLock()
	indices := make([]uint64, len(g.idxCache))
	copy(indices, g.idxCache)
	g.idxCacheLock.RUnlock()

	if len(indices) > 0 && num <= indices[len(indices)-1] {
		logger.Error("The same or more recent governance index exist. Skip writing governance",
			"newIdx", num, "govIdxes", indices)
		return nil
	}

	new := NewGovernanceSet()
	new.Import(data.Items())

	// merge delta into data
	if delta.Size() > 0 {
		new.Merge(delta.Items())
	}
	g.addGovernanceCache(num, new)
	return g.db.WriteGovernance(new.Items(), num)
}

func (g *Governance) searchCache(num uint64) (uint64, bool) {
	g.idxCacheLock.RLock()
	defer g.idxCacheLock.RUnlock()

	for i := len(g.idxCache) - 1; i >= 0; i-- {
		if g.idxCache[i] <= num {
			return g.idxCache[i], true
		}
	}
	return 0, false
}

func (g *Governance) ReadGovernance(num uint64) (uint64, map[string]interface{}, error) {
	if g.ChainConfig.Istanbul == nil {
		logger.Crit("Failed to read governance. ChainConfig.Istanbul == nil")
	}
	blockNum := CalcGovernanceInfoBlock(num, g.Epoch())
	// Check cache first
	if gBlockNum, ok := g.searchCache(blockNum); ok {
		if data, okay := g.getGovernanceCache(gBlockNum); okay {
			return gBlockNum, data, nil
		}
	}
	if g.db != nil {
		bn, result, err := g.db.ReadGovernanceAtNumber(num, g.Epoch())
		result = adjustDecodedSet(result)
		return bn, result, err
	} else {
		// For CI tests which don't have a database
		return 0, nil, nil
	}
}

func CalcGovernanceInfoBlock(num uint64, epoch uint64) uint64 {
	governanceInfoBlock := num - (num % epoch)
	if governanceInfoBlock >= epoch {
		governanceInfoBlock -= epoch
	}
	return governanceInfoBlock
}

func (g *Governance) GetGovernanceChange() map[string]interface{} {
	if g.changeSet.Size() > 0 {
		return g.changeSet.Items()
	}
	return nil
}

// WriteGovernanceForNextEpoch creates governance items for next epoch and writes them to the database.
// The governance items on next epoch will be the given `governance` items applied on the top of past epoch items.
func (gov *Governance) WriteGovernanceForNextEpoch(number uint64, governance []byte) {
	var epoch uint64
	var ok bool

	if epoch, ok = gov.GetGovernanceValue(params.Epoch).(uint64); !ok {
		if epoch, ok = gov.GetGovernanceValue(params.CliqueEpoch).(uint64); !ok {
			logger.Error("Couldn't find epoch from governance items")
			return
		}
	}

	// Store updated governance information if exist
	if number%epoch == 0 {
		if len(governance) > 0 {
			tempData := []byte("")
			tempItems := make(map[string]interface{})
			tempSet := NewGovernanceSet()
			if err := rlp.DecodeBytes(governance, &tempData); err != nil {
				logger.Error("Failed to decode governance data", "number", number, "err", err, "data", governance)
				return
			}
			if err := json.Unmarshal(tempData, &tempItems); err != nil {
				logger.Error("Failed to unmarshal governance data", "number", number, "err", err, "data", tempData)
				return

			}
			tempItems = adjustDecodedSet(tempItems)
			tempSet.Import(tempItems)

			_, govItems, err := gov.ReadGovernance(number)
			if err != nil {
				logger.Error("Failed to read governance", "number", number, "err", err)
				return
			}
			govSet := NewGovernanceSet()
			govSet.Import(govItems)

			// Store new governance items for next epoch to governance database
			if err := gov.WriteGovernance(number, govSet, tempSet); err != nil {
				logger.Crit("Failed to store new governance data", "number", number, "err", err)
			}
		}
	}
}

func (gov *Governance) removeDuplicatedVote(vote *GovernanceVote, number uint64) {
	gov.RemoveVote(vote.Key, vote.Value, number)
}

func (gov *Governance) UpdateCurrentSet(num uint64) {
	newNumber, newGovernanceSet, _ := gov.ReadGovernance(num)
	// Do the change only when the governance actually changed
	if newGovernanceSet != nil && newNumber > gov.actualGovernanceBlock.Load().(uint64) {
		gov.actualGovernanceBlock.Store(newNumber)
		gov.currentSet.Import(newGovernanceSet)
		gov.triggerChange(newGovernanceSet)
	}
}

func (gov *Governance) triggerChange(src map[string]interface{}) {
	for k, v := range src {
		if f := GovernanceItems[GovernanceKeyMap[k]].trigger; f != nil {
			f(gov, k, v)
		}
	}
}

func adjustDecodedSet(src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		x := reflect.ValueOf(v)
		if x.Kind() == reflect.Float64 {
			src[k] = uint64(v.(float64))
		}
		if GovernanceKeyMap[k] == params.GoverningNode {
			if reflect.TypeOf(v) == stringT {
				src[k] = common.HexToAddress(v.(string))
			} else {
				src[k] = v
			}
		}
	}
	return src
}

func (gov *Governance) GetGovernanceValue(key int) interface{} {
	if v, ok := gov.currentSet.GetValue(key); !ok {
		return nil
	} else {
		return v
	}
}

func (gov *Governance) VerifyGovernance(received []byte) error {
	change := []byte{}
	if rlp.DecodeBytes(received, &change) != nil {
		return ErrDecodeGovChange
	}

	rChangeSet := make(map[string]interface{})
	if json.Unmarshal(change, &rChangeSet) != nil {
		return ErrUnmarshalGovChange
	}
	rChangeSet = adjustDecodedSet(rChangeSet)

	if len(rChangeSet) != gov.changeSet.Size() {
		logger.Error("Verification Error", "len(receivedChangeSet)", len(rChangeSet), "len(changeSet)", gov.changeSet.Size())
		return ErrVoteValueMismatch
	}

	for k, v := range rChangeSet {
		if GovernanceKeyMap[k] == params.GoverningNode {
			if reflect.TypeOf(v) == stringT {
				v = common.HexToAddress(v.(string))
			}
		}

		have, _ := gov.changeSet.GetValue(GovernanceKeyMap[k])
		if have != v {
			logger.Error("Verification Error", "key", k, "received", rChangeSet[k], "have", have, "receivedType", reflect.TypeOf(rChangeSet[k]), "haveType", reflect.TypeOf(have))
			return ErrVoteValueMismatch
		}
	}
	return nil
}

type governanceJSON struct {
	BlockNumber     uint64                 `json:"blockNumber"`
	ChainConfig     *params.ChainConfig    `json:"chainConfig"`
	VoteMap         map[string]VoteStatus  `json:"voteMap"`
	NodeAddress     common.Address         `json:"nodeAddress"`
	GovernanceVotes []GovernanceVote       `json:"governanceVotes"`
	GovernanceTally []GovernanceTallyItem  `json:"governanceTally"`
	CurrentSet      map[string]interface{} `json:"currentSet"`
	ChangeSet       map[string]interface{} `json:"changeSet"`
}

func (gov *Governance) toJSON(num uint64) ([]byte, error) {
	ret := &governanceJSON{
		BlockNumber:     num,
		ChainConfig:     gov.ChainConfig,
		VoteMap:         gov.voteMap.Copy(),
		NodeAddress:     gov.nodeAddress.Load().(common.Address),
		GovernanceVotes: gov.GovernanceVotes.Copy(),
		GovernanceTally: gov.GovernanceTallies.Copy(),
		CurrentSet:      gov.currentSet.Items(),
		ChangeSet:       gov.changeSet.Items(),
	}
	j, _ := json.Marshal(ret)
	return j, nil
}

func (gov *Governance) UnmarshalJSON(b []byte) error {
	var j governanceJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	gov.ChainConfig = j.ChainConfig
	gov.voteMap.Import(j.VoteMap)
	gov.nodeAddress.Store(j.NodeAddress)
	gov.GovernanceVotes.Import(j.GovernanceVotes)
	gov.GovernanceTallies.Import(j.GovernanceTally)
	gov.currentSet.Import(adjustDecodedSet(j.CurrentSet))
	gov.changeSet.Import(adjustDecodedSet(j.ChangeSet))
	atomic.StoreUint64(&gov.lastGovernanceStateBlock, j.BlockNumber)

	return nil
}

func (gov *Governance) CanWriteGovernanceState(num uint64) bool {
	if num <= atomic.LoadUint64(&gov.lastGovernanceStateBlock) {
		return false
	}
	return true
}

func (gov *Governance) WriteGovernanceState(num uint64, isCheckpoint bool) error {
	if b, err := gov.toJSON(num); err != nil {
		logger.Error("Error in marshaling governance state", "err", err)
		return err
	} else {
		if err = gov.db.WriteGovernanceState(b); err != nil {
			logger.Error("Error in writing governance state", "err", err)
			return err
		} else {
			if isCheckpoint {
				atomic.StoreUint64(&gov.lastGovernanceStateBlock, num)
			}
			logger.Info("Successfully stored governance state", "num", num)
			return nil
		}
	}
}

// ReadGovernanceState reads field values of the Governance struct from database.
// It also updates params.stakingUpdateInterval and params.proposerUpdateInterval with the retrieved value.
func (gov *Governance) ReadGovernanceState() {
	b, err := gov.db.ReadGovernanceState()
	if err != nil {
		logger.Info("No governance state found in a database")
		return
	}
	gov.UnmarshalJSON(b)
	gov.updateGovernanceParams()

	logger.Info("Successfully loaded governance state from database", "blockNumber", atomic.LoadUint64(&gov.lastGovernanceStateBlock))
}

func (gov *Governance) SetBlockchain(bc blockChain) {
	gov.blockChain = bc
}

func (gov *Governance) SetTxPool(txpool txPool) {
	gov.TxPool = txpool
}

func GetGovernanceItemsFromChainConfig(config *params.ChainConfig) GovernanceSet {
	g := NewGovernanceSet()

	if config.Governance != nil {
		governance := config.Governance
		governanceMap := map[int]interface{}{
			params.GovernanceMode:          governance.GovernanceMode,
			params.GoverningNode:           governance.GoverningNode,
			params.UnitPrice:               config.UnitPrice,
			params.MintingAmount:           governance.Reward.MintingAmount.String(),
			params.Ratio:                   governance.Reward.Ratio,
			params.UseGiniCoeff:            governance.Reward.UseGiniCoeff,
			params.DeferredTxFee:           governance.Reward.DeferredTxFee,
			params.MinimumStake:            governance.Reward.MinimumStake.String(),
			params.StakeUpdateInterval:     governance.Reward.StakingUpdateInterval,
			params.ProposerRefreshInterval: governance.Reward.ProposerUpdateInterval,
		}

		for k, v := range governanceMap {
			if err := g.SetValue(k, v); err != nil {
				writeFailLog(k, err)
			}
		}
	}

	if config.Istanbul != nil {
		istanbul := config.Istanbul
		istanbulMap := map[int]interface{}{
			params.Epoch:         istanbul.Epoch,
			params.Policy:        istanbul.ProposerPolicy,
			params.CommitteeSize: istanbul.SubGroupSize,
		}

		for k, v := range istanbulMap {
			if err := g.SetValue(k, v); err != nil {
				writeFailLog(k, err)
			}
		}
	}
	return g
}

func writeFailLog(key int, err error) {
	msg := "Failed to set " + GovernanceKeyMapReverse[key]
	logger.Crit(msg, "err", err)
}

func AddGovernanceCacheForTest(e HeaderEngine, num uint64, config *params.ChainConfig) {
	// addGovernanceCache only exists and relevant in *Governance.
	if g, ok := e.(*Governance); ok {
		data := GetGovernanceItemsFromChainConfig(config)
		g.addGovernanceCache(num, data)
	}
}

func (gov *Governance) GovernanceMode() string {
	return gov.GetGovernanceValue(params.GovernanceMode).(string)
}

func (gov *Governance) GoverningNode() common.Address {
	return gov.GetGovernanceValue(params.GoverningNode).(common.Address)
}

func (gov *Governance) UnitPrice() uint64 {
	return gov.GetGovernanceValue(params.UnitPrice).(uint64)
}

func (gov *Governance) CommitteeSize() uint64 {
	return gov.GetGovernanceValue(params.CommitteeSize).(uint64)
}

func (gov *Governance) Epoch() uint64 {
	if ret := gov.GetGovernanceValue(params.Epoch); ret == nil {
		// When a node is initializing, value can be nil
		return gov.ChainConfig.Istanbul.Epoch
	}
	return gov.GetGovernanceValue(params.Epoch).(uint64)
}

func (gov *Governance) ProposerPolicy() uint64 {
	return gov.GetGovernanceValue(params.Policy).(uint64)
}

func (gov *Governance) DeferredTxFee() bool {
	return gov.GetGovernanceValue(params.DeferredTxFee).(bool)
}

func (gov *Governance) MinimumStake() string {
	return gov.GetGovernanceValue(params.MinimumStake).(string)
}

func (gov *Governance) MintingAmount() string {
	return gov.GetGovernanceValue(params.MintingAmount).(string)
}

func (gov *Governance) ProposerUpdateInterval() uint64 {
	return gov.GetGovernanceValue(params.ProposerRefreshInterval).(uint64)
}

func (gov *Governance) Ratio() string {
	return gov.GetGovernanceValue(params.Ratio).(string)
}

func (gov *Governance) StakingUpdateInterval() uint64 {
	return gov.GetGovernanceValue(params.StakeUpdateInterval).(uint64)
}

func (gov *Governance) UseGiniCoeff() bool {
	return gov.GetGovernanceValue(params.UseGiniCoeff).(bool)
}

func (gov *Governance) ChainId() uint64 {
	return gov.ChainConfig.ChainID.Uint64()
}

func (gov *Governance) InitialChainConfig() *params.ChainConfig {
	return gov.ChainConfig
}

func (g *Governance) GetVoteMapCopy() map[string]VoteStatus {
	return g.voteMap.Copy()
}

func (g *Governance) GetGovernanceTalliesCopy() []GovernanceTallyItem {
	return g.GovernanceTallies.Copy()
}

func (gov *Governance) CurrentSetCopy() map[string]interface{} {
	return gov.currentSet.Items()
}

func (gov *Governance) PendingChanges() map[string]interface{} {
	return gov.changeSet.Items()
}

func (gov *Governance) Votes() []GovernanceVote {
	return gov.GovernanceVotes.Copy()
}

func (gov *Governance) IdxCache() []uint64 {
	gov.idxCacheLock.RLock()
	defer gov.idxCacheLock.RUnlock()

	copiedCache := make([]uint64, len(gov.idxCache))
	copy(copiedCache, gov.idxCache)
	return copiedCache
}

func (gov *Governance) IdxCacheFromDb() []uint64 {
	res, _ := gov.db.ReadRecentGovernanceIdx(0)
	return res
}
