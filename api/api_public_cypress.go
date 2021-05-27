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

package api

import (
	"context"
	"errors"
	"strings"

	"github.com/klaytn/klaytn/accounts/abi"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/networks/rpc"
)

// cypressCreditABI is the input ABI used to generate the binding from.
const cypressCreditABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"getPhoto\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getNames\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"}]"

// CypressCredit contract is stored in the address zero.
var cypressCreditContractAddress = common.HexToAddress("0x0000000000000000000000000000000000000000")

var errNoCypressCreditContract = errors.New("no cypress credit contract")

type CreditOutput struct {
	Photo string `json:"photo"`
	Names string `json:"names"`
}

// callCypressCreditGetFunc executes funcName in CypressCreditContract and returns the output.
func (s *PublicBlockChainAPI) callCypressCreditGetFunc(ctx context.Context, parsed *abi.ABI, funcName string) (*string, error) {
	abiGet, err := parsed.Pack(funcName)
	if err != nil {
		return nil, err
	}

	args := CallArgs{
		To:   &cypressCreditContractAddress,
		Data: abiGet,
	}
	ret, err := s.Call(ctx, args, rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if err != nil {
		return nil, err
	}

	output := new(string)
	err = parsed.Unpack(output, funcName, ret)
	if err != nil {
		return nil, err
	}

	return output, nil
}

// GetCypressCredit calls getPhoto and getNames in the CypressCredit contract
// and returns all the results as a struct.
func (s *PublicBlockChainAPI) GetCypressCredit(ctx context.Context) (*CreditOutput, error) {
	answer, err := s.IsContractAccount(ctx, cypressCreditContractAddress, rpc.NewBlockNumberOrHashWithNumber(rpc.LatestBlockNumber))
	if err != nil {
		return nil, err
	}
	if !answer {
		return nil, errNoCypressCreditContract
	}

	parsed, err := abi.JSON(strings.NewReader(cypressCreditABI))
	if err != nil {
		return nil, err
	}

	photo, err := s.callCypressCreditGetFunc(ctx, &parsed, "getPhoto")
	if err != nil {
		return nil, err
	}

	names, err := s.callCypressCreditGetFunc(ctx, &parsed, "getNames")
	if err != nil {
		return nil, err
	}

	output := &CreditOutput{
		Photo: *photo,
		Names: *names,
	}

	return output, nil
}
