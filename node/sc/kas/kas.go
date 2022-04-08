// Copyright 2020 The klaytn Authors
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

package kas

import (
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/log"
)

var logger = log.NewModuleLogger(log.KAS)

type KASConfig struct {
	Url, XChainId, User, Pwd string

	Operator       common.Address
	Anchor         bool
	AnchorPeriod   uint64
	RequestTimeout time.Duration
}
