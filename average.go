package main

import (
	"github.com/antonholmquist/jason"
)

func averageJason(array []*jason.Value) float64 {
	var arr []float64
	for _, v := range array {
		val, _ := v.Float64()
		arr = append(arr, val)
	}
	return average(arr)
}

// average calculates average (mean) value for int/float array
// trimming zero values from the right.
//
// The whole array/average thing was added to support memstats.PauseNs
// array, which may be filled with zeroes on very beginning.
// Probably it would be better to use Weighted Moving Average and
// add some advanced arrays avarages support, but it's probably wouldn't
// be used much, but PauseNs will be for sure.
func average(arr []float64) float64 {
	// find rightmost non-zero and trim
	right := len(arr)
	for i := right; i > 0; i-- {
		if arr[i-1] != 0.0 {
			right = i
			break
		}
	}
	trimmed := arr[:right]

	// calculate mean
	var sum float64
	for _, v := range trimmed {
		sum += v
	}
	return sum / float64(len(trimmed))
}
