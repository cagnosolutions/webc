package main

import (
	"fmt"
	"net/http"

	"github.com/cagnosolutions/web"
	"github.com/cagnosolutions/web/tmpl"
	_ "github.com/davecheney/profile"
)

func main() {
	//defer profile.Start(profile.CPUProfile).Stop()
	web.Get("/", homeController)
	web.Get("/user", user)
	web.Get("/user/add", userAdd)
	web.Get("/user/:id", userId)
	web.Get("/reload", reloadTemplates)
	web.Get("/:slug", landing)
	web.Get("/login/:slug", multiLogin)
	web.Get("/logout/:slug", logout)
	web.Get("/protected/:slug", protected)
	web.Serve(":8080")
}

func homeController(w http.ResponseWriter, r *http.Request, c *web.Context) {
	msgK, msgV := c.GetFlash()
	tmpl.Render(w, "index.tmpl", tmpl.Model{
		msgK:    msgV,
		"name":  "Greg",
		"age":   28,
		"email": "gregpechiro@yahoo.com",
		"data":  "bubb hghb hghgv bj hbkjb jhb jhbjh bhjb jhb jhb jhbj bjk bjhb jhbj hbjhbjh jhb Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
	})
}

func landing(w http.ResponseWriter, r *http.Request, c *web.Context) {
	slug := c.GetPathVar("slug")
	_, msgV := c.GetFlash()
	fmt.Fprintf(w, "Your are at the landing page for %s. %s", slug, msgV)
}

func multiLogin(w http.ResponseWriter, r *http.Request, c *web.Context) {
	slug := c.GetPathVar("slug")
	c.SetAuth(true)
	c.SetFlash("success", "You are now logged into "+slug+". Enjoy")
	http.Redirect(w, r, "/"+slug, 303)
}

func logout(w http.ResponseWriter, r *http.Request, c *web.Context) {
	slug := c.GetPathVar("slug")
	c.SetAuth(false)
	c.SetFlash("success", "You are now logged out. Thanks for visiting")
	http.Redirect(w, r, "/"+slug, 303)
}

func protected(w http.ResponseWriter, r *http.Request, c *web.Context) {
	slug := c.GetPathVar("slug")
	c.CheckAuth(w, r, "/"+slug)
	fmt.Fprintf(w, "You are authorized to view page %s", slug)
}

func user(w http.ResponseWriter, r *http.Request, c *web.Context) {
	fmt.Fprintf(w, "page: user, addr: %s, user-agent: %s", r.RemoteAddr, r.UserAgent())
	return
}

func userAdd(w http.ResponseWriter, r *http.Request, c *web.Context) {
	fmt.Fprintf(w, "User Add Page")
}

func userId(w http.ResponseWriter, r *http.Request, c *web.Context) {
	fmt.Fprintf(w, "user id: %v", c.GetPathVar("id"))
}

func reloadTemplates(w http.ResponseWriter, r *http.Request, c *web.Context) {
	if r.FormValue("user") == "admin" && r.FormValue("pass") == "admin" {
		tmpl.Reload()
		c.SetFlash("alertSuccess", "Successfully reloaded templates")
	}
	http.Redirect(w, r, "/", 303)
	return
}
