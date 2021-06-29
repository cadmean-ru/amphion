package engine

import "github.com/cadmean-ru/amphion/common/a"

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

func newInputManager() *InputManager {
	return &InputManager{}
}