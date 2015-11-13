package main

import (
	"fmt"
	"net/http"

	. "github.com/cagnosolutions/webc"
	"github.com/davecheney/profile"
)

var ts = TemplateStoreInstance("templates", nil)

func main() {

	/* # # # # #  ENABLING PROFILING # # # # # */
	defer profile.Start(profile.CPUProfile).Stop()
	/*
		profileConfig := profile.Config{
			CPUProfile:     true,
			MemProfile:     true,
			ProfilePath:    ".",  // store profiles in current directory
			NoShutdownHook: true, // do not hook SIGINT
		}
		p := profile.Start(&profileConfig)
		defer p.Stop()
	*/
	/* # # # # #  ENABLING PROFILING # # # # # */

	ts.Cache("index.html", "home.html", "404.html")

	mux := MuxInstance()
	mux.Handle("GET", "/index", index)
	mux.Handle("GET", "/home", home)
	mux.Handle("GET", "/404", err404)
	mux.Handle("GET", "/user", user)
	mux.Handle("GET", "/user/add", userAdd)
	mux.Handle("GET", "/user/:id", userId)
	mux.Handle("GET", "/:slug", landing)
	mux.Handle("GET", "/login/:slug", multiLogin)
	mux.Handle("GET", "/logout/:slug", logout)
	mux.Handle("GET", "/protected/:slug", protected)

	http.ListenAndServe(":8080", mux)
}

func index(w http.ResponseWriter, r *http.Request, c *Context) {
	ts.Render(w, "index.html", M{"page": "INDEX PAGE"})
}

func home(w http.ResponseWriter, r *http.Request, c *Context) {
	ts.Render(w, "home.html", M{"page": "HOME PAGE"})
}

func err404(w http.ResponseWriter, r *http.Request, c *Context) {
	ts.Render(w, "404.html", M{"page": "ERROR 404 PAGE"})
}

func landing(w http.ResponseWriter, r *http.Request, c *Context) {
	slug := c.GetPathVar("slug")
	msgS := c.GetFlash()
	var msg string
	if len(msgS) == 2 {
		msg = msgS[1]
	}
	fmt.Fprintf(w, "Your are at the landing page for %s. %s", slug, msg)
}

func multiLogin(w http.ResponseWriter, r *http.Request, c *Context) {
	slug := c.GetPathVar("slug")
	c.SetAuth(true)
	c.SetFlash("success", "You are now logged into "+slug+". Enjoy")
	http.Redirect(w, r, "/"+slug, 303)
}

func logout(w http.ResponseWriter, r *http.Request, c *Context) {
	slug := c.GetPathVar("slug")
	c.SetAuth(false)
	c.SetFlash("success", "You are now logged out. Thanks for visiting")
	http.Redirect(w, r, "/"+slug, 303)
}

func protected(w http.ResponseWriter, r *http.Request, c *Context) {
	slug := c.GetPathVar("slug")
	c.CheckAuth(w, r, "/"+slug)
	fmt.Fprintf(w, "You are authorized to view page %s", slug)
}

func user(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "page: user, addr: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
	return
}

func userAdd(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "User Add Page")
}

func userId(w http.ResponseWriter, r *http.Request, c *Context) {
	fmt.Fprintf(w, "user id: %v", c.GetPathVar("id"))
}
