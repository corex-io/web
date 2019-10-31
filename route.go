package web

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

// Entry route entry point
type Entry struct {
	regex       *regexp.Regexp
	MyInterface Handler
}

// Multiplexer mux
type Multiplexer []*Entry

// NewMultiplexer new mux
func NewMultiplexer() Multiplexer {
	var mux []*Entry
	return mux
}

// Handle std http handle
func (mux *Multiplexer) Handle(path string, handler http.Handler) {
	warp := warpHandlerFunc(handler.ServeHTTP)
	mux.Route(path, warp)
}

// HandleFunc handlefunc
func (mux *Multiplexer) HandleFunc(path string, f http.HandlerFunc) {
	warp := warpHandlerFunc(f)
	mux.Route(path, warp)
}

// Route handle
func (mux *Multiplexer) Route(path string, handler Handler) {
	entry := Entry{
		regex:       regexp.MustCompile(path),
		MyInterface: handler,
	}
	*mux = append(*mux, &entry)
}

// RouteFunc route handlerFunc
func (mux *Multiplexer) RouteFunc(path string, f HandlerFunc) {
	mux.Route(path, f)
}

// FindRoute find router
func (mux Multiplexer) FindRoute(path string) *Entry {
	for _, m := range mux {
		if matchs := m.regex.FindStringSubmatch(path); matchs != nil {
			match := make(map[string]string, len(matchs))
			for idx, value := range m.regex.SubexpNames() {
				match[value] = matchs[idx]
			}
			delete(match, "")
			return m
		}
	}
	return nil
}
func (mux Multiplexer) String() string {
	var res []string
	for _, m := range mux {
		res = append(res, fmt.Sprintf("%s", m.regex))
	}
	return strings.Join(res, "\n")
}
