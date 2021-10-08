package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/engine"
)

type MouseMover struct {
	engine.UpdatingComponentImpl
	dragging bool
	mousePos a.IntVector2
}

func (m *MouseMover) OnStart() {
	engine.BindEventHandler(engine.EventMouseDown, m.handleMouseDown)
	engine.BindEventHandler(engine.EventTouchDown, m.handleMouseDown)
	engine.BindEventHandler(engine.EventMouseUp, m.handleMouseUp)
	engine.BindEventHandler(engine.EventTouchUp, m.handleMouseUp)
}

func (m *MouseMover) handleMouseDown(e engine.Event) bool {
	eventData := e.Data.(engine.MouseEventData)
	if eventData.SceneObject != m.SceneObject {
		return true
	}

	m.dragging = true
	m.mousePos = m.Engine.GetInputManager().GetCursorPosition()
	m.Engine.RequestUpdate()
	return true
}

func (m *MouseMover) handleMouseUp(_ engine.Event) bool {
	m.dragging = false
	return true
}

func (m *MouseMover) OnUpdate(_ engine.UpdateContext) {
	if !m.dragging {
		return
	}

	newMousePos := m.Engine.GetInputManager().GetCursorPosition()
	dPos := newMousePos.Sub(m.mousePos)
	m.mousePos = newMousePos
	m.SceneObject.Transform.Translate(dPos.ToFloat3())
	m.Engine.GetMessageDispatcher().DispatchDown(m.SceneObject, dispatch.NewMessageFromWithAnyData(m, engine.MessageRedraw, nil), engine.MessageMaxDepth)
	m.Engine.RequestRendering()
}

func (m *MouseMover) OnStop() {
	engine.UnbindEventHandler(engine.EventMouseUp, m.handleMouseUp)
	engine.UnbindEventHandler(engine.EventTouchUp, m.handleMouseUp)
	engine.UnbindEventHandler(engine.EventMouseDown, m.handleMouseDown)
	engine.UnbindEventHandler(engine.EventTouchDown, m.handleMouseDown)
}

func NewMouseMover() *MouseMover {
	return &MouseMover{}
}