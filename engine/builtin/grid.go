package builtin

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type GridLayout struct {
	engine.ComponentImpl
	Rows int `state:"rows"`
	Cols int `state:"cols"`
	RowPadding int `state:"rowPadding"`
	ColPadding int `state:"colPadding"`
}

func (l *GridLayout) LayoutChildren() {
	fmt.Println("layouring")
	i := 0
	children := l.SceneObject.GetChildren()
	colWidth := l.SceneObject.Transform.Size.X / float32(l.Cols)
	fRowPadding := float32(l.RowPadding)
	fColPadding := float32(l.ColPadding)
	var y = fRowPadding
	for r := 0; (l.Rows > 0 && r < l.Rows) || (l.Rows == 0 && i < len(children)); r++ {
		var maxRowHeight float32 = 0

		for c := 0; c < l.Cols; c++ {
			if i == len(children) {
				return
			}

			child := children[i]
			chSize := child.Transform.Size

			x := float32(c) * colWidth + fColPadding

			if chSize.Y > maxRowHeight {
				maxRowHeight = chSize.Y
			}

			newPos := a.NewVector3(x, y, child.Transform.Position.Z)
			child.Transform.Position = newPos

			child.Transform.Size = a.NewVector3(colWidth - fColPadding * 2, chSize.Y, 0)

			i++
		}

		y += maxRowHeight + fRowPadding * 2
	}
}

func (l *GridLayout) GetName() string {
	return engine.NameOfComponent(l)
}

func NewGridLayout() *GridLayout {
	return &GridLayout{
		Rows: 0,
		Cols: 1,
	}
}