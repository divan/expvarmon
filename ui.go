package main

// UI represents UI module
type UI interface {
	Init(UIData) error
	Close()
	Update(UIData)
}
