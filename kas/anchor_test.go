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
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/blockchain/types"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/crypto"
	"github.com/klaytn/klaytn/kas/mocks"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

var (
	errTest = errors.New("test error")
)

func testAnchorData() *types.AnchoringDataInternalType0 {
	return &types.AnchoringDataInternalType0{
		BlockHash:     common.HexToHash("0"),
		TxHash:        common.HexToHash("1"),
		ParentHash:    common.HexToHash("2"),
		ReceiptHash:   common.HexToHash("3"),
		StateRootHash: common.HexToHash("4"),
		BlockNumber:   big.NewInt(5),
		BlockCount:    big.NewInt(6),
		TxCount:       big.NewInt(7),
	}
}
func TestSendRequest(t *testing.T) {
	config := KASConfig{}
	anchor := NewKASAnchor(&config, nil, nil)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockHTTPClient(ctrl)
	anchor.client = m

	anchorData := testAnchorData()
	pl := dataToPayload(anchorData)

	// OK case
	{
		expectedRes := http.Response{Status: strconv.Itoa(http.StatusOK)}
		expectedRespBody := respBody{
			Code: 0,
		}
		bodyBytes, _ := json.Marshal(expectedRespBody)
		expectedRes.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))

		m.EXPECT().Do(gomock.Any()).Times(1).Return(&expectedRes, nil)
		resp, err := anchor.sendRequest(pl)

		assert.NoError(t, err)
		assert.Equal(t, expectedRespBody.Code, resp.Code)
	}

	// Error case
	{
		m.EXPECT().Do(gomock.Any()).Times(1).Return(nil, errTest)
		resp, err := anchor.sendRequest(pl)

		assert.Error(t, errTest, err)
		assert.Nil(t, resp)
	}
}

