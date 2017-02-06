package peer

import (
	"bytes"

	"github.com/hashicorp/go-msgpack/codec"
)

// messageType is the type of gossip messages Peer will send.
type messageType uint8

const (
	messageTypeJoin = iota
	messageTypeLeave
	messageTypePushPull
	messageTypeUserEvent
)

func encodeMessage(t messageType, msg interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(uint8(t))

	handle := codec.MsgpackHandle{}
	encoder := codec.NewEncoder(buf, &handle)
	err := encoder.Encode(msg)
	return buf.Bytes(), err
}

func decodeMessage(buf []byte, out interface{}) error {
	var handle codec.MsgpackHandle
	return codec.NewDecoder(bytes.NewReader(buf), &handle).Decode(out)
}

type messageJoin struct{}

type messageLeave struct{}

// messagePushPull is used when doing a state exchange.
type messagePushPull struct{}

// messageUserEvent is used for user-generated events.
type messageUserEvent struct{}
