package frontend

import "github.com/cadmean-ru/amphion/common"

const (
	MessageRender = 0
)

type Message struct {
	Code common.AByte
	Data interface{}
}

func NewFrontendMessage(code common.AByte) Message {
	return Message{
		Code: code,
		Data: nil,
	}
}

func NewFrontendMessageWithData(code common.AByte, data interface{}) Message {
	return Message{
		Code: code,
		Data: data,
	}
}