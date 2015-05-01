package main

import "time"

// UIData represents data to be passed to UI.
type UIData struct {
	Services      Services
	Vars          []VarName
	LastTimestamp time.Time
}

// NewUIData inits and return new data object.
func NewUIData(vars []VarName) *UIData {
	return &UIData{
		Vars: vars,
	}
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
