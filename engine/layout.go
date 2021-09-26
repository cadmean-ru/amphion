package engine

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
)

// Layout defines layout components.
type Layout interface {
	Component
	Measurable
	LayoutChildren()
}

type Measurable interface {
	MeasureContents() a.Vector3
}

// LayoutResponder defines interface for components that respond to layout changes.
//type LayoutResponder interface {
//	Component
//	OnMeasure(ctx *LayoutContext)
//	OnLayout(ctx *LayoutContext)
//}
//
//type LayoutContext struct {
//
//}

// LayoutImpl TODO: optimize maybe
type LayoutImpl struct {
	ComponentImpl
}

func (l *LayoutImpl) LayoutChildren() {
	//l.setInLayout(l.SceneObject)
	//l.measureObjects(l.SceneObject)
	//l.absoluteSizes(l.SceneObject)
	//l.wrapContents(l.SceneObject)
	//l.matchParents(l.SceneObject)
	l.firstPass(l.SceneObject)
	l.secondPass(l.SceneObject)
}

//func (l *LayoutImpl) setInLayout(object *SceneObject) {
//	object.Transform.inLayout = true
//	for _, c := range object.children {
//		l.setInLayout(c)
//	}
//}

func (l *LayoutImpl) firstPass(object *SceneObject) {
	for _, c := range object.children {
		if c.HasLayout() {
			l.measureObject(c)
			l.finalizeAbsoluteDimensions(c)
			l.finalizeWrapContentDimensions(c)
			continue
		}

		l.firstPass(c)
	}

	l.measureObject(object)
	l.finalizeAbsoluteDimensions(object)
	l.finalizeWrapContentDimensions(object)
}

func (l *LayoutImpl) secondPass(object *SceneObject) {
	l.finalizeMatchParentDimensions(object)
	l.finalizePosition(object)

	for _, c := range object.children {
		if c.HasLayout() {
			l.finalizeMatchParentDimensions(c)
			continue
		}

		l.secondPass(c)
	}
}

//func (l *LayoutImpl) thirdPass(object *SceneObject) {
//
//}

//func (l *LayoutImpl) measureObjects(object *SceneObject) {
//	for _, c := range object.children {
//		if c.HasLayout() {
//			c.Transform.measuredSize = c.layout.component.(Layout).MeasureContents()
//			continue
//		}
//
//		l.measureObjects(c)
//	}
//
//	l.measureObject(object)
//}

func (l *LayoutImpl) measureObject(object *SceneObject) {
	if object.HasLayout() {
		object.Transform.measuredSize = object.layout.component.(Layout).MeasureContents()
		return
	}

	size := a.ZeroVector()

	for _, v := range object.GetViews() {
		contents := v.MeasureContents()
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

//func (l *LayoutImpl) absoluteSizes(object *SceneObject) {
//	for _, c := range object.children {
//		if c.HasLayout() {
//			l.finalizeAbsoluteDimensions(c)
//			continue
//		}
//
//		l.absoluteSizes(c)
//	}
//
//	l.finalizeAbsoluteDimensions(object)
//}

func (l *LayoutImpl) finalizeAbsoluteDimensions(object *SceneObject) {
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

//func (l *LayoutImpl) wrapContents(object *SceneObject) {
//	for _, c := range object.children {
//		if c.HasLayout() {
//			l.finalizeWrapContentDimensions(c)
//			continue
//		}
//
//		l.wrapContents(c)
//	}
//
//	l.finalizeWrapContentDimensions(object)
//}

func (l *LayoutImpl) finalizeWrapContentDimensions(object *SceneObject) {
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

//func (l *LayoutImpl) matchParents(object *SceneObject) {
//	l.finalizeMatchParentDimensions(object)
//
//	for _, c := range object.children {
//		if c.HasLayout() {
//			l.finalizeMatchParentDimensions(c)
//			continue
//		}
//
//		l.matchParents(c)
//	}
//}

func (l *LayoutImpl) finalizeMatchParentDimensions(object *SceneObject) {
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

func (l *LayoutImpl) finalizePosition(object *SceneObject) {
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
