package main

import (
	"errors"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/antonholmquist/jason"
)

// ExpvarsUrl is the default url for fetching expvar info.
const ExpvarsURL = "/debug/vars"

// Expvar represents fetched expvar variable.
type Expvar struct {
	*jason.Object
}

func getBasicAuthEnv() (user, password string) {
	return os.Getenv("HTTP_USER"), os.Getenv("HTTP_PASSWORD")
}

// FetchExpvar fetches expvar by http for the given addr (host:port)
func FetchExpvar(addr string) (*Expvar, error) {
	e := &Expvar{&jason.Object{}}
	client := &http.Client{
		Timeout: 1 * time.Second, // TODO: make it configurable or left default?
	}

	req, _ := http.NewRequest("GET", addr, nil)
	if user, pass := getBasicAuthEnv(); user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}
	resp, err := client.Do(req)
	if err != nil {
		return e, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return e, errors.New("Vars not found. Did you import expvars?")
	}

	expvar, err := ParseExpvar(resp.Body)
	if err != nil {
		return e, err
	}
	e = expvar
	return e, nil
}

// ParseExpvar parses expvar data from reader.
func ParseExpvar(r io.Reader) (*Expvar, error) {
	object, err := jason.NewObjectFromReader(r)
	return &Expvar{object}, err
}
