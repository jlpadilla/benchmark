package generator

import "sync"

// Simulation settings
type Simulation struct {
	InsertChan chan *Record
	UpdateChan chan *Record
	DeleteChan chan string
	WG         *sync.WaitGroup
}

func NewSimulation() Simulation {
	s := Simulation{
		InsertChan: make(chan *Record, 100),
		UpdateChan: make(chan *Record, 100),
		DeleteChan: make(chan string, 100),
		WG:         &sync.WaitGroup{},
	}
	return s
}
