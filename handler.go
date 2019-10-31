package web

import (
	"net/http"
)

// HandlerFunc handler func
type HandlerFunc func(*Context)

// Init init
func (f HandlerFunc) Init(ctx *Context) {
}

// Prepare Prepare
func (f HandlerFunc) Prepare(ctx *Context) {
}

// CONNECT CONNECT
func (f HandlerFunc) CONNECT(ctx *Context) {
	f(ctx)
}

// OPTIONS OPTIONS
func (f HandlerFunc) OPTIONS(ctx *Context) {
	f(ctx)
}

// HEAD HEAD
func (f HandlerFunc) HEAD(ctx *Context) {
	f(ctx)
}

// GET GET
func (f HandlerFunc) GET(ctx *Context) {
	f(ctx)
}

// POST POST
func (f HandlerFunc) POST(ctx *Context) {
	f(ctx)
}

// DELETE DELETE
func (f HandlerFunc) DELETE(ctx *Context) {
	f(ctx)
}

// PUT PUT
func (f HandlerFunc) PUT(ctx *Context) {
	f(ctx)
}

// TRACE TRACE
func (f HandlerFunc) TRACE(ctx *Context) {
	f(ctx)
}

// PATCH PATCH
func (f HandlerFunc) PATCH(ctx *Context) {
	f(ctx)
}

// Finish Finish
func (f HandlerFunc) Finish(ctx *Context) {
}

// Handler interface
type Handler interface {
	Init(*Context)
	Prepare(*Context)
	CONNECT(*Context)
	OPTIONS(*Context)
	HEAD(*Context)
	GET(*Context)
	POST(*Context)
	DELETE(*Context)
	PUT(*Context)
	TRACE(*Context)
	PATCH(*Context)
	Finish(*Context)
}

// BaseHandler base httphandler
type BaseHandler struct {
}

// Init init
func (*BaseHandler) Init(*Context) {}

// Prepare prepare
func (*BaseHandler) Prepare(*Context) {}

// CONNECT method
func (handle *BaseHandler) CONNECT(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// OPTIONS method
func (handle *BaseHandler) OPTIONS(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// HEAD method
func (handle *BaseHandler) HEAD(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// GET method
func (handle *BaseHandler) GET(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// POST method
func (handle *BaseHandler) POST(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// DELETE method
func (handle *BaseHandler) DELETE(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// PUT method
func (handle *BaseHandler) PUT(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// TRACE method
func (handle *BaseHandler) TRACE(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// PATCH method
func (handle *BaseHandler) PATCH(ctx *Context) {
	ctx.Error(http.StatusMethodNotAllowed)
}

// Finish finish
func (*BaseHandler) Finish(*Context) {}

// WarpHTTPHandler xx
type warpHandlerFunc http.HandlerFunc

// Init init
func (warpHandlerFunc) Init(*Context) {}

// Prepare prepare
func (warpHandlerFunc) Prepare(*Context) {}

// CONNECT method
func (warp warpHandlerFunc) CONNECT(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// OPTIONS method
func (warp warpHandlerFunc) OPTIONS(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// HEAD method
func (warp warpHandlerFunc) HEAD(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// GET method
func (warp warpHandlerFunc) GET(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// POST method
func (warp warpHandlerFunc) POST(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// DELETE method
func (warp warpHandlerFunc) DELETE(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// PUT method
func (warp warpHandlerFunc) PUT(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// TRACE method
func (warp warpHandlerFunc) TRACE(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// PATCH method
func (warp warpHandlerFunc) PATCH(ctx *Context) {
	warp(ctx.ResponseWriter, ctx.Request)
}

// Finish finish
func (warpHandlerFunc) Finish(*Context) {}
