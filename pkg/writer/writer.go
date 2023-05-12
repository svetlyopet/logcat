package writer

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/svetlyopet/logcat/pkg/utils"
)

// Writer describes a writer
type Writer struct {
	File        *os.File
	Directory   string
	Flag        int
	Permissions os.FileMode
	WriteQueue  chan string
	QuitChan    chan bool
	Context     context.Context
}

// NewWriter creates and returns a new Writer object
func NewWriter(dir string, writeQueue chan string, ctx context.Context) Writer {
	writer := Writer{
		Directory:   dir,
		Flag:        os.O_APPEND | os.O_CREATE | os.O_WRONLY,
		Permissions: 0644,
		WriteQueue:  writeQueue,
		QuitChan:    make(chan bool),
		Context:     ctx,
	}

	return writer
}

// Start starts a writer
func (w *Writer) Start() {
	go func() {
		// create initial output file
		if err := w.Open(w.Directory); err != nil {
			log.Fatalf("failed to open file: %v", err)
		}

		// create channel for sending ticks to rotate output file
		tick := make(chan bool, 1)

		// create cron and set to send ticks every hour
		c := cron.New()
		_, err := c.AddFunc("0 * * * *", func() { tick <- true })
		if err != nil {
			log.Fatalf("failed to start cron timer: %v", err)
		}
		c.Start()

		for {
			select {
			// listen for incoming log entries from the workers
			case logEntry := <-w.WriteQueue:
				if err = w.Write(logEntry); err != nil {
					log.Printf("failed writing to file: %v", err)
				}
			// listen for ticks to rotate output file
			case <-tick:
				if err = w.Close(); err != nil {
					log.Printf("failed to close output file: %v : %v", w.File.Name(), err)
				}
				if err = w.Open(w.Directory); err != nil {
					log.Printf("failed to open output file: %v : %v", w.File.Name(), err)
				}
			// listen for context cancel func
			case <-w.QuitChan:
				log.Println("Stopping the writer")
				// close fd on output file
				if err = w.Close(); err != nil {
					log.Printf("failed to close output file: %v : %v", w.File.Name(), err)
				}

				return
			}
		}
	}()
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

	random := utils.GenerateRandomString(8)

	filename := "artifactory-traffic-" + timestamp + "-" + random + ".log"
	// looping 10 times should be sufficient to get a unique string from random func to have as file name
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(dir + filename); err != nil {
			random = utils.GenerateRandomString(8)
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

// Stop signals the writer to stop
func (w *Writer) Stop() {
	w.QuitChan <- true

	// give some time for the writer to finish writing log entries from the queue
	time.Sleep(time.Millisecond * 100)
}
