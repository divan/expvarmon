package main

import (
	"fmt"
	"strings"
	"time"

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

// VarValue represents arbitrary value for variable.
type VarValue interface{}

const (
	KindDefault VarKind = iota
	KindMemory
	KindDuration
	KindString
)

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

// Format returns human-readable var value representation.
func Format(v VarValue, kind VarKind) string {
	switch kind {
	case KindMemory:
		if _, ok := v.(int64); !ok {
			break
		}
		return fmt.Sprintf("%s", byten.Size(v.(int64)))
	case KindDuration:
		if _, ok := v.(int64); !ok {
			break
		}
		return fmt.Sprintf("%s", roundDuration(time.Duration(v.(int64))))
	}

	if f, ok := v.(float64); ok {
		return fmt.Sprintf("%.2f", f)
	}

	return fmt.Sprintf("%v", v)
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
