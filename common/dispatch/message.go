package dispatch

import "fmt"

const (
	MessageEmpty = -2147483648 + iota
	MessageWorkExec
	MessageExit
)

//Message holds an int message code and some data.
type Message struct {
	Id      int
	What    int
	StrData string
	IntData int
	AnyData interface{}
	Sender  interface{}
}

func (m *Message) String() string {
	return fmt.Sprintf("What: %d, StrData: %s, AnyData: %v", m.What, m.StrData, m.AnyData)
}

func NewMessage(what int) *Message {
	return &Message{
		What: what,
	}
}

func NewMessageFrom(from interface{}, what int) *Message {
	return &Message{
		What:   what,
		Sender: from,
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

func NewMessageFromWithAnyData(from interface{}, what int, data interface{}) *Message {
	return &Message{
		What:    what,
		AnyData: data,
		Sender:  from,
	}
}
