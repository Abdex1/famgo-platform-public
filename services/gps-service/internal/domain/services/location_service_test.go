// services/gps-service/internal/domain/services/location_service_test.go
// Tests for location service

package services

import (
	"testing"

	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/gps-service/internal/domain/valueobjects"
)

func TestFindNearbyDrivers(t *testing.T) {
	service := NewLocationService()

	referenceLocation, _ := valueobjects.NewGeolocation(0, 0, 0, 0, 0, 0, 1000)

	// Create test drivers
	nearby, _ := valueobjects.NewGeolocation(0, 0.05, 0, 0, 0, 0, 1000)
	driver1, _ := entities.NewDriverLocation("driver1", nearby)

	far, _ := valueobjects.NewGeolocation(2, 2, 0, 0, 0, 0, 1000)
	driver2, _ := entities.NewDriverLocation("driver2", far)

	drivers := []*entities.DriverLocation{driver1, driver2}

	results := service.FindNearbyDrivers(referenceLocation, drivers, 100, 10, 40)

	if len(results) != 1 {
		t.Errorf("FindNearbyDrivers() returned %d drivers, want 1", len(results))
	}

	if results[0].Driver.DriverID != "driver1" {
		t.Errorf("FindNearbyDrivers() returned wrong driver: %s", results[0].Driver.DriverID)
	}
}

func TestCalculateRoute(t *testing.T) {
	service := NewLocationService()

	from, _ := valueobjects.NewGeolocation(0, 0, 0, 0, 0, 0, 1000)
	to, _ := valueobjects.NewGeolocation(0, 1, 0, 0, 0, 0, 1000)

	route := service.CalculateRoute(from, to, 60)

	if route == nil {
		t.Fatalf("CalculateRoute() returned nil")
	}

	if route.DistanceKm <= 0 {
		t.Errorf("CalculateRoute() distance = %f, want > 0", route.DistanceKm)
	}

	if route.ETAMinutes <= 0 {
		t.Errorf("CalculateRoute() ETA = %f, want > 0", route.ETAMinutes)
	}
}

func TestAnalyzeTrajectory(t *testing.T) {
	service := NewLocationService()

	from, _ := valueobjects.NewGeolocation(0, 0, 0, 0, 0, 0, 0)
	to, _ := valueobjects.NewGeolocation(0, 1, 0, 0, 0, 0, 3600000) // 1 hour later

	trajectory := service.AnalyzeTrajectory(from, to, 150)

	if !trajectory.IsValid {
		t.Errorf("AnalyzeTrajectory() IsValid = false, want true")
	}

	if trajectory.IsAnomaly {
		t.Errorf("AnalyzeTrajectory() IsAnomaly = true, want false")
	}
}

func TestClusterNearbyDrivers(t *testing.T) {
	service := NewLocationService()

	loc1, _ := valueobjects.NewGeolocation(0, 0, 0, 0, 0, 0, 1000)
	driver1, _ := entities.NewDriverLocation("driver1", loc1)

	loc2, _ := valueobjects.NewGeolocation(0, 0.05, 0, 0, 0, 0, 1000)
	driver2, _ := entities.NewDriverLocation("driver2", loc2)

	loc3, _ := valueobjects.NewGeolocation(5, 5, 0, 0, 0, 0, 1000)
	driver3, _ := entities.NewDriverLocation("driver3", loc3)

	drivers := []*entities.DriverLocation{driver1, driver2, driver3}

	clusters := service.ClusterNearbyDrivers(drivers, 1.0)

	if len(clusters) < 2 {
		t.Errorf("ClusterNearbyDrivers() returned %d clusters, want >= 2", len(clusters))
	}
}

func TestLocationQuality(t *testing.T) {
	service := NewLocationService()

	tests := []struct {
		name     string
		accuracy float64
		maxAcc   float64
		wantOK   bool
	}{
		{"good accuracy", 10, 50, true},
		{"poor accuracy", 100, 50, false},
		{"threshold", 50, 50, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			location, _ := valueobjects.NewGeolocation(0, 0, 0, tt.accuracy, 0, 0, 1000)
			ok, err := service.UpdateLocationQuality(location, tt.maxAcc)
			if ok != tt.wantOK {
				t.Errorf("UpdateLocationQuality() = %v, want %v (err: %v)", ok, tt.wantOK, err)
			}
		})
	}
}

func TestCalculateGeohash(t *testing.T) {
	service := NewLocationService()

	hash, err := service.CalculateGeohash(0, 0, 8)

	if err != nil {
		t.Errorf("CalculateGeohash() error = %v", err)
	}

	if hash == "" {
		t.Errorf("CalculateGeohash() returned empty hash")
	}

	if len(hash) != 8 {
		t.Errorf("CalculateGeohash() length = %d, want 8", len(hash))
	}
}
