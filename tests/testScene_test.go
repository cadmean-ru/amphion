package tests

import "testing"

func TestStartSceneTest(t *testing.T) {
	startEngineTest()
	var scene = scene2(testEngineInstance)
	testEngineInstance.ShowScene(scene)
	testEngineInstance.WaitForStop()
}
