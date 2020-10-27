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
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/Shopify/sarama"
)

//go:generate mockgen -destination=./mocks/consumer_group_session_mock.go -package=mocks github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka ConsumerGroupSession
// ConsumerGroupSession is for mocking sarama.ConsumerGroupSession for better testing.
type ConsumerGroupSession interface {
	MarkOffset(topic string, partition int32, offset int64, metadata string)
	MarkMessage(msg *sarama.ConsumerMessage, metadata string)
}

var (
	eventNameErrorMsg          = "the event name must be either 'blockgroup' or 'tracegroup'"
	nilConsumerMessageErrorMsg = "the given message should not be nil"
	wrongHeaderNumberErrorMsg  = "the number of header is not expected"
	wrongHeaderKeyErrorMsg     = "the header key is not expected"
	missingSegmentErrorMsg     = "there is a missing segment"
	noHandlerErrorMsg          = "the handler does not exist for the given topic"
	emptySegmentErrorMsg       = "there is no segment in the segment slice"
	bufferOverflowErrorMsg     = "the number of items in buffer exceeded the maximum"
)

// TopicHandler is a handler function in order to consume published messages.
type TopicHandler func(message *sarama.ConsumerMessage) error

// Segment represents a message segment with the parsed headers.
type Segment struct {
	orig  *sarama.ConsumerMessage
	key   string
	total uint64
	index uint64
	value []byte
}

func (s *Segment) String() string {
	return fmt.Sprintf("key: %v, total: %v, index: %v, value %v", s.key, s.total, s.index, string(s.value))
}

// newSegment creates a new segment structure after parsing the headers.
func newSegment(msg *sarama.ConsumerMessage) (*Segment, error) {
	if msg == nil {
		return nil, errors.New(nilConsumerMessageErrorMsg)
	}

	if len(msg.Headers) != MsgHeaderLength {
		return nil, fmt.Errorf("%v [header length: %v]", wrongHeaderNumberErrorMsg, len(msg.Headers))
	}

	// check the existence of KeyTotalSegments header
	keyTotalSegments := string(msg.Headers[MsgHeaderTotalSegments].Key)
	if keyTotalSegments != KeyTotalSegments {
		return nil, fmt.Errorf("%v [expected: %v, actual: %v]", wrongHeaderKeyErrorMsg, KeyTotalSegments, keyTotalSegments)
	}

	// check the existence of MsgHeaderSegmentIdx header
	keySegmentIdx := string(msg.Headers[MsgHeaderSegmentIdx].Key)
	if keySegmentIdx != KeySegmentIdx {
		return nil, fmt.Errorf("%v [expected: %v, actual: %v]", wrongHeaderKeyErrorMsg, KeySegmentIdx, keySegmentIdx)
	}

	key := string(msg.Key)
	totalSegments := binary.BigEndian.Uint64(msg.Headers[MsgHeaderTotalSegments].Value)
	segmentIdx := binary.BigEndian.Uint64(msg.Headers[MsgHeaderSegmentIdx].Value)
	return &Segment{orig: msg, key: key, total: totalSegments, index: segmentIdx, value: msg.Value}, nil
}

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
	if event != EventBlockGroup && event != EventTraceGroup {
		return fmt.Errorf("%v [given: %v]", eventNameErrorMsg, event)
	}
	topic := c.config.GetTopicName(event)
	c.topics = append(c.topics, topic)
	c.handlers[topic] = handler
	return nil
}

func (c *Consumer) Errors() <-chan error {
	// c.config.SaramaConfig.Consumer.Return.Errors has to be set to true, and
	// read the errors from c.member.Errors() channel.
	// Currently, it leaves only error logs as default.
	return c.group.Errors()
}

// Subscribe subscribes the registered topics with the handlers until the consumer is closed.
func (c *Consumer) Subscribe(ctx context.Context) error {
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

// insertSegment inserts the given segment to the given buffer.
// Assumption:
//  1. it is guaranteed that the order of segments is correct.
//  2. the inserted messages may be duplicated.
//
// We can consider the following cases.
// case1. new segment with index 0 is inserted into newly created segment slice.
// case2. new consecutive segment is inserted into the right position.
// case3. duplicated segment is ignored.
// case4. new sparse segment shouldn't be given, so return an error.
func insertSegment(newSegment *Segment, buffer [][]*Segment) ([][]*Segment, error) {
	for idx, bufferedSegments := range buffer {
		numBuffered := len(bufferedSegments)
		if numBuffered > 0 && bufferedSegments[0].key == newSegment.key {
			// there is a missing segment which should not exist.
			if newSegment.index > uint64(numBuffered) {
				logger.Error("there may be a missing segment", "numBuffered", numBuffered, "newSegment", newSegment)
				return buffer, errors.New(missingSegmentErrorMsg)
			}

			// the segment is already inserted to buffer.
			if newSegment.index < uint64(numBuffered) {
				logger.Warn("the message is duplicated", "newSegment", newSegment)
				return buffer, nil
			}

			// insert the segment to the buffer.
			buffer[idx] = append(bufferedSegments, newSegment)
			return buffer, nil
		}
	}

	if newSegment.index == 0 {
		// create a segment slice and append it.
		buffer = append(buffer, []*Segment{newSegment})
	} else {
		// the segment may be already handled.
		logger.Warn("the message may be inserted already. drop the segment", "segment", newSegment)
	}
	return buffer, nil
}

// handleBufferedMessages handles all consecutive complete messages in the buffer.
func (c *Consumer) handleBufferedMessages(buffer [][]*Segment) ([][]*Segment, error) {
	for len(buffer) > 0 {
		// if any message exists in the buffer
		oldestMsg, firstSegment, buffered := buffer[0], buffer[0][0], len(buffer[0])
		if uint64(buffered) != firstSegment.total {
			// not ready for assembling messages
			return buffer, nil
		}

		// ready for assembling message
		var msgBuffer []byte
		for _, segment := range oldestMsg {
			msgBuffer = append(msgBuffer, segment.value...)
		}
		msg := &sarama.ConsumerMessage{
			Key:   []byte(firstSegment.key),
			Value: msgBuffer,
		}

		f, ok := c.handlers[firstSegment.orig.Topic]
		if !ok {
			return buffer, fmt.Errorf("%v: %v", noHandlerErrorMsg, msg.Topic)
		}

		if err := f(msg); err != nil {
			return buffer, err
		}

		buffer = buffer[1:]
	}

	return buffer, nil
}

// updateOffset updates offset after handling messages.
// The offset should be marked for the oldest message (which is not read) in the given buffer.
// If there is no segment in the buffer, the last consumed message offset should be marked.
func (c *Consumer) updateOffset(buffer [][]*Segment, lastMsg *sarama.ConsumerMessage, session ConsumerGroupSession) error {
	if len(buffer) > 0 {
		if len(buffer[0]) <= 0 {
			return errors.New(emptySegmentErrorMsg)
		}

		oldestMsg := buffer[0][0].orig
		// mark the offset as the oldest message has not been read
		session.MarkOffset(oldestMsg.Topic, oldestMsg.Partition, oldestMsg.Offset, "")
	} else {
		// mark the offset as the last message has been read
		session.MarkMessage(lastMsg, "")
	}

	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *Consumer) ConsumeClaim(cgs sarama.ConsumerGroupSession, cgc sarama.ConsumerGroupClaim) error {
	var buffer [][]*Segment
	for msg := range cgc.Messages() {
		if len(buffer) > c.config.MaxMessageNumber {
			return fmt.Errorf("%v: increasing buffer size may resolve this problem. [max: %v, current: %v]", bufferOverflowErrorMsg, c.config.MaxMessageNumber, len(buffer))
		}

		segment, err := newSegment(msg)
		if err != nil {
			return err
		}

		// insert a new message segment into the buffer
		buffer, err = insertSegment(segment, buffer)
		if err != nil {
			return err
		}

		// handle the buffered messages if any message can be reassembled
		buffer, err = c.handleBufferedMessages(buffer)
		if err != nil {
			return err
		}

		// mark offset of the oldest message to be read
		if err := c.updateOffset(buffer, msg, cgs); err != nil {
			return err
		}
	}
	return nil
}
