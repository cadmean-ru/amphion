package web

import (
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/rendering"
)

func drawPoint(p5 *p5, primitive rendering.IPrimitive) {
	t := primitive.GetTransform()
	pos := t.Position

	point := primitive.(*rendering.GeometryPrimitive)
	p5.fill(point.Appearance.FillColor)
	p5.point(pos.X, pos.Y)
}

func drawLine(p5 *p5, primitive rendering.IPrimitive) {
	t := primitive.GetTransform()
	pos := t.Position
	size := t.Size

	line := primitive.(*rendering.GeometryPrimitive)
	p5.fill(line.Appearance.StrokeColor)
	p5.strokeWeight(int(line.Appearance.StrokeWeight))
	x2, y2 := pos.X+size.X, pos.Y+size.Y
	p5.line(pos.X, pos.Y, x2, y2)
}

func drawRectangle(p5 *p5, primitive rendering.IPrimitive) {
	t := primitive.GetTransform()
	pos := t.Position
	size := t.Size

	rect := primitive.(*rendering.GeometryPrimitive)
	p5.fill(rect.Appearance.FillColor)
	p5.strokeWeight(int(rect.Appearance.StrokeWeight))
	p5.stroke(rect.Appearance.StrokeColor)
	p5.rect(pos.X, pos.Y, size.X, size.Y, int(rect.Appearance.CornerRadius))
}

func drawEllipse(p5 *p5, primitive rendering.IPrimitive) {
	t := primitive.GetTransform()
	pos := t.Position
	size := t.Size

	ellipse := primitive.(*rendering.GeometryPrimitive)
	p5.fill(ellipse.Appearance.FillColor)
	p5.strokeWeight(int(ellipse.Appearance.StrokeWeight))
	p5.stroke(ellipse.Appearance.StrokeColor)
	p5.ellipse(pos.X, pos.Y, size.X, size.Y)
}

func drawTriangle(p5 *p5, primitive rendering.IPrimitive) {
	t := primitive.GetTransform()
	pos := t.Position
	size := t.Size

	triangle := primitive.(*rendering.GeometryPrimitive)
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
}

var prevFontSize byte
func drawText(p5 *p5, primitive rendering.IPrimitive) {
	t := primitive.GetTransform()
	pos := t.Position
	size := t.Size

	tp := primitive.(*rendering.TextPrimitive)
	p5.fill(tp.Appearance.FillColor)
	p5.strokeWeight(int(tp.Appearance.StrokeWeight))
	if tp.TextAppearance.FontSize != prevFontSize {
		prevFontSize = tp.TextAppearance.FontSize
		p5.textSize(int(tp.TextAppearance.FontSize))
	}
	p5.textAlign(tp.HTextAlign, tp.VTextAlign)
	p5.text(tp.Text, pos.X, pos.Y, size.X, size.Y)
}

var images = map[string]*p5image{}
func drawImage(p5 *p5, primitive rendering.IPrimitive) {
	t := primitive.GetTransform()
	pos := t.Position
	size := t.Size

	ip := primitive.(*rendering.ImagePrimitive)
	if img, ok := images[ip.ImageUrl]; ok {
		if img.ready {
			p5.image(img, pos.X, pos.Y, size.X, size.Y)
		}
	} else {
		images[ip.ImageUrl] = p5.loadImage(ip.ImageUrl, func() {
			engine.GetInstance().RequestRendering()
		})
	}
}
