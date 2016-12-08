package expvarmon

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/antonholmquist/jason"
	"github.com/pyk/byten"
)

// VarName represents variable name.
//
// It has dot-separated format, like "memstats.Alloc",
// but can be used in different forms, hence it's own type.
//
// It also can have optional "kind:" modifier, like "mem:" or "duration:"
type VarName string

// VarKind specifies special kinds of values, affects formatting.
type VarKind int

const (
	KindDefault VarKind = iota
	KindMemory
	KindDuration
	KindString
	KindGCPauses
	KindGCIntervals
)

// Var represents arbitrary value for variable.
type Var interface {
	Kind() VarKind
	String() string
	Set(*jason.Value)
	SetNA() // mark as N/A
}

// IntVar represents variable which value can be represented as integer,
// and suitable for displaying with sparklines.
type IntVar interface {
	Var
	Value() int
}

// NewVar inits new Var object with the given name.
func NewVar(name VarName) Var {
	kind := name.Kind()

	switch kind {
	case KindDefault:
		return &Number{}
	case KindMemory:
		return &Memory{}
	case KindDuration:
		return &Duration{}
	case KindString:
		return &String{}
	case KindGCPauses:
		return &GCPauses{}
	case KindGCIntervals:
		return &GCIntervals{}
	default:
		return &Number{}
	}
}

// Number is a type for numeric values, obtained from JSON.
// In JSON it's always float64, so there is no straightforward way
// to separate float from int, so let's keep everything as float.
type Number struct {
	val float64
	na  bool
}

func (v *Number) Kind() VarKind { return KindDefault }
func (v *Number) String() string {
	if v.na {
		return "N/A"
	}
	// if fraction part is zero, assume int's integer
	if _, frac := math.Modf(v.val); frac == 0 {
		return fmt.Sprintf("%.0f", v.val)
	}
	// else, return as float
	return fmt.Sprintf("%.02f", v.val)
}
func (v *Number) Set(j *jason.Value) {
	v.na = false
	if n, err := j.Float64(); err == nil {
		v.val = n
	} else if n, err := j.Int64(); err == nil {
		v.val = float64(n)
	} else {
		v.val = 0
	}
}
func (v *Number) SetNA() {
	v.na = true
	v.val = 0
}

// Value implements IntVar for Number type.
func (v *Number) Value() int {
	return int(v.val)
}

// Memory represents memory information in bytes.
type Memory struct {
	bytes int64
	na    bool
}

func (v *Memory) Kind() VarKind { return KindMemory }
func (v *Memory) String() string {
	if v.na {
		return "N/A"
	}
	return fmt.Sprintf("%s", byten.Size(v.bytes))
}
func (v *Memory) Set(j *jason.Value) {
	v.na = false
	if n, err := j.Int64(); err == nil {
		v.bytes = n
	} else {
		v.bytes = 0
	}
}

// Value implements IntVar for Memory type.
func (v *Memory) Value() int {
	// TODO: check for possible overflows
	return int(v.bytes)
}
func (v *Memory) SetNA() {
	v.na = true
	v.bytes = 0
}

// Duration represents duration data (in ns)
type Duration struct {
	dur time.Duration
	na  bool
}

func (v *Duration) Kind() VarKind { return KindDuration }
func (v *Duration) String() string {
	if v.na {
		return "N/A"
	}
	return fmt.Sprintf("%s", round(time.Duration(v.dur)))
}

func (v *Duration) Set(j *jason.Value) {
	v.na = false
	if n, err := j.Int64(); err == nil {
		v.dur = time.Duration(n)
	} else if n, err := j.Float64(); err == nil {
		v.dur = time.Duration(int64(n))
	} else {
		v.dur = 0
	}
}

// Value implements IntVar for Duration type.
func (v *Duration) Value() int {
	// TODO: check for possible overflows
	return int(v.dur)
}
func (v *Duration) SetNA() {
	v.na = true
	v.dur = 0
}

// Strings represents string data.
type String struct {
	str string
}

func (v *String) Kind() VarKind  { return KindString }
func (v *String) String() string { return v.str }
func (v *String) Set(j *jason.Value) {
	if n, err := j.String(); err == nil {
		v.str = n
	} else {
		v.str = "N/A"
	}
}
func (v *String) SetNA() {
	v.str = "N/A"
}

// GCPauses represents GC pauses data.
//
// It uses memstat.PauseNS circular buffer, but lacks
// NumGC information, so we don't know what the start
// and the end. It's enough for most stats, though.
type GCPauses struct {
	pauses [256]uint64
	hist   *Histogram
	na     bool
}

func (v *GCPauses) Kind() VarKind { return KindGCPauses }
func (v *GCPauses) String() string {
	if v.na {
		return "N/A"
	}
	// return Mean by default
	return fmt.Sprintf("%v", round(time.Duration(v.hist.Mean())))
}
func (v *GCPauses) Value() int {
	if v.hist == nil {
		return 0
	}
	// return Mean by default
	return int(v.hist.Mean())
}
func (v *GCPauses) Set(j *jason.Value) {
	v.na = false
	v.pauses = [256]uint64{}
	if arr, err := j.Array(); err == nil {
		for i := 0; i < len(arr); i++ {
			p, _ := arr[i].Int64()
			v.pauses[i] = uint64(p)
		}
	}

	// we need histogram object
	// to access mean() method
	// number of bins doesn't matter here
	v.hist = v.Histogram(1)
}
func (v *GCPauses) Histogram(bins int) *Histogram {
	hist := NewHistogram(bins)
	for i := 0; i < 256; i++ {
		// we ignore zeros, since
		// its never the case, but
		// we have zeros on the very beginning
		if v.pauses[i] > 0 {
			hist.Add(v.pauses[i])
		}
	}
	v.hist = hist
	return v.hist
}
func (v *GCPauses) SetNA() {
	v.na = true
}

// GCIntervals represents GC pauses intervals.
//
// It uses memstat.PauseEnd circular buffer w/
// timestamps.
type GCIntervals struct {
	intervals [256]uint64
	hist      *Histogram
	na        bool
}

func (v *GCIntervals) Kind() VarKind { return KindGCIntervals }
func (v *GCIntervals) String() string {
	if v.na {
		return "N/A"
	}
	// return Mean by default
	return fmt.Sprintf("%v", round(time.Duration(v.hist.Mean())))
}
func (v *GCIntervals) Value() int {
	if v.na {
		return 0
	}
	// return Mean by default
	return int(v.hist.Mean())
}
func (v *GCIntervals) Set(j *jason.Value) {
	v.na = false
	v.intervals = [256]uint64{}
	// as original array contains UNIX timestamps,
	// we want to calculate diffs to previous values (interval)
	// and work with them
	duration := func(a, b int64) uint64 {
		dur := int64(a - b)
		if dur < 0 {
			dur = -dur
		}
		return uint64(dur)
	}
	var prev int64
	if arr, err := j.Array(); err == nil {
		// process first elem
		p, _ := arr[0].Int64()
		plast, _ := arr[255].Int64()
		v.intervals[0] = duration(p, plast)
		prev = p

		for i := 1; i < len(arr); i++ {
			p, _ := arr[i].Int64()
			if p == 0 {
				break
			}
			v.intervals[i] = duration(p, prev)
			prev = p
		}
	}

	// the same as for GCPauses, we need it for mean()
	// and number of bins doesn't matter here
	v.hist = v.Histogram(1)
}
func (v *GCIntervals) Histogram(bins int) *Histogram {
	hist := NewHistogram(bins)

	// we need to skip maximum value here
	// because it's always a diff between last and fist
	// elem in cicrular buffer (we don't know NumGC)
	var max uint64
	for i := 0; i < 256; i++ {
		if v.intervals[i] > max {
			max = v.intervals[i]
		}
	}

	for i := 0; i < 256; i++ {
		// we ignore zeros, since
		// its never the case, but
		// we have zeros on the very beginning
		if v.intervals[i] > 0 && v.intervals[i] < max {
			hist.Add(v.intervals[i])
		}
	}
	v.hist = hist
	return hist
}
func (v *GCIntervals) SetNA() {
	v.na = true
}

// TODO: add boolean, timestamp, types

// ToSlice converts "dot-separated" notation into the "slice of strings".
//
// "dot-separated" notation is a human-readable format, passed via args.
// "slice of strings" is used by Jason library.
//
// Example: "memstats.Alloc" => []string{"memstats", "Alloc"}
// Example: "mem:memstats.Alloc" => []string{"memstats", "Alloc"}
func (v VarName) ToSlice() []string {
	start := strings.IndexRune(string(v), ':') + 1
	slice := DottedFieldsToSliceEscaped(string(v)[start:])
	return slice
}

// Short returns short name, which is typically is the last word in the long names.
func (v VarName) Short() string {
	if v == "" {
		return ""
	}

	slice := v.ToSlice()
	return slice[len(slice)-1]
}

// Long returns long name, without kind: modifier.
func (v VarName) Long() string {
	if v == "" {
		return ""
	}

	start := strings.IndexRune(string(v), ':') + 1
	return string(v)[start:]
}

// Kind returns kind of variable, based on it's name
// modifiers ("mem:") or full names for special cases.
func (v VarName) Kind() VarKind {
	if v.Long() == "memstats.PauseNs" {
		return KindGCPauses
	}
	if v.Long() == "memstats.PauseEnd" {
		return KindGCIntervals
	}

	start := strings.IndexRune(string(v), ':')
	if start == -1 {
		return KindDefault
	}

	switch string(v)[:start] {
	case "mem":
		return KindMemory
	case "duration":
		return KindDuration
	case "str":
		return KindString
	}
	return KindDefault
}

// rate calculates rate per seconds for the duration.
func rate(d time.Duration, precision time.Duration) float64 {
	return float64(precision) / float64(d)
}

// round removes unneeded precision from the String() output for time.Duration.
func round(d time.Duration) time.Duration {
	r := time.Second
	if d < time.Second {
		r = time.Millisecond
	}
	if d < time.Millisecond {
		r = time.Microsecond
	}
	if d < time.Microsecond {
		r = time.Nanosecond
	}
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d = d - m
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}

func DottedFieldsToSliceEscaped(s string) []string {
	rv := make([]string, 0)
	lastSlash := false
	curr := ""
	for _, r := range s {
		// base case, dot not after slash
		if !lastSlash && r == '.' {
			if len(curr) > 0 {
				rv = append(rv, curr)
				curr = ""
			}
			continue
		} else if !lastSlash {
			// any character not after slash
			curr += string(r)
			if r == '\\' {
				lastSlash = true
			} else {
				lastSlash = false
			}
			continue
		} else if r == '\\' {
			// last was slash, and so is this
			lastSlash = false // 2 slashes = 0
			// we already appended a single slash on first
			continue
		} else if r == '.' {
			// we see \. but already appended \ last time
			// replace it with .
			curr = curr[:len(curr)-1] + "."
			lastSlash = false
		} else {
			// \ and any other character, ignore
			curr += string(r)
			lastSlash = false
			continue
		}
	}
	if len(curr) > 0 {
		rv = append(rv, curr)
	}
	return rv
}
