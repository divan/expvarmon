package main

import (
	"os"
	"testing"
)

const expvarsTestFile = "./expvars.json"

func TestExpvars(t *testing.T) {
	file, err := os.Open(expvarsTestFile)
	if err != nil {
		t.Fatalf("cannot open test file %v", err)
	}
	defer file.Close()

	vars, err := ParseExpvar(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(vars.Cmdline) != 3 {
		t.Fatalf("Cmdline should have 3 items, but has %d", len(vars.Cmdline))
	}
}
