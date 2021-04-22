package dispatch

type MessageHandler interface {
	OnMessage(msg *Message)
}

type MessageHandlerFunc func(message *Message)