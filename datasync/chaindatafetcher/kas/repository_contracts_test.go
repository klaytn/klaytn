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
	"math/big"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas/mocks"
	"github.com/stretchr/testify/assert"
)

func TestFilterKIPContracts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	number := new(big.Int).SetInt64(1)
	block := types.NewBlockWithHeader(&types.Header{Number: number})

	// assume that test contracts are deployed
	var receipts types.Receipts
	receipts = append(receipts, &types.Receipt{ContractAddress: common.HexToAddress("1"), Status: types.ReceiptStatusSuccessful})
	receipts = append(receipts, &types.Receipt{ContractAddress: common.HexToAddress("2"), Status: types.ReceiptStatusSuccessful})
	receipts = append(receipts, &types.Receipt{ContractAddress: common.HexToAddress("3"), Status: types.ReceiptStatusSuccessful})
	receipts = append(receipts, &types.Receipt{ContractAddress: common.HexToAddress("4"), Status: types.ReceiptStatusSuccessful})
	event := blockchain.ChainEvent{
		Block:    block,
		Receipts: receipts,
	}

	api := mocks.NewMockBlockchainAPI(ctrl)

	// First contract is KIP13 only
	setExpectation(api, &receipts[0].ContractAddress, ikip13Input, decodedTrue)
	setExpectation(api, &receipts[0].ContractAddress, ikip13InvalidInput, decodedFalse)
	setExpectation(api, &receipts[0].ContractAddress, ikip7Input, decodedFalse)
	setExpectation(api, &receipts[0].ContractAddress, ikip17Input, decodedFalse)

	// Second one is IKIP7
	setExpectation(api, &receipts[1].ContractAddress, ikip13Input, decodedTrue)
	setExpectation(api, &receipts[1].ContractAddress, ikip13InvalidInput, decodedFalse)
	setExpectation(api, &receipts[1].ContractAddress, ikip7Input, decodedTrue)
	setExpectation(api, &receipts[1].ContractAddress, ikip7MetadataInput, decodedTrue)

	// Third one is IKIP17
	setExpectation(api, &receipts[2].ContractAddress, ikip13Input, decodedTrue)
	setExpectation(api, &receipts[2].ContractAddress, ikip13InvalidInput, decodedFalse)
	setExpectation(api, &receipts[2].ContractAddress, ikip7Input, decodedFalse)
	setExpectation(api, &receipts[2].ContractAddress, ikip17Input, decodedTrue)
	setExpectation(api, &receipts[2].ContractAddress, ikip17MetadataInput, decodedTrue)

	// Last one is not KIP
	setExpectation(api, &receipts[3].ContractAddress, ikip13Input, decodedFalse)

	fts, nfts, others, err := filterKIPContracts(api, event)
	assert.Equal(t, 1, len(fts))
	assert.Equal(t, 1, len(nfts))
	assert.Equal(t, 2, len(others))
	assert.NoError(t, err)

}
