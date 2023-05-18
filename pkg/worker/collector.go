package worker

import (
	"log"
	"strings"
)

// Collector receives log entries and builds a work request for the workers and sends it in the WorkQueue
func Collector(line string, format LogFormat, workQueue chan WorkRequest) {
	// check if the input string should be processed
	split := strings.Split(line, format.Delimiter)
	if len(split) != format.NumFields {
		log.Printf("missmatch number of fields for line: %v : expected number of fields: %d, found %d\n", line, format.NumFields, len(split))
		return
	}

	// build the work requests for the workers
	work := WorkRequest{
		Line:      line,
		Delimiter: format.Delimiter,
	}

	// send the work request to the work queue to be picked up by the workers
	workQueue <- work
	return
}
