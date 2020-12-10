package engine

const (
	FrontCallbackContextChange = -100
	FrontCallbackMouseDown     = -101
	FrontCallbackKeyDown       = -102
	FrontCallbackMouseUp       = -103
)

type frontEndCallback struct {
	code int
	data string
}


