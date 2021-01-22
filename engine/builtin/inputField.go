package builtin

import (
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type InputField struct {
	engine.ViewImpl
	rendering.Appearance
	Padding        float32
	AllowMultiline bool
	inputView      *InputView
}

func (f *InputField) OnStart() {
	f.ViewImpl.OnStart()

	input := engine.NewSceneObject("input")

	f.inputView = NewInputView()
	f.inputView.TextColor = a.BlackColor()
	f.inputView.FontSize = 15
	f.inputView.AllowMultiline = f.AllowMultiline

	input.Transform.Position = a.NewVector3(f.Padding, f.Padding, 1)
	size := f.SceneObject.Transform.Size
	input.Transform.Size = a.NewVector3(size.X - f.Padding, size.Y - f.Padding, 0)

	input.AddComponent(f.inputView)
	input.AddComponent(NewRectBoundary())

	f.SceneObject.AddChild(input)
}

func (f *InputField) OnDraw(ctx engine.DrawingContext) {
	pr := rendering.NewGeometryPrimitive(rendering.PrimitiveRectangle)
	pr.Transform = transformToRenderingTransform(f.SceneObject.Transform)
	pr.Appearance = f.Appearance
	ctx.GetRenderer().SetPrimitive(f.PrimitiveId, pr, f.ShouldRedraw())
}

func (f *InputField) GetName() string {
	return engine.NameOfComponent(f)
}

func NewInputField() *InputField {
	return &InputField{
		Appearance: rendering.Appearance{
			FillColor:    a.WhiteColor(),
			StrokeColor:  a.BlackColor(),
			StrokeWeight: 2,
			CornerRadius: 10,
		},
	}
}