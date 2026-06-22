/*
# PHASE 4 — STANDARDIZE EVENT TYPES

---

# STEP 4 — CREATE EVENT TYPE GOVERNANCE

shared/contracts/events/catalog/event-types.go
*/
package catalog

const (

	// AUTH

	EventAuthLoginSucceeded      = "auth.login.succeeded"
	EventAuthLoginFailed         = "auth.login.failed"
	EventAuthLogoutSucceeded     = "auth.logout.succeeded"

	EventAuthTokenRefreshed      = "auth.token.refreshed"
	EventAuthSessionRevoked      = "auth.session.revoked"

	EventAuthOTPRequested        = "auth.otp.requested"
	EventAuthOTPVerified         = "auth.otp.verified"

	// RIDES

	EventRideRequested           = "ride.requested"
	EventRideAccepted            = "ride.accepted"
	EventRideCancelled           = "ride.cancelled"
	EventRideStarted             = "ride.started"
	EventRideCompleted           = "ride.completed"

	// DRIVER

	EventDriverOnline            = "driver.online"
	EventDriverOffline           = "driver.offline"
	EventDriverLocationUpdated   = "driver.location.updated"

	// PAYMENT

	EventPaymentAuthorized       = "payment.authorized"
	EventPaymentCaptured         = "payment.captured"
	EventPaymentFailed           = "payment.failed"

	// FRAUD

	EventFraudDetected           = "fraud.detected"

	// SAFETY

	EventSOSTriggered            = "safety.sos.triggered"
)
