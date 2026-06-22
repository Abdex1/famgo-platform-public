// services/ride-service/internal/domain/services/ride_service_test.go
package services

import (
	"testing"
	"time"
)

func TestFareCalculation(t *testing.T) {
	service := NewRideService()
	
	// 10 km, 20 minutes
	fare := service.CalculateFare(10.0, 20*time.Minute, 1.0, 0)
	
	if fare == nil {
		t.Fatalf("CalculateFare() returned nil")
	}
	
	if fare.TotalFare < 100 {
		t.Errorf("TotalFare = %f, expected >= 100", fare.TotalFare)
	}
}

func TestFareWithSurge(t *testing.T) {
	service := NewRideService()
	
	normalFare := service.CalculateFare(10.0, 20*time.Minute, 1.0, 0)
	surgeFare := service.CalculateFare(10.0, 20*time.Minute, 1.5, 0)
	
	if surgeFare.TotalFare <= normalFare.TotalFare {
		t.Errorf("Surge fare should be higher than normal fare")
	}
}

func TestValidateRideRequest(t *testing.T) {
	service := NewRideService()
	
	tests := []struct {
		name    string
		pickupLat float64
		pickupLng float64
		dropoffLat float64
		dropoffLng float64
		minDist  float64
		maxDist  float64
		wantValid bool
	}{
		{"valid", 9.0320, 38.7469, 9.0350, 38.7500, 0.1, 100, true},
		{"invalid pickup lat", 91, 38.7469, 9.0350, 38.7500, 0.1, 100, false},
		{"too short", 9.0320, 38.7469, 9.0321, 38.7470, 1, 100, false},
		{"too long", 0, 0, 45, 45, 0.1, 100, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, _ := service.ValidateRideRequest(
				tt.pickupLat, tt.pickupLng,
				tt.dropoffLat, tt.dropoffLng,
				tt.minDist, tt.maxDist,
			)
			if valid != tt.wantValid {
				t.Errorf("ValidateRideRequest() = %v, want %v", valid, tt.wantValid)
			}
		})
	}
}

func TestWaitTimeCalculation(t *testing.T) {
	service := NewRideService()
	
	// 5 km away, 60 km/h speed = 5 minutes
	waitTime := service.CalculateWaitTime(5, 60)
	expectedMinutes := 5.0
	
	if waitTime.Minutes() < expectedMinutes-1 || waitTime.Minutes() > expectedMinutes+1 {
		t.Errorf("WaitTime = %f minutes, want ~%f", waitTime.Minutes(), expectedMinutes)
	}
}
