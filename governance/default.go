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
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/ser/rlp"
	"github.com/klaytn/klaytn/storage/database"
	"github.com/pkg/errors"
	"math/big"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	ErrValueTypeMismatch  = errors.New("Value's type mismatch")
	ErrDecodeGovChange    = errors.New("Failed to decode received governance changes")
	ErrUnmarshalGovChange = errors.New("Failed to unmarshal received governance changes")
	ErrVoteValueMismatch  = errors.New("Received change mismatches with the value this node has!!")
	ErrNotInitialized     = errors.New("Cache not initialized")
	ErrItemNotFound       = errors.New("Failed to find governance item")
	ErrItemNil            = errors.New("Governance Item is nil")
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
	}

	GovernanceForbiddenKeyMap = map[string]int{
		"istanbul.policy":               params.Policy,
		"reward.stakingupdateinterval":  params.StakeUpdateInterval,
		"reward.proposerupdateinterval": params.ProposerRefreshInterval,

		// TODO-Klaytn-Issue3567 forbid voting except add, remove validator
		// bellow code is added only for temporarily until we solve data race of governance.
		// until then, voting except add, remove validator will be forbidden
		"governance.governancemode": params.GovernanceMode,
		"governance.governingnode":  params.GoverningNode,
		"istanbul.epoch":            params.Epoch,
		"istanbul.committeesize":    params.CommitteeSize,
		"governance.unitprice":      params.UnitPrice,
		"reward.mintingamount":      params.MintingAmount,
		"reward.ratio":              params.Ratio,
		"reward.useginicoeff":       params.UseGiniCoeff,
		"reward.deferredtxfee":      params.DeferredTxFee,
		"reward.minimumstake":       params.MinimumStake,
		"param.txgashumanreadable":  params.ConstTxGasHumanReadable,
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
type GovernanceSet map[string]interface{}

// Governance represents vote information given from istanbul.vote()
type GovernanceVote struct {
	Validator common.Address `json:"validator"`
	Key       string         `json:"key"`
	Value     interface{}    `json:"value"`
}

// GovernanceTally represents a tally for each governance item
type GovernanceTally struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Votes uint64      `json:"votes"`
}

type VoteStatus struct {
	Value  interface{} `json:"value"`
	Casted bool        `json:"casted"`
	Num    uint64      `json:"num"`
}

type Governance struct {
	ChainConfig    *params.ChainConfig
	governanceLock sync.RWMutex

	// Map used to keep multiple types of votes
	voteMap     map[string]VoteStatus
	voteMapLock sync.RWMutex

	nodeAddress      common.Address
	totalVotingPower uint64
	votingPower      uint64

	GovernanceVotes     []GovernanceVote
	GovernanceTally     []GovernanceTally
	GovernanceTallyLock sync.RWMutex

	db        database.DBManager
	itemCache common.Cache
	idxCache  []uint64

	// The block number when current governance information was changed
	actualGovernanceBlock uint64

	// The last block number at governance state was stored (used not to replay old votes)
	lastGovernanceStateBlock uint64

	currentSet   GovernanceSet
	currentSetMu sync.RWMutex

	changeSet GovernanceSet
	mu        sync.RWMutex

	TxPool *blockchain.TxPool

	blockChain *blockchain.BlockChain
}

func (gs GovernanceSet) SetValue(itemType int, value interface{}) error {
	key := GovernanceKeyMapReverse[itemType]

	if GovernanceItems[itemType].t != reflect.TypeOf(value) {
		return ErrValueTypeMismatch
	}
	gs[key] = value
	return nil
}

func NewGovernance(chainConfig *params.ChainConfig, dbm database.DBManager) *Governance {
	ret := Governance{
		ChainConfig:              chainConfig,
		voteMap:                  make(map[string]VoteStatus),
		db:                       dbm,
		itemCache:                newGovernanceCache(),
		currentSet:               GovernanceSet{},
		changeSet:                GovernanceSet{},
		lastGovernanceStateBlock: 0,
	}
	// nil is for testing or simple function usage
	if dbm != nil {
		if err := ret.initializeCache(); err != nil {
			// If this is the first time to run, store governance information for genesis block on database
			cfg := getGovernanceItemsFromChainConfig(chainConfig)
			if err := ret.WriteGovernance(0, cfg, nil); err != nil {
				logger.Crit("Error in writing governance information", "err", err)
			}
			// If failed again after writing governance, stop booting up
			if err = ret.initializeCache(); err != nil {
				logger.Crit("No governance cache index found in a database", "err", err)
			}
		}
		ret.ReadGovernanceState()
	}
	return &ret
}

func (g *Governance) SetNodeAddress(addr common.Address) {
	g.nodeAddress = addr
}

func (g *Governance) SetTotalVotingPower(t uint64) {
	atomic.StoreUint64(&g.totalVotingPower, t)
}

func (g *Governance) SetMyVotingPower(t uint64) {
	atomic.StoreUint64(&g.votingPower, t)
}

func (g *Governance) GetEncodedVote(addr common.Address, number uint64) []byte {
	// TODO-Klaytn-Governance Change this part to add all votes to the header at once
	g.voteMapLock.RLock()
	defer g.voteMapLock.RUnlock()

	if len(g.voteMap) > 0 {
		for key, val := range g.voteMap {
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
	}
	return nil
}

func (g *Governance) getKey(k string) string {
	return strings.Trim(strings.ToLower(k), " ")
}

// RemoveVote remove a vote from the voteMap to prevent repetitive addition of same vote
func (g *Governance) RemoveVote(key string, value interface{}, number uint64) {
	g.voteMapLock.Lock()
	defer g.voteMapLock.Unlock()

	key = g.getKey(key)
	if g.voteMap[key].Value == value {
		g.voteMap[key] = VoteStatus{
			Value:  value,
			Casted: true,
			Num:    number,
		}
	}
	if g.CanWriteGovernanceState(number) {
		g.WriteGovernanceState(number, false)
	}
}

func (g *Governance) ClearVotes(num uint64) {
	g.voteMapLock.Lock()
	defer g.voteMapLock.Unlock()

	g.GovernanceVotes = nil
	g.GovernanceTally = nil
	g.mu.Lock()
	g.changeSet = GovernanceSet{}
	g.mu.Unlock()
	g.voteMap = make(map[string]VoteStatus)
	logger.Info("Governance votes are cleared", "num", num)
}

// parseVoteValue parse vote.Value from []uint8 to appropriate type
func (g *Governance) ParseVoteValue(gVote *GovernanceVote) (*GovernanceVote, error) {
	var val interface{}
	k := GovernanceKeyMap[gVote.Key]

	// filter out if vote value is an interface list
	if reflect.TypeOf(gVote.Value) == reflect.TypeOf([]interface{}{}) {
		return nil, ErrValueTypeMismatch
	}

	switch k {
	case params.GovernanceMode, params.MintingAmount, params.MinimumStake, params.Ratio:
		val = string(gVote.Value.([]uint8))
	case params.GoverningNode, params.AddValidator, params.RemoveValidator:
		val = common.BytesToAddress(gVote.Value.([]uint8))
	case params.Epoch, params.CommitteeSize, params.UnitPrice, params.StakeUpdateInterval, params.ProposerRefreshInterval, params.ConstTxGasHumanReadable, params.Policy:
		gVote.Value = append(make([]byte, 8-len(gVote.Value.([]uint8))), gVote.Value.([]uint8)...)
		val = binary.BigEndian.Uint64(gVote.Value.([]uint8))
	case params.UseGiniCoeff, params.DeferredTxFee:
		gVote.Value = append(make([]byte, 8-len(gVote.Value.([]uint8))), gVote.Value.([]uint8)...)
		if binary.BigEndian.Uint64(gVote.Value.([]uint8)) != uint64(0) {
			val = true
		} else {
			val = false
		}
	default:
		logger.Warn("Unknown key was given", "key", k)
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
	gov.mu.Lock()
	defer gov.mu.Unlock()

	switch GovernanceKeyMap[vote.Key] {
	case params.GoverningNode:
		gov.changeSet[vote.Key] = vote.Value.(common.Address)
		return true
	case params.GovernanceMode, params.Ratio:
		gov.changeSet[vote.Key] = vote.Value.(string)
		return true
	case params.Epoch, params.StakeUpdateInterval, params.ProposerRefreshInterval, params.CommitteeSize, params.UnitPrice, params.ConstTxGasHumanReadable:
		gov.changeSet[vote.Key] = vote.Value.(uint64)
		return true
	case params.Policy:
		gov.changeSet[vote.Key] = vote.Value.(uint64)
		return true
	case params.MintingAmount, params.MinimumStake:
		gov.changeSet[vote.Key], _ = vote.Value.(string)
		return true
	case params.UseGiniCoeff, params.DeferredTxFee:
		gov.changeSet[vote.Key] = vote.Value.(bool)
		return true
	default:
		logger.Warn("Unknown key was given", "key", vote.Key)
	}
	return false
}

func GetDefaultGovernanceConfig(engine params.EngineType) *params.GovernanceConfig {
	gov := &params.GovernanceConfig{
		GovernanceMode: params.DefaultGovernanceMode,
		GoverningNode:  common.HexToAddress(params.DefaultGoverningNode),
		Reward:         GetDefaultRewardConfig(),
	}
	return gov
}

func GetDefaultIstanbulConfig() *params.IstanbulConfig {
	return &params.IstanbulConfig{
		Epoch:          params.DefaultEpoch,
		ProposerPolicy: params.DefaultProposerPolicy,
		SubGroupSize:   params.DefaultSubGroupSize,
	}
}

func GetDefaultRewardConfig() *params.RewardConfig {
	return &params.RewardConfig{
		MintingAmount:          big.NewInt(params.DefaultMintingAmount),
		Ratio:                  params.DefaultRatio,
		UseGiniCoeff:           params.DefaultUseGiniCoeff,
		DeferredTxFee:          params.DefaultDefferedTxFee,
		StakingUpdateInterval:  uint64(86400),
		ProposerUpdateInterval: uint64(3600),
		MinimumStake:           big.NewInt(2000000),
	}
}

func GetDefaultCliqueConfig() *params.CliqueConfig {
	return &params.CliqueConfig{
		Epoch:  params.DefaultEpoch,
		Period: params.DefaultPeriod,
	}
}

func CheckGenesisValues(c *params.ChainConfig) error {
	gov := NewGovernance(c, nil)

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

func (g *Governance) initializeCache() error {
	// get last n governance change block number
	indices, err := g.db.ReadRecentGovernanceIdx(params.GovernanceCacheLimit)
	if err != nil {
		return ErrNotInitialized
	}
	g.idxCache = indices
	// Put governance items into the itemCache
	for _, v := range indices {
		if num, data, err := g.ReadGovernance(v); err == nil {
			g.itemCache.Add(getGovernanceCacheKey(num), data)
			g.actualGovernanceBlock = num
		} else {
			logger.Crit("Couldn't read governance cache from database. Check database consistency", "index", v, "err", err)
		}
	}

	// the last one is the one to be used now
	ret, _ := g.itemCache.Get(getGovernanceCacheKey(g.actualGovernanceBlock))
	g.currentSetMu.Lock()
	g.currentSet = ret.(GovernanceSet)
	g.currentSetMu.Unlock()
	return nil
}

// getGovernanceCache returns cached governance config as a byte slice
func (g *Governance) getGovernanceCache(num uint64) (GovernanceSet, bool) {
	cKey := getGovernanceCacheKey(num)

	if ret, ok := g.itemCache.Get(cKey); ok && ret != nil {
		return ret.(GovernanceSet), true
	}
	return nil, false
}

func (g *Governance) addGovernanceCache(num uint64, data GovernanceSet) {
	// Don't update cache if num (block number) is smaller than the biggest number of cached block number
	if len(g.idxCache) > 0 && num <= g.idxCache[len(g.idxCache)-1] {
		return
	}
	cKey := getGovernanceCacheKey(num)
	g.itemCache.Add(cKey, data)
	g.addIdxCache(num)
}

// getGovernanceCacheKey returns cache key of the given block number
func getGovernanceCacheKey(num uint64) common.GovernanceCacheKey {
	v := fmt.Sprintf("%v", num)
	return common.GovernanceCacheKey(params.GovernanceCachePrefix + "_" + v)
}

func (g *Governance) addIdxCache(num uint64) {
	g.idxCache = append(g.idxCache, num)
	if len(g.idxCache) > params.GovernanceIdxCacheLimit {
		g.idxCache = g.idxCache[len(g.idxCache)-params.GovernanceIdxCacheLimit:]
	}
}

// Store new governance data on DB. This updates Governance cache too.
func (g *Governance) WriteGovernance(num uint64, data GovernanceSet, delta GovernanceSet) error {

	new := make(GovernanceSet)
	new = CopyGovernanceSet(new, data)

	// merge delta into data
	if delta != nil {
		new = CopyGovernanceSet(new, delta)
	}
	g.addGovernanceCache(num, new)
	return g.db.WriteGovernance(new, num)
}

func CopyGovernanceSet(dst GovernanceSet, src GovernanceSet) GovernanceSet {
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func (g *Governance) searchCache(num uint64) (uint64, bool) {
	for i := len(g.idxCache) - 1; i >= 0; i-- {
		if g.idxCache[i] <= num {
			return g.idxCache[i], true
		}
	}
	return 0, false
}

func (g *Governance) ReadGovernance(num uint64) (uint64, GovernanceSet, error) {
	if g.ChainConfig.Istanbul == nil {
		logger.Crit("Failed to read governance. ChainConfig.Istanbul == nil")
	}
	blockNum := CalcGovernanceInfoBlock(num, g.ChainConfig.Istanbul.Epoch)
	// Check cache first
	if gBlockNum, ok := g.searchCache(blockNum); ok {
		if data, okay := g.getGovernanceCache(gBlockNum); okay {
			return gBlockNum, data, nil
		}
	}
	if g.db != nil {
		bn, result, err := g.db.ReadGovernanceAtNumber(num, g.ChainConfig.Istanbul.Epoch)
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

func (g *Governance) GetGovernanceChange() GovernanceSet {
	g.mu.Lock()
	defer g.mu.Unlock()

	if len(g.changeSet) > 0 {
		return g.changeSet
	}
	return nil
}

func (gov *Governance) UpdateGovernance(number uint64, governance []byte) {
	var epoch uint64
	var ok bool

	if epoch, ok = gov.GetGovernanceValue(GovernanceKeyMapReverse[params.Epoch]).(uint64); !ok {
		if epoch, ok = gov.GetGovernanceValue(GovernanceKeyMapReverse[params.CliqueEpoch]).(uint64); !ok {
			logger.Error("Couldn't find epoch from governance items")
			return
		}
	}

	// Store updated governance information if exist
	if number%epoch == 0 {
		if len(governance) > 0 {
			tempData := []byte("")
			tempSet := GovernanceSet{}
			if err := rlp.DecodeBytes(governance, &tempData); err != nil {
				logger.Error("Failed to decode governance data", "number", number, "err", err, "data", governance)
				return
			}
			if err := json.Unmarshal(tempData, &tempSet); err != nil {
				logger.Error("Failed to unmarshal governance data", "number", number, "err", err, "data", tempData)
				return

			}
			tempSet = adjustDecodedSet(tempSet)

			// Store new currentSet to governance database
			gov.currentSetMu.RLock()
			if err := gov.WriteGovernance(number, gov.currentSet, tempSet); err != nil {
				logger.Crit("Failed to store new governance data", "number", number, "err", err)
			}
			gov.currentSetMu.RUnlock()
		}
	}
}

func (gov *Governance) removeDuplicatedVote(vote *GovernanceVote, number uint64) {
	gov.RemoveVote(vote.Key, vote.Value, number)
}

func (gov *Governance) UpdateCurrentGovernance(num uint64) {
	newNumber, newGovernanceSet, _ := gov.ReadGovernance(num)

	// Do the change only when the governance actually changed
	if newGovernanceSet != nil && newNumber != gov.actualGovernanceBlock {
		gov.actualGovernanceBlock = newNumber
		gov.currentSetMu.Lock()
		gov.currentSet = newGovernanceSet
		gov.currentSetMu.Unlock()
		gov.triggerChange(newGovernanceSet)
	}
}

func (gov *Governance) triggerChange(set GovernanceSet) {
	for k, v := range set {
		GovernanceItems[GovernanceKeyMap[k]].trigger(gov, k, v)
	}
}

func adjustDecodedSet(src GovernanceSet) GovernanceSet {
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

func (gov *Governance) GetGovernanceValue(key string) interface{} {
	gov.currentSetMu.RLock()
	defer gov.currentSetMu.RUnlock()

	if v, ok := gov.currentSet[key]; !ok {
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

	rChangeSet := make(GovernanceSet)
	if json.Unmarshal(change, &rChangeSet) != nil {
		return ErrUnmarshalGovChange
	}
	rChangeSet = adjustDecodedSet(rChangeSet)

	gov.mu.RLock()
	defer gov.mu.RUnlock()
	if len(rChangeSet) == len(gov.changeSet) {
		for k, v := range rChangeSet {
			if GovernanceKeyMap[k] == params.GoverningNode {
				if reflect.TypeOf(v) == stringT {
					v = common.HexToAddress(v.(string))
				}
			}
			if gov.changeSet[k] != v {
				logger.Error("Verification Error", "key", k, "received", rChangeSet[k], "have", gov.changeSet[k], "receivedType", reflect.TypeOf(rChangeSet[k]), "haveType", reflect.TypeOf(gov.changeSet[k]))
				return ErrVoteValueMismatch
			}
		}
	}
	return nil
}

type governanceJSON struct {
	BlockNumber     uint64                `json:"blockNumber"`
	ChainConfig     *params.ChainConfig   `json:"chainConfig"`
	VoteMap         map[string]VoteStatus `json:"voteMap"`
	NodeAddress     common.Address        `json:"nodeAddress"`
	GovernanceVotes []GovernanceVote      `json:"governanceVotes"`
	GovernanceTally []GovernanceTally     `json:"governanceTally"`
	CurrentSet      GovernanceSet         `json:"currentSet"`
	ChangeSet       GovernanceSet         `json:"changeSet"`
}

func (gov *Governance) toJSON(num uint64) ([]byte, error) {
	gov.governanceLock.RLock()
	defer gov.governanceLock.RUnlock()
	ret := &governanceJSON{
		BlockNumber:     num,
		ChainConfig:     gov.ChainConfig,
		VoteMap:         gov.voteMap,
		NodeAddress:     gov.nodeAddress,
		GovernanceVotes: gov.GovernanceVotes,
		GovernanceTally: gov.GovernanceTally,
		CurrentSet:      gov.currentSet,
		ChangeSet:       gov.changeSet,
	}
	j, _ := json.Marshal(ret)
	return j, nil
}

func (gov *Governance) UnmarshalJSON(b []byte) error {
	var j governanceJSON
	if err := json.Unmarshal(b, &j); err != nil {
		return err
	}
	gov.governanceLock.Lock()
	defer gov.governanceLock.Unlock()
	gov.ChainConfig = j.ChainConfig
	gov.voteMap = j.VoteMap
	gov.nodeAddress = j.NodeAddress
	gov.GovernanceVotes = j.GovernanceVotes
	gov.GovernanceTally = j.GovernanceTally
	gov.currentSet = adjustDecodedSet(j.CurrentSet)
	gov.changeSet = adjustDecodedSet(j.ChangeSet)
	gov.lastGovernanceStateBlock = j.BlockNumber

	return nil
}

func (gov *Governance) setVotesAndTally(votes []GovernanceVote, tally []GovernanceTally) {
	gov.governanceLock.Lock()
	defer gov.governanceLock.Unlock()
	gov.GovernanceVotes = make([]GovernanceVote, len(votes))
	gov.GovernanceTally = make([]GovernanceTally, len(tally))
	copy(gov.GovernanceVotes, votes)
	copy(gov.GovernanceTally, tally)
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

func (gov *Governance) ReadGovernanceState() {
	b, err := gov.db.ReadGovernanceState()
	if err != nil {
		logger.Info("No governance state found in a database")
		return
	}
	gov.UnmarshalJSON(b)
	params.SetStakingUpdateInterval(gov.ChainConfig.Governance.Reward.StakingUpdateInterval)
	params.SetProposerUpdateInterval(gov.ChainConfig.Governance.Reward.ProposerUpdateInterval)

	if gov.currentSet["param.txgashumanreadable"] != nil {
		params.TxGasHumanReadable = gov.currentSet["param.txgashumanreadable"].(uint64)
	}
	logger.Info("Successfully loaded governance state from database", "blockNumber", atomic.LoadUint64(&gov.lastGovernanceStateBlock))
}

func (gov *Governance) SetBlockchain(bc *blockchain.BlockChain) {
	gov.blockChain = bc
}

func (gov *Governance) SetTxPool(txpool *blockchain.TxPool) {
	gov.TxPool = txpool
}

func getGovernanceItemsFromChainConfig(config *params.ChainConfig) GovernanceSet {
	g := make(GovernanceSet)

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

func AddGovernanceCacheForTest(g *Governance, num uint64, config *params.ChainConfig) {
	// Don't update cache if num (block number) is smaller than the biggest number of cached block number

	data := getGovernanceItemsFromChainConfig(config)
	cKey := getGovernanceCacheKey(num)
	g.itemCache.Add(cKey, data)
	g.addIdxCache(num)
}
