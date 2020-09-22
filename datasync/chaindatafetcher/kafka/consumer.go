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

	"github.com/Shopify/sarama"
)

// Consumer subscribes a topic pushed by a producer in the kafka cluster.
type Consumer struct {
	cancelCh chan struct{}
	handler  map[string]func(*sarama.ConsumerMessage) error
	members  sarama.ConsumerGroup
	ctx      context.Context
	isActive bool
}

func NewConsumer(ctx context.Context, members sarama.ConsumerGroup) *Consumer {
	return &Consumer{
		cancelCh: make(chan struct{}),
		handler:  map[string]func(*sarama.ConsumerMessage) error{},
		ctx:      ctx,
		members:  members,
		isActive: false,
	}
}

func (c *Consumer) Subscribe(topic string, handler func(*sarama.ConsumerMessage) error) error {
	if c.handler[topic] != nil {
		return nil
	}
	c.handler[topic] = handler

	if c.isActive {
		c.cancelCh <- struct{}{}
	}
	go func() {
		defer c.members.Close()
		var topics []string
		for topic := range c.handler {
			topics = append(topics, topic)
		}
		res := make(chan error, 1)
		for {
			go func(err chan<- error) {
				err <- c.members.Consume(c.ctx, topics, c)
			}(res)
			select {
			case err := <-res:
				if err != nil {
					logger.Error("consumed messages return an error", "topic", topic, "err", err)
				}
			case <-c.cancelCh:
				return
			case <-c.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *Consumer) Setup(sess sarama.ConsumerGroupSession) error {
	logger.Info("consumer was initialized", "id", sess.MemberID())
	c.isActive = true
	return nil
}

func (c *Consumer) Cleanup(sess sarama.ConsumerGroupSession) error {
	logger.Info("consumer was cleaned up", "id", sess.MemberID())
	c.isActive = false
	return nil
}

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logger.Debug("Message claimed", "value", string(message.Value), "timestamp", message.Timestamp, "topic", message.Topic)
		c.handler[message.Topic](message)
		session.MarkMessage(message, "")
	}

	return nil
}
