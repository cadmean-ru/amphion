package main

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
)

type Mover struct {
	obj  *engine.SceneObject
	log  *engine.Logger
	prc  *builtin.ShapeView
	eng  *engine.AmphionEngine
	dir  bool
}

func (m *Mover) OnInit(ctx engine.InitContext) {
	m.obj = ctx.GetSceneObject()
	m.log = ctx.GetLogger()
	m.prc = m.obj.GetComponentByName("ShapeView").(*builtin.ShapeView)
	m.eng = ctx.GetEngine()
}

func (m *Mover) OnStart() {
	m.log.Info(m, "Start")
}

func (m *Mover) OnUpdate(ctx engine.UpdateContext) {
	maxX := m.eng.GetCurrentScene().Transform.Size.X - m.obj.Transform.Size.X

	if m.obj.Transform.Position.X <= 0 {
		m.dir = true
	} else if m.obj.Transform.Position.X >= maxX {
		m.obj.Transform.Position.X = maxX
		m.dir = false
	}
	dX := 100 * ctx.DeltaTime
	if m.dir {
		m.obj.Transform.Position.X += dX
	} else {
		m.obj.Transform.Position.X -= dX
	}
	m.prc.Redraw()
	m.eng.RequestRendering()
}

func (m *Mover) OnStop() {
	m.log.Info(m, "Stop")
}

func (m *Mover) GetName() string {
	return engine.NameOfComponent(m)
}
