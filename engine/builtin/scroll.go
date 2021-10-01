package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/common/shape"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"math"
)

type ScrollDirection byte

const (
	ScrollUp ScrollDirection = iota
	ScrollDown
	ScrollRight
	ScrollLeft
	ScrollNone
)

type Scroll struct {
	engine.LayoutImpl
	scrollDirectionY, scrollDirectionX ScrollDirection
	dOffsetY, dOffsetX float32
}

func (s *Scroll) OnStart() {
	engine.BindEventHandler(engine.EventMouseScroll, s.handleScroll)
}

func (s *Scroll) OnUpdate(_ engine.UpdateContext) {
	s.SceneObject.GetRenderingNode().SetClipArea2D(rendering.NewClipArea2D(shape.Rectangle, s.SceneObject.Transform.GlobalRect()))
}

func (s *Scroll) OnStop() {
	engine.UnbindEventHandler(engine.EventMouseScroll, s.handleScroll)
	s.SceneObject.GetRenderingNode().RemoveClipArea2D()
}

func (s *Scroll) LayoutChildren() {
	s.LayoutImpl.LayoutChildren()

	if s.scrollDirectionY == ScrollNone && s.scrollDirectionX == ScrollNone {
		return
	}

	if s.scrollDirectionY != ScrollNone {
		s.SceneObject.ForEachChild(func(child *engine.SceneObject) {
			child.Transform.Translate(0, s.dOffsetY)
		})
	}

	if s.scrollDirectionX != ScrollNone {
		s.SceneObject.ForEachChild(func(child *engine.SceneObject) {
			child.Transform.Translate(s.dOffsetX, 0)
		})
	}

	engine.GetMessageDispatcher().DispatchDown(s.SceneObject, dispatch.NewMessageFrom(s.SceneObject, engine.MessageRedraw), engine.MessageMaxDepth)
	//engine.ForceAllViewsRedraw()

	s.dOffsetY, s.dOffsetX = 0, 0
	s.scrollDirectionY, s.scrollDirectionX = ScrollNone, ScrollNone
}

func (s *Scroll) IsPointInside(point a.Vector3) bool {
	return s.SceneObject.Transform.GlobalRect().IsPointInside(point)
}

func (s *Scroll) IsPointInside2D(point a.Vector3) bool {
	return s.SceneObject.Transform.GlobalRect().IsPointInside2D(point)
}

func (s *Scroll) IsSolid() bool {
	return false
}

func (s *Scroll) handleScroll(event engine.AmphionEvent) bool {
	if !s.SceneObject.IsHovered() {
		return true
	}

	scrollAmount := event.Data.(a.Vector2)

	if scrollAmount.Y > 0 {
		s.scrollDirectionY = ScrollUp
	} else if scrollAmount.Y < 0 {
		s.scrollDirectionY = ScrollDown
	} else {
		s.scrollDirectionY = ScrollNone
	}

	if scrollAmount.X > 0 {
		s.scrollDirectionX = ScrollRight
	} else if scrollAmount.X < 0 {
		s.scrollDirectionX = ScrollLeft
	} else {
		s.scrollDirectionX = ScrollNone
	}

	//engine.LogDebug("Scroll: %v %d %d", scrollAmount, s.scrollDirectionX, s.scrollDirectionY)

	dScrollY := scrollAmount.Y
	dScrollX := scrollAmount.X

	var realArea = s.measureChildren()
	var sceneRect = s.SceneObject.Transform.Rect()

	//engine.LogDebug("%+v", realArea)

	var theScrolly, theScrollx float32

	if s.scrollDirectionY == ScrollUp {
		if realArea.Min().Y < viewPort.Min().Y {
			mouseScroll := float64(dScrollY)
			areaOffset := float64(viewPort.Min().Y - realArea.Min().Y)
			//engine.LogDebug("%f %f", mouseScroll, areaOffset)
			theScrolly = float32(math.Min(mouseScroll, areaOffset))
		}
	} else if s.scrollDirectionY == ScrollDown {
		if realArea.Max().Y > viewPort.Max().Y {
			mouseScroll := float64(dScrollY)
			areaOffset := float64(viewPort.Max().Y - realArea.Max().Y)
			//engine.LogDebug("%f %f %f", mouseScroll, areaOffset, viewPort.Y.GetLength())
			theScrolly = float32(math.Max(mouseScroll, areaOffset))
		}
	}

	if s.scrollDirectionX == ScrollLeft {
		if realArea.Max().X > viewPort.Max().X {
			mouseScroll := float64(dScrollX)
			areaOffset := float64(viewPort.Max().X - realArea.Max().X)
			//engine.LogDebug("%f %f %f", mouseScroll, areaOffset, viewPort.Y.GetLength())
			theScrollx = float32(math.Max(mouseScroll, areaOffset))
		}
	} else if s.scrollDirectionX == ScrollRight {
		if realArea.Min().X < viewPort.Min().X {
			mouseScroll := float64(dScrollX)
			areaOffset := float64(viewPort.Min().X - realArea.Min().X)
			//engine.LogDebug("%f %f", mouseScroll, areaOffset)
			theScrollx = float32(math.Min(mouseScroll, areaOffset))
		}
	}

	//engine.LogDebug("Scrolling %d %f", s.scrollDirectionY, theScrolly)

	s.dOffsetY, s.dOffsetX = theScrolly, theScrollx

	return true
}

func (s *Scroll) measureChildren() *common.RectBoundary {
	realArea := s.SceneObject.Transform.GlobalRect()

	s.SceneObject.ForEachChild(func(object *engine.SceneObject) {
		rect := object.Transform.GlobalRect()
		if rect.X.Min < realArea.X.Min {
			realArea.X.Min = rect.X.Min
		}
		if rect.X.Max > realArea.X.Max {
			realArea.X.Max = rect.X.Max
		}
		if rect.Y.Min < realArea.Y.Min {
			realArea.Y.Min = rect.Y.Min
		}
		if rect.Y.Max > realArea.Y.Max {
			realArea.Y.Max = rect.Y.Max
		}
	})

	return realArea
}

func NewScroll() *Scroll {
	return &Scroll{}
}
