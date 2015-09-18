package main

import (
	"github.com/cagnosolutions/web"
)

// GC Rate
// Production/development
// session id

func main() {
	// default
	mux := &web.Mux{}

	// custom config
	mux := &web.Mux{
		Rate:       HOUR / 2,
		SessId:     "GOSESS",
		LoggerPath: "/opt/logs",
	}

	mux.Get("/", homeController)

	mux.Serve(":8080")

}
