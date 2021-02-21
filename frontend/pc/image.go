package pc

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/all-core/gl"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
)

type ImageRenderer struct {
	*glPrimitiveRenderer
}

func (r *ImageRenderer) OnStart() {
	r.program = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(ImageVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(ImageFragShaderStr, gl.FRAGMENT_SHADER),
	)
}

func (r *ImageRenderer) OnRender(ctx *rendering.PrimitiveRenderingContext) {
	ip := ctx.Primitive.(*rendering.ImagePrimitive)
	state := ctx.State.(*glPrimitiveState)

	var texId uint32
	if state.tex == 0 {
		imagePath := ip.ImageUrl

		imageFile, err := os.Open(imagePath)
		if err != nil {
			log.Fatal(err)
		}

		defer imageFile.Close()

		img, _, err := image.Decode(imageFile)
		if err != nil {
			log.Fatal(err)
		}

		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

		gl.GenTextures(1, &texId)

		gl.BindTexture(gl.TEXTURE_2D, texId)

		gl.TexImage2D(
			gl.TEXTURE_2D,
			0,
			gl.SRGB_ALPHA,
			int32(rgba.Bounds().Size().X),
			int32(rgba.Bounds().Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)

		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

		gl.BindTexture(gl.TEXTURE_2D, 0)

		state.tex = texId
	} else {
		texId = state.tex
	}

	wSize := engine.GetScreenSize3()
	nPos := ip.Transform.Position.Ndc(wSize)
	brPos := ip.Transform.Position.Add(ip.Transform.Size).Ndc(wSize)

	drawTex(ctx, nPos, brPos, texId, r.program, nil)
}
