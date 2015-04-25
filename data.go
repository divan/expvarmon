package main

import "time"

// Data represents data to be passed to UI.
type Data struct {
	Services      Services
	TotalMemory   *Stack
	LastTimestamp time.Time
}

// NewData inits and return new data object.
func NewData() *Data {
	return &Data{
		TotalMemory: NewStack(140),
	}
}

// FindService returns existing service by port.
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
