package main

import "testing"

func TestPushWithFloatAndIntValue(t *testing.T) {
	s := NewStack()
	s.Push(VarValue(int64(0.0))) // from service.go:guessValue
	s.Push(VarValue(5.0))
	s.Push(VarValue(float64(15.0)))
	if _, ok := s.Max.(float64); !ok {
		t.Fatalf("Expected Max to be float64, but it's not")
	}
	s.Push(VarValue(int64(25.0)))
	if _, ok := s.Max.(int64); !ok {
		t.Fatalf("Expected Max to be int64, but it's not")
	}
}

func TestStack(t *testing.T) {
	size := 10
	s := NewStackWithSize(size)

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

	if s.Front().(int) != 14 {
		t.Fatalf("Front returns wrong value: expecting %d, got %d", 14, s.Front())
	}

	s1 := NewStackWithSize(3)
	s1.Push(true)
	s1.Push(false)
	s1.Push(true)

	ints1 := s1.IntValues()
	if len(ints1) != 3 {
		t.Fatalf("expecting len of to be %d, but got %d", 3, len(ints1))
	}
	if ints1[0] != 1 || ints1[1] != 0 || ints1[2] != 1 {
		t.Fatalf("bool values converted to int incorrectly: %v", ints1)
	}

	s2 := NewStackWithSize(3)
	s2.Push(0.1)
	s2.Push(0.5)
	s2.Push(0.03)

	ints2 := s2.IntValues()
	if len(ints2) != 3 {
		t.Fatalf("expecting len to be %d, but got %d", 3, len(ints2))
	}
	if ints2[0] != 10 || ints2[1] != 50 || ints2[2] != 3 {
		t.Fatalf("float values converted to int incorrectly: %v", ints2)
	}
}
