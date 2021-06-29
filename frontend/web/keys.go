package web

import (
	"github.com/cadmean-ru/amphion/engine"
	"strings"
)

var jsKeyNames = map[string]engine.KeyName {
	"ShiftLeft": engine.KeyLeftShift,
	"ShiftRight": engine.KeyRightShift,
	"AltLeft": engine.KeyLeftAlt,
	"AltRight": engine.KeyRightAlt,
	"MetaLeft": engine.KeyLeftSuper,
	"MetaRight": engine.KeyRightSuper,
	"ControlLeft": engine.KeyLeftControl,
	"ControlRight": engine.KeyRightControl,
	"Space": engine.KeySpace,
	"Enter": engine.KeyEnter,
	"NumEnter": engine.KeyNumEnter,
	"Backspace": engine.KeyBackspace,
	"ArrowLeft": engine.KeyLeftArrow,
	"ArrowRight": engine.KeyRightArrow,
	"ArrowUp": engine.KeyUpArrow,
	"ArrowDown": engine.KeyDownArrow,
	"Insert": engine.KeyInsert,
	"Delete": engine.KeyDelete,
	"Home": engine.KeyHome,
	"End": engine.KeyEnd,
	"PageUp": engine.KeyPageUp,
	"PageDown": engine.KeyPageDown,
	"CapsLock": engine.KeyCapsLock,
	"Tab": engine.KeyTab,
	"Escape": engine.KeyEscape,
	"Slash": engine.KeyName("/"),
	"Backslash": engine.KeyName("\\"),
	"Comma": engine.KeyName(","),
	"Period": engine.KeyName("."),
	"Quote": engine.KeyName("'"),
	"Semicolon": engine.KeyName(";"),
	"Minus": engine.KeyName("-"),
	"Equal": engine.KeyName("="),
	"BracketLeft": engine.KeyName("["),
	"BracketRight": engine.KeyName("]"),
	"Backquote": engine.KeyName("`"),
	"IntlBackslash": engine.KeyName("§"),
}

func getKeyName(jsCode string) string {
	if n, ok := jsKeyNames[jsCode]; ok {
		return string(n)
	}

	if strings.HasPrefix(jsCode, "Key") {
		return strings.ToLower(jsCode[3:])
	}

	if strings.HasPrefix(jsCode, "Digit") {
		return jsCode[5:]
	}

	return string(engine.KeyUnknown)
}
