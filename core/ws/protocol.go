package ws

import (
	"encoding/json"
)

// Message 消息
type Message struct {
	Method  MethodID
	Payload []byte
}

func NewMsg(ID MethodID, payload interface{}) *Message {
	bts, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}
	return &Message{
		Method:  ID,
		Payload: bts,
	}
}

// Bytes 消息字节
func (msg *Message) Bytes() []byte {
	bts, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return bts
}

// MsgFromBytes 消息字节
func MsgFromBytes(data []byte) (*Message, error) {
	msg := &Message{}
	err := json.Unmarshal(data, msg)
	return msg, err
}

// DecodePayload 消息内容
func (msg *Message) DecodePayload(val interface{}) error {
	err := json.Unmarshal(msg.Payload, &val)
	return err
}

// MethodID ...
type MethodID int32

const (
	ReqJoinMsg MethodID = iota
	RespRefuseMsg
	ReqServerMsg
	ReqProcessMsg
	ReqNodeInfoMsg
	ReqCommandMsg

	RespAddNodeMsg
	RespRemoveNodeMsg

	RespAddNodeInfoMsg
	RespRemoveNodeInfoMsg
	RespCommandMsg
)
