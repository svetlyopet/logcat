package worker

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/hpcloud/tail"

	"github.com/svetlyopet/logcat/pkg/writer"
)

var (
	wg sync.WaitGroup
)

// Dispatcher builds the workers, distributes the work to them and initializes
// the writer, who listens on a channel where the workers send their finished work
func Dispatcher(ctx context.Context, file string, outdir string, workers int) error {
	// create a worker queue which holds all workers who are available to take work
	WorkerQueue := make(chan chan WorkRequest, workers)

	// create an output queue that workers send output to
	WriteQueue := make(chan string, 100)

	// initialize the workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		worker := NewWorker(i+1, WorkerQueue, WriteQueue, ctx, &wg)
		worker.Start()
	}

	// add work to the work queue and the workers to the worker queue
	go func() {
		for {
			select {
			case work := <-WorkQueue:
				go func() {
					worker := <-WorkerQueue
					worker <- work
				}()
			}
		}
	}()

	// create a new Writer implementation
	writerImpl := writer.NewWriter(outdir, WriteQueue, ctx)
	writerImpl.Start()

	// start tailing the input file
	t, err := tail.TailFile(file, tail.Config{
		Follow: true,
		ReOpen: true,
		Location: &tail.SeekInfo{
			Offset: 0,
			Whence: io.SeekEnd,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to tail log file: %v", err)
	}

	for {
		select {
		// send log lines from the tail channel to the collector
		case line := <-t.Lines:
			Collector(line.Text, "|", 11)
		// listen on the context channel
		case <-ctx.Done():
			// stop tailing the input file
			if err = t.Stop(); err != nil {
				return fmt.Errorf("failed to close input file %v : %v", t.Filename, err)
			}

			// wait for all workers to finish
			wg.Wait()

			// stop the writer
			writerImpl.Stop()
			return nil
		}
	}
}
