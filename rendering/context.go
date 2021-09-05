package rendering

import "github.com/cadmean-ru/amphion/common/a"

type PrimitiveRenderingContext struct {
	Renderer      *ARenderer
	Primitive     Primitive
	PrimitiveKind byte
	PrimitiveId   int
	State         interface{}
	Redraw        bool
	ClipArea2D    *ClipArea2D
	Projection    a.Matrix4
}