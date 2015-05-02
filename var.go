package main

import "strings"

// VarName represents variable name.
//
// It has dot-separated format, like "memstats.Alloc",
// but can be used in different forms, hence it's own type.
//
// It also can have optional "kind:" modifier, like "mem:" or "duration:"
type VarName string

type VarKind int

const (
	KindDefault VarKind = iota
	KindMemory
	KindDuration
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
	slice := strings.FieldsFunc(string(v)[start:], func(r rune) bool { return r == '.' })
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
	}
	return KindDefault
}
