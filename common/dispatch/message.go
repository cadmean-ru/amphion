package dispatch

const (
	MessageEmpty = -2147483648 + iota
	MessageWorkExec
	MessageExit
)

type Message struct {
	Id      int
	What    int
	StrData string
	AnyData interface{}
	ReplyTo *Message
}

func NewMessage(what int) *Message {
	return &Message{
		What: what,
	}
}

func NewMessageWithStringData(what int, data string) *Message {
	return &Message{
		What:    what,
		StrData: data,
	}
}

func NewMessageWithAnyData(what int, data interface{}) *Message {
	return &Message{
		What:    what,
		AnyData: data,
	}
}
