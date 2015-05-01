package main

import (
	//"io"
	"runtime"
)

const ExpvarsUrl = "/debug/vars"

// ExpvarsSource implements Source interface for retrieving Expvars.
type ExpvarsSource struct {
	Ports []string
}

type Expvars map[string]Expvar

type Expvar struct {
	MemStats *runtime.MemStats `json:"memstats"`
	Cmdline  []string          `json:"cmdline"`

	Extra map[string]interface{} `json:"-"`

	Err error `json:"-,omitempty"`
}

func NewExpvarsSource(ports []string) *ExpvarsSource {
	return &ExpvarsSource{
		Ports: ports,
	}
}
