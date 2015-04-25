package main

import (
	"fmt"
	"github.com/pyk/byten"
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
	MemStats *runtime.MemStats

	Values map[string]*Stack

	Err error
}

// NewService returns new Service object.
func NewService(port string) *Service {
	return &Service{
		Name: port, // we have only port on start, so use it as name until resolved
		Port: port,

		Values: make(map[string]*Stack),
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
	s.MemStats = expvar.MemStats

	// Update name and cmdline only if empty
	if len(s.Cmdline) == 0 {
		s.Cmdline = strings.Join(expvar.Cmdline, " ")
		s.Name = BaseCommand(expvar.Cmdline)
	}

	// Put metrics data
	mem, ok := s.Values["memory"]
	if !ok {
		s.Values["memory"] = NewStack(1200)
		mem = s.Values["memory"]
	}
	mem.Push(int(s.MemStats.Alloc) / 1024)
}

// Addr returns fully qualified host:port pair for service.
//
// If host is not specified, 'localhost' is used.
func (s Service) Addr() string {
	return fmt.Sprintf("http://localhost:%s%s", s.Port, ExpvarsUrl)
}

// StatusLine returns status line for services with it's name and status.
func (s Service) StatusLine() string {
	if s.Err != nil {
		return fmt.Sprintf("[ERR] %s failed", s.Name)
	}

	return fmt.Sprintf("[R] %s", s.Name)
}

// Meminfo returns memory info string for the given service.
func (s Service) Meminfo() string {
	if s.Err != nil || s.MemStats == nil {
		return "N/A"
	}

	allocated := byten.Size(int64(s.MemStats.Alloc))
	sys := byten.Size(int64(s.MemStats.Sys))
	return fmt.Sprintf("Alloc/Sys: %s / %s", allocated, sys)
}
