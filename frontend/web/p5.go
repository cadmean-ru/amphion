// +build js

package web

import (
	"github.com/cadmean-ru/amphion/common/a"
	"syscall/js"
)

type p5 struct {
	clearJs          js.Value
	fillJs           js.Value
	strokeWeightJs   js.Value
	strokeJs         js.Value
	rectModeJs       js.Value
	rectModeCornerJs js.Value
	pointJs          js.Value
	lineJs           js.Value
	rectJs           js.Value
	ellipseJs        js.Value
	triangleJs       js.Value
	textJs           js.Value
	resizeCanvasJs   js.Value
	redrawJs         js.Value
	textSizeJs       js.Value
	loadImageJs      js.Value
	imageJs          js.Value
	textAlignJs      js.Value
	renderer         *P5Renderer
	onDraw           func(p5 *p5)
}

func (p *p5) prepare() {
	p.clearJs = js.Global().Get("clear")
	p.fillJs = js.Global().Get("fill")
	p.strokeWeightJs = js.Global().Get("strokeWeight")
	p.strokeJs = js.Global().Get("stroke")
	p.rectModeJs = js.Global().Get("rectMode")
	p.rectModeCornerJs = js.Global().Get("CORNER")
	p.pointJs = js.Global().Get("point")
	p.lineJs = js.Global().Get("line")
	p.rectJs = js.Global().Get("rect")
	p.ellipseJs = js.Global().Get("ellipse")
	p.triangleJs = js.Global().Get("triangle")
	p.textJs = js.Global().Get("text")
	p.resizeCanvasJs = js.Global().Get("resizeCanvas")
	p.redrawJs = js.Global().Get("redraw")
	p.textSizeJs = js.Global().Get("textSize")
	p.loadImageJs = js.Global().Get("loadImage")
	p.imageJs = js.Global().Get("image")
	p.textAlignJs = js.Global().Get("textAlign")
	js.Global().Set("goDraw", js.FuncOf(p.goDraw))
}

func (p *p5) clear() {
	p.clearJs.Invoke()
}

func (p *p5) fill(color a.Color) {
	p.fillJs.Invoke(color.R, color.G, color.B, color.A)
}

func (p *p5) strokeWeight(w int) {
	p.strokeWeightJs.Invoke(w)
}

func (p *p5) stroke(color a.Color) {
	p.strokeJs.Invoke(color.R, color.G, color.B, color.A)
}

func (p *p5) rectModeCorner() {
	p.rectModeJs.Invoke(p.rectModeCornerJs)
}

func (p *p5) point(x, y int) {
	p.pointJs.Invoke(x, y)
}

func (p *p5) line(x1, y1, x2, y2 int) {
	p.lineJs.Invoke(x1, y1, x2, y2)
}

func (p *p5) rect(x, y, sizeX, sizeY, cornerRadius int) {
	p.rectJs.Invoke(x, y, sizeX, sizeY, cornerRadius)
}

func (p *p5) ellipse(x, y, sizeX, sizeY int) {
	x1 := x + sizeX / 2
	y1 := y + sizeY / 2
	p.ellipseJs.Invoke(x1, y1, sizeX, sizeY)
}

func (p *p5) triangle(x1, y1, x2, y2, x3, y3 int) {
	p.triangleJs.Invoke(x1, y1, x2, y2, x3, y3)
}

func (p *p5) text(text string, x, y, sizeX, sizeY int) {
	p.textJs.Invoke(text, x, y, sizeX, sizeY)
}

func (p *p5) textSize(size int) {
	p.textSizeJs.Invoke(size)
}

func (p *p5) textAlign(hAlign, vAlign a.TextAlign) {
	var hAlignJs, vAlignJs js.Value

	switch hAlign {
	case a.TextAlignRight:
		hAlignJs = js.Global().Get("RIGHT")
	case a.TextAlignCenter:
		hAlignJs = js.Global().Get("CENTER")
	default:
		hAlignJs = js.Global().Get("LEFT")
	}

	switch vAlign {
	case a.TextAlignBottom:
		vAlignJs = js.Global().Get("BOTTOM")
	case a.TextAlignCenter:
		vAlignJs = js.Global().Get("CENTER")
	default:
		vAlignJs = js.Global().Get("TOP")
	}

	p.textAlignJs.Invoke(hAlignJs, vAlignJs)
}

func (p *p5) resizeCanvas(x, y int) {
	p.resizeCanvasJs.Invoke(x, y)
}

func (p *p5) loadImage(path string, callback func()) *p5image {
	img := &p5image{}
	imgJs := p.loadImageJs.Invoke(path, js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		img.ready = true
		callback()
		return nil
	}))
	img.value = imgJs
	return img
}

func (p *p5) image(img *p5image, x, y, w, h int) {
	p.imageJs.Invoke(img.value, x, y, w, h)
}

func (p *p5) redraw() {
	p.redrawJs.Invoke()
}

func (p *p5) goDraw(_ js.Value, _ []js.Value) interface{} {
	p.onDraw(p)
	return nil
}

type p5image struct {
	value js.Value
	ready bool
}

func (i *p5image) GetWidth() int {
	return i.value.Get("width").Int()
}

func (i *p5image) GetHeight() int {
	return i.value.Get("height").Int()
}