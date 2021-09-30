package cli

import "github.com/cadmean-ru/amphion/rendering"

type PrimitiveRenderingContext struct {
	GeometryPrimitiveData *GeometryPrimitiveData
	ImagePrimitiveData    *ImagePrimitiveData
	TextPrimitiveData     *TextPrimitiveData
	PrimitiveKind         int
	State                 interface{}
	Redraw                bool
	PrimitiveId           int
	Projection            [16]float32
}

func newCliPrimitiveRenderingContext(ctx *rendering.PrimitiveRenderingContext) *PrimitiveRenderingContext {
	cliCtx := &PrimitiveRenderingContext{
		PrimitiveKind: int(ctx.PrimitiveKind),
		PrimitiveId:   ctx.PrimitiveId,
		Redraw:        ctx.Redraw,
		Projection:    ctx.Projection,
	}

	t := ctx.Primitive.GetTransform()
	tlPos := t.Position.ToFloat()
	brPos := t.Position.Add(t.Size).ToFloat()

	switch ctx.PrimitiveKind {
	case rendering.PrimitiveTriangle, rendering.PrimitiveRectangle, rendering.PrimitiveEllipse, rendering.PrimitiveLine:
		gp := ctx.Primitive.(*rendering.GeometryPrimitive)

		cliCtx.GeometryPrimitiveData = &GeometryPrimitiveData{
			GeometryType: int(ctx.PrimitiveKind),
			TlPosition: NewVector3FromAVector3(tlPos),
			BrPosition: NewVector3FromAVector3(brPos),
			FillColor: NewVector4FromAColor(gp.Appearance.FillColor),
			StrokeColor: NewVector4FromAColor(gp.Appearance.StrokeColor),
			StrokeWeight: int(gp.Appearance.StrokeWeight),
			CornerRadius: int(gp.Appearance.CornerRadius),
		}
	case rendering.PrimitiveImage:
		ip := ctx.Primitive.(*rendering.ImagePrimitive)

		cliCtx.ImagePrimitiveData = &ImagePrimitiveData{
			TlPosition: NewVector3FromAVector3(tlPos),
			BrPosition: NewVector3FromAVector3(brPos),
			Index: ip.Index,
		}
		cliCtx.ImagePrimitiveData.Bitmaps = make([]*Bitmap, len(ip.Bitmaps))
		for i, b := range ip.Bitmaps {
			cliCtx.ImagePrimitiveData.Bitmaps[i] = newBitmap(b)
		}
	case rendering.PrimitiveText:
		tp := ctx.Primitive.(*rendering.TextPrimitive)

		cliCtx.TextPrimitiveData = &TextPrimitiveData{
			Text: tp.Text,
			TlPosition: NewVector3FromAIntVector3(t.Position),
			Size: NewVector3FromAIntVector3(t.Size),
			TextColor: NewVector4FromAColor(tp.Appearance.FillColor),
			Provider:   tp.TextProvider,
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