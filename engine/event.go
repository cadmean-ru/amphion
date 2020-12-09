package engine

const (
	EventUpdate      = -1
	EventRender      = -2
	EventCloseScene  = -3
	EventClick       = -4
	EventDoubleClick = -5
	EventNavigate    = -6
	EventStop        = -7
)

type AmphionEvent struct {
	Sender interface{}
	Code   int
	Data   interface{}
}

func NewAmphionEvent(from interface{}, code int, data interface{}) AmphionEvent {
	return AmphionEvent{
		Sender: from,
		Code:   code,
		Data:   data,
	}
}

type EventHandler func(event AmphionEvent) bool

//type EventBinder interface {
//	Prepare()
//	Bind(Code int, handler EventHandler)
//	GetHandlers(Code int) []EventHandler
//}

type EventBinder struct {
	handlers map[int][]EventHandler
}

func (b *EventBinder) Bind(code int, handler EventHandler) {
	if _, ok := b.handlers[code]; !ok {
		b.handlers[code] = make([]EventHandler, 1, 10)
		b.handlers[code][0] = handler
	} else {
		b.handlers[code] = append(b.handlers[code], handler)
	}
}

func (b *EventBinder) GetHandlers(code int) []EventHandler {
	if h, ok := b.handlers[code]; ok {
		return h
	}

	return make([]EventHandler, 0)
}

func (b *EventBinder) InvokeHandlers(event AmphionEvent) {
	for _, h := range b.GetHandlers(event.Code) {
		if !h(event) {
			break
		}
	}
}

func newEventBinder() *EventBinder {
	return &EventBinder{
		handlers: make(map[int][]EventHandler),
	}
}

//func PrepareFrontendCallbackHandlerForPlatform(platform common.Platform) {
//	switch platform.GetName() {
//	case "web":
//		PrepareJsCallbackHandler()
//		break
//	}
//}
//
//func PrepareJsCallbackHandler() {
//	js.Global().Set("frontEndCallback", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
//		e := args[0]
//		Code := e.Get("Code").Int()
//		jsEvent := e.Get("Data").String()
//		event := AmphionEvent{
//			Code: Code,
//			Data: jsEvent,
//		}
//		instance.handleFrontEndCallback(event)
//		return nil
//	}))
//}