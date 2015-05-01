package main

import (
	"errors"
	"io"
	"net/http"

	"github.com/antonholmquist/jason"
)

const ExpvarsUrl = "/debug/vars"

type Expvar *jason.Object

// FetchExpvar fetches expvar by http for the given addr (host:port)
func FetchExpvar(addr string) (*jason.Object, error) {
	var e jason.Object
	resp, err := http.Get(addr)
	if err != nil {
		return &e, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return &e, errors.New("Vars not found. Did you import expvars?")
	} else {
		expvar, err := ParseExpvar(resp.Body)
		e = *expvar
		if err != nil {
			return &e, err
		}
	}
	return &e, nil
}

func ParseExpvar(r io.Reader) (*jason.Object, error) {
	return jason.NewObjectFromReader(r)
}
