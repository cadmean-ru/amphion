package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
)

type MouseMover struct {
	object   *engine.SceneObject
	engine   *engine.AmphionEngine
	dragging bool
	mousePos common.IntVector3
}

func (m *MouseMover) OnInit(ctx engine.InitContext) {
	m.object = ctx.GetSceneObject()
	m.engine = ctx.GetEngine()
}

func (m *MouseMover) OnStart() {
	m.engine.BindEventHandler(engine.EventMouseUp, m.handleMouseUp)
}

func (m *MouseMover) handleMouseUp(event engine.AmphionEvent) bool {
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
	m.mousePos = engine.GetMousePosition()
	m.engine.RequestUpdate()

	return true
}

func (m *MouseMover) OnUpdate(ctx engine.UpdateContext) {
	if !m.dragging {
		return
	}

	newMousePos := engine.GetMousePosition()
	dPos := newMousePos.Sub(m.mousePos)
	m.mousePos = newMousePos
	m.object.Transform.Position = m.object.Transform.Position.Add(dPos.ToFloat())
	m.engine.GetMessageDispatcher().DispatchDown(m.object, engine.NewMessage(m, engine.MessageRedraw, nil))
	m.engine.RequestRendering()
}

func (m *MouseMover) OnStop() {
	m.engine.UnbindEventHandler(engine.EventMouseUp, m.handleMouseUp)
}

func (m *MouseMover) GetName() string {
	return "MouseMover"
}

func NewMouseMover() *MouseMover {
	return &MouseMover{}
}