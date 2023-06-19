package worker

import (
	"log"
	"sync"
	"testing"
	"time"
)

func TestNewDispatcher(t *testing.T) {
	serverName := "artifactory.domain"
	workers := 5
	workQueue := make(MockWorkQueue)
	outputQueue := make(MockOutputQueue)
	waitGroup := &sync.WaitGroup{}
	logger := log.New(&MockLogger{}, "", 0)

	dispatcherConfig := Dispatcher{
		ServerName:  serverName,
		Workers:     workers,
		WorkQueue:   workQueue,
		OutputQueue: outputQueue,
		WaitGroup:   waitGroup,
		Logger:      logger,
	}

	dispatcher := NewDispatcher(dispatcherConfig)

	if dispatcher.ServerName != serverName {
		t.Errorf("NewDispatcher() - Expected ServerName: %s, got: %s", serverName, dispatcher.ServerName)
	}
	if dispatcher.Workers != workers {
		t.Errorf("NewDispatcher() - Expected Workers: %d, got: %d", workers, dispatcher.Workers)
	}
	if dispatcher.WorkQueue != workQueue {
		t.Errorf("NewDispatcher() - Expected WorkQueue: %v, got: %v", workQueue, dispatcher.WorkQueue)
	}
	if dispatcher.OutputQueue != outputQueue {
		t.Errorf("NewDispatcher() - Expected OutputQueue: %v, got: %v", outputQueue, dispatcher.OutputQueue)
	}
	if dispatcher.WaitGroup != waitGroup {
		t.Errorf("NewDispatcher() - Expected WaitGroup: %v, got: %v", waitGroup, dispatcher.WaitGroup)
	}
	if dispatcher.Logger != logger {
		t.Errorf("NewDispatcher() - Expected Logger: %v, got: %v", logger, dispatcher.Logger)
	}
}

func TestDispatcher_Start(t *testing.T) {
	serverName := "artifactory.domain"
	workers := 5
	workQueue := make(MockWorkQueue)
	outputQueue := make(MockOutputQueue)
	waitGroup := &sync.WaitGroup{}
	logger := log.New(&MockLogger{}, "", 0)

	dispatcherConfig := Dispatcher{
		ServerName:  serverName,
		Workers:     workers,
		WorkQueue:   workQueue,
		OutputQueue: outputQueue,
		WaitGroup:   waitGroup,
		Logger:      logger,
	}

	dispatcher := NewDispatcher(dispatcherConfig)

	// Start the dispatcher
	dispatcher.Start()

	// Wait for some time to allow the goroutines to start
	time.Sleep(100 * time.Millisecond)

	defer func() {
		if r := recover(); r != nil {
			t.Error("Dispatcher.Start() - Less workers were started by dispatcher than expected")
		}
	}()

	// decrese wait group by number of workers that should be started
	for i := 0; i < workers; i++ {
		waitGroup.Done()
	}

	// if workers were present and counter is 0 after their wait groups were finished, pass the test
	waitGroup.Wait()
}

func TestDispatcher_Stop(t *testing.T) {
	serverName := "artifactory.domain"
	workers := 5
	workQueue := make(MockWorkQueue)
	outputQueue := make(MockOutputQueue)
	waitGroup := &sync.WaitGroup{}
	logger := log.New(&MockLogger{}, "", 0)

	dispatcherConfig := Dispatcher{
		ServerName:  serverName,
		Workers:     workers,
		WorkQueue:   workQueue,
		OutputQueue: outputQueue,
		WaitGroup:   waitGroup,
		Logger:      logger,
	}

	dispatcher := NewDispatcher(dispatcherConfig)

	// Start the dispatcher
	dispatcher.Start()

	// Wait for some time to allow the goroutines to start
	time.Sleep(100 * time.Millisecond)

	// Call the Stop method
	dispatcher.Stop()

	// Check if the work queue was closed
	_, ok := <-workQueue
	if ok {
		t.Error("Dispatcher.Stop() - Work queue was not closed")
	}

	// Check if all workers finished
	waitGroup.Wait()
}
