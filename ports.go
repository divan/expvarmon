package main

import (
	"errors"
	"strings"
)

// ParsePorts converts comma-separated ports into strings slice
func ParsePorts(s string) ([]string, error) {
	ports := strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
	if len(ports) == 0 {
		return nil, errors.New("no ports specified")
	}

	return ports, nil
}
