package engine

// Default component interface implementation
type ComponentImpl struct {
	__name__    string
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

func (c *ComponentImpl) GetName() string {
	return "DefaultName"
}

func (c ComponentImpl) NameOf(outer interface{}) string {
	if c.__name__ != "" {
		return c.__name__
	}

	c.__name__ = NameOfComponent(outer)
	return c.__name__
}
