package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
)

// OnClickListener detects clicks on the scene object and calls the specified engine.EventHandler.
// For click detection to work remember to add a component that implements engine.BoundaryComponent to the scene object.
type OnClickListener struct {
	engine.ComponentImpl
	OnClick engine.EventHandler `state:"onClick"`
}

func (l *OnClickListener) GetName() string {
	return engine.NameOfComponent(l)
}

func (l *OnClickListener) OnMessage(m engine.Message) bool {
	if m.Code != engine.MessageBuiltinEvent || l.OnClick == nil || m.Sender != l.SceneObject {
		return true
	}

	event := m.Data.(engine.AmphionEvent)
	if event.Code != engine.EventMouseDown && event.Code != engine.EventTouchDown {
		return true
	}

	return l.OnClick(event)
}

// Creates a new OnClickListener component with the specified event handler.
func NewOnClickListener(handler engine.EventHandler) *OnClickListener {
	return &OnClickListener{
		OnClick: handler,
	}
}
