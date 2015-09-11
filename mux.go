package web

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

type Controller func(http.ResponseWriter, *http.Request, *Context)

type Route struct {
	method, path string
	handle       Controller
	secure       bool
}

func RouteInstance(method, path string, handle Controller, secure bool) *Route {
	return &Route{
		method: method,
		path:   path,
		handle: handle,
		secure: secure,
	}
}

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

// ##### HELPERS ##### //

func GetUser(r *http.Request) string {
	return base64.StdEncoding.EncodeToString([]byte(r.RemoteAddr + r.UserAgent()))
}

func UUID4() string {
	u := make([]byte, 16)
	if _, err := rand.Read(u[:16]); err != nil {
		log.Println(err)
	}
	u[8] = (u[8] | 0x80) & 0xbf
	u[6] = (u[6] | 0x40) & 0x4f
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}
