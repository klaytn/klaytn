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
	"bytes"
	"errors"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/klaytn/klaytn/datasync/chaindatafetcher/kafka/mocks"

	"github.com/Shopify/sarama"
	"github.com/klaytn/klaytn/common"
	"github.com/stretchr/testify/assert"
)

func Test_newSegment_Success_LegacyMessage(t *testing.T) {
	value := common.MakeRandomBytes(100)
	rand.Seed(time.Now().UnixNano())
	total := uint64(10)
	idx := rand.Uint64() % total
	totalBytes := common.Int64ToByteBigEndian(total)
	idxBytes := common.Int64ToByteBigEndian(idx)
	key := "test-key"
	msg := &sarama.ConsumerMessage{
		Headers: []*sarama.RecordHeader{
			{Key: []byte(KeyTotalSegments), Value: totalBytes},
			{Key: []byte(KeySegmentIdx), Value: idxBytes},
		},
		Key:   []byte(key),
		Value: value,
	}

	segment, err := newSegment(msg)
	assert.NoError(t, err)
	assert.Equal(t, msg, segment.orig)
	assert.Equal(t, value, segment.value)
	assert.Equal(t, total, segment.total)
	assert.Equal(t, idx, segment.index)
	assert.Equal(t, key, segment.key)
}

func Test_newSegment_Success_Version1Message(t *testing.T) {
	value := common.MakeRandomBytes(100)
	rand.Seed(time.Now().UnixNano())
	total := uint64(10)
	idx := rand.Uint64() % total
	totalBytes := common.Int64ToByteBigEndian(total)
	idxBytes := common.Int64ToByteBigEndian(idx)
	key := "test-key"
	version := MsgVersion1_0
	versionBytes := []byte(version)
	producerId := GetDefaultProducerId()
	producerIdBytes := []byte(producerId)

	msg := &sarama.ConsumerMessage{
		Headers: []*sarama.RecordHeader{
			{Key: []byte(KeyTotalSegments), Value: totalBytes},
			{Key: []byte(KeySegmentIdx), Value: idxBytes},
			{Key: []byte(KeyVersion), Value: versionBytes},
			{Key: []byte(KeyProducerId), Value: producerIdBytes},
		},
		Key:   []byte(key),
		Value: value,
	}

	segment, err := newSegment(msg)
	assert.NoError(t, err)
	assert.Equal(t, msg, segment.orig)
	assert.Equal(t, value, segment.value)
	assert.Equal(t, total, segment.total)
	assert.Equal(t, idx, segment.index)
	assert.Equal(t, key, segment.key)
	assert.Equal(t, version, segment.version)
	assert.Equal(t, producerId, segment.producerId)
}

func Test_newSegment_Fail(t *testing.T) {
	type testcase struct {
		name   string
		input  *sarama.ConsumerMessage
		errMsg string
	}

	testcases := []testcase{
		{
			"nil consumer message",
			nil,
			nilConsumerMessageErrorMsg,
		},
		{
			"empty consumer message",
			&sarama.ConsumerMessage{},
			wrongHeaderNumberErrorMsg,
		},
		{
			"wrong first key with legacy header",
			&sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{Key: []byte("wrong-header-key")},
					{Key: []byte(KeyTotalSegments)},
				},
			},
			wrongHeaderKeyErrorMsg,
		},
		{
			"wrong second key with legacy header",
			&sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{Key: []byte(KeySegmentIdx)},
					{Key: []byte("wrong-header-key")},
				},
			},
			wrongHeaderKeyErrorMsg,
		},
		{
			"wrong first key with 1.0 header",
			&sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{Key: []byte("wrong-header-key")},
					{Key: []byte(KeyTotalSegments)},
					{Key: []byte(KeyVersion), Value: []byte(MsgVersion1_0)},
					{Key: []byte(KeyProducerId)},
				},
			},
			wrongHeaderKeyErrorMsg,
		},
		{
			"wrong second key with 1.0 header",
			&sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{Key: []byte(KeySegmentIdx)},
					{Key: []byte("wrong-header-key")},
					{Key: []byte(KeyVersion), Value: []byte(MsgVersion1_0)},
					{Key: []byte(KeyProducerId)},
				},
			},
			wrongHeaderKeyErrorMsg,
		},
		{
			"wrong third key with 1.0 header",
			&sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{Key: []byte(KeySegmentIdx)},
					{Key: []byte(KeyTotalSegments)},
					{Key: []byte("wrong-header-key")},
					{Key: []byte(KeyProducerId)},
				},
			},
			wrongHeaderKeyErrorMsg,
		},
		{
			"wrong fourth key with 1.0 header",
			&sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{Key: []byte(KeySegmentIdx)},
					{Key: []byte(KeyTotalSegments)},
					{Key: []byte(KeyVersion), Value: []byte(MsgVersion1_0)},
					{Key: []byte("wrong-header-key")},
				},
			},
			wrongHeaderKeyErrorMsg,
		},
		{
			"wrong version with 1.0 header",
			&sarama.ConsumerMessage{
				Headers: []*sarama.RecordHeader{
					{Key: []byte(KeySegmentIdx)},
					{Key: []byte(KeyTotalSegments)},
					{Key: []byte(KeyVersion), Value: []byte("wrong-version")},
					{Key: []byte(KeyProducerId)},
				},
			},
			wrongMsgVersionErrorMsg,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			seg, err := newSegment(tc.input)
			assert.Error(t, err)
			assert.True(t, strings.Contains(err.Error(), tc.errMsg))
			assert.Nil(t, seg)
		})
	}
}

func makeTestV1Segment(orig *sarama.ConsumerMessage, key string, total, index uint64, version, producerId string) *Segment {
	seg := makeTestSegment(orig, key, total, index)
	seg.version = version
	seg.producerId = producerId
	return seg
}

func makeTestSegment(orig *sarama.ConsumerMessage, key string, total, index uint64) *Segment {
	return &Segment{
		orig:       orig,
		key:        key,
		total:      total,
		index:      index,
		value:      common.MakeRandomBytes(5),
		version:    "",
		producerId: "",
	}
}

func Test_insertSegment_Success_OneMessage(t *testing.T) {
	// insert the message segments into the empty buffer in the order of m1s0, m1s1, m1s2
	var buffer [][]*Segment
	var err error

	msg1Key := "msg-1"
	total := uint64(3)
	m1s0 := makeTestSegment(nil, msg1Key, total, 0)
	m1s1 := makeTestSegment(nil, msg1Key, total, 1)
	m1s2 := makeTestSegment(nil, msg1Key, total, 2)

	buffer, err = insertSegment(m1s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m1s1, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m1s2, buffer)
	assert.NoError(t, err)

	// expected buffer
	// [m1s0][m1s1][m1s2]
	expected := [][]*Segment{{m1s0, m1s1, m1s2}}
	assert.Equal(t, expected, buffer)
}

func Test_insertSegment_Success_MultipleMessages(t *testing.T) {
	// insert the message segments into the empty buffer in the order of m1s0, m2s0, m3s0, m2s1, m1s1, m3s1
	var buffer [][]*Segment
	var err error

	msg1Key := "msg-1"
	total := uint64(3)
	m1s0 := makeTestSegment(nil, msg1Key, total, 0)
	m1s1 := makeTestSegment(nil, msg1Key, total, 1)

	msg2Key := "msg-2"
	m2s0 := makeTestSegment(nil, msg2Key, total, 0)
	m2s1 := makeTestSegment(nil, msg2Key, total, 1)

	msg3Key := "msg-3"
	m3s0 := makeTestSegment(nil, msg3Key, total, 0)
	m3s1 := makeTestSegment(nil, msg3Key, total, 1)

	buffer, err = insertSegment(m1s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m2s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m3s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m2s1, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m1s1, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m3s1, buffer)
	assert.NoError(t, err)

	// expected buffer
	// [m1s0][m1s1]
	// [m2s0][m2s1]
	// [m3s0][m3s1]
	expected := [][]*Segment{
		{m1s0, m1s1},
		{m2s0, m2s1},
		{m3s0, m3s1},
	}
	assert.Equal(t, expected, buffer)

}

func Test_insertSegment_Success_DuplicatedMessages(t *testing.T) {
	var buffer [][]*Segment
	var err error

	msg1Key := "msg-1"
	total := uint64(3)
	m1s0 := makeTestSegment(nil, msg1Key, total, 0)
	m1s1 := makeTestSegment(nil, msg1Key, total, 1)

	buffer, err = insertSegment(m1s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m1s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m1s1, buffer)
	assert.NoError(t, err)

	// expected buffer
	// [m1s0][m1s1]
	expected := [][]*Segment{
		{m1s0, m1s1},
	}
	assert.Equal(t, expected, buffer)

	msg2Key := "msg-2"
	m2s0 := makeTestSegment(nil, msg2Key, total, 0)
	m2s1 := makeTestSegment(nil, msg2Key, total, 1)

	buffer, err = insertSegment(m2s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m2s1, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m2s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m2s1, buffer)
	assert.NoError(t, err)

	// expected buffer
	// [m1s0][m1s1]
	// [m2s0][m2s1]
	expected = [][]*Segment{
		{m1s0, m1s1},
		{m2s0, m2s1},
	}
	assert.Equal(t, expected, buffer)
}

func Test_insertSegment_Success_IgnoreAlreadyInsertedSegment(t *testing.T) {
	var buffer [][]*Segment
	var err error

	msg1Key := "msg-1"
	total := uint64(3)
	m1s1 := makeTestSegment(nil, msg1Key, total, 1)

	buffer, err = insertSegment(m1s1, buffer)
	assert.NoError(t, err)

	msg2Key := "msg-2"
	m2s2 := makeTestSegment(nil, msg2Key, total, 2)

	buffer, err = insertSegment(m2s2, buffer)
	assert.NoError(t, err)

	// expected empty buffer
	var expected [][]*Segment
	assert.Equal(t, expected, buffer)
}

func Test_insertSegment_Fail_WrongSegmentError(t *testing.T) {
	var buffer [][]*Segment
	var err error

	msg1Key := "msg-1"
	total := uint64(3)
	m1s0 := makeTestSegment(nil, msg1Key, total, 0)
	m1s2 := makeTestSegment(nil, msg1Key, total, 2)

	buffer, err = insertSegment(m1s0, buffer)
	assert.NoError(t, err)
	buffer, err = insertSegment(m1s2, buffer)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), missingSegmentErrorMsg))
}

func TestConsumer_handleBufferedMessages_Success_NotCompleteMessage(t *testing.T) {
	testConsumer := &Consumer{}

	// not a complete message, so nothing is processed
	m1s0 := makeTestSegment(nil, "msg-1", 2, 0)
	buffer := [][]*Segment{{m1s0}}
	after, err := testConsumer.handleBufferedMessages(buffer)
	assert.NoError(t, err)
	assert.Equal(t, buffer, after)

	// a complete message in the middle
	m1s0 = makeTestSegment(nil, "msg-1", 2, 0)
	m2s0 := makeTestSegment(nil, "msg-2", 2, 0)
	m2s1 := makeTestSegment(nil, "msg-2", 2, 1)
	buffer = [][]*Segment{{m1s0}, {m2s0, m2s1}}
	after, err = testConsumer.handleBufferedMessages(buffer)
	assert.NoError(t, err)
	assert.Equal(t, buffer, after)
}

func TestConsumer_handleBufferedMessages_Success(t *testing.T) {
	testTopic := "test-topic"
	testOrig := &sarama.ConsumerMessage{Topic: testTopic}
	testConsumer := &Consumer{
		handlers: make(map[string]TopicHandler),
	}

	msg1Key := []byte("msg-1")
	msg2Key := []byte("msg-2")
	var msg1Expected []byte
	var msg2Expected []byte
	testConsumer.handlers[testTopic] = func(message *sarama.ConsumerMessage) error {
		if bytes.Equal(message.Key, msg1Key) {
			assert.Equal(t, msg1Expected, message.Value)
		}

		if bytes.Equal(message.Key, msg2Key) {
			assert.Equal(t, msg2Expected, message.Value)
		}
		return nil
	}

	emptyBuffer := [][]*Segment{}

	// a complete message with total 1
	m1s0 := makeTestSegment(testOrig, "msg-1", 1, 0)
	buffer := [][]*Segment{{m1s0}}
	msg1Expected = m1s0.value
	after, err := testConsumer.handleBufferedMessages(buffer)
	assert.NoError(t, err)
	assert.Equal(t, emptyBuffer, after)

	// a complete message with total 2
	m2s0 := makeTestSegment(testOrig, "msg-2", 2, 0)
	m2s1 := makeTestSegment(testOrig, "msg-2", 2, 1)
	buffer = [][]*Segment{{m2s0, m2s1}}
	msg2Expected = append(m2s0.value, m2s1.value...)
	after, err = testConsumer.handleBufferedMessages(buffer)
	assert.NoError(t, err)
	assert.Equal(t, emptyBuffer, after)

	// two complete messages
	buffer = [][]*Segment{{m1s0}, {m2s0, m2s1}}
	msg1Expected = m1s0.value
	msg2Expected = append(m2s0.value, m2s1.value...)
	after, err = testConsumer.handleBufferedMessages(buffer)
	assert.NoError(t, err)
	assert.Equal(t, emptyBuffer, after)
}

func TestConsumer_handleBufferedMessages_Fail(t *testing.T) {
	testTopic := "test-topic"
	testOrig := &sarama.ConsumerMessage{Topic: testTopic}
	testConsumer := &Consumer{
		handlers: make(map[string]TopicHandler),
	}

	// no handler is added
	m2s0 := makeTestSegment(testOrig, "msg-1", 2, 0)
	m2s1 := makeTestSegment(testOrig, "msg-1", 2, 1)
	buffer := [][]*Segment{{m2s0, m2s1}}
	_, err := testConsumer.handleBufferedMessages(buffer)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), noHandlerErrorMsg))
}

func TestConsumer_updateOffset(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	session := mocks.NewMockConsumerGroupSession(ctrl)

	testTopic := "test-topic"
	testOffset := int64(111)
	testPartition := int32(812)
	testMsg := &sarama.ConsumerMessage{Topic: testTopic, Partition: testPartition, Offset: testOffset}
	testConsumer := &Consumer{}

	// no segment in the buffer
	session.EXPECT().MarkMessage(gomock.Eq(testMsg), gomock.Eq(""))
	var buffer [][]*Segment
	err := testConsumer.updateOffset(buffer, testMsg, session)
	assert.NoError(t, err)

	// there is a segment in the buffer
	seg := makeTestSegment(testMsg, "msg-1", 3, 0)
	session.EXPECT().MarkOffset(gomock.Eq(seg.orig.Topic), gomock.Eq(seg.orig.Partition), gomock.Eq(seg.orig.Offset), gomock.Eq(""))
	buffer = [][]*Segment{{seg}}
	err = testConsumer.updateOffset(buffer, testMsg, session)
	assert.NoError(t, err)

	// wrong buffer
	buffer = [][]*Segment{{}}
	err = testConsumer.updateOffset(buffer, testMsg, session)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), emptySegmentErrorMsg))
}

func getTestConsumer() (string, *Consumer) {
	testTopic := "test-topic"
	testConsumer := &Consumer{
		handlers: make(map[string]TopicHandler),
	}

	testConsumer.handlers[testTopic] = func(message *sarama.ConsumerMessage) error {
		return nil
	}
	return testTopic, testConsumer
}

func Test_V1Segment_MultiProducers(t *testing.T) {
	{ // Legacy segments test
		var buffer [][]*Segment
		var err error

		msgKey := "msg"
		total := uint64(3)
		indices := []uint64{0, 1, 0, 2, 1, 0, 2, 1, 1 /* duplicated */, 2}

		topic, consumer := getTestConsumer()

		var msgs []*Segment
		for _, idx := range indices {
			msgs = append(msgs, makeTestSegment(&sarama.ConsumerMessage{Topic: topic}, msgKey, total, idx))
		}

		for idx, msg := range msgs {
			buffer, err = insertSegment(msg, buffer)
			if idx < 6 {
				assert.NoError(t, err)
			} else if idx == 6 {
				assert.Error(t, err)
				assert.True(t, strings.Contains(err.Error(), missingSegmentErrorMsg))
				break
			}
			buffer, err = consumer.handleBufferedMessages(buffer)
			assert.NoError(t, err)
		}
	}

	{ // V1 segments test
		var buffer [][]*Segment
		var err error

		msgKey := "msg"
		total := uint64(3)
		producerId := []string{"1", "1", "2", "1", "2", "3", "2", "3", "3", "3"}
		indices := []uint64{0, 1, 0, 2, 1, 0, 2, 1, 1 /* duplicated */, 2}

		topic, consumer := getTestConsumer()

		var msgs []*Segment
		for pi, si := range indices {
			msgs = append(msgs, makeTestV1Segment(&sarama.ConsumerMessage{Topic: topic}, msgKey, total, si, MsgVersion1_0, producerId[pi]))
		}

		for _, msg := range msgs {
			buffer, err = insertSegment(msg, buffer)
			assert.NoError(t, err)
			buffer, err = consumer.handleBufferedMessages(buffer)
			assert.NoError(t, err)
		}
	}
}

func TestConsumer_handleError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	testSegment := &Segment{
		orig: &sarama.ConsumerMessage{
			Topic:     "test-topic",
			Partition: int32(1),
			Offset:    int64(99),
		},
	}

	testSegment2 := &Segment{
		orig: &sarama.ConsumerMessage{
			Topic:     "test-topic",
			Partition: int32(1),
			Offset:    int64(100),
		},
	}

	session := mocks.NewMockConsumerGroupSession(ctrl)
	parentErr := errors.New("parent error")
	callbackErr := errors.New("callback error")

	tests := []struct {
		name      string
		setup     func()
		callback  func(s string) error
		buffer    [][]*Segment
		parentErr error
		expected  error
	}{
		{
			"no_item_in_buffer",
			func() {},
			func(s string) error { panic("should not be called") },
			[][]*Segment{},
			parentErr,
			parentErr,
		},
		{
			"no_callback",
			func() {},
			nil,
			[][]*Segment{{testSegment}},
			parentErr,
			parentErr,
		},
		{
			"callback_success_an_item_in_buffer",
			func() { session.EXPECT().MarkMessage(gomock.Eq(testSegment.orig), gomock.Eq("")).Times(1) },
			func(s string) error { return nil },
			[][]*Segment{{testSegment}},
			parentErr,
			nil,
		},
		{
			"callback_success_two_items_in_buffer",
			func() {
				session.EXPECT().MarkOffset(gomock.Eq(testSegment2.orig.Topic), gomock.Eq(testSegment2.orig.Partition), gomock.Eq(testSegment2.orig.Offset), gomock.Eq("")).Times(1)
			},
			func(s string) error { return nil },
			[][]*Segment{{testSegment}, {testSegment2}},
			parentErr,
			nil,
		},
		{
			"callback_failure",
			func() {},
			func(s string) error { return callbackErr },
			[][]*Segment{{testSegment}, {testSegment2}},
			parentErr,
			callbackErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consumer{}
			c.config = GetDefaultKafkaConfig()
			c.config.ErrCallback = tt.callback
			tt.setup()
			if err := c.handleError(tt.buffer, session, tt.parentErr); err != tt.expected {
				t.Errorf("handleError() error = %v, wantErr %v", err, tt.expected)
			}
		})
	}
}

func TestConsumer_renewExpireMsg(t *testing.T) {
	testTimer := time.NewTimer(1 * time.Second)
	testTimer.Stop()

	oldMsg := &sarama.ConsumerMessage{Key: []byte("old")}
	newMsg := &sarama.ConsumerMessage{Key: []byte("new")}

	type args struct {
		buffer    [][]*Segment
		timer     *time.Timer
		oldestMsg *sarama.ConsumerMessage
	}
	tests := []struct {
		name   string
		config *KafkaConfig
		args   args
		want   *sarama.ConsumerMessage
	}{
		{
			"no_expire_option",
			&KafkaConfig{ExpirationTime: time.Duration(0)},
			args{},
			nil,
		},
		{
			"all_handled",
			&KafkaConfig{ExpirationTime: time.Duration(1)},
			args{
				[][]*Segment{},
				testTimer,
				oldMsg,
			},
			nil,
		},
		{
			"oldest_handled",
			&KafkaConfig{ExpirationTime: time.Duration(1)},
			args{
				[][]*Segment{{&Segment{orig: newMsg}}},
				testTimer,
				oldMsg,
			},
			newMsg,
		},
		{
			"oldest_not_handled",
			&KafkaConfig{ExpirationTime: time.Duration(1)},
			args{
				[][]*Segment{{&Segment{orig: oldMsg}}},
				testTimer,
				oldMsg,
			},
			oldMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Consumer{
				config: tt.config,
			}
			if got := c.resetTimer(tt.args.buffer, tt.args.timer, tt.args.oldestMsg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("resetTimer() = %v, want %v", got, tt.want)
			}
		})
	}
}
