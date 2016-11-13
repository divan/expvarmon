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

	Vars map[VarName]Var

	Err           error
	Restarted     bool
	UptimeCounter int64
}

// NewService returns new Service object.
func NewService(url url.URL, varNames []VarName) *Service {
	vars := make(map[VarName]Var)
	for _, name := range varNames {
		v := NewVar(name)
		vars[name] = v
	}

	return &Service{
		Name: url.Host, // we have only port on start, so use it as name until resolved
		URL:  url,

		Vars: vars,
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

	// For all vars, fetch desired value from JSON
	for name, v := range s.Vars {
		value, err := expvar.GetValue(name.ToSlice()...)
		if err != nil {
			v.Set(nil)
		}
		v.Set(value)
	}
}

// guessValue attemtps to bruteforce all supported types.
// TODO: FIXME: remove this func
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

// Value returns current value for the given varable of this service.
func (s Service) Value(name VarName) string {
	v, ok := s.Vars[name]
	if !ok {
		return ""
	}
	return v.String()
}
