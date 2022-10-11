package plugin

import (
	"net"
	"strings"
	"utilware/tig"
)

func GetClientIP(c *tig.Context) string {
	xForwardedFor := c.Req.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(c.Req.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Req.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
