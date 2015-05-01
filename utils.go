package main

import (
	"errors"
	"strings"
)

// ParseVars returns parsed and validated slice of strings with
// variables names that will be used for monitoring.
func ParseVars(def, extra string) ([]string, error) {
	if def == "" && extra == "" {
		return nil, errors.New("no vars specified")
	}

	fields := func(s string) []string {
		return strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
	}

	var ret []string
	ret = append(ret, fields(def)...)
	ret = append(ret, fields(extra)...)
	return ret, nil
}

// dot2slice converts "dot-separated" notation into the
// "slice of strings".
//
// "dot-separated" notation is a human-readable format, passed via args.
// "slice of strings" is used by Jason library.
//
// Example: "memstats.Alloc" => []string{"memstats", "Alloc"}
func dot2slice(name string) []string {
	return strings.FieldsFunc(name, func(r rune) bool { return r == '.' })
}
