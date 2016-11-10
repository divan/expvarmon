package main

import "time"

// UIData represents data to be passed to UI.
type UIData struct {
	Services      []*Service
	Vars          []VarName
	LastTimestamp time.Time
	Stacks        map[VarName]*Stack
	Stats         map[VarName]*Stat
}

// NewUIData inits and return new data object.
func NewUIData(vars []VarName) *UIData {
	stacks := make(map[VarName]*Stack)
	stats := make(map[VarName]*Stat)
	for _, v := range vars {
		stacks[v] = NewStack()
		stats[v] = NewStat()
	}
	return &UIData{
		Vars:   vars,
		Stacks: stacks,
		Stats:  stats,
	}
}
