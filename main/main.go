package main

import (
	"fmt"
	"net/http"

	. "github.com/cagnosolutions/web"
)

func main() {
	mux := MuxInstance()
	mux.Handle("GET", "/user", user)
	mux.Handle("GET", "/user/:id", userId)
	http.ListenAndServe(":8080", mux)
}

func user(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "page: user, addr: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
	return
}

func userId(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "user id: %v", c.GetPathVar("id"))
}
