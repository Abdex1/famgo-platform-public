
// STEP 13 — GEO TELEMETRY

 
//packages/telemetry/go/geo/geo.go
  
package geo

import "time"

type DriverLocationMetric struct {
    DriverID   string
    Latitude   float64
    Longitude  float64
    Speed      float64
    Timestamp  time.Time
}
