package main

import (
	"errors"
	"path/filepath"
	"strings"
)

// ParseVars returns parsed and validated slice of strings with
// variables names that will be used for monitoring.
func ParseVars(vars string) ([]VarName, error) {
	if vars == "" {
		return nil, errors.New("no vars specified")
	}

	ss := strings.FieldsFunc(vars, func(r rune) bool { return r == ',' })
	var ret []VarName
	for _, s := range ss {
		ret = append(ret, VarName(s))
	}
	return ret, nil
}

// BaseCommand returns cleaned command name from Cmdline array.
//
// I.e. "./some.service/binary.name -arg 1 -arg" will be "binary.name".
func BaseCommand(cmdline []string) string {
	if len(cmdline) == 0 {
		return ""
	}
	return filepath.Base(cmdline[0])
}

// ParsePorts converts comma-separated ports into strings slice
func ParsePorts(s string) ([]string, error) {
	ports := strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
	if len(ports) == 0 {
		return nil, errors.New("no ports specified")
	}

	return ports, nil
}
