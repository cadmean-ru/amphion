package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
)

type Margin struct {
	engine.ComponentImpl
	margin        *common.EdgeInsets
	marginObj     *engine.SceneObject
	marginPadding *Padding
}

func (s *Margin) OnInit(ctx engine.InitContext) {
	s.ComponentImpl.OnInit(ctx)
}

func (s *Margin) OnStart() {
	currentParent := s.SceneObject.GetParent()

	s.marginObj = engine.NewSceneObject("beu")
	s.marginPadding = NewPadding(s.margin)
	s.marginObj.AddComponent(s.marginPadding)

	s.SceneObject.SetParent(s.marginObj)
	s.marginObj.SetParent(currentParent)

	s.marginObj.Transform.SetSize(
		s.SceneObject.Transform.ActualSize().X+s.margin.Left+s.margin.Right,
		s.SceneObject.Transform.ActualSize().Y+s.margin.Top+s.margin.Bottom,
	)
}

func (s *Margin) SetMargin(margin *common.EdgeInsets) {
	s.margin = margin
	s.marginObj.Transform.SetSize(
		s.SceneObject.Transform.ActualSize().X+s.margin.Left+s.margin.Right,
		s.SceneObject.Transform.ActualSize().Y+s.margin.Top+s.margin.Bottom,
	)
	s.marginPadding.SetPadding(s.margin)
}

func NewMargin(margin *common.EdgeInsets) *Margin {
	return &Margin{
		margin: margin,
	}
}
