package governance

import (
	"testing"

	"github.com/docker/docker/pkg/testutil/assert"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/params"
	"github.com/klaytn/klaytn/storage/database"
)

func newTestGovernanceApi() *PublicGovernanceAPI {
	config := params.CypressChainConfig
	config.Governance.KIP71 = params.GetDefaultKip71Config()
	govApi := NewGovernanceAPI(NewMixedEngine(config, database.NewMemoryDBManager()))
	govApi.governance.SetNodeAddress(common.HexToAddress("0x52d41ca72af615a1ac3301b0a93efa222ecc7541"))
	return govApi
}

func TestUpperBoundBaseFeeSet_LowerThan_LowerBoundBaseFee(t *testing.T) {
	govApi := newTestGovernanceApi()

	curLowerBoundBaseFee := govApi.governance.LowerBoundBaseFee()
	invalidUpperBoundBaseFee := curLowerBoundBaseFee - 100
	_, err := govApi.Vote("kip71.upperboundbasefee", invalidUpperBoundBaseFee)
	assert.Equal(t, err, errInvalidUpperBound)
}

func TestLowerBoundFeeSet_Exceeding_UpperBoundBaseFee(t *testing.T) {
	govApi := newTestGovernanceApi()

	curUpperBoundBaseFee := govApi.governance.UpperBoundBaseFee()
	invalidLowerBoundBaseFee := curUpperBoundBaseFee + 100
	_, err := govApi.Vote("kip71.lowerboundbasefee", invalidLowerBoundBaseFee)
	assert.Equal(t, err, errInvalidLowerBound)
}
