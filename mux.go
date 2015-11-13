package webc

import "net/http"

type Mux struct {
	routes []*Route
	ctx    *ContextStore
}

func MuxInstance() *Mux {
	return &Mux{
		routes: make([]*Route, 0),
		ctx:    ContextStoreInstance(HOUR / 2),
	}
}

func (m *Mux) Handle(method, path string, controller Controller) {
	m.routes = append(m.routes, RouteInstance(method, SliceString(path, '/'), controller, false))
}

func (m *Mux) SecureHandle(method, path string, controller Controller) {
	m.routes = append(m.routes, RouteInstance(method, SliceString(path, '/'), controller, true))
}

func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" || r.URL.Path == "/favicon.ico" {
		return
	}
	//fmt.Println(r.URL.Path)
	for _, route := range m.routes {
		if route.method == r.Method {
			path := SliceString(r.URL.Path, '/')
			if pathVars, ok := match(path, route.path); ok {
				ctx := m.ctx.GetContext(w, r)
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
