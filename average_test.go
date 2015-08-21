package main

import (
	"testing"
)

func TestAverage(t *testing.T) {
	avg := average(samplesPartial)
	want := 621090.75
	if avg != want {
		t.Fatalf("Average must be %v, but got %v", want, avg)
	}
}

var samplesPartial = []float64{507472, 433979, 610916, 931996, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
