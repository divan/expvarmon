package main

import "time"

// UIData represents data to be passed to UI.
type UIData struct {
	Services      Services
	Total         int
	LastTimestamp time.Time
}

// NewData inits and return new data object.
func NewData() *UIData {
	return &UIData{}
}

// FindService returns existing service by port.
func (d *UIData) FindService(port string) *Service {
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
