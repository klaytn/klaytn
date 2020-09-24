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
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/klaytn/klaytn/log"
)

var logger = log.NewModuleLogger(log.ChainDataFetcher)

// Kafka connects to the brokers in an existing kafka cluster.
type Kafka struct {
	config   *KafkaConfig
	producer sarama.SyncProducer
	admin    sarama.ClusterAdmin
}

func NewKafka(conf *KafkaConfig) (*Kafka, error) {
	producer, err := newSyncProducer(conf)
	if err != nil {
		logger.Error("Failed to create a new producer", "brokers", conf.brokers)
		return nil, err
	}

	admin, err := newClusterAdmin(conf)
	if err != nil {
		logger.Error("Failed to create a new cluster admin", "brokers", conf.brokers)
		return nil, err
	}

	return &Kafka{
		config:   conf,
		producer: producer,
		admin:    admin,
	}, nil
}

func (k *Kafka) CreateTopic(topic string) error {
	return k.admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     k.config.partitions,
		ReplicationFactor: k.config.replicas,
	}, false)
}

func (k *Kafka) Publish(topic string, msg interface{}) error {
	item := &sarama.ProducerMessage{
		Topic: topic,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	item.Value = sarama.StringEncoder(data)

	_, _, err = k.producer.SendMessage(item)
	return err
}

func newSyncProducer(config *KafkaConfig) (sarama.SyncProducer, error) {
	return sarama.NewSyncProducer(config.brokers, config.saramaConfig)
}

func newClusterAdmin(config *KafkaConfig) (sarama.ClusterAdmin, error) {
	return sarama.NewClusterAdmin(config.brokers, config.saramaConfig)
}
