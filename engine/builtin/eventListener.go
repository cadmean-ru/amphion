package builtin

import (
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/engine"
)

// EventListener is used to detect and handle scene object events.
// It reacts to the events with the specified code and invokes the specified event handler.
type EventListener struct {
	engine.ComponentImpl
	EventCode int                 `state:"eventCode"`
	Handler   engine.EventHandler `state:"handler"`
}

func (l *EventListener) OnMessage(msg *dispatch.Message) bool {
	if msg.What != engine.MessageBuiltinEvent || l.Handler == nil {
		return true
	}

	event := msg.AnyData.(engine.AmphionEvent)

	if event.Code != l.EventCode || event.Sender != l.SceneObject {
		return true
	}

	l.Handler(event)

	return false
}

// NewEventListener creates a new EventListener, that reacts to the events with the specified code and invokes the specified handler.
func NewEventListener(eventCode int, handler engine.EventHandler) *EventListener {
	return &EventListener{
		EventCode:     eventCode,
		Handler:       handler,
	}
}