package main

import (
	"strings"
)

type StringArray []string

func (a *StringArray) Set(s string) error {
	s = strings.TrimSuffix(s, "/debug/vars")
	*a = append(*a, s)
	return nil
}

func (a *StringArray) String() string {
	return strings.Join(*a, ",")
}
