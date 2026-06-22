
// 7. CREATE ENVELOPE BUILDER

// packages/event-bus/envelope/builder.go

package envelope

import (
	"time"

	"github.com/google/uuid"
)

func NewEnvelope(
	eventType string,
	version string,
	service string,
	domain string,
	payload any,
) EventEnvelope {

	now := time.Now().UTC()

	return EventEnvelope{
		EventID:        uuid.NewString(),
		EventType:      eventType,
		EventVersion:   version,
		Service:        service,
		Domain:         domain,
		OccurredAt:     now,
		CorrelationID:  uuid.NewString(),
		IdempotencyKey: uuid.NewString(),
		Payload:        payload,
	}
}
