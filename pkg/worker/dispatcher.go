package worker

import (
	"log"
	"sync"

	"github.com/hpcloud/tail"
)

// Dispatcher describes a dispatcher
type Dispatcher struct {
	Tail        *tail.Tail
	Workers     int
	WorkQueue   chan WorkRequest
	OutputQueue chan string
	WaitGroup   *sync.WaitGroup
	Logger      *log.Logger
}

// NewDispatcher creates and returns a Dispatcher object
func NewDispatcher(d Dispatcher) *Dispatcher {
	dispatcher := &Dispatcher{
		Tail:        d.Tail,
		Workers:     d.Workers,
		WorkQueue:   d.WorkQueue,
		OutputQueue: d.OutputQueue,
		WaitGroup:   d.WaitGroup,
		Logger:      d.Logger,
	}
	return dispatcher
}

// Start starts the workers, dispatches the work to them and initializes
// the writer, who listens on a channel where the workers send their finished work
func (d *Dispatcher) Start() {
	go func() {
		// start the workers
		for i := 0; i < d.Workers; i++ {
			d.WaitGroup.Add(1)
			worker := NewWorker(i+1, d.WorkQueue, d.OutputQueue, d.WaitGroup, d.Logger)
			worker.Start()
		}
	}()
}

// Stop closes the work channels and triggers the workers to stop gracefully
func (d *Dispatcher) Stop() {
	// close the work queue
	close(d.WorkQueue)

	// wait for all workers to finish
	d.WaitGroup.Wait()
}
