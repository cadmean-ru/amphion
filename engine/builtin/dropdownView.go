package builtin

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"math"
)

type DropdownView struct {
	ViewImpl
	arrow1Id         int64
	arrow2Id         int64
	items            []string
	selectedItem     string
	textView         *TextView
	optionsContainer *engine.SceneObject
	showItems        bool
}

func (d *DropdownView) OnStart() {
	d.arrow1Id = d.eng.GetRenderer().AddPrimitive()
	d.arrow2Id = d.eng.GetRenderer().AddPrimitive()

	comp := d.obj.GetComponentByName("TextView")
	if comp == nil {
		panic("No text view for dropdown")
	}
	d.textView = comp.(*TextView)
	d.textView.Appearance = rendering.Appearance{
		FillColor:    common.BlackColor(),
	}
	d.textView.TextAppearance = rendering.TextAppearance{
		FontSize: 16,
	}
	d.textView.SetText(d.selectedItem)

	d.optionsContainer = engine.NewSceneObject("OptionsContainer")
	d.obj.AddChild(d.optionsContainer)
	d.optionsContainer.AddComponent(NewBoundaryView())

	siz := d.obj.Transform.Size
	d.optionsContainer.Transform.Position = common.NewVector3(0, siz.Y, 0)
	d.optionsContainer.Transform.Size = common.NewVector3(siz.X, float64(35*len(d.items)), 0)

	for i, o := range d.items {
		var itemText = o
		item := engine.NewSceneObject(fmt.Sprintf("Item%d", i))
		item.Transform.Position = common.NewVector3(0, float64(i*35), 0)
		item.Transform.Size = common.NewVector3(siz.X, 35, 0)
		itemTextView := NewTextView(itemText)
		itemTextView.Appearance.FillColor = common.BlackColor()
		itemTextView.TextAppearance.FontSize = 14
		item.AddComponent(itemTextView)
		item.AddComponent(NewRectBoundary())
		item.AddComponent(NewOnClickListener(func(event engine.AmphionEvent) bool {
			newSelectedItem := itemText
			d.selectedItem = newSelectedItem
			d.textView.SetText(newSelectedItem)
			d.showItems = false
			d.optionsContainer.SetEnabled(false)
			d.eng.RequestRendering()
			return false
		}))
		//item.AddComponent(NewBoundaryView())
		//item.AddComponent(NewDropDownItemView(o, d))
		d.optionsContainer.AddChild(item)
	}

	d.optionsContainer.SetEnabled(false)
}

func (d *DropdownView) OnDraw(ctx engine.DrawingContext) {
	pos := d.obj.Transform.GetGlobalTopLeftPosition()
	rect := d.obj.Transform.GetGlobalRect()
	x1 := int(math.Round(rect.X.Min + rect.X.GetLength() * 0.9))
	x2 := int(math.Round(rect.X.Max))
	x3 := x1 + int(math.Round(common.NewFloatRange(float64(x1), float64(x2)).GetLength() / 2))
	y1 := int(math.Round(rect.Y.Min))
	y2 := int(math.Round(rect.Y.Max))
	z1 := int(math.Round(pos.Z + 1))

	lp1 := rendering.NewGeometryPrimitive(rendering.PrimitiveLine)
	lp1.Transform.Position = common.NewIntVector3(x1, y1, z1)
	lp1.Transform.Size = common.NewIntVector3(x3 - x1, y2 - y1, 0)

	lp2 := rendering.NewGeometryPrimitive(rendering.PrimitiveLine)
	lp2.Transform.Position = common.NewIntVector3(x3, y2, z1)
	lp2.Transform.Size = common.NewIntVector3(x2 - x3, y1 - y2, 0)

	ctx.GetRenderer().SetPrimitive(d.arrow1Id, lp1, d.ShouldRedraw())
	ctx.GetRenderer().SetPrimitive(d.arrow2Id, lp2, d.ShouldRedraw())
}

func (d *DropdownView) HandleClick() {
	d.showItems = !d.showItems
	if d.showItems {
		d.showDropdown()
	} else {
		d.hideDropdown()
	}
}

func (d *DropdownView) showDropdown() {
	siz := d.obj.Transform.Size
	d.optionsContainer.Transform.Position = common.NewVector3(0, siz.Y, 0)
	d.optionsContainer.Transform.Size = common.NewVector3(siz.X, float64(35*len(d.items)), 0)
	d.optionsContainer.SetEnabled(true)
	d.eng.RequestRendering()
}

func (d *DropdownView) hideDropdown() {
	d.optionsContainer.SetEnabled(false)
	d.eng.RequestRendering()
}

func (d *DropdownView) OnStop() {
	d.eng.GetRenderer().RemovePrimitive(d.arrow1Id)
	d.eng.GetRenderer().RemovePrimitive(d.arrow2Id)
}

func (d *DropdownView) GetName() string {
	return "DropdownView"
}

func NewDropdownView(options []string) *DropdownView {
	defText := ""
	if len(options) == 0 {
		defText = "Select"
	} else {
		defText = options[0]
	}

	return &DropdownView{
		items:        options,
		selectedItem: defText,
		showItems:    false,
	}
}