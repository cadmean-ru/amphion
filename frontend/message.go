package frontend

import (
	"github.com/cadmean-ru/amphion/common/a"
)

const (
	MessageRender = 0
	MessageExit   = 1
)

type Message struct {
	Code a.Byte
	Data interface{}
}

func NewFrontendMessage(code a.Byte) Message {
	return Message{
		Code: code,
		Data: nil,
	}
}

func NewFrontendMessageWithData(code a.Byte, data interface{}) Message {
	return Message{
		Code: code,
		Data: data,
	}
}