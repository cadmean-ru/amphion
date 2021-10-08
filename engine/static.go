package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
	"regexp"
	"strings"
)

// LogInfo prints a message to the log from the current component, formatting the msg using fmt.Sprintf.
func LogInfo(msg string, values ...interface{}) {
	instance.logger.Info(NameOfComponent(instance.currentComponent), fmt.Sprintf(msg, values...))
}

// LogWarning prints a warning to the log from the current component, formatting the msg using fmt.Sprintf.
func LogWarning(msg string, values ...interface{}) {
	instance.logger.Warning(instance.currentComponent, fmt.Sprintf(msg, values...))
}

// LogError prints an error to the log from the current component, formatting the msg using fmt.Sprintf.
func LogError(msg string, values ...interface{}) {
	instance.logger.Error(instance.currentComponent, fmt.Sprintf(msg, values...))
}

// LogDebug is same as LogInfo, but prints only if app is in debug mode.
func LogDebug(msg string, values ...interface{}) {
	if instance.currentApp == nil || instance.currentApp.Debug {
		LogInfo(msg, values...)
	}
}

// RequestUpdate - shortcut for (engine *AmphionEngine) RequestUpdate().
func RequestUpdate() {
	instance.RequestUpdate()
}

// RequestRendering - shortcut for (engine *AmphionEngine) RequestRendering().
func RequestRendering() {
	instance.RequestRendering()
}

// LoadPrefab - shortcut for (engine *AmphionEngine) LoadPrefab(resId int) (*SceneObject, error).
func LoadPrefab(resId a.ResId) (*SceneObject, error) {
	return instance.LoadPrefab(resId)
}

// RunTask - shortcut for (engine *AmphionEngine) RunTask(task Task).
func RunTask(task *Task) {
	instance.RunTask(task)
}

// BindEventHandler - shortcut for (engine *AmphionEngine) BindEventHandler(code int, handler EventHandler).
// Binds an event handler for the specified event code.
// The handler will be invoked in the event Loop goroutine, when the event with the specified code is raised.
func BindEventHandler(eventCode int, handler EventHandler) {
	instance.BindEventHandler(eventCode, handler)
}

// UnbindEventHandler - shortcut for (engine *AmphionEngine) UnbindEventHandler(code int, handler EventHandler).
// Unbinds the event handler for the specified event code.
// The handler will no longer be invoked, when the event with the specified code is raised.
func UnbindEventHandler(eventCode int, handler EventHandler) {
	instance.UnbindEventHandler(eventCode, handler)
}

//RaiseEvent raises a new event to be processed in the event goroutine.
func RaiseEvent(event Event) {
	instance.RaiseEvent(event)
}

// ExecuteOnFrontendThread - shortcut for (engine *AmphionEngine) ExecuteOnFrontendThread(action func()).
// Executes the specified action on frontend thread. 
func ExecuteOnFrontendThread(action func()) {
	instance.ExecuteOnFrontendThread(action)
}

//GetFrontendContext returns the global application context.
//See frontend.Context.
func GetFrontendContext() frontend.Context {
	return instance.GetGlobalContext()
}

//GetScreenSize returns the screen size as a.IntVector2.
//X and Y are the width and height of the screen.
func GetScreenSize() a.IntVector2 {
	return instance.GetGlobalContext().ScreenInfo.GetSize()
}

//GetScreenSize3 returns the screen size as a.IntVector3.
//X and Y are the width and height of the screen. Z = 1.
func GetScreenSize3() a.IntVector3 {
	s := instance.GetGlobalContext().ScreenInfo.GetSize()
	return a.NewIntVector3(s.X, s.Y, 1)
}

//GetFeaturesManager returns the current FeaturesManager.
func GetFeaturesManager() *FeaturesManager {
	return instance.GetFeaturesManager()
}

//GetResourceManager returns the current resource manager.
func GetResourceManager() frontend.ResourceManager {
	return instance.GetResourceManager()
}

//GetCurrentScene returns the currently displaying SceneObject object.
func GetCurrentScene() *SceneObject {
	return instance.GetCurrentScene()
}

//GetSceneContext returns the current SceneObject's context.
func GetSceneContext() *SceneContext {
	return instance.GetSceneContext()
}

//GetCurrentApp returns the current loaded app or nil if no app is loaded.
func GetCurrentApp() *frontend.App {
	return instance.GetCurrentApp()
}

//GetAppContext returns the current app's context.
func GetAppContext() *AppContext {
	return instance.GetAppContext()
}

//GetInputManager returns the current input manager.
func GetInputManager() *InputManager {
	return instance.GetInputManager()
}

// SetWindowTitle updates app's window title.
// On web sets the tab's title.
func SetWindowTitle(title string) {
	instance.SetWindowTitle(title)
}

//FindObjectByName searches for an object with the specified name through all the current SceneObject object tree.
//See SceneObject.FindObjectByName.
func FindObjectByName(name string, includeDirty ...bool) *SceneObject {
	if instance.currentScene == nil {
		return nil
	}
	return instance.currentScene.FindObjectByName(name, includeDirty...)
}

//FindComponentByName searches for a component with the specified name through all the current SceneObject object tree.
//See SceneObject.FindComponentByName.
func FindComponentByName(name string, includeDirty ...bool) Component {
	if instance.currentScene == nil {
		return nil
	}
	return instance.currentScene.FindComponentByName(name, includeDirty...)
}

//ForceAllViewsRedraw will request all view in the SceneObject to redraw on the next rendering cycle.
//It will not request rendering, you will need to call RequestRendering after that.
func ForceAllViewsRedraw() {
	instance.forceRedraw = true
}

//IsForcedToRedraw checks if all views redraw was requested in the next rendering cycle by calling ForceAllViewsRedraw.
func IsForcedToRedraw() bool {
	return instance.forceRedraw
}

//ComponentNameMatches checks if namePattern matches the given componentName.
func ComponentNameMatches(componentName, namePattern string) bool {
	tokens := strings.Split(componentName, ".")
	if len(tokens) < 2 {
		return false
	}

	shortName := tokens[len(tokens)-1] //The name after .
	if namePattern == shortName {
		return true
	}

	if matched, err := regexp.MatchString(namePattern, componentName); matched && err == nil {
		return true
	}

	return false
}

//IsInDebugMode checks if engine is currently in debug mode.
func IsInDebugMode() bool {
	return instance.currentApp != nil && instance.currentApp.Debug
}

// CloseScene closes the currently showing SceneObject asynchronously.
// It will call the provided callback function as soon as the SceneObject was closed.
// If no SceneObject is showing calls the callback function immediately.
func CloseScene(closeCallback func()) {
	instance.CloseScene(closeCallback)
}

// ShowScene shows the specified SceneObject object.
// Returns an error, if the engine is not yet ready or if another SceneObject is already showing.
func ShowScene(scene *SceneObject) error {
	return instance.ShowScene(scene)
}

// LoadScene loads SceneObject from a resource file asynchronously.
// If show is true, after loading also shows this SceneObject.
func LoadScene(sceneId a.ResId, show bool) {
	instance.LoadScene(sceneId, show)
}

// SwapScenes hides the current showing SceneObject and shows the currently loaded SceneObject (using LoadScene).
// The previously showing SceneObject will be properly stopped and the deleted.
// So, calling SwapScenes again will not swap the two scenes back.
// If no SceneObject is loaded, will not do anything.
func SwapScenes() {
	instance.SwapScenes()
}

//GetComponentsManager returns the current ComponentsManager.
func GetComponentsManager() *ComponentsManager {
	return instance.componentsManager
}

// NameOfComponent return the name of the given component suitable for serialization.
func NameOfComponent(component interface{}) string {
	return instance.GetComponentsManager().NameOfComponent(component)
}

//GetMessageDispatcher returns the MessageDispatcher for the current scene.
func GetMessageDispatcher() *MessageDispatcher {
	return instance.GetMessageDispatcher()
}