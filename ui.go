package main

import (
	"fmt"
	"time"

	"github.com/divan/termui"
)

// UI represents UI renderer.
type UI interface {
	Init(UIData) error
	Close()
	Update(UIData)
}

// TermUI is a termUI implementation of UI interface.
type TermUI struct {
	Title      *termui.Par
	Status     *termui.Par
	Services   *termui.List
	Lists      map[VarName]*termui.List
	Sparkline1 *termui.Sparklines
	Sparkline2 *termui.Sparklines
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
			list.ItemFgColor = colorByKind(name.Kind())
			list.Border.Label = name.Short()
			list.Height = len(data.Services) + 2
			t.Lists[name] = list
		}
	}

	makeSparkline := func(name VarName) *termui.Sparklines {
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
		s.Border.Label = fmt.Sprintf("Monitoring %s", name.Long())
		return s
	}
	t.Sparkline1 = makeSparkline(data.Vars[0])
	if len(data.Vars) > 1 {
		t.Sparkline2 = makeSparkline(data.Vars[1])
	}

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
		termui.NewRow(
			termui.NewCol(6, 0, t.Sparkline1),
			termui.NewCol(6, 0, t.Sparkline2)),
	)

	termui.Body.Align()

	return nil
}

// Update updates UI widgets from UIData.
func (t *TermUI) Update(data UIData) {
	t.Title.Text = fmt.Sprintf("monitoring %d services, press q to quit", len(data.Services))
	t.Status.Text = fmt.Sprintf("Last update: %v", data.LastTimestamp.Format("15:04:05 02/Jan/06"))

	// List with service names
	var services []string
	for _, service := range data.Services {
		services = append(services, StatusLine(service))
	}
	t.Services.Items = services

	// Lists with values
	for _, name := range data.Vars {
		var lines []string
		for _, service := range data.Services {
			lines = append(lines, service.Value(name))
		}
		t.Lists[name].Items = lines
	}

	// Sparklines
	for i, service := range data.Services {
		max := formatMax(service.Max(data.Vars[0]))
		t.Sparkline1.Lines[i].Title = fmt.Sprintf("%s%s", service.Name, max)
		t.Sparkline1.Lines[i].Data = service.Values(data.Vars[0])

		if len(data.Vars) > 1 {
			max = formatMax(service.Max(data.Vars[1]))
			t.Sparkline2.Lines[i].Title = fmt.Sprintf("%s%s", service.Name, max)
			t.Sparkline2.Lines[i].Data = service.Values(data.Vars[1])
		}
	}

	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
	termui.Render(termui.Body)
}

// Close shuts down UI module.
func (t *TermUI) Close() {
	termui.Close()
}

// StatusLine returns status line for service with it's name and status.
func StatusLine(s *Service) string {
	if s.Err != nil {
		return fmt.Sprintf("[ERR] %s failed", s.Name)
	}

	return fmt.Sprintf("[R] %s", s.Name)
}

func colorByKind(kind VarKind) termui.Attribute {
	switch kind {
	case KindMemory:
		return termui.ColorRed | termui.AttrBold
	case KindDuration:
		return termui.ColorYellow | termui.AttrBold
	case KindString:
		return termui.ColorGreen | termui.AttrBold
	default:
		return termui.ColorBlue | termui.AttrBold
	}
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

// DummyUI is an simple console UI mockup, for testing purposes.
type DummyUI struct{}

// Init implements UI.
func (*DummyUI) Init(UIData) error { return nil }

// Close implements UI.
func (*DummyUI) Close() {}

// Update implements UI.
func (*DummyUI) Update(data UIData) {
	if data.Services == nil {
		return
	}
	fmt.Println(time.Now().Format("15:04:05 02/01"))
	for _, service := range data.Services {
		fmt.Printf("%s: ", service.Name)
		if service.Err != nil {
			fmt.Printf("ERROR: %s", service.Err)
			continue
		}

		for _, name := range data.Vars {
			fmt.Printf("%s: %v, ", name.Short(), service.Value(name))
		}

		fmt.Printf("\n")
	}
}
