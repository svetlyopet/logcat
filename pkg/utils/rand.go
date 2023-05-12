package utils

import (
	"math/rand"
	"time"
)

// GenerateRandomString takes an int and generates a random string with length int
func GenerateRandomString(length int) string {
	charset := "abcde1234567890123456789012345678901234567890"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
