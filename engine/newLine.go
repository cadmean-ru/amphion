package engine

import "runtime"

//GetNewLineString returns the correct new line sequence for the current platform.
func GetNewLineString() string {
	switch runtime.GOOS {
	case "js":
		if instance.currentApp != nil && instance.currentApp.HostOS == "windows" {
			return "\r\n"
		} else {
			return "\n"
		}
	case "windows":
		return "\r\n"
	default:
		return "\n"
	}
}