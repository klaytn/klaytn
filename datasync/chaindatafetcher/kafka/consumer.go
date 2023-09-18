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
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
)

// Logger is the instance of a sarama.StdLogger interface that chaindatafetcher leaves the SDK level information.
// By default it is set to print all log messages as standard output, but you can set it to redirect wherever you want.
var Logger sarama.StdLogger = log.New(os.Stdout, "[Chaindatafetcher] ", log.LstdFlags)

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
	wrongMsgVersionErrorMsg    = "the message version is not supported"
	missingSegmentErrorMsg     = "there is a missing segment"
	noHandlerErrorMsg          = "the handler does not exist for the given topic"
	emptySegmentErrorMsg       = "there is no segment in the segment slice"
	bufferOverflowErrorMsg     = "the number of items in buffer exceeded the maximum"
	msgExpiredErrorMsg         = "the message is expired"
)

// TopicHandler is a handler function in order to consume published messages.
type TopicHandler func(message *sarama.ConsumerMessage) error

// Segment represents a message segment with the parsed headers.
type Segment struct {
	orig       *sarama.ConsumerMessage
	key        string
	total      uint64
	index      uint64
	value      []byte
	version    string
	producerId string
}

func (s *Segment) String() string {
	return fmt.Sprintf("key: %v, total: %v, index: %v, value: %v, version: %v, producerId: %v", s.key, s.total, s.index, string(s.value), s.version, s.producerId)
}

// newSegment creates a new segment structure after parsing the headers.
func newSegment(msg *sarama.ConsumerMessage) (*Segment, error) {
	if msg == nil {
		return nil, errors.New(nilConsumerMessageErrorMsg)
	}

	headerLen := len(msg.Headers)
	if headerLen != MsgHeaderLength && headerLen != LegacyMsgHeaderLength {
		return nil, fmt.Errorf("%v [header length: %v]", wrongHeaderNumberErrorMsg, headerLen)
	}

	version := ""
	producerId := ""

	if headerLen == MsgHeaderLength {
		keyVersion := string(msg.Headers[MsgHeaderVersion].Key)
		if keyVersion != KeyVersion {
			return nil, fmt.Errorf("%v [expected: %v, actual: %v]", wrongHeaderKeyErrorMsg, KeyVersion, keyVersion)
		}
		version = string(msg.Headers[MsgHeaderVersion].Value)
		switch version {
		case MsgVersion1_0:
			keyProducerId := string(msg.Headers[MsgHeaderProducerId].Key)
			if keyProducerId != KeyProducerId {
				return nil, fmt.Errorf("%v [expected: %v, actual: %v]", wrongHeaderKeyErrorMsg, KeyProducerId, keyProducerId)
			}
			producerId = string(msg.Headers[MsgHeaderProducerId].Value)
		default:
			return nil, fmt.Errorf("%v [available: %v]", wrongMsgVersionErrorMsg, MsgVersion1_0)
		}
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
	return &Segment{
		orig:       msg,
		key:        key,
		total:      totalSegments,
		index:      segmentIdx,
		value:      msg.Value,
		version:    version,
		producerId: producerId,
	}, nil
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
	Logger.Printf("[INFO] the chaindatafetcher consumer is created. [groupId: %s, config: %s]", groupId, config.String())
	return &Consumer{
		config:   config,
		group:    group,
		handlers: make(map[string]TopicHandler),
	}, nil
}

// Close stops the ConsumerGroup and detaches any running sessions. It is required to call
// this function before the object passes out of scope, as it will otherwise leak memory.
func (c *Consumer) Close() error {
	Logger.Println("[INFO] the chaindatafetcher consumer is closed")
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
	// If c.config.SaramaConfig.Consumer.Return.Errors is set to true, then
	// the errors while consuming the messages can be read from c.group.Errors() channel.
	// Otherwise, it leaves only error logs using sarama.Logger as default.
	return c.group.Errors()
}

// Subscribe subscribes the registered topics with the handlers until the consumer is closed.
func (c *Consumer) Subscribe(ctx context.Context) error {
	if len(c.handlers) == 0 || len(c.topics) == 0 {
		return errors.New("there is no registered handler")
	}

	// Iterate over consumer sessions.
	for {
		Logger.Println("[INFO] started to consume Kafka message")
		if err := c.group.Consume(ctx, c.topics, c); err == sarama.ErrClosedConsumerGroup {
			Logger.Println("[INFO] the consumer group is closed")
			return nil
		} else if err != nil {
			Logger.Printf("[ERROR] the consumption is failed [err: %s]\n", err.Error())
			return err
		}
		// TODO-Chaindatafetcher add retry logic and error callback here
	}
}

// The following 3 methods implements ConsumerGroupHandler, and they are called in the order of Setup, Cleanup and ConsumeClaim.
// In Subscribe function, Consume method triggers the functions to handle published messages.

// Setup is called at the beginning of a new session, before ConsumeClaim.
func (c *Consumer) Setup(s sarama.ConsumerGroupSession) error {
	return c.config.Setup(s)
}

// Cleanup is called at the end of a session, once all ConsumeClaim goroutines have exited
// but before the offsets are committed for the very last time.
func (c *Consumer) Cleanup(s sarama.ConsumerGroupSession) error {
	return c.config.Cleanup(s)
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
		if numBuffered > 0 && bufferedSegments[0].key == newSegment.key && bufferedSegments[0].producerId == newSegment.producerId {
			// there is a missing segment which should not exist.
			if newSegment.index > uint64(numBuffered) {
				Logger.Printf("[ERROR] there may be a missing segment [numBuffered: %d, newSegment: %s]\n", numBuffered, newSegment.String())
				return buffer, errors.New(missingSegmentErrorMsg)
			}

			// the segment is already inserted to buffer.
			if newSegment.index < uint64(numBuffered) {
				Logger.Printf("[WARN] the message is duplicated [newSegment: %s]\n", newSegment.String())
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
		Logger.Printf("[WARN] the message may be inserted already. drop the segment [segment: %s]\n", newSegment.String())
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
			Logger.Printf("[ERROR] getting handler is failed with the given topic. [topic: %s]\n", msg.Topic)
			return buffer, fmt.Errorf("%v: %v", noHandlerErrorMsg, msg.Topic)
		}

		if err := f(msg); err != nil {
			Logger.Printf("[ERROR] the handler is failed [key: %s]\n", string(msg.Key))
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
			Logger.Println("[ERROR] no segment exists in the given buffer slice")
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

// resetTimer resets the given timer if oldest message is changed and it returns oldest message if exists.
func (c *Consumer) resetTimer(buffer [][]*Segment, timer *time.Timer, oldestMsg *sarama.ConsumerMessage) *sarama.ConsumerMessage {
	if c.config.ExpirationTime <= time.Duration(0) {
		return nil
	}

	if len(buffer) == 0 {
		timer.Stop()
		return nil
	}

	if oldestMsg != buffer[0][0].orig {
		timer.Reset(c.config.ExpirationTime)
	}

	return buffer[0][0].orig
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
// Once the Messages() channel is closed, the Handler must finish its processing
// loop and exit.
func (c *Consumer) ConsumeClaim(cgs sarama.ConsumerGroupSession, cgc sarama.ConsumerGroupClaim) error {
	if buffer, err := c.consumeClaim(cgs, cgc); err != nil {
		return c.handleError(buffer, cgs, err)
	}
	return nil
}

func (c *Consumer) consumeClaim(cgs sarama.ConsumerGroupSession, cgc sarama.ConsumerGroupClaim) ([][]*Segment, error) {
	var (
		buffer          [][]*Segment // TODO-Chaindatafetcher better to introduce segment buffer structure with useful methods
		oldestMsg       *sarama.ConsumerMessage
		expirationTimer = time.NewTimer(c.config.ExpirationTime)
	)

	// make sure that the expirationTimer channel is empty
	if !expirationTimer.Stop() {
		<-expirationTimer.C
	}

	for {
		select {
		case <-expirationTimer.C:
			return buffer, errors.New(msgExpiredErrorMsg)
		case msg, ok := <-cgc.Messages():
			if !ok {
				return buffer, nil
			}

			if len(buffer) > c.config.MaxMessageNumber {
				return buffer, fmt.Errorf("%v: increasing buffer size may resolve this problem. [max: %v, current: %v]", bufferOverflowErrorMsg, c.config.MaxMessageNumber, len(buffer))
			}

			segment, err := newSegment(msg)
			if err != nil {
				return buffer, err
			}

			// insert a new message segment into the buffer
			buffer, err = insertSegment(segment, buffer)
			if err != nil {
				return buffer, err
			}

			// handle the buffered messages if any message can be reassembled
			buffer, err = c.handleBufferedMessages(buffer)
			if err != nil {
				return buffer, err
			}

			// reset the expiration timer if necessary and update the oldest message
			oldestMsg = c.resetTimer(buffer, expirationTimer, oldestMsg)

			// mark offset of the oldest message to be read
			if err := c.updateOffset(buffer, msg, cgs); err != nil {
				return buffer, err
			}
		}
	}
}

func (c *Consumer) handleError(buffer [][]*Segment, cgs ConsumerGroupSession, parentErr error) error {
	if len(buffer) <= 0 || c.config.ErrCallback == nil {
		return parentErr
	}

	oldestMsg := buffer[0][0].orig
	key := string(oldestMsg.Key)

	if err := c.config.ErrCallback(key); err != nil {
		return err
	}

	buffer = buffer[1:]
	return c.updateOffset(buffer, oldestMsg, cgs)
}
