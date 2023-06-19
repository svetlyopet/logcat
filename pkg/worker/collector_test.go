package worker

import (
	"testing"
)

func TestCollector(t *testing.T) {
	// Create a channel for work requests
	workQueue := make(chan WorkRequest, 1)

	// Define the test input
	line := "example|log|line"
	format := LogFormat{
		Delimiter: "|",
		NumFields: 3,
	}

	// Call the Collector function
	Collector(line, format, workQueue)

	// Check if the work request was added to the work queue
	select {
	case work := <-workQueue:
		// Check if the work request matches the expected values
		if work.Line != line {
			t.Errorf("Collector() - Expected line: %s, got: %s", line, work.Line)
		}
		if work.Delimiter != format.Delimiter {
			t.Errorf("Collector() - Expected delimiter: %s, got: %s", format.Delimiter, work.Delimiter)
		}
		if work.NumFields != format.NumFields {
			t.Errorf("Collector() - Expected numFields: %d, got: %d", format.NumFields, work.NumFields)
		}
	default:
		t.Error("Collector() - Work request was not added to the work queue")
	}
}
