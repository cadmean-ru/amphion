package scenes

import (
	. "github.com/cadmean-ru/amphion/engine"
	. "github.com/cadmean-ru/amphion/engine/builtin"
)

type Mover struct {
	UpdatingComponentImpl
	shape *ShapeView
	dir  bool
}

func (m *Mover) OnInit(ctx InitContext) {
	m.UpdatingComponentImpl.OnInit(ctx)
	m.shape = GetShapeView(m.SceneObject, true)
}

func (m *Mover) OnStart() {
	LogInfo("Start")
}

func (m *Mover) OnUpdate(ctx UpdateContext) {
	maxX := GetCurrentScene().Transform.WantedSize().X - m.SceneObject.Transform.WantedSize().X

	if m.SceneObject.Transform.WantedPosition().X <= 0 {
		m.dir = true
	} else if m.SceneObject.Transform.WantedPosition().X >= maxX {
		m.SceneObject.Transform.SetPosition(maxX, m.SceneObject.Transform.WantedPosition().Y)
		m.dir = false
	}
	dX := 100 * ctx.DeltaTime
	if m.dir {
		m.SceneObject.Transform.Translate(dX, 0)
	} else {
		m.SceneObject.Transform.Translate(-dX, 0)
	}
	m.SceneObject.Redraw()
	RequestRendering()
}

func (m *Mover) OnStop() {
	LogInfo("Stop")
}
