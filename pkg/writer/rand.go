package writer

import (
	"math/rand"
	"time"
)

// GenerateRandomString generates a random alphanumerical string with a requested length
func GenerateRandomString(length int) string {
	if length <= 0 {
		return ""
	}
	charset := "abcde1234567890123456789012345678901234567890"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
