package sqlx

import (
	"reflect"
	"strings"
	"sync"
	"time"
	"utilware/dep/cast"
	"utilware/logger"
)

type Logger struct {
	locker *sync.RWMutex
	output loggerOutput
}

type loggerOutput func(l *QueryLogger)

type QueryLogger struct {
	g *Logger

	Query string
	Args  []interface{}

	Start time.Time
	End   time.Time

	Err error
}

type emptyQueryLogger struct{}

func (l *emptyQueryLogger) done(err error) {}

type loggerDone interface {
	done(err error)
}

func (g *Logger) SetLogger(output loggerOutput) {
	g.locker.Lock()
	defer g.locker.Unlock()

	g.output = output
}

func (g *Logger) append(query string, args ...interface{}) loggerDone {
	if g.output == nil {
		return &emptyQueryLogger{}
	}

	l := &QueryLogger{
		g:     g,
		Query: query,
		Args:  args,
		Start: time.Now(),
	}

	return l
}

func (l *QueryLogger) done(err error) {
	l.End = time.Now()
	l.Err = err

	l.g.locker.RLock()
	defer l.g.locker.RUnlock()
	l.g.output(l)
}

func DefaultLoggerRawOutput(l *QueryLogger) {
	if l.Err != nil {
		logger.Error("[sqlx] [%s] [\"%s\" %+v] => %s", l.End.Sub(l.Start).String(),
			l.Query, l.Args, l.Err.Error())
	} else {
		logger.Debug("[sqlx] [%s] \"%s\" %+v", l.End.Sub(l.Start).String(),
			l.Query, l.Args)
	}
}

func DefaultLoggerBindOutput(l *QueryLogger) {
	if l.Err != nil {
		logger.Error("[sqlx] [%s] %s => %s", l.End.Sub(l.Start).String(),
			bindQuery(l.Query, l.Args...), l.Err.Error())
	} else {
		logger.Debug("[sqlx] [%s] %s", l.End.Sub(l.Start).String(),
			bindQuery(l.Query, l.Args...))
	}
}

// bindQuery 构建日志专用的SQL语句
// @param query string 查询语句 (e.g. SELECT * FROM users WHERE id=?)
// @param args ...interface{} 查询参数 (e.g. 1)
// @tip 需要注意的是，这里的参数是不会被转义的，所以请不要使用这个函数来构建SQL语句, 因为这样会导致SQL注入
func bindQuery(query string, args ...interface{}) string {
	if len(args) == 0 {
		return query
	}

	replaceFirst := func(s string) bool {
		if !strings.Contains(query, "?") {
			return false
		}

		query = strings.Replace(query, "?", s, 1)
		return true
	}

	for _, arg := range args {
		// 反射判断类型并转换然后填充进去
		switch reflect.TypeOf(arg).Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64: // 直接填充
			if !replaceFirst(cast.ToString(arg)) {
				return query
			}
		case reflect.Array, reflect.Slice: // 填充字符串
			val, ok := arg.([]interface{})
			if !ok {
				return query
			}

			l := make([]string, len(val))
			for i := 0; i < len(val); i++ {
				l[i] = cast.ToString(val[i])
			}

			if !replaceFirst(strings.Join(l, ", ")) {
				return query
			}
		default: // 需要添加引号
			if !replaceFirst("'" + cast.ToString(arg) + "'") {
				return query
			}
		}
	}

	return query
}
