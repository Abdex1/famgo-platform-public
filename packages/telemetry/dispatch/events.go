
// STEP 14 — DISPATCH TELEMETRY

 
//packages/telemetry/go/dispatch/events.go
  
package dispatch

const (
    MatchStarted  = "dispatch.match.started"
    MatchFound    = "dispatch.match.found"
    MatchTimeout  = "dispatch.match.timeout"
    SurgeApplied  = "dispatch.surge.applied"
)
