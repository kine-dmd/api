package api_time

import (
	"testing"
	"time"
)

func TestSystemTimeMatchesSystemTime(t *testing.T) {
	apiTime := SystemTime{}.CurrentTime()
	curTime := time.Now()

	print(apiTime.Sub(curTime).Nanoseconds())
	if apiTime.Sub(curTime).Nanoseconds() > 1000 {
		t.Fatal("API time and system time are more than 1 microsecond out.")
	}

}
