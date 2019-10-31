package web_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/corex-io/web"
	"github.com/corex-io/web/middleware"
)

type MyTest struct {
	web.BaseHandler
}

func (t *MyTest) GET(ctx *web.Context) {
	ctx.JSON([]byte(ctx.Host), 0, nil)
}
func (t *MyTest) POST(ctx *web.Context) {
	ctx.JSON(ctx.Host, 0, nil)
}
func (t *MyTest) PUT(ctx *web.Context) {
	ctx.JSON(map[string]interface{}{"1": "2", "3": time.Now()}, 0, nil)
}
func (t *MyTest) DELETE(ctx *web.Context) {
	ctx.JSON(`{"1": "2", "3": "4"}`, 0, nil)
}

type Upload struct {
	web.BaseHandler
}

func (upload *Upload) POST(ctx *web.Context) {
	fmt.Println(ctx.Request.Header)
	// filepath, cnt, err := ctx.RecvFile("test", ".")
	// fmt.Println(filepath, cnt, err)
	// ctx.Text([]byte(filepath))
	err := ctx.RecvFile2(".", "")
	fmt.Println(err)
}

func TestWeb(t *testing.T) {
	app := web.New(
		web.Address("127.0.0.1:9999"),
		web.StaticPath("/static", "."),
	)
	app.Route("^/mytest$", &MyTest{})
	app.DebugPprof()
	app.Handle("^/filesystem/", http.StripPrefix("/filesystem/", http.FileServer(http.Dir("."))))
	app.HandleFs("^/fs/", ".")
	app.Route("^/upload", &Upload{})
	app.Use(middleware.Trace(), middleware.AccessIP("127.0.0.1/32"))
	app.Init()
	app.Run()
}
