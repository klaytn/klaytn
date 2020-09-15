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

package kafka_test

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/Shopify/sarama"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/types"
	"github.com/stretchr/testify/suite"
)

type kafkaData struct {
	Name   string `json:"name"`
	Data   string `json:"data"`
	Number int    `json:"number"`
	IsGood bool   `json:"isGood"`
}

type KafkaSuite struct {
	suite.Suite
	kfk             types.EventBroker
	topic           string
	data            string
	consumerGroupId string
	brokers         []string
	replica         int16
	partitions      int32
}

func (s *KafkaSuite) SetupSuite() {
	s.topic = "test-topic"
	s.data = "test-data"
	s.consumerGroupId = "test_group_id"
	s.brokers = []string{"kafka:9094"}
	s.replica = 1
	s.partitions = 1
	s.kfk = kafka.New(s.consumerGroupId, s.brokers, s.replica, s.partitions)
}

func (s *KafkaSuite) TearDownSuite() {
	// Assume that the delete topic works properly.
	s.kfk.DeleteTopic(s.topic)
}

func (s *KafkaSuite) TestKafka_CreateAndDeleteTopics() {
	// no topic to be deleted
	err := s.kfk.DeleteTopic(s.topic)
	s.Error(err)
	s.True(strings.Contains(err.Error(), sarama.ErrUnknownTopicOrPartition.Error()))

	// created a topic successfully
	topic, err := s.kfk.CreateTopic(s.topic)
	s.NoError(err)
	s.Equal(s.topic, topic.Name)

	// failed to create a duplicated topic
	duplicatedTopic, err := s.kfk.CreateTopic(s.topic)
	s.Error(err)
	s.True(strings.Contains(err.Error(), sarama.ErrTopicAlreadyExists.Error()))
	s.Equal(topic, duplicatedTopic)

	// deleted a topic successfully
	s.Nil(s.kfk.DeleteTopic(s.topic))
}

func (s *KafkaSuite) TestKafka_Success_PubSub() {
	numTests := 10
	done := make(chan *kafkaData, numTests)
	s.kfk.CreateTopic(s.topic)

	// subscribe the test topic to check the data correctness
	err := s.kfk.Subscribe(s.topic, func(msg *sarama.ConsumerMessage) error {
		s.Equal(s.topic, msg.Topic)
		d := &kafkaData{}
		s.Nil(json.Unmarshal(msg.Value, d))
		done <- d

		return nil
	})
	s.NoError(err)

	// sleep for a while for the subscription to be set
	time.Sleep(3 * time.Second)

	t := time.NewTicker(10 * time.Second)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < numTests; i++ {
		// generate random kafkadata
		expected := &kafkaData{
			Name:   s.topic,
			Data:   s.data + strconv.Itoa(rand.Int()),
			Number: rand.Int(),
			IsGood: rand.Int()%2 == 0,
		}
		// publish the generated data
		s.Nil(s.kfk.Publish(s.topic, expected))

		select {
		case <-t.C:
			s.Fail("subscription timeout")
		case actual := <-done:
			s.Equal(expected, actual)
		}
	}
}

func (s *KafkaSuite) TestKafka_Success_ListTopic() {
	tc := []string{s.topic + "1", s.topic + "2"}
	for _, v := range tc {
		topic, err := s.kfk.CreateTopic(v)
		s.Nil(err)
		s.Equal(v, topic.Name)
	}
	// wait for the topics available across the brokers
	time.Sleep(2 * time.Second)

	topics, err := s.kfk.ListTopics()
	s.Nil(err)

	for _, tn := range tc {
		wasCreated := false
		for _, v := range topics {
			if v.Name == tn {
				wasCreated = true
				break
			}
		}
		s.True(wasCreated)
		s.kfk.DeleteTopic(tn)
	}
}

func TestKafkaSuite(t *testing.T) {
	suite.Run(t, new(KafkaSuite))
}
