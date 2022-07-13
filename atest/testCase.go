package atest

type TestCase struct {
	Frontend             TestingFrontend
	TestingDelegate      TestingDelegate
	PrepareSceneDelegate SceneTestingDelegate
	SceneTestingDelegate SceneTestingDelegate
}
