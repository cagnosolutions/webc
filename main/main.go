package main

import (
	"fmt"
	"net/http"
	. "github.com/cagnosolutions/web"
)

func main() {
	ctx := ContextInstance()

	ctx.Set("foo", map[string]string{"foo:1": "bar:1", "foo:2": "bar:2"})
	val, _ := ctx.Get("foo")
	fmt.Println(val)
	fmt.Println(val.(map[string]string)["foo:2"])

	ctx.Push("one")
	ctx.Push("two")
	ctx.Push("three")

	val2, _ := ctx.Get("stack")
	fmt.Println(val2)

	val3 := ctx.Pop()
	fmt.Println(val3)

	val4, _ := ctx.Get("stack")
	fmt.Println(val4)

	mux := MuxInstance()
	mux.Handle("GET", "/", index)
	mux.Handle("GET", "/home", home)
	mux.Handle("GET", "/user", user)
	mux.Handle("GET", "/user/:id", userId)

	http.ListenAndServe(":8080", mux)

}

func index(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "%s", GetUser(r))
	//fmt.Fprintf(w, "page: index, addr: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
	return
}

func home(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "page: home, addr: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
	return
}

func user(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "page: user, addr: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
	return
}

func userId(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "user id: %v", )
}
