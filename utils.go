package main

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	"github.com/bsiegert/ranges"
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
	var (
		ports []string
		err   error
	)
	// Try simple mode, ports only ("1234-1235,80")
	ports, err = parseRange(s)
	if err == nil {
		return ports, nil
	}

	var ErrParsePorts = fmt.Errorf("cannot parse ports argument")

	// else, try host:ports notation ("localhost:1234-1235,remote:2000,2345")
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ',' })
	for _, field := range fields {
		// split host:ports
		var host, portsRange string
		parts := strings.FieldsFunc(field, func(r rune) bool { return r == ':' })
		if len(parts) == 1 {
			host = "localhost"
		} else if len(parts) == 2 {
			host, portsRange = parts[0], parts[1]
		} else {
			return nil, ErrParsePorts
		}

		pp, err := parseRange(portsRange)
		if err != nil {
			return nil, ErrParsePorts
		}

		for _, p := range pp {
			addr := net.JoinHostPort(host, p)
			ports = append(ports, addr)
		}
	}

	return ports, nil
}

func parseRange(s string) ([]string, error) {
	portsInt, err := ranges.Parse(s)
	if err != nil {
		return nil, err
	}

	var ports []string
	for _, port := range portsInt {
		ports = append(ports, fmt.Sprintf("%d", port))
	}
	return ports, nil
}
