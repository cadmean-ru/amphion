package cli

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

type RendererDelegate interface {
	rendering.RendererDelegate
}

type PrimitiveRenderingContext struct {
	GeometryPrimitiveData *GeometryPrimitiveData
	ImagePrimitiveData    *ImagePrimitiveData
	PrimitiveKind         int
	State                 interface{}
	Redraw                bool
	PrimitiveId           int
}

func newCliPrimitiveRenderingContext(ctx *rendering.PrimitiveRenderingContext) *PrimitiveRenderingContext {
	cliCtx := &PrimitiveRenderingContext{
		PrimitiveKind: int(ctx.PrimitiveKind),
		Redraw: ctx.Redraw,
	}

	screenSize := engine.GetScreenSize3()
	t := ctx.Primitive.GetTransform()
	tlPosN := t.Position.Ndc(screenSize)
	brPosN := t.Position.Add(t.Size).Ndc(screenSize)

	switch ctx.PrimitiveKind {
	case rendering.PrimitiveTriangle, rendering.PrimitiveRectangle, rendering.PrimitiveEllipse, rendering.PrimitiveLine:
		gp := ctx.Primitive.(*rendering.GeometryPrimitive)

		fillColorN := gp.Appearance.FillColor.Normalize()
		strokeColorN := gp.Appearance.StrokeColor.Normalize()

		cliCtx.GeometryPrimitiveData = &GeometryPrimitiveData{
			GeometryType: int(ctx.PrimitiveKind),
			TlPositionN:  &Vector3{
				X: tlPosN.X,
				Y: tlPosN.Y,
				Z: float32(t.Position.Z),
			},
			BrPositionN:  &Vector3{
				X: brPosN.X,
				Y: brPosN.Y,
				Z: brPosN.Z,
			},
			FillColorN:   &Vector4{
				X: fillColorN.X,
				Y: fillColorN.Y,
				Z: fillColorN.Z,
				W: fillColorN.W,
			},
			StrokeColorN: &Vector4{
				X: strokeColorN.X,
				Y: strokeColorN.Y,
				Z: strokeColorN.Z,
				W: strokeColorN.W,
			},
			StrokeWeight: int(gp.Appearance.StrokeWeight),
			CornerRadius: int(gp.Appearance.CornerRadius),
		}
	case rendering.PrimitiveImage:
		ip := ctx.Primitive.(*rendering.ImagePrimitive)

		cliCtx.ImagePrimitiveData = &ImagePrimitiveData{
			TlPositionN:  &Vector3{
				X: tlPosN.X,
				Y: tlPosN.Y,
				Z: float32(t.Position.Z),
			},
			BrPositionN:  &Vector3{
				X: brPosN.X,
				Y: brPosN.Y,
				Z: brPosN.Z,
			},
			ImageUrl:    ip.ImageUrl,
		}
	}

	return cliCtx
}

type PrimitiveRendererDelegate interface {
	OnStart()
	OnSetPrimitive(ctx *PrimitiveRenderingContext)
	OnRender(ctx *PrimitiveRenderingContext)
	OnRemovePrimitive(ctx *PrimitiveRenderingContext)
	OnStop()
}
