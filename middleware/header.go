package middleware

import (
	"time"

	"github.com/corex-io/web"
)

// AppendHeader append header
func AppendHeader(key, value string) func(*web.Context) {
	return func(ctx *web.Context) {
		ctx.Request.Header.Add(key, value)
	}
}

// Trace trace
func Trace() func(*web.Context) {
	return func(ctx *web.Context) {
		traceID := ctx.Request.Header.Get("trace-id")
		if traceID == "" {
			traceID = time.Now().String()
		}
		uuid := time.Now().String()
		ctx.Request.Header.Add("trace-id", uuid)
		ctx.ResponseWriter.Header().Add("trace-id", uuid)
	}
}
