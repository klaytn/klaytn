package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainConfig_CheckConfigForkOrder(t *testing.T) {
	assert.Nil(t, BaobabChainConfig.CheckConfigForkOrder())
	assert.Nil(t, CypressChainConfig.CheckConfigForkOrder())
}
