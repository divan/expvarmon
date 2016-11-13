package main

import "testing"

func TestStack(t *testing.T) {
	size := 10
	s := NewStackWithSize(size)

	for i := 0; i < size+5; i++ {
		s.Push(&Number{float64(i), false})
		l := len(s.values)

		if l < size {
			if l != i+1 {
				t.Fatalf("len is incorrect. expecting %d, got %d", i, l)
			}
		} else {
			if l != size {
				t.Fatalf("len is incorrect. expecting %d, got %d", size, l)
			}
		}
	}

	got := s.Values()[9]
	if got != 14 {
		t.Fatalf("Front returns wrong value: expecting %d, got %d", 14, got)
	}
}
