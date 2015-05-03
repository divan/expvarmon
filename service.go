package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/antonholmquist/jason"
)

// Service represents constantly updating info about single service.
type Service struct {
	Port    string
	Name    string
	Cmdline string

	stacks map[VarName]*Stack

	Err       error
	Restarted bool
}

// NewService returns new Service object.
func NewService(port string, vars []VarName) *Service {
	values := make(map[VarName]*Stack)
	for _, name := range vars {
		values[VarName(name)] = NewStack()
	}

	return &Service{
		Name: port, // we have only port on start, so use it as name until resolved
		Port: port,

		stacks: values,
	}
}

// Update updates Service info from Expvar variable.
func (s *Service) Update(wg *sync.WaitGroup) {
	defer wg.Done()
	expvar, err := FetchExpvar(s.Addr())
	// check for restart
	if s.Err != nil && err == nil {
		s.Restarted = true
	}
	s.Err = err

	// Update Cmdline & Name only once
	if len(s.Cmdline) == 0 {
		cmdline, err := expvar.GetStringArray("cmdline")
		if err != nil {
			s.Err = err
		} else {
			s.Cmdline = strings.Join(cmdline, " ")
			s.Name = BaseCommand(cmdline)
		}
	}

	// For all vars, fetch desired value from Json and push to it's own stack.
	for name, stack := range s.stacks {
		value, err := expvar.GetValue(name.ToSlice()...)
		if err != nil {
			stack.Push(nil)
			continue
		}
		v := guessValue(value)
		if v != nil {
			stack.Push(v)
		}
	}
}

// guessValue attemtps to bruteforce all supported types.
func guessValue(value *jason.Value) interface{} {
	if v, err := value.Int64(); err == nil {
		return v
	} else if v, err := value.Float64(); err == nil {
		return v
	} else if v, err := value.Boolean(); err == nil {
		return v
	} else if v, err := value.String(); err == nil {
		return v
	}

	return nil
}

// Addr returns fully qualified host:port pair for service.
//
// If host is not specified, 'localhost' is used.
func (s Service) Addr() string {
	// Try as port only
	_, err := strconv.Atoi(s.Port)
	if err == nil {
		return fmt.Sprintf("http://localhost:%s%s", s.Port, ExpvarsURL)
	}

	host, port, err := net.SplitHostPort(s.Port)
	if err == nil {
		return fmt.Sprintf("http://%s:%s%s", host, port, ExpvarsURL)
	}

	return ""
}

// Value returns current value for the given var of this service.
//
// It also formats value, if kind is specified.
func (s Service) Value(name VarName) string {
	if s.Err != nil {
		return "N/A"
	}
	val, ok := s.stacks[name]
	if !ok {
		return "N/A"
	}

	v := val.Front()
	if v == nil {
		return "N/A"
	}

	return Format(v, name.Kind())
}

// Values returns slice of ints with recent
// values of the given var, to be used with sparkline.
func (s Service) Values(name VarName) []int {
	stack, ok := s.stacks[name]
	if !ok {
		return nil
	}

	return stack.IntValues()
}

// Max returns maximum recorded value for given service and var.
func (s Service) Max(name VarName) interface{} {
	val, ok := s.stacks[name]
	if !ok {
		return nil
	}

	v := val.Max
	if v == nil {
		return nil
	}

	return Format(v, name.Kind())
}
