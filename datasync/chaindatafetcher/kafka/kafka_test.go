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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/suite"
)

type KafkaSuite struct {
	suite.Suite
	kfk      *Kafka
	consumer sarama.Consumer
	topic    string
}

// In order to test KafkaSuite, any available kafka broker must be connectable with "kafka:9094".
// If no kafka broker is available, the KafkaSuite tests are skipped.
func (s *KafkaSuite) SetupTest() {
	config := GetDefaultKafkaConfig()
	config.Brokers = []string{"kafka:9094"}
	kfk, err := NewKafka(config)
	if err == sarama.ErrOutOfBrokers {
		s.T().Log("Failed connecting to brokers", config.Brokers)
		s.T().Skip()
	}
	s.NoError(err)
	s.kfk = kfk

	consumer, err := sarama.NewConsumer(config.Brokers, config.SaramaConfig)
	s.NoError(err)
	s.consumer = consumer
	s.topic = "test-topic"
}

func (s *KafkaSuite) TearDownTest() {
	s.kfk.Close()
}

func (s *KafkaSuite) TestKafka_CreateAndDeleteTopic() {
	// no topic to be deleted
	err := s.kfk.DeleteTopic(s.topic)
	s.Error(err)
	s.True(strings.Contains(err.Error(), sarama.ErrUnknownTopicOrPartition.Error()))

	// created a topic successfully
	err = s.kfk.CreateTopic(s.topic)
	s.NoError(err)

	// failed to create a duplicated topic
	err = s.kfk.CreateTopic(s.topic)
	s.Error(err)
	s.True(strings.Contains(err.Error(), sarama.ErrTopicAlreadyExists.Error()))

	// deleted a topic successfully
	s.Nil(s.kfk.DeleteTopic(s.topic))

	topics, err := s.kfk.ListTopics()
	if _, exist := topics[s.topic]; exist {
		s.Fail("topic must not exist")
	}
}

type kafkaData struct {
	Data []byte `json:"data"`
}

func (s *KafkaSuite) publishRandomData(topic string, numTests, testBytesSize int) []*kafkaData {
	var expected []*kafkaData
	for i := 0; i < numTests; i++ {
		testData := &kafkaData{common.MakeRandomBytes(testBytesSize)}
		s.NoError(s.kfk.Publish(topic, testData))
		expected = append(expected, testData)
	}
	return expected
}

func (s *KafkaSuite) subscribeData(topic, groupId string, numTests int, handler func(message *sarama.ConsumerMessage) error) {
	numCheckCh := make(chan struct{}, numTests)

	// make a test consumer group
	s.kfk.config.SaramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := NewConsumer(s.kfk.config, groupId)
	s.NoError(err)
	defer consumer.Close()

	// add handler for the test event group
	consumer.topics = append(consumer.topics, topic)
	consumer.handlers[topic] = func(message *sarama.ConsumerMessage) error {
		err := handler(message)
		numCheckCh <- struct{}{}
		return err
	}

	// subscribe the added topics
	go func() {
		err := consumer.Subscribe(context.Background())
		s.NoError(err)
	}()

	// wait for all data to be consumed
	timeout := time.NewTimer(3 * time.Second)
	for i := 0; i < numTests; i++ {
		select {
		case <-numCheckCh:
		case <-timeout.C:
			s.Fail("timeout")
		}
	}
}

func (s *KafkaSuite) TestKafka_Publish() {
	numTests := 10
	testBytesSize := 100

	s.kfk.CreateTopic(s.topic)

	expected := s.publishRandomData(s.topic, numTests, testBytesSize)

	// consume from the first partition and the first item
	partitionConsumer, err := s.consumer.ConsumePartition(s.topic, int32(0), int64(0))
	s.NoError(err)

	var actual []*kafkaData
	i := 0
	for msg := range partitionConsumer.Messages() {
		var dec *kafkaData
		json.Unmarshal(msg.Value, &dec)
		actual = append(actual, dec)
		i++
		if i == numTests {
			break
		}
	}

	s.True(len(actual) == numTests)
	for idx, v := range expected {
		s.Equal(v, actual[idx])
	}
}

func (s *KafkaSuite) TestKafka_Subscribe() {
	numTests := 10
	testBytesSize := 100

	topic := "test-subscribe"
	s.kfk.CreateTopic(topic)

	expected := s.publishRandomData(topic, numTests, testBytesSize)

	var actual []*kafkaData
	s.subscribeData(topic, "test-group-id", numTests, func(message *sarama.ConsumerMessage) error {
		var d *kafkaData
		json.Unmarshal(message.Value, &d)
		actual = append(actual, d)
		return nil
	})

	// compare the results with the published data
	s.Equal(expected, actual)
}

func (s *KafkaSuite) TestKafka_PubSubWith2Partitions() {
	numTests := 10
	testBytesSize := 100

	s.kfk.config.Partitions = 2
	defer func() { s.kfk.config.Partitions = DefaultPartitions }()

	topicPartition := "test-2-partition-topic"
	s.kfk.CreateTopic(topicPartition)

	// publish random data
	expected := s.publishRandomData(topicPartition, numTests, testBytesSize)

	var actual []*kafkaData
	s.subscribeData(topicPartition, "test-group-id", numTests, func(message *sarama.ConsumerMessage) error {
		var d *kafkaData
		json.Unmarshal(message.Value, &d)
		actual = append(actual, d)
		return nil
	})

	// the number of partitions is not 1, so the order may be changed after subscription.
	// compare the results with the published data
	s.Equal(len(expected), len(actual))

	for _, e := range expected {
		has := false
		for _, a := range actual {
			if reflect.DeepEqual(e, a) {
				has = true
			}
		}
		if !has {
			s.Fail("the expected data is not contained in the actual data", "expected", e)
		}
	}
}

func (s *KafkaSuite) TestKafka_PubSubWith2DifferentGroups() {
	numTests := 10
	testBytesSize := 100

	topic := "test-different-groups"
	s.kfk.CreateTopic(topic)

	// publish random data
	expected := s.publishRandomData(topic, numTests, testBytesSize)

	var actual []*kafkaData
	s.subscribeData(topic, "test-group-id-1", numTests, func(message *sarama.ConsumerMessage) error {
		var d *kafkaData
		json.Unmarshal(message.Value, &d)
		actual = append(actual, d)
		return nil
	})

	var actual2 []*kafkaData
	s.subscribeData(topic, "test-group-id-2", numTests, func(message *sarama.ConsumerMessage) error {
		var d *kafkaData
		json.Unmarshal(message.Value, &d)
		actual2 = append(actual2, d)
		return nil
	})

	// the number of partitions is not 1, so the order may be changed after subscription.
	// compare the results with the published data
	s.Equal(expected, actual)
	s.Equal(expected, actual2)
}

func (s *KafkaSuite) TestKafka_Consumer_AddTopicAndHandler() {
	consumer, err := NewConsumer(s.kfk.config, "test-group-id")
	s.NoError(err)

	blockGroupHandler := func(msg *sarama.ConsumerMessage) error { return nil }
	s.NoError(consumer.AddTopicAndHandler(EventBlockGroup, blockGroupHandler))
	traceGroupHandler := func(msg *sarama.ConsumerMessage) error { return nil }
	s.NoError(consumer.AddTopicAndHandler(EventTraceBroup, traceGroupHandler))

	blockGroupTopic := s.kfk.config.getTopicName(EventBlockGroup)
	traceGroupTopic := s.kfk.config.getTopicName(EventTraceBroup)
	expectedTopics := []string{blockGroupTopic, traceGroupTopic}
	s.Equal(expectedTopics, consumer.topics)
	s.Equal(reflect.ValueOf(blockGroupHandler).Pointer(), reflect.ValueOf(consumer.handlers[blockGroupTopic]).Pointer())
	s.Equal(reflect.ValueOf(traceGroupHandler).Pointer(), reflect.ValueOf(consumer.handlers[traceGroupTopic]).Pointer())
}

func (s *KafkaSuite) TestKafka_Consumer_AddTopicAndHandler_Error() {
	consumer, err := NewConsumer(s.kfk.config, "test-group-id")
	s.NoError(err)

	err = consumer.AddTopicAndHandler("not-available-event", nil)
	s.Error(err)
	s.True(strings.Contains(err.Error(), eventNameErrorMsg))
}

func TestKafkaSuite(t *testing.T) {
	suite.Run(t, new(KafkaSuite))
}
