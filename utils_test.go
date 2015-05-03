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

func TestPorts(t *testing.T) {
	arg := "1234,1235"
	ports, err := ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 2 || ports[0] != "1234" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "1234-1237,2000"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 5 || ports[0] != "1234" || ports[4] != "2000" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "localhost:2000-2002,remote:1234-1235"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 5 || ports[0] != "localhost:2000" || ports[4] != "remote:1235" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "localhost:2000-2002,remote:1234-1235,some:weird:1234-123input"
	_, err = ParsePorts(arg)
	if err == nil {
		t.Fatalf("err shouldn't be nil")
	}

	arg = "string,sdasd"
	_, err = ParsePorts(arg)
	if err == nil {
		t.Fatalf("err shouldn't be nil")
	}
}
