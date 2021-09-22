package utils

import (
	"fmt"
	"log"
	"time"
)

// Retry retries the given function the given number waiting the given amount of time between each retry
// Returns the lasy err if the the function keeps failing till the last attempt
func Retry(attempts int, sleep time.Duration, f func() error) (err error) {
	for i := 0; i < attempts; i++ {
		if i > 0 {
			log.Printf("attempt %d retrying after error: %v\n", i, err)
			time.Sleep(sleep)
			sleep *= 2
		}
		err = f()
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}

// DefaultRetry attemots to call the function 5 times with 5 second intervals
func DefaultRetry(f func() error) (err error) {
	return Retry(5, 5*time.Second, f)
}
