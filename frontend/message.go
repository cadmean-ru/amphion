package frontend

const (
	MessageRender = 0
	MessageExit   = 1
	MessageExec   = 2
)

type Message struct {
	Code byte
	Data interface{}
}

func NewFrontendMessage(code byte) Message {
	return Message{
		Code: code,
		Data: nil,
	}
}

func NewFrontendMessageWithData(code byte, data interface{}) Message {
	return Message{
		Code: code,
		Data: data,
	}
}