package main

import (
	"fmt"
	"time"
)

// DummyUI is an simple console UI mockup, for testing purposes.
type DummyUI struct{}

// Init implements UI.
func (*DummyUI) Init(UIData) error { return nil }

// Close implements UI.
func (*DummyUI) Close() {}

// Update implements UI.
func (*DummyUI) Update(data UIData) {
	if data.Services == nil {
		return
	}
	fmt.Println(time.Now().Format("15:04:05 02/01"))
	for _, service := range data.Services {
		fmt.Printf("%s: ", service.Name)
		if service.Err != nil {
			fmt.Printf("ERROR: %s\n", service.Err)
			continue
		}

		for _, name := range data.Vars {
			fmt.Printf("%s: %v, ", name.Short(), service.Value(name))
		}

		fmt.Printf("\n")
	}
}
