// +build linux windows darwin
// +build !android
// +build !ios

package atest

import "github.com/cadmean-ru/amphion/engine"

// TestingComponent is a component for running testing code in scene.
// It calls the specified delegate in OnStart method.
type TestingComponent struct {
	engine.ComponentImpl
	testDelegate TestingDelegate
}

func (t *TestingComponent) OnStart() {
	if t.testDelegate == nil {
		return
	}

	t.testDelegate(t.Engine)
}

func (t *TestingComponent) GetName() string {
	return engine.NameOfComponent(t)
}

func NewTestingComponent(testingDelegate TestingDelegate) *TestingComponent {
	return &TestingComponent{
		testDelegate:  testingDelegate,
	}
}