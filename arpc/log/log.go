// Copyright 2020 lesismal. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package log

import "utilware/logger"

var (
	// DefaultLogger is the default logger and is used by arpc
	DefaultLogger Logger = &Log{
		Debugf: logger.Debug,
		Infof:  logger.Info,
		Warnf:  logger.Warn,
		Errorf: logger.Error,
	}

	// LevelDebug is the debug level
	level = LevelDebug
)

const (
	LevelError Level = 1 << iota
	LevelWarn
	LevelInfo
	LevelDebug
)

type Level uint

type Log struct {
	Debugf func(format string, v ...interface{})
	Infof  func(format string, v ...interface{})
	Warnf  func(format string, v ...interface{})
	Errorf func(format string, v ...interface{})
}

func (l *Log) Info(format string, v ...interface{}) {
	l.Infof(format, v...)
}

func (l *Log) Debug(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

func (l *Log) Warn(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

func (l *Log) Error(format string, v ...interface{}) {
	l.Errorf(format, v...)
}

// Logger defines log interface
type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
}

// SetLogger sets default logger.
func SetLogger(l Logger) {
	DefaultLogger = l
}

func SetLevel(val Level) {
	level = val
}

// Debug uses DefaultLogger to log a message at LevelDebug.
func Debug(format string, v ...interface{}) {
	if DefaultLogger != nil && level&LevelDebug != 0 {
		DefaultLogger.Debug(format, v...)
	}
}

// Info uses DefaultLogger to log a message at LevelInfo.
func Info(format string, v ...interface{}) {
	if DefaultLogger != nil && level&LevelInfo != 0 {
		DefaultLogger.Info(format, v...)
	}
}

// Warn uses DefaultLogger to log a message at LevelWarn.
func Warn(format string, v ...interface{}) {
	if DefaultLogger != nil && level&LevelWarn != 0 {
		DefaultLogger.Warn(format, v...)
	}
}

// Error uses DefaultLogger to log a message at LevelError.
func Error(format string, v ...interface{}) {
	if DefaultLogger != nil && level&LevelError != 0 {
		DefaultLogger.Error(format, v...)
	}
}
