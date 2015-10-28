package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/gizak/termui.v1"
)

var (
	interval = flag.Duration("i", 5*time.Second, "Polling interval")
	urls     = flag.String("ports", "", "Ports/URLs for accessing services expvars (start-end,port2,port3,https://host:port)")
	varsArg  = flag.String("vars", "mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,mem:memstats.HeapInuse,duration:memstats.PauseNs,duration:memstats.PauseTotalNs", "Vars to monitor (comma-separated)")
	dummy    = flag.Bool("dummy", false, "Use dummy (console) output")
	self     = flag.Bool("self", false, "Monitor itself")
	endpoint = flag.String("endpoint", DefaultEndpoint, "URL endpoint for expvars")
)

func main() {
	flag.Usage = Usage
	flag.Parse()

	DefaultEndpoint = *endpoint

	// Process ports/urls
	ports, _ := ParsePorts(*urls)
	if *self {
		port, err := StartSelfMonitor()
		if err == nil {
			ports = append(ports, port)
		}
	}
	if len(ports) == 0 {
		fmt.Fprintln(os.Stderr, "no ports specified. Use -ports arg to specify ports of Go apps to monitor")
		Usage()
		os.Exit(1)
	}
	if *interval <= 0 {
		fmt.Fprintln(os.Stderr, "update interval is not valid. Valid examples: 5s, 1m, 1h30m")
		Usage()
		os.Exit(1)
	}

	// Process vars
	vars, err := ParseVars(*varsArg)
	if err != nil {
		log.Fatal(err)
	}

	// Init UIData
	data := NewUIData(vars)
	for _, port := range ports {
		service := NewService(port, vars)
		data.Services = append(data.Services, service)
	}

	// Start proper UI
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

// Usage reimplements flag.Usage
func Usage() {
	progname := os.Args[0]
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", progname)
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
Examples:
	%s -ports="80"
	%s -ports="23000-23010,http://example.com:80-81" -i=1m
	%s -ports="80,remoteapp:80" -vars="mem:memstats.Alloc,duration:Response.Mean,Counter"
	%s -ports="1234-1236" -vars="Goroutines" -self

For more details and docs, see README: http://github.com/divan/expvarmon
`, progname, progname, progname, progname)
}
