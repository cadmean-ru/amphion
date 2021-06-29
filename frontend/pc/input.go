package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
)

var keyNames = map[glfw.Key]engine.KeyName{
	glfw.KeyLeftShift:    engine.KeyLeftShift,
	glfw.KeyRightShift:   engine.KeyRightShift,
	glfw.KeyLeftAlt:      engine.KeyLeftAlt,
	glfw.KeyRightAlt:     engine.KeyRightAlt,
	glfw.KeyLeftSuper:    engine.KeyLeftSuper,
	glfw.KeyRightSuper:   engine.KeyRightSuper,
	glfw.KeyLeftControl:  engine.KeyLeftControl,
	glfw.KeyRightControl: engine.KeyRightControl,
	glfw.KeySpace:        engine.KeySpace,
	glfw.KeyEnter:        engine.KeyEnter,
	glfw.KeyKPEnter:      engine.KeyNumEnter,
	glfw.KeyBackspace:    engine.KeyBackspace,
	glfw.KeyLeft:         engine.KeyLeftArrow,
	glfw.KeyRight:        engine.KeyRightArrow,
	glfw.KeyUp:           engine.KeyUpArrow,
	glfw.KeyDown:         engine.KeyDownArrow,
	glfw.KeyInsert:       engine.KeyInsert,
	glfw.KeyDelete:       engine.KeyDelete,
	glfw.KeyHome:         engine.KeyHome,
	glfw.KeyEnd:          engine.KeyEnd,
	glfw.KeyPageUp:       engine.KeyPageUp,
	glfw.KeyPageDown:     engine.KeyPageDown,
	glfw.KeyCapsLock:     engine.KeyCapsLock,
	glfw.KeyTab:          engine.KeyTab,
	glfw.KeyEscape:       engine.KeyEscape,
}

func (f *Frontend) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, _ glfw.ModifierKey) {
	keyName := findKeyName(key, scancode)
	data := fmt.Sprintf("%s\n%s", keyName, "")

	var code int
	switch action {
	case glfw.Press:
		code = frontend.CallbackKeyDown
	case glfw.Release:
		code = frontend.CallbackKeyUp
	}

	f.disp.SendMessage(dispatch.NewMessageWithStringData(code, data))
}

func findKeyName(key glfw.Key, scancode int) string {
	if n, ok := keyNames[key]; ok {
		return string(n)
	}
	if runtime.GOOS == "darwin" && scancode == 63 {
		return string(engine.KeyFn)
	}
	return glfw.GetKeyName(key, scancode)
}

func (f *Frontend) charCallback(_ *glfw.Window, char rune) {
	f.disp.SendMessage(dispatch.NewMessageWithStringData(frontend.CallbackTextInput, string(char)))
}
