package web

import "strings"

// Options options
type Options struct {
	Address           string `yaml:"address" json:"address,omitempty"`
	CertFile          string
	KeyFile           string
	ReadTimeout       int
	ReadHeaderTimeout int
	WriteTimeout      int
	IdleTimeout       int
	MaxHeaderBytes    int
	StaticPaths       map[string]string //静态文件路径头 strings.Trim(path, staticPath)
}

// Option func
type Option func(*Options)

func newOptions(opts ...Option) Options {
	opt := Options{
		Address:     "127.0.0.1:8080",
		StaticPaths: map[string]string{},
	}

	for _, o := range opts {
		o(&opt)
	}

	return opt
}

// Address address
func Address(addr string) Option {
	return func(o *Options) {
		o.Address = addr
	}
}

// StaticPath StaticPath
func StaticPath(urlpath string, webpath ...string) Option {
	return func(o *Options) {
		if len(webpath) == 0 {
			o.StaticPaths[strings.Trim(urlpath, "^$")] = strings.Trim(urlpath, "^$")
			return
		}
		o.StaticPaths[strings.Trim(urlpath, "^$")] = strings.Trim(webpath[0], "^$")
	}
}
