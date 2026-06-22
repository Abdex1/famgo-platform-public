/*
# PHASE 8 — ENTERPRISE AUTH EVENTS

internal/domain/events/auth_events.go
*/
package events

type LoginSucceededEvent struct {
    UserID      string `json:"user_id"`
    Email       string `json:"email"`
    DeviceID    string `json:"device_id"`
    IPAddress   string `json:"ip_address"`
}

type LoginFailedEvent struct {
    Email       string `json:"email"`
    IPAddress   string `json:"ip_address"`
}

type TokenRefreshedEvent struct {
    UserID string `json:"user_id"`
}

type LogoutEvent struct {
    UserID string `json:"user_id"`
}

