package main

// DefaultSize specifies maximum number of items in stack.
//
// Values should be enough for sparklines on high-res terminals
// with minimal font size.
const DefaultSize = 1200

// Stack is a limited FIFO for holding sparkline values.
type Stack struct {
	Values []VarValue
	Len    int
	Max    VarValue
}

// NewStack inits new Stack with default size limit.
func NewStack() *Stack {
	return NewStackWithSize(DefaultSize)
}

// NewStackWithSize inits new Stack with size limit.
func NewStackWithSize(size int) *Stack {
	return &Stack{
		Values: make([]VarValue, size),
		Len:    size,
	}
}

// Push inserts data to stack, preserving constant length.
func (s *Stack) Push(val VarValue) {
	s.Values = append(s.Values, val)
	if len(s.Values) > s.Len {
		s.Values = s.Values[1:]
	}

	if s.Max == nil {
		s.Max = val
		return
	}

	switch val.(type) {
	case int64:
		switch s.Max.(type) {
		case int64:
			if val.(int64) > s.Max.(int64) {
				s.Max = val
			}
		case float64:
			if float64(val.(int64)) > s.Max.(float64) {
				s.Max = val
			}
		}
	case float64:
		switch s.Max.(type) {
		case int64:
			if val.(float64) > float64(s.Max.(int64)) {
				s.Max = val
			}
		case float64:
			if val.(float64) > s.Max.(float64) {
				s.Max = val
			}
		}
	}
}

// Front returns front value.
func (s *Stack) Front() VarValue {
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
