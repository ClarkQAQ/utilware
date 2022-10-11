package logger

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

// 日志等级
type LoggerLevel uint16

const (
	LevelMuted  LoggerLevel = 0           // 关闭输出
	LevelDebug  LoggerLevel = (1 << iota) // 调试信息
	LevelInfo                             // 普通信息
	LevelWarn                             // 警告消息
	LevelError                            // 错误消息
	LevelFatal                            // 致命错误
	LevelPrintf                           // 格式化输出

	defaultCaller int = 2 // 默认追踪调用层级
)

// 日志接口
type LoggerInterface struct {
	asyncFlag bool           // 是否异步
	cache     chan *Logger   // 日志队列
	writers   []LoggerWriter // 日志输出
}

// 日志初始化配置
type LoggerOptions struct {
	CacheSize int          // 日志异步队列缓存大小 (0: 同步输出)
	Writer    LoggerWriter // 默认日志输出接口
}

// 日志内容
type Logger struct {
	Level   LoggerLevel // 等级
	File    string      // 追溯文件
	Line    int         // 追溯行号
	Message string      // 消息
	Attrs   Attrs       // 扩展Attrs
	Time    time.Time   // 时间
}

// 日志输出接口
// @param log 日志内容
type LoggerWriter func(log *Logger)

// 日志参数
type Attrs map[string]interface{}

// 创建日志接口
// @return 日志接口实例
func NewLogger(opts ...*LoggerOptions) *LoggerInterface {
	l := &LoggerInterface{
		asyncFlag: false,
		writers:   []LoggerWriter{},
	}

	if len(opts) < 1 || opts[0] == nil {
		return l
	}

	if opts[0].CacheSize > 0 {
		l.cache = make(chan *Logger, opts[0].CacheSize)
		l.asyncLogger()
	}

	if opts[0].Writer != nil {
		l.Register(opts[0].Writer)
	}

	return l
}

// 运行日志队列监听
func (l *LoggerInterface) asyncLogger() {
	if l.asyncFlag {
		return
	}

	l.asyncFlag = true

	go func(l *LoggerInterface) {
		defer func() {
			l.asyncFlag = false
		}()

		for g := range l.cache {
			l.Writer(g)
		}
	}(l)
}

// 添加日志队列
// @param level 日志等级
// @param skipCaller 跳过调用者数量
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) Logger(level LoggerLevel, skipCaller int, format string, args ...interface{}) {
	g := &Logger{
		Level: level,
		Time:  time.Now(),
	}

	args, g.Attrs = splitAttrs(args)
	_, g.File, g.Line, _ = runtime.Caller(skipCaller)
	g.Message = fmt.Sprintf(fmt.Sprint(format), args...)

	if l.asyncFlag {
		l.cache <- g // 异步消息 写入队列
		return
	}

	l.Writer(g) // 默认同步消息 直接写入
}

// 写入日志到输出
// @param g 日志内容
func (l *LoggerInterface) Writer(g *Logger) {
	for _, writer := range l.writers {
		writer(g)
	}
}

// 关闭日志接口
// @return error 错误信息
func (l *LoggerInterface) Close() error {
	if !func() bool {
		defer func() { recover() }()
		close(l.cache)
		return true
	}() {
		return ErrChannelClosed
	}

	return nil
}

// 注册日志输出
// @param writer 日志输出接口
func (l *LoggerInterface) Register(writer LoggerWriter) {
	l.writers = append(l.writers, writer)
}

// Debug 日志 (调试信息)
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) Debug(format string, args ...interface{}) {
	l.Logger(LevelDebug, defaultCaller, format, args...)
}

// Info 日志 (普通信息)
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) Info(format string, args ...interface{}) {
	l.Logger(LevelInfo, defaultCaller, format, args...)
}

// Warn 日志 (警告信息)
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) Warn(format string, args ...interface{}) {
	l.Logger(LevelWarn, defaultCaller, format, args...)
}

// Error 日志 (错误信息)
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) Error(format string, args ...interface{}) {
	l.Logger(LevelError, defaultCaller, format, args...)
}

// Fatal 日志 (致命错误) (会调用 os.Exit(1) 退出程序)
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) Fatal(format string, args ...interface{}) {
	l.Logger(LevelFatal, defaultCaller, format, args...)
	panic(fmt.Sprintf(format, args...))
}

// Fatal 日志 (异常错误) (会调用 os.Exit(code) 抛出异常)
// @param code 退出码
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) FatalWithExit(code int, format string, args ...interface{}) {
	l.Logger(LevelFatal, defaultCaller, format, args...)
	os.Exit(code)
}

// Printf 输出日志 (普通信息)
// @param format 字符串
// @param args 参数 (最后一个参数可以是 attrs)
func (l *LoggerInterface) Printf(format string, args ...interface{}) {
	l.Logger(LevelPrintf, defaultCaller, format, args...)
}
