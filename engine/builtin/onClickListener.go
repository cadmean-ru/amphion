package builtin

import (
	"github.com/cadmean-ru/amphion/engine"
)

type OnClickListener struct {
	OnClick engine.EventHandler
	object  *engine.SceneObject
}

func (l *OnClickListener) GetName() string {
	return "OnClickListener"
}

func (l *OnClickListener) OnInit(ctx engine.InitContext) {
	l.object = ctx.GetSceneObject()
}

func (l *OnClickListener) OnStart() {
}

func (l *OnClickListener) OnStop() {

}

func (l *OnClickListener) OnMessage(m engine.Message) bool {
	if m.Code != engine.MessageBuiltinEvent || l.OnClick == nil || m.Sender != l.object {
		return true
	}

	event := m.Data.(engine.AmphionEvent)
	if event.Code != engine.EventMouseDown {
		return true
	}

	return l.OnClick(m.Data.(engine.AmphionEvent))
}

func NewOnClickListener(handler engine.EventHandler) *OnClickListener {
	return &OnClickListener{
		OnClick: handler,
	}
}
