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
	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/api"
	"github.com/klaytn/klaytn/blockchain/vm"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kas/mocks"
	"github.com/klaytn/klaytn/networks/rpc"
	"github.com/stretchr/testify/suite"
	"math/big"
	"testing"
	"time"
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

func setExpectation(m *mocks.MockBlockchainAPI, contract *common.Address, blockNumber *big.Int, data, result []byte) {
	arg := api.CallArgs{
		From: common.Address{},
		To:   contract,
		Data: data,
	}

	m.EXPECT().Call(gomock.Any(), gomock.Eq(arg), gomock.Eq(rpc.BlockNumber(blockNumber.Int64()))).Return(result, nil).Times(1)
}

func (s *SuiteContractCaller) TestContractCaller_IsKIP13_Success() {
	addr := common.HexToAddress("1")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address
	blockNumber := big.NewInt(1)

	// KIP13 validation
	setExpectation(s.api, &addr, blockNumber, ikip13Input, decodedTrue)
	setExpectation(s.api, &addr, blockNumber, ikip13InvalidInput, decodedFalse)

	isKIP13, err := s.caller.isKIP13(addr, blockNumber)
	s.NoError(err)
	s.True(isKIP13)
}

func (s *SuiteContractCaller) TestContractCaller_IsKIP7_Success() {
	addr := common.HexToAddress("2")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address
	blockNumber := big.NewInt(2)

	// IKIP7 and IKIP7Metadata validation
	setExpectation(s.api, &addr, blockNumber, ikip7Input, decodedTrue)
	setExpectation(s.api, &addr, blockNumber, ikip7MetadataInput, decodedTrue)

	isKIP7, err := s.caller.isKIP7(addr, blockNumber)
	s.NoError(err)
	s.True(isKIP7)
}

func (s *SuiteContractCaller) TestContractCaller_IsKIP17_Success() {
	addr := common.HexToAddress("3")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address
	blockNumber := big.NewInt(3)

	// IKIP17 and IKIP17Metadata validation
	setExpectation(s.api, &addr, blockNumber, ikip17Input, decodedTrue)
	setExpectation(s.api, &addr, blockNumber, ikip17MetadataInput, decodedTrue)

	isKIP17, err := s.caller.isKIP17(addr, blockNumber)
	s.NoError(err)
	s.True(isKIP17)
}

func (s *SuiteContractCaller) TestContractCaller_SupportsInterface_EmptyResult() {
	addr := common.HexToAddress("4")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address
	blockNumber := big.NewInt(4)
	fakeContract := common.HexToAddress("5")
	s.True(fakeContract != (common.Address{})) // make sure that the address is not empty address

	// Call method should return the byte slice result or an error. However, if a fallback function exists,
	// the method can return no result nor error. So, the next mocks a fallback function: both output and error are nil
	s.api.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).Times(1)

	// After this, it checks if contract exists or not. If contract exists, then it returns an error `abi: unmarshalling empty output`
	s.api.EXPECT().GetCode(gomock.Any(), gomock.Any(), gomock.Any()).Return(fakeContract.Bytes(), nil).Times(1)

	opts, cancel := getCallOpts(blockNumber, 1*time.Second)
	defer cancel()
	isSupported, err := s.caller.supportsInterface(addr, opts, [4]byte{})
	s.False(isSupported)
	s.Nil(err)
}

func (s *SuiteContractCaller) TestContractCaller_SupportsInterface_ExecutionReverted() {
	addr := common.HexToAddress("6")
	s.True(addr != (common.Address{})) // make sure that the address is not empty address
	blockNumber := big.NewInt(6)

	// evm execution reverted mocking
	s.api.EXPECT().Call(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, vm.ErrExecutionReverted).Times(1)

	opts, cancel := getCallOpts(blockNumber, 1*time.Second)
	defer cancel()
	isSupported, err := s.caller.supportsInterface(addr, opts, [4]byte{})
	s.False(isSupported)
	s.Nil(err)
}

func TestSuiteContractCaller(t *testing.T) {
	suite.Run(t, new(SuiteContractCaller))
}
