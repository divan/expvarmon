package main

import (
	"fmt"
	"time"

	"gopkg.in/gizak/termui.v1"
)

// TermUI is a termUI implementation of UI interface.
type TermUI struct {
	Title      *termui.Par
	Status     *termui.Par
	Services   *termui.List
	Lists      []*termui.List
	Sparkline1 *termui.Sparklines
	Sparkline2 *termui.Sparklines
}

// Init creates widgets, sets sizes and labels.
func (t *TermUI) Init(data UIData) error {
	err := termui.Init()
	if err != nil {
		return err
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
		list := termui.NewList()
		list.ItemFgColor = termui.ColorGreen
		list.Border.LabelFgColor = termui.ColorGreen | termui.AttrBold
		list.Border.Label = "Services"
		list.Height = len(data.Services) + 2
		return list
	}()

	t.Lists = make([]*termui.List, len(data.Vars))
	for i, name := range data.Vars {
		list := termui.NewList()
		list.ItemFgColor = colorByKind(name.Kind())
		list.Border.Label = name.Short()
		list.Border.LabelFgColor = termui.ColorGreen
		if i < 2 {
			list.Border.LabelFgColor = termui.ColorGreen | termui.AttrBold
		}
		list.Height = len(data.Services) + 2
		t.Lists[i] = list
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

	t.Relayout()

	return nil
}

// Update updates UI widgets from UIData.
func (t *TermUI) Update(data UIData) {
	t.Title.Text = fmt.Sprintf("monitoring %d services every %v, press q to quit", len(data.Services), *interval)
	t.Status.Text = fmt.Sprintf("Last update: %v", data.LastTimestamp.Format(time.Stamp))

	// List with service names
	var services []string
	for _, service := range data.Services {
		services = append(services, StatusLine(service))
	}
	t.Services.Items = services

	// Lists with values
	for i, name := range data.Vars {
		var lines []string
		for _, service := range data.Services {
			lines = append(lines, service.Value(name))
		}
		t.Lists[i].Items = lines
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

	t.Relayout()

	var widgets []termui.Bufferer
	widgets = append(widgets, t.Title, t.Status, t.Services, t.Sparkline1)
	for _, list := range t.Lists {
		widgets = append(widgets, list)
	}
	if t.Sparkline2 != nil {
		widgets = append(widgets, t.Sparkline2)
	}
	termui.Render(widgets...)
}

// Relayout recalculates widgets sizes and coords.
func (t *TermUI) Relayout() {
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

	// Second Row: lists
	num := len(t.Lists) + 1
	listW := tw / num

	// Services list must have visible names
	minNameWidth := 20
	t.Services.Width = minNameWidth
	if listW > minNameWidth {
		t.Services.Width = listW
	}

	// Recalculate widths for each list
	listW = (tw - t.Services.Width) / (num - 1)

	// Finally, enlarge services list, if there is a space left
	if listW*(num-1)+t.Services.Width < tw {
		t.Services.Width = tw - (listW * (num - 1))
	}

	t.Services.Y = th - h

	for i, list := range t.Lists {
		list.Y = th - h
		list.Width = listW
		list.Height = len(t.Lists[0].Items) + 2
		list.X = t.Services.X + t.Services.Width + i*listW
	}
	h -= t.Lists[0].Height

	// Third row: sparklines for two vars
	t.Sparkline1.Width = tw
	t.Sparkline1.Height = h
	t.Sparkline1.Y = th - h
	if t.Sparkline2 != nil {
		t.Sparkline1.Width = tw / 2

		t.Sparkline2.Width = tw / 2
		if tw%2 == 1 {
			t.Sparkline2.Width++
		}
		t.Sparkline2.X = t.Sparkline1.X + t.Sparkline1.Width
		t.Sparkline2.Height = h
		t.Sparkline2.Y = th - h
	}

}

// Close shuts down UI module.
func (t *TermUI) Close() {
	termui.Close()
}

// StatusLine returns status line for service with it's name and status.
func StatusLine(s *Service) string {
	if s.Err != nil {
		return fmt.Sprintf("[E] ⛔ %s failed", s.Name)
	}

	if s.Restarted {
		return fmt.Sprintf("[R] 🔥 %s", s.Name)
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
