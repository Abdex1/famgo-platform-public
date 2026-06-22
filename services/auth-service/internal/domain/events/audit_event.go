/*
# STEP 13 — CREATE AUDIT EVENT MODEL

services/auth-service/internal/domain/events/audit_event.go


# FILE: audit_event.go
*/
package events

import "time"

type AuditEvent struct {
	ID string

	UserID string

	Action string

	IPAddress string
	UserAgent string

	Metadata map[string]any

	CreatedAt time.Time
}
