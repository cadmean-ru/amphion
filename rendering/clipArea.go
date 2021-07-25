package rendering

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/shape"
)

type ClipArea2D struct {
	Rect  *common.RectBoundary
	Shape shape.Kind
}

func NewClipArea2DEmpty() *ClipArea2D {
	return &ClipArea2D{
		Rect:  common.NewRectBoundary(0, 0, 0, 0, 0, 0),
		Shape: shape.Empty,
	}
}

func NewClipArea2D(kind shape.Kind, rect *common.RectBoundary) *ClipArea2D {
	return &ClipArea2D{
		Rect:  rect,
		Shape: kind,
	}
}

type ClipStack2D struct {
	clips []*ClipArea2D
	bottom *ClipArea2D
}

func (s *ClipStack2D) Size() int {
	return len(s.clips)
}

func (s *ClipStack2D) Push(area *ClipArea2D) {
	s.clips = append(s.clips, area)
}

func (s *ClipStack2D) Peek() *ClipArea2D {
	if len(s.clips) == 0 {
		return s.bottom
	}

	return s.clips[len(s.clips)-1]
}

func (s *ClipStack2D) Pop() {
	if len(s.clips) == 0{
		return
	}

	s.clips = s.clips[:len(s.clips)-1]
}

func (s *ClipStack2D) Clear() {
	s.clips = make([]*ClipArea2D, 0, 10)
}

func NewClipStack2D() *ClipStack2D {
	return &ClipStack2D{
		clips:  make([]*ClipArea2D, 0, 10),
		bottom: NewClipArea2DEmpty(),
	}
}