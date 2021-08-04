package engine

import (
	"fmt"
	"github.com/cadmean-ru/amphion/common"
	"github.com/cadmean-ru/amphion/common/a"
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

func (logger *Logger) Info(subsystem interface{}, message string) {
	logger.loggerImpl.Log(Info, logger.subsystemNameOrEmptyString(subsystem), message)
}

func (logger *Logger) Warning(subsystem interface{}, message string) {
	logger.loggerImpl.Log(Warning, logger.subsystemNameOrEmptyString(subsystem), message)
}

func (logger *Logger) Error(subsystem interface{}, message string) {
	logger.loggerImpl.Log(Error, logger.subsystemNameOrEmptyString(subsystem), message)
}

func (logger *Logger) subsystemNameOrEmptyString(subsystem interface{}) string {
	if subsystem == nil {
		return ""
	}

	if comp, ok := subsystem.(Component); ok {
		return NameOfComponent(comp)
	} else if named, ok := subsystem.(a.NamedObject); ok {
		return named.GetName()
	} else {
		return fmt.Sprintf("%v", subsystem)
	}
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