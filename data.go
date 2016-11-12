package main

import "time"

// UIData represents data to be passed to UI.
type UIData struct {
	Services      []*Service
	Vars          []VarName
	LastTimestamp time.Time

	SparklineData []*SparklineData
}

// SparklineData holds additional data needed for sparklines.
type SparklineData struct {
	Stacks map[VarName]*Stack
	Stats  map[VarName]*Stat
}

// NewSparklineData inits new Sparkline data object.
func NewSparklineData(vars []VarName) *SparklineData {
	stacks := make(map[VarName]*Stack)
	stats := make(map[VarName]*Stat)
	for _, v := range vars {
		stacks[v] = NewStack()
		stats[v] = NewStat()
	}
	return &SparklineData{
		Stacks: stacks,
		Stats:  stats,
	}
}

// NewUIData inits and return new data object.
func NewUIData(vars []VarName, services []*Service) *UIData {
	sp := make([]*SparklineData, len(services))
	for i, _ := range services {
		sp[i] = NewSparklineData(vars)
	}
	return &UIData{
		Services:      services,
		Vars:          vars,
		SparklineData: sp,
	}
}
