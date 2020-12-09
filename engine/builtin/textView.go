package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type TextView struct {
	object         *engine.SceneObject
	engine         *engine.AmphionEngine
	Appearance     rendering.Appearance
	TextAppearance rendering.TextAppearance
	renderer       rendering.Renderer
	pId            int64
	text           common.AString
	changed        bool
	children       []engine.ViewComponent
}

func (t *TextView) GetName() string {
	return "TextView"
}

func (t *TextView) OnInit(ctx engine.InitContext) {
	t.object = ctx.GetSceneObject()
	t.engine = ctx.GetEngine()
	t.renderer = ctx.GetRenderer()
	t.children = make([]engine.ViewComponent, 0)
}

func (t *TextView) OnStart() {
	t.pId = t.renderer.AddPrimitive()
}

func (t *TextView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewTextPrimitive(t.text)
	pr.Transform = transformToRenderingTransform(t.object.Transform)
	pr.Appearance = t.Appearance
	pr.TextAppearance = t.TextAppearance
	ctx.GetRenderer().SetPrimitive(t.pId, pr, t.changed || t.engine.IsForcedToRedraw())
	t.changed = false
}

func (t *TextView) ForceRedraw() {
	t.changed = true
	t.engine.GetMessageDispatcher().DispatchDown(t.object, engine.NewMessage(t, engine.MessageRedraw, nil))
}

func (t *TextView) OnStop() {

}

func (t *TextView) SetText(text string) {
	t.text = common.AString(text)
	t.changed = true
}

func (t *TextView) GetText() string {
	return string(t.text)
}

func (t *TextView) OnMessage(message engine.Message) bool {
	if message.Code == engine.MessageRedraw {
		t.changed = true
	}

	return true
}

func NewTextView(text string) *TextView {
	return &TextView{
		text:    common.AString(text),
		changed: true,
	}
}
