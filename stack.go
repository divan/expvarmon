package main

// DefaultSize specifies maximum number of items in stack.
//
// Values should be enough for sparklines on high-res terminals
// with minimal font size.
const DefaultSize = 1200

// Stack is a limited FIFO for holding sparkline values.
type Stack struct {
	Values []interface{}
	Len    int
}

// NewStack inits new Stack with default size limit.
func NewStack() *Stack {
	return NewStackWithSize(DefaultSize)
}

// NewStackWithSize inits new Stack with size limit.
func NewStackWithSize(size int) *Stack {
	return &Stack{
		Values: make([]interface{}, size),
		Len:    size,
	}
}

// Push inserts data to stack, preserving constant length.
func (s *Stack) Push(val interface{}) {
	s.Values = append(s.Values, val)
	if len(s.Values) > s.Len {
		s.Values = s.Values[1:]
	}
}

// Front returns front value.
func (s *Stack) Front() interface{} {
	if len(s.Values) == 0 {
		return nil
	}
	return s.Values[len(s.Values)-1]
}

// IntValues returns stack values explicitly casted to int.
//
// Main case is to use with termui.Sparklines.
func (s *Stack) IntValues() []int {
	ret := make([]int, s.Len)
	for i, v := range s.Values {
		n, ok := v.(int64)
		if ok {
			ret[i] = int(n)
			continue
		}

		f, ok := v.(float64)
		if ok {
			// 12.34 (float) -> 1234 (int)
			ret[i] = int(f * 100)
			continue
		}

		b, ok := v.(bool)
		if ok {
			// false => 0, true = 1
			if b {
				ret[i] = 1
			} else {
				ret[i] = 0
			}
		}
	}
	return ret
}
