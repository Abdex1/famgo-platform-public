
// 8. ENTERPRISE EVENT NAMING GOVERNANCE

// packages/event-bus/governance/naming.go

package governance

const (
	EventRideCreated          = "ride.created.v1"
	EventRideAccepted         = "ride.accepted.v1"
	EventRideCancelled        = "ride.cancelled.v1"

	EventDriverLocationUpdated = "driver.location.updated.v1"

	EventPaymentCompleted     = "payment.completed.v1"

	EventAuthLoginSucceeded   = "auth.login.succeeded.v1"

	EventDispatchMatchingStarted = "dispatch.matching.started.v1"
	EventDispatchDriverMatched   = "dispatch.driver.matched.v1"
	EventDispatchDriverAssigned  = "dispatch.driver.assigned.v1"
	EventDispatchMatchingFailed  = "dispatch.matching.failed.v1"
	EventDispatchMatchingExpired = "dispatch.matching.expired.v1"
)
