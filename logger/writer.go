package logger

import (
	"fmt"
	"io"
	"path"
	"time"
)

const (
	DefaultLevel                    = LevelDebug | LevelInfo | LevelWarn | LevelError | LevelFatal | LevelPrintf | LevelAnsiCode // 默认日志等级
	DefaultLevelNoColor             = LevelDebug | LevelInfo | LevelWarn | LevelError | LevelFatal | LevelPrintf                 // 默认日志等级 (无颜色)
	LevelAnsiCode       LoggerLevel = 128                                                                                        // 颜色等级
)

type DefaultWriter struct {
	location *time.Location // 时区
	writer   io.Writer      // 输出
	level    LoggerLevel    // 日志等级
}

// 初始化默认日志输出
// @param w io.Writer 实例
// @return DefaultWriter 实例
func NewDefaultWriter(w io.Writer) *DefaultWriter {
	return &DefaultWriter{
		writer:   w,
		location: time.Local,
		level:    DefaultLevel,
	}
}

// 设置时区
// @param location 时区偏移
// @return DefaultWriter 实例
func (sw *DefaultWriter) Location(location *time.Location) *DefaultWriter {
	if location != nil {
		sw.location = location
	}

	return sw
}

// 设置日志等级
// @param level 日志等级
// @return DefaultWriter 实例
func (sw *DefaultWriter) Level(val LoggerLevel) *DefaultWriter {
	sw.level = val
	return sw
}

// 写入日志
// @param log 日志实例
func (sw *DefaultWriter) Write(log *Logger) {
	if sw.level&log.Level != 0 || LevelFatal&log.Level != 0 {
		fmt.Fprintln(sw.writer, sw.Format(log))
	}
}

// 格式化日志
// @param log 日志实例
// @return string 格式化后的日志
func (sw *DefaultWriter) Format(log *Logger) string {
	if sw.level&LevelAnsiCode != 0 {
		return sw.PrettyFormat(log)
	}

	return sw.PureFormat(log)
}

// 纯文本格式化
// @param log 日志实例
// @return string 格式化后的日志
func (sw *DefaultWriter) PureFormat(log *Logger) string {
	return fmt.Sprintf("%s %s %s %s %s",
		LevelName(log.Level),
		log.Time.In(sw.location).Format("06-01-02 15:04:05.000"),
		fmt.Sprintf("%s:%d", path.Base(log.File), log.Line),
		log.Message,
		sw.PrettyAttrs(log.Attrs))
}

// 带颜色格式化
// @param log 日志实例
// @return string 格式化后的日志
func (sw *DefaultWriter) PrettyFormat(log *Logger) string {
	return fmt.Sprintf("%s %s %s %s %s%s",
		levelPretty(log.Level),
		ANSICode(ANSIcyan, log.Time.In(sw.location).Format("06-01-02 15:04:05.000")),
		ANSICode(ANSImagenta, fmt.Sprintf("%s:%d", path.Base(log.File), log.Line)),
		log.Message,
		sw.PrettyAttrs(log.Attrs), ANSIreset)
}

func levelPretty(level LoggerLevel) string {
	switch level {
	case LevelInfo:
		return ANSICode(ANSIgreen, ANSICode(ANSIbold, LevelName(level)))
	case LevelWarn:
		return ANSICode(ANSIyellow, ANSICode(ANSIbold, LevelName(level)))
	case LevelError:
		return ANSICode(ANSIred, ANSICode(ANSIbold, LevelName(level)))
	case LevelDebug:
		return ANSICode(ANSIblue, ANSICode(ANSIbold, LevelName(level)))
	case LevelFatal:
		return ANSICode(ANSIred, ANSICode(ANSIflash, ANSICode(ANSIbold, LevelName(level))))
	case LevelPrintf:
		return ANSICode(ANSIblue, ANSICode(ANSIbold, LevelName(level)))
	default:
		return ANSICode(ANSIwhite, ANSICode(ANSIbold, LevelName(level)))
	}
}

// 解析Attrs
// @param attrs Attrs 实例
// @return string 解析后的Attrs字符串
func (sw *DefaultWriter) PrettyAttrs(attrs Attrs) (result string) {
	for key, val := range attrs {
		result = fmt.Sprintf("%s %s=%v", result, key, val)
	}

	return
}
