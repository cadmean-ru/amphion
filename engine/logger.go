package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
)

type LogLevel uint

const (
	Info    LogLevel = 0
	Warning LogLevel = 1
	Error   LogLevel = 2
)

type ILogger interface {
	Log(level LogLevel, subsystem string, message string)
}

type Logger struct {
	loggerImpl ILogger
}

func (logger *Logger) Info(subsystem NamedObject, message string) {
	logger.loggerImpl.Log(Info, logger.subsystemNameOrEmptyString(subsystem), message)
}

func (logger *Logger) Warning(subsystem NamedObject, message string) {
	logger.loggerImpl.Log(Warning, logger.subsystemNameOrEmptyString(subsystem), message)
}

func (logger *Logger) Error(subsystem NamedObject, message string) {
	logger.loggerImpl.Log(Error, logger.subsystemNameOrEmptyString(subsystem), message)
}

func (logger *Logger) subsystemNameOrEmptyString(subsystem NamedObject) string {
	if subsystem == nil {
		return ""
	}

	return subsystem.GetName()
}

func getLoggerImpl(p common.Platform) ILogger {
	switch p.GetName() {
	default:
		return &ConsoleLogger{}
	}
}

func GetLoggerForPlatform(p common.Platform) *Logger {
	return &Logger{
		loggerImpl: getLoggerImpl(p),
	}
}

type ConsoleLogger struct {}

func (logger *ConsoleLogger) Log(level LogLevel, subsystem string, message string) {
	var lvlStr string
	switch level {
	case Info:
		lvlStr = "Info"
		break
	case Warning:
		lvlStr = "Warning"
		break
	case Error:
		lvlStr = "Error"
		break
	}

	fmt.Println(fmt.Sprintf("%s: %s: %s", lvlStr, subsystem, message))
}