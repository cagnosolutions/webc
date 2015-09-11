package web

import "github.com/scottcagno/safemap"

type Context struct {
	sm *safemap.SafeMap
}

func ContextInstance() *Context {
	ctx := &Context{
		sm: safemap.SafeMapInstance(32),
	}
	ctx.sm.Set("stack", make([]string, 0))
	return ctx
}

func (c *Context) Set(key string, val interface{}) {
	c.sm.Set(key, val)
}

func (c *Context) Get(key string) (interface{}, bool) {
	return c.sm.Get(key)
}

func (c *Context) Del(key string) {
	c.sm.Del(key)
}

func (c *Context) Push(val string) {
	stack, _ := c.sm.Get("stack")
	stack = append(stack.([]string), val)
	c.sm.Set("stack", stack)
}

func (c *Context) Pop() string {
	stack, _ := c.sm.Get("stack")
	var val string
	val, stack = stack.([]string)[0], stack.([]string)[1:]
	c.sm.Set("stack", stack)
	return val
}
