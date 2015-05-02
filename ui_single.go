package main

import (
	"fmt"

	"github.com/divan/termui"
)

// TermUISingle is a termUI implementation of UI interface.
type TermUISingle struct {
	Title      *termui.Par
	Status     *termui.Par
	Sparklines map[VarName]*termui.Sparkline
	Sparkline  *termui.Sparklines
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

	for _, name := range data.Vars {
		_, ok := t.Sparklines[name]
		if !ok {
			spl := termui.NewSparkline()
			spl.Height = 1
			spl.TitleColor = colorByKind(name.Kind())
			spl.LineColor = colorByKind(name.Kind())
			spl.Title = name.Long()
			t.Sparklines[name] = &spl
		}
	}

	t.Sparkline = func() *termui.Sparklines {
		var spls []termui.Sparkline
		for _, sp := range t.Sparklines {
			spls = append(spls, *sp)
		}
		s := termui.NewSparklines(spls...)
		s.Height = 2*len(spls) + 2
		s.HasBorder = true
		s.Border.Label = fmt.Sprintf("Monitoring")
		return s
	}()

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, t.Title),
			termui.NewCol(6, 0, t.Status)),
		termui.NewRow(termui.NewCol(12, 0, t.Sparkline)),
	)

	termui.Body.Align()

	return nil
}

// Update updates UI widgets from UIData.
func (t *TermUISingle) Update(data UIData) {
	// single mode assumes we have one service only to monitor
	service := data.Services[0]

	t.Title.Text = fmt.Sprintf("monitoring %s, press q to quit", service.Name)
	t.Status.Text = fmt.Sprintf("Last update: %v", data.LastTimestamp.Format("15:04:05 02/Jan/06"))

	// Sparklines
	for i, name := range data.Vars {
		if name.Kind() == KindString {
			continue
		}
		spl := &t.Sparkline.Lines[i]
		spl.Title = fmt.Sprintf("%s: %v", name.Long(), service.Value(name))
		spl.TitleColor = colorByKind(name.Kind())
		spl.LineColor = colorByKind(name.Kind())
		spl.Data = service.Values(name)
	}

	termui.Body.Width = termui.TermWidth()
	termui.Body.Align()
	termui.Render(termui.Body)
}

// Close shuts down UI module.
func (t *TermUISingle) Close() {
	termui.Close()
}
