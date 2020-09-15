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

type Consumer struct {
	cancel   chan bool
	handler  map[string]func(*sarama.ConsumerMessage) error
	consumer sarama.ConsumerGroup
	ctx      context.Context
	isActive bool
}

func NewConsumer(ctx context.Context, consumer sarama.ConsumerGroup) *Consumer {
	return &Consumer{
		cancel:   make(chan bool),
		handler:  map[string]func(*sarama.ConsumerMessage) error{},
		ctx:      ctx,
		consumer: consumer,
		isActive: false,
	}
}

func (r *Consumer) Subscribe(topic string, handler func(*sarama.ConsumerMessage) error) error {
	if r.handler[topic] != nil {
		return nil
	}
	r.handler[topic] = handler

	if r.isActive {
		r.cancel <- true
	}
	go func() {
		defer r.consumer.Close()
		var topics []string
		for topic := range r.handler {
			topics = append(topics, topic)
		}
		h := func(err chan<- error) {
			err <- r.consumer.Consume(r.ctx, topics, r)
		}
		res := make(chan error, 1)
		for {
			go h(res)
			select {
			case err := <-res:
				if err != nil {
					logger.Error("consumed messages return an error", "topic", topic, "err", err)
				}
			case <-r.cancel:
				return
			case <-r.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (r *Consumer) Setup(sess sarama.ConsumerGroupSession) error {
	logger.Info("consumer was initialized", "id", sess.MemberID())
	r.isActive = true
	return nil
}

func (r *Consumer) Cleanup(sess sarama.ConsumerGroupSession) error {
	logger.Info("consumer was cleaned up", "id", sess.MemberID())
	r.isActive = false
	return nil
}

func (r *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		logger.Debug("Message claimed", "value", string(message.Value), "timestamp", message.Timestamp, "topic", message.Topic)
		r.handler[message.Topic](message)
		session.MarkMessage(message, "")
	}

	return nil
}
