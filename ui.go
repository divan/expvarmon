package main

// UI represents UI module
type UI interface {
	Init(UIData)
	Close()
	Update(UIData)
}
