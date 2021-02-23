package rendering

import "github.com/cadmean-ru/amphion/common/a"

type Layer struct {
	primitives []int
	translation a.IntVector3
}

func (l *Layer) GetTranslation() a.IntVector3 {
	return l.translation
}

func (l *Layer) SetTranslation(t a.IntVector3) {
	l.translation = t
}

func (l *Layer) addPrimitiveId(id int) {
	l.primitives = append(l.primitives, id)
}

func (l *Layer) removePrimitiveId(id int) {
	index := -1

	for i, pi := range l.primitives {
		if pi == id {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	l.primitives[index] = l.primitives[len(l.primitives) - 1]
	l.primitives = l.primitives[:len(l.primitives) - 1]
}

func newLayer() *Layer {
	return &Layer{
		primitives:  make([]int, 0, 10),
		translation: a.IntVector3{},
	}
}