package builtin

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type InputField struct {
	ViewImpl
	rendering.Appearance
	Padding        float64
	AllowMultiline bool
	inputView      *InputView
}

func (f *InputField) OnStart() {
	f.ViewImpl.OnStart()

	input := engine.NewSceneObject("input")

	f.inputView = NewInputView()
	f.inputView.Appearance.FillColor = common.BlackColor()
	f.inputView.TextAppearance.FontSize = 15
	f.inputView.AllowMultiline = f.AllowMultiline

	input.Transform.Position = common.NewVector3(f.Padding, f.Padding, 1)
	size := f.obj.Transform.Size
	input.Transform.Size = common.NewVector3(size.X - f.Padding, size.Y - f.Padding, 0)

	input.AddComponent(f.inputView)
	input.AddComponent(NewRectBoundary())

	f.obj.AddChild(input)
}

func (f *InputField) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(rendering.PrimitiveRectangle)
	pr.Transform = transformToRenderingTransform(f.obj.Transform)
	pr.Appearance = f.Appearance
	ctx.GetRenderer().SetPrimitive(f.pId, pr, f.ShouldRedraw())
}

func (f *InputField) GetName() string {
	return "InputField"
}

func NewInputField() *InputField {
	return &InputField{
		Appearance: rendering.Appearance{
			FillColor:    common.WhiteColor(),
			StrokeColor:  common.BlackColor(),
			StrokeWeight: 2,
			CornerRadius: 10,
		},
	}
}