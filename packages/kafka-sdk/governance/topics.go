
//# STEP 5 — DEFINE TOPIC GOVERNANCE

//packages/kafka-sdk/governance/topics.go

package governance

const (
	TopicRideCreated          = "ride.created.v1"
	TopicRideAccepted         = "ride.accepted.v1"
	TopicRideCancelled        = "ride.cancelled.v1"
	TopicDriverLocationUpdate = "driver.location.updated.v1"
	TopicPaymentCompleted     = "payment.completed.v1"
	TopicAuthLoginSucceeded   = "auth.login.succeeded.v1"
)
