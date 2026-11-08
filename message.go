// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0
package rbus

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
)

const (
	FLAGS_REQUEST       = 1 << iota //nolint:revive // exported name kept for API compatibility
	FLAGS_RESPONSE                  //nolint:revive // exported name kept for API compatibility
	FLAGS_UNDELIVERABLE             //nolint:revive // exported name kept for API compatibility
	FLAGS_TAINTED                   //nolint:revive // exported name kept for API compatibility
	FLAGS_RAW_BINARY                //nolint:revive // exported name kept for API compatibility
	FLAGS_ENCRYPTED                 //nolint:revive // exported name kept for API compatibility

	headerVersion     = 2
	headerMarker      = 0xaaaa
	headerMaxTopicLen = 128
	headerMin         = 32
)

var (
	ErrInvalidMsg = errors.New("invalid message")
)

type Message struct {
	SeqNum     uint32
	CtrlData   uint32
	Flags      uint32
	Topic      string
	ReplyTopic string
	Payload    []byte
	Times      [5]uint32
}

func (m Message) Encode() ([]byte, error) {
	if len(m.Payload) > math.MaxUint32 {
		return nil, fmt.Errorf("%w: payload too large", ErrInvalidMsg)
	}
	if len(m.Topic) > headerMaxTopicLen {
		return nil, fmt.Errorf("%w: topic too large", ErrInvalidMsg)
	}
	if len(m.ReplyTopic) > headerMaxTopicLen {
		return nil, fmt.Errorf("%w: reply topic too large", ErrInvalidMsg)
	}

	buf := make([]byte, 0, 64+headerMaxTopicLen*2)

	buf = binary.BigEndian.AppendUint16(buf, headerMarker)
	buf = binary.BigEndian.AppendUint16(buf, headerVersion)
	buf = binary.BigEndian.AppendUint16(buf, 0) // Write a placeholder
	buf = binary.BigEndian.AppendUint32(buf, m.SeqNum)
	buf = binary.BigEndian.AppendUint32(buf, m.Flags)
	buf = binary.BigEndian.AppendUint32(buf, m.CtrlData)
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(m.Payload))) //nolint:gosec // bounded by check above
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(m.Topic)))   //nolint:gosec // bounded by headerMaxTopicLen
	buf = append(buf, []byte(m.Topic)...)
	buf = binary.BigEndian.AppendUint32(buf, uint32(len(m.ReplyTopic))) //nolint:gosec // bounded by headerMaxTopicLen
	buf = append(buf, []byte(m.ReplyTopic)...)

	// Only include timing information if one or more are non zero
	timing := make([]byte, 0, len(m.Times)*4)
	send := false
	for _, time := range m.Times {
		timing = binary.BigEndian.AppendUint32(timing, time)
		if time != 0 {
			send = true
		}
	}
	if send {
		buf = append(buf, timing...)
	}
	buf = binary.BigEndian.AppendUint16(buf, headerMarker)

	// Write the actual length.
	binary.BigEndian.PutUint16(buf[4:6], uint16(len(buf))) //nolint:gosec // header is bounded by topic length limits

	buf = append(buf, m.Payload...)

	return buf, nil
}

func Decode(buf []byte) (m *Message, err error) {
	// Since error should be few and far between so catching seem like a good approach.
	defer func() {
		if recover() != nil {
			m = nil
			err = fmt.Errorf("%w: buffer too small", ErrInvalidMsg)
		}
	}()

	m = new(Message)

	if headerMarker != binary.BigEndian.Uint16(buf) {
		return nil, fmt.Errorf("%w: invalid leading marker", ErrInvalidMsg)
	}
	headerLen := binary.BigEndian.Uint16(buf[2:3])

	header := buf[4 : headerLen-1]
	m.Payload = buf[headerLen:]

	if headerMarker != binary.BigEndian.Uint16(header[len(header)-2:]) {
		return nil, fmt.Errorf("%w: invalid trailing marker", ErrInvalidMsg)
	}
	header = header[:len(header)-2]

	m.SeqNum = binary.BigEndian.Uint32(header)
	header = header[4:]

	m.Flags = binary.BigEndian.Uint32(header)
	header = header[4:]

	m.CtrlData = binary.BigEndian.Uint32(header)
	header = header[4:]

	payloadLen := binary.BigEndian.Uint32(header)
	header = header[4:]
	if uint64(payloadLen) != uint64(len(m.Payload)) {
		return nil, fmt.Errorf("%w: invalid payload length", ErrInvalidMsg)
	}

	topicLen := binary.BigEndian.Uint32(header)
	header = header[4:]
	m.Topic = string(header[:topicLen])
	header = header[topicLen:]

	replyTopicLen := binary.BigEndian.Uint32(header)
	header = header[4:]
	m.ReplyTopic = string(header[:replyTopicLen])
	header = header[replyTopicLen:]

	if len(header) == 0 {
		return m, nil
	}

	if len(header) == 5*4 {
		for i := range 5 {
			t := binary.BigEndian.Uint32(header)
			header = header[4:]
			m.Times[i] = t
		}
		return m, nil
	}

	return nil, fmt.Errorf("%w: unknown protocol version", ErrInvalidMsg)
}
