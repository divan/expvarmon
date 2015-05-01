package main

import (
	"flag"
	"log"
	"time"

	"github.com/divan/termui"
)

var (
	interval    = flag.Duration("i", 5*time.Second, "Polling interval")
	portsArg    = flag.String("ports", "40001,40002,40000,40004,1233,1234,1235", "Ports for accessing services expvars")
	defaultVars = flag.String("vars", "memstats.Alloc,memstats.Sys", "Default vars to monitor")
	extraVars   = flag.String("extravars", "", "Comma-separated extra vars exported with expvars")
	dummy       = flag.Bool("dummy", false, "Use dummy (console) output")
)

func main() {
	flag.Parse()
	ports, err := ParsePorts(*portsArg)
	if err != nil {
		log.Fatal("cannot parse ports:", err)
	}

	vars, err := ParseVars(*defaultVars, *extraVars)
	if err != nil {
		log.Fatal(err)
	}

	data := NewUIData(vars)
	for _, port := range ports {
		service := NewService(port, vars)
		data.Services = append(data.Services, service)
	}

	var ui UI = &TermUI{}
	if *dummy {
		ui = &DummyUI{}
	}
	if err := ui.Init(*data); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	tick := time.NewTicker(*interval)
	evtCh := termui.EventCh()

	update := func() {
		for _, port := range ports {
			service := data.FindService(port)
			if service == nil {
				continue
			}

			service.Update()
		}

		data.LastTimestamp = time.Now()

		ui.Update(*data)
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
