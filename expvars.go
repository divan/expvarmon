package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
)

// Expvars holds all vars we support via expvars. It implements Source interface.

const ExpvarsUrl = "/debug/vars"

type ExpvarsSource struct {
	Ports []string
}

type Expvars map[string]Expvar

type Expvar struct {
	MemStats   *runtime.MemStats `json:"memstats"`
	Cmdline    []string          `json:"cmdline"`
	Goroutines int64             `json:"goroutines,omitempty"`

	Err error `json:"-,omitempty"`
}

func NewExpvarsSource(ports []string) *ExpvarsSource {
	return &ExpvarsSource{
		Ports: ports,
	}
}

func (e *ExpvarsSource) Update() (interface{}, error) {
	vars := make(Expvars)
	for _, port := range e.Ports {
		addr := fmt.Sprintf("http://localhost:%s%s", port, ExpvarsUrl)
		resp, err := http.Get(addr)
		if err != nil {
			expvar := &Expvar{}
			expvar.Err = err
			vars[port] = *expvar
			continue
		}
		if resp.StatusCode == http.StatusNotFound {
			expvar := &Expvar{}
			expvar.Err = fmt.Errorf("Page not found. Did you import expvars?")
			vars[port] = *expvar
			continue
		}
		defer resp.Body.Close()

		expvar, err := ParseExpvar(resp.Body)
		if err != nil {
			expvar = &Expvar{}
			expvar.Err = err
		}

		vars[port] = *expvar
	}

	return vars, nil
}

// ParseExpvar unmarshals data to Expvar variable.
func ParseExpvar(r io.Reader) (*Expvar, error) {
	var vars Expvar
	dec := json.NewDecoder(r)
	err := dec.Decode(&vars)
	if err != nil {
		return nil, err
	}

	return &vars, err
}
