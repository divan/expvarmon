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
	Meminfo      *termui.List
	MemSparkline *termui.Sparklines
}

func (t *TermUI) Init(data Data) {
	err := termui.Init()
	if err != nil {
		log.Fatal(err)
	}

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
		l := termui.NewList()
		l.ItemFgColor = termui.ColorGreen
		l.Border.Label = "Services"
		return l
	}()
	t.Meminfo = func() *termui.List {
		l := termui.NewList()
		l.ItemFgColor = termui.ColorBlue | termui.AttrBold
		l.Border.Label = "Memory Usage"
		return l
	}()
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
		s.Height = 2*data.Total + 2
		s.HasBorder = true
		s.Border.Label = "Memory Track"
		return s
	}()

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, t.Title),
			termui.NewCol(6, 0, t.Status)),
		termui.NewRow(
			termui.NewCol(3, 0, t.Services),
			termui.NewCol(9, 0, t.Meminfo)),
		termui.NewRow(
			termui.NewCol(12, 0, t.MemSparkline)),
	)

	termui.Body.Align()
}

func (t *TermUI) Update(data Data) {
	t.Title.Text = fmt.Sprintf("monitoring %d services, press q to quit", data.Total)
	t.Status.Text = fmt.Sprintf("Last update: %v", data.LastTimestamp.Format("15:04:05 02/Jan/06"))

	var services []string
	var meminfos []string
	for _, service := range data.Services {
		services = append(services, service.StatusLine())
		meminfos = append(meminfos, service.Meminfo())
	}
	t.Services.Items = services
	t.Services.Height = data.Total + 2

	t.Meminfo.Items = meminfos
	t.Meminfo.Height = data.Total + 2

	// Sparklines
	for i, service := range data.Services {
		t.MemSparkline.Lines[i].Title = service.Name
		t.MemSparkline.Lines[i].Data = service.Values["memory"].Values
	}

	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
	termui.Render(termui.Body)
	/*

		goroutines := termui.NewList()
		goroutines.Y = 3
		goroutines.X = meminfo.X + meminfo.Width
		goroutines.Width = termui.TermWidth() - goroutines.X
		goroutines.Height = total + 2
		goroutines.ItemFgColor = termui.ColorGreen
		goroutines.Border.Label = "Goroutines"

		var totalAlloc int64
		for _, service := range data.Services {
			if service.Err != nil {
				names.Items = append(names.Items, fmt.Sprintf("[ERR] %s failed", service.Name))
				meminfo.Items = append(meminfo.Items, "N/A")
				goroutines.Items = append(goroutines.Items, "N/A")
				continue
			}
			alloc := byten.Size(int64(service.MemStats.Alloc))
			heap := byten.Size(int64(service.MemStats.HeapAlloc))
			totalAlloc += int64(service.MemStats.Alloc)

			name := fmt.Sprintf("[R] %s", service.Name)
			meminfos := fmt.Sprintf("%s/%s", alloc, heap)

			names.Items = append(names.Items, name)
			meminfo.Items = append(meminfo.Items, meminfos)
		}


		data.TotalMemory.Push(int(totalAlloc / 1024))

		spl3 := termui.NewSparkline()
		spl3.Data = data.TotalMemory.Values
		spl3.Height = termui.TermHeight() - 3 - (total + 2) - 3
		spl3.LineColor = termui.ColorYellow

		spls2 := termui.NewSparklines(spl3)
		spls2.Y = 3 + (total + 2)
		spls2.Height = termui.TermHeight() - spls2.Y
		spls2.Width = termui.TermWidth()
		spls2.Border.FgColor = termui.ColorCyan
		spls2.Border.Label = fmt.Sprintf("Total Memory Usage: %s", byten.Size(totalAlloc))

		termui.Render(p, p1, names, meminfo, goroutines, spls2, spls)
	*/
}

func (t *TermUI) Close() {
	termui.Close()
}
