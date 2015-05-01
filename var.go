package main

import "strings"

// VarName represents variable name.
//
// It has dot-separated format, like "memstats.Alloc",
// but can be used in different forms, hence it's own type.
type VarName string

// ToSlice converts "dot-separated" notation into the "slice of strings".
//
// "dot-separated" notation is a human-readable format, passed via args.
// "slice of strings" is used by Jason library.
//
// Example: "memstats.Alloc" => []string{"memstats", "Alloc"}
func (v VarName) ToSlice() []string {
	return strings.FieldsFunc(string(v), func(r rune) bool { return r == '.' })
}

// Short returns short name, which is typically is the last word in the long names.
func (v VarName) Short() string {
	if v == "" {
		return ""
	}

	slice := v.ToSlice()
	return slice[len(slice)-1]
}
