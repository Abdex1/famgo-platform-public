
// STEP 11 — AUDIT TELEMETRY

 
//packages/telemetry/go/audit/audit.go
  
package audit

import (
    "context"
    "time"
)

type Event struct {
    Action      string
    ActorID     string
    ResourceID  string
    Resource    string
    Timestamp   time.Time
    Metadata    map[string]any
}

type Publisher interface {
    Publish(ctx context.Context, event Event) error
}
