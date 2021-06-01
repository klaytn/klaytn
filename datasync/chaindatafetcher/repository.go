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

package chaindatafetcher

import (
	"time"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
)

const (
	InsertRetryInterval = 500 * time.Millisecond
	InsertMaxRetry      = 100
)

type HandleChainEventFn func(blockchain.ChainEvent, types.RequestType) error

//go:generate mockgen -destination=./mocks/repository_mock.go -package=mocks github.com/klaytn/klaytn/datasync/chaindatafetcher Repository
//go:generate mockgen -destination=./mocks/checkpoint_db_mock.go -package=mocks github.com/klaytn/klaytn/datasync/chaindatafetcher CheckpointDB
//go:generate mockgen -destination=./mocks/component_setter_mock.go -package=mocks github.com/klaytn/klaytn/datasync/chaindatafetcher ComponentSetter

type Repository interface {
	HandleChainEvent(event blockchain.ChainEvent, dataType types.RequestType) error
}

type CheckpointDB interface {
	ReadCheckpoint() (int64, error)
	WriteCheckpoint(checkpoint int64) error
}

type ComponentSetter interface {
	SetComponent(component interface{})
}
