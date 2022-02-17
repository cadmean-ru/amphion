package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type Padding struct {
	engine.LayoutImpl
	EdgeInsets *common.EdgeInsets `state:"edgeInsets"`
	AutoAdjust bool               `state:"autoAdjust"`
}

func (s *Padding) LayoutChildren() {
	s.LayoutImpl.LayoutChildren()

	innerRect := s.calculateInnerRect()

	for _, c := range s.SceneObject.GetChildren() {
		childRect := c.Transform.Rect()
		if innerRect.IsRectInside(childRect) {
			continue
		}

		if s.AutoAdjust {
			c.Transform.SetRect(innerRect)
		} else {
			childPos := c.Transform.LocalPosition()
			childPos.X = common.ClampFloat32(childPos.X, innerRect.X.Min, innerRect.X.Max)
			childPos.Y = common.ClampFloat32(childPos.Y, innerRect.Y.Min, innerRect.Y.Max)
			childPos.Z = common.ClampFloat32(childPos.Z, innerRect.Z.Min, innerRect.Z.Max)
			c.Transform.SetPosition(childPos)

			childSize := c.Transform.ActualSize()
			childSize.X = common.ClampFloat32(childSize.X, 0, innerRect.X.Max-childPos.X)
			childSize.Y = common.ClampFloat32(childSize.Y, 0, innerRect.Y.Max-childPos.Y)
			childSize.Z = common.ClampFloat32(childSize.Z, 0, innerRect.Z.Max-childPos.Z)
			c.Transform.SetSize(childSize)
		}

		c.Redraw()
	}
}

func (s *Padding) MeasureContents() a.Vector3 {
	return s.LayoutImpl.MeasureContents().Add(a.NewVector3(s.EdgeInsets.Left+s.EdgeInsets.Right, s.EdgeInsets.Top+s.EdgeInsets.Bottom, 0))
}

func (s *Padding) calculateInnerRect() *common.Rect {
	var rect = s.SceneObject.Transform.Rect()
	rect.Move(a.NewVector3(s.EdgeInsets.Left, s.EdgeInsets.Top, 0))
	rect.ShrinkMax(a.NewVector3(s.EdgeInsets.Left+s.EdgeInsets.Right, s.EdgeInsets.Top+s.EdgeInsets.Bottom, 0))
	return rect
}

func (s *Padding) SetPadding(padding *common.EdgeInsets) {
	s.EdgeInsets = padding
	engine.RequestRendering()
}

func (s *Padding) SetAutoAdjust(adjust bool) {
	s.AutoAdjust = adjust
	engine.RequestRendering()
}

func NewPadding(insets *common.EdgeInsets) *Padding {
	return &Padding{
		EdgeInsets: insets,
		AutoAdjust: true,
	}
}
