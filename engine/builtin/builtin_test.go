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
		testObject.AddComponent(NewEventListener(engine.EventFocusGain, func(event engine.AmphionEvent) bool {
			engine.LogInfo("Focus gained")
			gained = true
			return true
		}))
		testObject.AddComponent(NewEventListener(engine.EventFocusLoose, func(event engine.AmphionEvent) bool {
			engine.LogInfo("Focus lost")
			lost = true
			return true
		}))
	}, func(e *engine.AmphionEngine, testScene, testObject *engine.SceneObject) {
		atest.SimulateClickOnObject(testObject)
		time.Sleep(1 * time.Second)
		atest.SimulateClick(0, 0)
		atest.Stop()
	})

	atest.WaitForStop()

	assert.True(t, gained)
	assert.True(t, lost)
}
