// services/gps-service/internal/domain/valueobjects/geolocation_test.go
// Tests for geolocation value object

package valueobjects

import (
	"math"
	"testing"
)

func TestNewCoordinates(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lng     float64
		wantErr bool
	}{
		{"valid coords", 9.0320, 38.7469, false}, // Addis Ababa
		{"north pole", 90, 0, false},
		{"south pole", -90, 0, false},
		{"date line west", 0, -180, false},
		{"date line east", 0, 180, false},
		{"invalid lat high", 91, 0, true},
		{"invalid lat low", -91, 0, true},
		{"invalid lng high", 0, 181, true},
		{"invalid lng low", 0, -181, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			coords, err := NewCoordinates(tt.lat, tt.lng)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCoordinates() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !coords.IsValid() {
				t.Errorf("NewCoordinates() returned invalid coordinates")
			}
		})
	}
}

func TestDistanceCalculation(t *testing.T) {
	// Addis Ababa
	addisAbaba, _ := NewGeolocation(9.0320, 38.7469, 0, 0, 0, 0, 1000)
	// Dire Dawa (roughly 500 km away)
	direDawa, _ := NewGeolocation(9.6412, 41.8656, 0, 0, 0, 0, 1000)

	distance := addisAbaba.DistanceToKm(direDawa)

	// Should be approximately 300 km
	if distance < 200 || distance > 400 {
		t.Errorf("Distance calculation off: got %.2f km, expected ~300 km", distance)
	}
}

func TestBearingCalculation(t *testing.T) {
	point1, _ := NewGeolocation(0, 0, 0, 0, 0, 0, 1000)
	point2, _ := NewGeolocation(0, 1, 0, 0, 0, 0, 1000)

	bearing := point1.BearingTo(point2)

	// Should be roughly 90 degrees (East)
	if bearing < 85 || bearing > 95 {
		t.Errorf("Bearing calculation off: got %.2f degrees, expected ~90", bearing)
	}
}

func TestETACalculation(t *testing.T) {
	point1, _ := NewGeolocation(0, 0, 0, 0, 0, 0, 1000)
	point2, _ := NewGeolocation(0, 1, 0, 0, 0, 0, 1000)

	// At 60 km/h, roughly 111 km should take ~1.85 hours
	eta := point1.EstimatedArrivalTime(point2, 60)

	if eta < 60 || eta > 150 {
		t.Errorf("ETA calculation off: got %.2f minutes, expected ~110", eta)
	}
}

func TestWithinRadius(t *testing.T) {
	center, _ := NewGeolocation(0, 0, 0, 0, 0, 0, 1000)
	nearby, _ := NewGeolocation(0, 0.5, 0, 0, 0, 0, 1000)
	far, _ := NewGeolocation(1, 1, 0, 0, 0, 0, 1000)

	tests := []struct {
		name     string
		location *Geolocation
		radius   float64
		want     bool
	}{
		{"within 100 km", nearby, 100, true},
		{"outside 10 km", nearby, 10, false},
		{"far location", far, 50, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := center.IsWithinRadius(tt.location, tt.radius)
			if got != tt.want {
				t.Errorf("IsWithinRadius() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInterpolation(t *testing.T) {
	point1, _ := NewGeolocation(0, 0, 0, 0, 0, 0, 0)
	point2, _ := NewGeolocation(0, 2, 0, 0, 0, 0, 1000)

	// Midpoint
	mid := point1.Interpolate(point2, 0.5)

	if mid == nil {
		t.Fatalf("Interpolate() returned nil")
	}

	// Longitude should be approximately 1
	if mid.Coordinates.Longitude < 0.9 || mid.Coordinates.Longitude > 1.1 {
		t.Errorf("Interpolation longitude off: got %.2f, want ~1", mid.Coordinates.Longitude)
	}
}

func TestGeolocationValidation(t *testing.T) {
	tests := []struct {
		name    string
		lat     float64
		lng     float64
		alt     float64
		acc     float64
		spd     float64
		heading float64
		ts      int64
		want    bool
	}{
		{"valid", 0, 0, 0, 10, 0, 0, 1000, true},
		{"negative accuracy", 0, 0, 0, -1, 0, 0, 1000, false},
		{"negative speed", 0, 0, 0, 10, -1, 0, 1000, false},
		{"invalid heading high", 0, 0, 0, 10, 0, 361, 1000, false},
		{"invalid heading low", 0, 0, 0, 10, 0, -1, 1000, false},
		{"zero timestamp", 0, 0, 0, 10, 0, 0, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geo, err := NewGeolocation(tt.lat, tt.lng, tt.alt, tt.acc, tt.spd, tt.heading, tt.ts)
			if err != nil {
				if tt.want {
					t.Errorf("NewGeolocation() unexpected error: %v", err)
				}
				return
			}
			if got := geo.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDegreesToRadiansConversion(t *testing.T) {
	tests := []struct {
		degrees float64
		want    float64
	}{
		{0, 0},
		{90, math.Pi / 2},
		{180, math.Pi},
		{360, 2 * math.Pi},
		{-90, -math.Pi / 2},
	}

	for _, tt := range tests {
		t.Run("conversion", func(t *testing.T) {
			got := degreesToRadians(tt.degrees)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("degreesToRadians(%f) = %f, want %f", tt.degrees, got, tt.want)
			}
		})
	}
}
