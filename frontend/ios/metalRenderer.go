package ios

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/rendering"
)

type MetalRenderer struct {
	r cli.RendererCLI
}

func (m *MetalRenderer) Prepare() {
	fmt.Println("Preparing in go")
	m.r.Prepare()
}

func (m *MetalRenderer) AddPrimitive() int {
	fmt.Println("Adding primitive in go")
	return m.r.AddPrimitive()
}

func (m *MetalRenderer) SetPrimitive(id int, primitive rendering.IPrimitive, shouldRerender bool) {
	fmt.Println("Setting primitive in go 0")

	if !shouldRerender {
		return
	}

	fmt.Println("Setting primitive in go 1")

	screenSize := a.NewIntVector3(engine.GetInstance().GetGlobalContext().ScreenInfo.GetWidth(), engine.GetInstance().GetGlobalContext().ScreenInfo.GetHeight(), 0)
	t := primitive.GetTransform()
	tlPosN := t.Position.Ndc(screenSize)
	brPosN := t.Position.Add(t.Size).Ndc(screenSize)

	fmt.Printf("%+v\n", screenSize)
	fmt.Printf("%+v\n", engine.GetInstance().GetCurrentScene().Transform.Size)
	fmt.Printf("%+v\n", tlPosN)
	fmt.Printf("%+v\n", brPosN)


	switch primitive.GetType() {
	case rendering.PrimitiveRectangle, rendering.PrimitiveEllipse, rendering.PrimitiveTriangle:
		gp := primitive.(*rendering.GeometryPrimitive)

		fillColorN := gp.Appearance.FillColor.Normalize()
		strokeColorN := gp.Appearance.StrokeColor.Normalize()

		data := cli.GeometryPrimitiveData{
			GeometryType: int(primitive.GetType()),
			TlPositionN:  &cli.Vector3{
				X: tlPosN.X,
				Y: tlPosN.Y,
				Z: float32(t.Position.Z),
			},
			BrPositionN:  &cli.Vector3{
				X: brPosN.X,
				Y: brPosN.Y,
				Z: brPosN.Z,
			},
			FillColorN:   &cli.Vector4{
				X: fillColorN.X,
				Y: fillColorN.Y,
				Z: fillColorN.Z,
				W: fillColorN.W,
			},
			StrokeColorN: &cli.Vector4{
				X: strokeColorN.X,
				Y: strokeColorN.Y,
				Z: strokeColorN.Z,
				W: strokeColorN.W,
			},
			StrokeWeight: int(gp.Appearance.StrokeWeight),
			CornerRadius: int(gp.Appearance.CornerRadius),
		}

		fmt.Println("Setting primitive in go 2")
		m.r.SetGeometryPrimitive(id, &data)
	}
}

func (m *MetalRenderer) RemovePrimitive(id int) {
	m.r.RemovePrimitive(id)
}

func (m *MetalRenderer) PerformRendering() {
	fmt.Println("Performing rendering in go")
	m.r.PerformRendering()
}

func (m *MetalRenderer) Clear() {
	m.r.Clear()
}

func (m *MetalRenderer) Stop() {
	m.r.Stop()
}

