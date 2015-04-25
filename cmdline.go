package main

import (
	"path/filepath"
)

// BaseCommand returns cleaned command name from Cmdline array.
//
// I.e. "./some.service/binary.name -arg 1 -arg" will be "binary.name".
func BaseCommand(cmdline []string) string {
	if len(cmdline) == 0 {
		return ""
	}
	return filepath.Base(cmdline[0])
}
