package retry

import (
	"errors"
	"fmt"
	"time"

	"github.com/minc-org/minc/pkg/log"
)

// Retry retries a function `fn` up to `maxRetries` times with exponential backoff.
func Retry(fn func() error, maxRetries int, initialDelay time.Duration) error {
	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil // Success, exit retry loop
		}

		log.Warn(fmt.Sprintf("Attempt: %d", attempt), "failed:", err)

		if attempt < maxRetries {
			sleepDuration := initialDelay * time.Duration(attempt) // Exponential backoff
			log.Info("Retrying in", "duration", sleepDuration)
			time.Sleep(sleepDuration)
		} else {
			return errors.New("operation failed after multiple attempts: " + err.Error())
		}
	}
	return nil
}
