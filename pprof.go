package web

import (
	"os"
	"runtime"
	"runtime/debug"
	"runtime/trace"
)

// StartTrace 运行trace
var startTrace = func(ctx *Context) {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	ctx.Text([]byte("start trace, file=trace.out"))
}

// StopTrace 停止trace
var stopTrace = func(ctx *Context) {
	trace.Stop()
}

// StartGC 手动触发GC
var startGC = func(ctx *Context) {
	runtime.GC()
}

// StopGC stop gc
var stopGC = func(ctx *Context) {
	debug.SetGCPercent(-1)
}
