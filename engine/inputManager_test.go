package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInputManager_IsKeyPressed(t *testing.T) {
	ass := assert.New(t)

	keyPressed := KeyEnter
	keyNotPressed := KeyName("a")

	m := newInputManager()

	m.reportKeyPressed(keyPressed)

	ass.True(m.IsKeyPressed(keyPressed))
	ass.False(m.IsKeyPressed(keyNotPressed))
}