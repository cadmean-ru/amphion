package dispatch

type Looper interface {
	Loop()
	GetMessageDispatcher() MessageDispatcher
}

type LooperImpl struct {
	queue    *MessageQueue
	handlers map[int]MessageHandler
}

func (l *LooperImpl) SendMessage(message *Message) {
	l.queue.Enqueue(message)
}

func (l *LooperImpl) Execute(item WorkItem) {
	l.queue.Enqueue(NewMessageWithAnyData(MessageWorkExec, item))
}

func (l *LooperImpl) Loop() {
	l.queue.LockMainChannel()

	for !l.queue.IsEmpty() {
		msg := l.queue.Dequeue()
		if handler, ok := l.handlers[msg.What]; ok {
			handler.OnMessage(msg)
		}
	}

	l.queue.UnlockMainChannel()
}

func (l *LooperImpl) SetMessageHandler(what int, handler MessageHandler) {
	l.handlers[what] = handler
}

func (l *LooperImpl) RemoveMessageHandler(what int) {
	delete(l.handlers, what)
}

func (l *LooperImpl) GetMessageDispatcher() MessageDispatcher {
	return l
}

func (l *LooperImpl) GetWorkDispatcher() WorkDispatcher {
	return l
}

func NewLooperImpl(messageQueueBuffer uint) *LooperImpl {
	return &LooperImpl{
		queue: NewMessageQueue(messageQueueBuffer),
	}
}

func NewLooperImplCompat(messageQueueBuffer int) *LooperImpl {
	return &LooperImpl{
		queue: NewMessageQueue(uint(messageQueueBuffer)),
	}
}