/*
# PHASE 7 — DEAD LETTER GOVERNANCE

---

# STEP 7 — CREATE DLQ CONTRACTS


shared/contracts/events/dlq/dlq-envelope.go
*/
package dlq

import (
	"time"
)

type DeadLetterEvent struct {
	OriginalTopic     string         `json:"original_topic"`
	ConsumerGroup     string         `json:"consumer_group"`

	FailureReason     string         `json:"failure_reason"`

	RetryCount        int            `json:"retry_count"`

	FailedAt          time.Time      `json:"failed_at"`

	OriginalPayload   []byte         `json:"original_payload"`
}
