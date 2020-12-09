package tests

import "testing"

func TestStartSceneTest(t *testing.T) {
	startEngineTest()
	var scene = createTestScene()
	testEngineInstance.ShowScene(scene)
}
