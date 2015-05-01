package main

import (
	"fmt"
	"log"

	"github.com/divan/termui"
)

// TermUI is a termUI implementation of UI interface.
type TermUI struct {
	Title        *termui.Par
	Status       *termui.Par
	Services     *termui.List
	Values       map[string]*termui.List
	MemSparkline *termui.Sparklines
}

func (t *TermUI) Init(data UIData) {
	err := termui.Init()
	if err != nil {
		log.Fatal(err)
	}

	t.Values = make(map[string]*termui.List)

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
		_, ok := t.Values[name]
		if !ok {
			list := termui.NewList()
			list.ItemFgColor = termui.ColorBlue | termui.AttrBold
			list.Border.Label = name
			list.Height = len(data.Services) + 2
			t.Values[name] = list
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

	col := termui.NewCol(2, 0, t.Services)
	col1 := termui.NewCol(3, 0, t.Values[data.Vars[0]])
	col2 := termui.NewCol(3, 0, t.Values[data.Vars[1]])
	col3 := termui.NewCol(2, 0, t.Values[data.Vars[2]])
	col4 := termui.NewCol(2, 0, t.Values[data.Vars[3]])
	valuesRow := termui.NewRow(col, col1, col2, col3, col4)
	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, t.Title),
			termui.NewCol(6, 0, t.Status)),
		valuesRow,
		termui.NewRow(
			termui.NewCol(12, 0, t.MemSparkline)),
	)

	termui.Body.Align()
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
		t.Values[name].Items = lines
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

func (t *TermUI) Close() {
	termui.Close()
}
