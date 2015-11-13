package main

import (
	"fmt"
	"net/http"

	"github.com/cagnosolutions/webc"
	"github.com/cagnosolutions/webc/tmpl"
	_ "github.com/davecheney/profile"
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	webc.Get("/", homeController)
	webc.Get("/user", user)
	webc.Get("/user/add", userAdd)
	webc.Get("/user/:id", userId)
	webc.Get("/reload", reloadTemplates)
	webc.Get("/:slug", landing)
	webc.Get("/login/:slug", multiLogin)
	webc.Get("/logout/:slug", logout)
	webc.Get("/protected/:slug", protected)
	webc.Serve(":8080")
}

func homeController(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	msgK, msgV := c.GetFlash()
	tmpl.Render(w, "index.tmpl", tmpl.Model{
		msgK:    msgV,
		"name":  "Greg",
		"age":   28,
		"email": "gregpechiro@yahoo.com",
		"data":  "bubb hghb hghgv bj hbkjb jhb jhbjh bhjb jhb jhb jhbj bjk bjhb jhbj hbjhbjh jhb Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})
}

func landing(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	slug := c.GetPathVar("slug")
	_, msgV := c.GetFlash()
	fmt.Fprintf(w, "Your are at the landing page for %s. %s", slug, msgV)
}

func multiLogin(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	slug := c.GetPathVar("slug")
	c.SetAuth(true)
	c.SetFlash("success", "You are now logged into "+slug+". Enjoy")
	http.Redirect(w, r, "/"+slug, 303)
}

func logout(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	slug := c.GetPathVar("slug")
	c.SetAuth(false)
	c.SetFlash("success", "You are now logged out. Thanks for visiting")
	http.Redirect(w, r, "/"+slug, 303)
}

func protected(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	slug := c.GetPathVar("slug")
	c.CheckAuth(w, r, "/"+slug)
	fmt.Fprintf(w, "You are authorized to view page %s", slug)
}

func user(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	fmt.Fprintf(w, "page: user, addr: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
	return
}

func userAdd(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	fmt.Fprintf(w, "User Add Page")
}

func userId(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	fmt.Fprintf(w, "user id: %v", c.GetPathVar("id"))
}

func reloadTemplates(w http.ResponseWriter, r *http.Request, c *webc.Context) {
	if r.FormValue("user") == "admin" && r.FormValue("pass") == "admin" {
		tmpl.Reload()
		c.SetFlash("alertSuccess", "Successfully reloaded templates")
	}
	http.Redirect(w, r, "/", 303)
	return
}
