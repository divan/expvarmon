package main

import (
	"fmt"
	"log"

	"github.com/gizak/termui"
	"github.com/pyk/byten"
)

// TermUI is a termUI implementation of UI interface.
type TermUI struct {
}

func (t *TermUI) Init() {
	err := termui.Init()
	if err != nil {
		log.Fatal(err)
	}

	termui.UseTheme("helloworld")
}

func (t *TermUI) Update(data Data) {
	total := len(data.Services)
	text := fmt.Sprintf("monitoring %d services, press q to quit", total)

	p := termui.NewPar(text)
	p.Height = 3
	p.Width = termui.TermWidth() / 2
	p.TextFgColor = termui.ColorWhite
	p.Border.Label = "Services Monitor"
	p.Border.FgColor = termui.ColorCyan

	text1 := fmt.Sprintf("Last update: %v", data.LastTimestamp.Format("15:04:05 02/Jan/06"))
	p1 := termui.NewPar(text1)
	p1.Height = 3
	p1.X = p.X + p.Width
	p1.Width = termui.TermWidth() - p1.X
	p1.TextFgColor = termui.ColorWhite
	p1.Border.Label = "Status"
	p1.Border.FgColor = termui.ColorCyan

	names := termui.NewList()
	names.Y = 3
	names.ItemFgColor = termui.ColorYellow
	names.Border.Label = "Services"
	names.Height = total + 2
	names.Width = termui.TermWidth() / 4

	meminfo := termui.NewList()
	meminfo.Y = 3
	meminfo.X = names.X + names.Width
	meminfo.Width = meminfo.X + termui.TermWidth()/3
	meminfo.Height = total + 2
	meminfo.ItemFgColor = termui.ColorBlue
	meminfo.Border.Label = "Memory Usage (Alloc/HeapAlloc)"

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
		alloc := byten.Size(int64(service.Memstats.Alloc))
		heap := byten.Size(int64(service.Memstats.HeapAlloc))
		totalAlloc += int64(service.Memstats.Alloc)

		name := fmt.Sprintf("[R] %s", service.Name)
		meminfos := fmt.Sprintf("%s/%s", alloc, heap)
		goroutine := fmt.Sprintf("%d", service.Goroutines)

		names.Items = append(names.Items, name)
		meminfo.Items = append(meminfo.Items, meminfos)
		goroutines.Items = append(goroutines.Items, goroutine)
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

	termui.Render(p, p1, names, meminfo, goroutines, spls2)
}

func (t *TermUI) Close() {
	termui.Close()
}
