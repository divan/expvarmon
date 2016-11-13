package main

import "testing"

func TestRecalcBins(t *testing.T) {
	testBins := func(t *testing.T, w, wantBins, wantWidth int) {
		bins, binWidth := recalcBins(w)
		if bins != wantBins {
			t.Fatalf("Expect bins to be %v, but got %v (width: %v)", wantBins, bins, w)
		}
		if binWidth != wantWidth {
			t.Fatalf("Expect bin width to be %v, but got %v (width: %v)", wantWidth, binWidth, w)
		}
	}

	testBins(t, 10, 2, 5)
	testBins(t, 60, 12, 5)
	testBins(t, 100, 20, 5)
}
