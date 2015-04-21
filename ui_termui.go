package main

import (
	"fmt"
	"github.com/gizak/termui"
	"github.com/pyk/byten"
	"log"
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
	text := fmt.Sprintf("monitoring %d services, press q to quit", len(data.Services))

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

	ls := termui.NewList()

	for _, service := range data.Services {
		if service.Err != nil {
			ls.Items = append(ls.Items, fmt.Sprintf("[E] %s failed", service.Name))
			continue
		}
		alloc := byten.Size(int64(service.Memstats.Alloc))
		sys := byten.Size(int64(service.Memstats.Sys))
		ls.Items = append(ls.Items, fmt.Sprintf("[R] %s: %s/%s goroutines: %d", service.Name, alloc, sys, service.Goroutines))
	}
	ls.ItemFgColor = termui.ColorYellow
	ls.Border.Label = "Services"
	ls.Height = 10
	ls.Width = 100
	ls.Width = termui.TermWidth()
	ls.Y = 3

	dat := []int{4, 2, 1, 6, 3, 9, 1, 4, 2, 15, 14, 9, 8, 6, 10, 13, 15, 12, 10, 5, 3, 6, 1, 7, 10, 10, 14, 13, 6}
	spl0 := termui.NewSparkline()
	spl0.Data = dat[3:]
	spl0.LineColor = termui.ColorGreen

	spls0 := termui.NewSparklines(spl0)
	spls0.Height = 2
	spls0.Width = 20
	spls0.X = 60
	spls0.Y = 3
	spls0.HasBorder = false

	termui.Render(p, p1, ls, spls0)
}

func (t *TermUI) Close() {
	termui.Close()
}
