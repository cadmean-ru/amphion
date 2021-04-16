package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
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
func RunTask(task Task) {
	instance.RunTask(task)
}

// BindEventHandler - shortcut for (engine *AmphionEngine) BindEventHandler(code int, handler EventHandler).
// Binds an event handler for the specified event code.
// The handler will be invoked in the event loop goroutine, when the event with the specified code is raised.
func BindEventHandler(eventCode int, handler EventHandler) {
	instance.BindEventHandler(eventCode, handler)
}

// UnbindEventHandler - shortcut for (engine *AmphionEngine) UnbindEventHandler(code int, handler EventHandler).
// Unbinds the event handler for the specified event code.
// The handler will no longer be invoked, when the event with the specified code is raised.
func UnbindEventHandler(eventCode int, handler EventHandler) {
	instance.UnbindEventHandler(eventCode, handler)
}

// ExecuteOnFrontendThread - shortcut for (engine *AmphionEngine) ExecuteOnFrontendThread(action func()).
// Executes the specified action on frontend thread. 
func ExecuteOnFrontendThread(action func()) {
	instance.ExecuteOnFrontendThread(action)
}

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