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
}
