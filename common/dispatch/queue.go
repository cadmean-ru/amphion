package dispatch

import "fmt"

//MessageQueue is a non blocking message buffer.
//It is used to communicate between the engine and the frontend.
//There are two channels under the hood - main and secondary.
//The size of main channel is specified in the constructor function,
//whereas the size of the secondary channel equals half the size of the main channel.
//You can lock the main channel, so new messages wont be send to it. Then all messages in it can be processed.
//While the main channel is locked, messages are sent to the secondary channel.
//As soon as the main channel is unlocked, all messages from the secondary channel are sent to main channel.
type MessageQueue struct {
	messageChan      chan *Message
	secMessageChan   chan *Message
	size             uint
	bufferSize       uint
	locked           bool
}

//Enqueue sends the specified message to the end of the queue.
//Depending on the current state of the MessageQueue it sent to either main or secondary channel.
func (q *MessageQueue) Enqueue(message *Message) {
	if q.size >= q.bufferSize {
		return
	}

	//if message.What == MessageRender && q.hasRenderMessage {
	//	return
	//}

	if q.locked {
		select {
		case q.secMessageChan <- message:
		default:
			fmt.Println("Warning! Secondary message buffer full!")
		}
		return
	}

	select {
	case q.messageChan <- message:
		q.size++
	default:
		fmt.Println("Warning! Message buffer full!")
		return
	}
}

//Dequeue Removes the first message from the queue.
//If queue is empty returns an empty Message.
func (q *MessageQueue) Dequeue() *Message {
	if q.size <= 0 {
		return NewMessage(MessageEmpty)
	}

	var msg = <-q.messageChan
	q.size--

	return msg
}

//DequeueBlocking Removes the first message from the queue.
//Id the queue is empty waits for a message to come.
func (q *MessageQueue) DequeueBlocking() *Message {
	var msg = <-q.messageChan
	q.size--
	return msg
}

//GetSize returns the current size of the queue.
func (q *MessageQueue) GetSize() uint {
	return q.size
}

//IsEmpty indicate if the queue is empty, i.e. size == 0.
func (q *MessageQueue) IsEmpty() bool {
	return q.size == 0
}

//LockMainChannel locks the main channel.
//While it is locked new messages will be sent to the secondary channel.
//The locked state remains until UnlockMainChannel is called.
func (q *MessageQueue) LockMainChannel() {
	q.locked = true
}

//UnlockMainChannel unlocks the main channel.
//All messages from the secondary channel are sent to the main one.
//After that new messages are sent to the main channel again.
func (q *MessageQueue) UnlockMainChannel() {
	q.locked = false
	for {
		select {
		case msg := <-q.secMessageChan:
			q.Enqueue(msg)
		default:
			return
		}
	}
}

//Close closes the underlying channels.
func (q *MessageQueue) Close() {
	close(q.messageChan)
	close(q.secMessageChan)
}

//NewMessageQueue creates a new instance of MessageQueue with the specified buffer size.
//Returns pointer to the MessageQueue instance.
func NewMessageQueue(bufferSize uint) *MessageQueue {
	return &MessageQueue{
		messageChan:    make(chan *Message, bufferSize),
		secMessageChan: make(chan *Message, bufferSize/2),
		size:           0,
		bufferSize:     bufferSize,
	}
}
