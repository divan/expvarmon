package main

// DefaultSize specifies maximum number of items in stack.
//
// Values should be enough for sparklines on high-res terminals
// with minimal font size.
const DefaultSize = 1200

// Stack is a limited FIFO for holding sparkline values.
type Stack struct {
	Values []int
	Len    int
}

// NewStack inits new Stack with default size limit.
func NewStack() *Stack {
	return NewStackWithSize(DefaultSize)
}

// NewStackWithSize inits new Stack with size limit.
func NewStackWithSize(size int) *Stack {
	return &Stack{
		Values: make([]int, size),
		Len:    size,
	}
}

// Push inserts data to stack, preserving constant length.
func (s *Stack) Push(v IntVar) {
	val := v.Value()
	s.Values = append(s.Values, val)
	if len(s.Values) > s.Len {
		// TODO: check if underlying array is growing constantly
		s.Values = s.Values[1:]
	}
}

// IntValues returns stack values explicitly casted to int.
//
// Main case is to use with termui.Sparklines.
func (s *Stack) IntValues() []int {
	return s.Values
}

// TODO: implement trim and resize
