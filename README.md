# ExpvarMon

TermUI based Go apps monitor using [expvars](http://golang.org/pkg/expvar/) variables (/debug/vars). Quickest way to monitor your Go app.

## Introduction

Go apps console monitoring tool. Minimal configuration efforts. Quick and fast monitoring solution for one or multiple services.

## Demo

### Multiple apps mode
<img src="./demo/demo_multi.png" alt="Multi mode" width="800">

### Single mode
<img src="./demo/demo_single.png" alt="Single mode" width="800">

You can monitor arbitrary number of services and variables:

<a href="./demo/demo_1var.png" target="_blank"><img src="./demo/demo_1var.png" alt="1 var" width="350"></a> <a href="./demo/demo_small.png" target="_blank"><img src="./demo/demo_small.png" alt="25 apps" width="350"></a>

## Install

Just run go install:

    go install github.com/divan/expvarmon

## Usage

### Prepare your app

First, you have to add [expvars](http://golang.org/pkg/expvar/) support into your Go program. It's as simple as:

    import _ "expvar"

and note the port your app is listening on. It it's not, just add two lines:

    import "net/http"
    ...
    http.ListenAndServe(":1234", nil)

and expvar will add handler for "localhost:1234/debug/vars" to your app.

By default, expvars adds to variables: *memstats* and *cmdline*. It's enough to monitor memory memory and garbage collector status in your app.

### Run expvarmon


Just run expvarmon with -ports="1234" parameter:

    expvarmon -ports="1234"
    
That's it.

## Advanced usage

If you need to monitor more (or less) vars, you can specify them with -vars command line flag.

    expvarmon -help
    Usage of ./expvarmon:
    -dummy=false: Use dummy (console) output
    -i=5s: Polling interval
    -ports="": Ports for accessing services expvars (start-end,port2,port3)
    -self=false: Monitor itself
    -vars="mem:memstats.Alloc,mem:memstats.Sys,mem:memstats.HeapAlloc,mem:memstats.HeapInuse,memstats.EnableGC,memstats.NumGC,duration:memstats.PauseTotalNs": Vars to monitor (comma-separated)
    Examples:
        ./expvarmon -ports="80"
        ./expvarmon -ports="23000-23010,80" -i=1m
        ./expvarmon -ports="80,remoteapp:80" -vars="mem:memstats.Alloc,duration:Response.Mean,Counter"
        ./expvarmon -ports="1234-1236" -vars="Goroutines" -self
    For more details and docs, see README: http://github.com/divan/expvarmon

So, yes, you can specify multiple ports, using '-' for ranges, and specify host(s) for remote apps.

You can also monitor expvarmon itself, using -self flag.
