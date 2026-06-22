// 18. CREATE EVENT POLICY GOVERNANCE

// packages/event-bus/governance/policies.go

package governance

const (
	MaxRetryAttempts       = 5
	DefaultPartitions      = 12
	DefaultReplication     = 3
	RequireTracePropagation = true
	RequireIdempotency     = true
)
