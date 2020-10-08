package main

import (
	"fmt"
	"time"

	"github.com/gizak/termui"
)

// TermUISingle is a termUI implementation of UI interface.
type TermUISingle struct {
	Title      *termui.Paragraph
	Status     *termui.Paragraph
	Sparklines map[VarName]*termui.Sparkline
	Sparkline  *termui.Sparklines
	Pars       []*termui.Paragraph
}

// Init creates widgets, sets sizes and labels.
func (t *TermUISingle) Init(data UIData) error {
	err := termui.Init()
	if err != nil {
		return err
	}

	t.Sparklines = make(map[VarName]*termui.Sparkline)

	t.Title = func() *termui.Paragraph {
		p := termui.NewParagraph("")
		p.Height = 3
		p.TextFgColor = termui.ColorWhite
		p.Border = true
		p.BorderLabel = "Services Monitor"
		p.BorderFg = termui.ColorCyan
		return p
	}()
	t.Status = func() *termui.Paragraph {
		p := termui.NewParagraph("")
		p.Height = 3
		p.TextFgColor = termui.ColorWhite
		p.Border = true
		p.BorderLabel = "Status"
		p.BorderFg = termui.ColorCyan
		return p
	}()

	t.Pars = make([]*termui.Paragraph, len(data.Vars))
	for i, name := range data.Vars {
		par := termui.NewParagraph("")
		par.TextFgColor = colorByKind(name.Kind())
		par.Border = true
		par.BorderLabel = name.Short()
		par.BorderLabelFg = termui.ColorGreen
		par.Height = 3
		t.Pars[i] = par
	}

	sparklines := make([]termui.Sparkline, 0, len(data.Vars))
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
		s.Border = true
		s.BorderLabel = "Monitoring"
		return s
	}()

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
		spl := &t.Sparkline.Lines[i]

		max := formatMax(service.Max(name))
		spl.Title = fmt.Sprintf("%s: %v%s", name.Long(), service.Value(name), max)
		spl.TitleColor = colorByKind(name.Kind())
		spl.LineColor = colorByKind(name.Kind())

		if name.Kind() == KindString {
			continue
		}
		spl.Data = service.Values(name)
	}

	t.Relayout()

	widgets := []termui.Bufferer{t.Title, t.Status, t.Sparkline}
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
	t.Sparkline.Width = tw
	t.Sparkline.Height = h
	t.Sparkline.Y = th - h
}

func formatMax(max interface{}) string {
	var str string
	if max != nil {
		str = fmt.Sprintf(" (max: %v)", max)
	}
	return str
}
