package ios

import (
	"fmt"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend/cli"
	"github.com/cadmean-ru/amphion/rendering"
)

type MetalRenderer struct {
}

func (m *MetalRenderer) SetPrimitive(id int, primitive rendering.IPrimitive, shouldRerender bool) {
	fmt.Println("Setting primitive in go 0")

	if !shouldRerender {
		return
	}

	fmt.Println("Setting primitive in go 1")

	screenSize := engine.GetScreenSize3()
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

		fmt.Printf("%+v\n", fillColorN)
		fmt.Printf("%+v\n", gp.Appearance.FillColor)
		fmt.Printf("%+v\n", strokeColorN)

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
		fmt.Println(data)
	case rendering.PrimitiveImage:
		ip := primitive.(*rendering.ImagePrimitive)
		data := cli.ImagePrimitiveData{
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
			ImageUrl:    ip.ImageUrl,
		}
		fmt.Println(data)

	}
}

