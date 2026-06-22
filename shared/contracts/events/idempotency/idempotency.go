/*
# STEP 8 — CREATE IDEMPOTENCY CONTRACTS


shared/contracts/events/idempotency/idempotency.go
*/
package idempotency

type EventDeduplication struct {
	IdempotencyKey string `json:"idempotency_key"`
	EventID        string `json:"event_id"`
}
