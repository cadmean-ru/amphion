package engine

// Prints a message to the log from the current component.
func LogInfo(msg string) {
	instance.logger.Info(instance.currentComponent, msg)
}

// Prints a warning to the log from the current component.
func LogWarning(msg string) {
	instance.logger.Warning(instance.currentComponent, msg)
}

// Prints an error to the log from the current component.
func LogError(msg string) {
	instance.logger.Error(instance.currentComponent, msg)
}

// Same as LogInfo, but prints only if app is in debug mode.
func LogDebug(msg string) {
	if instance.currentApp == nil || instance.currentApp.Debug {
		LogInfo(msg)
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