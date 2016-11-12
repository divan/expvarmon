package main

// Stat holds basic statistics data for
// integer data used for sparklines.
type Stat struct {
	max IntVar
	// TODO: implement running median
}

// NewStat inits new Stat object.
func NewStat() *Stat {
	return &Stat{
		max: &Number{},
	}
}

// Update updates stats on each push.
func (s *Stat) Update(v IntVar) {
	if v.Value() > s.max.Value() {
		s.max = v
	}
}

// Max returns maximum recorded value.
func (s *Stat) Max() IntVar {
	return s.max
}
