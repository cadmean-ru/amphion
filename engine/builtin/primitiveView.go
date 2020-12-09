package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type ShapeView struct {
	object    *engine.SceneObject
	engine    *engine.AmphionEngine
	redraw    bool
	pType     common.AByte
	renderer  rendering.Renderer
	pId       int64
	children  engine.ViewComponent

	Appearance rendering.Appearance
}

func (c *ShapeView) OnInit(ctx engine.InitContext) {
	c.object = ctx.GetSceneObject()
	c.engine = ctx.GetEngine()
	c.renderer = ctx.GetRenderer()
}

func (c *ShapeView) OnStart() {
	c.pId = c.renderer.AddPrimitive()
}

func (c *ShapeView) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(c.pType)
	pr.Transform = transformToRenderingTransform(c.object.Transform)
	pr.Appearance = c.Appearance
	ctx.GetRenderer().SetPrimitive(c.pId, pr, c.redraw || c.engine.IsForcedToRedraw())
	//ctx.GetRenderer().SetPrimitive(c.pId, pr, true)
	c.redraw = false
}

func (c *ShapeView) OnStop() {}

func (c *ShapeView) GetName() string {
	return "ShapeView"
}

func (c *ShapeView) ForceRedraw() {
	c.redraw = true
	c.engine.GetMessageDispatcher().DispatchDown(c.object, engine.NewMessage(c, engine.MessageRedraw, nil))
}

func (c *ShapeView) OnMessage(message engine.Message) bool {
	if message.Code == engine.MessageRedraw {
		c.redraw = true
	}
	return true
}

func NewShapeView(pType common.AByte) *ShapeView {
	if pType < 0 || pType > 5 {
		panic("Invalid primitive type")
	}

	return &ShapeView{
		pType:      pType,
		Appearance: rendering.DefaultAppearance(),
		redraw:     true,
	}
}