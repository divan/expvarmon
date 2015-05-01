package main

import "testing"

func TestUtils(t *testing.T) {
	str := "memstats.Alloc,memstats.Sys"

	vars, err := ParseVars(str)
	if err != nil {
		t.Fatalf("Err not nil: %v", err)
	}

	if len(vars) != 2 {
		t.Fatalf("vars should contain 2 elements, but has %d", len(vars))
	}

	str = "memstats.Alloc,memstats.Sys,goroutines,Counter.A"

	vars, err = ParseVars(str)
	if err != nil {
		t.Fatalf("Err not nil: %v", err)
	}

	if len(vars) != 4 {
		t.Fatalf("vars should contain 4 elements, but has %d", len(vars))
	}
}
