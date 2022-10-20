package logger

import (
	"errors"
	"fmt"
)

var (
	ErrChannelClosed = errors.New("maybe channel is closed")
)

// 终端控制符
const (
	ANSIwhite   = "\033[37m"
	ANSIreset   = "\033[0m"
	ANSIbold    = "\033[1m"
	ANSIred     = "\033[31m"
	ANSIblue    = "\033[34m"
	ANSIgreen   = "\033[32m"
	ANSIcyan    = "\033[36m"
	ANSIyellow  = "\033[33m"
	ANSImagenta = "\033[35m"
	ANSIflash   = "\033[5m"
	ANSIarrow   = "\033[%dA"
	ANSIkclear  = "\033[K"
)

// 包装ANSI控制符
// @param ansiCode ANSI控制符
// @param val 格式化参数
func ANSICode(ansiCcode string, val string) string {
	return ansiCcode + val + ANSIreset
}

// 包装ANSI控制符
// @param ansiCode ANSI控制符
// @param val 格式化参数
func ANSICodef(ansiCcode string, format string, val ...interface{}) string {
	return ansiCcode + fmt.Sprintf(format, val...) + ANSIreset
}

// 等级名称
// @param level 等级
func LevelName(level LoggerLevel) string {
	switch level {
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelDebug:
		return "DEBUG"
	case LevelFatal:
		return "FATAL"
	case LevelPrintf:
		return "PRINTF"
	default:
		return "UNKNOWN"
	}
}

// 分离Attrs参数
func splitAttrs(v []interface{}) ([]interface{}, Attrs) {
	valLen := len(v)

	if valLen == 0 {
		return v, nil
	}

	attrs, ok := v[valLen-1].(Attrs)
	if !ok {
		return v, nil
	}

	return v[:valLen-1], attrs
}
