package frontend

const (
	MessageRender   = iota
	MessageExit
	MessageExec
	MessageNavigate
	MessageTitle
)

//Deprecated
//Use dispatch package instead
type Message struct {
	Code byte
	Data interface{}
}

//Deprecated
//Use dispatch package instead
func NewFrontendMessage(code byte) Message {
	return Message{
		Code: code,
		Data: nil,
	}
}

//Deprecated
//Use dispatch package instead
func NewFrontendMessageWithData(code byte, data interface{}) Message {
	return Message{
		Code: code,
		Data: data,
	}
}