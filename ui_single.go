package main

import (
	"fmt"
	"time"

	"gopkg.in/gizak/termui.v1"
)

// TermUISingle is a termUI implementation of UI interface.
type TermUISingle struct {
	Title      *termui.Par
	Status     *termui.Par
	Sparklines map[VarName]*termui.Sparkline
	Sparkline  *termui.Sparklines
	Pars       []*termui.Par
	BarChart   *termui.BarChart
	GCStats    *termui.Par
	BarChart2  *termui.BarChart

	bins int // histograms' bins count
}

// Init creates widgets, sets sizes and labels.
func (t *TermUISingle) Init(data UIData) error {
	err := termui.Init()
	if err != nil {
		return err
	}

	t.Sparklines = make(map[VarName]*termui.Sparkline)

	termui.UseTheme("helloworld")

	t.Title = func() *termui.Par {
		p := termui.NewPar("")
		p.Height = 3
		p.TextFgColor = termui.ColorWhite
		p.Border.Label = "Services Monitor"
		p.Border.FgColor = termui.ColorCyan
		return p
	}()
	t.Status = func() *termui.Par {
		p := termui.NewPar("")
		p.Height = 3
		p.TextFgColor = termui.ColorWhite
		p.Border.Label = "Status"
		p.Border.FgColor = termui.ColorCyan
		return p
	}()

	t.Pars = make([]*termui.Par, len(data.Vars))
	for i, name := range data.Vars {
		par := termui.NewPar("")
		par.TextFgColor = colorByKind(name.Kind())
		par.Border.Label = name.Short()
		par.Border.LabelFgColor = termui.ColorGreen
		par.Height = 3
		t.Pars[i] = par
	}

	var sparklines []termui.Sparkline
	for _, name := range data.Vars {
		spl := termui.NewSparkline()
		spl.Height = 1
		spl.TitleColor = colorByKind(name.Kind())
		spl.LineColor = colorByKind(name.Kind())
		spl.Title = name.Long()
		sparklines = append(sparklines, spl)
	}

	t.Sparkline = func() *termui.Sparklines {
		s := termui.NewSparklines(sparklines...)
		s.Height = 2*len(sparklines) + 2
		s.HasBorder = true
		s.Border.Label = fmt.Sprintf("Monitoring")
		return s
	}()

	if data.HasGCPauses {
		t.BarChart = func() *termui.BarChart {
			bc := termui.NewBarChart()
			bc.Border.Label = "Bar Chart"
			bc.TextColor = termui.ColorGreen
			bc.BarColor = termui.ColorGreen
			bc.NumColor = termui.ColorBlack
			return bc
		}()
		t.GCStats = func() *termui.Par {
			p := termui.NewPar("")
			p.Height = 4
			p.Width = len("Max: 123ms") // example
			p.HasBorder = false
			p.TextFgColor = termui.ColorGreen
			return p
		}()
	}

	if data.HasGCIntervals {
		t.BarChart2 = func() *termui.BarChart {
			bc := termui.NewBarChart()
			bc.Border.Label = "Bar Chart"
			bc.TextColor = termui.ColorGreen
			bc.BarColor = termui.ColorGreen
			bc.NumColor = termui.ColorBlack
			return bc
		}()
	}

	t.Relayout()

	return nil
}

// Update updates UI widgets from UIData.
func (t *TermUISingle) Update(data UIData) {
	// single mode assumes we have one service only to monitor
	service := data.Services[0]

	t.Title.Text = fmt.Sprintf("monitoring %s every %v, press q to quit", service.Name, *interval)
	t.Status.Text = fmt.Sprintf("Last update: %v", data.LastTimestamp.Format(time.Stamp))

	// Pars
	for i, name := range data.Vars {
		t.Pars[i].Text = service.Value(name)
	}

	// Sparklines
	for i, name := range data.Vars {
		v, ok := service.Vars[name].(IntVar)
		if !ok {
			continue
		}
		data.SparklineData[0].Stacks[name].Push(v)
		data.SparklineData[0].Stats[name].Update(v)

		spl := &t.Sparkline.Lines[i]

		max := data.SparklineData[0].Stats[name].Max().String()
		spl.Title = fmt.Sprintf("%s: %v (max: %v)", name.Long(), service.Value(name), max)
		spl.TitleColor = colorByKind(name.Kind())
		spl.LineColor = colorByKind(name.Kind())

		spl.Data = data.SparklineData[0].Stacks[name].Values()
	}

	// BarChart
	if data.HasGCPauses {
		var gcpauses *GCPauses
		for _, v := range service.Vars {
			if v.Kind() == KindGCPauses {
				gcpauses = v.(*GCPauses)
				break
			}
		}
		hist := gcpauses.Histogram(t.bins)
		values, counts := hist.BarchartData()
		vals := make([]int, 0, len(counts))
		labels := make([]string, 0, len(counts))
		for i := 0; i < len(counts); i++ {
			vals = append(vals, int(counts[i]))
			d := round(time.Duration(values[i]))
			labels = append(labels, d.String())
		}
		t.BarChart.Data = vals
		t.BarChart.DataLabels = labels
		t.BarChart.Border.Label = "GC Pauses (last 256)"

		t.GCStats.Text = fmt.Sprintf("Min: %v\nAvg: %v\n95p: %v\nMax: %v",
			round(time.Duration(hist.Min())),
			round(time.Duration(hist.Mean())),
			round(time.Duration(hist.Quantile(0.95))),
			round(time.Duration(hist.Max())),
		)
	}

	if data.HasGCIntervals {
		var gcintervals *GCIntervals
		for _, v := range service.Vars {
			if v.Kind() == KindGCIntervals {
				gcintervals = v.(*GCIntervals)
				break
			}
		}
		hist := gcintervals.Histogram(t.bins)
		values, counts := hist.BarchartData()
		vals := make([]int, 0, len(counts))
		labels := make([]string, 0, len(counts))
		for i := 0; i < len(counts); i++ {
			vals = append(vals, int(counts[i]))
			d := round(time.Duration(values[i]))
			labels = append(labels, d.String())
		}
		t.BarChart2.Data = vals
		t.BarChart2.DataLabels = labels
		t.BarChart2.Border.Label = "Intervals between GC (last 256)"
	}

	t.Relayout()

	var widgets []termui.Bufferer
	widgets = append(widgets, t.Title, t.Status, t.Sparkline)
	if data.HasGCPauses {
		widgets = append(widgets, t.BarChart, t.GCStats)
	}
	if data.HasGCIntervals {
		widgets = append(widgets, t.BarChart2)
	}
	for _, par := range t.Pars {
		widgets = append(widgets, par)
	}
	termui.Render(widgets...)
}

// Close shuts down UI module.
func (t *TermUISingle) Close() {
	termui.Close()
}

// Relayout recalculates widgets sizes and coords.
func (t *TermUISingle) Relayout() {
	tw, th := termui.TermWidth(), termui.TermHeight()
	h := th

	// First row: Title and Status pars
	firstRowH := 3
	t.Title.Height = firstRowH
	t.Title.Width = tw / 2
	if tw%2 == 1 {
		t.Title.Width++
	}
	t.Status.Height = firstRowH
	t.Status.Width = tw / 2
	t.Status.X = t.Title.X + t.Title.Width
	h -= firstRowH

	// Second row: lists
	secondRowH := 3
	num := len(t.Pars)
	parW := tw / num
	for i, par := range t.Pars {
		par.Y = th - h
		par.Width = parW
		par.Height = secondRowH
		par.X = i * parW
	}
	if num*parW < tw {
		t.Pars[num-1].Width = tw - ((num - 1) * parW)
	}
	h -= secondRowH

	// Third row: Sparklines
	calcHeight := len(t.Sparkline.Lines) * 2
	if calcHeight > (h / 2) {
		calcHeight = h / 2
	}

	t.Sparkline.Width = tw
	t.Sparkline.Height = calcHeight
	t.Sparkline.Y = th - h

	// Fourth row: Barcharts
	var barchartWidth, charts int
	if t.BarChart != nil {
		charts++
	}
	if t.BarChart2 != nil {
		charts++
	}

	if charts > 0 {
		barchartWidth = tw / charts
		bins, binWidth := recalcBins(barchartWidth)
		t.bins = bins

		if t.BarChart != nil {
			t.BarChart.Width = barchartWidth
			t.BarChart.Height = h - calcHeight
			t.BarChart.Y = th - t.BarChart.Height
			t.BarChart.BarWidth = binWidth

			t.GCStats.X = barchartWidth - t.GCStats.Width - 1
			t.GCStats.Y = t.BarChart.Y + 1
		}

		if t.BarChart2 != nil {
			t.BarChart2.Width = barchartWidth
			t.BarChart2.X = 0 + barchartWidth
			t.BarChart2.Height = h - calcHeight
			t.BarChart2.Y = th - t.BarChart.Height
			t.BarChart2.BarWidth = binWidth
		}
	}
}

// recalcBins attempts to select optimal value for the number
// of bins for histograms.
//
// Optimal range is 10-30, but we must try to keep bins' width
// no less then 5 to fit the labels ("123ms"). Hence some heuristics.
//
// Should be called on resize or creation.
func recalcBins(tw int) (int, int) {
	var (
		bins, w  int
		minWidth = 5
		minBins  = 10
		maxBins  = 30
	)
	w = minWidth

	tryWidth := func(w int) int {
		return tw / w
	}

	bins = tryWidth(w)
	for bins > maxBins {
		w++
		bins = tryWidth(w)
	}

	for bins < minBins && w > minWidth {
		w--
		bins = tryWidth(w)
	}

	return bins, w
}
