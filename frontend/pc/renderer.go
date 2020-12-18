// +build windows linux darwin

package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/rendering"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

type OpenGLPrimitive interface {
	GetType() common.AByte
}

type OpenGLRenderer struct {
	defaultProgram uint32
	window         *glfw.Window
	idgen          *common.IdGenerator
	primitives     map[int64]*glContainer
	shouldDelete   bool
	wSize          common.IntVector3
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

	for _, p := range r.primitives {
		switch p.primitive.GetType() {
		case rendering.PrimitiveRectangle:
			r.drawRectangle(p)
		}
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

func (r *OpenGLRenderer) drawRectangle(p *glContainer) {
	gp := p.primitive.(*rendering.GeometryPrimitive)

	if p.vbo == 0 {
		temp := make([]uint32, 2)
		gl.GenBuffers(2, &temp[0])
		p.vbo = temp[0]
		p.ebo = temp[1]
		gl.GenVertexArrays(1, &p.vao)
	}

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

	gl.UseProgram(r.defaultProgram)
	gl.BindVertexArray(p.vao)

	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
}
