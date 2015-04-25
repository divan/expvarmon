package main

import "testing"

func TestStack(t *testing.T) {
	size := 10
	s := NewStack(size)

	for i := 0; i < size+5; i++ {
		s.Push(i)
		l := len(s.Values)

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
}
