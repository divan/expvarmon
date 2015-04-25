package main

// UI represents UI module
type UI interface {
	Init(Data)
	Close()
	Update(Data)
}
