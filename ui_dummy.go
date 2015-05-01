package main

import (
	"fmt"
)

// DummyUI is an simple console UI mockup, for testing purposes.
type DummyUI struct{}

func (*DummyUI) Init(UIData) error { return nil }
func (*DummyUI) Close()            {}
func (*DummyUI) Update(data UIData) {
	if data.Services == nil {
		return
	}
	for _, service := range data.Services {
		fmt.Printf("%s: ", service.Name)
		if service.Err != nil {
			fmt.Printf("ERROR: %s", service.Err)
			continue
		}

		/*
			if service.Goroutines != 0 {
				fmt.Printf("goroutines: %d", service.Goroutines)
			}
		*/
		fmt.Printf("\n")
	}
}
