package middleware

import (
	"time"
	"utilware/tig"
)

func Timer(l func(msg string, v ...interface{})) tig.HandlerFunc {
	return func(c *tig.Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()

		l("[%d] [%s] %s %s", c.Status(),
			c.Req.Method, c.Req.RequestURI, time.Since(t))
	}
}
