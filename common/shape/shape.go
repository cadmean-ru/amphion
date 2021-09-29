package shape

import "github.com/cadmean-ru/amphion/common"

type Shape interface {
	common.Boundary
	Kind() Kind
}

func New(kind Kind, boundary *common.RectBoundary) Shape {
	switch kind {
	case Ellipse:
		return NewEllipseShape(boundary)
	case Rectangle:
		return NewRectangleShape(boundary)
	case Triangle:
		return NewTriangleShape(boundary)
	default:
		return EmptyShape{}
	}
}