package engine

import "fmt"

const (
	MessageRedraw = iota
	MessageBuiltinEvent
	MessageUpdate
	MessageUpdateStop
)

const MessageMaxDepth = -1

type Message struct {
	Sender interface{}
	Code   byte
	Data   interface{}
}

func (m Message) ToString() string {
	return fmt.Sprintf("Message from %t. Code: %d", m.Sender, m.Code)
}

func NewMessage(sender interface{}, code byte, data interface{}) Message {
	return Message{
		Sender: sender,
		Code:   code,
		Data:   data,
	}
}

type MessageListenerComponent interface {
	Component
	MessageListener
}

// Interface for receiving messages from dispatcher
type MessageListener interface {
	// Receives a message and returns whether to continue message propagation or not
	OnMessage(message Message) bool
}


type messageListenerContainer struct {
	listener  MessageListener
	parent    *messageListenerContainer
	children  []*messageListenerContainer
}

func (c *messageListenerContainer) findListenerContainer(listener MessageListener) *messageListenerContainer {
	if c.listener == listener {
		return c
	}

	for _, c1 := range c.children {
		if c2 := c1.findListenerContainer(listener); c2 != nil {
			return c2
		}
	}

	return nil
}

func (c *messageListenerContainer) sendMessageDown(message Message, counter int) {
	if c.listener.OnMessage(message) && (counter == MessageMaxDepth || counter > 0) {
		for _, c1 := range c.children {
			c1.sendMessageDown(message, counter-1)
		}
	}
}

func (c *messageListenerContainer) sendMessageUp(message Message) {
	if c.listener.OnMessage(message) && c.parent != nil {
		c.parent.sendMessageUp(message)
	}
}


type MessageDispatcher struct {
	root    *messageListenerContainer
}

// Adds listener into the message tree
func (d *MessageDispatcher) AddListener(listener MessageListener, parent MessageListener) {
	c := d.root.findListenerContainer(parent)
	if c == nil {
		return
	}

	c.children = append(c.children, &messageListenerContainer{
		listener: listener,
		parent:   c,
		children: make([]*messageListenerContainer, 0, 1),
	})
}

// Removes listener from message tree. You cannot remove the root element.
func (d *MessageDispatcher) RemoveListener(listener MessageListener) {
	c := d.root.findListenerContainer(listener)
	if c == nil {
		return
	}

	if c.parent == nil {
		return
	}

	p := c.parent
	var removeIndex int
	for i, c1 := range p.children {
		if c1.listener == listener {
			removeIndex = i
			break
		}
	}

	if len(p.children) > 0 {
		c.children[removeIndex] = c.children[len(c.children)-1]
		c.children = c.children[:len(c.children)-1]
	}
}

// Sends message from the root down in the message tree
func (d *MessageDispatcher) Dispatch(message Message, maxDepth int) {
	if !d.shouldSendMessage(message) {
		return
	}

	d.root.sendMessageDown(message, maxDepth)
}

// Sends this message up in the message tree
func (d *MessageDispatcher) DispatchDown(from MessageListener, message Message, maxDepth int) {
	if !d.shouldSendMessage(message) {
		return
	}

	c := d.root.findListenerContainer(from)
	if c == nil {
		return
	}
	c.sendMessageDown(message, maxDepth)
}

// Sends specified message down in the message tree.
func (d *MessageDispatcher) DispatchUp(from MessageListener, message Message) {
	if !d.shouldSendMessage(message) {
		return
	}

	c := d.root.findListenerContainer(from)
	if c == nil {
		return
	}
	c.sendMessageUp(message)
}

func (d *MessageDispatcher) DispatchDirectly(listener MessageListener, message Message) {
	if !d.shouldSendMessage(message) {
		return
	}

	listener.OnMessage(message)
}

func (d *MessageDispatcher) shouldSendMessage(message Message) bool {
	if message.Code == MessageRedraw && instance.state == StateRendering {
		return false
	}
	return true
}

func NewMessageDispatcher(root MessageListener) *MessageDispatcher {
	return &MessageDispatcher{
		root: &messageListenerContainer{
			listener: root,
			parent:   nil,
			children: make([]*messageListenerContainer, 0, 1),
		},
	}
}

func newMessageDispatcherForScene(scene *SceneObject) *MessageDispatcher {
	var dis = &MessageDispatcher{
		root: &messageListenerContainer{
			listener: instance,
			parent:   nil,
			children: make([]*messageListenerContainer, 0),
		},
	}

	addSceneObject(dis, scene, instance)

	return dis
}

func addSceneObject(dis *MessageDispatcher, obj *SceneObject, p MessageListener) {
	dis.AddListener(obj, p)
	for _, o := range obj.children {
		addSceneObject(dis, o, obj)
	}
}