package main

import (
	"fmt"
	"github.com/pyk/byten"
)

// DummyUI is an simple console UI mockup, for testing purposes.
type DummyUI struct{}

func (*DummyUI) Init(Data) {}
func (*DummyUI) Close()    {}
func (*DummyUI) Update(data Data) {
	if data.Services == nil {
		return
	}
	for _, service := range data.Services {
		fmt.Printf("%s: ", service.Name)
		if service.Err != nil {
			fmt.Printf("ERROR: %s", service.Err)
			continue
		}

		if service.MemStats != nil {
			alloc := byten.Size(int64(service.MemStats.Alloc))
			sys := byten.Size(int64(service.MemStats.Sys))
			fmt.Printf("%s/%s ", alloc, sys)
		}

		/*
			if service.Goroutines != 0 {
				fmt.Printf("goroutines: %d", service.Goroutines)
			}
		*/
		fmt.Printf("\n")
	}
}
