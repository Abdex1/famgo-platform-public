
// STEP 12 — SECURITY TELEMETRY

 
//packages/telemetry/go/security/events.go
  
package security

const (
    EventOTPRequested    = "security.otp.requested"
    EventOTPVerified     = "security.otp.verified"
    EventOTPFailed       = "security.otp.failed"
    EventMFAEnabled      = "security.mfa.enabled"
    EventSessionCreated  = "security.session.created"
    EventTokenRotated    = "security.token.rotated"
)
