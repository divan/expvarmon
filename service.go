package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

// Service represents constantly updating info about single service.
type Service struct {
	Port    string
	Name    string
	Cmdline string

	values map[VarName]*Stack

	Err error
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

		values: values,
	}
}

// Update updates Service info from Expvar variable.
func (s *Service) Update(wg *sync.WaitGroup) {
	defer wg.Done()
	expvar, err := FetchExpvar(s.Addr())
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
	for name, stack := range s.values {
		value, err := expvar.GetInt64(name.ToSlice()...)
		if err != nil {
			continue
		}
		stack.Push(int(value))
	}
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

// StatusLine returns status line for services with it's name and status.
func (s Service) StatusLine() string {
	if s.Err != nil {
		return fmt.Sprintf("[ERR] %s failed", s.Name)
	}

	return fmt.Sprintf("[R] %s", s.Name)
}

// Value returns current value for the given var of this service.
func (s Service) Value(name VarName) string {
	if s.Err != nil {
		return "N/A"
	}
	val, ok := s.values[name]
	if !ok {
		return "N/A"
	}
	if val.Front() == 0 {
		return "N/A"
	}

	return fmt.Sprintf("%d", val.Front())
}

// Values returns slice of ints with recent values of the given var,
// to be used with sparkline.
func (s Service) Values(name VarName) []int {
	if s.Err != nil {
		return nil
	}
	val, ok := s.values[name]
	if !ok {
		return nil
	}

	return val.Values
}
