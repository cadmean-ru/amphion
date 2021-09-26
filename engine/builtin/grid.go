package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/require"
	"github.com/cadmean-ru/amphion/engine"
	"math"
)

type GridRowDefinition struct {
	Height     float32
	maxHeight  float32
	fillHeight  float32
}

func (r *GridRowDefinition) ToMap() a.SiMap {
	return a.SiMap{
		"height": r.Height,
	}
}

func (r *GridRowDefinition) FromMap(siMap a.SiMap) {
	height := siMap["height"]
	if engine.IsSpecialValueString(height) {
		r.Height = require.Float32(engine.GetSpecialValueFromString(height))
	} else {
		r.Height = require.Float32(height)
	}
}

func (r *GridRowDefinition) actualHeight() float32 {
	switch r.Height {
	case a.WrapContent:
		return r.maxHeight
	case a.FillParent:
		return r.fillHeight
	default:
		return r.Height
	}
}

type GridColumnDefinition struct {
	Width     float32
	maxWidth  float32
	fillWidth  float32
}

func (c *GridColumnDefinition) ToMap() a.SiMap {
	return a.SiMap{
		"width": c.Width,
	}
}

func (c *GridColumnDefinition) FromMap(siMap a.SiMap) {
	width := siMap["width"]
	if engine.IsSpecialValueString(width) {
		c.Width = require.Float32(engine.GetSpecialValueFromString(width))
	} else {
		c.Width = require.Float32(width)
	}
}

func (c *GridColumnDefinition) actualWidth() float32 {
	switch c.Width {
	case a.WrapContent:
		return c.maxWidth
	case a.FillParent:
		return c.fillWidth
	default:
		return c.Width
	}
}

const (
	GridVertical byte = iota
	GridHorizontal
)

//GridLayout is a simple layout component, that organizes its children into a grid.
//Each child is placed in cell after previous.
//If the orientation is vertical, the children are first places in rows from left to right.
//When a row is filled, the next row below is filled and so on.
//If the orientation is horizontal first columns are filled from top to bottom.
//When the column is filled, the next to the right is used.
//To use the grid you need to specify columns and rows definitions using AddColumn and AddRow methods respectively.
//For each row/column you can specify its height/width. You can set either a value in pixels or
//the special values a.FillParent and a.WrapContent.
//For a.WrapContent the row/col will take the size of the biggest child.
//For a.FillParent the row/col will try to fill all the available space not occupied by other rows/cols.
type GridLayout struct {
	engine.LayoutImpl
	Orientation     byte                   `state:"orientation"`
	AutoExpansion   bool                   `state:"autoExpansion"`
	AutoShrinking   bool                   `state:"autoShrinking"`
	AutoSizeAdjust  bool                   `state:"autoSizeAdjust"`
	RowPadding      float32                 `state:"rowPadding"`
	ColumnPadding   float32                 `state:"columnPadding"`
	Rows            []*GridRowDefinition    `state:"rows"`
	Columns         []*GridColumnDefinition `state:"columns"`
}

//AddRow add new row definition with the given height to the grid.
//The height may be eiter a value in pixels or a.WrapContent or a.FillParent.
func (l *GridLayout) AddRow(height float32) {
	l.Rows = append(l.Rows, &GridRowDefinition{Height: height})

	engine.RequestRendering()
}

//AddColumn add new column definition with the given width to the grid.
//The width may be eiter a value in pixels or a.WrapContent or a.FillParent.
func (l *GridLayout) AddColumn(width float32) {
	l.Columns = append(l.Columns, &GridColumnDefinition{Width: width})

	engine.RequestRendering()
}

//RemoveRowAt removes row at the given index.
//Note that the children in that row will not be removed, but moved to the next row if one exists.
func (l *GridLayout) RemoveRowAt(index int) {
	if index < 0 || index > len(l.Rows) {
		return
	}

	l.Rows = append(l.Rows[:index], l.Rows[index+1:]...)

	engine.RequestRendering()
}

//RemoveColumnAt removes column at the given index.
//Note that the children in that column will not be removed, but moved to the column row if one exists.
func (l *GridLayout) RemoveColumnAt(index int) {
	if index < 0 || index > len(l.Columns) {
		return
	}

	l.Columns = append(l.Columns[:index], l.Columns[index+1:]...)

	engine.RequestRendering()
}

//RemoveRow removes the last row of the grid. Objects in that row will no linger be visible.
func (l *GridLayout) RemoveRow() {
	l.RemoveRowAt(len(l.Rows) - 1)

	engine.RequestRendering()
}

//RemoveColumn removes the last column if the grid. Objects in that column will no linger be visible.
func (l *GridLayout) RemoveColumn() {
	l.RemoveColumnAt(len(l.Columns) - 1)

	engine.RequestRendering()
}

//GetRowsCount returns the number of rows in the grid.
func (l *GridLayout) GetRowsCount() int {
	return len(l.Rows)
}

//GetColumnsCount returns the number of columns in the grid.
func (l *GridLayout) GetColumnsCount() int {
	return len(l.Columns)
}

//SetRowHeight sets the height of the row at the given index.
func (l *GridLayout) SetRowHeight(rowIndex int, height float32) {
	if rowIndex < 0 || rowIndex >= len(l.Rows) {
		return
	}

	l.Rows[rowIndex].Height = height

	engine.RequestRendering()
}

//SetColumnWidth sets the width of the column at the given index.
func (l *GridLayout) SetColumnWidth(colIndex int, width float32) {
	if colIndex < 0 || colIndex >= len(l.Columns) {
		return
	}

	l.Columns[colIndex].Width = width

	engine.RequestRendering()
}

//GetRows returns the slice of all row definitions.
//Modifying the returned slice wont modify the actual rows of the grid.
func (l *GridLayout) GetRows() []*GridRowDefinition {
	rowsCopy := make([]*GridRowDefinition, len(l.Rows))
	copy(rowsCopy, l.Rows)
	return rowsCopy
}

//GetColumns returns the slice of all column definitions.
//Modifying the returned slice wont modify the actual columns of the grid.
func (l *GridLayout) GetColumns() []*GridColumnDefinition {
	colsCopy := make([]*GridColumnDefinition, len(l.Columns))
	copy(colsCopy, l.Columns)
	return colsCopy
}

//LayoutChildren implements the engine.Layout interface.
func (l *GridLayout) LayoutChildren() {
	l.LayoutImpl.LayoutChildren()

	children := l.SceneObject.GetChildren()

	l.expandIfNeeded(children)

	l.shrinkIfNeeded(children)

	cellsCount := l.calculateRowsColsMax(children)

	l.hideNotFittingObjects(children, cellsCount)

	l.setFillParentRowsCols()

	l.adjustSizeIfNeeded(l.layout(children))
}

func (l *GridLayout) MeasureContents() a.Vector3 {
	size := a.ZeroVector()

	for _, c := range l.Columns {
		if c.Width == a.WrapContent {
			size.X += c.maxWidth
		} else if c.Width != a.FillParent && c.Width != a.MatchParent {
			size.X += c.Width
		}
	}

	for _, c := range l.Rows {
		if c.Height == a.WrapContent {
			size.Y += c.maxHeight
		} else if c.Height != a.FillParent && c.Height != a.MatchParent {
			size.Y += c.Height
		}
	}

	size.X += l.ColumnPadding * float32(len(l.Columns) - 1)
	size.Y += l.RowPadding * float32(len(l.Rows) - 1)

	return size
}

func (l *GridLayout) expandIfNeeded(children []*engine.SceneObject) {
	if !l.AutoExpansion {
		return
	}

	if l.Orientation == GridVertical {
		if len(l.Columns) == 0 {
			l.AddColumn(a.WrapContent)
		}

		requiredRows := int(math.Ceil(float64(len(children)) / float64(len(l.Columns))))
		rowsToAdd := requiredRows - len(l.Rows)

		for r := 0; r < rowsToAdd; r++ {
			l.AddRow(a.WrapContent)
		}
	} else {
		if len(l.Rows) == 0 {
			l.AddRow(a.WrapContent)
		}

		requiredCols := int(math.Ceil(float64(len(children)) / float64(len(l.Rows))))
		colsToAdd := requiredCols - len(l.Columns)

		for c := 0; c < colsToAdd; c++ {
			l.AddColumn(a.WrapContent)
		}
	}
}

func (l *GridLayout) shrinkIfNeeded(children []*engine.SceneObject) {
	if l.AutoShrinking {
		if l.Orientation == GridVertical {
			requiredRows := int(math.Ceil(float64(len(children)) / float64(len(l.Columns))))

			for r := len(l.Rows) - 1; r > requiredRows; r-- {
				l.RemoveRow()
			}
		} else {
			requiredCols := int(math.Ceil(float64(len(children)) / float64(len(l.Rows))))

			for c := len(l.Columns) - 1; c > requiredCols; c-- {
				l.RemoveColumn()
			}
		}
	}
}

func (l *GridLayout) calculateRowsColsMax(children []*engine.SceneObject) int {
	var forr func(action func(r, c int, row *GridRowDefinition, col *GridColumnDefinition) bool)
	if l.Orientation == GridVertical {
		forr = l.forVertical
	} else {
		forr = l.forHorizontal
	}

	cellsCount := 0

	forr(func(r, c int, row *GridRowDefinition, col *GridColumnDefinition) bool {
		if cellsCount >= len(children) {
			return false
		}

		child := children[cellsCount]
		chSize := child.Transform.ActualSize()
		chSize1 := child.Transform.WantedSize()

		if chSize1.Y != a.MatchParent && chSize.Y > row.maxHeight {
			row.maxHeight = chSize.Y
		}
		if chSize1.X != a.MatchParent && chSize.X > col.maxWidth {
			col.maxWidth = chSize.X
		}

		if row.Height == a.FillParent {
			row.fillHeight = 0
		}
		if col.Width == a.FillParent {
			col.fillWidth = 0
		}

		cellsCount++
		return true
	})

	return cellsCount
}

func (l *GridLayout) hideNotFittingObjects(children []*engine.SceneObject, cellsCount int) {
	for i := cellsCount; i < len(children); i++ {
		child := children[i]
		child.Transform.SetSize(a.ZeroVector())
	}
}

func (l *GridLayout) calculateRowsMetrics() (float32, int) {
	var totalRowsHeight float32
	var rowFillCount int

	for r := 0; r < len(l.Rows); r++ {
		row := l.Rows[r]
		totalRowsHeight += row.actualHeight()
		if row.Height == a.FillParent {
			rowFillCount++
		}
	}

	return totalRowsHeight, rowFillCount
}

func (l *GridLayout) calculateColumnsMetrics() (float32, int) {
	var totalColsWidth float32
	var colFillCount int

	for c := 0; c < len(l.Columns); c++ {
		col := l.Columns[c]
		totalColsWidth += col.actualWidth()
		if col.Width == a.FillParent {
			colFillCount++
		}
	}

	return totalColsWidth, colFillCount
}

func (l *GridLayout) setFillParentRowsCols() {
	gridSize := l.SceneObject.Transform.ActualSize()

	totalRowsHeight, rowFillCount := l.calculateRowsMetrics()
	totalColsWidth, colFillCount := l.calculateColumnsMetrics()

	availableColsSpace := common.MaxFloat32(gridSize.X-totalColsWidth-l.ColumnPadding*2*float32(l.GetColumnsCount()), 0)
	availableRowsSpace := common.MaxFloat32(gridSize.Y-totalRowsHeight-l.RowPadding*2*float32(l.GetRowsCount()), 0)

	for r := 0; r < len(l.Rows); r++ {
		for c := 0; c < len(l.Columns); c++ {
			row := l.Rows[r]
			col := l.Columns[c]

			if row.Height == a.FillParent {
				row.fillHeight = availableRowsSpace / float32(rowFillCount)
			}

			if col.Width == a.FillParent {
				col.fillWidth = availableColsSpace / float32(colFillCount)
			}
		}
	}
}

func (l *GridLayout) layout(children []*engine.SceneObject) (float32, float32) {
	if l.Orientation == GridVertical {
		return l.layoutVertical(children)
	} else {
		return l.layoutHorizontal(children)
	}
}

func (l *GridLayout) layoutVertical(children []*engine.SceneObject) (x float32, y float32) {
	var i = 0
	y = 0

	for r := 0; r < len(l.Rows); r++ {
		row := l.Rows[r]
		x = 0

		for c := 0; c < len(l.Columns); c++ {
			if i >= len(children) {
				break
			}

			child := children[i]
			col := l.Columns[c]

			l.adjustChildTransformIfNeeded(child, row, col, x, y)

			x += col.actualWidth() + l.ColumnPadding

			i++
		}

		y += row.actualHeight() + l.RowPadding
	}

	return
}

func (l *GridLayout) layoutHorizontal(children []*engine.SceneObject) (x float32, y float32) {
	var i = 0
	x = 0

	for c := 0; c < len(l.Columns); c++ {
		col := l.Columns[c]
		y = 0

		for r := 0; r < len(l.Rows); r++ {
			if i >= len(children) {
				break
			}

			child := children[i]
			row := l.Rows[r]

			l.adjustChildTransformIfNeeded(child, row, col, x, y)

			y += row.actualHeight() + l.RowPadding

			i++
		}

		x += col.actualWidth() + l.ColumnPadding
	}

	return
}

func (l *GridLayout) adjustChildTransformIfNeeded(child *engine.SceneObject,
	row *GridRowDefinition, col *GridColumnDefinition,
	x, y float32) {

	pos := a.NewVector3(x, y, l.SceneObject.Transform.WantedPosition().Z)
	size := a.NewVector3(
		col.actualWidth(),
		row.actualHeight(),
		child.Transform.WantedSize().Z,
	)

	if child.Transform.WantedPosition().Equals(pos) && child.Transform.WantedSize().Equals(size) {
		return
	}

	child.Transform.SetSize(size)
	child.Transform.SetPosition(pos)

	child.Redraw()
}

func (l *GridLayout) adjustSizeIfNeeded(x, y float32) {
	if !l.AutoSizeAdjust {
		return
	}

	l.SceneObject.Transform.SetSize(x, y)
}

func (l *GridLayout) forVertical(action func(r, c int, row *GridRowDefinition, col *GridColumnDefinition) bool) {
	l.forRows(func(r int, row *GridRowDefinition) bool {
		br := false
		l.forColumns(func(c int, col *GridColumnDefinition) bool {
			if br = action(r, c, row, col); !br {
				return false
			}
			return true
		})
		return br
	})
}

func (l *GridLayout) forHorizontal(action func(r, c int, row *GridRowDefinition, col *GridColumnDefinition) bool) {
	l.forColumns(func(c int, col *GridColumnDefinition) bool {
		br := false
		l.forRows(func(r int, row *GridRowDefinition) bool {
			if br = action(r, c, row, col); !br {
				return false
			}
			return true
		})
		return br
	})
}

func (l *GridLayout) forRows(action func(r int, row *GridRowDefinition) bool) {
	for r := 0; r < len(l.Rows); r++ {
		if !action(r, l.Rows[r]) {
			break
		}
	}
}

func (l *GridLayout) forColumns(action func(c int, col *GridColumnDefinition) bool) {
	for c := 0; c < len(l.Columns); c++ {
		if !action(c, l.Columns[c]) {
			break
		}
	}
}

//NewGridLayout creates a new empty GridLayout.
func NewGridLayout() *GridLayout {
	return &GridLayout{
		Rows:          make([]*GridRowDefinition, 0, 10),
		Columns:       make([]*GridColumnDefinition, 0, 10),
		AutoExpansion: true,
		AutoShrinking: true,
	}
}

//NewGridLayoutSized creates a new GridLayout with the given number of rows and columns.
func NewGridLayoutSized(rows, cols int) *GridLayout {
	grid := NewGridLayout()
	for i := 0; i < rows; i++ {
		grid.AddRow(a.WrapContent)
	}
	for i := 0; i < cols; i++ {
		grid.AddColumn(a.WrapContent)
	}
	return grid
}
