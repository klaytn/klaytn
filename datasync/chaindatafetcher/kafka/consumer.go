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
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
)

var eventNameErrorMsg = "the event name must be either 'blockgroup' or 'tracegroup'"

// TopicHandler is a handler function in order to consume published messages.
type TopicHandler func(message *sarama.ConsumerMessage) error

// Consumer is a reference structure to subscribe block or trace group produced by EN.
type Consumer struct {
	config   *KafkaConfig
	group    sarama.ConsumerGroup
	topics   []string
	handlers map[string]TopicHandler
}

func NewConsumer(config *KafkaConfig, groupId string) (*Consumer, error) {
	group, err := sarama.NewConsumerGroup(config.Brokers, groupId, config.SaramaConfig)
	if err != nil {
		return nil, err
	}
	return &Consumer{
		config:   config,
		group:    group,
		handlers: make(map[string]TopicHandler),
	}, nil
}

// Close stops the ConsumerGroup and detaches any running sessions. It is required to call
// this function before the object passes out of scope, as it will otherwise leak memory.
func (c *Consumer) Close() error {
	return c.group.Close()
}

// AddTopicAndHandler adds a topic associated the given event and its handler function to consume published messages of the topic.
func (c *Consumer) AddTopicAndHandler(event string, handler TopicHandler) error {
	if event != EventBlockGroup && event != EventTraceBroup {
		return fmt.Errorf("%v [given: %v]", eventNameErrorMsg, event)
	}
	topic := c.config.getTopicName(event)
	c.topics = append(c.topics, topic)
	c.handlers[topic] = handler
	return nil
}

// Subscribe subscribes the registered topics with the handlers until the consumer is closed.
func (c *Consumer) Subscribe(ctx context.Context) error {
	// TODO-ChainDataFetcher consider error handling if necessary
	// c.config.SaramaConfig.Consumer.Return.Errors has to be set to true, and
	// read the errors from c.member.Errors() channel.
	// Currently, it leaves only error logs as default.

	if len(c.handlers) == 0 || len(c.topics) == 0 {
		return errors.New("there is no registered handler")
	}

	// Iterate over consumer sessions.
	for {
		if err := c.group.Consume(ctx, c.topics, c); err == sarama.ErrClosedConsumerGroup {
			return nil
		} else if err != nil {
			return err
		}
	}
}

// The following 3 methods implements ConsumerGroupHandler, and they are called in the order of Setup, Cleanup and ConsumeClaim.
// In Subscribe function, Consume method triggers the functions to handle published messages.

// Setup is called at the beginning of a new session, before ConsumeClaim.
func (c *Consumer) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is called at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (c *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *Consumer) ConsumeClaim(cgs sarama.ConsumerGroupSession, cgc sarama.ConsumerGroupClaim) error {
	for msg := range cgc.Messages() {
		f, ok := c.handlers[msg.Topic]
		if !ok {
			return fmt.Errorf("the handler does not exist for the given topic: %v", msg.Topic)
		}
		if err := f(msg); err != nil {
			return err
		}
		// mark the message as consumed
		cgs.MarkMessage(msg, "")
	}
	return nil
}
