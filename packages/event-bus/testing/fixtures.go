// 19. CREATE CONTRACT TEST FIXTURES

// packages/event-bus/testing/fixtures.go

package testing

import (
	"time"

	"github.com/Abdex1/FamGo-platform/packages/event-bus/envelope"
)

func SampleEnvelope() envelope.EventEnvelope {
	return envelope.EventEnvelope{
		EventID:      "evt_123",
		EventType:    "ride.created.v1",
		EventVersion: "v1",
		Service:      "ride-service",
		OccurredAt:   time.Now().UTC(),
	}
}
