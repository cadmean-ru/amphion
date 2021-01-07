package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type MouseMover struct {
	object   *engine.SceneObject
	engine   *engine.AmphionEngine
	dragging bool
	mousePos a.IntVector2
}

func (m *MouseMover) OnInit(ctx engine.InitContext) {
	m.object = ctx.GetSceneObject()
	m.engine = ctx.GetEngine()
}

func (m *MouseMover) OnStart() {
	m.engine.BindEventHandler(engine.EventMouseUp, m.handleMouseUp)
}

func (m *MouseMover) handleMouseUp(_ engine.AmphionEvent) bool {
	m.dragging = false
	return true
}

func (m *MouseMover) OnMessage(msg engine.Message) bool {
	if msg.Code != engine.MessageBuiltinEvent || msg.Sender != m.object {
		return true
	}

	event := msg.Data.(engine.AmphionEvent)
	if event.Code != engine.EventMouseDown {
		return true
	}

	m.dragging = true
	m.mousePos = m.engine.GetInputManager().GetMousePosition()
	m.engine.RequestUpdate()

	return true
}

func (m *MouseMover) OnUpdate(_ engine.UpdateContext) {
	if !m.dragging {
		return
	}

	newMousePos := m.engine.GetInputManager().GetMousePosition()
	dPos := newMousePos.Sub(m.mousePos)
	m.mousePos = newMousePos
	m.object.Transform.Position = m.object.Transform.Position.Add(dPos.ToFloat3())
	m.engine.GetMessageDispatcher().DispatchDown(m.object, engine.NewMessage(m, engine.MessageRedraw, nil), engine.MessageMaxDepth)
	m.engine.RequestRendering()
}

func (m *MouseMover) OnStop() {
	m.engine.UnbindEventHandler(engine.EventMouseUp, m.handleMouseUp)
}

func (m *MouseMover) GetName() string {
	return engine.NameOfComponent(m)
}

func NewMouseMover() *MouseMover {
	return &MouseMover{}
}