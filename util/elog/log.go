package elog

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"strconv"
	"time"
)

type logType struct {
	debug    bool
	timeZone *time.Location
	stack    bool
	logEvent map[string](func(string, string))
}

func New(debug bool, timeZone int64, stack bool, defaultLogEvent func(tp string, log string)) *logType {
	timeLocation := time.FixedZone("time", int((time.Duration(timeZone) * time.Hour).Seconds()))
	if defaultLogEvent != nil {
		return &logType{debug, timeLocation, stack, map[string](func(string, string)){"default": defaultLogEvent}}
	}
	return &logType{debug, timeLocation, stack, make(map[string](func(string, string)))}
}

func callerName() string {
	pc, _, _, _ := runtime.Caller(3)
	return runtime.FuncForPC(pc).Name()
}

func (lgt *logType) StringTimeZone(str string) error {
	if data := strings.Split(str, ":"); len(data) == 2 {
		if n, err := strconv.Atoi(data[1]); err == nil {
			lgt.timeZone = time.FixedZone(data[0], int((time.Duration(n) * time.Hour).Seconds()))
			return nil
		}
	}
	return errors.New("split error")
}

func (lgt *logType) FixedZone(x string, i int) {
	lgt.timeZone = time.FixedZone(x, int((time.Duration(i) * time.Hour).Seconds()))
}

func (lgt *logType) TimeZone(t *time.Location) {
	lgt.timeZone = t
}

func (lgt *logType) LogEventSet(name string, fc func(log_type string, log string)) {
	lgt.logEvent[name] = fc
}

func (lgt *logType) LogEventDel(name string) error {
	if _, has := lgt.logEvent[name]; !has {
		return errors.New("event does not exist")
	}
	delete(lgt.logEvent, name)
	return nil
}

func (lgt *logType) logEventCall(log_type string, log_value string) {
	for _, fc := range lgt.logEvent {
		fc(log_type, log_value)
	}
}

func (lgt *logType) logPrefix(tag string) string {
	return time.Now().In(lgt.timeZone).Format("[" + tag + "] 2006-01-02 15:04:05 ")
}

func (lgt *logType) logString(list []interface{}) (value string) {
	if lgt.stack == true {
		value += fmt.Sprintf(" [$%v] ", callerName())
	}

	for _, v := range list {
		value += fmt.Sprintf("%v", v)
	}

	return value
}

func (lgt *logType) Info(list ...interface{}) {
	var s string
	s += lgt.logPrefix("info")
	s += lgt.logString(list)
	fmt.Println(s)
	lgt.logEventCall("info", s)
}
func (lgt *logType) Warning(list ...interface{}) {
	var s string
	s += lgt.logPrefix("warning")
	s += lgt.logString(list)
	fmt.Println(s)
	lgt.logEventCall("warning", s)
}
func (lgt *logType) Debug(list ...interface{}) {
	var s string
	s += lgt.logPrefix("debug")
	s += lgt.logString(list)
	fmt.Println(s)
	lgt.logEventCall("debug", s)
}
func (lgt *logType) Error(list ...interface{}) {
	var s string
	s += lgt.logPrefix("error")
	s += lgt.logString(list)
	fmt.Println(s)
	lgt.logEventCall("error", s)
}
