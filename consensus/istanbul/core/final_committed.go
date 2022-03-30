// Modifications Copyright 2018 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// This file is derived from quorum/consensus/istanbul/core/final_committed.go (2018/06/04).
// Modified and improved for the klaytn development.

package core

import (
	"time"

	"github.com/klaytn/klaytn/common"
)

func (c *core) handleFinalCommitted() error {
	logger := c.logger.NewWith("state", c.state)
	logger.Trace("Received a final committed proposal")

	if !c.committedTime.IsZero() {
		c.blockCommitTimeGauge.Update(int64(time.Since(c.committedTime)))
	}
	c.preprepareStartTime = time.Time{}
	c.prepreparedTime = time.Time{}
	c.preparedTime = time.Time{}
	c.committedTime = time.Time{}

	c.startNewRound(common.Big0)
	return nil
}
