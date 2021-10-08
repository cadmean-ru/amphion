package engine

import (
	"github.com/cadmean-ru/amphion/common/dispatch"
)

const (
	MessageRedraw = iota
	MessageBuiltinEvent
	MessageUpdate
	MessageUpdateStop
)

const MessageMaxDepth = -1

type MessageListenerComponent interface {
	Component
	MessageListener
}

// MessageListener is the interface for receiving messages from dispatcher.
type MessageListener interface {
	// OnMessage receives and handles a message and returns whether to continue message propagation or not.
	OnMessage(message *dispatch.Message) bool
}

type MessageDispatcher struct {
	root *SceneObject
}

// Dispatch sends message from the root down in the message tree.
func (d *MessageDispatcher) Dispatch(message *dispatch.Message, maxDepth int) {
	d.sendMessageDown(d.root, message, maxDepth)
}

// DispatchDown sends this message up in the message tree
func (d *MessageDispatcher) DispatchDown(from MessageListener, message *dispatch.Message, maxDepth int) {
	c := d.findObject(d.root, from)
	if c == nil {
		return
	}
	d.sendMessageDown(c, message, maxDepth)
}

// DispatchUp sends specified message down in the message tree.
func (d *MessageDispatcher) DispatchUp(from MessageListener, message *dispatch.Message) {
	c := d.findObject(d.root, from)
	if c == nil {
		return
	}
	d.sendMessageUp(c, message)
}

func (d *MessageDispatcher) sendMessageDown(object *SceneObject, message *dispatch.Message, counter int) {
	if object.OnMessage(message) && (counter == MessageMaxDepth || counter > 0) {
		for _, c := range object.children {
			nextCounter := 0
			if counter == MessageMaxDepth {
				nextCounter = MessageMaxDepth
			} else {
				nextCounter = counter - 1
			}
			d.sendMessageDown(c, message, nextCounter)
		}
	}
}

func (d *MessageDispatcher) findObject(object *SceneObject, listener MessageListener) *SceneObject {
	if object == listener {
		return object
	}

	for _, c := range object.children {
		if found := d.findObject(c, listener); found != nil {
			return found
		}
	}

	return nil
}

func (d *MessageDispatcher) sendMessageUp(object *SceneObject, message *dispatch.Message) {
	if object.OnMessage(message) && object.parent != nil {
		d.sendMessageUp(object.parent, message)
	}
}

func (d *MessageDispatcher) DispatchDirectly(listener MessageListener, message *dispatch.Message) {
	listener.OnMessage(message)
}

func newMessageDispatcherForScene(scene *SceneObject) *MessageDispatcher {
	return &MessageDispatcher{
		root: scene,
	}
}

func IsBuiltinEventMessage(msg *dispatch.Message) bool {
	return msg.What == MessageBuiltinEvent
}

func EventFromMessage(msg *dispatch.Message) Event {
	return msg.AnyData.(Event)
}