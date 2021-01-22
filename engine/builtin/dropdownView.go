package builtin

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type OnSelectHandler func(item string)

type DropdownView struct {
	engine.ViewImpl
	//arrow1Id         int64
	//arrow2Id         int64
	arrowId          int
	items            []string
	selectedItem     string
	textView         *TextView
	optionsContainer *engine.SceneObject
	showItems        bool
	OnSelect         OnSelectHandler
}

func (d *DropdownView) OnStart() {
	d.arrowId = d.Engine.GetRenderer().AddPrimitive()
	//d.arrow2Id = d.eng.GetRenderer().AddPrimitive()

	siz := d.SceneObject.Transform.Size

	bg := NewShapeView(rendering.PrimitiveRectangle)
	bg.StrokeWeight = 2
	bg.StrokeColor = a.BlackColor()
	bg.FillColor = a.WhiteColor()
	bg.CornerRadius = 10
	d.SceneObject.AddComponent(bg)

	textObj := engine.NewSceneObject("selected item")
	textObj.Transform.Position = a.NewVector3(5, 5, 1)
	textObj.Transform.Size = a.NewVector3(siz.X - 10, siz.Y - 10, 0)

	d.textView = NewTextView(d.selectedItem)
	d.textView.FontSize = 16

	textObj.AddComponent(d.textView)
	d.SceneObject.AddChild(textObj)

	d.optionsContainer = engine.NewSceneObject("OptionsContainer")
	optionsBg := NewShapeView(rendering.PrimitiveRectangle)
	optionsBg.FillColor = a.WhiteColor()
	optionsBg.CornerRadius = 10
	d.optionsContainer.AddComponent(optionsBg)

	d.optionsContainer.Transform.Position = a.NewVector3(0, siz.Y, 1)
	d.optionsContainer.Transform.Size = a.NewVector3(siz.X, float32(35*len(d.items)) + 5, 0)

	for i, o := range d.items {
		var itemText = o
		item := engine.NewSceneObject(fmt.Sprintf("Item%d", i))
		item.Transform.Position = a.NewVector3(10, float32(i*35) + 5, 1)
		item.Transform.Size = a.NewVector3(siz.X, 35, 0)
		itemTextView := NewTextView(itemText)
		itemTextView.TextColor = a.BlackColor()
		itemTextView.FontSize = 14
		item.AddComponent(itemTextView)
		item.AddComponent(NewRectBoundary())
		item.AddComponent(NewOnClickListener(func(event engine.AmphionEvent) bool {
			newSelectedItem := itemText
			d.selectedItem = newSelectedItem
			d.textView.SetText(newSelectedItem)
			d.showItems = false
			d.optionsContainer.SetEnabled(false)
			if d.OnSelect != nil {
				d.OnSelect(newSelectedItem)
			}
			d.Engine.RequestRendering()
			return false
		}))
		d.optionsContainer.AddChild(item)
	}

	d.SceneObject.AddChild(d.optionsContainer)
	d.optionsContainer.SetEnabled(false)
}

func (d *DropdownView) OnDraw(ctx engine.DrawingContext) {
	//pos := d.obj.Transform.GetGlobalTopLeftPosition()
	//rect := d.obj.Transform.GetGlobalRect()
	//x1 := int(math.Round(float64(rect.X.Max))) - 25
	//x2 := int(math.Round(float64(rect.X.Max))) - 5
	////x3 := x1 + int(math.Round(common.NewFloatRange(float64(x1), float64(x2)).GetLength() / 2))
	//y1 := int(math.Round(float64(rect.Y.Min))) + 12
	//y2 := int(math.Round(float64(rect.Y.Max))) - 12
	//z1 := int(math.Round(float64(pos.Z + 1)))

	//pr := rendering.NewImagePrimitive(4)
	//pr.Transform.Position = a.NewIntVector3(x1, y1, z1)
	//pr.Transform.Size = a.NewIntVector3(x2 - x1, y2 - y1, 0)
	//
	//ctx.GetRenderer().SetPrimitive(d.arrowId, pr, d.ShouldRedraw())
}

func (d *DropdownView) HandleClick() {
	d.showItems = !d.showItems
	if d.showItems {
		d.showDropdown()
	} else {
		d.hideDropdown()
	}
}

func (d *DropdownView) SetItems(items []string) {
	d.items = items
	if len(d.items) > 0 {
		d.selectedItem = d.items[0]
	}
	d.ForceRedraw()
}

func (d *DropdownView) SetSelectedItem(item string) {
	d.selectedItem = item
	d.ForceRedraw()
}

func (d *DropdownView) GetSelectedItem() string {
	return d.selectedItem
}

func (d *DropdownView) showDropdown() {
	siz := d.SceneObject.Transform.Size
	d.optionsContainer.Transform.Position = a.NewVector3(0, siz.Y, 0)
	d.optionsContainer.Transform.Size = a.NewVector3(siz.X, float32(35*len(d.items)), 0)
	d.optionsContainer.SetEnabled(true)
	d.Engine.RequestRendering()
}

func (d *DropdownView) hideDropdown() {
	d.optionsContainer.SetEnabled(false)
	d.Engine.RequestRendering()
}

func (d *DropdownView) OnStop() {
	d.Engine.GetRenderer().RemovePrimitive(d.arrowId)
}

func (d *DropdownView) GetName() string {
	return engine.NameOfComponent(d)
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