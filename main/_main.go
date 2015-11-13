package main

import (
	"github.com/cagnosolutions/webc"
)

// GC Rate
// Production/development
// session id

func main() {
	// default
	mux := &webc.Mux{}

	// custom config
	mux := &webc.Mux{
		Rate:       HOUR / 2,
		SessId:     "GOSESS",
		LoggerPath: "/opt/logs",
	}

	mux.Get("/", homeController)

	mux.Serve(":8080")

}
