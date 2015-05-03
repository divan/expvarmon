package main

import (
	_ "expvar"
	"net/http"
)

func StartHttp(bind string) {
	// silently discard error here
	http.ListenAndServe(":"+bind, nil)
}
