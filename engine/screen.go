package engine

type ScreenInfo struct {
	width, height int
}

func (s ScreenInfo) GetWidth() int {
	return s.width
}

func (s ScreenInfo) GetHeight() int {
	return s.height
}

func (s ScreenInfo) FromMap(m map[string]interface{}) {
	s.width = m["width"].(int)
	s.height = m["height"].(int)
}

func newScreenInfo(width, height int) ScreenInfo {
	return ScreenInfo{
		width:  width,
		height: height,
	}
}