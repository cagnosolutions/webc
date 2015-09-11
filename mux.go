package web

import "net/http"

type Mux struct {
	routes []*Route
	ctx    *Context
}

func MuxInstance() *Mux {
	return &Mux{
		routes: make([]*Route, 0),
		ctx:    ContextInstance(),
	}
}

func (m *Mux) Handle(method, path string, controller Controller) {
	m.routes = append(m.routes, RouteInstance(method, path, controller, false))
}

func (m *Mux) SecureHandle(method, path string, controller Controller) {
	m.routes = append(m.routes, RouteInstance(method, path, controller, true))
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		return
	}
	for _, route := range m.routes {
		if route.path == r.URL.Path && route.method == r.Method {
			route.handle(w, r, m.ctx)
			return
		}
	}
	return
}
