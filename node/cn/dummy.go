// Copyright 2018 The klaytn Authors
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

package cn

import (
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/event"
)

type EmptyTxPool struct {
	txFeed event.Feed
}

func (re *EmptyTxPool) AddRemotes([]*types.Transaction) []error {
	return nil
}

func (re *EmptyTxPool) Pending() (map[common.Address]types.Transactions, error) {
	return map[common.Address]types.Transactions{}, nil
}

func (re *EmptyTxPool) SubscribeNewTxsEvent(newtxch chan<- blockchain.NewTxsEvent) event.Subscription {
	return re.txFeed.Subscribe(newtxch)
}
