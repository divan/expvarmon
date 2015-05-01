package main

import (
	"errors"
	"strings"
)

// ParseVars returns parsed and validated slice of strings with
// variables names that will be used for monitoring.
func ParseVars(def, extra string) ([]VarName, error) {
	if def == "" && extra == "" {
		return nil, errors.New("no vars specified")
	}

	fields := func(s string) []VarName {
		ss := strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
		ret := []VarName{}
		for _, str := range ss {
			ret = append(ret, VarName(str))
		}
		return ret
	}

	var ret []VarName
	ret = append(ret, fields(def)...)
	ret = append(ret, fields(extra)...)
	return ret, nil
}
