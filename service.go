package main

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

type Services []*Service

// Service represents constantly updating info about single service.
type Service struct {
	Port string
	Name string

	Cmdline  string
	Memstats *runtime.MemStats

	Err error
}

// NewService returns new Service object.
func NewService(port string) *Service {
	return &Service{
		Name: port, // we have only port on start, so use it as name until resolved
		Port: port,
	}
}

// Update updates Service info from Expvar variable.
func (s *Service) Update() {
	expvar := &Expvar{}
	resp, err := http.Get(s.Addr())
	defer resp.Body.Close()
	if err != nil {
		expvar.Err = err
	} else if resp.StatusCode == http.StatusNotFound {
		expvar.Err = fmt.Errorf("Vars not found. Did you import expvars?")
	} else {
		expvar, err = ParseExpvar(resp.Body)
		if err != nil {
			expvar = &Expvar{Err: err}
		}
	}

	s.Err = expvar.Err
	s.Memstats = expvar.MemStats

	// Update name and cmdline only if empty
	if len(s.Cmdline) == 0 {
		s.Cmdline = strings.Join(expvar.Cmdline, " ")
		s.Name = BaseCommand(expvar.Cmdline)
	}
}

// Addr returns fully qualified host:port pair for service.
//
// If host is not specified, 'localhost' is used.
func (s Service) Addr() string {
	return fmt.Sprintf("http://localhost:%s%s", s.Port, ExpvarsUrl)
}
