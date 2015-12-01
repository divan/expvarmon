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

	// single \. escapes the dot
	v = VarName(`bleve.indexes.bench\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 5 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench.bleve" ||
		slice[3] != "index" || slice[4] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// double \\. escapes backslash, not dot
	v = VarName(`bleve.indexes.bench\\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 6 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\" ||
		slice[3] != "bleve" || slice[4] != "index" || slice[5] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// triple \\\. escapes backslash then dot
	v = VarName(`bleve.indexes.bench\\\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 5 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\.bleve" ||
		slice[3] != "index" || slice[4] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// quadruple \\\\. escapes two backslashes, not dot
	v = VarName(`bleve.indexes.bench\\\\.bleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 6 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\\\" ||
		slice[3] != "bleve" || slice[4] != "index" || slice[5] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}

	// unsupported \x passes through unaltered
	v = VarName(`bleve.indexes.bench\xbleve.index.lookup_queue_len`)

	slice = v.ToSlice()
	if len(slice) != 5 || slice[0] != "bleve" || slice[1] != "indexes" || slice[2] != "bench\\xbleve" ||
		slice[3] != "index" || slice[4] != "lookup_queue_len" {
		t.Fatalf("ToSlice failed: %v", slice)
	}
}
