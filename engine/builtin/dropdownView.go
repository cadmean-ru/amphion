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
	//arrow1Id         int64
	//arrow2Id         int64
	arrowId          int64
	items            []string
	selectedItem     string
	textView         *TextView
	optionsContainer *engine.SceneObject
	showItems        bool
}

func (d *DropdownView) OnStart() {
	d.arrowId = d.eng.GetRenderer().AddPrimitive()
	//d.arrow2Id = d.eng.GetRenderer().AddPrimitive()

	siz := d.obj.Transform.Size

	bg := NewShapeView(rendering.PrimitiveRectangle)
	bg.StrokeWeight = 2
	bg.StrokeColor = common.BlackColor()
	bg.FillColor = common.WhiteColor()
	bg.CornerRadius = 10
	d.obj.AddComponent(bg)

	textObj := engine.NewSceneObject("selected item")
	textObj.Transform.Position = common.NewVector3(5, 5, 1)
	textObj.Transform.Size = common.NewVector3(siz.X - 10, siz.Y - 10, 0)

	d.textView = NewTextView(d.selectedItem)
	d.textView.Appearance = rendering.Appearance{
		FillColor:    common.BlackColor(),
	}
	d.textView.TextAppearance = rendering.TextAppearance{
		FontSize: 16,
	}

	textObj.AddComponent(d.textView)
	d.obj.AddChild(textObj)

	d.optionsContainer = engine.NewSceneObject("OptionsContainer")
	optionsBg := NewShapeView(rendering.PrimitiveRectangle)
	optionsBg.Appearance.FillColor = common.WhiteColor()
	optionsBg.Appearance.CornerRadius = 10
	d.optionsContainer.AddComponent(optionsBg)

	d.optionsContainer.Transform.Position = common.NewVector3(0, siz.Y, 1)
	d.optionsContainer.Transform.Size = common.NewVector3(siz.X, float64(35*len(d.items)) + 5, 0)

	for i, o := range d.items {
		var itemText = o
		item := engine.NewSceneObject(fmt.Sprintf("Item%d", i))
		item.Transform.Position = common.NewVector3(10, float64(i*35) + 5, 1)
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
		d.optionsContainer.AddChild(item)
	}

	d.obj.AddChild(d.optionsContainer)
	d.optionsContainer.SetEnabled(false)
}

func (d *DropdownView) OnDraw(ctx engine.DrawingContext) {
	pos := d.obj.Transform.GetGlobalTopLeftPosition()
	rect := d.obj.Transform.GetGlobalRect()
	x1 := int(math.Round(rect.X.Max)) - 40
	x2 := int(math.Round(rect.X.Max)) - 20
	//x3 := x1 + int(math.Round(common.NewFloatRange(float64(x1), float64(x2)).GetLength() / 2))
	y1 := int(math.Round(rect.Y.Min)) + 12
	y2 := int(math.Round(rect.Y.Max)) - 12
	z1 := int(math.Round(pos.Z + 1))
	//
	//lp1 := rendering.NewGeometryPrimitive(rendering.PrimitiveLine)
	//lp1.Transform.Position = common.NewIntVector3(x1, y1, z1)
	//lp1.Transform.Size = common.NewIntVector3(x3 - x1, y2 - y1, 0)
	//lp1.Appearance.StrokeWeight = 3
	//lp1.Appearance.StrokeColor = common.NewColor(0xc4, 0xc4, 0xc4, 0xff)
	//
	//lp2 := rendering.NewGeometryPrimitive(rendering.PrimitiveLine)
	//lp2.Transform.Position = common.NewIntVector3(x3, y2, z1)
	//lp2.Transform.Size = common.NewIntVector3(x2 - x3, y1 - y2, 0)
	//lp2.Appearance.StrokeWeight = 3
	//lp2.Appearance.StrokeColor = common.NewColor(0xc4, 0xc4, 0xc4, 0xff)
	//
	//ctx.GetRenderer().SetPrimitive(d.arrow1Id, lp1, d.ShouldRedraw())
	//ctx.GetRenderer().SetPrimitive(d.arrow2Id, lp2, d.ShouldRedraw())

	pr := rendering.NewImagePrimitive(4)
	pr.Transform.Position = common.NewIntVector3(x1, y1, z1)
	pr.Transform.Size = common.NewIntVector3(x2 - x1, y2 - y1, 0)

	ctx.GetRenderer().SetPrimitive(d.arrowId, pr, d.ShouldRedraw())
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
	d.eng.GetRenderer().RemovePrimitive(d.arrowId)
	//d.eng.GetRenderer().RemovePrimitive(d.arrow2Id)
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