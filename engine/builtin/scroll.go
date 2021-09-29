package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
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
	engine.ComponentImpl
	dScrollY, dScrollX float32
	scrollDirectionY, scrollDirectionX ScrollDirection
}

func (s *Scroll) OnStart() {
	engine.BindEventHandler(engine.EventMouseScroll, s.handleScroll)
}

func (s *Scroll) OnStop() {
	engine.UnbindEventHandler(engine.EventMouseScroll, s.handleScroll)
}

//func (s *Scroll) LayoutChildren() {
//
//}

func (s *Scroll) handleScroll(event engine.AmphionEvent) bool {
	scrollAmount := event.Data.(a.Vector2)

	if scrollAmount.Y < 0 {
		s.scrollDirectionY = ScrollUp
	} else if scrollAmount.Y > 0 {
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

	s.dScrollY = scrollAmount.Y
	s.dScrollX = scrollAmount.X

	var realArea = s.measureChildren()
	var sceneRect = engine.GetCurrentScene().Transform.Rect()

	//engine.LogDebug("%+v", realArea)

	var theScrolly, theScrollx float32

	if s.scrollDirectionY == ScrollUp {
		if realArea.Min().Y < 0 {
			mouseScroll := float64(s.dScrollY)
			areaOffset := float64(realArea.Min().Y)
			//engine.LogDebug("%f %f", mouseScroll, areaOffset)
			theScrolly = float32(math.Max(mouseScroll, areaOffset))
		}
	} else if s.scrollDirectionY == ScrollDown {
		if realArea.Max().Y > sceneRect.Y.GetLength() {
			mouseScroll := float64(s.dScrollY)
			areaOffset := float64(realArea.Max().Y - sceneRect.Y.GetLength())
			//engine.LogDebug("%f %f %f", mouseScroll, areaOffset, sceneRect.Y.GetLength())
			theScrolly = float32(math.Min(mouseScroll, areaOffset))
		}
	}

	if s.scrollDirectionX == ScrollLeft {
		if realArea.Min().X < 0 {
			mouseScroll := float64(s.dScrollX)
			areaOffset := float64(realArea.Min().X)
			//engine.LogDebug("%f %f", mouseScroll, areaOffset)
			theScrollx = float32(math.Max(mouseScroll, areaOffset))
		}
	} else if s.scrollDirectionX == ScrollRight {
		if realArea.Max().X > sceneRect.X.GetLength() {
			mouseScroll := float64(s.dScrollX)
			areaOffset := float64(realArea.Max().X - sceneRect.X.GetLength())
			//engine.LogDebug("%f %f %f", mouseScroll, areaOffset, sceneRect.Y.GetLength())
			theScrollx = float32(math.Min(mouseScroll, areaOffset))
		}
	}

	//engine.LogDebug("Scrolling %d %f", s.scrollDirectionY, theScrolly)

	if s.scrollDirectionY != ScrollNone {
		s.SceneObject.Transform.Translate(0, -theScrolly)

		s.SceneObject.ForEachObject(func(object *engine.SceneObject) {
			object.Redraw()
		})
		s.SceneObject.Redraw()
	}

	if s.scrollDirectionX != ScrollNone {
		s.SceneObject.Transform.Translate(-theScrollx, 0)

		s.SceneObject.ForEachObject(func(object *engine.SceneObject) {
			object.Redraw()
		})
		s.SceneObject.Redraw()
	}

	s.scrollDirectionY, s.scrollDirectionX = ScrollNone, ScrollNone

	return true
}

func (s *Scroll) measureChildren() *common.RectBoundary {
	realArea := common.NewRectBoundary(0, 0, 0, 0, -999, 999)

	s.SceneObject.ForEachObject(func(object *engine.SceneObject) {
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

func (s *Scroll) GetName() string {
	return engine.NameOfComponent(s)
}

func NewScroll() *Scroll {
	return &Scroll{}
}
