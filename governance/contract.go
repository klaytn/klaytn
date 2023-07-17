// Copyright 2022 The klaytn Authors
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

package governance

import (
	"errors"
	"math/big"

	"github.com/klaytn/klaytn/accounts/abi/bind/backends"
	"github.com/klaytn/klaytn/common"
	govcontract "github.com/klaytn/klaytn/contracts/gov"
	"github.com/klaytn/klaytn/params"
)

var (
	errContractEngineNotReady = errors.New("ContractEngine is not ready")
	errParamsAtFail           = errors.New("headerGov EffectiveParams() failed")
	errGovParamNotExist       = errors.New("GovParam does not exist")
	errInvalidGovParam        = errors.New("GovParam conversion failed")
)

type ContractEngine struct {
	currentParams *params.GovParamSet

	// for headerGov.EffectiveParams() and BlockChain()
	headerGov *Governance
}

func NewContractEngine(headerGov *Governance) *ContractEngine {
	e := &ContractEngine{
		currentParams: params.NewGovParamSet(),
		headerGov:     headerGov,
	}

	return e
}

// CurrentParams effective at upcoming block (head+1)
func (e *ContractEngine) CurrentParams() *params.GovParamSet {
	return e.currentParams
}

// Parameters effective at requested block (num)
func (e *ContractEngine) EffectiveParams(num uint64) (*params.GovParamSet, error) {
	pset, err := e.contractGetAllParamsAt(num)
	if err != nil {
		return nil, err
	}
	return pset, nil
}

// if UpdateParam fails, leave currentParams as-is
func (e *ContractEngine) UpdateParams(num uint64) error {
	// request the parameters required for generating the next block
	pset, err := e.contractGetAllParamsAt(num + 1)
	if err != nil {
		return err
	}

	e.currentParams = pset
	return nil
}

// contractGetAllParamsAt sets evmCtx.BlockNumber as num
func (e *ContractEngine) contractGetAllParamsAt(num uint64) (*params.GovParamSet, error) {
	chain := e.headerGov.BlockChain()
	if chain == nil {
		logger.Crit("headerGov.BlockChain() is nil")
		return nil, errContractEngineNotReady
	}

	config := chain.Config()
	if !config.IsKoreForkEnabled(new(big.Int).SetUint64(num)) {
		logger.Trace("ContractEngine disabled: hardfork block not passed")
		return params.NewGovParamSet(), nil
	}

	addr, err := e.contractAddrAt(num)
	if err != nil {
		return nil, err
	}
	if common.EmptyAddress(addr) {
		logger.Trace("ContractEngine disabled: GovParamContract address not set")
		return params.NewGovParamSet(), nil
	}

	caller := backends.NewBlockchainContractCaller(chain)
	contract, _ := govcontract.NewGovParamCaller(addr, caller)

	names, values, err := contract.GetAllParamsAt(nil, new(big.Int).SetUint64(num))
	if err != nil {
		logger.Warn("ContractEngine disabled: getAllParams call failed", "err", err)
		return params.NewGovParamSet(), nil
	}

	if len(names) != len(values) {
		logger.Warn("ContractEngine disabled: getAllParams result invalid", "len(names)", len(names), "len(values)", len(values))
		return params.NewGovParamSet(), nil
	}

	bytesMap := make(map[string][]byte)
	for i := 0; i < len(names); i++ {
		bytesMap[names[i]] = values[i]
	}
	return params.NewGovParamSetBytesMapTolerant(bytesMap), nil
}

// contractAddrAt returns the GovParamContract address effective at given block number
func (e *ContractEngine) contractAddrAt(num uint64) (common.Address, error) {
	headerParams, err := e.headerGov.EffectiveParams(num)
	if err != nil {
		logger.Error("headerGov.EffectiveParams failed", "err", err, "num", num)
		return common.Address{}, errParamsAtFail
	}

	// this happens when GovParamContract has not been voted
	param, ok := headerParams.Get(params.GovParamContract)
	if !ok {
		logger.Debug("Could not find GovParam contract address")
		return common.Address{}, nil
	}

	addr, ok := param.(common.Address)
	if !ok {
		logger.Error("Could not convert GovParam contract address into common.Address", "param", param)
		return common.Address{}, errInvalidGovParam
	}

	return addr, nil
}
