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
	instance.logger.Info(instance.currentComponent, fmt.Sprintf(msg, values...))
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
func RaiseEvent(event AmphionEvent) {
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

//GetCurrentScene returns the currently displaying scene object.
func GetCurrentScene() *SceneObject {
	return instance.GetCurrentScene()
}

//GetSceneContext returns the current scene's context.
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

//FindObjectByName searches for an object with the specified name through all the current scene object tree.
//See SceneObject.FindObjectByName.
func FindObjectByName(name string, includeDirty ...bool) *SceneObject {
	if instance.currentScene == nil {
		return nil
	}
	return instance.currentScene.FindObjectByName(name, includeDirty...)
}

//FindComponentByName searches for a component with the specified name through all the current scene object tree.
//See SceneObject.FindComponentByName.
func FindComponentByName(name string, includeDirty ...bool) Component {
	if instance.currentScene == nil {
		return nil
	}
	return instance.currentScene.FindComponentByName(name, includeDirty...)
}

//ForceAllViewsRedraw will request all view in the scene to redraw on the next rendering cycle.
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

func CloseScene(closeCallback func()) {
	instance.CloseScene(closeCallback)
}

func ShowScene(scene *SceneObject) error {
	return instance.ShowScene(scene)
}

func LoadScene(sceneId a.ResId, show bool) {
	instance.LoadScene(sceneId, show)
}

func SwapScenes() {
	instance.SwapScenes()
}