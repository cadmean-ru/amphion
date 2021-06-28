package pc

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/frontend"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var keyNames = map[glfw.Key]string {
	glfw.KeyLeftShift: "LeftShift",
	glfw.KeyRightShift: "RightShift",
	glfw.KeyLeftAlt: "LeftAlt",
	glfw.KeyRightAlt: "RightAlt",
	glfw.KeyLeftSuper: "LeftSuper",
	glfw.KeyRightSuper: "RightSuper",
	glfw.KeyLeftControl: "LeftControl",
	glfw.KeyRightControl: "RightControl",
	glfw.KeySpace: "Space",
	glfw.KeyEnter: "Enter",
	glfw.KeyKPEnter: "NumEnter",
	glfw.KeyBackspace: "Backspace",
	glfw.KeyLeft: "LeftArrow",
	glfw.KeyRight: "RightArrow",
	glfw.KeyUp: "UpArrow",
	glfw.KeyDown: "DownArrow",
	glfw.KeyInsert: "Insert",
	glfw.KeyDelete: "Delete",
	glfw.KeyHome: "Home",
	glfw.KeyEnd: "End",
	glfw.KeyPageUp: "PageUp",
	glfw.KeyPageDown: "PageDown",
	glfw.KeyCapsLock: "CapsLock",
	glfw.KeyTab: "Tab",
	glfw.KeyEscape: "Escape",
}

func (f *Frontend) keyCallback(_ *glfw.Window, key glfw.Key, scancode int, action glfw.Action, _ glfw.ModifierKey) {
	keyName := findKeyName(key, scancode)
	data := fmt.Sprintf("%s\n%s", keyName, "")

	fmt.Printf("Key: %v, Scancode: %v, Keyname: %s\n", key, scancode, keyName)

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
		return n
	}
	return glfw.GetKeyName(key, scancode)
}

func (f *Frontend) charCallback(_ *glfw.Window, char rune) {
	fmt.Printf("Char: %s\n", string(char))
	f.disp.SendMessage(dispatch.NewMessageWithStringData(frontend.CallbackRuneInput, string(char)))
}