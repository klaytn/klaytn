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
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas/mocks"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/stretchr/testify/suite"
)

var (
	decodedTrue         = hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001")
	decodedFalse        = hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000000")
	ikip13Input         = hexutil.MustDecode("0x01ffc9a701ffc9a700000000000000000000000000000000000000000000000000000000")
	ikip13InvalidInput  = hexutil.MustDecode("0x01ffc9a7ffffffff00000000000000000000000000000000000000000000000000000000")
	ikip7Input          = hexutil.MustDecode("0x01ffc9a76578737100000000000000000000000000000000000000000000000000000000")
	ikip7MetadataInput  = hexutil.MustDecode("0x01ffc9a7a219a02500000000000000000000000000000000000000000000000000000000")
	ikip17Input         = hexutil.MustDecode("0x01ffc9a780ac58cd00000000000000000000000000000000000000000000000000000000")
	ikip17MetadataInput = hexutil.MustDecode("0x01ffc9a75b5e139f00000000000000000000000000000000000000000000000000000000")
)

type SuiteContractCaller struct {
	suite.Suite
	ctrl   *gomock.Controller
	api    *mocks.MockBlockchainAPI
	caller *contractCaller
}

func (s *SuiteContractCaller) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.api = mocks.NewMockBlockchainAPI(s.ctrl)
	s.caller = newContractCaller(s.api)
}

func (s *SuiteContractCaller) TearDownTest() {
	s.ctrl.Finish()
}

func setExpectation(m *mocks.MockBlockchainAPI, contract *common.Address, data, result []byte) {
	arg := api.CallArgs{
		From: common.Address{},
		To:   contract,
		Data: data,
	}

	m.EXPECT().Call(gomock.Any(), gomock.Eq(arg), gomock.Eq(rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber))).Return(result, nil).Times(1)
}

func (s *SuiteContractCaller) TestContractCaller_IsKIP13_Success() {
	addr := common.HexToAddress("1")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address

	// KIP13 validation
	setExpectation(s.api, &addr, ikip13Input, decodedTrue)
	setExpectation(s.api, &addr, ikip13InvalidInput, decodedFalse)

	isKIP13, err := s.caller.isKIP13(addr, nil)
	s.NoError(err)
	s.True(isKIP13)
}

func (s *SuiteContractCaller) TestContractCaller_IsKIP7_Success() {
	addr := common.HexToAddress("2")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address

	// IKIP7 and IKIP7Metadata validation
	setExpectation(s.api, &addr, ikip7Input, decodedTrue)
	setExpectation(s.api, &addr, ikip7MetadataInput, decodedTrue)

	isKIP7, err := s.caller.isKIP7(addr, nil)
	s.NoError(err)
	s.True(isKIP7)
}

func (s *SuiteContractCaller) TestContractCaller_IsKIP17_Success() {
	addr := common.HexToAddress("3")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address

	// IKIP17 and IKIP17Metadata validation
	setExpectation(s.api, &addr, ikip17Input, decodedTrue)
	setExpectation(s.api, &addr, ikip17MetadataInput, decodedTrue)

	isKIP17, err := s.caller.isKIP17(addr, nil)
	s.NoError(err)
	s.True(isKIP17)
}

func TestSuiteContractCaller(t *testing.T) {
	suite.Run(t, new(SuiteContractCaller))
}
