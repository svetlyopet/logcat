package worker

// Collector receives log entries and builds a work request for the workers and sends it in the WorkQueue
func Collector(line string, format LogFormat, workQueue chan WorkRequest) {
	// build the work requests for the workers
	work := WorkRequest{
		Line:      line,
		Delimiter: format.Delimiter,
		NumFields: format.NumFields,
	}

	// send the work request to the work queue to be picked up by the workers
	workQueue <- work
	return
}
