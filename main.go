package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/divan/termui"
)

var (
	interval = flag.Duration("i", 5*time.Second, "Polling interval")
	portsArg = flag.String("ports", "1234", "Ports for accessing services expvars (comma-separated)")
	varsArg  = flag.String("vars", "mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,mem:memstats.HeapInuse,memstats.EnableGC,memstats.NumGC,duration:memstats.PauseTotalNs", "Vars to monitor (comma-separated)")
	dummy    = flag.Bool("dummy", false, "Use dummy (console) output")
	bind     = flag.String("expvar", "1234", "Port to listen to be able monitor itself")
)

func main() {
	flag.Parse()
	ports, err := ParsePorts(*portsArg)
	if err != nil {
		log.Fatal("cannot parse ports:", err)
	}

	vars, err := ParseVars(*varsArg)
	if err != nil {
		log.Fatal(err)
	}

	go StartHttp(*bind)

	data := NewUIData(vars)
	for _, port := range ports {
		service := NewService(port, vars)
		data.Services = append(data.Services, service)
	}

	var ui UI
	if len(data.Services) > 1 {
		ui = &TermUI{}
	} else {
		ui = &TermUISingle{}
	}
	if *dummy {
		ui = &DummyUI{}
	}

	if err := ui.Init(*data); err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	tick := time.NewTicker(*interval)
	evtCh := termui.EventCh()

	UpdateAll(ui, data)
	for {
		select {
		case <-tick.C:
			UpdateAll(ui, data)
		case e := <-evtCh:
			if e.Type == termui.EventKey && e.Ch == 'q' {
				return
			}
			if e.Type == termui.EventResize {
				ui.Update(*data)
			}
		}
	}
}

// UpdateAll collects data from expvars and refreshes UI.
func UpdateAll(ui UI, data *UIData) {
	var wg sync.WaitGroup
	for _, service := range data.Services {
		wg.Add(1)
		go service.Update(&wg)
	}
	wg.Wait()

	data.LastTimestamp = time.Now()

	ui.Update(*data)
}
