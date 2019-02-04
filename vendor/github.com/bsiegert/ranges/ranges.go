/*-
 * Copyright (c) 2011
 *	Benny Siegert <bsiegert@gmail.com>
 *
 * Provided that these terms and disclaimer and all copyright notices
 * are retained or reproduced in an accompanying document, permission
 * is granted to deal in this work without restriction, including un-
 * limited rights to use, publicly perform, distribute, sell, modify,
 * merge, give away, or sublicence.
 *
 * This work is provided "AS IS" and WITHOUT WARRANTY of any kind, to
 * the utmost extent permitted by applicable law, neither express nor
 * implied; without malicious intent or gross negligence. In no event
 * may a licensor, author or contributor be held liable for indirect,
 * direct, other damage, loss, or other issues arising in any way out
 * of dealing in the work, even if advised of the possibility of such
 * damage or existence of a defect, except proven that it results out
 * of said person's immediate fault when using the work as intended.
 */

// Package ranges contains tools for working with integer ranges.
// 
// An "integer range" allows to give a set of numbers as a string,
// which can be parsed by a call to Parse. The result can be obtained
// as a slice of integers by calling Expand or be tested against with
// Contains.
package ranges

import (
	"fmt"
	"strconv"
	"strings"
)

// An IntRange is a single component of an integer range expression.
type IntRange struct {
	Lo, Hi int
}

// Contains returns true if ir contains value.
func (ir *IntRange) Contains(value int) bool {
	return value >= ir.Lo && value <= ir.Hi
}

// Expand returns a sorted slice of integers that contains all the numbers
// in ir.
func (ir *IntRange) Expand() []int {
	e := make([]int, 0, ir.Hi-ir.Lo+1)
	for i := ir.Lo; i <= ir.Hi; i++ {
		e = append(e, i)
	}
	return e
}

// Clean exchanges the upper and lower bound if the upper bound is
// smaller than the lower one.
func (ir *IntRange) Clean() {
	if ir.Hi < ir.Lo {
		ir.Hi, ir.Lo = ir.Lo, ir.Hi
	}
}

// IntRanges is a slice of multiple integer ranges, allowing the
// expression of non-contiguous ranges (for example "1,3-4").
type IntRanges []IntRange

// Contains returns true if ir contains value.
func (ir *IntRanges) Contains(value int) bool {
	for _, r := range *ir {
		if r.Contains(value) {
			return true
		}
	}
	return false
}

// Expand returns a slice of integers that contains all the numbers in ir.
// If ir has been cleaned by calling Clean, the slice will be sorted.
func (ir *IntRanges) Expand() []int {
	// This guess for the length is as good as any other.
	e := make([]int, 0, 2*len(*ir))

	for _, r := range *ir {
		e = append(e, r.Expand()...)
	}
	return e
}

// Len returns the number of distinct ranges in ir.
func (ir *IntRanges) Len() int {
	return len(*ir)
}

// Less returns true if the lower bound of the i-th element is smaller
// than the one of the j-th element.
func (ir *IntRanges) Less(i, j int) bool {
	return (*ir)[i].Lo < (*ir)[j].Lo
}

// Swap swaps the i-th and the j-th element.
func (ir *IntRanges) Swap(i, j int) {
	(*ir)[i], (*ir)[j] = (*ir)[j], (*ir)[i]
}

func Parse(r string) ([]int, error) {
	var expanded []int

	for _, item := range strings.Split(r, ",") {
		lohi := strings.Split(item, "-")
		switch len(lohi) {
		case 1:
			v, err := strconv.Atoi(item)
			if err != nil {
				return nil, err
			}
			expanded = append(expanded, v)
		case 2:
			lo, err := strconv.Atoi(lohi[0])
			if err != nil {
				return nil, err
			}
			hi, err := strconv.Atoi(lohi[1])
			if err != nil {
				return nil, err
			}
			for i := lo; i <= hi; i++ {
				expanded = append(expanded, i)
			}
		default:
			return nil, fmt.Errorf("invalid range: %s", item)
		}
	}
	return expanded, nil
}
