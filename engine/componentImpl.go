package engine

// ComponentImpl is a default component interface implementation
type ComponentImpl struct {
	Engine      *AmphionEngine
	SceneObject *SceneObject
	Logger      *Logger
}

func (c *ComponentImpl) OnInit(ctx InitContext) {
	c.Engine = ctx.GetEngine()
	c.SceneObject = ctx.GetSceneObject()
	c.Logger = ctx.GetLogger()
}

func (c *ComponentImpl) OnStart() {

}

func (c *ComponentImpl) OnStop() {

}