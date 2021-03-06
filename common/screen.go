package common

import "github.com/cadmean-ru/amphion/common/a"

type ScreenInfo struct {
	width, height int
}

func (s ScreenInfo) GetWidth() int {
	return s.width
}

func (s ScreenInfo) GetHeight() int {
	return s.height
}

func (s ScreenInfo) GetSize() a.IntVector2 {
	return a.NewIntVector2(s.width, s.height)
}

func (s ScreenInfo) FromMap(m map[string]interface{}) {
	s.width = m["width"].(int)
	s.height = m["height"].(int)
}

func NewScreenInfo(width, height int) ScreenInfo {
	return ScreenInfo{
		width:  width,
		height: height,
	}
}