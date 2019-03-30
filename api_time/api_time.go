package api_time

import "time"

// Need a time interface for mocking time in tests
type ApiTime interface {
	CurrentTime() time.Time
}

type SystemTime struct{}

func (SystemTime) CurrentTime() time.Time {
	return time.Now()
}
