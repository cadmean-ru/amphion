//+build js

package web

import "github.com/cadmean-ru/amphion/rendering"

type p5Container struct {
	primitive rendering.Primitive
	redraw    bool
}

func newP5Container() *p5Container {
	return &p5Container{}
}