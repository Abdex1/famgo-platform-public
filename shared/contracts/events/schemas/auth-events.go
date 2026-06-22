/*
# PHASE 10 — ENTERPRISE AUTH EVENT CONTRACTS

Now we define REAL contracts.

---

# STEP 10 — CREATE AUTH EVENTS


shared/contracts/events/schemas/auth-events.go
*/
package schemas

type LoginSucceededEvent struct {
	UserID       string `json:"user_id"`
	Email        string `json:"email"`

	IP           string `json:"ip"`
	DeviceID     string `json:"device_id"`

	SessionID    string `json:"session_id"`
}

type LoginFailedEvent struct {
	Email        string `json:"email"`
	IP           string `json:"ip"`
	Reason       string `json:"reason"`
}

type TokenRefreshedEvent struct {
	UserID       string `json:"user_id"`
	SessionID    string `json:"session_id"`
}

type SessionRevokedEvent struct {
	UserID       string `json:"user_id"`
	SessionID    string `json:"session_id"`
}
