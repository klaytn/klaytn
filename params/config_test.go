// Copyright 2022 The klaytn Authors
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

package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainConfig_CheckConfigForkOrder(t *testing.T) {
	assert.Nil(t, BaobabChainConfig.CheckConfigForkOrder())
	assert.Nil(t, CypressChainConfig.CheckConfigForkOrder())
}

func TestChainConfig_Copy(t *testing.T) {
	a := CypressChainConfig
	b := a.Copy()

	b.UnitPrice = 0x1111
	assert.NotEqual(t, a.UnitPrice, b.UnitPrice)

	b.Istanbul.Epoch = 0x2222
	assert.NotEqual(t, a.Istanbul.Epoch, b.Istanbul.Epoch)

	b.Governance.Reward = &RewardConfig{Ratio: "11/22/33"}
	assert.NotEqual(t, a.Governance.Reward.Ratio, b.Governance.Reward.Ratio)
}

func BenchmarkChainConfig_Copy(b *testing.B) {
	a := CypressChainConfig
	for i := 0; i < b.N; i++ {
		a.Copy()
	}
}
