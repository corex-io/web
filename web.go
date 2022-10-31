package web

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/pprof"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/corex-io/log"
	"golang.org/x/net/http2"
)

// Web service
type Web struct {
	opts Options
	Log  log.Logger
	Mids []func(*Context)
	Mux  Multiplexer
	*http.Server
	sync.Pool
}

// New new service
func New(opts ...Option) *Web {
	options := newOptions(opts...)
	web := Web{
		opts: options,
		Log:  log.DefaultStdLog(),
		Mux:  NewMultiplexer(),
	}
	return &web
}

// Use append middleware
func (s *Web) Use(f ...func(*Context)) {
	s.Mids = append(s.Mids, f...)
}

// Init initialises options.
func (s *Web) Init(opts ...Option) {
	// process options
	for _, o := range opts {
		o(&s.opts)
	}
}

// SetLog set log
func (s *Web) SetLog(log log.Logger) {
	s.Log = log
}

// Load load config
func (s *Web) Load(v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &s.opts)
}

// Run run
func (s *Web) Run(ctx context.Context) error {
	s.Server = &http.Server{
		Addr:    s.opts.Address,
		Handler: s,
		// TLSConfig:         nil,
		ReadTimeout:       time.Duration(s.opts.ReadTimeout) * time.Second,
		ReadHeaderTimeout: time.Duration(s.opts.ReadHeaderTimeout) * time.Second,
		WriteTimeout:      time.Duration(s.opts.WriteTimeout) * time.Second,
		IdleTimeout:       time.Duration(s.opts.IdleTimeout) * time.Second,
		MaxHeaderBytes:    s.opts.MaxHeaderBytes,
	}

	serverhttp2 := &http2.Server{
		IdleTimeout: time.Duration(s.opts.IdleTimeout) * time.Second,
	}
	if err := http2.ConfigureServer(s.Server, serverhttp2); err != nil {
		s.Log.Errorf("%v", err)
	}
	s.Pool.New = func() interface{} {
		return &Context{}
	}
	s.Log.Debugf("http serve [%s]...", s.opts.Address)

	go func(ctx context.Context) {
		select {
		case <-ctx.Done():
			_ = s.Shutdown(ctx)
		}
	}(ctx)

	if s.opts.CertFile != "" && s.opts.KeyFile != "" {
		return s.Server.ListenAndServeTLS(s.opts.CertFile, s.opts.KeyFile)
	}
	return s.Server.ListenAndServe()
}

// Close close
func (s *Web) Close() error {
	return s.Shutdown(context.Background())
}

func (s *Web) String() string {
	return "web-httpd"
}

// Handle http handle
func (s *Web) Handle(path string, handler http.Handler) {
	s.Mux.Handle(path, handler)
}

// HandleFunc http Handle func
func (s *Web) HandleFunc(path string, f http.HandlerFunc) {
	s.Mux.HandleFunc(path, f)
}

// HandleFs filesystem
func (s *Web) HandleFs(srtipPath, path string) {
	s.Handle(srtipPath, http.StripPrefix(strings.Trim(srtipPath, "^$"), http.FileServer(http.Dir(path))))
}

// Route handle
func (s *Web) Route(path string, handle Handler) {
	s.Mux.Route(path, handle)
}

// RouteFunc route handlerfunc
func (s *Web) RouteFunc(path string, f HandlerFunc) {
	s.Mux.RouteFunc(path, f)
}

func (s *Web) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	ctx := &Context{
		Request:        req,
		ResponseWriter: resp,
		Timestamp:      time.Now(),
		Logger:         s.Log,
	}

	defer func(ctx *Context) {
		if err := recover(); err != nil {
			ctx.Error(http.StatusInternalServerError)
			s.Log.Errorf("%v, %v", string(debug.Stack()), err)
		}
		if ctx.statusCode == 0 {
			ctx.statusCode = http.StatusOK
		}
		s.Log.Infof("%s %d %s (%s) %s", req.Method, ctx.statusCode, ctx.URL.String(), ctx.Remote(), time.Since(ctx.Timestamp))
	}(ctx)

	if err := ctx.ParseForm(); err != nil {
		s.Log.Errorf("parse form fail: %v", err)
		return
	}

	for _, mid := range s.Mids {
		mid(ctx)
		if ctx.IsFinish() {
			return
		}
	}
	// handler static file
	for staticPath, webPath := range s.opts.StaticPaths {
		if filepath.HasPrefix(ctx.URL.Path, staticPath) {
			file := filepath.Join(webPath, strings.TrimPrefix(ctx.URL.Path, staticPath))
			if _, err := os.Stat(file); os.IsNotExist(err) {
				ctx.statusCode = http.StatusNotFound
			}
			http.ServeFile(ctx.ResponseWriter, ctx.Request, file)
			return
		}
	}

	entry := s.Mux.FindRoute(ctx.URL.Path)
	if entry == nil {
		ctx.Error(http.StatusNotFound)
		return
	}

	entry.MyInterface.Init(ctx)
	if ctx.IsFinish() {
		return
	}

	entry.MyInterface.Prepare(ctx)
	if ctx.IsFinish() {
		return
	}
	switch ctx.Method {
	case http.MethodConnect:
		entry.MyInterface.CONNECT(ctx)
	case http.MethodOptions:
		entry.MyInterface.OPTIONS(ctx)
	case http.MethodGet:
		entry.MyInterface.GET(ctx)
	case http.MethodHead:
		entry.MyInterface.HEAD(ctx)
	case http.MethodPost:
		entry.MyInterface.POST(ctx)
	case http.MethodPut:
		entry.MyInterface.PUT(ctx)
	case http.MethodDelete:
		entry.MyInterface.DELETE(ctx)
	case http.MethodPatch:
		entry.MyInterface.PATCH(ctx)
	case http.MethodTrace:
		entry.MyInterface.TRACE(ctx)
	}
	if ctx.IsFinish() {
		return
	}
	entry.MyInterface.Finish(ctx)
}

// DebugPprof debugPprof
func (s *Web) DebugPprof() {
	s.HandleFunc("^/_routeList$", func(resp http.ResponseWriter, req *http.Request) {
		resp.Write([]byte(s.Mux.String()))
	})
	s.HandleFunc("^/debug/pprof/$", pprof.Index)
	s.HandleFunc("^/debug/pprof/allocs$", pprof.Index)
	s.HandleFunc("^/debug/pprof/block$", pprof.Index)
	s.HandleFunc("^/debug/pprof/cmdline$", pprof.Cmdline)
	s.HandleFunc("^/debug/pprof/goroutine$", pprof.Index)
	s.HandleFunc("^/debug/pprof/heap$", pprof.Index)
	s.HandleFunc("^/debug/pprof/mutex$", pprof.Index)
	s.HandleFunc("^/debug/pprof/profile$", pprof.Profile)
	s.HandleFunc("^/debug/pprof/threadcreate$", pprof.Profile)
	s.HandleFunc("^/debug/pprof/symbol$", pprof.Symbol)
	s.HandleFunc("^/debug/pprof/trace$", pprof.Trace)
	s.RouteFunc("^/debug/trace/start$", startTrace)
	s.RouteFunc("^/debug/trace/stop$", stopTrace)
	s.RouteFunc("^/debug/gc/start$", startGC)
	s.RouteFunc("^/debug/gc/stop$", stopGC)
}
