package webc

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cagnosolutions/webc/util"
)

const (
	MIN     = 60
	HOUR    = MIN * 60
	DAY     = HOUR * 24
	WEEK    = DAY * 7
	MONTH   = DAY * 30
	YEAR    = WEEK * 52
	SESSION = 0
)

type contextStore struct {
	contexts map[string]*Context
	Ctxid    string
	Rate     int64
	sync.Mutex
}

func NewContextStore(ctxid string, rate int64) *contextStore {
	c := &contextStore{
		contexts: make(map[string]*Context),
		Ctxid:    ctxid,
		Rate:     rate,
	}
	return c
}

func (cs *contextStore) get(w http.ResponseWriter, r *http.Request) *Context {
	uuid, ok := cs.getId(r)
	if ok { // uuid (cookie) found
		cs.Lock()
		if ctx, ok := cs.contexts[uuid]; ok {
			// update and return context
			ctx.ts = time.Now()
			cs.Unlock()
			return ctx
		}
		// context dead, create new one based on uuid and return new context
		ctx := freshContext(uuid)
		cs.contexts[uuid] = ctx
		cs.Unlock()
		return ctx
	}
	// uuid (cookie) not foud, create and set a new one
	cookie := cs.freshCookie(uuid)
	http.SetCookie(w, &cookie)
	// add or over-write any context with same uuid, and return context
	cs.Lock()
	ctx := freshContext(uuid)
	cs.contexts[uuid] = ctx
	cs.Unlock()
	return ctx
}

func (cs *contextStore) gc() {
	cs.Lock()
	defer cs.Unlock()
	for uuid, ctx := range cs.contexts {
		if (ctx.ts.Unix() + cs.Rate) < time.Now().Unix() {
			delete(cs.contexts, uuid)
		} else {
			break
		}
	}
	time.AfterFunc(time.Duration(cs.Rate)*time.Second, func() { cs.gc() })
}

func (cs *contextStore) viewContexts() {
	for k, v := range cs.contexts {
		fmt.Printf("key: %v\nval: %v\n\n", k, v)
	}
}

func (cs *contextStore) getId(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(cs.Ctxid)
	if err != nil && err == http.ErrNoCookie || cookie.Value == "" {
		return util.UUID4(), false
	}
	return cookie.Value, true
}

func (cs *contextStore) freshCookie(uuid string) http.Cookie {
	return http.Cookie{
		Name:     cs.Ctxid,
		Value:    uuid,
		Path:     "/",
		Expires:  time.Now().AddDate(3, 0, 0), // 3 years in the future
		HttpOnly: true,
	}
}

func freshContext(uuid string) *Context {
	return &Context{
		uuid:    uuid,
		ts:      time.Now(),
		items:   make(map[string]interface{}, 0),
		path:    make(map[string]string, 0),
		flash:   make([]string, 0),
		session: make(map[string]interface{}, 0),
		role:    "",
	}
}

type Context struct {
	uuid    string
	ts      time.Time
	items   map[string]interface{}
	path    map[string]string
	flash   []string
	session map[string]interface{}
	role    string
}

func (c *Context) SetPathVars(m map[string]string) {
	c.path = m
}

func (c *Context) GetPathVars() map[string]string {
	return c.path
}

func (c *Context) GetPathVar(k string) string {
	return c.path[k]
}

func (c *Context) Set(k string, v interface{}) {
	c.items[k] = v
}

func (c *Context) Get(k string) interface{} {
	return c.items[k]
}

func (c *Context) GetAll() map[string]interface{} {
	return c.items
}

func (c *Context) Del(k string) {
	delete(c.items, k)
}

func (c *Context) SetFlash(k, msg string) {
	c.flash = []string{k, msg}
}

func (c *Context) GetFlashSlice() []string {
	flash := c.flash
	c.flash = []string{}
	return flash
}

func (c *Context) GetFlash() (string, string) {
	flash := c.flash
	c.flash = []string{}
	if len(flash) == 2 {
		return flash[0], flash[1]
	}
	return "", ""
}

func (c *Context) GetSession() map[string]interface{} {
	return c.session
}

func (c *Context) SetSession(session map[string]interface{}) {
	c.session = session
}

func (c *Context) GetFromSession(key string) interface{} {
	return c.session[key]
}

func (c *Context) SetToSession(key string, val interface{}) {
	c.session[key] = val
}

func (c *Context) SetRole(role string) {
	c.role = role
}

func (c *Context) GetRole() string {
	return c.role
}

func (c *Context) CheckAuth(w http.ResponseWriter, r *http.Request, path string, requiredRoles ...string) bool {
	if len(requiredRoles) < 1 {
		return true
	}
	if c.role == "" {
		c.SetFlash("alertError", "Your are not logged in!")
		http.Redirect(w, r, path, 303)
		return false
	}
	for _, role := range requiredRoles {
		if c.role == role {
			return true
		}
	}
	c.SetFlash("alertError", "You are not authorized")
	http.Redirect(w, r, "/", 303)
	return false
}

func (c *Context) Login(role string) {
	c.role = role
}

func (c *Context) Logout() {
	c.session = make(map[string]interface{}, 0)
	c.role = ""
}
