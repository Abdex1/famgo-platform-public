/*
# PHASE 3 — STANDARDIZE TOPIC NAMING

---

# STEP 3 — CREATE TOPIC GOVERNANCE

shared/contracts/events/topics/topics.go
*/
package topics

const (
	TopicAuthEvents            = "auth.events.v1"
	TopicRideEvents            = "ride.events.v1"
	TopicDriverEvents          = "driver.events.v1"
	TopicPaymentEvents         = "payment.events.v1"
	TopicDispatchEvents        = "dispatch.events.v1"
	TopicFraudEvents           = "fraud.events.v1"
	TopicNotificationEvents    = "notification.events.v1"
	TopicSafetyEvents          = "safety.events.v1"
	TopicAuditEvents           = "audit.events.v1"
	TopicAnalyticsEvents       = "analytics.events.v1"
)

const (
	TopicAuthDLQ               = "auth.events.dlq.v1"
	TopicRideDLQ               = "ride.events.dlq.v1"
	TopicPaymentDLQ            = "payment.events.dlq.v1"
)
