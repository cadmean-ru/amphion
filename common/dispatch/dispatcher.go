package dispatch

type CallbackHandler func(message *Message)

type MessageDispatcher interface {
	SendMessage(message *Message)
}

type WorkDispatcher interface {
	Execute(item WorkItem)
}
