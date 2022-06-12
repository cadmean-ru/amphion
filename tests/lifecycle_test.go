package tests

import (
	"github.com/cadmean-ru/amphion/atest"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type LifecycleTester struct {
	engine.ComponentImpl
	counter int
}

func (s *LifecycleTester) OnInit(ctx engine.InitContext) {
	s.ComponentImpl.OnInit(ctx)
	s.counter++
	engine.LogInfo("Initialized")
}

func (s *LifecycleTester) OnStart() {
	s.counter++
	engine.LogInfo("Started")
}

func (s *LifecycleTester) OnStop() {
	s.counter++
	engine.LogInfo("Stopped")
}

func NewLifecycleTester() *LifecycleTester {
	return &LifecycleTester{}
}

func TestComponentLifecycle(t *testing.T) {
	lifecycleTester := NewLifecycleTester()
	childLifecycleTester := NewLifecycleTester()
	var err error

	atest.RunEngineTest(t, func(e *engine.AmphionEngine) {
		scene := engine.NewSceneObject("lifecycle test scene")
		scene.AddComponent(lifecycleTester)
		err = e.ShowScene(scene)
		if err != nil {
			return
		}

		time.Sleep(3 * time.Second)

		testObj := engine.NewSceneObject("lifecycle test object")
		testObj.AddComponent(childLifecycleTester)
		scene.AddChild(testObj)

		time.Sleep(3 * time.Second)

		testObj.SetEnabled(false)

		time.Sleep(3 * time.Second)

		testObj.SetEnabled(true)

		time.Sleep(3 * time.Second)

		e.Stop()
	})

	atest.WaitForStop()

	assert.Nil(t, err)
	assert.Equal(t, 3, lifecycleTester.counter)
	assert.Equal(t, 5, childLifecycleTester.counter)
}
