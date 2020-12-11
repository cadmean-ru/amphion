package engine

import (
	"reflect"
)

const (
	EventUpdate      = -1
	EventRender      = -2
	EventCloseScene  = -3
	EventMouseDown   = -4
	EventDoubleClick = -5
	EventNavigate    = -6
	EventStop        = -7
	EventKeyDown     = -8
	EventMouseUp     = -9
)

type AmphionEvent struct {
	Sender interface{}
	Code   int
	Data   interface{}
}

func NewAmphionEvent(from interface{}, code int, data interface{}) AmphionEvent {
	return AmphionEvent{
		Sender: from,
		Code:   code,
		Data:   data,
	}
}

type EventHandler func(event AmphionEvent) bool

type EventBinder struct {
	handlers map[int][]EventHandler
}

func (b *EventBinder) Bind(code int, handler EventHandler) {
	if _, ok := b.handlers[code]; !ok {
		b.handlers[code] = make([]EventHandler, 1, 10)
		b.handlers[code][0] = handler
	} else {
		b.handlers[code] = append(b.handlers[code], handler)
	}
}

func (b *EventBinder) Unbind(code int, handler EventHandler) {
	if handlers, ok := b.handlers[code]; !ok {
		var index = -1
		for i, h := range handlers {
			p1 := reflect.ValueOf(h).Pointer()
			p2 := reflect.ValueOf(handler).Pointer()
			if p1 == p2 {
				index = i
				break
			}
		}
		if index != -1 {
			handlers[index] = handlers[len(handlers) - 1]
			handlers = handlers[:len(handlers) - 1]
		}
	}
}

func (b *EventBinder) GetHandlers(code int) []EventHandler {
	if h, ok := b.handlers[code]; ok {
		return h
	}

	return make([]EventHandler, 0)
}

func (b *EventBinder) InvokeHandlers(event AmphionEvent) {
	for _, h := range b.GetHandlers(event.Code) {
		if !h(event) {
			break
		}
	}
}

func newEventBinder() *EventBinder {
	return &EventBinder{
		handlers: make(map[int][]EventHandler),
	}
}

type KeyEvent struct {
	Key, Code string
}