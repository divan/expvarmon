package main

type Histogram struct {
	Bins    []Bin
	Maxbins int
	Total   uint64
}

type Bin struct {
	Value, Count uint64
}

// NewHistogram returns a new Histogram with a maximum of n bins.
//
// There is no "optimal" bin count, but somewhere between 20 and 80 bins
// should be sufficient.
func NewHistogram(n int) *Histogram {
	return &Histogram{
		Bins:    make([]Bin, 0),
		Maxbins: n,
		Total:   0,
	}
}

func (h *Histogram) Add(n uint64) {
	defer h.trim()
	h.Total++
	for i := range h.Bins {
		if h.Bins[i].Value == n {
			h.Bins[i].Count++
			return
		}

		if h.Bins[i].Value > n {
			newbin := Bin{Value: n, Count: 1}
			head := append(make([]Bin, 0), h.Bins[0:i]...)

			head = append(head, newbin)
			tail := h.Bins[i:]
			h.Bins = append(head, tail...)
			return
		}
	}

	h.Bins = append(h.Bins, Bin{Count: 1, Value: n})
}

func (h *Histogram) Quantile(q uint64) int64 {
	count := q * h.Total
	for i := range h.Bins {
		count -= h.Bins[i].Count

		if count <= 0 {
			return int64(h.Bins[i].Value)
		}
	}

	return -1
}

// CDF returns the value of the cumulative distribution function
// at x
func (h *Histogram) CDF(x uint64) uint64 {
	var count uint64
	for i := range h.Bins {
		if h.Bins[i].Value <= x {
			count += h.Bins[i].Count
		}
	}

	return count / h.Total
}

// Mean returns the sample mean of the distribution
func (h *Histogram) Mean() float64 {
	if h.Total == 0 {
		return 0
	}

	sum := 0.0

	for i := range h.Bins {
		sum += float64(h.Bins[i].Value * h.Bins[i].Count)
	}

	return sum / float64(h.Total)
}

// Variance returns the variance of the distribution
func (h *Histogram) Variance() float64 {
	if h.Total == 0 {
		return 0
	}

	sum := 0.0
	mean := h.Mean()

	for i := range h.Bins {
		sum += float64(h.Bins[i].Count * (h.Bins[i].Value - uint64(mean)) * (h.Bins[i].Value - uint64(mean)))
	}

	return sum / float64(h.Total)
}

func (h *Histogram) Count() uint64 {
	return h.Total
}

// trim merges adjacent bins to decrease the bin count to the maximum value
func (h *Histogram) trim() {
	for len(h.Bins) > h.Maxbins {
		// Find closest bins in terms of value
		minDelta := uint64(1e10)
		minDeltaIndex := 0
		for i := range h.Bins {
			if i == 0 {
				continue
			}

			if delta := h.Bins[i].Value - h.Bins[i-1].Value; delta < minDelta {
				minDelta = delta
				minDeltaIndex = i
			}
		}

		// We need to merge bins minDeltaIndex-1 and minDeltaIndex
		totalCount := h.Bins[minDeltaIndex-1].Count + h.Bins[minDeltaIndex].Count
		mergedbin := Bin{
			Value: (h.Bins[minDeltaIndex-1].Value*
				h.Bins[minDeltaIndex-1].Count +
				h.Bins[minDeltaIndex].Value*
					h.Bins[minDeltaIndex].Count) /
				totalCount, // weighted average
			Count: totalCount, // summed heights
		}
		head := append(make([]Bin, 0), h.Bins[0:minDeltaIndex-1]...)
		tail := append([]Bin{mergedbin}, h.Bins[minDeltaIndex+1:]...)
		h.Bins = append(head, tail...)
	}
}

func (h *Histogram) BarchartData() ([]uint64, []int) {
	values := make([]uint64, len(h.Bins))
	counts := make([]int, len(h.Bins))
	for i := 0; i < len(h.Bins); i++ {
		values[i] = h.Bins[i].Value
		counts[i] = int(h.Bins[i].Count)
	}
	return values, counts
}
