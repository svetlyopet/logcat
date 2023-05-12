package worker

import (
	"context"
	"log"
	"sync"

	"github.com/svetlyopet/logcat/pkg/parser"
)

// Worker describes a worker
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	OutputQueue chan string
	Context     context.Context
	WaitGroup   *sync.WaitGroup
}

// NewWorker creates and returns a new Worker object.
func NewWorker(id int, workerQueue chan chan WorkRequest, outputQueue chan string, ctx context.Context, wg *sync.WaitGroup) Worker {
	// Create and return the worker
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		OutputQueue: outputQueue,
		Context:     ctx,
		WaitGroup:   wg,
	}

	return worker
}

// Start starts a worker
func (w *Worker) Start() {
	go func() {
		defer w.WaitGroup.Done()

		for {
			// worker adds itself to the WorkerQueue
			w.WorkerQueue <- w.Work

			select {
			// get work from the Work channel
			case work := <-w.Work:
				logEntry, err := parser.Parse(work.Line)
				if err != nil {
					log.Printf("error while parsing line: %v : %v\n", work.Line, err)
					return
				}
				if logEntry == "" {
					return
				}
				w.OutputQueue <- logEntry
			// listen for context cancel func
			case <-w.Context.Done():
				log.Printf("workerID %d stopping\n", w.ID)
				return
			}
		}
	}()
}
