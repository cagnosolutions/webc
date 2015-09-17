package web

import (
	"net/http"
	"strings"
)

type controller func(http.ResponseWriter, *http.Request, *Context)

type route struct {
	method string
	path   []string
	handle controller
}

var defaultMux = instance()

func Get(path string, handler controller) {
	defaultMux.handle("GET", path, handler)
}

func Put(path string, handler controller) {
	defaultMux.handle("PUT", path, handler)
}

func Post(path string, handler controller) {
	defaultMux.handle("Post", path, handler)
}

func Delete(path string, handler controller) {
	defaultMux.handle("DELETE", path, handler)
}

func Serve(host string) {
	defaultMux.ctx.gc()
	if err := http.ListenAndServe(host, defaultMux); err != nil {
		panic(err)
	}
}

type mux struct {
	routes []*route
	ctx    *contextStore
}

func instance() *mux {
	return &mux{
		routes: make([]*route, 0),
		ctx:    &contextStore{contexts: make(map[string]*Context, 0)},
	}
}

func (m *mux) handle(method, path string, handler controller) {
	m.routes = append(m.routes, &route{method, SliceString(path, '/'), handler})
}

func (m *mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// for now: ignore options and favicon
	if r.Method == "OPTIONS" || r.URL.Path == "/favicon.ico" {
		return
	}

	// handle static mapping
	if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/static/") {
		http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
		return
	}

	for _, route := range m.routes {
		if route.method == r.Method {
			path := SliceString(r.URL.Path, '/')
			if pathVars, ok := match(path, route.path); ok {
				ctx := m.ctx.get(w, r)
				ctx.SetPathVars(pathVars)
				route.handle(w, r, ctx)
				return
			}
		}
	}
	return
}

func match(req []string, pat []string) (map[string]string, bool) {
	v := make(map[string]string)
	if len(req) == len(pat) {
		for i := 0; i < len(pat); i++ {
			if req[i] != pat[i] {
				if pat[i][0] == ':' {
					key := pat[i][1:len(pat[i])]
					v[key] = req[i]
				} else {
					return nil, false
				}
			}
		}
		return v, true
	}
	return nil, false
}
