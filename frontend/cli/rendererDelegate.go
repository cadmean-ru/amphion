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
	TextPrimitiveData     *TextPrimitiveData
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
			TlPositionN:  NewVector3(
				tlPosN.X,
				tlPosN.Y,
				float32(t.Position.Z),
			),
			BrPositionN:  NewVector3(
				brPosN.X,
				brPosN.Y,
				brPosN.Z,
			),
			FillColorN:   NewVector4(
				fillColorN.X,
				fillColorN.Y,
				fillColorN.Z,
				fillColorN.W,
			),
			StrokeColorN: NewVector4(
				strokeColorN.X,
				strokeColorN.Y,
				strokeColorN.Z,
				strokeColorN.W,
			),
			StrokeWeight: int(gp.Appearance.StrokeWeight),
			CornerRadius: int(gp.Appearance.CornerRadius),
		}
	case rendering.PrimitiveImage:
		ip := ctx.Primitive.(*rendering.ImagePrimitive)

		cliCtx.ImagePrimitiveData = &ImagePrimitiveData{
			TlPositionN:  NewVector3(
				tlPosN.X,
				tlPosN.Y,
				float32(t.Position.Z),
			),
			BrPositionN:  NewVector3(
				brPosN.X,
				brPosN.Y,
				brPosN.Z,
			),
			ImageUrl:    ip.ImageUrl,
		}
	case rendering.PrimitiveText:
		tp := ctx.Primitive.(*rendering.TextPrimitive)

		cliCtx.TextPrimitiveData = &TextPrimitiveData{
			Text:        tp.Text,
			TlPosition:  NewVector3(
				float32(t.Position.X),
				float32(t.Position.Y),
				float32(t.Position.Z),
			),
			Size:  NewVector3(
				float32(t.Size.X),
				float32(t.Size.Y),
				float32(t.Size.Z),
			),
			TextColorN:   NewVector4FromAVector4(tp.Appearance.FillColor.Normalize()),
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
