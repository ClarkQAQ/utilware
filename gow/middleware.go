package gow

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var (
	timeZone = time.FixedZone("UTC", int((0 * time.Hour).Seconds()))
)

func TimeZone(t *time.Location) {
	timeZone = t
}

// print stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:]) // skip first 3 caller

	var str strings.Builder
	str.WriteString(message + "\ntraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if fmt.Sprintf("%v", err) == "close" {
					panic(nil)
				}
				log.SetPrefix(time.Now().In(timeZone).Format("[error] 2006-01-02 15:04:05 "))
				log.SetFlags(0)
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()

		c.Next()
	}
}

func Traceback(callback func(string)) HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				if fmt.Sprintf("%v", err) == "close" {
					panic(nil)
				}
				log := time.Now().In(timeZone).Format("[error] 2006-01-02 15:04:05 " + trace(fmt.Sprintf("%s", err)))
				fmt.Println(log)
				if callback != nil {
					callback(log)
				}
			}
		}()

		c.Next()
	}
}

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		
		log.SetPrefix(time.Now().In(timeZone).Format("[web] 2006-01-02 15:04:05 "))
		log.SetFlags(0)
		// Calculate resolution time
		log.Printf("[%s] [%d] %s in %v", c.Req.Method, c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
