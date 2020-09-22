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
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hashicorp/go-uuid"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
	"github.com/klaytn/klaytn/log"
)

var (
	logger = log.NewModuleLogger(log.ChainDataFetcher)
	kb     *KafkaBroker
	once   sync.Once
)

type KafkaBroker struct {
	producer   sarama.AsyncProducer
	admin      sarama.ClusterAdmin
	brokers    []string
	handlers   map[string]func(*sarama.ConsumerMessage) error
	consumer   *Consumer
	replicas   int16
	partitions int32
}

func New(groupID string, brokerList []string, replicas int16, partitions int32) types.EventBroker {
	once.Do(func() {
		kb = &KafkaBroker{
			brokers:    brokerList,
			handlers:   map[string]func(*sarama.ConsumerMessage) error{},
			replicas:   replicas,
			partitions: partitions,
		}
		kb.newClusterAdmin()
		kb.newProducer()

		// TODO-ChainDataFetcher context has to be passed by outside.
		kb.consumer = NewConsumer(context.Background(), kb.newConsumerGroup(groupID))
	})

	return kb
}

func (k *KafkaBroker) Publish(topic string, msg interface{}) error {
	k.CreateTopic(topic)
	item := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(topic),
	}
	if v, ok := msg.(types.IKey); ok {
		item.Key = sarama.StringEncoder(v.Key())
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	item.Value = sarama.StringEncoder(data)

	k.producer.Input() <- item

	return nil
}

func (k *KafkaBroker) Subscribe(topic string, handler interface{}) error {
	k.CreateTopic(topic)
	h, ok := handler.(func(*sarama.ConsumerMessage) error)
	if !ok {
		return fmt.Errorf("unsupported type. type: %v", reflect.TypeOf(handler))
	}

	return k.consumer.Subscribe(topic, h)
}

func (k *KafkaBroker) CreateTopic(topic string) (types.Topic, error) {
	err := k.admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     k.partitions,
		ReplicationFactor: k.replicas,
	}, false)

	return types.Topic{Name: topic}, err
}

func (k *KafkaBroker) DeleteTopic(topic string) error {
	return k.admin.DeleteTopic(topic)
}

func (k *KafkaBroker) ListTopics() ([]types.Topic, error) {
	topics, err := k.admin.ListTopics()
	if err != nil {
		return nil, err
	}

	var ret []types.Topic
	for k := range topics {
		ret = append(ret, types.Topic{Name: k})
	}

	return ret, nil
}

func (k *KafkaBroker) newProducer() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(k.brokers, config)
	if err != nil {
		logger.Crit("Failed to start Sarama producer", "err", err, "config", config)
	}

	k.producer = producer
}

func (k *KafkaBroker) newConsumerGroup(groupID string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Group.Session.Timeout = 6 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 2 * time.Second

	id, _ := uuid.GenerateUUID()
	config.ClientID = fmt.Sprintf("%s-%s", groupID, id)

	consumer, err := sarama.NewConsumerGroup(k.brokers, groupID, config)
	if err != nil {
		logger.Crit("NewConsumerGroup is failed", "err", err, "groupId", groupID, "config", config)
	}

	return consumer
}

func (k *KafkaBroker) newClusterAdmin() {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	admin, err := sarama.NewClusterAdmin(k.brokers, config)
	if err != nil {
		logger.Crit("NewClusterAdmin is failed", "err", err, "config", config)
	}
	k.admin = admin
}
