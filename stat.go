package expvarmon

// Stat holds basic statistics data for
// integer data used for sparklines.
type Stat struct {
    max    int
    maxStr string
}

// NewStat inits new Stat object.
func NewStat() *Stat {
    return &Stat{}
}

// Update updates stats on each push.
func (s *Stat) Update(v IntVar) {
    if v.Value() > s.max {
        s.max = v.Value()
        s.maxStr = v.String()
    }
}

// Max returns maximum recorded value.
func (s *Stat) Max() string {
    return s.maxStr
}
