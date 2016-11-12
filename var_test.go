package main

import (
	"github.com/antonholmquist/jason"
	"strings"
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

func TestVarNew(t *testing.T) {
	v := NewVar(VarName("mem:field.Subfield"))
	if v.Kind() != KindMemory {
		t.Fatalf("Expect Memory, got %v", v.Kind())
	}
	v = NewVar(VarName("str:field.Subfield"))
	if v.Kind() != KindString {
		t.Fatalf("Expect String, got %v", v.Kind())
	}
	v = NewVar(VarName("duration:field.Subfield"))
	if v.Kind() != KindDuration {
		t.Fatalf("Expect Duration, got %v", v.Kind())
	}
}

func str2val(t *testing.T, s string) *jason.Value {
	val, err := jason.NewValueFromReader(strings.NewReader(s))
	if err != nil {
		t.Fatal(err)
	}
	return val
}

func TestVarNumber(t *testing.T) {
	v := &Number{}
	testNumber := func(t *testing.T, v *Number, json string, intval int, str string) {
		v.Set(str2val(t, json))
		if want := intval; v.Value() != want {
			t.Fatalf("Expect value to be %d, got %d", want, v.Value())
		}
		if want := str; v.String() != want {
			t.Fatalf("Expect value to be %s, got %s", want, v.String())
		}
	}

	testNumber(t, v, "142", 142, "142.00")
	testNumber(t, v, "13.24", 13, "13.24")
	testNumber(t, v, "true", 0, "0.00")
	testNumber(t, v, "\"success\"", 0, "0.00")
}

func TestVarMemory(t *testing.T) {
	v := &Memory{}
	testMemory := func(t *testing.T, v *Memory, json string, intval int, str string) {
		v.Set(str2val(t, json))
		if want := intval; v.Value() != want {
			t.Fatalf("Expect value to be %d, got %d", want, v.Value())
		}
		if want := str; v.String() != want {
			t.Fatalf("Expect value to be %s, got %s", want, v.String())
		}
	}

	testMemory(t, v, "12", 12, "12B")
	testMemory(t, v, "1024", 1024, "1.0KB")
	testMemory(t, v, "1048576", 1048576, "1.0MB")
	testMemory(t, v, "1073741824", 1073741824, "1.0GB")
	testMemory(t, v, "6815744", 6815744, "6.5MB")
	testMemory(t, v, "128849018880", 128849018880, "120GB")
}

func TestVarDuration(t *testing.T) {
	v := &Duration{}
	testDuration := func(t *testing.T, v *Duration, json string, intval int, str string) {
		v.Set(str2val(t, json))
		if want := intval; v.Value() != want {
			t.Fatalf("Expect value to be %d, got %d", want, v.Value())
		}
		if want := str; v.String() != want {
			t.Fatalf("Expect value to be %s, got %s", want, v.String())
		}
	}

	testDuration(t, v, "12", 12, "12ns")
	testDuration(t, v, "1000", 1e3, "1µs")
	testDuration(t, v, "2000", 2*1e3, "2µs")
	testDuration(t, v, "1000000", 1e6, "1ms")
	testDuration(t, v, "2000000", 2*1e6, "2ms")
	testDuration(t, v, "155000000", 155*1e6, "155ms")
	testDuration(t, v, "1000000000", 1e9, "1s")
	testDuration(t, v, "13000000000", 13*1e9, "13s")
	testDuration(t, v, "60000000000", 60*1e9, "1m0s")
	testDuration(t, v, "90000000000", 90*1e9, "1m30s")
	testDuration(t, v, "172800000000000", 48*3600*1e9, "48h0m0s")
	testDuration(t, v, "63072000000000000", 2*365*24*3600*1e9, "17520h0m0s")
}

func TestVarString(t *testing.T) {
	v := &String{}
	v.Set(str2val(t, "\"success\""))
	if want := "success"; v.String() != want {
		t.Fatalf("Expect value to be %s, got %s", want, v.String())
	}
	v.Set(str2val(t, "123"))
	if want := "N/A"; v.String() != want {
		t.Fatalf("Expect value to be %s, got %s", want, v.String())
	}
}
