package worker

import (
	"log"
	"sync"

	"github.com/svetlyopet/logcat/pkg/parser"
)

// Worker describes a worker
type Worker struct {
	ID          int
	ServerName  string
	WorkQueue   chan WorkRequest
	OutputQueue chan string
	WaitGroup   *sync.WaitGroup
	Logger      *log.Logger
}

// NewWorker creates and returns a new Worker object.
func NewWorker(id int, serverName string, workQueue chan WorkRequest, outputQueue chan string, waitGroup *sync.WaitGroup, logger *log.Logger) Worker {
	// Create and return the worker
	worker := Worker{
		ID:          id,
		ServerName:  serverName,
		WorkQueue:   workQueue,
		OutputQueue: outputQueue,
		WaitGroup:   waitGroup,
		Logger:      logger,
	}

	return worker
}

// Start starts a worker
func (w *Worker) Start() {
	go func() {
		defer w.WaitGroup.Done()

		for {
			// get work from the Work channel until we receive signal that channel is closed
			work, ok := <-w.WorkQueue
			if !ok {
				w.Logger.Printf("stoping worker %d", w.ID)
				return
			}

			// do the work
			logEntry, err := parser.Parse(work.Line, work.Delimiter, work.NumFields, w.ServerName)
			if err != nil {
				w.Logger.Printf("error while parsing line: \"%v\" : %v\n", work.Line, err)
				continue
			}
			if logEntry == "" {
				continue
			}

			// send the finished work to the output channel
			w.OutputQueue <- logEntry
		}
	}()
}
