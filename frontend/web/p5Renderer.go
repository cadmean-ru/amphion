//+build js

package web

import (
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
	"sort"
)

type P5Renderer struct {
	p5           *p5
	idgen        *common.IdGenerator
	primitives   map[int]*p5Container
	prevFontSize byte
	images       map[string]*p5image
	front        *Frontend
}

func (r *P5Renderer) Prepare() {
	r.p5.renderer = r
	r.p5.prepare()
	r.p5.onDraw = r.drawP5
}

func (r *P5Renderer) AddPrimitive() int {
	id := r.idgen.NextId()
	r.primitives[id] = newP5Container()
	return id
}

func (r *P5Renderer) SetPrimitive(id int, primitive rendering.IPrimitive, shouldRerender bool) {
	if !shouldRerender {
		return
	}

	if _, ok := r.primitives[id]; ok {
		r.primitives[id].primitive = primitive
		r.primitives[id].redraw = true
	}
}

func (r *P5Renderer) RemovePrimitive(id int) {
	delete(r.primitives, id)
}

func (r *P5Renderer) PerformRendering() {
	r.p5.redraw()
}

func (r *P5Renderer) drawP5(p5 *p5) {
	p5.clear()

	p5.rectModeCorner()

	count := 0
	for _, p := range r.primitives {
		if p.primitive == nil {
			continue
		}

		count++
	}

	list := make([]*p5Container, count)

	i := 0
	for _, p := range r.primitives {
		if p.primitive == nil {
			continue
		}
		list[i] = p
		i++
	}

	sort.Slice(list, func(i, j int) bool {
		z1 := list[i].primitive.GetTransform().Position.Z
		z2 := list[j].primitive.GetTransform().Position.Z
		return z1 < z2
	})

	for _, p := range list {
		t := p.primitive.GetTransform()
		pos := t.Position
		size := t.Size

		switch p.primitive.GetType() {
		case rendering.PrimitivePoint:
			point := p.primitive.(*rendering.GeometryPrimitive)
			p5.fill(point.Appearance.FillColor)
			p5.point(pos.X, pos.Y)
		case rendering.PrimitiveLine:
			line := p.primitive.(*rendering.GeometryPrimitive)
			p5.fill(line.Appearance.StrokeColor)
			p5.strokeWeight(int(line.Appearance.StrokeWeight))
			x2, y2 := pos.X+size.X, pos.Y+size.Y
			p5.line(pos.X, pos.Y, x2, y2)
		case rendering.PrimitiveRectangle:
			rect := p.primitive.(*rendering.GeometryPrimitive)
			p5.fill(rect.Appearance.FillColor)
			p5.strokeWeight(int(rect.Appearance.StrokeWeight))
			p5.stroke(rect.Appearance.StrokeColor)
			p5.rect(pos.X, pos.Y, size.X, size.Y, int(rect.Appearance.CornerRadius))
		case rendering.PrimitiveEllipse:
			ellipse := p.primitive.(*rendering.GeometryPrimitive)
			p5.fill(ellipse.Appearance.FillColor)
			p5.strokeWeight(int(ellipse.Appearance.StrokeWeight))
			p5.stroke(ellipse.Appearance.StrokeColor)
			p5.ellipse(pos.X, pos.Y, size.X, size.Y)
		case rendering.PrimitiveTriangle:
			triangle := p.primitive.(*rendering.GeometryPrimitive)
			p5.fill(triangle.Appearance.FillColor)
			p5.strokeWeight(int(triangle.Appearance.StrokeWeight))
			p5.stroke(triangle.Appearance.StrokeColor)
			tx1 := pos.X
			ty1 := pos.Y + size.Y
			tx2 := pos.X + (size.X / 2)
			ty2 := pos.Y
			tx3 := pos.X + size.X
			ty3 := pos.Y + size.Y
			p5.triangle(tx1, ty1, tx2, ty2, tx3, ty3)
		case rendering.PrimitiveText:
			tp := p.primitive.(*rendering.TextPrimitive)
			p5.fill(tp.Appearance.FillColor)
			p5.strokeWeight(int(tp.Appearance.StrokeWeight))
			if tp.TextAppearance.FontSize != r.prevFontSize {
				r.prevFontSize = tp.TextAppearance.FontSize
				p5.textSize(int(tp.TextAppearance.FontSize))
			}
			p5.textAlign(tp.HTextAlign, tp.VTextAlign)
			p5.text(tp.Text, pos.X, pos.Y, size.X, size.Y)
		case rendering.PrimitiveImage:
			ip := p.primitive.(*rendering.ImagePrimitive)
			if img, ok := r.images[ip.ImageUrl]; ok {
				if img.ready {
					p5.image(img, pos.X, pos.Y, size.X, size.Y)
				}
			} else {
				r.images[ip.ImageUrl] = p5.loadImage(ip.ImageUrl, func() {
					engine.GetInstance().RequestRendering()
				})
			}
		case rendering.PrimitiveBezier:

		}

		p.redraw = false
	}
}

func (r *P5Renderer) Clear() {
	r.primitives = make(map[int]*p5Container)
	r.idgen = common.NewIdGenerator()
	r.p5.redraw()
}

func (r *P5Renderer) Stop() {

}

func newP5Renderer(front *Frontend) *P5Renderer {
	return &P5Renderer{
		p5:         &p5{},
		primitives: make(map[int]*p5Container),
		idgen:      common.NewIdGenerator(),
		images:     make(map[string]*p5image),
		front:      front,
	}
}
