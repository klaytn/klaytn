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

package kafka

import (
	"fmt"
	"math/big"

	"github.com/klaytn/klaytn/blockchain/vm"

	"github.com/klaytn/klaytn/blockchain"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
)

type traceGroupResult struct {
	BlockNumber      *big.Int              `json:"blockNumber"`
	InternalTxTraces []*vm.InternalTxTrace `json:"result"`
}

func (r *traceGroupResult) Key() string {
	return r.BlockNumber.String()
}

type blockGroupResult struct {
	BlockNumber *big.Int               `json:"blockNumber"`
	Result      map[string]interface{} `json:"result"`
}

func (r *blockGroupResult) Key() string {
	return r.BlockNumber.String()
}

type repository struct {
	blockchain *blockchain.BlockChain
	kafka      *Kafka
}

func NewRepository(config *KafkaConfig) (*repository, error) {
	kafka, err := NewKafka(config)
	if err != nil {
		logger.Error("Failed to create a new Kafka structure", "err", err, "config", config)
		return nil, err
	}
	return &repository{
		kafka: kafka,
	}, nil
}

func (r *repository) SetComponent(component interface{}) {
	switch c := component.(type) {
	case *blockchain.BlockChain:
		r.blockchain = c
	}
}

func (r *repository) HandleChainEvent(event blockchain.ChainEvent, dataType types.RequestType) error {
	switch dataType {
	case types.RequestTypeBlockGroup:
		result := &blockGroupResult{
			BlockNumber: event.Block.Number(),
			Result:      makeBlockGroupOutput(r.blockchain, event.Block, event.Receipts),
		}
		return r.kafka.Publish(r.kafka.getTopicName(EventBlockGroup), result)
	case types.RequestTypeTraceGroup:
		if len(event.InternalTxTraces) > 0 {
			result := &traceGroupResult{
				BlockNumber:      event.Block.Number(),
				InternalTxTraces: event.InternalTxTraces,
			}
			return r.kafka.Publish(r.kafka.getTopicName(EventTraceGroup), result)
		}
		return nil
	default:
		return fmt.Errorf("not supported type. [blockNumber: %v, reqType: %v]", event.Block.NumberU64(), dataType)
	}
}
