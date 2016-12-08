package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

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

func TestExtractUrlAndPorts(t *testing.T) {
	var rawurl, ports string
	rawurl, ports = extractURLAndPorts("40000-40002")
	if rawurl != "http://localhost" || ports != "40000-40002" {
		t.Fatalf("extract url and ports failed: %v, %v", rawurl, ports)
	}

	rawurl, ports = extractURLAndPorts("https://example.com:1234")
	if rawurl != "https://example.com" || ports != "1234" {
		t.Fatalf("extract url and ports failed: %v, %v", rawurl, ports)
	}

	rawurl, ports = extractURLAndPorts("http://user:passwd@example.com:1234-1256")
	if rawurl != "http://user:passwd@example.com" || ports != "1234-1256" {
		t.Fatalf("extract url and ports failed: %v, %v", rawurl, ports)
	}

	rawurl, ports = extractURLAndPorts("https://example.com:1234-1256/_endpoint")
	if rawurl != "https://example.com/_endpoint" || ports != "1234-1256" {
		t.Fatalf("extract url and ports failed: %v, %v", rawurl, ports)
	}
}

func TestPorts(t *testing.T) {
	arg := "1234,1235"
	ports, err := ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 2 || ports[0].Host != "localhost:1234" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "1234-1237,2000"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 5 || ports[0].Host != "localhost:1234" || ports[4].Host != "localhost:2000" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	arg = "40000-40002,localhost:2000-2002,remote:1234-1235,https://example.com:1234-1236"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 11 ||
		ports[0].Host != "localhost:40000" ||
		ports[3].Host != "localhost:2000" ||
		ports[7].Host != "remote:1235" ||
		ports[7].Path != "/debug/vars" ||
		ports[10].Host != "example.com:1236" ||
		ports[10].Scheme != "https" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}

	// Test Auth
	arg = "http://user:pass@localhost:2000-2002"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	pass, isSet := ports[0].User.Password()
	if len(ports) != 3 ||
		ports[0].User.Username() != "user" ||
		pass != "pass" || !isSet ||
		ports[0].Scheme != "http" {
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

	// Test endpoints
	arg = "localhost:2000,https://example.com:1234/_custom_expvars"
	ports, err = ParsePorts(arg)
	if err != nil {
		t.Fatal(err)
	}
	if ports[0].Path != "/debug/vars" || ports[1].Path != "/_custom_expvars" {
		t.Fatalf("ParsePorts returns wrong data: %v", ports)
	}
}

func TestVarsFlagStringRep(t *testing.T) {
	testTable := []struct {
		input, expected string
	}{
		{"test", "test"},
		{"test ", "test"},
		{" test ", "test"},
		{", test ,", "test"},
		{",test", "test"},
		{"test,", "test"},
	}
	for _, tc := range testTable {
		t.Run(tc.input, func(t *testing.T) {
			v := VarsFlag{}
			v.Set(tc.input)
			if r := v.String(); r != tc.expected {
				t.Errorf("got [%s], want [%s]", r, tc.expected)

			}
		})
	}
}

func TestVarsFlagMultiInputStringRep(t *testing.T) {
	testTable := []struct {
		input    []string
		expected string
	}{
		{[]string{""}, ""},
		{[]string{"test1", "test2"}, "test1,test2"},
		{[]string{"test1 ", "test2"}, "test1,test2"},
		{[]string{"test1 ", " test2"}, "test1,test2"},
		{[]string{" test1 ", " test2 "}, "test1,test2"},
		{[]string{",test1,", " test2 "}, "test1,test2"},
		{[]string{",test1,", ",test2,"}, "test1,test2"},
	}
	for _, tc := range testTable {
		name := strings.Join(tc.input, ",")
		t.Run(name, func(t *testing.T) {
			v := VarsFlag{}
			for _, s := range tc.input {
				v.Set(s)
			}
			if r := v.String(); r != tc.expected {
				t.Errorf("got [%v], want [%v]", r, tc.expected)
			}
		})
	}
}

// This test is for uniqVarNameSlice function
func TestUniqVarNameSlice(t *testing.T) {
	testTable := []struct {
		input, expected []VarName
	}{
		{[]VarName{""}, []VarName{""}},
		{[]VarName{"test1", "test2"}, []VarName{"test1", "test2"}},
		{[]VarName{"test1", "test1"}, []VarName{"test1"}},
		{[]VarName{"test1", "test1", ""}, []VarName{"test1", ""}},
	}
	for i, tc := range testTable {
		name := fmt.Sprintf("Case %d", i)
		t.Run(name, func(t *testing.T) {
			if r := uniqVarNameSlice(tc.input); !reflect.DeepEqual(r, tc.expected) {
				t.Errorf("got [%v], want [%v]", r, tc.expected)
			}
		})
	}
}

func TestVarNamesErrors(t *testing.T) {
	testTable := []struct {
		input    []string
		expected error
	}{
		{[]string{""}, ErrNoVarsSpecified},
		{[]string{}, ErrNoVarsSpecified},
		{[]string{"test1"}, nil},
		{[]string{"test1", "test1"}, nil},
	}
	for i, tc := range testTable {
		name := fmt.Sprintf("Case %d", i)
		t.Run(name, func(t *testing.T) {
			v := VarsFlag{}
			for _, s := range tc.input {
				v.Set(s)
			}

			if _, err := v.VarNames(); err != tc.expected {
				t.Errorf("got [%v], want [%v]", err, tc.expected)
			}
		})
	}

}

// This test is for VarNames method
func TestUniqueVarNames(t *testing.T) {
	testTable := []struct {
		input    []string
		expected []VarName
	}{
		{[]string{"test1"}, []VarName{"test1"}},
		{[]string{"test1", "test2"}, []VarName{"test1", "test2"}},
		{[]string{"test1", "test1"}, []VarName{"test1"}},
		{[]string{"test1,test2"}, []VarName{"test1", "test2"}},
	}
	for _, tc := range testTable {
		name := strings.Join(tc.input, ",")
		t.Run(name, func(t *testing.T) {
			v := VarsFlag{}
			for _, s := range tc.input {
				v.Set(s)
			}
			if r, _ := v.VarNames(); !reflect.DeepEqual(r, tc.expected) {
				t.Errorf("got [%#v], want [%#v]", r, tc.expected)
			}
		})
	}
}
