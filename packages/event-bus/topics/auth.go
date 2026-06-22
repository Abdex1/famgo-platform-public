/*
# PHASE 2 — KAFKA AUTH EVENT PUBLISHING

# =========================================================

# STEP 4 — CREATE AUTH TOPICS

Create:

```txt
packages/event-bus/topics/auth.go
*/
package topics

const (
	AuthLoginSucceeded = "auth.login.succeeded.v1"
	AuthLoginFailed    = "auth.login.failed.v1"

	AuthLogout         = "auth.logout.v1"

	AuthTokenRefreshed = "auth.token.refreshed.v1"

	AuthOTPRequested   = "auth.otp.requested.v1"

	AuthOTPVerified    = "auth.otp.verified.v1"
)
