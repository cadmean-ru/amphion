package engine

const (
	FrontCallbackContextChange = -100
	FrontCallbackClick         = -101
)

type frontEndCallback struct {
	code int
	data string
}


