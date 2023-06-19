package writer

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

type MockLogger struct{}

func (l *MockLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func TestWriter_Start(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(dir)

	// Create a Writer instance with a mock logger
	writeQueue := make(chan string)
	doneChan := make(chan bool)
	logger := log.New(&MockLogger{}, "", 0)
	writer := NewWriter(Writer{
		Directory:  dir,
		WriteQueue: writeQueue,
		DoneChan:   doneChan,
		Logger:     logger,
	})

	// Start the writer
	go func() {
		err := writer.Start()
		if err != nil {
			t.Error("Start() returned an error:", err)
		}
	}()

	// Write a log entry to the write queue
	writeQueue <- "Log entry 1"

	// Close the write queue to trigger stopping the writer
	close(writeQueue)

	// Wait for the writer to finish
	<-doneChan

	// Verify that the log file was created and contains the expected content
	fileInfo, err := os.Stat(dir)
	if err != nil {
		t.Fatal("Failed to get directory info:", err)
	}
	if !fileInfo.IsDir() {
		t.Fatal("Expected directory, got file")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("Failed to read directory:", err)
	}
	if len(files) != 1 {
		t.Fatal("Expected 1 file, got", len(files))
	}

	filePath := dir + "/" + files[0].Name()
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal("Failed to read file:", err)
	}

	expectedContent := "Log entry 1\n"
	if string(fileContent) != expectedContent {
		t.Errorf("Unexpected file content. Expected: %q, Got: %q", expectedContent, string(fileContent))
	}
}

func TestWriter_Stop(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(dir)

	// Create a Writer instance with a mock logger
	writeQueue := make(chan string)
	doneChan := make(chan bool)
	logger := log.New(&MockLogger{}, "", 0)
	writer := NewWriter(Writer{
		Directory:  dir,
		WriteQueue: writeQueue,
		DoneChan:   doneChan,
		Logger:     logger,
	})

	// Start the writer
	go writer.Start()

	// Wait some time to have the go routine running
	time.Sleep(time.Millisecond * 100)

	// Stop the writer
	writer.Stop()

	// Verify that the done channel receives a signal
	select {
	case <-doneChan:
		// Writer stopped successfully
	case <-time.After(time.Second):
		t.Error("Writer did not stop within the timeout")
	}
}

func TestWriter_Write(t *testing.T) {
	// Create a temporary directory for testing
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal("Failed to create temporary directory:", err)
	}
	defer os.RemoveAll(dir)

	// Create a Writer instance with a mock logger
	writeQueue := make(chan string)
	doneChan := make(chan bool)
	logger := log.New(&MockLogger{}, "", log.LstdFlags)
	writer := NewWriter(Writer{
		Directory:  dir,
		WriteQueue: writeQueue,
		DoneChan:   doneChan,
		Logger:     logger,
	})

	// Start the writer
	go func() {
		err := writer.Start()
		if err != nil {
			t.Error("Start() returned an error:", err)
		}
	}()

	// Write a log entry to the write queue
	writeQueue <- "Log entry 1"

	// Close the write queue to trigger stopping the writer
	close(writeQueue)

	// Wait for the writer to finish
	<-doneChan

	// Verify that the log file was created and contains the expected content
	fileInfo, err := os.Stat(dir)
	if err != nil {
		t.Fatal("Failed to get directory info:", err)
	}
	if !fileInfo.IsDir() {
		t.Fatal("Expected directory, got file")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal("Failed to read directory:", err)
	}
	if len(files) != 1 {
		t.Fatal("Expected 1 file, got", len(files))
	}

	filePath := dir + "/" + files[0].Name()
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal("Failed to read file:", err)
	}

	expectedContent := "Log entry 1\n"
	if string(fileContent) != expectedContent {
		t.Errorf("Unexpected file content. Expected: %q, Got: %q", expectedContent, string(fileContent))
	}
}

func TestMockLogger_Write(t *testing.T) {
	logger := MockLogger{}

	// Call the Write method
	logger.Write([]byte("test message"))
}
