// Modifications Copyright 2020 The klaytn Authors
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.
//
// InternalTxTracer is a full blown transaction tracer that extracts and reports all
// the internal calls made by a transaction, along with any useful information.
//
// This file is derived from eth/tracers/internal/tracers/call_tracer.js (2018/06/04).
// Modified and improved for the klaytn development.

package vm

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

func (i InternalTxTrace) MarshalJSON() ([]byte, error) {
	type internalTxTrace struct {
		Type  string          `json:"type"`
		From  *common.Address `json:"from,omitempty"`
		To    *common.Address `json:"to,omitempty"`
		Value *string         `json:"value,omitempty"`

		Gas     *hexutil.Uint64 `json:"gas,omitempty"`
		GasUsed *hexutil.Uint64 `json:"gasUsed,omitempty"`

		Input  *string `json:"input,omitempty"`  // hex string
		Output *string `json:"output,omitempty"` // hex string
		Error  *string `json:"error,omitempty"`

		Time  *string            `json:"time,omitempty"`
		Calls []*InternalTxTrace `json:"calls,omitempty"`

		Reverted *RevertedInfo `json:"reverted,omitempty"`
	}

	var enc internalTxTrace
	enc.Type = i.Type
	enc.From = i.From
	enc.To = i.To
	if i.Value != "" {
		enc.Value = &i.Value
	}
	if i.Gas != 0 {
		gas := hexutil.Uint64(i.Gas)
		enc.Gas = &gas
	}
	if i.GasUsed != 0 {
		gasUsed := hexutil.Uint64(i.GasUsed)
		enc.GasUsed = &gasUsed
	}
	if i.Input != "" {
		enc.Input = &i.Input
	}
	if i.Output != "" {
		enc.Output = &i.Output
	}
	if i.Error != nil {
		errStr := i.Error.Error()
		enc.Error = &errStr
	}
	if i.Time != time.Duration(0) {
		timeStr := i.Time.String()
		enc.Time = &timeStr
	}
	enc.Calls = i.Calls
	enc.Reverted = i.Reverted

	return json.Marshal(&enc)
}

func (i *InternalTxTrace) UnmarshalJSON(input []byte) error {
	type internalTxTrace struct {
		Type  string          `json:"type"`
		From  *common.Address `json:"from,omitempty"`
		To    *common.Address `json:"to,omitempty"`
		Value *string         `json:"value,omitempty"`

		Gas     *hexutil.Uint64 `json:"gas,omitempty"`
		GasUsed *hexutil.Uint64 `json:"gasUsed,omitempty"`

		Input  *string `json:"input,omitempty"`  // hex string
		Output *string `json:"output,omitempty"` // hex string
		Error  *string `json:"error,omitempty"`

		Time  *string            `json:"time,omitempty"`
		Calls []*InternalTxTrace `json:"calls,omitempty"`

		Reverted *RevertedInfo `json:"reverted,omitempty"`
	}
	var dec internalTxTrace
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	i.Type = dec.Type
	if dec.From != nil {
		i.From = dec.From
	}
	if dec.To != nil {
		i.To = dec.To
	}
	if dec.Value != nil {
		i.Value = *dec.Value
	}
	if dec.Gas != nil {
		i.Gas = uint64(*dec.Gas)
	}
	if dec.GasUsed != nil {
		i.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.Input != nil {
		i.Input = *dec.Input
	}
	if dec.Output != nil {
		i.Output = *dec.Output
	}
	if dec.Error != nil {
		i.Error = errors.New(*dec.Error)
	}
	if dec.Time != nil {
		t, err := time.ParseDuration(*dec.Time)
		if err != nil {
			return err
		}
		i.Time = t
	}
	i.Calls = dec.Calls
	i.Reverted = dec.Reverted
	return nil
}

func (r RevertedInfo) MarshalJSON() ([]byte, error) {
	type RevertedInfo struct {
		Contract *common.Address `json:"contract,omitempty"`
		Message  *string         `json:"message,omitempty"`
	}
	var enc RevertedInfo
	enc.Contract = r.Contract
	if r.Message != "" {
		enc.Message = &r.Message
	}
	return json.Marshal(&enc)
}

func (r *RevertedInfo) UnmarshalJSON(data []byte) error {
	type RevertedInfo struct {
		Contract *common.Address `json:"contract,omitempty"`
		Message  *string         `json:"message,omitempty"`
	}
	var dec RevertedInfo
	if err := json.Unmarshal(data, &dec); err != nil {
		return err
	}
	if dec.Message != nil {
		r.Message = *dec.Message
	}
	r.Contract = dec.Contract
	return nil
}
