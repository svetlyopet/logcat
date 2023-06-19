package writer

import (
	"math/rand"
	"testing"
	"time"
)

func TestGenerateRandomString(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	// Test case 1: Generate a random string with length 10
	length := 10
	randomString := GenerateRandomString(length)
	if len(randomString) != length {
		t.Errorf("Expected length: %d, Got: %d", length, len(randomString))
	}

	// Test case 2: Generate a random string with length 20
	length = 20
	randomString = GenerateRandomString(length)
	if len(randomString) != length {
		t.Errorf("Expected length: %d, Got: %d", length, len(randomString))
	}

	// Test case 3: Generate a random string with length 0
	length = 0
	randomString = GenerateRandomString(length)
	if len(randomString) != length {
		t.Errorf("Expected length: %d, Got: %d", length, len(randomString))
	}

	// Test case 4: Generate a random string with negative length
	length = -5
	randomString = GenerateRandomString(length)
	if randomString != "" {
		t.Errorf("Expected empty string, Got: %s", randomString)
	}
}
