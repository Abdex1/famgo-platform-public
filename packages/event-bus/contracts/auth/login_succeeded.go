
// 14. CREATE CONTRACTS — AUTH DOMAIN

packages/event-bus/contracts/auth/login_succeeded.go

package auth

type LoginSucceeded struct {
	UserID       string `json:"user_id"`
	DeviceID     string `json:"device_id"`
	IPAddress    string `json:"ip_address"`
	UserAgent    string `json:"user_agent"`
	OccurredAt   string `json:"occurred_at"`
}
