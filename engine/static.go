package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common/a"
	"github.com/cadmean-ru/amphion/frontend"
)

// Prints a message to the log from the current component, formatting the msg using fmt.Sprintf.
func LogInfo(msg string, values ...interface{}) {
	instance.logger.Info(instance.currentComponent, fmt.Sprintf(msg, values...))
}

// Prints a warning to the log from the current component, formatting the msg using fmt.Sprintf.
func LogWarning(msg string, values ...interface{}) {
	instance.logger.Warning(instance.currentComponent, fmt.Sprintf(msg, values...))
}

// Prints an error to the log from the current component, formatting the msg using fmt.Sprintf.
func LogError(msg string, values ...interface{}) {
	instance.logger.Error(instance.currentComponent, fmt.Sprintf(msg, values...))
}

// Same as LogInfo, but prints only if app is in debug mode.
func LogDebug(msg string, values ...interface{}) {
	if instance.currentApp == nil || instance.currentApp.Debug {
		LogInfo(msg, values...)
	}
}

// Shortcut for (engine *AmphionEngine) RequestUpdate().
func RequestUpdate() {
	instance.RequestUpdate()
}

// Shortcut for (engine *AmphionEngine) RequestRendering().
func RequestRendering() {
	instance.RequestRendering()
}

// Shortcut for (engine *AmphionEngine) LoadPrefab(resId int) (*SceneObject, error).
func LoadPrefab(resId a.ResId) (*SceneObject, error) {
	return instance.LoadPrefab(resId)
}

// Shortcut for (engine *AmphionEngine) RunTask(task Task).
func RunTask(task Task) {
	instance.RunTask(task)
}

// Shortcut for (engine *AmphionEngine) BindEventHandler(code int, handler EventHandler).
// Binds an event handler for the specified event code.
// The handler will be invoked in the event loop goroutine, when the event with the specified code is raised.
func BindEventHandler(eventCode int, handler EventHandler) {
	instance.BindEventHandler(eventCode, handler)
}

// Shortcut for (engine *AmphionEngine) UnbindEventHandler(code int, handler EventHandler).
// Unbinds the event handler for the specified event code.
// The handler will no longer be invoked, when the event with the specified code is raised.
func UnbindEventHandler(eventCode int, handler EventHandler) {
	instance.UnbindEventHandler(eventCode, handler)
}

// Shortcut for (engine *AmphionEngine) ExecuteOnFrontendThread(action func()).
// Executes the specified action on frontend thread. 
func ExecuteOnFrontendThread(action func()) {
	instance.ExecuteOnFrontendThread(action)
}

func GetFrontendContext() frontend.Context {
	return instance.GetGlobalContext()
}

func GetScreenSize() a.IntVector2 {
	return instance.GetGlobalContext().ScreenInfo.GetSize()
}

func GetScreenSize3() a.IntVector3 {
	s := instance.GetGlobalContext().ScreenInfo.GetSize()
	return a.NewIntVector3(s.X, s.Y, 0)
}