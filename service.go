package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/antonholmquist/jason"
	//"github.com/pyk/byten"
)

type Services []*Service

// Service represents constantly updating info about single service.
type Service struct {
	Port    string
	Name    string
	Cmdline string

	values map[string]*Stack

	Err error
}

// NewService returns new Service object.
func NewService(port string) *Service {
	return &Service{
		Name: port, // we have only port on start, so use it as name until resolved
		Port: port,

		values: make(map[string]*Stack),
	}
}

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
		expvar, err := jason.NewObjectFromReader(resp.Body)
		e = *expvar
		if err != nil {
			return &e, err
		}
	}
	return &e, nil
}

// Update updates Service info from Expvar variable.
func (s *Service) Update() {
	expvar, err := FetchExpvar(s.Addr())
	if err != nil {
		s.Err = err
	}

	cmdline, err := expvar.GetStringArray("cmdline")
	if err != nil {
		s.Err = err
	} else {
		s.updateCmdline(cmdline)
	}

	alloc, err := expvar.GetInt64("memstats", "Alloc")
	if err != nil {
		s.Err = err
	} else {
		s.updateMem(alloc)
	}
}

func (s *Service) updateCmdline(cmdline []string) {
	// Update name and cmdline only if empty
	if len(s.Cmdline) == 0 {
		s.Cmdline = strings.Join(cmdline, " ")
		s.Name = BaseCommand(cmdline)
	}
}

func (s *Service) updateMem(alloc int64) {
	// Put metrics data
	mem, ok := s.values["mem.alloc"]
	if !ok {
		s.values["mem.alloc"] = NewStack(1200)
		mem = s.values["mem.alloc"]
	}
	mem.Push(int(alloc))
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

func (s Service) Value(key string) string {
	if s.Err != nil {
		return "N/A"
	}
	val, ok := s.values[key]
	if !ok {
		return "N/A"
	}
	if val.Front() == 0 {
		return "N/A"
	}

	//allocated := byten.Size(int64(val.Front()))
	//return fmt.Sprintf("Alloc: %s", allocated)
	return fmt.Sprintf("%d", val.Front())
}

func (s Service) Values(key string) []int {
	if s.Err != nil {
		return nil
	}
	val, ok := s.values[key]
	if !ok {
		return nil
	}

	return val.Values
}
