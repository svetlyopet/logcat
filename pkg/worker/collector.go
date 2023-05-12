package worker

import (
	"log"
	"strings"
)

var (
	WorkQueue = make(chan WorkRequest, 100)
)

// The Collector receives log entries and builds a work request for the workers
func Collector(line string, delimiter string, num int) {
	split := strings.Split(line, delimiter)
	if len(split) != num {
		log.Printf("missmatch number of fields for line: %v : expected number of fields: %d, found %d\n", line, num, len(split))
		return
	}

	work := WorkRequest{
		Line: line,
	}

	WorkQueue <- work
	return
}
