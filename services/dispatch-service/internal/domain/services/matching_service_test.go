// services/dispatch-service/internal/domain/services/matching_service_test.go
package services

import (
	"testing"
)

func TestMatchDrivers(t *testing.T) {
	service := NewMatchingService(0.40, 0.30, 0.20, 0.10, 50.0, 3.5, 25.0)

	drivers := []*DriverInfo{
		{
			DriverID:       "driver1",
			Latitude:       9.0320,
			Longitude:      38.7469,
			IsOnline:       true,
			AcceptanceRate: 95.0,
			Rating:         4.8,
			Distance:       2.5,
			ETA:            8.0,
		},
		{
			DriverID:       "driver2",
			Latitude:       9.0350,
			Longitude:      38.7500,
			IsOnline:       true,
			AcceptanceRate: 88.0,
			Rating:         4.2,
			Distance:       5.0,
			ETA:            15.0,
		},
		{
			DriverID:       "driver3",
			Latitude:       9.0200,
			Longitude:      38.7400,
			IsOnline:       false,
			AcceptanceRate: 92.0,
			Rating:         4.9,
			Distance:       1.5,
			ETA:            5.0,
		},
	}

	scores, err := service.MatchDrivers(9.0320, 38.7469, drivers, 5)

	if err != nil {
		t.Errorf("MatchDrivers() error = %v", err)
	}

	if len(scores) != 2 {
		t.Errorf("MatchDrivers() returned %d drivers, want 2 (offline driver filtered)", len(scores))
	}

	if scores[0].DriverID != "driver1" {
		t.Errorf("Top driver = %s, want driver1", scores[0].DriverID)
	}

	if scores[0].TotalScore < 50 || scores[0].TotalScore > 100 {
		t.Errorf("Score out of range: %.2f", scores[0].TotalScore)
	}
}

func TestCalculateProximityScore(t *testing.T) {
	service := NewMatchingService(0.40, 0.30, 0.20, 0.10, 50.0, 3.5, 25.0)

	tests := []struct {
		distance float64
		want     float64
	}{
		{0, 100.0},
		{12.5, 50.0},
		{25.0, 0.0},
		{30.0, 0.0},
	}

	for _, tt := range tests {
		t.Run("proximity", func(t *testing.T) {
			got := service.calculateProximityScore(tt.distance)
			if got != tt.want {
				t.Errorf("calculateProximityScore(%f) = %f, want %f", tt.distance, got, tt.want)
			}
		})
	}
}

func TestValidateDriversForMatching(t *testing.T) {
	service := NewMatchingService(0.40, 0.30, 0.20, 0.10, 50.0, 3.5, 25.0)

	tests := []struct {
		name       string
		driver     *DriverInfo
		wantValid  bool
	}{
		{
			"valid driver",
			&DriverInfo{IsOnline: true, AcceptanceRate: 95.0, Rating: 4.8, Distance: 5.0},
			true,
		},
		{
			"offline driver",
			&DriverInfo{IsOnline: false, AcceptanceRate: 95.0, Rating: 4.8, Distance: 5.0},
			false,
		},
		{
			"low acceptance rate",
			&DriverInfo{IsOnline: true, AcceptanceRate: 40.0, Rating: 4.8, Distance: 5.0},
			false,
		},
		{
			"low rating",
			&DriverInfo{IsOnline: true, AcceptanceRate: 95.0, Rating: 3.0, Distance: 5.0},
			false,
		},
		{
			"too far",
			&DriverInfo{IsOnline: true, AcceptanceRate: 95.0, Rating: 4.8, Distance: 30.0},
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := service.ValidateDriversForMatching(tt.driver)
			if valid != tt.wantValid {
				t.Errorf("ValidateDriversForMatching() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}
