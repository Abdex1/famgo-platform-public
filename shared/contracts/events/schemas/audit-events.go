/*
# PHASE 11 — ENTERPRISE AUDIT PIPELINE

---

# STEP 11 — CREATE AUDIT CONTRACTS

shared/contracts/events/schemas/audit-events.go
*/
package schemas

import "time"

type AuditEvent struct {
	Action         string    `json:"action"`
	ActorID        string    `json:"actor_id"`

	ResourceType   string    `json:"resource_type"`
	ResourceID     string    `json:"resource_id"`

	IPAddress      string    `json:"ip_address"`

	Timestamp      time.Time `json:"timestamp"`
}
