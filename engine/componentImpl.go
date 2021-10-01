package engine

// ComponentImpl is a default Component interface implementation.
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

//UpdatingComponentImpl is a default implementation of interface UpdatingComponent.
type UpdatingComponentImpl struct {
	ComponentImpl
}

func (c *UpdatingComponentImpl) OnUpdate(_ UpdateContext) {

}

func (c *UpdatingComponentImpl) OnLateUpdate(_ UpdateContext) {

}