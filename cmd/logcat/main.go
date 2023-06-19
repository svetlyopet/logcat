package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/hpcloud/tail"

	"github.com/svetlyopet/logcat/pkg/worker"
	"github.com/svetlyopet/logcat/pkg/writer"
)

var (
	// variables to store cmd args
	file   string
	outdir string

	// create work queue for the workers and write queue for the writer
	workQueue  = make(chan worker.WorkRequest, 100)
	writeQueue = make(chan string, 100)

	// create a done channel for the writer
	doneChan = make(chan bool)

	// create a wait group to track the workers
	wg sync.WaitGroup
)

// PrintHelp prints out to stdout help information about this program and exits
func PrintHelp() {
	fmt.Println("Usage: logcat -file [FILEPATH] -outdir [DIRECTORY]")
	fmt.Println("Example: logcat -file /opt/artifactory/var/log/artifactory-requests.log -outdir /tmp")
	os.Exit(1)
}

func main() {
	// parse the cli flags
	flag.StringVar(&file, "file", "", "Path to file we are parsing")
	flag.StringVar(&outdir, "outdir", "", "Directory for writing billing logs to")
	flag.Parse()

	// ensure file and outdir are absolute paths
	if !strings.HasPrefix(file, "/") || !strings.HasPrefix(file, "/") {
		fmt.Println("Path to file and directory must be absolute path")
		PrintHelp()
	}

	// create a logger
	logger := log.New(os.Stdout, "logcat: ", log.Ldate|log.Ltime)

	// use context to handle sys signals
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		<-sigCh
		cancel()
	}()

	// start reading lines from the file we are monitoring
	t, err := tail.TailFile(file, tail.Config{
		Follow: true,
		ReOpen: true,
		Logger: logger,
		Location: &tail.SeekInfo{
			Offset: io.SeekCurrent,
			Whence: io.SeekEnd,
		},
	})
	if err != nil {
		logger.Fatalf("failed to tail log file: %v", err)
	}

	// define the log format and number of fields that should be present in the log file we are reading from
	// this is used by the collector which does the sanity check for input log lines
	logFormat := worker.LogFormat{
		Delimiter: "|",
		NumFields: 11,
	}

	// create a config for the work dispatcher
	dispatcherConfig := worker.Dispatcher{
		ServerName:  "artifactory.domain",
		Workers:     5,
		WorkQueue:   workQueue,
		OutputQueue: writeQueue,
		WaitGroup:   &wg,
		Logger:      logger,
	}

	// create a work Dispatcher implementation
	dispatcherImpl := worker.NewDispatcher(dispatcherConfig)
	dispatcherImpl.Start()

	writerConfig := writer.Writer{
		Directory:   outdir,
		Flag:        os.O_CREATE | os.O_APPEND | os.O_WRONLY,
		Permissions: 0644,
		WriteQueue:  writeQueue,
		DoneChan:    doneChan,
		Logger:      logger,
	}

	// create a Writer implementation
	writerImpl := writer.NewWriter(writerConfig)
	err = writerImpl.Start()
	if err != nil {
		logger.Fatalf("failed to initialize writer: %v", err)
	}

	for {
		select {
		case line := <-t.Lines:
			// send log lines from the tail channel to the collector
			worker.Collector(line.Text, logFormat, workQueue)
		case <-ctx.Done():
			// gracefully stop everything
			if err = t.Stop(); err != nil {
				logger.Printf("failed to gracefully stop tailing input file: %v", err)
			}
			dispatcherImpl.Stop()
			writerImpl.Stop()

			// wait until a done signal is sent from the writer
			<-writerImpl.DoneChan
			logger.Printf("logcat stopped successfully")
			return
		}
	}
}
