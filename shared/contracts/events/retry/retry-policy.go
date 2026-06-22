/*
# PHASE 6 — ENTERPRISE RETRY POLICY

---

# STEP 6 — CREATE RETRY GOVERNANCE


shared/contracts/events/retry/retry-policy.go
*/
package retry

import "time"

type RetryPolicy struct {
	MaxRetries      int
	InitialBackoff  time.Duration
	MaxBackoff      time.Duration
	Multiplier      float64
}

var DefaultRetryPolicy = RetryPolicy{
	MaxRetries:     5,
	InitialBackoff: 2 * time.Second,
	MaxBackoff:     2 * time.Minute,
	Multiplier:     2,
}
