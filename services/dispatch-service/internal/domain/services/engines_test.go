package services

import (
	"context"
	"testing"

	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/entities"
	"github.com/Abdex1/FamGo-platform/services/dispatch-service/internal/domain/ports"
)

type fakeDiscovery struct {
	drivers []ports.DriverCandidate
}

func (f fakeDiscovery) FindDriversWithinRadius(
	_ context.Context,
	_, _, _ float64,
	_ int,
) ([]ports.DriverCandidate, error) {
	return f.drivers, nil
}

func TestDriverCandidateEngine_FindWithinRadius(t *testing.T) {
	engine := NewDriverCandidateEngine(fakeDiscovery{
		drivers: []ports.DriverCandidate{
			{DriverID: "d1", DistanceKm: 1.2, IsOnline: true, AcceptanceRate: 90, Rating: 4.8},
			{DriverID: "d2", DistanceKm: 0.8, IsOnline: false, AcceptanceRate: 90, Rating: 4.9},
			{DriverID: "d3", DistanceKm: 2.0, IsOnline: true, AcceptanceRate: 80, Rating: 4.5},
		},
	})

	candidates, err := engine.FindWithinRadius(context.Background(), 9.03, 38.74, 5, 10)
	if err != nil {
		t.Fatalf("FindWithinRadius() error = %v", err)
	}
	if len(candidates) != 2 {
		t.Fatalf("expected 2 online drivers, got %d", len(candidates))
	}

	sorted := engine.SortByDistance(candidates)
	if sorted[0].DriverID != "d1" {
		t.Fatalf("expected nearest online driver d1, got %s", sorted[0].DriverID)
	}
}

func TestMatchingEngine_Run(t *testing.T) {
	matchingService := NewMatchingService(0.4, 0.3, 0.2, 0.1, 50, 3.5, 25)
	engine := NewMatchingEngine(
		NewDriverCandidateEngine(fakeDiscovery{
			drivers: []ports.DriverCandidate{
				{DriverID: "near", DistanceKm: 1.0, ETAMinutes: 4, IsOnline: true, AcceptanceRate: 95, Rating: 4.9, Latitude: 9.031, Longitude: 38.747},
				{DriverID: "far", DistanceKm: 4.0, ETAMinutes: 12, IsOnline: true, AcceptanceRate: 92, Rating: 4.7, Latitude: 9.04, Longitude: 38.75},
			},
		}),
		matchingService,
		3,
	)

	request, err := entities.NewDispatchRequest("ride-1", "rider-1", 9.03, 38.7469, 9.05, 38.76, 5, 25, 3)
	if err != nil {
		t.Fatalf("NewDispatchRequest() error = %v", err)
	}

	result, err := engine.Run(context.Background(), request)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.PrimaryDriverID != "near" {
		t.Fatalf("expected near driver, got %s", result.PrimaryDriverID)
	}
	if len(result.ProposedDrivers) == 0 {
		t.Fatal("expected proposed drivers")
	}
}

func TestAssignmentEngine_AssignFirstAvailable(t *testing.T) {
	engine := NewAssignmentEngine(60, ports.NoOpPoolingHook{})
	request, _ := entities.NewDispatchRequest("ride-2", "rider-2", 9.03, 38.74, 9.05, 38.76, 5, 25, 3)

	matchResult := &MatchingEngineResult{
		PrimaryDriverID: "driver-a",
		ProposedDrivers: []string{"driver-a", "driver-b"},
	}

	assignment, err := engine.AssignFirstAvailable(context.Background(), request, matchResult)
	if err != nil {
		t.Fatalf("AssignFirstAvailable() error = %v", err)
	}
	if assignment.AssignedDriver != "driver-a" {
		t.Fatalf("expected driver-a, got %s", assignment.AssignedDriver)
	}
	if request.Status != entities.StatusMatched {
		t.Fatalf("expected matched status, got %s", request.Status)
	}
}

func TestAssignmentEngine_Reassign(t *testing.T) {
	engine := NewAssignmentEngine(60, ports.NoOpPoolingHook{})
	request, _ := entities.NewDispatchRequest("ride-3", "rider-3", 9.03, 38.74, 9.05, 38.76, 5, 25, 3)
	_ = request.StartMatching()
	_ = request.Match("driver-a", []string{"driver-a", "driver-b", "driver-c"})
	_ = request.Reject("driver declined")

	next, err := engine.Reassign(request, "driver-a")
	if err != nil {
		t.Fatalf("Reassign() error = %v", err)
	}
	if next != "driver-b" {
		t.Fatalf("expected driver-b, got %s", next)
	}
}
