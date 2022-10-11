package middleware

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"utilware/tig"
)

func Recovery(errorHandler func(c *tig.Context, trace string)) tig.HandlerFunc {
	if errorHandler == nil {
		errorHandler = func(c *tig.Context, _ string) {
			c.String(http.StatusInternalServerError, "Internal Server Error")
		}
	}

	return func(c *tig.Context) {
		defer func() {
			if r := recover(); r != nil {

				msg := fmt.Sprint(r)

				var pcs [32]uintptr
				n := runtime.Callers(3, pcs[:])

				var str strings.Builder
				str.WriteString(msg + "\ntraceback:")
				for _, pc := range pcs[:n] {
					fn := runtime.FuncForPC(pc)
					file, line := fn.FileLine(pc)
					str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
				}

				errorHandler(c, str.String())
			}

			if c.Index() < -1 {
				panic(nil)
			}
		}()

		c.Next()
	}
}
