package params

import (
	"bytes"
	"errors"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/log"
)

var (
	errUnknownGovParamKey  = errors.New("Unknown governance param key")
	errUnknownGovParamName = errors.New("Unknown governance param name")
	errBadGovParamValue    = errors.New("Malformed governance param value")
)

type govParamType struct {
	// The canonical type
	canonicalType reflect.Type

	// Parse arbitrary typed value into canonical type
	// Return false if not possible
	// Used to parse or normalize database content.
	parseValue func(v interface{}) (interface{}, bool)

	// Parse byte array into canonical type
	// Return false if not possible
	// Used to parse solidity contract content.
	parseBytes func(b []byte) (interface{}, bool)

	// Application-specific checks.
	// It is safe to assume that type of 'v' is canonicalType
	validate func(v interface{}) bool
}

func (ty *govParamType) ParseValue(v interface{}) (interface{}, bool) {
	if x, ok := ty.parseValue(v); ok {
		return x, ty.Validate(x)
	} else {
		return nil, false
	}
}

func (ty *govParamType) ParseBytes(b []byte) (interface{}, bool) {
	if x, ok := ty.parseBytes(b); ok {
		return x, ty.Validate(x)
	} else {
		return nil, false
	}
}

func (ty *govParamType) Validate(v interface{}) bool {
	// return ty.canonicalType == reflect.TypeOf(v) && ty.validate(v)
	if ty.canonicalType != reflect.TypeOf(v) {
		return false
	}
	if ty.validate != nil && !ty.validate(v) {
		return false
	}
	return true
}

var (
	govModeNames = map[string]int{
		"none":   GovernanceMode_None,
		"single": GovernanceMode_Single,
		"ballot": GovernanceMode_Ballot,
	}

	parseValueString = func(v interface{}) (interface{}, bool) {
		s, ok := v.(string)
		return s, ok
	}
	parseBytesString = func(b []byte) (interface{}, bool) {
		return string(b), true
	}
	validatePass = func(v interface{}) bool {
		return true
	}

	uint64ByteLen = int(reflect.TypeOf(uint64(0)).Size())

	govParamTypeGovMode = &govParamType{
		canonicalType: reflect.TypeOf("single"),
		parseValue:    parseValueString,
		parseBytes:    parseBytesString,
		validate: func(v interface{}) bool {
			_, ok := govModeNames[v.(string)]
			return ok
		},
	}

	govParamTypeAddress = &govParamType{
		canonicalType: reflect.TypeOf(common.Address{}),
		parseValue: func(v interface{}) (interface{}, bool) {
			switch v.(type) {
			case string:
				s := v.(string)
				return common.HexToAddress(s), common.IsHexAddress(s)
			case common.Address:
				return v, true
			default:
				return nil, false
			}
		},
		parseBytes: func(b []byte) (interface{}, bool) {
			return common.BytesToAddress(b), len(b) == common.AddressLength
		},
		validate: validatePass,
	}

	govParamTypeUint64 = &govParamType{
		canonicalType: reflect.TypeOf(uint64(0)),
		parseValue: func(v interface{}) (interface{}, bool) {
			switch v.(type) {
			case int:
				return uint64(v.(int)), v.(int) >= 0
			case uint:
				return uint64(v.(uint)), true
			case uint64:
				return v.(uint64), true
			case float64:
				return uint64(v.(float64)), v.(float64) >= 0
			default:
				return nil, false
			}
		},
		parseBytes: func(b []byte) (interface{}, bool) {
			// Must not exceed uint64 range
			return new(big.Int).SetBytes(b).Uint64(), len(b) <= uint64ByteLen
		},
		validate: validatePass,
	}

	govParamTypeBigInt = &govParamType{
		canonicalType: reflect.TypeOf(""),
		parseValue: func(v interface{}) (interface{}, bool) {
			switch v.(type) {
			case string:
				return v.(string), true
			case *big.Int:
				return v.(*big.Int).String(), true
			default:
				return nil, false
			}
		},
		parseBytes: parseBytesString,
		validate: func(v interface{}) bool {
			if n, ok := new(big.Int).SetString(v.(string), 10); ok {
				return n.Sign() >= 0 // must be non-negative.
			}
			return false
		},
	}

	govParamTypeRatio = &govParamType{
		canonicalType: reflect.TypeOf("12/34/54"),
		parseValue:    parseValueString,
		parseBytes:    parseBytesString,
		validate: func(v interface{}) bool {
			strs := strings.Split(v.(string), "/")
			if len(strs) != 3 {
				return false
			}
			sum := 0
			for _, s := range strs {
				n, err := strconv.Atoi(s)
				if err != nil || n < 0 {
					return false
				}
				sum += n
			}
			return sum == 100
		},
	}

	govParamTypeBool = &govParamType{
		canonicalType: reflect.TypeOf(true),
		parseValue: func(v interface{}) (interface{}, bool) {
			b, ok := v.(bool)
			return b, ok
		},
		parseBytes: func(b []byte) (interface{}, bool) {
			if bytes.Compare(b, []byte{0x01}) == 0 {
				return true, true
			} else if bytes.Compare(b, []byte{0x00}) == 0 {
				return false, true
			} else {
				return nil, false
			}
		},
		validate: validatePass,
	}
)

var govParamTypes = map[int]*govParamType{
	GovernanceMode:            govParamTypeGovMode,
	GoverningNode:             govParamTypeAddress,
	Epoch:                     govParamTypeUint64,
	Policy:                    govParamTypeUint64,
	CommitteeSize:             govParamTypeUint64,
	UnitPrice:                 govParamTypeUint64,
	MintingAmount:             govParamTypeBigInt,
	Ratio:                     govParamTypeRatio,
	UseGiniCoeff:              govParamTypeBool,
	DeferredTxFee:             govParamTypeBool,
	MinimumStake:              govParamTypeBigInt,
	StakeUpdateInterval:       govParamTypeUint64,
	ProposerRefreshInterval:   govParamTypeUint64,
	LowerBoundBaseFee:         govParamTypeUint64,
	UpperBoundBaseFee:         govParamTypeUint64,
	GasTarget:                 govParamTypeUint64,
	MaxBlockGasUsedForBaseFee: govParamTypeUint64,
	BaseFeeDenominator:        govParamTypeUint64,
}

var govParamNames = map[string]int{
	"governance.governancemode":       GovernanceMode,
	"governance.governingnode":        GoverningNode,
	"istanbul.epoch":                  Epoch,
	"istanbul.policy":                 Policy,
	"istanbul.committeesize":          CommitteeSize,
	"governance.unitprice":            UnitPrice,
	"reward.mintingamount":            MintingAmount,
	"reward.ratio":                    Ratio,
	"reward.useginicoeff":             UseGiniCoeff,
	"reward.deferredtxfee":            DeferredTxFee,
	"reward.minimumstake":             MinimumStake,
	"reward.stakingupdateinterval":    StakeUpdateInterval,
	"reward.proposerupdateinterval":   ProposerRefreshInterval,
	"kip71.lowerboundbasefee":         LowerBoundBaseFee,
	"kip71.upperboundbasefee":         UpperBoundBaseFee,
	"kip71.gastarget":                 GasTarget,
	"kip71.maxblockgasusedforbasefee": MaxBlockGasUsedForBaseFee,
	"kip71.basefeedenominator":        BaseFeeDenominator,
}

var govParamNamesReverse = map[int]string{}

func init() {
	for name, key := range govParamNames {
		govParamNamesReverse[key] = name
	}
}

// GovParamSet is an immutable set of governance parameters
// with various convenience getters.
type GovParamSet struct {
	// Items in canonical type.
	// Only type checked and validated values will be stored.
	items map[int]interface{}
}

func NewGovParamSet() *GovParamSet {
	return &GovParamSet{
		items: make(map[int]interface{}),
	}
}

// Return a new GovParamSet that contains keys from both input sets.
// If a key belongs to both sets, the value from `update` is used.
func NewGovParamSetMerged(base *GovParamSet, update *GovParamSet) *GovParamSet {
	p := NewGovParamSet()
	for key, value := range base.items {
		p.items[key] = value
	}
	for key, value := range update.items {
		p.items[key] = value
	}
	return p
}

func NewGovParamSetStrMap(items map[string]interface{}) (*GovParamSet, error) {
	p := NewGovParamSet()

	for name, value := range items {
		key, ok := govParamNames[name]
		if !ok {
			return nil, errUnknownGovParamName
		}
		err := p.set(key, value)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func NewGovParamSetIntMap(items map[int]interface{}) (*GovParamSet, error) {
	p := NewGovParamSet()

	for key, value := range items {
		err := p.set(key, value)
		if err != nil {
			return nil, err
		}
	}

	return p, nil
}

func NewGovParamSetBytesMap(items map[string][]byte) (*GovParamSet, error) {
	p := NewGovParamSet()

	for name, value := range items {
		key, ok := govParamNames[name]
		if !ok {
			return nil, errUnknownGovParamName
		}
		err := p.setBytes(key, value)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

func NewGovParamSetChainConfig(config *ChainConfig) (*GovParamSet, error) {
	items := make(map[int]interface{})
	if config.Istanbul != nil {
		items[Epoch] = config.Istanbul.Epoch
		items[Policy] = config.Istanbul.ProposerPolicy
		items[CommitteeSize] = config.Istanbul.SubGroupSize
	}
	items[UnitPrice] = config.UnitPrice
	if config.Governance != nil {
		items[GoverningNode] = config.Governance.GoverningNode
		items[GovernanceMode] = config.Governance.GovernanceMode
		if config.Governance.Reward != nil {
			if config.Governance.Reward.MintingAmount != nil {
				items[MintingAmount] = config.Governance.Reward.MintingAmount.String()
			}
			items[Ratio] = config.Governance.Reward.Ratio
			items[UseGiniCoeff] = config.Governance.Reward.UseGiniCoeff
			items[DeferredTxFee] = config.Governance.Reward.DeferredTxFee
			items[StakeUpdateInterval] = config.Governance.Reward.StakingUpdateInterval
			items[ProposerRefreshInterval] = config.Governance.Reward.ProposerUpdateInterval
			if config.Governance.Reward.MinimumStake != nil {
				items[MinimumStake] = config.Governance.Reward.MinimumStake.String()
			}
		}
		if config.Governance.KIP71 != nil {
			items[LowerBoundBaseFee] = config.Governance.KIP71.LowerBoundBaseFee
			items[UpperBoundBaseFee] = config.Governance.KIP71.UpperBoundBaseFee
			items[GasTarget] = config.Governance.KIP71.GasTarget
			items[MaxBlockGasUsedForBaseFee] = config.Governance.KIP71.MaxBlockGasUsedForBaseFee
			items[BaseFeeDenominator] = config.Governance.KIP71.BaseFeeDenominator
		}
	}

	return NewGovParamSetIntMap(items)
}

func (p *GovParamSet) set(key int, value interface{}) error {
	ty, ok := govParamTypes[key]
	if !ok {
		return errUnknownGovParamKey
	}
	parsed, ok := ty.ParseValue(value)
	if !ok {
		logger.Error("Bad GovParam value",
			"key", govParamNamesReverse[key], "value", value)
		return errBadGovParamValue
	}
	p.items[key] = parsed
	return nil
}

func (p *GovParamSet) setBytes(key int, bytes []byte) error {
	ty, ok := govParamTypes[key]
	if !ok {
		return errUnknownGovParamKey
	}
	parsed, ok := ty.ParseBytes(bytes)
	if !ok {
		logger.Error("Bad GovParam value",
			"key", govParamNamesReverse[key], "bytes", hexutil.Encode(bytes))
		return errBadGovParamValue
	}
	p.items[key] = parsed
	return nil
}

func (p *GovParamSet) StrMap() map[string]interface{} {
	m := map[string]interface{}{}
	for key, value := range p.items {
		m[govParamNamesReverse[key]] = value
	}
	return m
}

func (p *GovParamSet) IntMap() map[int]interface{} {
	m := map[int]interface{}{}
	for key, value := range p.items {
		m[key] = value
	}
	return m
}

// Returns a parameter value and a boolean indicating success.
func (p *GovParamSet) Get(key int) (interface{}, bool) {
	v, ok := p.items[key]
	return v, ok
}

// Return a parameter value or return a nil if the key does not exist.
func (p *GovParamSet) MustGet(key int) interface{} {
	if v, ok := p.Get(key); ok {
		return v
	} else {
		logger := log.NewModuleLogger(log.Governance)
		logger.Crit("Attempted to get missing GovParam item", "key", key, "name", govParamNamesReverse[key])
		return nil
	}
}

func (p *GovParamSet) ToIstanbulConfig() *IstanbulConfig {
	var ret IstanbulConfig
	if _, ok := p.Get(Epoch); ok {
		ret.Epoch = p.Epoch()
	}
	if _, ok := p.Get(Policy); ok {
		ret.ProposerPolicy = p.Policy()
	}
	if _, ok := p.Get(CommitteeSize); ok {
		ret.SubGroupSize = p.CommitteeSize()
	}

	return &ret
}

func (p *GovParamSet) ToRewardConfig() *RewardConfig {
	var ret RewardConfig
	if _, ok := p.Get(MintingAmount); ok {
		ret.MintingAmount = p.MintingAmountBig()
	}
	if _, ok := p.Get(Ratio); ok {
		ret.Ratio = p.Ratio()
	}
	if _, ok := p.Get(UseGiniCoeff); ok {
		ret.UseGiniCoeff = p.UseGiniCoeff()
	}
	if _, ok := p.Get(DeferredTxFee); ok {
		ret.DeferredTxFee = p.DeferredTxFee()
	}
	if _, ok := p.Get(StakeUpdateInterval); ok {
		ret.StakingUpdateInterval = p.StakeUpdateInterval()
	}
	if _, ok := p.Get(ProposerRefreshInterval); ok {
		ret.ProposerUpdateInterval = p.ProposerRefreshInterval()
	}
	if _, ok := p.Get(MinimumStake); ok {
		ret.MinimumStake = p.MinimumStakeBig()
	}

	return &ret
}

func (p *GovParamSet) ToKIP71Config() *KIP71Config {
	var ret KIP71Config
	if _, ok := p.Get(LowerBoundBaseFee); ok {
		ret.LowerBoundBaseFee = p.LowerBoundBaseFee()
	}
	if _, ok := p.Get(UpperBoundBaseFee); ok {
		ret.UpperBoundBaseFee = p.UpperBoundBaseFee()
	}
	if _, ok := p.Get(GasTarget); ok {
		ret.GasTarget = p.GasTarget()
	}
	if _, ok := p.Get(MaxBlockGasUsedForBaseFee); ok {
		ret.MaxBlockGasUsedForBaseFee = p.MaxBlockGasUsedForBaseFee()
	}
	if _, ok := p.Get(BaseFeeDenominator); ok {
		ret.BaseFeeDenominator = p.BaseFeeDenominator()
	}

	return &ret
}

func (p *GovParamSet) ToGovernanceConfig() *GovernanceConfig {
	var ret GovernanceConfig
	if _, ok := p.Get(GoverningNode); ok {
		ret.GoverningNode = p.GoverningNode()
	}
	if _, ok := p.Get(GovernanceMode); ok {
		ret.GovernanceMode = p.GovernanceModeStr()
	}
	ret.Reward = p.ToRewardConfig()
	ret.KIP71 = p.ToKIP71Config()

	return &ret
}

func (p *GovParamSet) ToChainConfig() *ChainConfig {
	var ret ChainConfig
	if _, ok := p.Get(UnitPrice); ok {
		ret.UnitPrice = p.UnitPrice()
	}
	ret.Istanbul = p.ToIstanbulConfig()
	ret.Governance = p.ToGovernanceConfig()

	return &ret
}

// Nominal getters. Shortcut for MustGet() + type assertion.

func (p *GovParamSet) GovernanceModeStr() string {
	return p.MustGet(GovernanceMode).(string)
}

func (p *GovParamSet) GovernanceModeInt() int {
	return govModeNames[p.GovernanceModeStr()]
}

func (p *GovParamSet) GoverningNode() common.Address {
	return p.MustGet(GoverningNode).(common.Address)
}

func (p *GovParamSet) Epoch() uint64 {
	return p.MustGet(Epoch).(uint64)
}

func (p *GovParamSet) Policy() uint64 {
	return p.MustGet(Policy).(uint64)
}

func (p *GovParamSet) CommitteeSize() uint64 {
	return p.MustGet(CommitteeSize).(uint64)
}

func (p *GovParamSet) UnitPrice() uint64 {
	return p.MustGet(UnitPrice).(uint64)
}

func (p *GovParamSet) MintingAmountStr() string {
	return p.MustGet(MintingAmount).(string)
}

func (p *GovParamSet) MintingAmountBig() *big.Int {
	n, _ := new(big.Int).SetString(p.MintingAmountStr(), 10)
	return n
}

func (p *GovParamSet) Ratio() string {
	return p.MustGet(Ratio).(string)
}

func (p *GovParamSet) UseGiniCoeff() bool {
	return p.MustGet(UseGiniCoeff).(bool)
}

func (p *GovParamSet) DeferredTxFee() bool {
	return p.MustGet(DeferredTxFee).(bool)
}

func (p *GovParamSet) MinimumStakeStr() string {
	return p.MustGet(MinimumStake).(string)
}

func (p *GovParamSet) MinimumStakeBig() *big.Int {
	n, _ := new(big.Int).SetString(p.MinimumStakeStr(), 10)
	return n
}

func (p *GovParamSet) StakeUpdateInterval() uint64 {
	return p.MustGet(StakeUpdateInterval).(uint64)
}

func (p *GovParamSet) ProposerRefreshInterval() uint64 {
	return p.MustGet(ProposerRefreshInterval).(uint64)
}

func (p *GovParamSet) Timeout() uint64 {
	return p.MustGet(Timeout).(uint64)
}

func (p *GovParamSet) LowerBoundBaseFee() uint64 {
	return p.MustGet(LowerBoundBaseFee).(uint64)
}

func (p *GovParamSet) UpperBoundBaseFee() uint64 {
	return p.MustGet(UpperBoundBaseFee).(uint64)
}

func (p *GovParamSet) GasTarget() uint64 {
	return p.MustGet(GasTarget).(uint64)
}

func (p *GovParamSet) MaxBlockGasUsedForBaseFee() uint64 {
	return p.MustGet(MaxBlockGasUsedForBaseFee).(uint64)
}

func (p *GovParamSet) BaseFeeDenominator() uint64 {
	return p.MustGet(BaseFeeDenominator).(uint64)
}
