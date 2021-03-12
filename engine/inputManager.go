package engine

import "github.com/cadmean-ru/amphion/common/a"

type InputManager struct {
	lastReportedMousePosition a.IntVector2
}

func (m *InputManager) reportCursorPosition(pos a.IntVector2) {
	m.lastReportedMousePosition = pos
}

func (m *InputManager) GetCursorPosition() a.IntVector2 {
	return m.lastReportedMousePosition
}

func newInputManager() *InputManager {
	return &InputManager{}
}