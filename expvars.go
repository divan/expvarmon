package main

import (
	"encoding/json"
	"io"
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

	Err error `json:"-,omitempty"`
}

func NewExpvarsSource(ports []string) *ExpvarsSource {
	return &ExpvarsSource{
		Ports: ports,
	}
}

// ParseExpvar unmarshals data to Expvar variable.
// TODO: implement Unmarshaller/Decode for Expvar
func ParseExpvar(r io.Reader) (*Expvar, error) {
	var vars Expvar
	dec := json.NewDecoder(r)
	err := dec.Decode(&vars)
	if err != nil {
		return nil, err
	}

	return &vars, err
}
