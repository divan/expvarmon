package main

// UI represents UI module
type UI interface {
	Init()
	Close()
	Update(Data)
}
