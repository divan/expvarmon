# ExpvarMon

[![Build Status](https://drone.io/github.com/divan/expvarmon/status.png)](https://drone.io/github.com/divan/expvarmon/latest)

TermUI based Go apps monitor using [expvars](http://golang.org/pkg/expvar/) variables (/debug/vars). Quickest way to monitor your Go app.

## Introduction

Go apps console monitoring tool. Minimal configuration efforts. Quick and fast monitoring solution for one or multiple services.

## Features

* Single- and multi-apps mode
* Local and remote apps support
* Arbitrary number of apps and vars to monitor (from 1 to 30+, depends on size of your terminal)
* Track restarted/failed apps
* Show maximum value
* Supports: Integer, float, duration, memory, string, bool variables
* Sparkline charts for integer, duration and memory data
* Auto-resize on font-size change or window resize
* Uses amazing [TermUI](https://github.com/gizak/termui) library by [gizak](https://github.com/gizak)

## Demo

### Multiple apps mode
<img src="./demo/demo_multi.png" alt="Multi mode" width="800">

### Single app mode
<img src="./demo/demo_single.png" alt="Single mode" width="800">

You can monitor arbitrary number of services and variables:

<a href="./demo/demo_1var.png" target="_blank"><img src="./demo/demo_1var.png" alt="1 var" width="350"></a> <a href="./demo/demo_small.png" target="_blank"><img src="./demo/demo_small.png" alt="25 apps" width="350"></a>

## Purpose

This app targets debug/develop sessions when you need an instant way to monitor you app(s). It's not intended to monitor apps in production (though it may be convenient).
Also it doesn't use any storage engines and doesn't send notifications.

## Install

Just run go get:

    go get github.com/divan/expvarmon

## Usage

### Prepare your app

First, you have to add [expvars](http://golang.org/pkg/expvar/) support into your Go program. It's as simple as:

    import _ "expvar"

and note the port your app is listening on. It it's not, just add two lines:

    import "net/http"
    ...
    http.ListenAndServe(":1234", nil)

and expvar will add handler for "localhost:1234/debug/vars" to your app.

By default, expvars adds to variables: *memstats* and *cmdline*. It's enough to monitor memory and garbage collector status in your app.

### Run expvarmon


Just run expvarmon with -ports="1234" parameter:

    expvarmon -ports="1234"
    
That's it.

More examples:

    ./expvarmon -ports="80"
    ./expvarmon -ports="23000-23010,80" -i=1m
    ./expvarmon -ports="80,remoteapp:80" -vars="mem:memstats.Alloc,duration:Response.Mean,Counter"
    ./expvarmon -ports="1234-1236" -vars="Goroutines" -self

## Advanced usage

If you need to monitor more (or less) vars, you can specify them with -vars command line flag.

    $ expvarmon -help
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

### Vars

Expvarmon doesn't restrict you to monitor only memstats. You can publish you own counters and variables using expvar.Publish() method or using expvar wrappers libraries. Just pass you variables names as they appear in JSON to -var command line flag.

Notation is dot-separated, for example: **memstats.Alloc** for .MemStats.Alloc field. Quick link to runtime.MemStats documentation: http://golang.org/pkg/runtime/#MemStats

Expvar allows to export only basic types - structs, ints, floats, bools and strings. Ints are used for sparklines, and displayed as is. But you can spefify modifier to make sure it will be rendered properly.

Vars are specified as a comma-separated list of var identifiers with (optional) modifiers.

| Modifier | Description |
| --------- | ----------- |
| mem:      | renders int64 as memory string (KB, MB, etc) |
| duration: | renders int64 as time.Duration (1s, 2ms, 12h23h) |
| str:      | doesn't display sparklines chart for this value, just display as string |

## TODO

* ports autodiscovery for given hostname
* more tests coverage
* better usage of color highligting (for max values or failed apps), after relevant patches will be merged to TermUI

