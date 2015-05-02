package main

import (
	"testing"
)

func TestVarName(t *testing.T) {
	v := VarName("memstats.Alloc")

	slice := v.ToSlice()
	if len(slice) != 2 || slice[0] != "memstats" || slice[1] != "Alloc" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	short := v.Short()
	if short != "Alloc" {
		t.Fatalf("Expecting Short() to be 'Alloc', but got: %s", short)
	}

	kind := v.Kind()
	if kind != KindDefault {
		t.Fatalf("Expecting kind to be %v, but got: %v", KindDefault, kind)
	}

	v = VarName("mem:memstats.Alloc")

	slice = v.ToSlice()
	if len(slice) != 2 || slice[0] != "memstats" || slice[1] != "Alloc" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	short = v.Short()
	if short != "Alloc" {
		t.Fatalf("Expecting Short() to be 'Alloc', but got: %s", short)
	}

	kind = v.Kind()
	if kind != KindMemory {
		t.Fatalf("Expecting kind to be %v, but got: %v", KindMemory, kind)
	}

	v = VarName("duration:ResponseTimes.API.Users")
	kind = v.Kind()
	if kind != KindDuration {
		t.Fatalf("Expecting kind to be %v, but got: %v", KindDuration, kind)
	}
}
