package builtin

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
)

type OnSelectHandler func(item string)

type DropdownView struct {
	engine.ViewImpl
	arrowId          int
	items            []string
	selectedItem     string
	textView         *TextView
	optionsContainer *engine.SceneObject
	showItems        bool
	OnSelect         OnSelectHandler
}

func (d *DropdownView) OnStart() {
	d.arrowId = d.Context.GetRenderingNode().AddPrimitive()
	//d.arrow2Id = d.eng.GetRenderer().AddPrimitive()

	siz := d.SceneObject.Transform.WantedSize()

	bg := NewShapeView(ShapeRectangle)
	bg.StrokeWeight = 2
	bg.StrokeColor = a.Black()
	bg.FillColor = a.White()
	bg.CornerRadius = 10
	d.SceneObject.AddComponent(bg)

	textObj := engine.NewSceneObject("selected item")
	textObj.Transform.SetPosition(5, 5, 1)
	textObj.Transform.SetSize(siz.X - 10, siz.Y - 10)

	d.textView = NewTextView(d.selectedItem)
	d.textView.FontSize = 16

	textObj.AddComponent(d.textView)
	d.SceneObject.AddChild(textObj)

	d.optionsContainer = engine.NewSceneObject("OptionsContainer")
	optionsBg := NewShapeView(ShapeRectangle)
	optionsBg.FillColor = a.White()
	optionsBg.CornerRadius = 10
	d.optionsContainer.AddComponent(optionsBg)

	d.optionsContainer.Transform.SetPosition(0, siz.Y, 1)
	d.optionsContainer.Transform.SetSize(siz.X, float32(35*len(d.items)) + 5)

	for i, o := range d.items {
		var itemText = o
		item := engine.NewSceneObject(fmt.Sprintf("Item%d", i))
		item.Transform.SetPosition(10, float32(i*35) + 5, 1)
		item.Transform.SetSize(siz.X, 35)
		itemTextView := NewTextView(itemText)
		itemTextView.TextColor = a.Black()
		itemTextView.FontSize = 14
		item.AddComponent(itemTextView)
		item.AddComponent(NewRectBoundary())
		item.AddComponent(NewOnClickListener(func(event engine.Event) bool {
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

func (d *DropdownView) OnDraw(_ engine.DrawingContext) {
	//pos := d.obj.Transform.GlobalTopLeftPosition()
	//rect := d.obj.Transform.GlobalRect()
	//x1 := int(math.Round(float64(rect.X.Max))) - 25
	//x2 := int(math.Round(float64(rect.X.Max))) - 5
	////x3 := x1 + int(math.Round(common.NewFloatRange(float64(x1), float64(x2)).GetLength() / 2))
	//y1 := int(math.Round(float64(rect.Y.Min))) + 12
	//y2 := int(math.Round(float64(rect.Y.Max))) - 12
	//z1 := int(math.Round(float64(pos.Z + 1)))

	//pr := rendering.NewImagePrimitive(4)
	//pr.Transform.position = a.NewIntVector3(x1, y1, z1)
	//pr.Transform.size = a.NewIntVector3(x2 - x1, y2 - y1, 0)
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
	d.Redraw()
}

func (d *DropdownView) SetSelectedItem(item string) {
	d.selectedItem = item
	d.Redraw()
}

func (d *DropdownView) GetSelectedItem() string {
	return d.selectedItem
}

func (d *DropdownView) showDropdown() {
	siz := d.SceneObject.Transform.WantedSize()
	d.optionsContainer.Transform.SetPosition(0, siz.Y)
	d.optionsContainer.Transform.SetSize(siz.X, float32(35*len(d.items)))
	d.optionsContainer.SetEnabled(true)
	d.Engine.RequestRendering()
}

func (d *DropdownView) hideDropdown() {
	d.optionsContainer.SetEnabled(false)
	d.Engine.RequestRendering()
}

func (d *DropdownView) OnStop() {
	d.Context.GetRenderingNode().RemovePrimitive(d.arrowId)
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