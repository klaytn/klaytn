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

package gasprice

import (
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	mock_api "github.com/klaytn/klaytn/api/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGasPrice_NewOracle(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBackend := mock_api.NewMockBackend(mockCtrl)
	params := Config{}
	oracle := NewOracle(mockBackend, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Blocks: 2}
	oracle = NewOracle(mockBackend, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 2, oracle.checkBlocks)
	assert.Equal(t, 1, oracle.maxEmpty)
	assert.Equal(t, 10, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Percentile: -1}
	oracle = NewOracle(mockBackend, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 0, oracle.percentile)

	params = Config{Percentile: 101}
	oracle = NewOracle(mockBackend, params)

	assert.Nil(t, oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 100, oracle.percentile)

	params = Config{Percentile: 101, Default: big.NewInt(123)}
	oracle = NewOracle(mockBackend, params)

	assert.Equal(t, big.NewInt(123), oracle.lastPrice)
	assert.Equal(t, 1, oracle.checkBlocks)
	assert.Equal(t, 0, oracle.maxEmpty)
	assert.Equal(t, 5, oracle.maxBlocks)
	assert.Equal(t, 100, oracle.percentile)
}

func TestGasPrice_SuggestPrice(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockBackend := mock_api.NewMockBackend(mockCtrl)
	params := Config{}
	oracle := NewOracle(mockBackend, params)

	price, err := oracle.SuggestPrice(nil)
	assert.Nil(t, price)
	assert.Nil(t, err)

	params = Config{Default: big.NewInt(123)}
	oracle = NewOracle(mockBackend, params)

	price, err = oracle.SuggestPrice(nil)
	assert.Equal(t, big.NewInt(123), price)
	assert.Nil(t, err)
}
