
// STEP 8 — CREATE RETRY ENGINE

// packages/kafka-sdk/consumer/retry_handler.go

package consumer

import (
	"time"

	"github.com/cenkalti/backoff/v4"
)

func NewRetryPolicy() backoff.BackOff {
	exponential := backoff.NewExponentialBackOff()

	exponential.InitialInterval = 500 * time.Millisecond
	exponential.MaxInterval = 30 * time.Second
	exponential.MaxElapsedTime = 5 * time.Minute
	return exponential
}
