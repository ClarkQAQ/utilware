package router

import (
	"utilware/arpc"
	"utilware/arpc/util"
)

// Recover returns the recovery middleware handler.
func Recover() arpc.HandlerFunc {
	return func(ctx *arpc.Context) {
		defer util.Recover()
		ctx.Next()
	}
}
