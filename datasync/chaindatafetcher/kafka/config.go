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
	"time"

	"github.com/Shopify/sarama"
	"github.com/klaytn/klaytn/common"
	"github.com/klaytn/klaytn/common/hexutil"
)

const (
	EventBlockGroup = "blockgroup"
	EventTraceGroup = "tracegroup"
)

const (
	topicProjectName = "klaytn"
	topicServiceName = "chaindatafetcher"
	topicVersion     = "v1"
)

const (
	DefaultReplicas             = 1
	DefaultPartitions           = 1
	DefaultTopicEnvironmentName = "local"
	DefaultTopicResourceName    = "en-0"
	DefaultMaxMessageBytes      = 1000000
	DefaultRequiredAcks         = 1
	DefaultSegmentSizeBytes     = 1000000 // 1 MB
	DefaultMaxMessageNumber     = 100     // max number of messages in buffer
	DefaultKafkaMessageVersion  = MsgVersion1_0
	DefaultProducerIdPrefix     = "producer-"
	DefaultExpirationTime       = time.Duration(0)
)

var (
	DefaultSetup   = func(s sarama.ConsumerGroupSession) error { return nil }
	DefaultCleanup = func(s sarama.ConsumerGroupSession) error { return nil }
)

type KafkaConfig struct {
	SaramaConfig         *sarama.Config `json:"-"` // kafka client configurations.
	MsgVersion           string         // MsgVersion is the version of Kafka message.
	ProducerId           string         // ProducerId is for the identification of the message publisher.
	Brokers              []string       // Brokers is a list of broker URLs.
	TopicEnvironmentName string
	TopicResourceName    string
	Partitions           int32 // Partitions is the number of partitions of a topic.
	Replicas             int16 // Replicas is a replication factor of kafka settings. This is the number of the replicated partitions in the kafka cluster.
	SegmentSizeBytes     int   // SegmentSizeBytes is the size of kafka message segment
	// (number of partitions) * (average size of segments) * buffer size should not be greater than memory size.
	// default max number of messages is 100
	MaxMessageNumber int // MaxMessageNumber is the maximum number of consumer messages.

	ExpirationTime time.Duration
	ErrCallback    func(string) error
	Setup          func(s sarama.ConsumerGroupSession) error
	Cleanup        func(s sarama.ConsumerGroupSession) error
}

func GetDefaultKafkaConfig() *KafkaConfig {
	// TODO-ChainDataFetcher add more configuration if necessary
	config := sarama.NewConfig()
	// The following configurations should be true
	config.Producer.Return.Errors = true
	config.Producer.Return.Successes = true
	config.Version = sarama.V2_4_0_0
	config.Producer.MaxMessageBytes = DefaultMaxMessageBytes
	config.Producer.RequiredAcks = sarama.RequiredAcks(DefaultRequiredAcks)
	return &KafkaConfig{
		SaramaConfig:         config,
		TopicEnvironmentName: DefaultTopicEnvironmentName,
		TopicResourceName:    DefaultTopicResourceName,
		Partitions:           DefaultPartitions,
		Replicas:             DefaultReplicas,
		SegmentSizeBytes:     DefaultSegmentSizeBytes,
		MaxMessageNumber:     DefaultMaxMessageNumber,
		MsgVersion:           DefaultKafkaMessageVersion,
		ProducerId:           GetDefaultProducerId(),
		ExpirationTime:       DefaultExpirationTime,
		Setup:                DefaultSetup,
		Cleanup:              DefaultCleanup,
		ErrCallback:          nil,
	}
}

func GetDefaultProducerId() string {
	rb := common.MakeRandomBytes(8)
	randomString := hexutil.Encode(rb)
	return DefaultProducerIdPrefix + randomString[2:]
}

func (c *KafkaConfig) GetTopicName(event string) string {
	return fmt.Sprintf("%v.%v.%v.%v.%v.%v", c.TopicEnvironmentName, topicProjectName, topicServiceName, c.TopicResourceName, event, topicVersion)
}

func (c *KafkaConfig) String() string {
	return fmt.Sprintf("brokers: %v, topicEnvironment: %v, topicResourceName: %v, partitions: %v, replicas: %v, maxMessageBytes: %v, requiredAcks: %v, segmentSize: %v, msgVersion: %v, producerId: %v",
		c.Brokers, c.TopicEnvironmentName, c.TopicResourceName, c.Partitions, c.Replicas, c.SaramaConfig.Producer.MaxMessageBytes, c.SaramaConfig.Producer.RequiredAcks, c.SegmentSizeBytes, c.MsgVersion, c.ProducerId)
}
