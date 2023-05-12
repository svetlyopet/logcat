package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/svetlyopet/logcat/pkg/utils"
	"github.com/svetlyopet/logcat/pkg/worker"
)

var (
	file   string
	outdir string
)

func main() {
	// use context to handle sys signals
	ctx, cancel := context.WithCancel(context.Background())

	// parse cli flags
	flag.StringVar(&file, "file", "", "Path to file we are parsing")
	flag.StringVar(&outdir, "outdir", "", "Directory for writing billing logs to")
	flag.Parse()

	// ensure file and outdir are absolute paths
	if !utils.IsAbsolutePath(file) || !utils.IsAbsolutePath(outdir) {
		fmt.Println("Path to file and directory must be absolute path")
		utils.PrintHelp()
	}

	// ensure outdir path is correct format
	outdir = utils.AppendSlashIfMissing(outdir)

	// handle system signals
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		<-sigCh
		cancel()
	}()

	// initialize the work dispatcher and start processing work
	if err := worker.Dispatcher(ctx, file, outdir, 5); err != nil {
		log.Fatalf("failed to start the work dispatcher: %v", err)
	}
}
