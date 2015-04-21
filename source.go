package main

type Source interface {
	Update() (interface{}, error)
}
