package builtin

import (
	"github.com/cadmean-ru/amphion/atest"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFocusEvent(t *testing.T) {
	var gained, lost bool

	atest.RunEngineTestWithScene(t, func(e *engine.AmphionEngine, testScene, testObject *engine.SceneObject) {
		testObject.AddComponent(NewShapeView(ShapeRectangle))
		testObject.AddComponent(NewRectBoundary())
		testObject.AddComponent(NewEventListener(engine.EventFocusGain, func(event engine.Event) bool {
			engine.LogInfo("Focus gained")
			gained = true
			return true
		}))
		testObject.AddComponent(NewEventListener(engine.EventFocusLose, func(event engine.Event) bool {
			engine.LogInfo("Focus lost")
			lost = true
			return true
		}))
	}, func(e *engine.AmphionEngine, testScene, testObject *engine.SceneObject) {
		atest.SimulateClickOnObject(testObject, engine.MouseLeft)
		time.Sleep(1 * time.Second)
		atest.SimulateClick(0, 0, engine.MouseLeft)
		atest.Stop()
	})

	atest.WaitForStop()

	assert.True(t, gained)
	assert.True(t, lost)
}

func TestClickEvent(t *testing.T) {
	var clicked bool
	atest.RunEngineTestWithScene(t, func(e *engine.AmphionEngine, testScene, testObject *engine.SceneObject) {
		testObject.AddComponent(NewRectBoundary())
		testObject.AddComponent(NewEventListener(engine.EventMouseDown, func(event engine.Event) bool {
			engine.LogInfo("Clicked")
			clicked = true
			atest.Stop()
			return true
		}))
	}, func(e *engine.AmphionEngine, testScene, testObject *engine.SceneObject) {
		atest.SimulateClickOnObject(testObject, engine.MouseLeft)

	})

	atest.WaitForStop()

	assert.True(t, clicked)
}
