package main

import "time"

// UIData represents data to be passed to UI.
type UIData struct {
	Services      []*Service
	Vars          []VarName
	UITheme       string
	LastTimestamp time.Time
}

// NewUIData inits and return new data object.
func NewUIData(vars []VarName, theme string) *UIData {
	return &UIData{
		Vars:    vars,
		UITheme: theme,
	}
}
