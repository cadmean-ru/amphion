package engine

import (
	"reflect"
)

const (
	EventEmpty = -1 - iota
	EventMouseDown
	EventNavigate
	EventStop
	EventKeyDown
	EventMouseUp
	EventAppHide
	EventAppShow
	EventClipboardPaste
	EventClipboardCopy
	EventMouseIn
	EventMouseOut
	EventDropFile
	EventMouseMove
	EventAppLoaded
	EventFocusGain
	EventFocusLose
	EventMouseScroll
	EventTouchDown
	EventTouchUp
	EventTouchMove
	EventKeyUp
	EventTextInput
	EventOrientationChange
)

type Event struct {
	Sender interface{}
	Code   int
	Data   interface{}
}

func (e Event) StringData() string {
	return e.Data.(string)
}

func (e Event) MouseEventData() MouseEventData {
	return e.Data.(MouseEventData)
}

func (e Event) KeyEventData() KeyEventData {
	return e.Data.(KeyEventData)
}

func NewAmphionEvent(from interface{}, code int, data interface{}) Event {
	return Event{
		Sender: from,
		Code:   code,
		Data:   data,
	}
}

// EventHandler handles event. Returns whether to continue event propagation or not.
type EventHandler func(event Event) bool

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

func (b *EventBinder) InvokeHandlers(event Event) {
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

// Represents file data returned as a result of the user selecting file or dragging it.
type InputFileData struct {
	Name string
	Data []byte
	Mime string
}