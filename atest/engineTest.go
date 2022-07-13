package atest

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/dispatch"
	"github.com/cadmean-ru/amphion/engine"
	"github.com/cadmean-ru/amphion/engine/builtin"
	"github.com/cadmean-ru/amphion/frontend"
	"testing"
	"time"
)

var eng *engine.AmphionEngine

type TestingDelegate func(e *engine.AmphionEngine)

type SceneTestingDelegate func(e *engine.AmphionEngine, testScene, testObject *engine.SceneObject)

// RunEngineTest starts the Amphion engine and the calls the specified delegate.
func RunEngineTest(t *testing.T, delegate TestingDelegate) {
	t.Logf("Starting engine test")

	front := NewHeadlessFrontend()
	front.Init()

	eng = engine.Initialize(front)

	front.Run()

	eng.Start()

	delegate(eng)
}

// RunEngineTestWithScene first starts the Amphion engine.
// It then creates the default testing scene (See MakeTestScene).
// The prepareDelegate is called before the scene is shown.
// Here you can perform some setup like adding new objects and components to the scene.
// The testingDelegate is where you actually call the tested code.
func RunEngineTestWithScene(t *testing.T, prepareDelegate, testingDelegate SceneTestingDelegate) {
	RunEngineTest(t, func(e *engine.AmphionEngine) {
		var scene, testObject *engine.SceneObject

		scene, testObject = MakeTestScene(func(e *engine.AmphionEngine) {
			if testingDelegate != nil {
				testingDelegate(e, scene, testObject)
			}
		})

		if prepareDelegate != nil {
			prepareDelegate(e, scene, testObject)
		}

		err := e.ShowScene(scene)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func RunEngineTestCase(t *testing.T, testCase *TestCase) {
	t.Logf("Starting engine test")

	testCase.Frontend.Init()
	eng = engine.Initialize(testCase.Frontend)

	go func() {
		eng.Start()

		if testCase.TestingDelegate != nil {
			testCase.TestingDelegate(eng)
		} else {
			var scene, testObject *engine.SceneObject
			scene, testObject = MakeTestScene(func(e *engine.AmphionEngine) {
				if testCase.SceneTestingDelegate != nil {
					testCase.SceneTestingDelegate(e, scene, testObject)
				}
			})

			if testCase.PrepareSceneDelegate != nil {
				testCase.PrepareSceneDelegate(eng, scene, testObject)
			}

			err := eng.ShowScene(scene)
			if err != nil {
				t.Fatal(err)
			}
		}
	}()

	testCase.Frontend.Run()
}

// MakeTestScene creates the default testing scene, that contains only one child object of size (100; 100; 100)
// located in the center with the TestingComponent attached to it.
func MakeTestScene(delegate TestingDelegate) (*engine.SceneObject, *engine.SceneObject) {
	scene := engine.NewSceneObject("test scene")
	scene.Transform.SetSize(500, 500, 0)
	scene.AddComponent(builtin.NewAbsoluteLayout())

	testObject := engine.NewSceneObject("test object")
	testObject.Transform.SetPositionCentered()
	testObject.Transform.SetSize(100, 100, 100)
	testObject.Transform.SetPivotCentered()
	testObject.AddComponent(NewTestingComponent(delegate))

	scene.AddChild(testObject)

	return scene, testObject
}

// SimulateCallback simulates a frontend callback with the specified code and data.
func SimulateCallback(callback *dispatch.Message) {
	instance.SimulateCallback(callback)
}

// SimulateClick simulates user's click at the specified position on the screen.
func SimulateClick(x, y int, button engine.MouseButton) {
	data := fmt.Sprintf("%d;%d;%d", x, y, button)
	instance.SimulateCallback(dispatch.NewMessageWithStringData(frontend.CallbackMouseDown, data))
	time.Sleep(100)
	instance.SimulateCallback(dispatch.NewMessageWithStringData(frontend.CallbackMouseUp, data))
}

// SimulateClickOnObject simulates user's click in the center of the specified object.
func SimulateClickOnObject(o *engine.SceneObject, button engine.MouseButton) {
	rect := o.Transform.GlobalRect()
	x := int(rect.X.Min + rect.X.GetLength()/2)
	y := int(rect.Y.Min + rect.Y.GetLength()/2)
	SimulateClick(x, y, button)
}

// WaitForStop blocks the calling goroutine until the engine is stopped.
func WaitForStop() {
	eng.WaitForStop()
}

// Stop stops the testing instance of Amphion engine.
func Stop() {
	eng.Stop()
}
