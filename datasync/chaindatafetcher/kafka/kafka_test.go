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
	"strings"
	"testing"

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

func (s *KafkaSuite) TestKafka_Publish() {
	numTests := 10
	testBytesSize := 100

	s.kfk.CreateTopic(s.topic)

	var expected []*kafkaData
	for i := 0; i < numTests; i++ {
		testData := &kafkaData{common.MakeRandomBytes(testBytesSize)}
		s.NoError(s.kfk.Publish(s.topic, testData))
		expected = append(expected, testData)
	}

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
	numCheckCh := make(chan struct{}, numTests)
	testBytesSize := 100

	testEvent := EventBlockGroup
	topic := s.kfk.config.getTopicName(testEvent)
	s.kfk.CreateTopic(topic)

	// publish random data
	var expected []*kafkaData
	for i := 0; i < numTests; i++ {
		testData := &kafkaData{common.MakeRandomBytes(testBytesSize)}
		expected = append(expected, testData)
		s.NoError(s.kfk.Publish(topic, testData))
	}

	// make a test consumer group
	s.kfk.config.SaramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := NewConsumer(s.kfk.config, "test-groupId")
	s.NoError(err)
	defer consumer.Close()

	// add handler for the test event group
	var actual []*kafkaData
	err = consumer.AddTopicAndHandler(EventBlockGroup, func(message *sarama.ConsumerMessage) error {
		var d *kafkaData
		json.Unmarshal(message.Value, &d)
		actual = append(actual, d)
		numCheckCh <- struct{}{}
		return nil
	})
	s.NoError(err)

	// subscribe the added topics
	go func() {
		err := consumer.Subscribe(context.Background())
		s.NoError(err)
	}()

	// wait for all data to be consumed
	for i := 0; i < numTests; i++ {
		<-numCheckCh
	}

	// compare the results with the published data
	s.Equal(expected, actual)
}

func TestKafkaSuite(t *testing.T) {
	suite.Run(t, new(KafkaSuite))
}
