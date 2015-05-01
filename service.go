package main

import (
	"fmt"
	"strings"

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
func NewService(port string, vars []string) *Service {
	values := make(map[string]*Stack)
	for _, name := range vars {
		values[name] = NewStack()
	}
	return &Service{
		Name: port, // we have only port on start, so use it as name until resolved
		Port: port,

		values: values,
	}
}

// Update updates Service info from Expvar variable.
func (s *Service) Update() {
	expvar, err := FetchExpvar(s.Addr())
	s.Err = err

	cmdline, err := expvar.GetStringArray("cmdline")
	if err != nil {
		s.Err = err
	} else {
		s.updateCmdline(cmdline)
	}

	for name, stack := range s.values {
		value, err := expvar.GetInt64(dot2slice(name)...)
		if err != nil {
			continue
		}
		stack.Push(int(value))
	}
}

func (s *Service) updateCmdline(cmdline []string) {
	// Update name and cmdline only if empty
	// TODO: move it to Update() with sync.Once
	if len(s.Cmdline) == 0 {
		s.Cmdline = strings.Join(cmdline, " ")
		s.Name = BaseCommand(cmdline)
	}
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
