package main

import (
	"os"
	"testing"
)

const (
	expvarsTestFile    = "./expvars.json"
	expvarsAdvTestFile = "./expvars_advanced.json"
)

func TestExpvars(t *testing.T) {
	file, err := os.Open(expvarsTestFile)
	if err != nil {
		t.Fatalf("cannot open test file %v", err)
	}
	defer file.Close()

	expvar, err := ParseExpvar(file)
	if err != nil {
		t.Fatal(err)
	}

	cmdline, err := expvar.GetStringArray("cmdline")
	if err != nil {
		t.Fatal(err)
	}
	if len(cmdline) != 3 {
		t.Fatalf("Cmdline should have 3 items, but has %d", len(cmdline))
	}

	alloc, err := expvar.GetInt64(VarName("memstats.Alloc").ToSlice()...)
	if err != nil {
		t.Fatal(err)
	}
	if alloc == 0 {
		t.Fatalf("Alloc should be greater than 0")
	}

	pauses, err := expvar.GetInt64Array(VarName("memstats.PauseNs").ToSlice()...)
	if err != nil {
		t.Fatal(err)
	}
	if len(pauses) == 0 {
		t.Fatalf("Pauses length should be greater than 0")
	}
}

func TestExpvarsAdvanced(t *testing.T) {
	file, err := os.Open(expvarsAdvTestFile)
	if err != nil {
		t.Fatalf("cannot open test file %v", err)
	}
	defer file.Close()

	expvar, err := ParseExpvar(file)
	if err != nil {
		t.Fatal(err)
	}

	goroutines, err := expvar.GetInt64("goroutines")
	if err != nil {
		t.Fatal(err)
	}
	if goroutines != 10 {
		t.Fatalf("Expecting 'goroutines' to be %d, but got %d", 10, goroutines)
	}

	counterA, err := expvar.GetFloat64(VarName("counters.A").ToSlice()...)
	if err != nil {
		t.Fatal(err)
	}
	if counterA != 123.12 {
		t.Fatalf("Expecting 'counters.A' to be %f, but got %f", 123.12, counterA)
	}
}
