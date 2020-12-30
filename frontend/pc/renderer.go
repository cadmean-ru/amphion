// +build windows linux darwin
// +build !android

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
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

//go:generate ../../build/darwin/test --generate shaders -i ./shaders -o ./shaders.gen.go --package pc

type OpenGLPrimitive interface {
	GetType() a.Byte
	GetTransform() rendering.Transform
}

type OpenGLRenderer struct {
	defaultProgram uint32
	ellipseProgram uint32
	imageProgram   uint32
	textProgram    uint32
	shapeProgram   uint32
	window         *glfw.Window
	idgen          *common.IdGenerator
	primitives     map[int64]*glContainer
	shouldDelete   bool
	wSize          common.IntVector3
	fonts          map[string]*glFont
	projection     [16]float32
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
		createAndCompileShaderOrPanic(EllipseVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(EllipseFragShaderStr, gl.FRAGMENT_SHADER),
	)

	r.imageProgram = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(ImageVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(ImageFragShaderStr, gl.FRAGMENT_SHADER),
	)

	r.textProgram = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(TextVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(TextFragShaderStr, gl.FRAGMENT_SHADER),
	)

	r.shapeProgram = createAndLinkProgramOrPanic(
		createAndCompileShaderOrPanic(ShapeVertexShaderStr, gl.VERTEX_SHADER),
		createAndCompileShaderOrPanic(ShapeFragShaderStr, gl.FRAGMENT_SHADER),
	)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)

	gl.ClearColor(1, 1, 1, 1)

	gl.Clear(gl.COLOR_BUFFER_BIT)

	r.calculateProjection()

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
	//r.primitives[id].free()
	delete(r.primitives, id)
}

func (r *OpenGLRenderer) PerformRendering() {
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	count := 0
	for _, p := range r.primitives {
		if p.primitive == nil {
			//fmt.Println(fmt.Sprintf("primitive data was never set (id: %d)", id))
			continue
		}

		count++
	}

	list := make([]*glContainer, count)

	i := 0
	for _, p := range r.primitives {
		if p.primitive == nil {
			continue
		}
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

}

func (r *OpenGLRenderer) handleWindowResize() {
	gl.Viewport(0, 0, int32(r.wSize.X), int32(r.wSize.Y))
	r.calculateProjection()
}

func (r *OpenGLRenderer) drawRectangle(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	p.gen()

	gl.UseProgram(r.shapeProgram)

	if p.redraw {
		gl.BindVertexArray(p.vao)

		ntlPos := gp.Transform.Position.Ndc(r.wSize) // normalized top left
		nbrPos := gp.Transform.Position.Add(gp.Transform.Size).Ndc(r.wSize) // normalized bottom right

		//tlPos := gp.Transform.Position.ToInt32() // top left
		//brPos := gp.Transform.Position.Add(gp.Transform.Size).ToInt32() // bottom right
		//
		color := gp.Appearance.FillColor
		r1 := float32(color.R)
		g1 := float32(color.G)
		b1 := float32(color.B)
		a1 := float32(color.A)
		//
		//strokeColor := gp.Appearance.StrokeColor
		//r2 := int32(strokeColor.R)
		//g2 := int32(strokeColor.G)
		//b2 := int32(strokeColor.B)
		//a2 := int32(strokeColor.A)
		//
		//sw := int32(gp.Appearance.StrokeWeight)
		//cr := int32(gp.Appearance.CornerRadius)

		vertices := []float32 {
			ntlPos.X, ntlPos.Y, 0, r1, g1, b1, a1,
			ntlPos.X, nbrPos.Y, 0, r1, g1, b1, a1,
			nbrPos.X, nbrPos.Y, 0, r1, g1, b1, a1,
			nbrPos.X, ntlPos.Y, 0, r1, g1, b1, a1,
		}
		
		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		const stride int32 = 28

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 4, gl.FLOAT, false, stride, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.BindVertexArray(p.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (r *OpenGLRenderer) drawTriangle(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	p.gen()

	gl.UseProgram(r.shapeProgram)

	if p.redraw {
		gl.BindVertexArray(p.vao)

		nPos := gp.Transform.Position.Ndc(r.wSize)
		nSize := gp.Transform.Position.Add(gp.Transform.Size).Ndc(r.wSize)
		midX := nPos.X + ((nSize.X - nPos.X) / 2)

		color := gp.Appearance.FillColor
		r1 := float32(color.R)
		g1 := float32(color.G)
		b1 := float32(color.B)
		a1 := float32(color.A)

		vertices := []float32 {
			nPos.X,  nSize.Y, 0, r1, g1, b1, a1,
			midX,    nPos.Y,  0, r1, g1, b1, a1,
			nSize.X, nSize.Y, 0, r1, g1, b1, a1,
		}

		indices := []uint32 {
			0, 1, 2,
		}

		const stride int32 = 28

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 4, gl.FLOAT, false, stride, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

	gl.BindVertexArray(p.vao)

	gl.DrawElements(gl.TRIANGLES, 3, gl.UNSIGNED_INT, nil)
}

func (r *OpenGLRenderer) drawEllipse(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	p.gen()

	gl.UseProgram(r.ellipseProgram)

	if p.redraw {
		gl.BindVertexArray(p.vao)

		tlPos := gp.Transform.Position.Ndc(r.wSize)
		brPos := gp.Transform.Position.Add(gp.Transform.Size).Ndc(r.wSize)

		color := gp.Appearance.FillColor
		r1 := float32(color.R)
		g1 := float32(color.G)
		b1 := float32(color.B)
		a1 := float32(color.A)

		vertices := []float32 {
			tlPos.X, tlPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
			tlPos.X, brPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
			brPos.X, brPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
			brPos.X, tlPos.Y, 0, r1, g1, b1, a1, tlPos.X, tlPos.Y, 0, brPos.X, brPos.Y, 0,
		}

		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		const stride int32 = 52

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, stride, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 4, gl.FLOAT, false, stride, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.VertexAttribPointer(2, 3, gl.FLOAT, false, stride, gl.PtrOffset(28))
		gl.EnableVertexAttribArray(2)

		gl.VertexAttribPointer(3, 3, gl.FLOAT, false, stride, gl.PtrOffset(40))
		gl.EnableVertexAttribArray(3)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	}

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
	tp := p.primitive.(*rendering.TextPrimitive)

	fontName := tp.TextAppearance.Font
	if fontName == "" {
		fontName = "Arial"
	}

	var font *glFont
	var ok bool
	var err error
	if font, ok = r.fonts[fontName]; !ok {
		font, err = loadFont(fontName)
		if err != nil {
			panic(err)
		}
		r.fonts[fontName] = font
	}

	//p.gen()
	//
	//if p.redraw {
	//	gl.BindVertexArray(p.vao)
	//	gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
	//	gl.BufferData(gl.ARRAY_BUFFER, 4*6*4, nil, gl.DYNAMIC_DRAW)
	//	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 4*4, nil)
	//	gl.EnableVertexAttribArray(0)
	//	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	//	gl.BindVertexArray(0)
	//}

	color := tp.Appearance.FillColor.Normalize()

	//gl.UseProgram(r.textProgram)
	//gl.ActiveTexture(gl.TEXTURE0)
	//gl.BindVertexArray(p.vao)

	x := tp.Transform.Position.X
	y := tp.Transform.Position.Y

	runes := []rune(tp.Text)

	for i, c := range runes {
		if c == ' ' {
			x += 10
			continue
		}

		ch, ok1 := font.characters[c]

		if !ok1 {
			//fmt.Printf("Char not found %+v\n", c)
			continue
		}

		//fmt.Println(string(c))

		xpos := x + ch.bearing.X
		ypos := y

		pos := common.NewIntVector3(xpos, ypos, 0)
		npos := pos.Ndc(r.wSize)

		w := ch.size.X
		h := ch.size.Y

		brpos := pos.Add(common.NewIntVector3(w, h, 0))
		nbrpos := brpos.Ndc(r.wSize)

		r.drawTex(p, npos, nbrpos, ch.textureId, r.textProgram, func() {
			gl.Uniform3f(gl.GetUniformLocation(r.textProgram, gl.Str("textColor\x00")), color.X, color.Y, color.Z)
		})

		k := 0
		if i < len(tp.Text) - 1 {
			k = font.kern(c, runes[i + 1])
		}
		x += ch.advance.X + k
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
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

	nPos := ip.Transform.Position.Ndc(r.wSize)
	brPos := ip.Transform.Position.Add(ip.Transform.Size).Ndc(r.wSize)

	r.drawTex(p, nPos, brPos, texId, r.imageProgram, nil)
}

func (r *OpenGLRenderer) drawTex(p *glContainer, nPos, nbrPos common.Vector3, texId, progId uint32, beforeDraw func()) {
	p.gen()

	//if p.redraw {
		gl.BindVertexArray(p.vao)

		vertices := []float32 {
			nPos.X,   nPos.Y,   0,	0, 0, // top left
			nPos.X,   nbrPos.Y, 0,	0, 1, // bottom left
			nbrPos.X, nbrPos.Y, 0,	1, 1, // top right
			nbrPos.X, nPos.Y,   0,	1, 0, // bottom right
		}

		indices := []uint32 {
			0, 1, 2,
			0, 3, 2,
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, p.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.DYNAMIC_DRAW)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, p.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.DYNAMIC_DRAW)

		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 20, nil)
		gl.EnableVertexAttribArray(0)

		gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 20, gl.PtrOffset(12))
		gl.EnableVertexAttribArray(1)

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindVertexArray(0)
	//}

	gl.UseProgram(progId)
	gl.BindVertexArray(p.vao)
	gl.BindTexture(gl.TEXTURE_2D, texId)

	if beforeDraw != nil {
		beforeDraw()
	}

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}

func (r *OpenGLRenderer) calculateProjection() {
	xs := float32(r.wSize.X)
	ys := float32(r.wSize.Y)
	c1 := 2 / xs
	c2 := 2 / ys
	r.projection = [16]float32 {
		c1, 0,  0,  -1,
		0,  c2, 0,  -1,
		0,  0,  1,  0,
		0,  0,  0,  1,
	}

	//fmt.Println(r.projection)

	r.setProjectionUniform(r.shapeProgram)
}

func (r *OpenGLRenderer) setProjectionUniform(program uint32) {
	gl.UseProgram(r.shapeProgram)
	gl.UniformMatrix4fv(gl.GetUniformLocation(program, gl.Str("uProjection\x00")), 1, false, &r.projection[0])
	gl.UseProgram(0)
}