package util

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func GenerateULID() (ulid.ULID, error) {
	// Initialize the entropy source using the current time
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	// Generate a new ULID using the current timestamp
	timestamp := time.Now()
	ulid, err := ulid.New(ulid.Timestamp(timestamp), entropy)
	if err != nil {
		return ulid, err
	}
	return ulid, nil
}
