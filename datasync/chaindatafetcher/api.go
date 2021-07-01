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
	"errors"
	"sync/atomic"

	"github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
)

type PublicChainDataFetcherAPI struct {
	f *ChainDataFetcher
}

func NewPublicChainDataFetcherAPI(f *ChainDataFetcher) *PublicChainDataFetcherAPI {
	return &PublicChainDataFetcherAPI{f: f}
}

func (api *PublicChainDataFetcherAPI) StartFetching() error {
	return api.f.startFetching()
}

func (api *PublicChainDataFetcherAPI) StopFetching() error {
	return api.f.stopFetching()
}

func (api *PublicChainDataFetcherAPI) StartRangeFetching(start, end uint64, reqType uint) error {
	return api.f.startRangeFetching(start, end, types.RequestType(reqType))
}

func (api *PublicChainDataFetcherAPI) StopRangeFetching() error {
	return api.f.stopRangeFetching()
}

func (api *PublicChainDataFetcherAPI) Status() string {
	return api.f.status()
}

func (api *PublicChainDataFetcherAPI) ReadCheckpoint() (int64, error) {
	return api.f.checkpointDB.ReadCheckpoint()
}

func (api *PublicChainDataFetcherAPI) WriteCheckpoint(checkpoint int64) error {
	isRunning := atomic.LoadUint32(&api.f.fetchingStarted)
	if isRunning == running {
		return errors.New("call stopFetching before writing checkpoint manually")
	}

	api.f.checkpoint = checkpoint
	return api.f.checkpointDB.WriteCheckpoint(checkpoint)
}

// GetConfig returns the configuration setting of the launched chaindata fetcher.
func (api *PublicChainDataFetcherAPI) GetConfig() *ChainDataFetcherConfig {
	return api.f.config
}
