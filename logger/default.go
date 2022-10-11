package logger

import (
	"fmt"
	"os"
)

var (
	globalWriter = NewDefaultWriter(os.Stdout)
	globalLogger = NewLogger(&LoggerOptions{
		Writer: globalWriter.Write,
	})
)

func GlobalInterface() *LoggerInterface {
	return globalLogger
}

func GlobalWriter() *DefaultWriter {
	return globalWriter
}

func Info(format string, args ...interface{}) {
	globalLogger.Logger(LevelInfo, defaultCaller, format, args...)
}

func Debug(format string, args ...interface{}) {
	globalLogger.Logger(LevelDebug, defaultCaller, format, args...)
}

func Warn(format string, args ...interface{}) {
	globalLogger.Logger(LevelWarn, defaultCaller, format, args...)
}

func Error(format string, args ...interface{}) {
	globalLogger.Logger(LevelError, defaultCaller, format, args...)
}

func Fatal(format string, args ...interface{}) {
	globalLogger.Logger(LevelFatal, defaultCaller, format, args...)
	panic(fmt.Sprintf(format, args...))
}

func FatalWithExit(code int, format string, args ...interface{}) {
	globalLogger.Logger(LevelFatal, defaultCaller, format, args...)
	os.Exit(code)
}

func Printf(format string, args ...interface{}) {
	globalLogger.Logger(LevelPrintf, defaultCaller, format, args...)
}

func Timer() *TimerLogger {
	return globalLogger.Timer()
}

func Progress(length int, total float64, unit string) *ProgressLogger {
	return globalLogger.Progress(length, total, unit)
}

func Register(writer LoggerWriter) {
	globalLogger.Register(writer)
}
