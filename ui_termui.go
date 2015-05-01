package main

import (
	"fmt"

	"github.com/divan/termui"
)

// TermUI is a termUI implementation of UI interface.
type TermUI struct {
	Title        *termui.Par
	Status       *termui.Par
	Services     *termui.List
	Lists        map[VarName]*termui.List
	MemSparkline *termui.Sparklines
}

// Init creates widgets, sets sizes and labels.
func (t *TermUI) Init(data UIData) error {
	err := termui.Init()
	if err != nil {
		return err
	}

	t.Lists = make(map[VarName]*termui.List)

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
	t.Services = func() *termui.List {
		list := termui.NewList()
		list.ItemFgColor = termui.ColorGreen
		list.Border.Label = "Services"
		list.Height = len(data.Services) + 2
		return list
	}()

	for _, name := range data.Vars {
		_, ok := t.Lists[name]
		if !ok {
			list := termui.NewList()
			list.ItemFgColor = termui.ColorBlue | termui.AttrBold
			list.Border.Label = name.Short()
			list.Height = len(data.Services) + 2
			t.Lists[name] = list
		}
	}

	t.MemSparkline = func() *termui.Sparklines {
		var sparklines []termui.Sparkline
		for _, service := range data.Services {
			spl := termui.NewSparkline()
			spl.Height = 1
			spl.LineColor = termui.ColorGreen
			spl.Title = service.Name
			sparklines = append(sparklines, spl)
		}

		s := termui.NewSparklines(sparklines...)
		s.Height = 2*len(data.Services) + 2
		s.HasBorder = true
		s.Border.Label = "Memory Track"
		return s
	}()

	cellW, firstW := calculateCellWidth(len(data.Vars) + 1)
	col := termui.NewCol(firstW, 0, t.Services)
	cols := []*termui.Row{col}
	for _, name := range data.Vars {
		cols = append(cols, termui.NewCol(cellW, 0, t.Lists[name]))
	}
	listsRow := termui.NewRow(cols...)

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, t.Title),
			termui.NewCol(6, 0, t.Status)),
		listsRow,
		termui.NewRow(termui.NewCol(12, 0, t.MemSparkline)),
	)

	termui.Body.Align()

	return nil
}

func (t *TermUI) Update(data UIData) {
	t.Title.Text = fmt.Sprintf("monitoring %d services, press q to quit", len(data.Services))
	t.Status.Text = fmt.Sprintf("Last update: %v", data.LastTimestamp.Format("15:04:05 02/Jan/06"))

	var services []string
	for _, service := range data.Services {
		services = append(services, service.StatusLine())
	}
	t.Services.Items = services

	for _, name := range data.Vars {
		var lines []string
		for _, service := range data.Services {
			lines = append(lines, service.Value(name))
		}
		t.Lists[name].Items = lines
	}

	// Sparklines
	for i, service := range data.Services {
		t.MemSparkline.Lines[i].Title = service.Name
		t.MemSparkline.Lines[i].Data = service.Values("memstats.Alloc")
	}

	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
	termui.Render(termui.Body)
}

// Close shuts down UI module.
func (t *TermUI) Close() {
	termui.Close()
}

// GridSz defines grid size used in TermUI
const GridSz = 12

// calculateCellWidth does some heuristics to calculate optimal cells width
// for all cells, and adjust first width (with service names) if needed.
func calculateCellWidth(num int) (cellW int, firstW int) {
	cellW = GridSz / num
	firstW = cellW + (GridSz - (num * cellW))
	return
}
