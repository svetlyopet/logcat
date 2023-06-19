package worker

import (
	"log"
	"sync"
	"testing"
)

type MockWorkQueue chan WorkRequest
type MockOutputQueue chan string
type MockLogger struct{}

func (l *MockLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestNewWorker(t *testing.T) {
	id := 1
	serverName := "artifactory.domain"
	workQueue := make(MockWorkQueue)
	outputQueue := make(MockOutputQueue)
	waitGroup := &sync.WaitGroup{}
	logger := log.New(&MockLogger{}, "", 0)

	// Create a new worker
	worker := NewWorker(id, serverName, workQueue, outputQueue, waitGroup, logger)

	// Check if the worker is created correctly
	if worker.ID != id {
		t.Errorf("NewWorker() - Expected ID: %d, got: %d", id, worker.ID)
	}
	if worker.ServerName != serverName {
		t.Errorf("NewWorker() - Expected ServerName: %s, got: %s", serverName, worker.ServerName)
	}
	if worker.WorkQueue != workQueue {
		t.Error("NewWorker() - WorkQueue is not set correctly")
	}
	if worker.OutputQueue != outputQueue {
		t.Error("NewWorker() - OutputQueue is not set correctly")
	}
	if worker.WaitGroup != waitGroup {
		t.Error("NewWorker() - WaitGroup is not set correctly")
	}
	if worker.Logger != logger {
		t.Error("NewWorker() - Logger is not set correctly")
	}
}

func TestWorker_Start(t *testing.T) {
	id := 1
	serverName := "artifactory.domain"
	workQueue := make(MockWorkQueue)
	outputQueue := make(MockOutputQueue)
	waitGroup := &sync.WaitGroup{}
	logger := log.New(&MockLogger{}, "", 0)

	// Create a new worker
	worker := NewWorker(id, serverName, workQueue, outputQueue, waitGroup, logger)
	waitGroup.Add(1)

	// Start the worker
	worker.Start()

	// Send a work request to the work queue
	workQueue <- WorkRequest{
		Line:      "2023-06-15T12:34:56.789Z|abcdefgh12345678|1.2.3.4|user|GET|/api/docker/registry-docker-remote/v2/alpine/curl/manifests/latest|200|-1|1234|567|user-agent123",
		Delimiter: "|",
		NumFields: 11,
	}

	// Check if the work request is processed and sent to the output queue
	select {
	case output := <-outputQueue:
		// Check if the output matches the expected value
		expectedOutput := `{"billing_timestamp":"2023-06-15 12:00:00.000","server_name":"artifactory.domain","service":"artifactory","action":"download","ip":"1.2.3.4","repository":"registry-docker-remote","project":"default","artifactory_path":"alpine/curl/manifests/latest","user_name":"user","consumption_unit":"bytes","quantity":1234}`
		if output != expectedOutput {
			t.Errorf("Worker.Start() - Expected output: %s, got: %s", expectedOutput, output)
		}
	}

	// Close the work queue
	close(workQueue)

	// Wait for the worker to finish
	waitGroup.Wait()
}

func TestMockLogger_Write(t *testing.T) {
	logger := MockLogger{}

	// Call the Write method
	logger.Write([]byte("test message"))
}
