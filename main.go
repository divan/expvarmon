package main

import (
	"flag"
	"log"
	"time"

	"github.com/divan/termui"
)

var (
	interval  = flag.Duration("i", 1*time.Second, "Polling interval")
	portsArg  = flag.String("ports", "40001,40002,40000,40004,1233,1234,1235", "Ports for accessing services expvars")
	extraVars = flag.String("extravars", "Goroutines,Counters.A", "Extra vars exported with expvars package")
	dummy     = flag.Bool("dummy", false, "Use dummy (console) output")
)

func main() {
	flag.Parse()
	ports, err := ParsePorts(*portsArg)
	if err != nil {
		log.Fatal("cannot parse ports:", err)
	}

	data := *NewData()
	var source = NewExpvarsSource(ports)
	for _, port := range ports {
		service := NewService(port)
		data.Services = append(data.Services, service)
	}
	data.Total = len(data.Services)

	var ui UI = &TermUI{}
	if *dummy {
		ui = &DummyUI{}
	}
	ui.Init(data)
	defer ui.Close()

	tick := time.NewTicker(*interval)
	evtCh := termui.EventCh()

	update := func() {
		for _, port := range source.Ports {
			service := data.FindService(port)
			if service == nil {
				continue
			}

			service.Update()
		}

		data.LastTimestamp = time.Now()

		ui.Update(data)
	}
	update()
	for {
		select {
		case <-tick.C:
			update()
		case e := <-evtCh:
			if e.Type == termui.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == termui.EventResize {
				termui.Body.Width = termui.TermWidth()
				termui.Body.Align()
			}
		}
	}
}
