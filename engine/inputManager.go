package engine

import (
	"github.com/cadmean-ru/amphion/common/a"
	"runtime"
)

type InputManager struct {
	lastReportedMousePosition a.IntVector2
	pressedKeys []KeyName
}

func (m *InputManager) reportCursorPosition(pos a.IntVector2) {
	m.lastReportedMousePosition = pos
}

func (m *InputManager) reportKeyPressed(name KeyName) {
	for _, n := range m.pressedKeys {
		if n == name {
			return
		}
	}

	m.pressedKeys = append(m.pressedKeys, name)
}

func (m *InputManager) reportKeyReleased(name KeyName) {
	index := -1
	for i, n := range m.pressedKeys {
		if n == name {
			index = i
			break
		}
	}

	if index == -1 {
		return
	}

	if len(m.pressedKeys) > 1 {
		m.pressedKeys[index] = m.pressedKeys[len(m.pressedKeys)-1]
	}
	m.pressedKeys = m.pressedKeys[:len(m.pressedKeys)-1]
}

func (m *InputManager) GetCursorPosition() a.IntVector2 {
	return m.lastReportedMousePosition
}

func (m *InputManager) IsKeyPressed(name KeyName) bool {
	for _, n := range m.pressedKeys {
		if n == name {
			return true
		}
	}
	return false
}

func (m *InputManager) IsMainCombinationKeyPressed() bool {
	if runtime.GOOS == "darwin" {
		return m.IsSuperPressed()
	}
	return m.IsControlPressed()
}

func (m *InputManager) IsControlPressed() bool {
	return m.IsKeyPressed(KeyLeftControl) || m.IsKeyPressed(KeyRightControl)
}

func (m *InputManager) IsShiftPressed() bool {
	return m.IsKeyPressed(KeyLeftShift) || m.IsKeyPressed(KeyRightShift)
}

func (m *InputManager) IsSuperPressed() bool {
	return m.IsKeyPressed(KeyLeftSuper) || m.IsKeyPressed(KeyRightSuper)
}

func newInputManager() *InputManager {
	return &InputManager{}
}