package writer

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

// Writer describes a writer
type Writer struct {
	File        *os.File
	Directory   string
	Flag        int
	Permissions os.FileMode
	WriteQueue  chan string
	DoneChan    chan bool
	Logger      *log.Logger
}

// NewWriter creates and returns a new Writer object
func NewWriter(w Writer) Writer {
	writer := Writer{
		Directory:   w.Directory,
		Flag:        os.O_APPEND | os.O_CREATE | os.O_WRONLY,
		Permissions: 0644,
		WriteQueue:  w.WriteQueue,
		DoneChan:    w.DoneChan,
		Logger:      w.Logger,
	}

	return writer
}

// Start starts a writer
func (w *Writer) Start() error {
	// create initial output file
	if err := w.Open(w.Directory); err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}

	// create channel for sending ticks to rotate output file
	tick := make(chan bool, 1)

	// create cron and set to send ticks every hour
	c := cron.New()
	_, err := c.AddFunc("0 * * * *", func() { tick <- true })
	if err != nil {
		return fmt.Errorf("failed to start cron timer: %v", err)
	}
	c.Start()

	go func() {
		for {
			select {
			// listen for incoming log entries from the workers
			case logEntry, ok := <-w.WriteQueue:
				if !ok {
					if err = w.File.Close(); err != nil {
						w.Logger.Fatalf("failed to close output file: %v : %v", w.File.Name(), err)
					}
					w.Logger.Printf("stopping the writer")
					w.DoneChan <- true
					return
				}
				if err = w.Write(logEntry); err != nil {
					w.Logger.Printf("failed writing to file: %v", err)
				}

			// listen for ticks to rotate output file
			case <-tick:
				if err = w.Close(); err != nil {
					w.Logger.Fatalf("failed to close output file: %v : %v", w.File.Name(), err)
				}
				if err = w.Open(w.Directory); err != nil {
					w.Logger.Fatalf("failed to open output file: %v : %v", w.File.Name(), err)
				}
			}
		}
	}()
	return nil
}

// Close closes the file descriptor of the current file
func (w *Writer) Close() error {
	if _, err := w.File.Stat(); err != nil {
		return err
	}
	if w.File != nil {
		if err := w.File.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Open creates a new file with unique name in the dir path for writing the output
func (w *Writer) Open(dir string) error {
	t := time.Now()
	timestamp := t.Format("2006-01-02")

	random := GenerateRandomString(8)

	filename := "artifactory-traffic-" + timestamp + "-" + random + ".log"
	// looping 10 times should be sufficient to get a unique string from random func to have as file name
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(dir + filename); err != nil {
			random = GenerateRandomString(8)
			filename = "artifactory-traffic-" + timestamp + "-" + random + ".log"
		} else {
			break
		}
	}

	var err error

	w.File, err = os.OpenFile(dir+filename, w.Flag, w.Permissions)
	if err != nil {
		return err
	}
	return nil
}

// Write writes input strings to the last open file
// When the last open file does not exist, a new one is created
func (w *Writer) Write(line string) error {
	// check if the output file created by the writer exists
	// create a new one if its missing
	if _, err := os.Stat(w.File.Name()); err != nil {
		if err = w.Open(w.Directory); err != nil {
			return err
		}
	}

	// if File is set write to it
	if w.File != nil {
		if _, err := w.File.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return nil
}

// Stop closes the write queue of the writer which triggers a graceful stop
func (w *Writer) Stop() {
	// close the work queue
	close(w.WriteQueue)
}
