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

type Mux struct {
	routes []*route
	ctx    *contextStore
	static http.Handler
}

func NewMux(ctxid string, rate int64) *Mux {
	return &Mux{
		routes: make([]*route, 0),
		ctx:    NewContextStore(ctxid, rate),
		static: http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	}
}

func (m *Mux) Serve(host string) {
	m.ctx.gc()
	if err := http.ListenAndServe(host, m); err != nil {
		panic(err)
	}
}

func (m *Mux) handle(method, path string, handler controller) {
	m.routes = append(m.routes, &route{method, SliceString(path, '/'), handler})
}

func (m *Mux) Get(path string, handler controller) {
	m.handle("GET", path, handler)
}

func (m *Mux) Put(path string, handler controller) {
	m.handle("PUT", path, handler)
}

func (m *Mux) Post(path string, handler controller) {
	m.handle("POST", path, handler)
}

func (m *Mux) Delete(path string, handler controller) {
	m.handle("DELETE", path, handler)
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// for now: ignore options and favicon
	if r.Method == "OPTIONS" || r.URL.Path == "/favicon.ico" {
		return
	}

	//handle static mapping
	if r.Method == "GET" && strings.HasPrefix(r.URL.Path, "/static/") {
		m.static.ServeHTTP(w, r)
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

/*
func match(req []string, pat []string) (map[string]string, bool) {
	vals := make(map[string]string)
	if len(req) != len(pat) {
		return nil, false
	}
	for i, v := range pat {
		if v[0] == ':' {
			vals[v[1:len(v)]] = req[i]
		} else if v != req[i] {
			return nil, false
		}
	}

	return vals, true
}
*/

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
