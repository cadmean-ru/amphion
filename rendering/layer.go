package rendering

import "github.com/cadmean-ru/amphion/common/a"

type Layer struct {
	index       int
	translation a.IntVector3
}

func (l *Layer) GetIndex() int {
	return l.index
}

func (l *Layer) GetTranslation() a.IntVector3 {
	return l.translation
}

func (l *Layer) SetTranslation(t a.IntVector3) {
	l.translation = t
}

func newLayer(index int) *Layer {
	return &Layer{
		index:       index,
		translation: a.IntVector3{},
	}
}
