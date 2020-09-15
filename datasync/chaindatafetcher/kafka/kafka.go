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
	"errors"
	"fmt"
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
		kb.consumer = NewConsumer(context.Background(), kb.newConsumer(groupID))
	})

	return kb
}

func (r *KafkaBroker) Publish(topic string, msg interface{}) error {
	r.CreateTopic(topic)
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

	r.producer.Input() <- item

	return nil
}

func (r *KafkaBroker) Subscribe(topic string, handler interface{}) error {
	r.CreateTopic(topic)
	h, ok := handler.(func(*sarama.ConsumerMessage) error)
	if !ok {
		return errors.New("unsupported type")
	}

	return r.consumer.Subscribe(topic, h)
}

func (r *KafkaBroker) CreateTopic(topic string) (types.Topic, error) {
	err := r.admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     r.partitions,
		ReplicationFactor: r.replicas,
	}, false)

	return types.Topic{Name: topic}, err
}

func (r *KafkaBroker) DeleteTopic(topic string) error {
	return r.admin.DeleteTopic(topic)
}

func (r *KafkaBroker) ListTopics() ([]types.Topic, error) {
	topics, err := r.admin.ListTopics()
	if err != nil {
		return nil, err
	}

	var ret []types.Topic
	for k := range topics {
		ret = append(ret, types.Topic{Name: k})
	}

	return ret, nil
}

func (r *KafkaBroker) newProducer() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	producer, err := sarama.NewAsyncProducer(r.brokers, config)
	if err != nil {
		logger.Crit("Failed to start Sarama producer", "err", err)
	}

	r.producer = producer
}

func (r *KafkaBroker) newConsumer(groupID string) sarama.ConsumerGroup {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Group.Session.Timeout = 6 * time.Second
	config.Consumer.Group.Heartbeat.Interval = 2 * time.Second

	id, _ := uuid.GenerateUUID()
	config.ClientID = fmt.Sprintf("%s-%s", groupID, id)

	consumer, err := sarama.NewConsumerGroup(r.brokers, groupID, config)
	if err != nil {
		logger.Crit("NewConsumerGroup is failed", "err", err)
	}

	return consumer
}

func (r *KafkaBroker) newClusterAdmin() {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion

	admin, err := sarama.NewClusterAdmin(r.brokers, config)
	if err != nil {
		logger.Crit("NewClusterAdmin is failed", "err", err)
	}
	r.admin = admin
}
