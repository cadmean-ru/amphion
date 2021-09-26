package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

// Layout defines interface for layout components.
type Layout interface {
	Component
	Measurable
	LayoutChildren()
}

//Measurable defines interface for views and layouts that is used to get the size of their contents.
//The measured size is used if the wanted size of object has special value a.WrapContent.
type Measurable interface {
	MeasureContents() a.Vector3
}

//LayoutImpl is the default implementation of Layout interface.
//It handles the special values in objects' sizes and positions.
type LayoutImpl struct {
	ComponentImpl
}

func (l *LayoutImpl) LayoutChildren() {
	l.FirstPass(l.SceneObject)
	l.SecondPass(l.SceneObject)
}

func (l *LayoutImpl) FirstPass(object *SceneObject) {
	for _, c := range object.children {
		if c.HasLayout() {
			l.MeasureObject(c)
			l.FinalizeAbsoluteDimensions(c)
			l.FinalizeWrapContentDimensions(c)
			continue
		}

		l.FirstPass(c)
	}

	l.MeasureObject(object)
	l.FinalizeAbsoluteDimensions(object)
	l.FinalizeWrapContentDimensions(object)
}

func (l *LayoutImpl) SecondPass(object *SceneObject) {
	l.FinalizeMatchParentDimensions(object)
	l.FinalizePosition(object)

	for _, c := range object.children {
		if c.HasLayout() {
			l.FinalizeMatchParentDimensions(c)
			continue
		}

		l.SecondPass(c)
	}
}

func (l *LayoutImpl) MeasureObject(object *SceneObject) {
	if !l.ObjectNeedsToBeMeasured(object) {
		return
	}

	if object.HasLayout() {
		object.Transform.measuredSize = object.layout.component.(Layout).MeasureContents()
		return
	}

	size := a.ZeroVector()

	if view := object.GetView(); view != nil {
		contents := view.MeasureContents()
		if contents.X > size.X {
			size.X = contents.X
		}
		if contents.Y > size.Y {
			size.Y = contents.Y
		}
		if contents.Z > size.Z {
			size.Z = contents.Z
		}
	}

	for _, c := range object.children {
		childRect := common.NewRectBoundaryFromPositionAndSize(c.Transform.positionForMeasurement(), c.Transform.actualSize)
		if childRect.X.Max > size.X {
			size.X = childRect.X.Max
		}
		if childRect.Y.Max > size.Y {
			size.Y = childRect.Y.Max
		}
		if childRect.Y.Max > size.Y {
			size.Y = childRect.Y.Max
		}
	}

	object.Transform.measuredSize = size
}

func (l *LayoutImpl) ObjectNeedsToBeMeasured(object *SceneObject) bool {
	return object.Transform.size.X == a.WrapContent ||
		object.Transform.size.Y == a.WrapContent ||
		object.Transform.size.Z == a.WrapContent
}

func (l *LayoutImpl) FinalizeAbsoluteDimensions(object *SceneObject) {
	if !IsSpecialTransformValue(object.Transform.size.X) {
		object.Transform.actualSize.X = object.Transform.size.X
	}
	if !IsSpecialTransformValue(object.Transform.size.Y) {
		object.Transform.actualSize.Y = object.Transform.size.Y
	}
	if !IsSpecialTransformValue(object.Transform.size.Z) {
		object.Transform.actualSize.Z = object.Transform.size.Z
	}
}

func (l *LayoutImpl) FinalizeWrapContentDimensions(object *SceneObject) {
	if object.Transform.size.X == a.WrapContent {
		object.Transform.actualSize.X = object.Transform.measuredSize.X
	}
	if object.Transform.size.Y == a.WrapContent {
		object.Transform.actualSize.Y = object.Transform.measuredSize.Y
	}
	if object.Transform.size.Z == a.WrapContent {
		object.Transform.actualSize.Z = object.Transform.measuredSize.Z
	}
}

func (l *LayoutImpl) FinalizeMatchParentDimensions(object *SceneObject) {
	if object.parent == nil {
		return
	}

	if object.Transform.size.X == a.MatchParent {
		object.Transform.actualSize.X = object.parent.Transform.actualSize.X
	}
	if object.Transform.size.Y == a.MatchParent {
		object.Transform.actualSize.Y = object.parent.Transform.actualSize.Y
	}
	if object.Transform.size.Z == a.MatchParent {
		object.Transform.actualSize.Z = object.parent.Transform.actualSize.Z
	}
}

func (l *LayoutImpl) FinalizePosition(object *SceneObject) {
	if object.parent != nil {
		pos := object.Transform.position
		pr := object.parent.Transform.Rect()
		if pos.X == a.CenterInParent {
			pos.X = pr.X.GetLength() / 2
		}
		if pos.Y == a.CenterInParent {
			pos.Y = pr.Y.GetLength() / 2
		}
		if pos.Z == a.CenterInParent {
			pos.Z = pr.Z.GetLength() / 2
		}
		object.Transform.actualPosition = pos
	} else {
		object.Transform.actualPosition = object.Transform.positionForMeasurement()
	}
}

func (l *LayoutImpl) MeasureContents() a.Vector3 {
	return a.ZeroVector()
}
