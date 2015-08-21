package main

import (
	"net/url"
	"strings"
	"sync"

	"github.com/antonholmquist/jason"
)

var (
	// uptimeCounter is a variable used for tracking uptime status.
	// It should be always incrementing and included into default expvar vars.
	// Could be replaced with something different or made configurable in
	// the future.
	uptimeCounter = VarName("memstats.PauseTotalNs").ToSlice()
)

// Service represents constantly updating info about single service.
type Service struct {
	URL     url.URL
	Name    string
	Cmdline string

	stacks map[VarName]*Stack

	Err           error
	Restarted     bool
	UptimeCounter int64
}

// NewService returns new Service object.
func NewService(url url.URL, vars []VarName) *Service {
	values := make(map[VarName]*Stack)
	for _, name := range vars {
		values[VarName(name)] = NewStack()
	}

	return &Service{
		Name: url.Host, // we have only port on start, so use it as name until resolved
		URL:  url,

		stacks: values,
	}
}

// Update updates Service info from Expvar variable.
func (s *Service) Update(wg *sync.WaitGroup) {
	defer wg.Done()
	expvar, err := FetchExpvar(s.URL)
	// check for restart
	if s.Err != nil && err == nil {
		s.Restarted = true
	}
	s.Err = err

	// if memstat.PauseTotalNs less than s.UptimeCounter
	// then service was restarted
	c, err := expvar.GetInt64(uptimeCounter...)
	if err != nil {
		s.Err = err
	} else {
		if s.UptimeCounter > c {
			s.Restarted = true
		}
		s.UptimeCounter = c
	}

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
	} else if v, err := value.Array(); err == nil {
		// if we get an array, calculate average

		// empty array, treat as zero
		if len(v) == 0 {
			return 0
		}

		avg := averageJason(v)

		// cast to int64 for Int64 values
		if _, err := v[0].Int64(); err == nil {
			return int64(avg)
		}

		return avg
	}

	return nil
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
