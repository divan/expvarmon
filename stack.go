package main

type Stack struct {
	Values []int
	Len    int
}

func NewStack(size int) *Stack {
	return &Stack{
		Values: make([]int, size),
		Len:    size,
	}
}

func (s *Stack) Push(val int) {
	s.Values = append(s.Values, val)
	if len(s.Values) > s.Len {
		s.Values = s.Values[1:]
	}
}
