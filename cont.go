package web

import "net/http"

type Controller func(http.ResponseWriter, *http.Request, *Context)
