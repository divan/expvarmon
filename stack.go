package main

// DefaultSize specifies maximum number of items in stack.
//
// Values should be enough for sparklines on high-res terminals
// with minimal font size.
const DefaultSize = 1200

// Stack is a limited FIFO for holding sparkline values.
type Stack struct {
	values []int
	len    int
}

// NewStack inits new Stack with default size limit.
func NewStack() *Stack {
	return NewStackWithSize(DefaultSize)
}

// NewStackWithSize inits new Stack with size limit.
func NewStackWithSize(size int) *Stack {
	return &Stack{
		values: make([]int, size),
		len:    size,
	}
}

// Push inserts data to stack, preserving constant length.
func (s *Stack) Push(v IntVar) {
	val := v.Value()
	s.values = append(s.values, val)
	if len(s.values) > s.len {
		// TODO: check if underlying array is growing constantly
		s.values = s.values[1:]
	}
}

// Values returns stack values explicitly casted to int.
//
// Main case is to use with termui.Sparklines.
func (s *Stack) Values() []int {
	return s.values
}

// TODO: implement trim and resize
