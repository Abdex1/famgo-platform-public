
// STEP 15 — ML TELEMETRY

 
// packages/telemetry/go/ml/inference.go
  
package ml

import "time"

type InferenceMetric struct {
    ModelName      string
    ModelVersion   string
    Latency        time.Duration
    Success        bool
}
