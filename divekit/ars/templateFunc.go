package ars

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// now is the current date and time - it doesn't change between calls.
var now time.Time

// Now retuns the current date and time in the specified format - it doesn't change between calls.
func Now(format string) string {
	if now.IsZero() {
		now = time.Now()
	}
	return now.Format(format)
}

// Creation returns the current date and time in the specified format - it changes between calls.
func Creation(format string) string {
	return time.Now().Format(format)
}

// Hash creates a SHA-256 hash of the input string and returns it as a hex-encoded string
// If the input is empty, a new UUID is generated and used as input.
func Hash(input string) string {
	if input == "" {
		input = Uuid()
	}
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// Uuid generates a new uuid.
func Uuid() string {
	return uuid.New().String()
}

// autoIncrementCounter is a counter for autoincrementFunc
var autoIncrementCounter int

// Autoincrement increments the counter and returns the new value
func Autoincrement() int {
	autoIncrementCounter++
	return autoIncrementCounter
}
