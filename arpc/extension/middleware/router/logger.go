package router

import (
	"time"

	"utilware/arpc"
	"utilware/arpc/log"
)

// Logger returns the logger middleware.
func Logger() arpc.HandlerFunc {
	return func(ctx *arpc.Context) {
		t := time.Now()

		ctx.Next()

		cmd := ctx.Message.Cmd()
		flag := ctx.Message.Buffer[arpc.HeaderIndexFlag]
		method := ctx.Message.Method()
		addr := ctx.Client.Conn.RemoteAddr()
		cost := time.Since(t).Milliseconds()

		switch cmd {
		case arpc.CmdRequest, arpc.CmdNotify:
			log.Info("[%#x:%#x] [%d] %s %s %dms", cmd, flag, ctx.Message.Seq(), method, addr, cost)
		default:
			log.Error("invalid cmd: %d,\tdropped", cmd)
			ctx.Done()
		}
	}
}
