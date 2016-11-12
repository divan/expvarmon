package main

import (
	"fmt"
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
)

// Var represents arbitrary value for variable.
type Var interface {
	Kind() VarKind
	String() string
	Set(*jason.Value)
}

// IntVar represents variable which value can be represented as integer,
// and suitable for displaying with sparklines.
type IntVar interface {
	Var
	Value() int
}

// Number is a type for numeric values, obtained from JSON.
// In JSON it's always float64, so there is no straightforward way
// to separate float from int, so let's keep everything as float.
type Number struct {
	val float64
}

func (v *Number) Kind() VarKind { return KindDefault }
func (v *Number) String() string {
	return fmt.Sprintf("%.02f", v.val)
}
func (v *Number) Set(j *jason.Value) {
	if n, err := j.Float64(); err == nil {
		v.val = n
	} else if n, err := j.Int64(); err == nil {
		v.val = float64(n)
	} else {
		v.val = 0
	}
}

// Value implements IntVar for Number type.
func (v *Number) Value() int {
	return int(v.val)
}

// Memory represents memory information in bytes.
type Memory struct {
	bytes int64
}

func (v *Memory) Kind() VarKind { return KindMemory }
func (v *Memory) String() string {
	return fmt.Sprintf("%s", byten.Size(v.bytes))
}
func (v *Memory) Set(j *jason.Value) {
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

// Duration represents duration data (in ns)
type Duration struct {
	dur time.Duration
}

func (v *Duration) Kind() VarKind { return KindDuration }
func (v *Duration) String() string {
	return fmt.Sprintf("%s", roundDuration(time.Duration(v.dur)))
}

func (v *Duration) Set(j *jason.Value) {
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

// TODO: add boolean, timestamp, gcpauses, gcendtimes types

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
	default:
		return &Number{}
	}
}

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

// Kind returns kind of variable, based on it's name modifiers ("mem:")
func (v VarName) Kind() VarKind {
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

// roundDuration removes unneeded precision from the String() output for time.Duration.
func roundDuration(d time.Duration) time.Duration {
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
