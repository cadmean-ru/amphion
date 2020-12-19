// +build windows linux darwin
// +build !android

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"sort"
)

type OpenGLPrimitive interface {
	GetType() common.AByte
	GetTransform() rendering.Transform
}

type OpenGLRenderer struct {
	defaultProgram uint32
	ellipseProgram uint32
	imageProgram   uint32
	window         *glfw.Window
	idgen          *common.IdGenerator
	primitives     map[int64]*glContainer
	shouldDelete   bool
	wSize          common.IntVector3
	fonts          map[string]*glFont
}

func (r *OpenGLRenderer) Prepare() {
	var err error

	r.window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		panic(err)
	}

	fmt.Println(gl.GoStr(gl.GetString(gl.VERSION)))
	fmt.Println(gl.GoStr(gl.GetString(gl.VENDOR)))
	fmt.Println(gl.GoStr(gl.GetString(gl.RENDERER)))

	gl.Viewport(0, 0, int32(r.wSize.X), int32(r.wSize.Y))

	r.defaultProgram = createAndLinkDefaultProgramOrPanic()

	r.ellipseProgram = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(DefaultVertexShaderText, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(EllipseFragShaderText, gl.FRAGMENT_SHADER),
	)

	r.imageProgram = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(ImageVertexShader, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(ImageFragShader, gl.FRAGMENT_SHADER),
	)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.ClearColor(1, 1, 1, 1)

	gl.Clear(gl.COLOR_BUFFER_BIT)

	r.window.SwapBuffers()

	r.idgen = common.NewIdGenerator()
}

func (r *OpenGLRenderer) AddPrimitive() int64 {
	id := r.idgen.NextId()
	r.primitives[id] = newGlContainer()
	return id
}

func (r *OpenGLRenderer) SetPrimitive(id int64, primitive interface{}, shouldRedraw bool) {
	if !shouldRedraw {
		return
	}

	if _, ok := r.primitives[id]; ok {
		r.primitives[id].primitive = primitive.(OpenGLPrimitive)
		r.primitives[id].redraw = true
	} else {
		//panic(fmt.Sprintf("Primitive with id %d was not found.\nAdded primitives:\n%+v", id, r.primitives))
		fmt.Printf("Warning! Primitive with id %d was not found.\n", id)
	}
}

func (r *OpenGLRenderer) RemovePrimitive(id int64) {
	r.primitives[id].free()
	delete(r.primitives, id)
}

func (r *OpenGLRenderer) PerformRendering() {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	list := make([]*glContainer, len(r.primitives))

	i := 0
	for _, p := range r.primitives {
		list[i] = p
		i++
	}

	// TODO: replace this sort with something like binary search tree
	sort.Slice(list, func(i, j int) bool {
		z1 := list[i].primitive.GetTransform().Position.Z
		z2 := list[j].primitive.GetTransform().Position.Z
		return z1 < z2
	})

	for _, p := range list {
		switch p.primitive.GetType() {
		case rendering.PrimitivePoint:
			break
		case rendering.PrimitiveLine:
			r.drawLine(p)
		case rendering.PrimitiveRectangle:
			r.drawRectangle(p)
		case rendering.PrimitiveEllipse:
			r.drawEllipse(p)
		case rendering.PrimitiveTriangle:
			r.drawTriangle(p)
		case rendering.PrimitiveText:
			r.drawText(p)
		case rendering.PrimitiveImage:
			r.drawImage(p)
		case rendering.PrimitiveBezier:
			break
		}

		p.redraw = false
	}

	r.window.SwapBuffers()
}

func (r *OpenGLRenderer) Clear() {
	r.primitives = make(map[int64]*glContainer)
	r.idgen = common.NewIdGenerator()
}

func (r *OpenGLRenderer) Stop() {
	glfw.Terminate()
}

func (r *OpenGLRenderer) handleWindowResize() {
	gl.Viewport(0, 0, int32(r.wSize.X), int32(r.wSize.Y))
}

func (r *OpenGLRenderer) drawRectangle(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	p.gen()

	gl.UseProgram(r.defaultProgram)

	if p.redraw {
		gl.BindVertexArray(p.vao)

		nPos := gp.Transform.Position.Ndc(r.wSize)
		nSize := gp.Transform.Position.Add(gp.Transform.Size).Ndc(r.wSize)

		vertices := []float32 {
			nPos.X, nPos.Y, 0,
			nPos.X, nSize.Y, 0,
			nSize.X, nSize.Y, 0,
			nSize.X, nPos.Y, 0,
		}

		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
		gl.EnableVertexAttribArray(0)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	if p.colorLoc == -1 {
		p.colorLoc = gl.GetUniformLocation(r.defaultProgram, gl.Str("ourColor\x00"))
	}

	color := gp.Appearance.FillColor.Normalize()
	gl.Uniform4f(p.colorLoc, color.X, color.Y, color.Z, color.W)

	gl.BindVertexArray(p.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (r *OpenGLRenderer) drawTriangle(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	p.gen()

	if p.redraw {
		gl.BindVertexArray(p.vao)

		nPos := gp.Transform.Position.Ndc(r.wSize)
		nSize := gp.Transform.Position.Add(gp.Transform.Size).Ndc(r.wSize)
		midX := nPos.X + ((nSize.X - nPos.X) / 2)

		vertices := []float32 {
			nPos.X, nSize.Y, 0,
			midX, nPos.Y, 0,
			nSize.X, nSize.Y, 0,
		}

		indices := []uint32 {
			0, 1, 2,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
		gl.EnableVertexAttribArray(0)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.UseProgram(r.defaultProgram)

	if p.colorLoc == -1 {
		p.colorLoc = gl.GetUniformLocation(r.defaultProgram, gl.Str("ourColor\x00"))
	}

	color := gp.Appearance.FillColor.Normalize()
	gl.Uniform4f(p.colorLoc, color.X, color.Y, color.Z, color.W)

	gl.BindVertexArray(p.vao)

	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, nil)
}

func (r *OpenGLRenderer) drawEllipse(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	p.gen()

	gl.UseProgram(r.ellipseProgram)

	if p.redraw {
		gl.BindVertexArray(p.vao)

		nPos := gp.Transform.Position.Ndc(r.wSize)
		nSize := gp.Transform.Position.Add(gp.Transform.Size).Ndc(r.wSize)

		p.other["nPos"] = nPos
		p.other["nSize"] = nSize

		vertices := []float32 {
			nPos.X, nPos.Y, 0,
			nPos.X, nSize.Y, 0,
			nSize.X, nSize.Y, 0,
			nSize.X, nPos.Y, 0,
		}

		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
		gl.EnableVertexAttribArray(0)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	if p.colorLoc == -1 {
		p.colorLoc = gl.GetUniformLocation(r.ellipseProgram, gl.Str("ourColor\x00"))
	}

	if _, ok := p.other["tlLoc"]; !ok {
		p.other["tlLoc"] = gl.GetUniformLocation(r.ellipseProgram, gl.Str("tlPos\x00"))
		p.other["brLoc"] = gl.GetUniformLocation(r.ellipseProgram, gl.Str("brPos\x00"))
	}

	color := gp.Appearance.FillColor.Normalize()
	gl.Uniform4f(p.colorLoc, color.X, color.Y, color.Z, color.W)

	tlLoc := p.other["tlLoc"].(int32)
	brLoc := p.other["brLoc"].(int32)
	tlPos := p.other["nPos"].(common.Vector3)
	brPos := p.other["nSize"].(common.Vector3)

	gl.Uniform3f(tlLoc, tlPos.X, tlPos.Y, tlPos.Z)
	gl.Uniform3f(brLoc, brPos.X, brPos.Y, brPos.Z)

	gl.BindVertexArray(p.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (r *OpenGLRenderer) drawLine(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	p.gen()

	gl.UseProgram(r.defaultProgram)

	if p.redraw {
		gl.BindVertexArray(p.vao)

		nPos := gp.Transform.Position.Ndc(r.wSize)
		nSize := gp.Transform.Position.Add(gp.Transform.Size).Ndc(r.wSize)

		p.other["nPos"] = nPos
		p.other["nSize"] = nSize

		vertices := []float32 {
			nPos.X, nPos.Y, 0,
			nSize.X, nSize.Y, 0,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 12, nil)
		gl.EnableVertexAttribArray(0)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	if p.colorLoc == -1 {
		p.colorLoc = gl.GetUniformLocation(r.defaultProgram, gl.Str("ourColor\x00"))
	}

	color := gp.Appearance.StrokeColor.Normalize()
	gl.Uniform4f(p.colorLoc, color.X, color.Y, color.Z, color.W)

	gl.BindVertexArray(p.vao)

	gl.DrawArrays(gl.LINE, 0, 2)
}

func (r *OpenGLRenderer) drawText(p *glContainer) {
	//tp := p.primitive.(*rendering.TextPrimitive)
	//
	//fontName := tp.TextAppearance.Font
	//if fontName == "" {
	//	fontName = "Arial"
	//}
	//
	//var font *glFont
	//var ok bool
	//var err error
	//if font, ok = r.fonts[fontName]; !ok {
	//	font, err = loadFont(fontName)
	//	if err != nil {
	//		panic(err)
	//	}
	//	r.fonts[fontName] = font
	//}
	//
	//if _, ok = p.other["japanText"]; !ok {
	//	p.other["japanText"] = v41.NewText(font.japanFont, 1, 1)
	//}
	//
	//text := p.other["japanText"].(*v41.Text)
	//
	//if p.redraw {
	//	text.SetString(string(tp.Text))
	//	color := tp.Appearance.FillColor.Normalize()
	//	text.SetColor(mgl32.Vec3{color.X, color.Y, color.Z})
	//	nPos := tp.Transform.Position.Ndc(r.wSize)
	//	text.SetPosition(mgl32.Vec2{nPos.X, nPos.Y})
	//}
	//
	//text.Draw()
}

func (r *OpenGLRenderer) drawImage(p *glContainer) {
	ip := p.primitive.(*rendering.ImagePrimitive)

	var texId uint32
	if _, ok := p.other["tex"]; !ok {

		files, err := ioutil.ReadDir("./res")
		if err != nil {
			log.Fatal(err)
		}

		imagePath := "./res/" + files[ip.ResIndex].Name()

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

		p.other["tex"] = texId
	} else {
		texId = p.other["tex"].(uint32)
	}

	p.gen()

	if p.redraw {
		gl.BindVertexArray(p.vao)

		nPos := ip.Transform.Position.Ndc(r.wSize)
		nSize := ip.Transform.Position.Add(ip.Transform.Size).Ndc(r.wSize)

		vertices := []float32 {
			nPos.X,  nPos.Y,  0,	0, 1,  // top left
			nPos.X,  nSize.Y, 0,	0, 0,  // bottom left
			nSize.X, nSize.Y, 0,	1, 0,  // top right
			nSize.X, nPos.Y,  0,	1, 1,  // bottom right
		}

		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.UseProgram(r.imageProgram)
	gl.BindVertexArray(p.vao)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}