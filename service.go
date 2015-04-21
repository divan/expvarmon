package main

import (
	"runtime"
	"time"
)

// Data represents data to be passed to UI.
type Data struct {
	Services      Services
	TotalMemory   *Stack
	LastTimestamp time.Time
}

func NewData() *Data {
	return &Data{
		TotalMemory: NewStack(140),
	}
}

type Services []*Service

// Service represents constantly updating info about single service.
type Service struct {
	Name       string
	Port       string
	IsAlive    bool
	Cmdline    string
	Memstats   *runtime.MemStats
	Goroutines int64

	Err error
}

// NewService returns new Service object.
func NewService(port string) *Service {
	return &Service{
		Name: port, // we have only port on start, so use it as name until resolved
		Port: port,
	}
}

func (d *Data) FindService(port string) *Service {
	if d.Services == nil {
		return nil
	}
	for _, service := range d.Services {
		if service.Port == port {
			return service
		}
	}

	return nil
}
