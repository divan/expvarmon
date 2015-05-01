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
}

func TestExpvarsAdvanced(t *testing.T) {
	file, err := os.Open(expvarsAdvTestFile)
	if err != nil {
		t.Fatalf("cannot open test file %v", err)
	}
	defer file.Close()
	/*
		vars, err := ParseExpvar(file)
		if err != nil {
			t.Fatal(err)
		}

		if len(vars.Extra) != 2 {
			t.Fatalf("Got:", vars)
			t.Fatalf("vars should have 2 items, but has %d", len(vars.Extra))
		}

		if int(vars.Extra["goroutines"].(float64)) != 10 {
			t.Logf("Expecting 'goroutines' to be %d, but got %d", 10, vars.Extra["goroutines"])
		}

		counters := vars.Extra["counters"].(map[string]interface{})
		counterA := counters["A"].(float64)
		counterB := counters["B"].(float64)
		if counterA != 123.12 {
			t.Logf("Expecting 'counter.A' to be %f, but got %f", 123.12, counterA)
		}
		if int(counterB) != 245342 {
			t.Logf("Expecting 'counter.B' to be %d, but got %d", 245342, counterB)
		}
	*/
}
