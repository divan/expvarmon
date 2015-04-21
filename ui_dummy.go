package main

import (
	"fmt"
	"github.com/pyk/byten"
)

// DummyUI is an simple console UI mockup, for testing purposes.
type DummyUI struct{}

func (u *DummyUI) Init()  {}
func (u *DummyUI) Close() {}
func (u *DummyUI) Update(data Data) {
	if data.Services == nil {
		return
	}
	for _, service := range data.Services {
		fmt.Printf("%s: ", service.Name)
		if service.Err != nil {
			fmt.Printf("ERROR: %s", service.Err)
			continue
		}

		if service.Memstats != nil {
			alloc := byten.Size(int64(service.Memstats.Alloc))
			sys := byten.Size(int64(service.Memstats.Sys))
			fmt.Printf("%s/%s ", alloc, sys)
		}

		if service.Goroutines != 0 {
			fmt.Printf("goroutines: %d", service.Goroutines)
		}
		fmt.Printf("\n")
	}
}
