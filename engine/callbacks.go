package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/frontend"
	"strconv"
	"strings"
)

type frontendCallbackHandler struct {
	handlersMap map[int]dispatch.MessageHandlerFunc
}

func (c *frontendCallbackHandler) SendMessage(message *dispatch.Message) {
	if handler, ok := c.handlersMap[message.What]; ok {
		handler(message)
	}
}

func handleMouseDown(callback *dispatch.Message) {
	pos := parseCursorPositionData(callback.StrData)
	instance.inputManager.reportCursorPosition(pos)
	var eventCode int
	if callback.What == frontend.CallbackMouseDown {
		eventCode = EventMouseDown
	} else {
		eventCode = EventTouchDown
	}
	instance.handleClickEvent(pos, eventCode)
}

func handleMouseUp(callback *dispatch.Message) {
	pos := parseCursorPositionData(callback.StrData)
	instance.inputManager.reportCursorPosition(pos)
	var eventCode int
	if callback.What == frontend.CallbackMouseUp {
		eventCode = EventMouseUp
	} else {
		eventCode = EventTouchUp
	}
	event := NewAmphionEvent(instance, eventCode, MouseEventData{
		MousePosition: pos,
		SceneObject:   nil,
	})
	instance.updateRoutine.enqueueEventAndRequestUpdate(event)
}

func handleMouseMove(callback *dispatch.Message) {
	pos := parseCursorPositionData(callback.StrData)
	instance.inputManager.reportCursorPosition(pos)
	var eventCode int
	if callback.What == frontend.CallbackMouseMove {
		instance.handleMouseMove(pos)
		eventCode = EventMouseMove
	} else {
		eventCode = EventTouchMove
	}
	event := NewAmphionEvent(instance, eventCode, MouseEventData{
		MousePosition: pos,
		SceneObject:   nil,
	})
	instance.updateRoutine.enqueueEventAndRequestUpdate(event)
}

func parseCursorPositionData(data string) a.IntVector2 {
	coords := strings.Split(data, ";")
	if len(coords) != 2 {
		panic("Invalid click callback Data")
	}
	x, err := strconv.ParseInt(coords[0], 10, 32)
	if err != nil {
		panic("Invalid click callback Data")
	}
	y, err := strconv.ParseInt(coords[1], 10, 32)
	if err != nil {
		panic("Invalid click callback Data")
	}
	return a.NewIntVector2(int(x), int(y))
}

func handleContextChange(_ *dispatch.Message) {
	instance.globalContext = instance.front.GetContext()
	instance.configureScene(instance.currentScene)
	instance.forceRedraw = true
	instance.RequestRendering()
}

func handleKeyDown(callback *dispatch.Message) {
	tokens := strings.Split(callback.StrData, "\n")
	if len(tokens) != 2 {
		panic("Invalid key down callback Data")
	}
	event := NewAmphionEvent(instance, EventKeyDown, KeyEvent{
		Key:  tokens[0],
		Code: tokens[1],
	})
	instance.updateRoutine.enqueueEventAndRequestUpdate(event)
}

func handleAppHide(_ *dispatch.Message) {
	instance.updateRoutine.enqueueEventAndRequestUpdate(NewAmphionEvent(instance, EventAppHide, nil))
	instance.suspend = true
}

func handleAppShow(_ *dispatch.Message) {
	instance.suspend = false
	instance.updateRoutine.enqueueEventAndRequestUpdate(NewAmphionEvent(instance, EventAppShow, nil))
	instance.RequestRendering()
}

func handleMouseScroll(callback *dispatch.Message) {
	var x, y float32
	n, err := fmt.Sscanf(callback.StrData, "%f:%f", &x, &y)
	if n != 2 || err != nil {
		panic("Invalid scroll callback data")
	}
	instance.RaiseEvent(NewAmphionEvent(instance, EventMouseScroll, a.Vector2{X: x, Y: y}))
}

func handleFrontendReady(_ *dispatch.Message) {
	instance.logger.Info(instance, "Frontend ready")
	instance.startingWg.Done()
}

func newFrontendCallbackHandler() dispatch.MessageDispatcher {
	return &frontendCallbackHandler{
		handlersMap: map[int]dispatch.MessageHandlerFunc{
			frontend.CallbackMouseDown: handleMouseDown,
			frontend.CallbackTouchDown: handleMouseDown,
			frontend.CallbackMouseUp: handleMouseUp,
			frontend.CallbackTouchUp: handleMouseUp,
			frontend.CallbackMouseMove: handleMouseMove,
			frontend.CallbackTouchMove: handleMouseMove,
			frontend.CallbackContextChange: handleContextChange,
			frontend.CallbackKeyDown: handleKeyDown,
			frontend.CallbackAppHide: handleAppHide,
			frontend.CallbackAppShow: handleAppShow,
			frontend.CallbackMouseScroll: handleMouseScroll,
			frontend.CallbackReady: handleFrontendReady,
		},
	}
}