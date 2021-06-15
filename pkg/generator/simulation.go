package generator

import "sync"

// Transaction settings
type Transaction struct {
	// Configurable fields
	BatchSize  int
	GoRoutines int
	InsertType string
	// Internal fields
	InsertChan chan *Record
	UpdateChan chan *Record
	DeleteChan chan string
	WG         *sync.WaitGroup

	// Operations
	BatchInsert     func(instance string)
	StartConnectors func()
}
