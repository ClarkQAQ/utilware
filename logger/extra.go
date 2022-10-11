package logger

import (
	"fmt"
	"time"
)

const (
	maxProgressLeght int = 100 // 最大进度条长度
	secondFloat64        = float64(time.Second)
)

type TimerLogger struct {
	*LoggerInterface           // 日志接口
	t                time.Time // 开始时间 (当前时间)
}

// 时间记录
// @return 记时器
func (l *LoggerInterface) Timer() *TimerLogger {
	return &TimerLogger{
		LoggerInterface: l,
		t:               time.Now(),
	}
}

// 时间记录结束
// @param format 格式化字符串
// @param args 参数
func (t *TimerLogger) End(format string, args ...interface{}) {
	t.Logger(LevelPrintf, defaultCaller, "(TIME:%s) %s", time.Since(t.t), fmt.Sprintf(format, args...))
}

type ProgressLogger struct {
	*LoggerInterface

	startTime time.Time // 开始时间
	length    int       // 进度条长度
	total     float64   // 任务总数
	current   float64   // 当前任务数
	unit      string    // 单位
}

// 进度条
// @description 用于显示任务进度, 比如: 上传文件, 下载文件, 复制文件等
// @param length 进度条长度
// @param total 任务总数
// @param unit 进度单位
// @return 进度条
func (l *LoggerInterface) Progress(length int, total float64, unit string) *ProgressLogger {
	if length > maxProgressLeght || length < 0 {
		length = maxProgressLeght
	}

	return &ProgressLogger{
		LoggerInterface: l,
		startTime:       time.Now(),
		length:          length,
		total:           total,
		unit:            unit,
	}
}

// 进度条更新
// @description 用于更新进度条, 在输出前使用 fmt.Printf("\033[1A\033[K") 即可伪单行控制台输出
// @param append 追加任务数
// @param message 进度条信息
func (p *ProgressLogger) Append(append float64, message string) {
	p.Set(p.current+append, message)
}

// 进度条设置
// @description 用于更新进度条, 在输出前使用 fmt.Printf("\033[1A\033[K") 即可伪单行控制台输出
// @param current 当前任务数
// @param message 进度条信息
func (p *ProgressLogger) Set(current float64, message string) {
	p.current = current

	if p.current > p.total {
		p.current = p.total
	}

	percent := p.Percent()
	totalTime := time.Since(p.startTime)

	if p.length < 1 {
		return
	}

	pg := ""
	for i := 0; i < int(percent)/(maxProgressLeght/p.length); i++ {
		pg += "#"
	}

	speed := secondFloat64 / float64(totalTime/time.Duration(p.total))

	p.Logger(LevelPrintf, defaultCaller, "[%-"+fmt.Sprint(p.length)+"s] [%.2f%%] [%.2f%s/%.2f%s - %.2f%s/s] %s",
		pg, percent, p.current, p.unit, p.total, p.unit, speed, p.unit, message)

	if percent >= 100 {
		p.Logger(LevelPrintf, defaultCaller, "[DONE] [TOTAL: %.2f TIME: %s SPEED: %.2f%s/s] %s",
			p.total, time.Duration(totalTime), speed, p.unit, message)
	}
}

// 获取进度百分比
// @return int 百分比 (0-100)
func (log *ProgressLogger) Percent() float64 {
	return log.current / log.total * 100
}

// 获取进度条长度/任务数量
// @return int 进度条长度/任务数量
func (log *ProgressLogger) Current() float64 {
	return log.current
}
